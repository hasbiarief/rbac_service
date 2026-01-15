package audit

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	audit := router.Group("/audit")
	{
		audit.GET("", handler.GetAuditLogs)
		audit.POST("", handler.CreateAuditLog)
		audit.GET("/stats", handler.GetAuditStats)
		audit.GET("/users/:userId", handler.GetUserAuditLogs)
		audit.GET("/identity/:identity", handler.GetUserAuditLogsByIdentity)
	}
}
