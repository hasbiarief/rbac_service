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
	Unit         *handlers.UnitHandler
	UnitContext  *handlers.UnitContextHandler
}

func SetupRoutes(r *gin.Engine, h *Handlers, jwtSecret string, redis *redis.Client) {
	// Terapkan pembatasan rate secara global (batas berbeda untuk endpoint berbeda)
	r.Use(middleware.SmartRateLimit())

	// Pemeriksaan kesehatan sistem
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Rute API
	api := r.Group("/api/v1")
	{
		// Rute autentikasi (publik)
		setupAuthRoutes(api, h.Auth)

		// Rute langganan publik
		setupPublicSubscriptionRoutes(api, h.Subscription)

		// Rute yang dilindungi - memerlukan autentikasi
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
		{
			// Rute modul
			setupModuleRoutes(protected, h.Module)

			// Rute pengguna (termasuk akses modul pengguna dan manajemen kata sandi)
			setupUserRoutes(protected, h.User)

			// Rute perusahaan
			setupCompanyRoutes(protected, h.Company)

			// Rute peran (dasar dan lanjutan)
			setupRoleRoutes(protected, h.Role)

			// Rute langganan yang dilindungi
			setupProtectedSubscriptionRoutes(protected, h.Subscription)

			// Rute audit
			setupAuditRoutes(protected, h.Audit)

			// Rute cabang
			setupBranchRoutes(protected, h.Branch)

			// Rute unit
			setupUnitRoutes(protected, h.Unit)

			// Rute unit context
			setupUnitContextRoutes(protected, h.UnitContext)
		}
	}
}

func setupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		// Login dengan validasi - mendukung user_identity dan email
		auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)

		// Login dengan email (endpoint alternatif)
		auth.POST("/login-email", middleware.ValidateRequest(validation.LoginEmailValidation), authHandler.LoginWithEmail)

		// Refresh token dengan validasi
		auth.POST("/refresh", middleware.ValidateRequest(validation.RefreshValidation), authHandler.RefreshToken)

		// Logout dengan validasi
		auth.POST("/logout", middleware.ValidateRequest(validation.LogoutValidation), authHandler.Logout)

		// Periksa token pengguna (untuk validasi token frontend)
		auth.GET("/check-tokens", authHandler.CheckUserTokens)

		// Endpoint admin untuk manajemen sesi
		auth.GET("/session-count", authHandler.GetUserSessionCount)
		auth.POST("/cleanup-expired", authHandler.CleanupExpiredTokens)
	}
}

func setupModuleRoutes(protected *gin.RouterGroup, moduleHandler *handlers.ModuleHandler) {
	modules := protected.Group("/modules")
	{
		// Daftar modul dengan validasi paginasi
		modules.GET("", middleware.ValidateRequest(validation.ModuleListValidation), moduleHandler.GetModules)

		// Dapatkan modul berdasarkan ID dengan validasi ID
		modules.GET("/:id", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleByID)

		// Buat modul dengan validasi body
		modules.POST("", middleware.ValidateRequest(validation.CreateModuleValidation), moduleHandler.CreateModule)

		// Update modul dengan validasi
		modules.PUT("/:id", middleware.ValidateRequest(validation.UpdateModuleValidation), moduleHandler.UpdateModule)
		modules.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), moduleHandler.DeleteModule)
		modules.GET("/tree", moduleHandler.GetModuleTree)
		modules.GET("/:id/children", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleChildren)
		modules.GET("/:id/ancestors", middleware.ValidateRequest(validation.IDValidation), moduleHandler.GetModuleAncestors)
	}
}

// Fungsi helper untuk membuat pointer integer
func intPtr(i int) *int {
	return &i
}

func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	{
		// Daftar pengguna dengan validasi paginasi
		users.GET("", middleware.ValidateRequest(validation.UserListValidation), userHandler.GetUsers)

		// Validasi ID untuk operasi pengguna tunggal
		users.GET("/:id", middleware.ValidateRequest(validation.IDValidation), userHandler.GetUserByID)

		// Buat pengguna dengan validasi
		users.POST("", middleware.ValidateRequest(validation.CreateUserValidation), userHandler.CreateUser)

		// Update pengguna dengan validasi
		users.PUT("/:id", middleware.ValidateRequest(middleware.ValidationRules{
			Params: validation.IDValidation.Params,
			Body:   validation.UpdateUserValidation.Body,
		}), userHandler.UpdateUser)
		users.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), userHandler.DeleteUser)

		// Akses modul pengguna
		users.GET("/:id/modules", middleware.ValidateRequest(validation.IDValidation), userHandler.GetUserModules)

		// Validasi identitas
		users.GET("/identity/:identity/modules", middleware.ValidateRequest(validation.IdentityValidation), userHandler.GetUserModulesByIdentity)

		// Validasi pemeriksaan akses
		users.POST("/check-access", middleware.ValidateRequest(validation.AccessCheckValidation), userHandler.CheckAccess)

		// Validasi perubahan kata sandi
		users.PUT("/:id/password", middleware.ValidateRequest(validation.ChangePasswordValidation), userHandler.ChangeUserPassword)
		// users.PUT("/me/password", middleware.ValidateRequest(validation.MyPasswordChangeValidation), userHandler.ChangeMyPassword)
	}
}

