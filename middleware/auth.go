package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gin-scalable-api/pkg/response"
	"gin-scalable-api/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(jwtSecret string, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Add debug information as string
			debugInfo := fmt.Sprintf("header_length=%d, parts_count=%d, raw_header='%s'",
				len(authHeader), len(tokenParts), authHeader)
			if len(tokenParts) > 0 {
				debugInfo += fmt.Sprintf(", first_part='%s'", tokenParts[0])
			}
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", debugInfo)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Initialize token service
		tokenService := token.NewTokenService(redis)

		// Get token metadata from Redis
		metadata, err := tokenService.GetAccessToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Check if token is expired
		if time.Now().Unix() > metadata.ExpiresAt {
			response.Error(c, http.StatusUnauthorized, "Token expired", "")
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", metadata.UserID)
		c.Set("abilities", metadata.Abilities)

		c.Next()
	}
}
