package kafka

import (
	"context"
	"github.com/IBM/sarama"
)

type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler Handler) (err error)
	Close() error
}
