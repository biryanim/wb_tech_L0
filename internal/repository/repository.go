package repository

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
	CreateDelivery(ctx context.Context, orderID uuid.UUID, delivery *model.Delivery) (int64, error)
	CreatePayment(ctx context.Context, orderID uuid.UUID, payment *model.Payment) (int64, error)
	CreateItem(ctx context.Context, orderID uuid.UUID, items *model.Item) error
}
