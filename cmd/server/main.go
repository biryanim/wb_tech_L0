package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/biryanim/wb_tech_L0/internal/metric"
	"github.com/biryanim/wb_tech_L0/internal/middleware"
	"github.com/biryanim/wb_tech_L0/internal/tracing"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/api"
	"github.com/biryanim/wb_tech_L0/internal/client/cache/lrucache"
	"github.com/biryanim/wb_tech_L0/internal/client/db/pg"
	"github.com/biryanim/wb_tech_L0/internal/client/db/transaction"
	kafkaConsumer "github.com/biryanim/wb_tech_L0/internal/client/kafka/consumer"
	kafkaProducer "github.com/biryanim/wb_tech_L0/internal/client/kafka/producer"
	"github.com/biryanim/wb_tech_L0/internal/config"
	"github.com/biryanim/wb_tech_L0/internal/config/env"
	orderRepo "github.com/biryanim/wb_tech_L0/internal/repository/order"
	"github.com/biryanim/wb_tech_L0/internal/service"
	orderSaverConsumer "github.com/biryanim/wb_tech_L0/internal/service/consumer/ordersaver"
	"github.com/biryanim/wb_tech_L0/internal/service/order"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const cacheCap = 5

func main() {
	ctx := context.Background()

	err := config.Load("local.env")
	if err != nil {
		log.Fatal(err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load pg config: %v", err)
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Fatalf("failed to load http config: %v", err)
	}

	kafkaConsumerConfig, err := env.NewKafkaConsumerConfig()
	if err != nil {
		log.Fatalf("failed to load kafka consumer config: %v", err)
	}

	jaegerConfig, err := env.NewJaegerConfig()
	if err != nil {
		log.Fatalf("failed to load jaeger config: %v", err)
	}

	err = metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

	tracer, err := tracing.InitTracer(jaegerConfig.URL(), jaegerConfig.ServiceName())
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}

	consumerGroup, err := sarama.NewConsumerGroup(
		kafkaConsumerConfig.Brokers(),
		kafkaConsumerConfig.GroupID(),
		kafkaConsumerConfig.Config(),
	)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}
	consumerGroupHandler := kafkaConsumer.NewGroupHandler(tracer)
	consumer := kafkaConsumer.NewConsumer(consumerGroup, consumerGroupHandler)
	defer func() {
		err = consumer.Close()
		if err != nil {
			log.Fatalf("failed to close consumer: %v", err)
		}
	}()

	producer, err := newSyncProducer(kafkaConsumerConfig.Brokers())
	if err != nil {
		log.Fatalf("failed to create sync producer: %v", err)
	}
	syncProducer := kafkaProducer.NewProducer(producer, tracer)
	defer func() {
		err = syncProducer.Close()
		if err != nil {
			log.Fatalf("failed to close sync producer: %v", err)
		}
	}()

	cfg, err := pgxpool.ParseConfig(pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	dbcClient, err := pg.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize db client: %v", err)
	}
	defer func() {
		err = dbcClient.Close()
		if err != nil {
			log.Fatalf("failed to close db client: %v", err)
		}
	}()

	cacheClient := lrucache.New(cacheCap)

	txManager := transaction.NewTransactionManager(dbcClient.DB())
	orderRepository := orderRepo.NewRepository(dbcClient)
	ordSaverConsumer := orderSaverConsumer.NewService(orderRepository, consumer, txManager, cacheClient, syncProducer)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer wg.Done()
		err = ordSaverConsumer.RunConsumer(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("failed to run consumer: %s", err.Error())
		}
	}()

	orderService := order.NewService(orderRepository, txManager, cacheClient)
	orderImpl := api.NewImplementation(orderService)

	err = restoreCache(ctx, cacheCap, orderService)
	if err != nil {
		log.Fatalf("failed to restore cache: %s", err.Error())
	}

	router := gin.Default()
	router.Use(otelgin.Middleware("order_service"))
	router.Use(middleware.MetricMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("order/:order_uid", orderImpl.GetOrder)
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	httpServer := &http.Server{
		Addr:    httpConfig.Address(),
		Handler: router,
	}

	go func() {
		defer wg.Done()
		err = httpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()

	gracefulShutdown(ctx, cancel, httpServer, wg)
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, httpServer *http.Server, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-waitSignal():
		log.Println("terminating: caught signal")
	}

	if httpServer != nil {
		_ = httpServer.Shutdown(ctx)
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}

func restoreCache(ctx context.Context, cacheCap int, serv service.OrderService) error {
	err := serv.RestoreCache(ctx, cacheCap)
	if err != nil {
		return err
	}
	return nil
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
