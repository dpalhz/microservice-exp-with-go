package dto

import (
	"github.com/dpalhz/microservice-exp-with-go/internal/services/auth/domain"
	"github.com/google/uuid"
)

type VerifyRequest struct {
	UserID  uuid.UUID                  `json:"user_id"`
	Code    string                     `json:"code"`
	Purpose domain.VerificationPurpose `json:"purpose"`
}

type VerifyResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
