package usecase

import (
	"context"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type EventProducer interface {
	ProduceUserRegistered(ctx context.Context, user *domain.User) error
	ProduceEmailVerification(ctx context.Context, email string, userID uuid.UUID, code string, purpose domain.VerificationPurpose) error
}

type TokenGenerator interface {
	GenerateTokens(userID uuid.UUID) (string, string, error)
}

type SessionStore interface {
	Store(ctx context.Context, token string, userID uuid.UUID) error
	GetUserID(ctx context.Context, token string) (uuid.UUID, error)
}

type VerificationRepository interface {
	Create(ctx context.Context, v *domain.VerificationCode) error
	FindValid(ctx context.Context, userID uuid.UUID, purpose domain.VerificationPurpose, code string) (*domain.VerificationCode, error)
	Update(ctx context.Context, v *domain.VerificationCode) error
}
