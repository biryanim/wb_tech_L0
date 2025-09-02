package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/api"
	"github.com/biryanim/wb_tech_L0/internal/client/cache/lru_cache"
	"github.com/biryanim/wb_tech_L0/internal/client/db/pg"
	"github.com/biryanim/wb_tech_L0/internal/client/db/transaction"
	kafkaConsumer "github.com/biryanim/wb_tech_L0/internal/client/kafka/consumer"
	"github.com/biryanim/wb_tech_L0/internal/config"
	"github.com/biryanim/wb_tech_L0/internal/config/env"
	orderRepo "github.com/biryanim/wb_tech_L0/internal/repository/order"
	orderSaverConsumer "github.com/biryanim/wb_tech_L0/internal/service/consumer/order_saver"
	"github.com/biryanim/wb_tech_L0/internal/service/order"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const cacheCap = 30

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

	consumerGroup, err := sarama.NewConsumerGroup(
		kafkaConsumerConfig.Brokers(),
		kafkaConsumerConfig.GroupID(),
		kafkaConsumerConfig.Config(),
	)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}
	consumerGroupHandler := kafkaConsumer.NewGroupHandler()
	consumer := kafkaConsumer.NewConsumer(consumerGroup, consumerGroupHandler)
	defer consumer.Close()

	dbcClient, err := pg.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to initialize db client: %v", err)
	}
	defer dbcClient.Close()

	cacheClient := lru_cache.New(cacheCap)

	txManager := transaction.NewTransactionManager(dbcClient.DB())
	orderRepository := orderRepo.NewRepository(dbcClient)
	ordSaverConsumer := orderSaverConsumer.NewService(orderRepository, consumer, txManager, cacheClient)

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

	router := gin.Default()
	router.GET("order/:order_uid", orderImpl.GetOrder)

	httpServer := &http.Server{
		Addr:    httpConfig.Address(),
		Handler: router,
	}

	go func() {
		defer wg.Done()
		err = httpServer.ListenAndServe()
		fmt.Println(err)
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
		httpServer.Shutdown(ctx)
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
