package producer

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"github.com/pkg/errors"
	"log"
	"time"
)

var _ kafka.Producer = (*producer)(nil)

type producer struct {
	syncProducer sarama.SyncProducer
}

func NewProducer(syncProducer sarama.SyncProducer) *producer {
	return &producer{
		syncProducer: syncProducer,
	}
}

func (p *producer) SendToDLQ(ctx context.Context, message *kafka.DLQMessage, topic string) error {
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
		return errors.Wrap(err, "failed to marshal DLQ message")
	}

	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Topic),
		Value: sarama.ByteEncoder(payload),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to send message to DLQ topic %s", topic)
	}

	log.Printf("DLQ message sent: partition=%d, offset=%d", partition, offset)

	return nil
}

func (p *producer) Close() error {
	return p.syncProducer.Close()
}
