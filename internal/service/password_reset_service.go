package service

import (
	"context"
	"errors"
	"time"

	"github.com/Nucleussss/auth-service/internal/db/models"
	"github.com/Nucleussss/auth-service/internal/repositories"
	"github.com/Nucleussss/auth-service/pkg/logger"

	// import utils
	"github.com/Nucleussss/auth-service/internal/utils"
)

type PasswordResetService interface {
	RequestReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type passwordResetService struct {
	logger            logger.Logger
	userRepo          repositories.UserRepository
	passwordResetRepo repositories.PasswordResetRepository
	emailService      EmailService
	tokenExpiry       time.Duration
}

func NewPasswordResetService(
	logger logger.Logger,
	userRepo repositories.UserRepository,
	passwordResetRepo repositories.PasswordResetRepository,
	emailService EmailService,
	tokenExpiry time.Duration,
) PasswordResetService {
	return &passwordResetService{
		logger:            logger,
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		emailService:      emailService,
		tokenExpiry:       tokenExpiry,
	}
}

func (s *passwordResetService) RequestReset(ctx context.Context, email string) (string, error) {
	// var op = "PasswordResetService.RequestReset"

	user, err := s.userRepo.FindbyEmail(ctx, email)
	if err != nil {
		return " ", err
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return " ", err
	}

	// Create a new password reset record
	reset := &models.PasswordReset{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
	}

	// Save the password reset record to the database
	if err = s.passwordResetRepo.Create(ctx, reset); err != nil {
		return " ", err
	}

	if err := s.emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		return "", err
	}

	return token, nil

}

func (s *passwordResetService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	reset, err := s.passwordResetRepo.FindValidToken(ctx, token)
	if err != nil {
		return err
	}
	if reset == nil {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, reset.UserID, hashedPassword); err != nil {
		return err
	}

	return s.passwordResetRepo.Delete(ctx, token)
}
