package role

import (
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	roleManagement := api.Group("/role-management")
	{
		roleManagement.POST("/assign-user-role",
			middleware.ValidateJSON(&AssignRoleRequest{}),
			handler.AssignUserRole,
		)
		roleManagement.POST("/bulk-assign-roles",
			middleware.ValidateJSON(&BulkAssignRoleRequest{}),
			handler.BulkAssignUserRole,
		)
		roleManagement.PUT("/role/:roleId/modules",
			middleware.ValidateJSON(&UpdateRolePermissionsRequest{}),
			handler.UpdateRoleModules,
		)
		roleManagement.DELETE("/user/:userId/role/:roleId", handler.RemoveUserRole)
		roleManagement.GET("/role/:roleId/users", handler.GetUsersByRole)
		roleManagement.GET("/user/:userId/roles", handler.GetUserRoles)
		roleManagement.GET("/user/:userId/access-summary", handler.GetUserAccessSummary)
	}

	roles := api.Group("/roles")
	{
		roles.GET("", handler.GetRoles)
		roles.GET("/:id", handler.GetRoleByID)
		roles.POST("",
			middleware.ValidateJSON(&CreateRoleRequest{}),
			handler.CreateRole,
		)
		roles.PUT("/:id",
			middleware.ValidateJSON(&UpdateRoleRequest{}),
			handler.UpdateRole,
		)
		roles.DELETE("/:id", handler.DeleteRole)
	}
}
