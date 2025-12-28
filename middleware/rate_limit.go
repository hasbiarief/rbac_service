package middleware

import (
	"net/http"

	"gin-scalable-api/pkg/ratelimiter"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func RateLimit() gin.HandlerFunc {
	limiter := ratelimiter.NewLimiter()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !limiter.Allow(clientIP) {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}