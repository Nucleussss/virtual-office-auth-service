package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Nucleussss/auth-service/internal/db/models"
	"github.com/Nucleussss/auth-service/internal/repositories"
	"github.com/Nucleussss/auth-service/internal/utils"
	"github.com/Nucleussss/auth-service/pkg/logger"
	"github.com/google/uuid"
)

type AuthService struct {
	repo   repositories.UserRepository
	logger logger.Logger
}

func NewAuthService(repo repositories.UserRepository, logger logger.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
	}

}

func (s *AuthService) Register(ctx context.Context, user *models.RegisterRequest) error {
	const op = "AuthService.Register"

	s.logger.Infof("%s: Attemp to registration for %s", op, user.Email)

	// Check if Email already exists
	exists, err := s.repo.ExistsbyEmail(ctx, user.Email)
	if err != nil {
		s.logger.Errorf("%s: Error checking if user exists: %s %v", op, user.Email, err)
		return fmt.Errorf("Registration unavaible")
	}

	if exists {
		s.logger.Errorf("%s: Duplicate Email found for: %s", op, user.Email)
		return fmt.Errorf("Email already registered")
	}

	// Hash the password before storing it in the database
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Errorf("%s: Error hashing password: %v", op, err)
		return fmt.Errorf("Password hashing failed")
	}

	userToCreateNewUser := &models.CreateNewUser{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: hashPassword,
	}

	// Create the new user in the database
	err = s.repo.Create(ctx, userToCreateNewUser)
	if err != nil {
		s.logger.Errorf("%s: Error creating user: %s %v", op, user.Email, err)
		return fmt.Errorf("Registration failed")
	}

	s.logger.Infof("%s: Successfully registered user: %s", op, user.Email)
	return nil
}

func (s *AuthService) Login(ctx context.Context, userLoginRequest *models.LoginRequest) (string, error) {
	const op = "handlers.LoginHandler"
	s.logger.Infof("%s: Attempting to login with email: %s", op, userLoginRequest.Email)

	// Find by Email
	user, err := s.repo.FindbyEmail(ctx, userLoginRequest.Email)
	if err != nil {
		s.logger.Errorf("%s: Failed to find user by email: %v", op, err)
		return "", fmt.Errorf("Failed to find user by email")
	}

	// Verify the password hash
	if err := utils.VerifyPassword(user.PasswordHash, userLoginRequest.Password); err != nil {
		s.logger.Errorf("%s: Failed to verify password: %v", op, err)
		return "", fmt.Errorf("Failed to verify password")
	}

	// load the JWT secret from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Parse the string ID into UUID
	// userUUID, err := uuid.Parse(user.ID)
	// if err != nil {
	// 	return "", fmt.Errorf("Failed to parse UUID from token")
	// }

	// Generate a JWT token for the user
	token, err := utils.GenerateJWTToken(user.ID, jwtSecret)
	if err != nil {
		s.logger.Errorf("%s: Failed to generate JWT token: %v", op, err)
		return "", fmt.Errorf("Failed to generate JWT token")
	}

	s.logger.Infof("%s: Successfully logged in user: %s", op, userLoginRequest.Email)

	return token, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const op = "handlers.GetProfileHandler"
	s.logger.Infof("%s: Attempting to get Profile with id: %s", op, userID)

	// Find the user by ID
	user, err := s.repo.FindbyID(ctx, userID)
	if err != nil {
		s.logger.Errorf("%s: Failed to find user by ID: %v", op, err)
		return nil, fmt.Errorf("Failed to find user")
	}

	return user, nil
}
