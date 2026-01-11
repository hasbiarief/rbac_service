package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/response"
	"gin-scalable-api/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(jwtSecret string, redis *redis.Client) gin.HandlerFunc {
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
			// Add debug information as string
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

		// Set user context
		c.Set("user_id", metadata.UserID)
		c.Set("abilities", metadata.Abilities)

		c.Next()
	}
}

// UnitAwareAuthMiddleware provides enhanced authentication with unit context
func UnitAwareAuthMiddleware(jwtSecret string, redis *redis.Client, db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First run standard auth
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "")
			c.Abort()
			return
		}

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
		tokenService := token.NewSimpleTokenService(redis)

		metadata, err := tokenService.GetAccessToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

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
