package main

import (
	"fmt"
	"net/http"

	"time"

	"github.com/Nucleussss/auth-service/internal/config"
	"github.com/Nucleussss/auth-service/internal/db"
	"github.com/Nucleussss/auth-service/internal/handlers"
	"github.com/Nucleussss/auth-service/internal/middleware"
	"github.com/Nucleussss/auth-service/internal/repositories"
	"github.com/Nucleussss/auth-service/internal/service"
	"github.com/Nucleussss/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from .env file
	config := config.LoadConfig()

	// Initialize logger
	log := logger.NewLogger()

	// Initialize database connection
	dbconn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	defer dbconn.Close()
	log.Infof("Database connected successfully")

	// Initialize user repository
	userRepo := repositories.NewUserRepository(dbconn)

	// Initialize auth service
	authService := service.NewAuthService(userRepo, log)

	// Initialize password reset repository
	passwordResetRepo := repositories.NewPasswordRepository(dbconn)

	// Initialize email service
	emailService := service.NewSmtpEmailService(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUser,
		config.SMTPPass,
		config.SMTPFrom,
		log,
	)

	// Parse token expiration duration
	duration, err := time.ParseDuration(config.TokenExpiration)
	if err != nil {
		log.Fatalf("Error parsing token expiration duration: %v", err)
	}

	// Initialize password reset service
	passwordResetService := service.NewPasswordResetService(
		log,
		userRepo,
		passwordResetRepo,
		emailService,
		duration,
	)

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(authService, passwordResetService, log)

	// Initialize gin router
	router := gin.Default()

	// Register routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	//
	router.POST("/request-password-reset", authHandler.RequestPasswordReset)
	router.POST("/reset-password", authHandler.ResetPassword)

	// protected API group
	api := router.Group("/api")
	api.Use(middleware.JWTMiddleware(config.JWTSecret, log))
	{
		api.GET("/profile", authHandler.Profile)
	}

	// Start the server
	addr := fmt.Sprintf(":%s", config.ServerPort)
	log.Infof("Server is running on port %s", config.ServerPort)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
