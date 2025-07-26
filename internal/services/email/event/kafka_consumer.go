package event

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	authdomain "github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/domain"
	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/repository"
)

type KafkaConsumer struct {
	consumer *ckafka.Consumer
	repo     *repository.PostgresLogRepository
	log      *slog.Logger
}

type EmailVerificationEvent struct {
	UserID  uuid.UUID                      `json:"user_id"`
	Email   string                         `json:"email"`
	Code    string                         `json:"code"`
	Purpose authdomain.VerificationPurpose `json:"purpose"`
}

func NewKafkaConsumer(c *ckafka.Consumer, repo *repository.PostgresLogRepository, log *slog.Logger) *KafkaConsumer {
	return &KafkaConsumer{consumer: c, repo: repo, log: log}
}

func (kc *KafkaConsumer) Start(ctx context.Context, topic string) error {
	if err := kc.consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			kc.log.Error("consumer error", slog.String("error", err.Error()))
			continue
		}
		var evt EmailVerificationEvent
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			kc.log.Error("unmarshal", slog.String("error", err.Error()))
			continue
		}
		// simulate email send
		el := domain.EmailLog{To: evt.Email, Subject: "Your Code", Body: evt.Code, Status: "SENT", CreatedAt: msg.Timestamp}
		if err := kc.repo.Create(ctx, &el); err != nil {
			kc.log.Error("save log", slog.String("error", err.Error()))
		}
	}
}
