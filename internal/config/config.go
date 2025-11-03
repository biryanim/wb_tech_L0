package config

import (
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

// PGConfig defines methods for retrieving PostgreSQL database configuration.
type PGConfig interface {
	DSN() string
}

// HTTPConfig defines methods for retrieving HTTP server configuration.
type HTTPConfig interface {
	Address() string
}

// KafkaConsumerConfig defines methods for retrieving Kafka consumer configuration.
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

// JaegerConfig defines methods for retrieving Jaeger tracing configuration.
type JaegerConfig interface {
	URL() string
	ServiceName() string
}

// Load reads and parses environment variables from the specified file path using godotenv.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
