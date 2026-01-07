package routes

import (
	"gin-scalable-api/internal/handlers"
	"gin-scalable-api/internal/validation"
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handlers struct {
	Auth         *handlers.AuthHandler
	Module       *handlers.ModuleHandler
	User         *handlers.UserHandler
	Company      *handlers.CompanyHandler
	Role         *handlers.RoleHandler
	Subscription *handlers.SubscriptionHandler
	Audit        *handlers.AuditHandler
	Branch       *handlers.BranchHandler
}

func SetupRoutes(r *gin.Engine, h *Handlers, jwtSecret string, redis *redis.Client) {
	// Apply smart rate limiting globally (different limits for different endpoints)
	r.Use(middleware.SmartRateLimit())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		setupAuthRoutes(api, h.Auth)

		// Public subscription routes
		setupPublicSubscriptionRoutes(api, h.Subscription)

		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
		{
			// Module routes
			setupModuleRoutes(protected, h.Module)

			// User routes (including user module access and password management)
			setupUserRoutes(protected, h.User)

			// Company routes
			setupCompanyRoutes(protected, h.Company)

			// Role routes (basic and advanced)
			setupRoleRoutes(protected, h.Role)

			// Protected subscription routes
			setupProtectedSubscriptionRoutes(protected, h.Subscription)

			// Audit routes
			setupAuditRoutes(protected, h.Audit)

			// Branch routes
			setupBranchRoutes(protected, h.Branch)
		}
	}
}

func setupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		// Login with validation - support both user_identity and email
		auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)

		// Login with email (alternative endpoint)
		auth.POST("/login-email", middleware.ValidateRequest(validation.LoginEmailValidation), authHandler.LoginWithEmail)

		// Refresh token with validation
		auth.POST("/refresh", middleware.ValidateRequest(validation.RefreshValidation), authHandler.RefreshToken)

		// Logout with validation
		auth.POST("/logout", middleware.ValidateRequest(validation.LogoutValidation), authHandler.Logout)

		// Check user tokens (for frontend token validation)
		auth.GET("/check-tokens", authHandler.CheckUserTokens)

		// Admin endpoints for session management
		auth.GET("/session-count", authHandler.GetUserSessionCount)
		auth.POST("/cleanup-expired", authHandler.CleanupExpiredTokens)
	}
}

func setupModuleRoutes(protected *gin.RouterGroup, moduleHandler *handlers.ModuleHandler) {
	modules := protected.Group("/modules")
	{
		// List modules with pagination validation
		modules.GET("", middleware.ValidateRequest(validation.ModuleListValidation), moduleHandler.GetModules)

		// Get module by ID with ID validation
		modules.GET("/:id", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleByID)

		// Create module with body validation
		modules.POST("", middleware.ValidateRequest(validation.CreateModuleValidation), moduleHandler.CreateModule)

		// Update module with validation
		modules.PUT("/:id", middleware.ValidateRequest(validation.UpdateModuleValidation), moduleHandler.UpdateModule)
		modules.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), moduleHandler.DeleteModule)
		modules.GET("/tree", moduleHandler.GetModuleTree)
		modules.GET("/:id/children", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleChildren)
		modules.GET("/:id/ancestors", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleAncestors)
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		// List users with pagination validation
		users.GET("", middleware.ValidateRequest(validation.UserListValidation), userHandler.GetUsers)

		// ID validation for single user operations
		users.GET("/:id", middleware.ValidateRequest(validation.IDValidation), userHandler.GetUserByID)

		// Create user with validation
		users.POST("", middleware.ValidateRequest(validation.CreateUserValidation), userHandler.CreateUser)

		// Update user with validation
		users.PUT("/:id", middleware.ValidateRequest(middleware.ValidationRules{
			Params: validation.IDValidation.Params,
			Body:   validation.UpdateUserValidation.Body,
		}), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), userHandler.DeleteUser)

		// User module access
		users.GET("/:id/modules", middleware.ValidateRequest(validation.IDValidation), userHandler.GetUserModules)

		// Identity validation
		users.GET("/identity/:identity/modules", middleware.ValidateRequest(validation.IdentityValidation), userHandler.GetUserModulesByIdentity)

		// Access check validation
		users.POST("/check-access", middleware.ValidateRequest(validation.AccessCheckValidation), userHandler.CheckAccess)

		// Password change validation
		users.PUT("/:id/password", middleware.ValidateRequest(validation.PasswordChangeValidation), userHandler.ChangeUserPassword)
		// users.PUT("/me/password", middleware.ValidateRequest(validation.MyPasswordChangeValidation), userHandler.ChangeMyPassword)
	}
}

func setupCompanyRoutes(protected *gin.RouterGroup, companyHandler *handlers.CompanyHandler) {
	companies := protected.Group("/companies")
	{
		companies.GET("", companyHandler.GetCompanies)

		// ID validation for single company operations
		companies.GET("/:id", middleware.ValidateRequest(validation.IDValidation), companyHandler.GetCompanyByID)

		// Create company with validation
		companies.POST("", middleware.ValidateRequest(validation.CreateCompanyValidation), companyHandler.CreateCompany)

		// Update company with validation
		companies.PUT("/:id", middleware.ValidateRequest(validation.UpdateCompanyValidation), companyHandler.UpdateCompany)
		companies.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), companyHandler.DeleteCompany)
	}
}

