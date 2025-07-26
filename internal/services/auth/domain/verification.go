package domain

import (
	"github.com/google/uuid"
	"time"
)

type VerificationPurpose string

const (
	PurposeRegister  VerificationPurpose = "register"
	PurposeLogin     VerificationPurpose = "login_2fa"
	PurposeEnable2FA VerificationPurpose = "enable_2fa"
)

// VerificationCode stores OTP for various purposes
// gorm model.
type VerificationCode struct {
	ID           uuid.UUID           `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID           `gorm:"type:uuid;index"`
	Code         string              `gorm:"size:6"`
	Purpose      VerificationPurpose `gorm:"type:varchar(20)"`
	Status       string              `gorm:"type:varchar(10)"`
	AttemptCount int
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UsedAt       *time.Time
}
