package service

import (
	"github.com/Nucleussss/auth-service/internal/repositories"
	"github.com/Nucleussss/auth-service/pkg/logger"
)

type PasswordResetService struct {
	repo   repositories.UserRepository
	logger logger.Logger
}
