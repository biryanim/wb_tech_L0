package repository

import (
	"context"

	"github.com/biryanim/wb_tech_L0/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	CreateDelivery(ctx context.Context, orderID string, delivery *model.Delivery) (string, error)
	CreatePayment(ctx context.Context, orderID string, payment *model.Payment) (string, error)
	CreateItem(ctx context.Context, orderID string, items *model.Item) error

	GetOrder(ctx context.Context, orderID string) (*model.Order, error)
	GetDelivery(ctx context.Context, orderID string) (*model.Delivery, error)
	GetPayment(ctx context.Context, orderID string) (*model.Payment, error)
	ListItems(ctx context.Context, orderID string) ([]*model.Item, error)
}
