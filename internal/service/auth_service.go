package service

import (
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

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger.NewLogger(),
	}
}

func (s *AuthService) Register(name, email, password string) error {
	const op = "AuthService.Register"

	s.logger.Infof("%s: Attemp to registration for %s", op, email)

	// Check if Email already exists
	exists, err := s.repo.ExistsbyEmail(email)
	if err != nil {
		s.logger.Fatalf("%s: Error checking if user exists: %s %v", op, email, err)
		return fmt.Errorf("Registration unavaible")
	}

	if exists {
		s.logger.Infof("%s: Duplicate Email found for: %s", op, email)
		return fmt.Errorf("Email already registered")
	}

	// Hash the password before storing it in the database
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		s.logger.Fatalf("%s: Error hashing password: %v", op, err)
		return fmt.Errorf("Password hashing failed")
	}

	// Create the new user in the database
	err = s.repo.Create(name, email, hashPassword)
	if err != nil {
		s.logger.Fatalf("%s: Error creating user: %s %v", op, email, err)
		return fmt.Errorf("Registration failed")
	}

	s.logger.Infof("%s: Successfully registered user: %s", op, email)
	return nil
}

func (s *AuthService) Login(name, email, password string) (string, error) {
	const op = "handlers.LoginHandler"
	s.logger.Infof("%s: Attempting to login with name: %s", op, name)

	// Find by Email
	user, err := s.repo.FindbyEmail(email)
	if err != nil {
		s.logger.Errorf("%s: Failed to find user by email: %v", op, err)
		return "", fmt.Errorf("Failed to find user by email")
	}

	// Verify the password hash
	if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
		s.logger.Errorf("%s: Failed to verify password: %v", op, err)
		return "", fmt.Errorf("Failed to verify password")
	}

	// load the JWT secret from environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Parse the string ID into UUID
	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return "", fmt.Errorf("Failed to parse UUID from token")
	}

	// Generate a JWT token for the user
	token, err := utils.GenerateJWTToken(userUUID, jwtSecret)
	if err != nil {
		s.logger.Errorf("%s: Failed to generate JWT token: %v", op, err)
		return "", fmt.Errorf("Failed to generate JWT token")
	}

	s.logger.Infof("%s: Successfully logged in user: %s", op, name)
	return token, nil
}

func (s *AuthService) GetProfile(userID uuid.UUID) (*models.User, error) {
	const op = "handlers.GetProfileHandler"
	s.logger.Infof("%s: Attempting to get Profile with id: %s", op, userID)

	// Find the user by ID
	user, err := s.repo.FindbyID(userID)
	if err != nil {
		s.logger.Errorf("%s: Failed to find user by ID: %v", op, err)
		return nil, fmt.Errorf("Failed to find user")
	}

	return user, nil
}