func setupCompanyRoutes(protected *gin.RouterGroup, companyHandler *handlers.CompanyHandler) {
	companies := protected.Group("/companies")
	{
		companies.GET("", companyHandler.GetCompanies)

		// Validasi ID untuk operasi perusahaan tunggal
		companies.GET("/:id", middleware.ValidateRequest(validation.IDValidation), companyHandler.GetCompanyByID)

		// Buat perusahaan dengan validasi
		companies.POST("", middleware.ValidateRequest(validation.CreateCompanyValidation), companyHandler.CreateCompany)

		// Update perusahaan dengan validasi
		companies.PUT("/:id", middleware.ValidateRequest(validation.UpdateCompanyValidation), companyHandler.UpdateCompany)
		companies.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), companyHandler.DeleteCompany)
	}
}

func setupRoleRoutes(protected *gin.RouterGroup, roleHandler *handlers.RoleHandler) {
	// Manajemen peran dasar
	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.GetRoles)

		// Validasi ID untuk operasi peran tunggal
		roles.GET("/:id", middleware.ValidateRequest(validation.IDValidation), roleHandler.GetRoleByID)

		// Buat peran dengan validasi
		roles.POST("", middleware.ValidateRequest(validation.CreateRoleValidation), roleHandler.CreateRole)

		// Update peran dengan validasi
		roles.PUT("/:id", middleware.ValidateRequest(validation.UpdateRoleValidation), roleHandler.UpdateRole)
		roles.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), roleHandler.DeleteRole)
	}

	// Sistem manajemen peran lanjutan
	roleManagement := protected.Group("/role-management")
	{
		// Validasi penugasan peran pengguna
		roleManagement.POST("/assign-user-role", middleware.ValidateRequest(validation.AssignUserRoleValidation), roleHandler.AssignUserRole)

		// Validasi penugasan peran massal - method ini dikomentari karena tidak ada di handler
		// roleManagement.POST("/bulk-assign-roles", middleware.ValidateRequest(validation.BulkAssignRolesValidation), roleHandler.BulkAssignRoles)

		// Validasi update modul peran
		roleManagement.PUT("/role/:roleId/modules", middleware.ValidateRequest(validation.UpdateRoleModulesValidation), roleHandler.UpdateRoleModules)

		// Operasi Pengguna/Peran
		roleManagement.DELETE("/user/:userId/role/:roleId", middleware.ValidateRequest(validation.RoleUserValidation), roleHandler.RemoveUserRole)
		roleManagement.GET("/role/:roleId/users", middleware.ValidateRequest(validation.RoleIDValidation), roleHandler.GetUsersByRole)
		roleManagement.GET("/user/:userId/roles", middleware.ValidateRequest(validation.UserIDValidation), roleHandler.GetUserRoles)
		roleManagement.GET("/user/:userId/access-summary", middleware.ValidateRequest(validation.UserIDValidation), roleHandler.GetUserAccessSummary)
	}
}

func setupPublicSubscriptionRoutes(api *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	// Endpoint publik (tidak memerlukan autentikasi) - gunakan path berbeda untuk menghindari konflik
	plans := api.Group("/plans")
	{
		plans.GET("", subscriptionHandler.GetAllPlans)
		plans.GET("/:id", subscriptionHandler.GetPlanByID)
	}
}

func setupProtectedSubscriptionRoutes(protected *gin.RouterGroup, subscriptionHandler *handlers.SubscriptionHandler) {
	subscription := protected.Group("/subscription")
	{
		// Manajemen langganan
		subscription.GET("/subscriptions", subscriptionHandler.GetAllSubscriptions)

		// Buat langganan dengan validasi
		subscription.POST("/subscriptions", middleware.ValidateRequest(validation.CreateSubscriptionValidation), subscriptionHandler.CreateSubscription)

		// Validasi ID untuk operasi langganan
		subscription.GET("/subscriptions/:id", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetSubscriptionByID)
		subscription.PUT("/subscriptions/:id", middleware.ValidateRequest(validation.UpdateSubscriptionValidation), subscriptionHandler.UpdateSubscription)
		subscription.POST("/subscriptions/:id/renew", middleware.ValidateRequest(validation.RenewSubscriptionValidation), subscriptionHandler.RenewSubscription)
		subscription.POST("/subscriptions/:id/cancel", middleware.ValidateRequest(validation.CancelSubscriptionValidation), subscriptionHandler.CancelSubscription)
		subscription.POST("/subscriptions/:id/mark-paid", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.MarkPaymentAsPaid)

		// Manajemen langganan perusahaan
		subscription.GET("/companies/:id/subscription", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetCompanySubscription)
		subscription.GET("/companies/:id/status", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.GetCompanySubscriptionStatus)

		// Pemeriksaan akses modul - gunakan path berbeda untuk menghindari konflik
		subscription.GET("/module-access/:companyId/:moduleId", subscriptionHandler.CheckModuleAccess)

		// Utilitas langganan
		subscription.GET("/stats", subscriptionHandler.GetSubscriptionStats)
		subscription.GET("/expiring", subscriptionHandler.GetExpiringSubscriptions)
		subscription.POST("/update-expired", subscriptionHandler.UpdateExpiredSubscriptions)
	}

	// Admin routes untuk manajemen subscription plan
	admin := protected.Group("/admin")
	{
		// Manajemen subscription plan (admin only)
		admin.POST("/subscription-plans", middleware.ValidateRequest(validation.CreateSubscriptionPlanValidation), subscriptionHandler.CreateSubscriptionPlan)
		admin.PUT("/subscription-plans/:id", middleware.ValidateRequest(validation.UpdateSubscriptionPlanValidation), subscriptionHandler.UpdateSubscriptionPlan)
		admin.DELETE("/subscription-plans/:id", middleware.ValidateRequest(validation.IDValidation), subscriptionHandler.DeleteSubscriptionPlan)
	}
}

