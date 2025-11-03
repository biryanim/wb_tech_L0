package producer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"github.com/pkg/errors"
)

var _ kafka.Producer = (*producer)(nil)

type producer struct {
	syncProducer sarama.SyncProducer
	tracer       trace.Tracer
}

func NewProducer(syncProducer sarama.SyncProducer, trace trace.Tracer) *producer {
	return &producer{
		syncProducer: syncProducer,
		tracer:       trace,
	}
}

func (p *producer) SendToDLQ(ctx context.Context, message *kafka.DLQMessage, topic string) error {
	ctx, span := p.tracer.Start(ctx, "send_to_dlq")
	defer span.End()

	span.SetAttributes(
		attribute.String("kafka.topic", topic),
		attribute.String("original.topic", message.Topic),
		attribute.String("error.reason", message.ErrorReason),
		attribute.Int64("original.offset", message.Offset),
	)

	dlqPayload := map[string]interface{}{
		"original_message":   string(message.OriginalMessage),
		"original_topic":     message.Topic,
		"error_reason":       message.ErrorReason,
		"timestamp":          message.Timestamp,
		"original_offset":    message.Offset,
		"original_partition": message.Partition,
		"dlq_received_at":    time.Now().UnixMilli(),
	}

	payload, err := json.Marshal(dlqPayload)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "marshal"))
		return errors.Wrap(err, "failed to marshal DLQ message")
	}

	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Topic),
		Value: sarama.ByteEncoder(payload),
	})

	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("error.type", "send"))
		return errors.Wrapf(err, "failed to send message to DLQ topic %s", topic)
	}

	log.Printf("DLQ message sent: partition=%d, offset=%d", partition, offset)

	return nil
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
