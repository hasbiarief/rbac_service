package utils

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/errors"
	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseHelper provides utility methods for consistent API responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// Success sends a success response
func (h *ResponseHelper) Success(c *gin.Context, message string, data interface{}) {
	response.Success(c, http.StatusOK, message, data)
}

// Created sends a created response
func (h *ResponseHelper) Created(c *gin.Context, message string, data interface{}) {
	response.Success(c, http.StatusCreated, message, data)
}

// NoContent sends a no content response
func (h *ResponseHelper) NoContent(c *gin.Context, message string) {
	response.Success(c, http.StatusNoContent, message, nil)
}

// HandleError handles different types of errors and sends appropriate responses
func (h *ResponseHelper) HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		response.Error(c, appErr.Code, appErr.Message, appErr.Details)
		return
	}

	// Default to internal server error for unknown errors
	response.Error(c, http.StatusInternalServerError, constants.MsgNotFound, err.Error())
}

// ValidationError sends a validation error response
func (h *ResponseHelper) ValidationError(c *gin.Context, details string) {
	response.Error(c, http.StatusBadRequest, constants.MsgValidationFailed, details)
}

// NotFound sends a not found response
func (h *ResponseHelper) NotFound(c *gin.Context, resource string) {
	response.Error(c, http.StatusNotFound, resource+" not found", "")
}

// Unauthorized sends an unauthorized response
func (h *ResponseHelper) Unauthorized(c *gin.Context, message string) {
	response.Error(c, http.StatusUnauthorized, message, "")
}

// Forbidden sends a forbidden response
func (h *ResponseHelper) Forbidden(c *gin.Context, message string) {
	response.Error(c, http.StatusForbidden, message, "")
}

// Conflict sends a conflict response
func (h *ResponseHelper) Conflict(c *gin.Context, message, details string) {
	response.Error(c, http.StatusConflict, message, details)
}

// BadRequest sends a bad request response
func (h *ResponseHelper) BadRequest(c *gin.Context, message string) {
	response.Error(c, http.StatusBadRequest, message, "")
}

// InternalServerError sends an internal server error response
func (h *ResponseHelper) InternalServerError(c *gin.Context, details string) {
	response.Error(c, http.StatusInternalServerError, "Internal server error", details)
}
