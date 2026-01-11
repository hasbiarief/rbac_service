package routes

import (
	"gin-scalable-api/internal/handlers"
	"gin-scalable-api/internal/validation"
	"gin-scalable-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUnitRoutes(router *gin.RouterGroup, unitHandler *handlers.UnitHandler) {
	units := router.Group("/units")
	{
		// Unit CRUD operations
		units.GET("", middleware.ValidateRequest(validation.UnitListValidation), unitHandler.GetUnits)
		units.POST("", middleware.ValidateRequest(validation.CreateUnitValidation), unitHandler.CreateUnit)
		units.GET("/:id", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitByID)
		units.PUT("/:id", middleware.ValidateRequest(validation.UpdateUnitValidation), unitHandler.UpdateUnit)
		units.DELETE("/:id", middleware.ValidateRequest(validation.IDValidation), unitHandler.DeleteUnit)

		// Unit statistics
		units.GET("/:id/stats", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitWithStats)

		// Unit roles management
		units.GET("/:id/roles", middleware.ValidateRequest(validation.IDValidation), unitHandler.GetUnitRoles)
		units.POST("/:unit_id/roles/:role_id", unitHandler.AssignRoleToUnit)
		units.DELETE("/:unit_id/roles/:role_id", unitHandler.RemoveRoleFromUnit)

		// Unit permissions management
		units.GET("/:unit_id/roles/:role_id/permissions", unitHandler.GetUnitPermissions)

		// Bulk operations
		units.POST("/copy-permissions", middleware.ValidateRequest(validation.CopyUnitPermissionsValidation), unitHandler.CopyPermissions)
	}

	// Unit Role operations
	unitRoles := router.Group("/unit-roles")
	{
		unitRoles.PUT("/:unit_role_id/permissions", middleware.ValidateRequest(validation.BulkUpdateUnitRoleModulesValidation), unitHandler.UpdateUnitPermissions)
	}

	// Branch-specific unit operations
	branches := router.Group("/branches")
	{
		branches.GET("/:branch_id/units/hierarchy", unitHandler.GetUnitHierarchy)
	}

	// User-specific operations
	users := router.Group("/users")
	{
		users.GET("/:user_id/effective-permissions", unitHandler.GetUserEffectivePermissions)
	}
}
