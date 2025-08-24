package middleware

import (
	"strings"

	"github.com/Nucleussss/auth-service/internal/utils"
	"github.com/Nucleussss/auth-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func JWTMiddleware(secretKey string, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		op := "middleware.JWTMiddleware"

		// get the header from the request
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Errorf("%s: Authorization header is missing", op)
			c.JSON(401, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate the token
		claims, err := utils.ValidateJWTToken(tokenString, secretKey)
		if err != nil {
			log.Errorf("%s: Invalid JWT token", op)
			c.JSON(401, gin.H{
				"error": "Invalid JWT token",
			})
			c.Abort()
			return
		}

		// set the claims to the context
		// userID := int((*claims)["user_id"].(float64))

		userIDstr, ok := (*claims)["user_id"].(string)
		if !ok {
			log.Errorf("%s: failed to get user_id", op)
			c.JSON(401, gin.H{
				"error": "failed to get user_id",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDstr)
		if err != nil {
			log.Errorf("%s: failed to parse user_id to uuid", op)
			c.JSON(401, gin.H{
				"error": "failed to parse user_id to uuid",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
