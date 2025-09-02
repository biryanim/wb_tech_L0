package env

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	fmt.Println(brokersStr)
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

func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
