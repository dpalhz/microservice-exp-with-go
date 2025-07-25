package usecase

import (
	"context"
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
)

type UserUsecase struct {
	userRepo  UserRepository
	producer  EventProducer
	tokenGen  TokenGenerator
	log       *slog.Logger
}

func NewUserUsecase(repo UserRepository, producer EventProducer, tokenGen TokenGenerator, log *slog.Logger) *UserUsecase {
	return &UserUsecase{userRepo: repo, producer: producer, tokenGen: tokenGen, log: log}
}

func (uc *UserUsecase) Register(ctx context.Context, fullName, email, password string) (*domain.User, error) {
	user := &domain.User{
		FullName: fullName,
		Email:    email,
		Password: password,
	}

	if err := uc.userRepo.Create(ctx, user); err!= nil {
		uc.log.Error("Gagal membuat pengguna di repo", slog.String("error", err.Error()))
		return nil, err
	}

	if err := uc.producer.ProduceUserRegistered(ctx, user); err!= nil {
		// Log error tapi jangan gagalkan registrasi
		uc.log.Error("Gagal memproduksi event UserRegistered", slog.String("error", err.Error()))
	}

	return user, nil
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err!= nil {
		return "", "", err
	}

	if!user.CheckPassword(password) {
		return "", "", domain.ErrInvalidCredentials
	}

	return uc.tokenGen.GenerateTokens(user.ID)
}