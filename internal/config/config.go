package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost          string `env:"DB_HOST"`
	DBPort          string `env:"DB_PORT"`
	DBUser          string `env:"DB_USER"`
	DBPassword      string `env:"DB_PASSWORD"`
	DBName          string `env:"DB_NAME"`
	ServerPort      string `env:"SERVER_PORT"`
	JWTSecret       string `env:"JWT_SECRET"`
	JWTExpiration   string `env:"JWT_EXPIRATION"`
	TokenExpiration string `env:"TOKEN_EXPIRATION"`
	SMTPHost        string `env:"SMTP_HOST"`
	SMTPPort        string `env:"SMTP_PORT"`
	SMTPUser        string `env:"SMTP_USER"`
	SMTPPass        string `env:"SMTP_PASS"`
	SMTPFrom        string `env:"SMTP_FROM"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		ServerPort:      os.Getenv("SERVER_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiration:   os.Getenv("JWT_EXPIRATION"),
		TokenExpiration: os.Getenv("TOKEN_EXPIRATION"),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        os.Getenv("SMTP_PORT"),
		SMTPUser:        os.Getenv("SMTP_USER"),
		SMTPPass:        os.Getenv("SMTP_PASS"),
		SMTPFrom:        os.Getenv("SMTP_FROM"),
	}
}
