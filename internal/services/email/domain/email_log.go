package domain

import (
	"github.com/google/uuid"
	"time"
)

type EmailLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	To        string
	Subject   string
	Body      string
	Status    string
	ErrorMsg  string
	CreatedAt time.Time
}
