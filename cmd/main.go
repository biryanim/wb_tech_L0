package main

import (
	"context"
	"github.com/biryanim/wb_tech_L0/internal/client/cache/lru_cache"
	"github.com/biryanim/wb_tech_L0/internal/client/db/pg"
	"github.com/biryanim/wb_tech_L0/internal/client/db/transaction"
	"github.com/biryanim/wb_tech_L0/internal/config"
	"github.com/biryanim/wb_tech_L0/internal/config/env"
	orderRepo "github.com/biryanim/wb_tech_L0/internal/repository/order"
	kafkaConsumer "github.com/biryanim/wb_tech_L0/internal/client/kafka/consumer"
	"log"
)

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

	dbcClient, err := pg.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to initialize db client: %v", err)
	}

	cacheClient := lru_cache.New(30)

	txManager := transaction.NewTransactionManager(dbcClient.DB())
	orderRepository := orderRepo.NewRepository(dbcClient)

}
