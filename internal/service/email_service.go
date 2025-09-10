package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
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
}

func NewSmtpEmailService(host, port, username, password, from string) EmailService {
	return &smtpEmailService{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
		fromEmail:    from,
	}
}

func (s *smtpEmailService) SendPasswordResetEmail(email, resetToken string) error {
	// Gmail requires TLS
	tlsconfig := &tls.Config{
		ServerName: s.smtpHost,
	}

	// Connect to SMTP server
	conn, err := tls.Dial("tcp", s.smtpHost+":"+s.smtpPort, tlsconfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.smtpHost)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Set sender and recipient
	if err = client.Mail(s.fromEmail); err != nil {
		return err
	}
	if err = client.Rcpt(email); err != nil {
		return err
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Email content
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Password Reset\r\n\r\nReset token: %s",
		s.fromEmail, email, resetToken)

	_, err = w.Write([]byte(msg))
	return err
}
