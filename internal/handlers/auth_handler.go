package handlers

import (
	"net/http"

	"github.com/Nucleussss/auth-service/internal/db/models"
	"github.com/Nucleussss/auth-service/internal/service"
	"github.com/Nucleussss/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService          *service.AuthService
	passwordResetService *service.PasswordResetService
	logger               logger.Logger
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.NewLogger(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	const op = "handlers.Register"
	var req models.RegisterRequest

	// parse the JSON body from the request
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("%s: failed to parse JSON body: %v", op, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request",
			"detail": err.Error(),
		})
		return
	}

	// register the user
	if err := h.authService.Register(req.Name, req.Email, req.PasswordHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "registration failed",
			"detail": err.Error(),
		})
		return
	}

	// send a success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"data": gin.H{
			"name":  req.Name,
			"email": req.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	const op = "handlers.Login"
	var req models.LoginRequest

	// bind the request body to the LoginRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("%s: failed to bind JSON: %v", op, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request",
			"detail": err.Error(),
		})
		return
	}

	// validate the user credentials
	token, err := h.authService.Login(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "invalid credentials",
			"detail": err.Error(),
		})
		return
	}

	// return a success response with the user data
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data": gin.H{
			"name":  req.Name,
			"email": req.Email,
			"token": token,
		},
	})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	const op = "handlers.GetProfile"
	h.logger.Infof("%s: received request to get profile for user ID %d", op, c.MustGet("user_id").(uuid.UUID))

	// get the user ID from the request header
	userID := c.MustGet("user_id").(uuid.UUID)

	// retrieve the user data from the database
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		h.logger.Errorf("%s: failed to find user by ID %d", op, userID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// return the user data as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req models.PasswordResetRequest

	//
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid password reset request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON"})
		return
	}

	_, err := h.passwordResetService.RequestReset(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Errorf("Password reset request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"massage": "If email is correct, you will receive a password reset"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.NewPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid new password request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind JSON"})
		return
	}

	err := h.passwordResetService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		h.logger.Errorf("Password reset failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	c.JSON(200, gin.H{"massage": "Password updated succesfully"})

}