func setupAuditRoutes(protected *gin.RouterGroup, auditHandler *handlers.AuditHandler) {
	audit := protected.Group("/audit")
	{
		audit.GET("/logs", auditHandler.GetAuditLogs)

		// Buat log audit
		audit.POST("/logs", auditHandler.CreateAuditLog)

		// Validasi ID pengguna
		audit.GET("/users/:userId/logs", middleware.ValidateRequest(validation.UserIDValidation), auditHandler.GetUserAuditLogs)

		// Validasi identitas
		audit.GET("/users/identity/:identity/logs", middleware.ValidateRequest(validation.IdentityValidation), auditHandler.GetUserAuditLogsByIdentity)
		audit.GET("/stats", auditHandler.GetAuditStats)
	}
}

func setupBranchRoutes(protected *gin.RouterGroup, branchHandler *handlers.BranchHandler) {
	branches := protected.Group("/branches")
	{
		branches.GET("", branchHandler.GetBranches)

		// Validasi ID untuk operasi cabang tunggal
		branches.GET("/:id", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchByID)
		branches.GET("/:id/hierarchy", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchHierarchy)
		branches.GET("/:id/children", middleware.ValidateRequest(validation.IDValidation), branchHandler.GetBranchChildren)

		// Buat cabang dengan validasi
		branches.POST("", middleware.ValidateRequest(validation.CreateBranchValidation), branchHandler.CreateBranch)

		// Update cabang dengan validasi
		branches.PUT("/:id", middleware.ValidateRequest(validation.UpdateBranchValidation), branchHandler.UpdateBranch)
		branches.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), branchHandler.DeleteBranch)

		// Validasi ID perusahaan
		branches.GET("/company/:companyId", middleware.ValidateRequest(validation.CompanyIDValidation), branchHandler.GetCompanyBranches)
	}
}

func setupUnitRoutes(protected *gin.RouterGroup, unitHandler *handlers.UnitHandler) {
	units := protected.Group("/units")
	{
		// Unit CRUD operations
		units.GET("", middleware.ValidateRequest(validation.UnitListValidation), unitHandler.GetUnits)
		units.POST("", middleware.ValidateRequest(validation.CreateUnitValidation), unitHandler.CreateUnit)
		units.GET("/:id", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitByID)
		units.PUT("/:id", middleware.ValidateRequest(validation.UpdateUnitValidation), unitHandler.UpdateUnit)
		units.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), unitHandler.DeleteUnit)

		// Unit statistics
		units.GET("/:id/stats", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitWithStats)

		// Unit roles management - using :id consistently
		units.GET("/:id/roles", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitRoles)
		units.POST("/:id/roles/:role_id", unitHandler.AssignRoleToUnit)
		units.DELETE("/:id/roles/:role_id", unitHandler.RemoveRoleFromUnit)

		// Unit permissions management - using :id consistently
		units.GET("/:id/roles/:role_id/permissions", unitHandler.GetUnitPermissions)

		// Bulk operations
		units.POST("/copy-permissions", middleware.ValidateRequest(validation.CopyUnitPermissionsValidation), unitHandler.CopyPermissions)
	}

	// Unit Role operations
	unitRoles := protected.Group("/unit-roles")
	{
		unitRoles.PUT("/:unit_role_id/permissions", middleware.ValidateRequest(validation.BulkUpdateUnitRoleModulesValidation), unitHandler.UpdateUnitPermissions)
	}

	// Branch-specific unit operations
	branches := protected.Group("/branches")
	{
		branches.GET("/:id/units/hierarchy", unitHandler.GetUnitHierarchy)
	}

	// User-specific operations
	users := protected.Group("/users")
	{
		users.GET("/:id/effective-permissions", unitHandler.GetUserEffectivePermissions)
	}
}
func setupUnitContextRoutes(protected *gin.RouterGroup, unitContextHandler *handlers.UnitContextHandler) {
	// Unit context routes
	auth := protected.Group("/auth")
	{
		// Get current user's unit context
		auth.GET("/my-unit-context", unitContextHandler.GetMyUnitContext)

		// Get current user's unit permissions
		auth.GET("/my-unit-permissions", unitContextHandler.GetMyUnitPermissions)
	}
}
