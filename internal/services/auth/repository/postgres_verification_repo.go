package repository

import (
	"context"
	"errors"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresVerificationRepository struct {
	db *gorm.DB
}

func NewPostgresVerificationRepository(db *gorm.DB) *PostgresVerificationRepository {
	return &PostgresVerificationRepository{db: db}
}

func (r *PostgresVerificationRepository) Create(ctx context.Context, v *domain.VerificationCode) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *PostgresVerificationRepository) FindValid(ctx context.Context, userID uuid.UUID, purpose domain.VerificationPurpose, code string) (*domain.VerificationCode, error) {
	var vc domain.VerificationCode
	if err := r.db.WithContext(ctx).Where("user_id = ? AND purpose = ? AND code = ? AND status = ? AND expires_at > now()", userID, purpose, code, "UNUSED").First(&vc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &vc, nil
}

func (r *PostgresVerificationRepository) Update(ctx context.Context, v *domain.VerificationCode) error {
	return r.db.WithContext(ctx).Save(v).Error
}
