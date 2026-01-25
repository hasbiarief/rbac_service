package middleware

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/pkg/rbac"

	"github.com/gin-gonic/gin"
)

// APIDocPermissionChecker provides utility functions for checking API documentation permissions
type APIDocPermissionChecker struct {
	rbacService *rbac.RBACService
}

// NewAPIDocPermissionChecker creates a new permission checker instance
func NewAPIDocPermissionChecker(db *sql.DB) *APIDocPermissionChecker {
	return &APIDocPermissionChecker{
		rbacService: rbac.NewRBACService(db),
	}
}

// CheckCollectionPermission checks if user has permission for collection operations
func (p *APIDocPermissionChecker) CheckCollectionPermission(userID int64, permission string) error {
	hasPermission, err := p.rbacService.HasPermission(userID, constants.ModuleAPIDocCollections, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// CheckEndpointPermission checks if user has permission for endpoint operations
func (p *APIDocPermissionChecker) CheckEndpointPermission(userID int64, permission string) error {
	hasPermission, err := p.rbacService.HasPermission(userID, constants.ModuleAPIDocEndpoints, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// CheckEnvironmentPermission checks if user has permission for environment operations
func (p *APIDocPermissionChecker) CheckEnvironmentPermission(userID int64, permission string) error {
	hasPermission, err := p.rbacService.HasPermission(userID, constants.ModuleAPIDocEnvironments, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// CheckExportPermission checks if user has permission for export operations
func (p *APIDocPermissionChecker) CheckExportPermission(userID int64, permission string) error {
	hasPermission, err := p.rbacService.HasPermission(userID, constants.ModuleAPIDocExport, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// CheckMultiplePermissions checks if user has all specified permissions
func (p *APIDocPermissionChecker) CheckMultiplePermissions(userID int64, permissions []PermissionCheck) error {
	for _, perm := range permissions {
		hasPermission, err := p.rbacService.HasPermission(userID, perm.ModuleID, perm.Permission)
		if err != nil {
			return fmt.Errorf("failed to check permission for module %d: %w", perm.ModuleID, err)
		}
		if !hasPermission {
			return fmt.Errorf("insufficient permissions for module %d (%s)", perm.ModuleID, perm.Permission)
		}
	}
	return nil
}

// PermissionCheck represents a permission requirement
type PermissionCheck struct {
	ModuleID   int64
	Permission string
}

// GetUserPermissions retrieves all user permissions for API documentation modules
func (p *APIDocPermissionChecker) GetUserPermissions(userID int64) (*APIDocUserPermissions, error) {
	permissions, err := p.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	apiDocPerms := &APIDocUserPermissions{
		UserID: userID,
	}

	// Check collections permissions
	if collectionPerm, exists := permissions.Modules[constants.ModuleAPIDocCollections]; exists {
		apiDocPerms.Collections = ModulePermissions{
			CanRead:    collectionPerm.CanRead,
			CanWrite:   collectionPerm.CanWrite,
			CanDelete:  collectionPerm.CanDelete,
			CanApprove: collectionPerm.CanApprove,
		}
	}

	// Check endpoints permissions
	if endpointPerm, exists := permissions.Modules[constants.ModuleAPIDocEndpoints]; exists {
		apiDocPerms.Endpoints = ModulePermissions{
			CanRead:    endpointPerm.CanRead,
			CanWrite:   endpointPerm.CanWrite,
			CanDelete:  endpointPerm.CanDelete,
			CanApprove: endpointPerm.CanApprove,
		}
	}

	// Check environments permissions
	if environmentPerm, exists := permissions.Modules[constants.ModuleAPIDocEnvironments]; exists {
		apiDocPerms.Environments = ModulePermissions{
			CanRead:    environmentPerm.CanRead,
			CanWrite:   environmentPerm.CanWrite,
			CanDelete:  environmentPerm.CanDelete,
			CanApprove: environmentPerm.CanApprove,
		}
	}

	// Check export permissions
	if exportPerm, exists := permissions.Modules[constants.ModuleAPIDocExport]; exists {
		apiDocPerms.Export = ModulePermissions{
			CanRead:    exportPerm.CanRead,
			CanWrite:   exportPerm.CanWrite,
			CanDelete:  exportPerm.CanDelete,
			CanApprove: exportPerm.CanApprove,
		}
	}

	return apiDocPerms, nil
}

// APIDocUserPermissions represents user permissions for API documentation modules
type APIDocUserPermissions struct {
	UserID       int64             `json:"user_id"`
	Collections  ModulePermissions `json:"collections"`
	Endpoints    ModulePermissions `json:"endpoints"`
	Environments ModulePermissions `json:"environments"`
	Export       ModulePermissions `json:"export"`
}

// ModulePermissions represents permissions for a specific module
type ModulePermissions struct {
	CanRead    bool `json:"can_read"`
	CanWrite   bool `json:"can_write"`
	CanDelete  bool `json:"can_delete"`
	CanApprove bool `json:"can_approve"`
}

// HasAnyPermission checks if user has any permission for API documentation
func (p *APIDocUserPermissions) HasAnyPermission() bool {
	return p.Collections.HasAnyPermission() ||
		p.Endpoints.HasAnyPermission() ||
		p.Environments.HasAnyPermission() ||
		p.Export.HasAnyPermission()
}

// HasAnyPermission checks if module has any permission
func (m *ModulePermissions) HasAnyPermission() bool {
	return m.CanRead || m.CanWrite || m.CanDelete || m.CanApprove
}

// Utility functions for getting permission checker from context

// GetPermissionChecker retrieves the permission checker from gin context
func GetPermissionChecker(c *gin.Context) (*APIDocPermissionChecker, error) {
	rbacService, exists := c.Get("rbac_service")
	if !exists {
		return nil, fmt.Errorf("RBAC service not available in context")
	}

	rbacSvc, ok := rbacService.(*rbac.RBACService)
	if !ok {
		return nil, fmt.Errorf("invalid RBAC service type")
	}

	return &APIDocPermissionChecker{rbacService: rbacSvc}, nil
}

// GetUserIDFromContext retrieves the user ID from gin context
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userID, exists := c.Get("user_id_int64")
	if !exists {
		return 0, fmt.Errorf("user ID not available in context")
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	return userIDInt64, nil
}

// RequireAnyAPIDocPermission middleware ensures user has at least one API documentation permission
func RequireAnyAPIDocPermission(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		checker := NewAPIDocPermissionChecker(db)
		permissions, err := checker.GetUserPermissions(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !permissions.HasAnyPermission() {
			c.JSON(403, gin.H{"error": constants.MsgInsufficientPermissions})
			c.Abort()
			return
		}

		// Set permissions in context for use by handlers
		c.Set("api_doc_permissions", permissions)
		c.Next()
	}
}

// GetAPIDocPermissions retrieves API documentation permissions from context
func GetAPIDocPermissions(c *gin.Context) (*APIDocUserPermissions, error) {
	permissions, exists := c.Get("api_doc_permissions")
	if !exists {
		return nil, fmt.Errorf("API documentation permissions not available in context")
	}

	apiDocPerms, ok := permissions.(*APIDocUserPermissions)
	if !ok {
		return nil, fmt.Errorf("invalid API documentation permissions type")
	}

	return apiDocPerms, nil
}
