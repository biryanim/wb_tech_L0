package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

// Handler is a function type that processes Kafka consumer messages within a given context.
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

// Consumer defines methods for consuming messages from Kafka topics.
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler Handler) (err error)
	Close() error
}

// Producer defines the interface for sending messages to a Dead Letter Queue.
type Producer interface {
	SendToDLQ(ctx context.Context, message *DLQMessage, topic string) error
	Close() error
}

// DLQMessage represents a message that failed processing and is being sent to the Dead Letter Queue.
type DLQMessage struct {
	OriginalMessage []byte
	Topic           string
	ErrorReason     string
	Timestamp       int64
	Partition       int32
	Offset          int64
}
