package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string `env:"DB_HOST"`
	DBPort        string `env:"DB_PORT"`
	DBUser        string `env:"DB_USER"`
	DBPassword    string `env:"DB_PASSWORD"`
	DBName        string `env:"DB_NAME"`
	ServerPort    string `env:"SERVER_PORT"`
	JWTSecret     string `env:"JWT_SECRET"`
	JWTExpiration string `env:"JWT_EXPIRATION"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
	}
}
