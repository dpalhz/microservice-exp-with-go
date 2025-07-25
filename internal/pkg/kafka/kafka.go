package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerConfig struct {
	BootstrapServers string
	Topic            string
}

// NewProducer membuat instance produsen Kafka baru.
// Konfigurasi diambil dari ProducerConfig untuk memudahkan injeksi dependensi.
func NewProducer(cfg ProducerConfig) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.BootstrapServers})
}

// NewConsumer membuat instance konsumen Kafka baru.
func NewConsumer(bootstrapServers, groupID string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
}
