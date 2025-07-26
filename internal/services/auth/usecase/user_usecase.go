package usecase

import (
	"context"
	"log/slog"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
)

type UserUsecase struct {
	userRepo UserRepository
	producer EventProducer
	tokenGen TokenGenerator
	session  SessionStore
	log      *slog.Logger
}

func NewUserUsecase(repo UserRepository, producer EventProducer, tokenGen TokenGenerator, session SessionStore, log *slog.Logger) *UserUsecase {
	return &UserUsecase{userRepo: repo, producer: producer, tokenGen: tokenGen, session: session, log: log}
}

func (uc *UserUsecase) Register(ctx context.Context, fullName, email, password string) (*domain.User, error) {
	user := &domain.User{
		FullName: fullName,
		Email:    email,
		Password: password,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.log.Error("failed to create user", slog.String("error", err.Error()))
		return nil, err
	}

	if err := uc.producer.ProduceUserRegistered(ctx, user); err != nil {
		// Log the error but do not return it, as user creation should not fail due to event production issues.	
		uc.log.Error("failed to produce user registered event", slog.String("error", err.Error()))
	}

	return user, nil
}

func (uc *UserUsecase) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !user.CheckPassword(password) {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, refreshToken, err = uc.tokenGen.GenerateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	if err := uc.session.Store(ctx, refreshToken, user.ID); err != nil {
		uc.log.Error("failed to store refresh token", slog.String("error", err.Error()))
	}

	return accessToken, refreshToken, nil
}
