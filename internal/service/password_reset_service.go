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
	var op = "PasswordResetService.RequestReset"

	user, err := s.userRepo.FindbyEmail(ctx, email)
	if err != nil {
		s.logger.Errorf("%s: Failed find email in the database: %v ", op, err)
		return " ", err
	}

	// Generate a secure token for the password reset request
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
		s.logger.Errorf("%s: Failed Save the password reset record to the database: %v ", op, err)
		return " ", err
	}

	// Send the password reset email to the user
	if err := s.emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		return "", err
	}

	return token, nil
}

// explain this function.
func (s *passwordResetService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	var op = "PasswordResetService.ResetPassword"

	// Find the password reset record in the database using the token
	reset, err := s.passwordResetRepo.FindValidToken(ctx, token)
	if err != nil {
		s.logger.Errorf("%s: Failed to find password reset record: %v", op, err)
		return err
	}
	if reset == nil {
		s.logger.Errorf("%s: Failed to find password reset record: %v", op, err)
		return errors.New("invalid or expired token")
	}

	// Hash the new password before updating it in the database
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorf("%s: Failed to hash password %v", op, err)
		return err
	}

	// Update the user's password in the database with the
	if err := s.userRepo.UpdatePassword(ctx, reset.UserID, hashedPassword); err != nil {
		s.logger.Errorf("%s: Failed to update user password %v", op, err)
		return err
	}

	return s.passwordResetRepo.Delete(ctx, token)
}
