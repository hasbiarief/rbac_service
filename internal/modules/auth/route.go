package auth

import (
	"fmt"
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*LoginRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	if req.UserIdentity == "" {
		response.Error(c, http.StatusBadRequest, "Invalid request", "user_identity must be provided")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	if req.UserAgent == nil {
		req.UserAgent = &userAgent
	}
	if req.IP == nil {
		req.IP = &ip
	}

	authResponse, err := h.service.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

func (h *Handler) LoginWithEmail(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*LoginEmailRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	if req.UserAgent == nil {
		req.UserAgent = &userAgent
	}
	if req.IP == nil {
		req.IP = &ip
	}

	loginReq := &LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: req.UserAgent,
		IP:        req.IP,
	}

	authResponse, err := h.service.Login(loginReq)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*RefreshTokenRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	authResponse, err := h.service.RefreshToken(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Token refresh failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgTokenRefreshed, authResponse)
}

func (h *Handler) Logout(c *gin.Context) {
	validatedBody, exists := c.Get("validated_body")
	if !exists {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "validation failed")
		return
	}

	req, ok := validatedBody.(*LogoutRequest)
	if !ok {
		response.Error(c, http.StatusBadRequest, "Invalid request format", "invalid body structure")
		return
	}

	if req.Token == "" && req.UserID == 0 {
		response.Error(c, http.StatusBadRequest, "Bad request", "either token or user_id is required")
		return
	}

	var err error
	if req.Token != "" {
		err = h.service.Logout(req.Token)
	} else {
		err = h.service.LogoutByUserID(req.UserID)
	}

	if err != nil {
		response.Error(c, http.StatusBadRequest, "Logout failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLogoutSuccess, nil)
}

func (h *Handler) CheckUserTokens(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_id is required")
		return
	}

	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid user_id format")
		return
	}

	tokensResponse, err := h.service.CheckUserTokens(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token check completed", tokensResponse)
}

func (h *Handler) GetUserSessionCount(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_id is required")
		return
	}

	var userID int64
	if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err != nil {
		response.Error(c, http.StatusBadRequest, "Bad request", "invalid user_id format")
		return
	}

	accessCount, err := h.service.GetUserSessionCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	refreshCount, err := h.service.GetUserRefreshTokenCount(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Session count retrieved", map[string]interface{}{
		"user_id":             userID,
		"access_token_count":  accessCount,
		"refresh_token_count": refreshCount,
		"total_sessions":      accessCount,
	})
}

func (h *Handler) CleanupExpiredTokens(c *gin.Context) {
	err := h.service.CleanupExpiredTokens()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Cleanup failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expired tokens cleaned up successfully", nil)
}

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	auth := api.Group("/auth")
	{
		auth.POST("/login",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LoginRequest{},
			}),
			handler.Login,
		)
		auth.POST("/login-email",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LoginEmailRequest{},
			}),
			handler.LoginWithEmail,
		)
		auth.POST("/refresh",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &RefreshTokenRequest{},
			}),
			handler.RefreshToken,
		)
		auth.POST("/logout",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LogoutRequest{},
			}),
			handler.Logout,
		)
		auth.GET("/check-tokens", handler.CheckUserTokens)
		auth.GET("/session-count", handler.GetUserSessionCount)
		auth.POST("/cleanup-expired", handler.CleanupExpiredTokens)
	}
}
