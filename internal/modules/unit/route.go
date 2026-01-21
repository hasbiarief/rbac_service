package unit

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	// Unit CRUD routes
	units := router.Group("/units")
	{
		units.GET("", handler.GetUnits)
		units.POST("", handler.CreateUnit)
		units.GET("/:id", handler.GetUnitByID)
		units.PUT("/:id", handler.UpdateUnit)
		units.DELETE("/:id", handler.DeleteUnit)
		units.GET("/:id/stats", handler.GetUnitWithStats)

		// Unit role management
		units.POST("/:id/roles/:role_id", handler.AssignRoleToUnit)
		units.DELETE("/:id/roles/:role_id", handler.RemoveRoleFromUnit)
		units.GET("/:id/roles", handler.GetUnitRoles)
		units.GET("/:id/roles/:role_id/permissions", handler.GetUnitPermissions)
	}

	// Branch unit hierarchy
	branches := router.Group("/branches")
	{
		branches.GET("/:id/units/hierarchy", handler.GetUnitHierarchy)
	}

	// Unit role permissions
	unitRoles := router.Group("/unit-roles")
	{
		unitRoles.PUT("/:unit_role_id/permissions", handler.UpdateUnitPermissions)
	}

	// Permission management
	router.POST("/units/copy-permissions", handler.CopyPermissions)
	router.POST("/units/copy-unit-role-permissions", handler.CopyUnitRolePermissions)
	router.GET("/units/unit-role-info", handler.GetUnitRoleInfo)

	// User effective permissions
	users := router.Group("/users")
	{
		users.GET("/:id/effective-permissions", handler.GetUserEffectivePermissions)
	}
}
