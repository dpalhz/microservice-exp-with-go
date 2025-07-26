package dto

import "github.com/google/uuid"

// RegisterRequest represents a user registration payload.
type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse is returned after a successful registration.
type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