func setupRoleRoutes(protected *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	// Basic role management
	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.GetRoles)

		// ID validation for single role operations
		roles.GET("/:id", middleware.ValidateRequest(validation.IDValidation), roleHandler.GetRoleByID)

		// Create role with validation
		roles.POST("", middleware.ValidateRequest(validation.CreateRoleValidation), roleHandler.CreateRole)

		// Update role with validation
		roles.PUT("/:id", middleware.ValidateRequest(validation.UpdateRoleValidation), roleHandler.UpdateRole)
		roles.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), roleHandler.DeleteRole)
	}

	// Advanced role management system
	roleManagement := protected.Group("/role-management")
	{
		// Assign user role validation
		roleManagement.POST("/assign-user-role", middleware.ValidateRequest(validation.AssignUserRoleValidation), roleHandler.AssignUserRole)

		// Bulk assign roles validation
		roleManagement.POST("/bulk-assign-roles", middleware.ValidateRequest(validation.BulkAssignRolesValidation), roleHandler.BulkAssignRoles)

		// Update role modules validation
		roleManagement.PUT("/role/:roleId/modules", middleware.ValidateRequest(validation.UpdateRoleModulesValidation), roleHandler.UpdateRoleModules)

		// User/Role operations
		roleManagement.DELETE("/user/:userId/role/:roleId", middleware.ValidateRequest(validation.RoleUserValidation), roleHandler.RemoveUserRole)
		roleManagement.GET("/role/:roleId/users", middleware.ValidateRequest(validation.RoleIDValidation), roleHandler.GetUsersByRole)
		roleManagement.GET("/user/:userId/roles", middleware.ValidateRequest(validation.UserIDValidation), roleHandler.GetUserRoles)
		roleManagement.GET("/user/:userId/access-summary", middleware.ValidateRequest(validation.UserIDValidation), roleHandler.GetUserAccessSummary)
	}
}

func setupPublicSubscriptionRoutes(api *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	// Public endpoints (no auth required) - use different path to avoid conflicts
	plans := api.Group("/plans")
	{
		plans.GET("", subscriptionHandler.GetAllPlans)
		plans.GET("/:id", subscriptionHandler.GetPlanByID)
	}
}

func setupProtectedSubscriptionRoutes(protected *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	subscription := protected.Group("/subscription")
	{
		// Subscription management
		subscription.GET("/subscriptions", subscriptionHandler.GetAllSubscriptions)

		// Create subscription with validation
		subscription.POST("/subscriptions", middleware.ValidateRequest(validation.CreateSubscriptionValidation), subscriptionHandler.CreateSubscription)

		// ID validation for subscription operations
		subscription.GET("/subscriptions/:id", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetSubscriptionByID)
		subscription.PUT("/subscriptions/:id", middleware.ValidateRequest(validation.UpdateSubscriptionValidation), subscriptionHandler.UpdateSubscription)
		subscription.POST("/subscriptions/:id/renew", middleware.ValidateRequest(validation.RenewSubscriptionValidation), subscriptionHandler.RenewSubscription)
		subscription.POST("/subscriptions/:id/cancel", middleware.ValidateRequest(validation.CancelSubscriptionValidation), subscriptionHandler.CancelSubscription)
		subscription.POST("/subscriptions/:id/mark-paid", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.MarkPaymentAsPaid)

		// Company subscription management
		subscription.GET("/companies/:id/subscription", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetCompanySubscription)
		subscription.GET("/companies/:id/status", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetCompanySubscriptionStatus)

		// Module access check - use different path to avoid conflict
		subscription.GET("/module-access/:companyId/:moduleId", subscriptionHandler.CheckModuleAccess)

		// Subscription utilities
		subscription.GET("/stats", subscriptionHandler.GetSubscriptionStats)
		subscription.GET("/expiring", subscriptionHandler.GetExpiringSubscriptions)
		subscription.POST("/update-expired", subscriptionHandler.UpdateExpiredSubscriptions)
	}
}

func setupAuditRoutes(protected *gin.RouterGroup, auditHandler *handlers.AuditHandler) {
	audit := protected.Group("/audit")
	{
		audit.GET("/logs", auditHandler.GetAuditLogs)

		// Create audit log
		audit.POST("/logs", auditHandler.CreateAuditLog)

		// User ID validation
		audit.GET("/users/:userId/logs", middleware.ValidateRequest(validation.UserIDValidation), auditHandler.GetUserAuditLogs)

		// Identity validation
		audit.GET("/users/identity/:identity/logs", middleware.ValidateRequest(validation.IdentityValidation), auditHandler.GetUserAuditLogsByIdentity)
		audit.GET("/stats", auditHandler.GetAuditStats)
	}
}

func setupBranchRoutes(protected *gin.RouterGroup, branchHandler *handlers.BranchHandler) {
	branches := protected.Group("/branches")
	{
		branches.GET("", branchHandler.GetBranches)

		// ID validation for single branch operations
		branches.GET("/:id", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchByID)
		branches.GET("/:id/hierarchy", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchHierarchy)
		branches.GET("/:id/children", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchChildren)

		// Create branch with validation
		branches.POST("", middleware.ValidateRequest(validation.CreateBranchValidation), branchHandler.CreateBranch)

		// Update branch with validation
		branches.PUT("/:id", middleware.ValidateRequest(validation.UpdateBranchValidation), branchHandler.UpdateBranch)
		branches.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), branchHandler.DeleteBranch)

		// Company ID validation
		branches.GET("/company/:companyId", middleware.ValidateRequest(validation.CompanyIDValidation), branchHandler.GetCompanyBranches)
	}
}
