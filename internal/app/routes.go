package app

import (
	"gin-scalable-api/middleware"

	// Module imports
	auditModule "gin-scalable-api/internal/modules/audit"
	authModule "gin-scalable-api/internal/modules/auth"
	branchModule "gin-scalable-api/internal/modules/branch"
	companyModule "gin-scalable-api/internal/modules/company"
	moduleModule "gin-scalable-api/internal/modules/module"
	roleModule "gin-scalable-api/internal/modules/role"
	subscriptionModule "gin-scalable-api/internal/modules/subscription"
	unitModule "gin-scalable-api/internal/modules/unit"
	userModule "gin-scalable-api/internal/modules/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetupNewModuleRoutes sets up routes using module-based structure
func SetupNewModuleRoutes(r *gin.Engine, h *NewModuleHandlers, jwtSecret string, redis *redis.Client) {
	// Global middleware
	r.Use(middleware.SmartRateLimit())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")

	// Auth module routes (public)
	authModule.RegisterRoutes(api, h.Auth)

	// Subscription plans (public)
	subscriptionModule.RegisterRoutes(api, h.Subscription)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
	{
		// Register all module routes
		userModule.RegisterRoutes(protected, h.User)
		roleModule.RegisterRoutes(protected, h.Role)
		companyModule.RegisterRoutes(protected, h.Company)
		branchModule.RegisterRoutes(protected, h.Branch)
		moduleModule.RegisterRoutes(protected, h.Module)
		unitModule.RegisterRoutes(protected, h.Unit)
		auditModule.RegisterRoutes(protected, h.Audit)
	}
}
