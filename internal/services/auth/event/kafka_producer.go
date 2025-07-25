package event

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/google/uuid"
)

type KafkaEventProducer struct {
	producer *kafka.Producer
	topic    string
	log      *slog.Logger
}

func NewKafkaEventProducer(p *kafka.Producer, topic string, log *slog.Logger) *KafkaEventProducer {
	return &KafkaEventProducer{producer: p, topic: topic, log: log}
}

type UserRegisteredEvent struct {
	UserID       uuid.UUID    `json:"user_id"`
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
	if err!= nil {
		return err
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: int32(kafka.PartitionAny)},
		Value:          payload,
	}, nil)
}