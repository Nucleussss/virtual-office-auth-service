package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWTToken(userID uuid.UUID, secretKey string) (string, error) {
	expString := os.Getenv("JWT_EXPIRATION")

	// convert the expiration string to an integer duration in hours
	expDuration, err := strconv.Atoi(expString)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRATION: %v", err)
	}

	// create a new JWT token with the specified claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * time.Duration(expDuration)).Unix(),
	})

	// sign the token with the secret key and return it
	return token.SignedString([]byte(secretKey))
}

func ValidateJWTToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		// check if the token is signed with HS256 algorithm
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// return the secret key as the token's key
		return []byte(secretKey), nil
	})

	// check if there was an error during parsing or if the token is expired
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("JWT token is invalid: %v", err)
	}

	// extract the claims from the token and return them
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	// return the claims from the token
	return &claims, nil
}
