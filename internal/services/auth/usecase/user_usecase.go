package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
)

type UserUsecase struct {
	userRepo UserRepository
	producer EventProducer
	tokenGen TokenGenerator
	session  SessionStore
	verRepo  VerificationRepository
	log      *slog.Logger
}

func NewUserUsecase(repo UserRepository, producer EventProducer, tokenGen TokenGenerator, session SessionStore, verRepo VerificationRepository, log *slog.Logger) *UserUsecase {
	return &UserUsecase{userRepo: repo, producer: producer, tokenGen: tokenGen, session: session, verRepo: verRepo, log: log}
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (uc *UserUsecase) VerifyCode(ctx context.Context, userID uuid.UUID, code string, purpose domain.VerificationPurpose) (string, string, error) {
	vc, err := uc.verRepo.FindValid(ctx, userID, purpose, code)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}
	vc.Status = "USED"
	now := time.Now()
	vc.UsedAt = &now
	if err := uc.verRepo.Update(ctx, vc); err != nil {
		uc.log.Error("failed update verification", slog.String("error", err.Error()))
	}

	if purpose == domain.PurposeEnable2FA {
		user, err := uc.userRepo.FindByID(ctx, vc.UserID)
		if err == nil {
			user.MFAEnabled = true
			_ = uc.userRepo.Update(ctx, user)
		}
	}

	accessToken, refreshToken, err := uc.tokenGen.GenerateTokens(userID)
	if err != nil {
		return "", "", err
	}
	if err := uc.session.Store(ctx, refreshToken, userID); err != nil {
		uc.log.Error("failed to store refresh token", slog.String("error", err.Error()))
	}
	return accessToken, refreshToken, nil
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

	// generate verification code for registration
	code := generateCode()
	vc := &domain.VerificationCode{
		ID:        uuid.New(),
		UserID:    user.ID,
		Code:      code,
		Purpose:   domain.PurposeRegister,
		Status:    "UNUSED",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}
	if err := uc.verRepo.Create(ctx, vc); err == nil {
		if err := uc.producer.ProduceEmailVerification(ctx, user.Email, user.ID, code, domain.PurposeRegister); err != nil {
			uc.log.Error("send email event failed", slog.String("error", err.Error()))
		}
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

	if user.MFAEnabled {
		code := generateCode()
		vc := &domain.VerificationCode{
			ID:        uuid.New(),
			UserID:    user.ID,
			Code:      code,
			Purpose:   domain.PurposeLogin,
			Status:    "UNUSED",
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}
		if err := uc.verRepo.Create(ctx, vc); err == nil {
			if err := uc.producer.ProduceEmailVerification(ctx, user.Email, user.ID, code, domain.PurposeLogin); err != nil {
				uc.log.Error("send email event failed", slog.String("error", err.Error()))
			}
		}
		return "", "", domain.ErrVerificationRequired
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
