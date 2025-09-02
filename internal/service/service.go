package service

import (
	"context"
	"github.com/biryanim/wb_tech_L0/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderService interface {
	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
}
