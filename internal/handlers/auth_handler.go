package handlers

import (
	"fmt"
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService interfaces.AuthServiceInterface
}

func NewAuthHandler(authService interfaces.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.LoginRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Validate that either email or user_identity is provided
	if req.Email == "" && req.UserIdentity == "" {
		response.Error(c, http.StatusBadRequest, "Invalid request", "either email or user_identity must be provided")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	// Set user agent and IP if not provided
	if req.UserAgent == nil {
		req.UserAgent = &userAgent
	}
	if req.IP == nil {
		req.IP = &ip
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

func (h *AuthHandler) LoginWithEmail(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.LoginRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Validate that email is provided for this endpoint
	if req.Email == "" {
		response.Error(c, http.StatusBadRequest, "Invalid request", "email must be provided")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	// Set user agent and IP if not provided
	if req.UserAgent == nil {
		req.UserAgent = &userAgent
	}
	if req.IP == nil {
		req.IP = &ip
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.RefreshTokenRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	authResponse, err := h.authService.RefreshToken(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgTokenRefreshed, authResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected DTO struct
	req, ok := validatedBody.(*dto.LogoutRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Validate that either token or user_id is provided
	if req.Token == "" && req.UserID == 0 {
		response.Error(c, http.StatusBadRequest, "Bad request", "either token or user_id is required")
		return
	}

	var err error
	if req.Token != "" {
		// Logout using token (preferred method)
		err = h.authService.Logout(req.Token)
	} else {
		// Logout using user_id (fallback when token is expired)
		err = h.authService.LogoutByUserID(req.UserID)
	}

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Logout failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLogoutSuccess, nil)
}

func (h *AuthHandler) CheckUserTokens(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_id is required")
		return
	}

	// Convert to int64
	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid user_id format")
		return
	}

	tokensResponse, err := h.authService.CheckUserTokens(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token check completed", tokensResponse)
}

func (h *AuthHandler) GetUserSessionCount(c *gin.Context) {
	// Get user ID from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_id is required")
		return
	}

	// Convert to int64
	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid user_id format")
		return
	}

	accessCount, err := h.authService.GetUserSessionCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	refreshCount, err := h.authService.GetUserRefreshTokenCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Session count retrieved", map[string]interface{}{
		"user_id":             userID,
		"access_token_count":  accessCount,
		"refresh_token_count": refreshCount,
		"total_sessions":      accessCount, // Access tokens represent active sessions
	})
}

func (h *AuthHandler) CleanupExpiredTokens(c *gin.Context) {
	err := h.authService.CleanupExpiredTokens()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Cleanup failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expired tokens cleaned up successfully", nil)
}
