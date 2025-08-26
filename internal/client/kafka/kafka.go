package kafka

import (
	"context"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
