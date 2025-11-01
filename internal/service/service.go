package service

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/model"
)

// ConsumerService defines methods for managing Kafka message consumption.
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

// OrderService defines methods for managing order operations and cache management.
type OrderService interface {
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	RestoreCache(ctx context.Context, limit int) error
}
