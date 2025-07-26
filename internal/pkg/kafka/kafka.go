package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerConfig struct {
	BootstrapServers string
	Topic            string
}

// NewProducer creates a new Kafka producer instance.
// It initializes the producer with the provided configuration.
func NewProducer(cfg ProducerConfig) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.BootstrapServers})
}

// NewConsumer creates a new Kafka consumer instance.
// It initializes the consumer with the provided configuration.	
func NewConsumer(bootstrapServers, groupID string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
}
