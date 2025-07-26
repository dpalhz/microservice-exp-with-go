package dto

import "github.com/google/uuid"

type Enable2FARequest struct {
	UserID uuid.UUID `json:"user_id"`
	Code   string    `json:"code"`
}
