package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, error string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   error,
	})
}

// ErrorWithAutoStatus automatically determines the appropriate HTTP status code based on error message
func ErrorWithAutoStatus(c *gin.Context, message string, error string) {
	statusCode := DetermineStatusCode(error)
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   error,
	})
}

// DetermineStatusCode determines appropriate HTTP status code based on error message
func DetermineStatusCode(errorMsg string) int {
	errorLower := strings.ToLower(errorMsg)

	// 400 Bad Request - Client errors, validation errors
	if strings.Contains(errorLower, "validation failed") ||
		strings.Contains(errorLower, "invalid") ||
		strings.Contains(errorLower, "required") ||
		strings.Contains(errorLower, "bad request") ||
		strings.Contains(errorLower, "malformed") ||
		strings.Contains(errorLower, "missing") {
		return http.StatusBadRequest
	}

	// 401 Unauthorized - Authentication errors
	if strings.Contains(errorLower, "unauthorized") ||
		strings.Contains(errorLower, "invalid credentials") ||
		strings.Contains(errorLower, "authentication failed") ||
		strings.Contains(errorLower, "token") {
		return http.StatusUnauthorized
	}

	// 403 Forbidden - Authorization errors
	if strings.Contains(errorLower, "forbidden") ||
		strings.Contains(errorLower, "access denied") ||
		strings.Contains(errorLower, "permission denied") ||
		strings.Contains(errorLower, "not authorized") {
		return http.StatusForbidden
	}

	// 404 Not Found - Resource not found
	if strings.Contains(errorLower, "not found") ||
		strings.Contains(errorLower, "does not exist") ||
		strings.Contains(errorLower, "no rows") {
		return http.StatusNotFound
	}

	// 409 Conflict - Duplicate entries, constraint violations
	if strings.Contains(errorLower, "already exists") ||
		strings.Contains(errorLower, "duplicate") ||
		strings.Contains(errorLower, "conflict") ||
		strings.Contains(errorLower, "constraint") ||
		strings.Contains(errorLower, "unique") {
		return http.StatusConflict
	}

	// 422 Unprocessable Entity - Business logic errors
	if strings.Contains(errorLower, "cannot") ||
		strings.Contains(errorLower, "unable to") ||
		strings.Contains(errorLower, "failed to parse") ||
		strings.Contains(errorLower, "business rule") {
		return http.StatusUnprocessableEntity
	}

	// Default to 500 Internal Server Error for unexpected errors
	return http.StatusInternalServerError
}
