package consumer

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
)

// GroupHandler implements sarama.ConsumerGroupHandler for processing Kafka consumer group messages.
type GroupHandler struct {
	msgHandler kafka.Handler
}

// NewGroupHandler creates and returns a new GroupHandler instance.
func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

// Setup is called at the beginning of a new consumer group session before any messages are consumed.
func (c *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called at the end of a consumer group session after all messages have been processed.
func (c *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a claimed partition, handling each message with the configured handler.
func (c *GroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			log.Printf("message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)

			err := c.msgHandler(session.Context(), message)
			if err != nil {
				log.Printf("error handling message: %v", err)
				continue
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			log.Printf("session context done\n")
			return nil
		}
	}
}
