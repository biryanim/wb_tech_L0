package consumer

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"github.com/pkg/errors"
)

var _ kafka.Consumer = (*consumer)(nil)

type consumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

// NewConsumer creates and returns a new Kafka consumer with the provided consumer group and handler.
func NewConsumer(consumerGroup sarama.ConsumerGroup, consumerGroupHandler *GroupHandler) *consumer {
	return &consumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

// Consume starts consuming messages from the specified topic and processes them with the provided handler.
func (c *consumer) Consume(ctx context.Context, topicName string, handler kafka.Handler) error {
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx, topicName)
}

// Close closes the consumer group and releases associated resources.
func (c *consumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *consumer) consume(ctx context.Context, topicName string) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(topicName, "."), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		log.Printf("rebalancing...\n")
	}
}
