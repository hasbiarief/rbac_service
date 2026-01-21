package subscription

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Public plan routes (read-only)
	plans := router.Group("/subscription-plans")
	{
		plans.GET("", handler.GetAllPlans)
		plans.GET("/:id", handler.GetPlanByID)
	}
}

// RegisterProtectedRoutes registers protected subscription routes
func RegisterProtectedRoutes(router *gin.RouterGroup, handler *Handler) {
	// Admin/Protected plan routes
	adminPlans := router.Group("/admin/subscription-plans")
	{
		adminPlans.POST("", handler.CreateSubscriptionPlan)
		adminPlans.PUT("/:id", handler.UpdateSubscriptionPlan)
		adminPlans.DELETE("/:id", handler.DeleteSubscriptionPlan)
	}

	// Plan modules management (separate group to avoid conflicts)
	planModules := router.Group("/admin/plan-modules")
	{
		planModules.GET("/:plan_id", handler.GetPlanModules)
		planModules.POST("/:plan_id", handler.AddModulesToPlan)
		planModules.DELETE("/:plan_id/:module_id", handler.RemoveModuleFromPlan)
	}

	// Subscription routes (all protected)
	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.GET("", handler.GetAllSubscriptions)
		subscriptions.POST("", handler.CreateSubscription)
		subscriptions.GET("/:id", handler.GetSubscriptionByID)
		subscriptions.PUT("/:id", handler.UpdateSubscription)
	}

	// Company subscription routes
	companies := router.Group("/companies")
	{
		companies.GET("/:id/subscription", handler.GetCompanySubscription)
	}
}
