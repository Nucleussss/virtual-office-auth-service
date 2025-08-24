package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type NewPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
