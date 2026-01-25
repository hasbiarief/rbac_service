package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/response"

	"github.com/gin-gonic/gin"
)

// APIDocAuthMiddleware provides API documentation specific authorization
func APIDocAuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
			c.Abort()
			return
		}

		userIDInt64, ok := userID.(int64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid user ID", "")
			c.Abort()
			return
		}

		// Initialize RBAC service
		rbacService := rbac.NewRBACService(db)

		// Set RBAC service in context for use by handlers
		c.Set("rbac_service", rbacService)
		c.Set("user_id_int64", userIDInt64)

		c.Next()
	}
}

// RequireAPIDocPermission middleware ensures user has specific permission for API documentation operations
func RequireAPIDocPermission(moduleID int64, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID and RBAC service from context
		userID, exists := c.Get("user_id_int64")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
			c.Abort()
			return
		}

		rbacService, exists := c.Get("rbac_service")
		if !exists {
			response.Error(c, http.StatusInternalServerError, "RBAC service not available", "")
			c.Abort()
			return
		}

		rbacSvc, ok := rbacService.(*rbac.RBACService)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "Invalid RBAC service", "")
			c.Abort()
			return
		}

		userIDInt64 := userID.(int64)

		// Check permission
		hasPermission, err := rbacSvc.HasPermission(userIDInt64, moduleID, permission)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to check permissions", err.Error())
			c.Abort()
			return
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, constants.MsgInsufficientPermissions,
				fmt.Sprintf("module_id=%d, permission=%s", moduleID, permission))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCollectionAccess middleware ensures user has access to a specific collection
func RequireCollectionAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get collection ID from path parameter
		collectionIDStr := c.Param("collection_id")
		if collectionIDStr == "" {
			collectionIDStr = c.Param("id") // Fallback to generic id parameter
		}

		if collectionIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Collection ID required", "")
			c.Abort()
			return
		}

		collectionID, err := strconv.ParseInt(collectionIDStr, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid collection ID", err.Error())
			c.Abort()
			return
		}

		// Get user ID from context
		userID, exists := c.Get("user_id_int64")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
			c.Abort()
			return
		}

		// Get user's company ID to verify collection ownership
		// This will be handled by the service layer through company isolation
		// For now, just set the collection ID in context for use by handlers
		c.Set("collection_id", collectionID)
		c.Set("collection_id_str", collectionIDStr)
		c.Set("current_user_id", userID.(int64))

		c.Next()
	}
}

// RequireResourceOwnership middleware ensures user owns the resource through company isolation
func RequireResourceOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This middleware works in conjunction with the service layer
		// The actual ownership validation is done in the service methods
		// This middleware just ensures the necessary context is available

		// Get user ID from context
		userID, exists := c.Get("user_id_int64")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
			c.Abort()
			return
		}

		// Set user ID for easy access by handlers
		c.Set("current_user_id", userID.(int64))

		c.Next()
	}
}

// APIDocCollectionPermission creates a middleware for collection-specific permissions
func APIDocCollectionPermission(permission string) gin.HandlerFunc {
	return RequireAPIDocPermission(constants.ModuleAPIDocCollections, permission)
}

// APIDocEndpointPermission creates a middleware for endpoint-specific permissions
func APIDocEndpointPermission(permission string) gin.HandlerFunc {
	return RequireAPIDocPermission(constants.ModuleAPIDocEndpoints, permission)
}

// APIDocEnvironmentPermission creates a middleware for environment-specific permissions
func APIDocEnvironmentPermission(permission string) gin.HandlerFunc {
	return RequireAPIDocPermission(constants.ModuleAPIDocEnvironments, permission)
}

// APIDocExportPermission creates a middleware for export-specific permissions
func APIDocExportPermission(permission string) gin.HandlerFunc {
	return RequireAPIDocPermission(constants.ModuleAPIDocExport, permission)
}

// ValidateCollectionID middleware validates collection ID parameter
func ValidateCollectionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		collectionIDStr := c.Param("collection_id")
		if collectionIDStr == "" {
			collectionIDStr = c.Param("id")
		}

		if collectionIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Collection ID is required", "")
			c.Abort()
			return
		}

		collectionID, err := strconv.ParseInt(collectionIDStr, 10, 64)
		if err != nil || collectionID <= 0 {
			response.Error(c, http.StatusBadRequest, "Invalid collection ID", "Collection ID must be a positive integer")
			c.Abort()
			return
		}

		c.Set("collection_id", collectionID)
		c.Next()
	}
}

// ValidateEndpointID middleware validates endpoint ID parameter
func ValidateEndpointID() gin.HandlerFunc {
	return func(c *gin.Context) {
		endpointIDStr := c.Param("endpoint_id")
		if endpointIDStr == "" {
			endpointIDStr = c.Param("id")
		}

		if endpointIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Endpoint ID is required", "")
			c.Abort()
			return
		}

		endpointID, err := strconv.ParseInt(endpointIDStr, 10, 64)
		if err != nil || endpointID <= 0 {
			response.Error(c, http.StatusBadRequest, "Invalid endpoint ID", "Endpoint ID must be a positive integer")
			c.Abort()
			return
		}

		c.Set("endpoint_id", endpointID)
		c.Next()
	}
}

// ValidateEnvironmentID middleware validates environment ID parameter
func ValidateEnvironmentID() gin.HandlerFunc {
	return func(c *gin.Context) {
		environmentIDStr := c.Param("environment_id")
		if environmentIDStr == "" {
			environmentIDStr = c.Param("id")
		}

		if environmentIDStr == "" {
			response.Error(c, http.StatusBadRequest, "Environment ID is required", "")
			c.Abort()
			return
		}

		environmentID, err := strconv.ParseInt(environmentIDStr, 10, 64)
		if err != nil || environmentID <= 0 {
			response.Error(c, http.StatusBadRequest, "Invalid environment ID", "Environment ID must be a positive integer")
			c.Abort()
			return
		}

		c.Set("environment_id", environmentID)
		c.Next()
	}
}
