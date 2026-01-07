package middleware

import (
	"net/http"
	"strings"

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

// RateLimitForCheckAccess creates a more lenient rate limiter for check-access endpoint
func RateLimitForCheckAccess() gin.HandlerFunc {
	// More lenient: 30 requests per second with burst of 100
	limiter := ratelimiter.NewLimiterWithConfig(30, 100)

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

// SmartRateLimit applies different rate limits based on endpoint
func SmartRateLimit() gin.HandlerFunc {
	defaultLimiter := ratelimiter.NewLimiter()                      // 10 req/sec, burst 50
	checkAccessLimiter := ratelimiter.NewLimiterWithConfig(30, 100) // 30 req/sec, burst 100

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.Request.URL.Path

		var limiter *ratelimiter.IPRateLimiter

		// Use more lenient rate limiting for check-access endpoint
		if strings.Contains(path, "/check-access") {
			limiter = checkAccessLimiter
		} else {
			limiter = defaultLimiter
		}

		if !limiter.Allow(clientIP) {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}
