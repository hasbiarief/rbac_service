package auth

import (
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler struct
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Handler methods

// @Summary      Login dengan user identity
// @Description  Autentikasi pengguna menggunakan user_identity atau email dan mengembalikan access token dan refresh token
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginRequest  true  "Login credentials (user_identity atau email dengan password)"
// @Success      200          {object}  response.Response{data=auth.LoginResponse}  "Login berhasil"
// @Failure      400          {object}  response.Response  "Bad request - format request tidak valid"
// @Failure      401          {object}  response.Response  "Unauthorized - kredensial tidak valid"
// @Router       /api/v1/auth/login [post]
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

	// Get UserAgent and IP from request context
	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	authResponse, err := h.service.Login(req, userAgent, ip)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

// @Summary      Login dengan email
// @Description  Autentikasi pengguna menggunakan email dan password, mengembalikan access token dan refresh token
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginEmailRequest  true  "Login credentials dengan email"
// @Success      200          {object}  response.Response{data=auth.LoginResponse}  "Login berhasil"
// @Failure      400          {object}  response.Response  "Bad request - format request tidak valid"
// @Failure      401          {object}  response.Response  "Unauthorized - kredensial tidak valid"
// @Router       /api/v1/auth/login-email [post]
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

	// Get UserAgent and IP from request context
	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	authResponse, err := h.service.LoginWithEmail(req, userAgent, ip)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, constants.MsgLoginSuccess, authResponse)
}

// @Summary      Refresh access token
// @Description  Memperbarui access token menggunakan refresh token yang valid
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        token  body      auth.RefreshTokenRequest  true  "Refresh token"
// @Success      200    {object}  response.Response{data=auth.RefreshTokenResponse}  "Token berhasil diperbarui"
// @Failure      400    {object}  response.Response  "Bad request - format request tidak valid"
// @Failure      401    {object}  response.Response  "Unauthorized - refresh token tidak valid atau expired"
// @Router       /api/v1/auth/refresh [post]
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

// @Summary      Logout user
// @Description  Menghapus session pengguna berdasarkan token atau user_id
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        logout  body      auth.LogoutRequest  true  "Logout request (token atau user_id)"
// @Success      200     {object}  response.Response  "Logout berhasil"
// @Failure      400     {object}  response.Response  "Bad request - token atau user_id harus disediakan"
// @Router       /api/v1/auth/logout [post]
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

// @Summary      Check user tokens
// @Description  Memeriksa status token pengguna berdasarkan user_identity
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        user_identity  query     string  true  "User identity"
// @Success      200            {object}  response.Response  "Token check berhasil"
// @Failure      400            {object}  response.Response  "Bad request - user_identity diperlukan"
// @Failure      404            {object}  response.Response  "User tidak ditemukan"
// @Failure      500            {object}  response.Response  "Internal server error"
// @Router       /api/v1/auth/check-tokens [get]
func (h *Handler) CheckUserTokens(c *gin.Context) {
	userIdentity := c.Query("user_identity")
	if userIdentity == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_identity is required")
		return
	}

	// Get user by user_identity to get the user ID
	user, err := h.service.repo.GetByUserIdentity(userIdentity)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	tokensResponse, err := h.service.CheckUserTokens(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token check completed", tokensResponse)
}

// @Summary      Get user session count
// @Description  Mendapatkan jumlah session aktif pengguna berdasarkan user_identity
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        user_identity  query     string  true  "User identity"
// @Success      200            {object}  response.Response{data=map[string]interface{}}  "Session count berhasil diambil"
// @Failure      400            {object}  response.Response  "Bad request - user_identity diperlukan"
// @Failure      404            {object}  response.Response  "User tidak ditemukan"
// @Failure      500            {object}  response.Response  "Internal server error"
// @Router       /api/v1/auth/session-count [get]
func (h *Handler) GetUserSessionCount(c *gin.Context) {
	userIdentity := c.Query("user_identity")
	if userIdentity == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_identity is required")
		return
	}

	// Get user by user_identity to get the user ID
	user, err := h.service.repo.GetByUserIdentity(userIdentity)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	accessCount, err := h.service.GetUserSessionCount(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	refreshCount, err := h.service.GetUserRefreshTokenCount(user.ID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Session count retrieved", map[string]interface{}{
		"user_identity":       userIdentity,
		"user_id":             user.ID,
		"access_token_count":  accessCount,
		"refresh_token_count": refreshCount,
		"total_sessions":      accessCount,
	})
}

// @Summary      Cleanup expired tokens
// @Description  Membersihkan token yang sudah expired dari database
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response  "Cleanup berhasil"
// @Failure      500  {object}  response.Response  "Internal server error"
// @Router       /api/v1/auth/cleanup-expired [post]
func (h *Handler) CleanupExpiredTokens(c *gin.Context) {
	err := h.service.CleanupExpiredTokens()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Cleanup failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Expired tokens cleaned up successfully", nil)
}

// @Summary      Get user profile
// @Description  Mendapatkan profil pengguna berdasarkan user_identity dan application_code
// @Tags         🔐 Authentication
// @Accept       json
// @Produce      json
// @Param        user_identity     query     string  true  "User identity"
// @Param        application_code  query     string  true  "Application code"
// @Success      200               {object}  response.Response{data=auth.ProfileResponse}  "Profile berhasil diambil"
// @Failure      400               {object}  response.Response  "Bad request - parameter diperlukan"
// @Failure      404               {object}  response.Response  "User atau application tidak ditemukan"
// @Router       /api/v1/auth/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userIdentity := c.Query("user_identity")
	applicationCode := c.Query("application_code")

	if userIdentity == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "user_identity is required")
		return
	}

	if applicationCode == "" {
		response.Error(c, http.StatusBadRequest, "Bad request", "application_code is required")
		return
	}

	req := &ProfileRequest{
		UserIdentity:    userIdentity,
		ApplicationCode: applicationCode,
	}

	if err := ValidateProfileRequest(req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	profileResponse, err := h.service.GetUserProfile(req)
	if err != nil {
		response.ErrorWithAutoStatus(c, "Failed to get profile", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Profile successfully retrieved", profileResponse)
}

// Route registration
func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	auth := api.Group("/auth")
	{
		// POST /api/v1/auth/login - Login with user_identity
		auth.POST("/login",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LoginRequest{},
			}),
			handler.Login,
		)

		// POST /api/v1/auth/login-email - Login with email
		auth.POST("/login-email",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LoginEmailRequest{},
			}),
			handler.LoginWithEmail,
		)

		// POST /api/v1/auth/refresh - Refresh access token
		auth.POST("/refresh",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &RefreshTokenRequest{},
			}),
			handler.RefreshToken,
		)

		// POST /api/v1/auth/logout - Logout user
		auth.POST("/logout",
			middleware.ValidateRequest(middleware.ValidationRules{
				Body: &LogoutRequest{},
			}),
			handler.Logout,
		)

		// GET /api/v1/auth/check-tokens - Check user tokens
		auth.GET("/check-tokens", handler.CheckUserTokens)

		// GET /api/v1/auth/session-count - Get user session count
		auth.GET("/session-count", handler.GetUserSessionCount)

		// POST /api/v1/auth/cleanup-expired - Cleanup expired tokens
		auth.POST("/cleanup-expired", handler.CleanupExpiredTokens)

		// GET /api/v1/auth/profile - Get user profile by application
		auth.GET("/profile", handler.GetProfile)
	}
}
