package service

import (
	"fmt"
	"net/smtp"

	"github.com/Nucleussss/auth-service/pkg/logger"
)

type EmailService interface {
	SendPasswordResetEmail(email, resetToken string) error
	// Other email methods can be added here
}

type smtpEmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	logger       logger.Logger
}

func NewSmtpEmailService(host, port, username, password, from string, logger logger.Logger) EmailService {
	return &smtpEmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    from,
		logger:       logger,
	}
}

func (s *smtpEmailService) SendPasswordResetEmail(email, resetToken string) error {
	const op = "emailService.SendPasswordResetEmail"

	// Create the full address with host and port
	smtpAddress := s.smtpHost + ":" + s.smtpPort

	// Authenticate
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// Email content
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Password Reset\r\n\r\nReset token: %s",
		s.fromEmail, email, resetToken)

	//	Send the message
	err := smtp.SendMail(smtpAddress, auth, s.fromEmail, []string{email}, []byte(msg))
	if err != nil {
		s.logger.Errorf("Invalid password reset request: %v ", op, err)
		return err
	}

	return nil
}
