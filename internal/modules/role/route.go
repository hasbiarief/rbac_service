package role

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	roleManagement := api.Group("/role-management")
	{
		roleManagement.POST("/assign-user-role", handler.AssignUserRole)
		roleManagement.PUT("/role/:roleId/modules", handler.UpdateRoleModules)
		roleManagement.DELETE("/user/:userId/role/:roleId", handler.RemoveUserRole)
		roleManagement.GET("/role/:roleId/users", handler.GetUsersByRole)
		roleManagement.GET("/user/:userId/roles", handler.GetUserRoles)
		roleManagement.GET("/user/:userId/access-summary", handler.GetUserAccessSummary)

		// Debug endpoints
		roleManagement.GET("/debug/all-assignments", handler.GetAllUserRoleAssignments)
		roleManagement.GET("/debug/user/:userId/roles", handler.GetUserRolesByUserID)
		roleManagement.GET("/debug/role-users-mapping", handler.GetRoleUsersMapping)
	}

	roles := api.Group("/roles")
	{
		roles.GET("", handler.GetRoles)
		roles.GET("/:id", handler.GetRoleByID)
		roles.POST("", handler.CreateRole)
		roles.PUT("/:id", handler.UpdateRole)
		roles.DELETE("/:id", handler.DeleteRole)
	}
}
