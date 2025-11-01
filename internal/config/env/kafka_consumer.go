package env

import (
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

// NewKafkaConsumerConfig creates and returns a new Kafka consumer configuration by reading broker addresses and group ID from environment variables.
func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("kafka group id not found")
	}

	return &kafkaConsumerConfig{brokers, groupID}, nil
}

// Brokers returns the list of Kafka broker addresses.
func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

// GroupID returns the consumer group identifier.
func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

// Config returns a configured Sarama consumer configuration with version 2.6.0.0, round-robin rebalancing, and oldest offset strategy.
func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
