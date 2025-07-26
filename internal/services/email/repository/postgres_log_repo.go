package repository

import (
	"context"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/email/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresLogRepository struct{ db *gorm.DB }

func NewPostgresLogRepository(db *gorm.DB) *PostgresLogRepository {
	return &PostgresLogRepository{db: db}
}

func (r *PostgresLogRepository) Create(ctx context.Context, log *domain.EmailLog) error {
	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(log).Error
}
