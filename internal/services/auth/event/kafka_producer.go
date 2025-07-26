package event

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/pkg/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/google/uuid"
)

type KafkaEventProducer struct {
	producer *ckafka.Producer
	topic    string
	log      *slog.Logger
}

// NewKafkaEventProducer creates a new Kafka event producer.	
// It initializes the producer with the provided Kafka configuration and logger.
func NewKafkaEventProducer(p *ckafka.Producer, cfg kafka.ProducerConfig, log *slog.Logger) *KafkaEventProducer {
	return &KafkaEventProducer{producer: p, topic: cfg.Topic, log: log}
}

type UserRegisteredEvent struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	RegisteredAt time.Time `json:"registered_at"`
}

func (p *KafkaEventProducer) ProduceUserRegistered(ctx context.Context, user *domain.User) error {
	event := UserRegisteredEvent{
		UserID:       user.ID,
		Email:        user.Email,
		FullName:     user.FullName,
		RegisteredAt: user.CreatedAt,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &p.topic, Partition: int32(ckafka.PartitionAny)},
		Value:          payload,
	}, nil)
}
