package routes

import (
	"gin-scalable-api/internal/handlers"
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
	// Apply rate limiting globally
	r.Use(middleware.RateLimit())

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
		// Login with validation
		loginValidation := middleware.ValidationRules{
			Body: &struct {
				Email    string `json:"email" validate:"required,email"`
				Password string `json:"password" validate:"required,min=6"`
			}{},
		}
		auth.POST("/login", middleware.ValidateRequest(loginValidation), authHandler.Login)

		// Refresh token with validation
		refreshValidation := middleware.ValidationRules{
			Body: &struct {
				RefreshToken string `json:"refresh_token" validate:"required"`
			}{},
		}
		auth.POST("/refresh", middleware.ValidateRequest(refreshValidation), authHandler.RefreshToken)

		// Logout with validation
		logoutValidation := middleware.ValidationRules{
			Body: &struct {
				Token string `json:"token" validate:"required"`
			}{},
		}
		auth.POST("/logout", middleware.ValidateRequest(logoutValidation), authHandler.Logout)
	}
}

func setupModuleRoutes(protected *gin.RouterGroup, moduleHandler *handlers.ModuleHandler) {
	modules := protected.Group("/modules")
	{
		// List modules with pagination validation
		listValidation := middleware.ValidationRules{
			Query: []middleware.QueryValidation{
				{Name: "page", Type: "int", Default: 1, Min: intPtr(1)},
				{Name: "limit", Type: "int", Default: 10, Min: intPtr(1), Max: intPtr(100)},
				{Name: "search", Type: "string"},
				{Name: "category", Type: "string"},
			},
		}
		modules.GET("", middleware.ValidateRequest(listValidation), moduleHandler.GetModules)

		// Get module by ID with ID validation
		idValidation := middleware.ValidationRules{
			Params: []middleware.ParamValidation{
				{Name: "id", Type: "int", Required: true, Min: intPtr(1)},
			},
		}
		modules.GET("/:id", middleware.ValidateRequest(idValidation), moduleHandler.GetModuleByID)

		// Create module with body validation
		createValidation := middleware.ValidationRules{
			Body: &struct {
				Category         string `json:"category" validate:"required,min=2,max=50"`
				Name             string `json:"name" validate:"required,min=2,max=100"`
				URL              string `json:"url" validate:"required,min=1,max=255"`
				Icon             string `json:"icon" validate:"max=50"`
				Description      string `json:"description" validate:"max=500"`
				ParentID         *int64 `json:"parent_id"`
				SubscriptionTier string `json:"subscription_tier" validate:"required,oneof=basic pro enterprise"`
			}{},
		}
		modules.POST("", middleware.ValidateRequest(createValidation), moduleHandler.CreateModule)

		modules.PUT("/:id", middleware.ValidateRequest(idValidation), moduleHandler.UpdateModule)
		modules.DELETE("/:id", middleware.ValidateRequest(idValidation), moduleHandler.DeleteModule)
		modules.GET("/tree", moduleHandler.GetModuleTree)
		modules.GET("/:id/children", middleware.ValidateRequest(idValidation), moduleHandler.GetModuleChildren)
		modules.GET("/:id/ancestors", middleware.ValidateRequest(idValidation), moduleHandler.GetModuleAncestors)
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
		listValidation := middleware.ValidationRules{
			Query: []middleware.QueryValidation{
				{Name: "page", Type: "int", Default: 1, Min: intPtr(1)},
				{Name: "limit", Type: "int", Default: 10, Min: intPtr(1), Max: intPtr(100)},
				{Name: "search", Type: "string"},
				{Name: "is_active", Type: "bool"},
			},
		}
		users.GET("", middleware.ValidateRequest(listValidation), userHandler.GetUsers)

		// ID validation for single user operations
		idValidation := middleware.ValidationRules{
			Params: []middleware.ParamValidation{
				{Name: "id", Type: "int", Required: true, Min: intPtr(1)},
			},
		}
		users.GET("/:id", middleware.ValidateRequest(idValidation), userHandler.GetUserByID)

		// Create user with validation
		createUserValidation := middleware.ValidationRules{
			Body: &struct {
				Name         string  `json:"name" validate:"required,min=2,max=100"`
				Email        string  `json:"email" validate:"required,email,max=255"`
				UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
				Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
			}{},
		}
		users.POST("", middleware.ValidateRequest(createUserValidation), userHandler.CreateUser)

		users.PUT("/:id", middleware.ValidateRequest(idValidation), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.ValidateRequest(idValidation), userHandler.DeleteUser)

		// User module access
		users.GET("/:id/modules", middleware.ValidateRequest(idValidation), userHandler.GetUserModules)

		// Identity validation
		identityValidation := middleware.ValidationRules{
			Params: []middleware.ParamValidation{
				{Name: "identity", Type: "string", Required: true, Min: intPtr(3), Max: intPtr(50)},
			},
		}
		users.GET("/identity/:identity/modules", middleware.ValidateRequest(identityValidation), userHandler.GetUserModulesByIdentity)

		// Access check validation
		accessCheckValidation := middleware.ValidationRules{
			Body: &struct {
				UserID   int64 `json:"user_id" validate:"required,min=1"`
				ModuleID int64 `json:"module_id" validate:"required,min=1"`
			}{},
		}
		users.POST("/check-access", middleware.ValidateRequest(accessCheckValidation), userHandler.CheckAccess)

		// Password change validation
		passwordChangeValidation := middleware.ValidationRules{
			Params: []middleware.ParamValidation{
				{Name: "id", Type: "int", Required: true, Min: intPtr(1)},
			},
			Body: &struct {
				CurrentPassword string `json:"current_password" validate:"required,min=6"`
				NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
				ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
			}{},
		}
		users.PUT("/:id/password", middleware.ValidateRequest(passwordChangeValidation), userHandler.ChangeUserPassword)
		users.PUT("/me/password", middleware.ValidateRequest(middleware.ValidationRules{
			Body: passwordChangeValidation.Body,
		}), userHandler.ChangeMyPassword)
	}
}

