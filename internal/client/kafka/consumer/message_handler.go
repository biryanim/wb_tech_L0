package consumer

import (
	"github.com/IBM/sarama"
	"github.com/biryanim/wb_tech_L0/internal/client/kafka"
	"log"
)

type GroupHandler struct {
	msgHandler kafka.Handler
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{}
}

func (c *GroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *GroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

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
