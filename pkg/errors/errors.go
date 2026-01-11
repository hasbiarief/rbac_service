package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors with HTTP status codes
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Predefined error constructors
func NewBadRequestError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: detail,
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewConflictError(message, details string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NewInternalServerError(details string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Details: details,
	}
}

func NewValidationError(details string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
		Details: details,
	}
}

// Business logic specific errors
func NewBranchHasChildrenError() *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: "Cannot delete branch with children",
	}
}

func NewParentBranchNotFoundError() *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: "Parent branch not found",
	}
}
