package order

import (
	"context"
	"fmt"
	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	"github.com/biryanim/wb_tech_L0/internal/client/db"
	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/biryanim/wb_tech_L0/internal/repository"
	"github.com/biryanim/wb_tech_L0/internal/service"
)

var _ service.OrderService = (*serv)(nil)

type serv struct {
	orderRepository repository.OrderRepository
	txManager       db.TxManager
	cache           cache.Client
}

func NewService(orderRepository repository.OrderRepository, txManager db.TxManager, cache cache.Client) *serv {
	return &serv{
		orderRepository: orderRepository,
		txManager:       txManager,
		cache:           cache,
	}
}

func (s *serv) GetOrder(ctx context.Context, orderID string) (*model.Order, error) {
	if cached := s.cache.Get(orderID); cached != nil {
		if ord, ok := cached.(*model.Order); ok {
			return ord, nil
		}
		return nil, fmt.Errorf("cache miss")
	}

	orderModel, err := s.orderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if orderModel == nil {
		return nil, fmt.Errorf("order not found")
	}

	order := &model.Order{
		OrderUID:          orderModel.OrderUID,
		TrackNumber:       orderModel.TrackNumber,
		Entry:             orderModel.Entry,
		Locale:            orderModel.Locale,
		InternalSignature: orderModel.InternalSignature,
		CustomerID:        orderModel.CustomerID,
		DeliveryService:   orderModel.DeliveryService,
		ShardKey:          orderModel.ShardKey,
		SmID:              orderModel.SmID,
		DateCreated:       orderModel.DateCreated,
		OofShard:          orderModel.OofShard,
	}

	deliveryModel, err := s.orderRepository.GetDelivery(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	if deliveryModel != nil {
		order.Delivery = model.Delivery{
			Name:    deliveryModel.Name,
			Phone:   deliveryModel.Phone,
			Zip:     deliveryModel.Zip,
			City:    deliveryModel.City,
			Address: deliveryModel.Address,
			Region:  deliveryModel.Region,
			Email:   deliveryModel.Email,
		}
	}

	paymentModel, err := s.orderRepository.GetPayment(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if paymentModel != nil {
		order.Payment = model.Payment{
			Transaction:  paymentModel.Transaction,
			RequestID:    paymentModel.RequestID,
			Currency:     paymentModel.Currency,
			Provider:     paymentModel.Provider,
			Amount:       paymentModel.Amount,
			PaymentDt:    paymentModel.PaymentDt,
			Bank:         paymentModel.Bank,
			DeliveryCost: paymentModel.DeliveryCost,
			GoodsTotal:   paymentModel.GoodsTotal,
			CustomFee:    paymentModel.CustomFee,
		}
	}

	items, err := s.orderRepository.ListItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	for _, itemModel := range items {
		order.Items = append(order.Items, model.Item{
			ChrtID:      itemModel.ChrtID,
			TrackNumber: itemModel.TrackNumber,
			Price:       itemModel.Price,
			Rid:         itemModel.Rid,
			Name:        itemModel.Name,
			Size:        itemModel.Size,
			Sale:        itemModel.Sale,
			TotalPrice:  itemModel.TotalPrice,
			NmID:        itemModel.NmID,
			Brand:       itemModel.Brand,
			Status:      itemModel.Status,
		})
	}

	s.cache.Set(orderID, order)

	return order, nil
}
