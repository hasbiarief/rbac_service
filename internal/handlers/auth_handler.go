package handlers

import (
	"fmt"
	"gin-scalable-api/internal/service"
	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
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

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		UserIdentity string `json:"user_identity" validate:"required"`
		Password     string `json:"password" validate:"required,min=6"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Convert to service request
	loginReq := &service.LoginRequest{
		UserIdentity: req.UserIdentity,
		Password:     req.Password,
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	authResponse, err := h.authService.Login(loginReq, userAgent, ip)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", authResponse)
}

func (h *AuthHandler) LoginWithEmail(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	// Convert to service request - use email as user_identity
	loginReq := &service.LoginRequest{
		UserIdentity: req.Email,
		Password:     req.Password,
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	authResponse, err := h.authService.Login(loginReq, userAgent, ip)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", authResponse)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected struct
	req, ok := validatedBody.(*struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	})
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	authResponse, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", authResponse)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Get validated body from context (set by validation middleware)
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	// Type assert to the expected struct - support both token and user_id
	req, ok := validatedBody.(*struct {
		Token  string `json:"token"`
		UserID int64  `json:"user_id"`
	})
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

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
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
