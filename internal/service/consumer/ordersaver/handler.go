package ordersaver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"github.com/biryanim/wb_tech_L0/internal/model"
	"github.com/biryanim/wb_tech_L0/internal/validator"
)

// OrderSaveHandler processes a Kafka message containing order data, persists it to the database within a transaction, and caches the order.
func (s *service) OrderSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	order := &model.Order{}
	err := json.Unmarshal(msg.Value, order)
	if err != nil {
		return s.sendToDLQ(ctx, msg, fmt.Sprintf("failed to unmarshal order: %v", err))
	}

	validationResult := validator.ValidateOrder(order)
	if !validationResult.Valid {
		errorMsg := fmt.Sprintf("validation failed: %v", validationResult.Errors)
		log.Printf("Validation error for order %s: %s", order.OrderUID, errorMsg)
		err = s.sendToDLQ(ctx, msg, errorMsg)
		return err
	}

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

func (s *service) sendToDLQ(ctx context.Context, msg *sarama.ConsumerMessage, errorReason string) error {
	dlqMsg := &kafka.DLQMessage{
		OriginalMessage: msg.Value,
		Topic:           msg.Topic,
		ErrorReason:     errorReason,
		Timestamp:       msg.Timestamp.UnixMilli(),
		Partition:       msg.Partition,
		Offset:          msg.Offset,
	}

	err := s.producer.SendToDLQ(ctx, dlqMsg, "dlq-order")
	if err != nil {
		return err
	}
	return nil
}