func setupCompanyRoutes(protected *gin.RouterGroup, companyHandler *handlers.CompanyHandler) {
	companies := protected.Group("/companies")
	{
		companies.GET("", companyHandler.GetCompanies)
		companies.GET("/:id", companyHandler.GetCompanyByID)
		companies.POST("", companyHandler.CreateCompany)
		companies.PUT("/:id", companyHandler.UpdateCompany)
		companies.DELETE("/:id", companyHandler.DeleteCompany)
	}
}

func setupRoleRoutes(protected *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	// Basic role management
	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.GetRoles)
		roles.GET("/:id", roleHandler.GetRoleByID)
		roles.POST("", roleHandler.CreateRole)
		roles.PUT("/:id", roleHandler.UpdateRole)
		roles.DELETE("/:id", roleHandler.DeleteRole)
	}

	// Advanced role management system
	roleManagement := protected.Group("/role-management")
	{
		roleManagement.POST("/assign-user-role", roleHandler.AssignUserRole)
		roleManagement.POST("/bulk-assign-roles", roleHandler.BulkAssignRoles)
		roleManagement.PUT("/role/:roleId/modules", roleHandler.UpdateRoleModules)
		roleManagement.DELETE("/user/:userId/role/:roleId", roleHandler.RemoveUserRole)
		roleManagement.GET("/role/:roleId/users", roleHandler.GetUsersByRole)
		roleManagement.GET("/user/:userId/roles", roleHandler.GetUserRoles)
		roleManagement.GET("/user/:userId/access-summary", roleHandler.GetUserAccessSummary)
	}
}

func setupPublicSubscriptionRoutes(api *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	subscription := api.Group("/subscription")
	{
		// Public endpoints (no auth required)
		subscription.GET("/plans", subscriptionHandler.GetAllPlans)
		subscription.GET("/plans/:id", subscriptionHandler.GetPlanByID)
	}
}

func setupProtectedSubscriptionRoutes(protected *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	subscription := protected.Group("/subscription")
	{
		// Subscription management
		subscription.GET("/subscriptions", subscriptionHandler.GetAllSubscriptions)
		subscription.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		subscription.GET("/subscriptions/:id", subscriptionHandler.GetSubscriptionByID)
		subscription.PUT("/subscriptions/:id", subscriptionHandler.UpdateSubscription)
		subscription.POST("/subscriptions/:id/renew", subscriptionHandler.RenewSubscription)
		subscription.POST("/subscriptions/:id/cancel", subscriptionHandler.CancelSubscription)

		// Company subscription management
		subscription.GET("/companies/:id/subscription", subscriptionHandler.GetCompanySubscription)
		subscription.GET("/companies/:id/status", subscriptionHandler.GetCompanySubscriptionStatus)

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
		audit.POST("/logs", auditHandler.CreateAuditLog)
		audit.GET("/users/:userId/logs", auditHandler.GetUserAuditLogs)
		audit.GET("/users/identity/:identity/logs", auditHandler.GetUserAuditLogsByIdentity)
		audit.GET("/stats", auditHandler.GetAuditStats)
	}
}

func setupBranchRoutes(protected *gin.RouterGroup, branchHandler *handlers.BranchHandler) {
	branches := protected.Group("/branches")
	{
		branches.GET("", branchHandler.GetBranches)
		branches.GET("/:id", branchHandler.GetBranchByID)
		branches.POST("", branchHandler.CreateBranch)
		branches.PUT("/:id", branchHandler.UpdateBranch)
		branches.DELETE("/:id", branchHandler.DeleteBranch)
		branches.GET("/company/:companyId", branchHandler.GetCompanyBranches)
		branches.GET("/:id/children", branchHandler.GetBranchChildren)
	}
}
