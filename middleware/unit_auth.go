package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/response"
	"gin-scalable-api/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// UnitAuthMiddleware provides unit-aware authentication
func UnitAuthMiddleware(jwtSecret string, redis *redis.Client, db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			debugInfo := fmt.Sprintf("header_length=%d, parts_count=%d, raw_header='%s'",
				len(authHeader), len(tokenParts), authHeader)
			if len(tokenParts) > 0 {
				debugInfo += fmt.Sprintf(", first_part='%s'", tokenParts[0])
			}
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", debugInfo)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Initialize token service
		tokenService := token.NewSimpleTokenService(redis)

		// Get token metadata from Redis
		metadata, err := tokenService.GetAccessToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Check if token is expired
		if time.Now().Unix() > metadata.ExpiresAt {
			response.Error(c, http.StatusUnauthorized, "Token expired", "")
			c.Abort()
			return
		}

		// Set basic user context
		c.Set("user_id", metadata.UserID)
		c.Set("abilities", metadata.Abilities)

		// Load unit-aware permissions if database connection is available
		if dbConn, ok := db.(interface{ GetDB() interface{} }); ok {
			if sqlDB, ok := dbConn.GetDB().(*sql.DB); ok {
				unitRBACService := rbac.NewUnitRBACService(sqlDB)

				// Get comprehensive unit permissions
				unitPermissions, err := unitRBACService.GetUserUnitPermissions(metadata.UserID)
				if err == nil {
					// Set unit context
					c.Set("unit_permissions", unitPermissions)
					c.Set("company_id", unitPermissions.CompanyID)
					c.Set("branch_id", unitPermissions.BranchID)
					c.Set("unit_id", unitPermissions.UnitID)
					c.Set("effective_units", unitPermissions.EffectiveUnits)
					c.Set("is_unit_admin", unitPermissions.IsUnitAdmin)
					c.Set("is_branch_admin", unitPermissions.IsBranchAdmin)
					c.Set("is_company_admin", unitPermissions.IsCompanyAdmin)
					c.Set("unit_roles", unitPermissions.UnitRoles)
				}
			}
		}

		c.Next()
	}
}

// RequireUnitAccess middleware ensures user has access to a specific unit
func RequireUnitAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get unit ID from path parameter
		unitIDStr := c.Param("unit_id")
		if unitIDStr == "" {
			unitIDStr = c.Param("id") // Fallback to generic id parameter
		}

		if unitIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Unit ID required", "")
			c.Abort()
			return
		}

		unitID, err := strconv.ParseInt(unitIDStr, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid unit ID", err.Error())
			c.Abort()
			return
		}

		// Check if user has unit permissions loaded
		unitPermissions, exists := c.Get("unit_permissions")
		if !exists {
			response.Error(c, http.StatusForbidden, "Unit permissions not available", "")
			c.Abort()
			return
		}

		permissions, ok := unitPermissions.(*rbac.UnitUserPermissions)
		if !ok {
			response.Error(c, http.StatusForbidden, "Invalid unit permissions", "")
			c.Abort()
			return
		}

		// Check if user can access this unit
		hasAccess := false

		// Company/Branch admins have access to all units in their scope
		if permissions.IsCompanyAdmin || permissions.IsBranchAdmin {
			hasAccess = true
		} else {
			// Check if unit is in user's effective units
			for _, effectiveUnitID := range permissions.EffectiveUnits {
				if effectiveUnitID == unitID {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			response.Error(c, http.StatusForbidden, "Access denied to unit", fmt.Sprintf("unit_id=%d", unitID))
			c.Abort()
			return
		}

		// Set current unit context
		c.Set("current_unit_id", unitID)
		c.Next()
	}
}

// RequireUnitPermission middleware ensures user has specific permission for a module in unit context
func RequireUnitPermission(moduleID int64, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get unit permissions
		unitPermissions, exists := c.Get("unit_permissions")
		if !exists {
			response.Error(c, http.StatusForbidden, "Unit permissions not available", "")
			c.Abort()
			return
		}

		permissions, ok := unitPermissions.(*rbac.UnitUserPermissions)
		if !ok {
			response.Error(c, http.StatusForbidden, "Invalid unit permissions", "")
			c.Abort()
			return
		}

		// Check module permission
		modulePerm, exists := permissions.Modules[moduleID]
		if !exists {
			response.Error(c, http.StatusForbidden, "No access to module", fmt.Sprintf("module_id=%d", moduleID))
			c.Abort()
			return
		}

		hasPermission := false
		switch permission {
		case "read":
			hasPermission = modulePerm.CanRead
		case "write":
			hasPermission = modulePerm.CanWrite
		case "delete":
			hasPermission = modulePerm.CanDelete
		case "approve":
			hasPermission = modulePerm.CanApprove
		default:
			response.Error(c, http.StatusBadRequest, "Invalid permission type", permission)
			c.Abort()
			return
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, "Insufficient permissions",
				fmt.Sprintf("module_id=%d, permission=%s", moduleID, permission))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireUnitAdmin middleware ensures user is a unit administrator
func RequireUnitAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isUnitAdmin, exists := c.Get("is_unit_admin")
		if !exists || !isUnitAdmin.(bool) {
			response.Error(c, http.StatusForbidden, "Unit administrator access required", "")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireBranchAdmin middleware ensures user is a branch administrator
func RequireBranchAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isBranchAdmin, exists := c.Get("is_branch_admin")
		if !exists || !isBranchAdmin.(bool) {
			response.Error(c, http.StatusForbidden, "Branch administrator access required", "")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireCompanyAdmin middleware ensures user is a company administrator
func RequireCompanyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isCompanyAdmin, exists := c.Get("is_company_admin")
		if !exists || !isCompanyAdmin.(bool) {
			response.Error(c, http.StatusForbidden, "Company administrator access required", "")
			c.Abort()
			return
		}
		c.Next()
	}
}

// FilterByUnitAccess middleware filters results based on user's unit access
func FilterByUnitAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get unit permissions
		unitPermissions, exists := c.Get("unit_permissions")
		if exists {
			if permissions, ok := unitPermissions.(*rbac.UnitUserPermissions); ok {
				// Set accessible units for filtering in handlers
				c.Set("accessible_units", permissions.EffectiveUnits)
			}
		}
		c.Next()
	}
}
