package order_saver

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/model"
)

func (s *service) OrderSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	order := &model.Order{}
	err := json.Unmarshal(msg.Value, order)
	if err != nil {
		return err
	}

	//var id uuid.UUID
	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		_, err = s.orderRepository.CreateOrder(ctx, order)
		if err != nil {
			return err
		}

		_, err = s.orderRepository.CreateDelivery(ctx, order.OrderUID, &order.Delivery)
		if err != nil {
			return err
		}

		_, err = s.orderRepository.CreatePayment(ctx, order.OrderUID, &order.Payment)
		if err != nil {
			return err
		}

		for _, item := range order.Items {
			err = s.orderRepository.CreateItem(ctx, order.OrderUID, &item)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.cache.Set(order.OrderUID, order)

	return nil
}
