package auth

import (
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
)

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
