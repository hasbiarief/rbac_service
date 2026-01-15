package subscription

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Public plan routes
	plans := router.Group("/subscription-plans")
	{
		plans.GET("", handler.GetAllPlans)
		plans.GET("/:id", handler.GetPlanByID)
		plans.POST("", handler.CreateSubscriptionPlan)
		plans.PUT("/:id", handler.UpdateSubscriptionPlan)
		plans.DELETE("/:id", handler.DeleteSubscriptionPlan)
	}

	// Subscription routes
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
