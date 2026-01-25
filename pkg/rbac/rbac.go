package rbac

import (
	"database/sql"
	"fmt"
	"strings"
)

type UserPermissions struct {
	UserID  int64
	Roles   []string
	Modules map[int64]ModulePermission
}

type ModulePermission struct {
	ModuleID   int64
	CanRead    bool
	CanWrite   bool
	CanDelete  bool
	CanApprove bool
}

type RBACService struct {
	db *sql.DB
}

func NewRBACService(db *sql.DB) *RBACService {
	return &RBACService{db: db}
}

// GetUserPermissions retrieves all permissions for a user with subscription filtering
func (r *RBACService) GetUserPermissions(userID int64) (*UserPermissions, error) {
	permissions := &UserPermissions{
		UserID:  userID,
		Roles:   []string{},
		Modules: make(map[int64]ModulePermission),
	}

	// Get user roles
	roleQuery := `
		SELECT DISTINCT r.name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(roleQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roleName string
		if err := rows.Scan(&roleName); err != nil {
			continue
		}
		permissions.Roles = append(permissions.Roles, roleName)
	}

	// Get user's company ID for subscription filtering
	var companyID int64
	companyQuery := `SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1`
	err = r.db.QueryRow(companyQuery, userID).Scan(&companyID)
	if err != nil {
		// If no company found, use basic permissions only
		return r.getUserBasicPermissions(userID, permissions)
	}

	// Get user module permissions with subscription filtering
	permQuery := `
		SELECT 
			rm.module_id,
			MAX(CASE WHEN rm.can_read THEN 1 ELSE 0 END) as can_read,
			MAX(CASE WHEN rm.can_write THEN 1 ELSE 0 END) as can_write,
			MAX(CASE WHEN rm.can_delete THEN 1 ELSE 0 END) as can_delete,
			MAX(CASE WHEN rm.can_approve THEN 1 ELSE 0 END) as can_approve
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
		JOIN subscriptions s ON pm.plan_id = s.plan_id
		WHERE ur.user_id = $1
			AND s.company_id = $2
			AND s.status = 'active'
			AND s.end_date > CURRENT_DATE
			AND m.is_active = true
		GROUP BY rm.module_id
	`

	permRows, err := r.db.Query(permQuery, userID, companyID)
	if err != nil {
		return r.getUserBasicPermissions(userID, permissions)
	}
	defer permRows.Close()

	for permRows.Next() {
		var moduleID int64
		var canRead, canWrite, canDelete, canApprove int

		if err := permRows.Scan(&moduleID, &canRead, &canWrite, &canDelete, &canApprove); err != nil {
			continue
		}

		permissions.Modules[moduleID] = ModulePermission{
			ModuleID:   moduleID,
			CanRead:    canRead == 1,
			CanWrite:   canWrite == 1,
			CanDelete:  canDelete == 1,
			CanApprove: canApprove == 1,
		}
	}

	// If no subscription modules found, fallback to basic
	if len(permissions.Modules) == 0 {
		return r.getUserBasicPermissions(userID, permissions)
	}

	return permissions, nil
}

// getUserBasicPermissions returns basic tier permissions only
func (r *RBACService) getUserBasicPermissions(userID int64, permissions *UserPermissions) (*UserPermissions, error) {
	permQuery := `
		SELECT 
			rm.module_id,
			MAX(CASE WHEN rm.can_read THEN 1 ELSE 0 END) as can_read,
			MAX(CASE WHEN rm.can_write THEN 1 ELSE 0 END) as can_write,
			MAX(CASE WHEN rm.can_delete THEN 1 ELSE 0 END) as can_delete,
			MAX(CASE WHEN rm.can_approve THEN 1 ELSE 0 END) as can_approve
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		WHERE ur.user_id = $1
			AND m.is_active = true
			AND (m.subscription_tier = 'basic' OR m.subscription_tier IS NULL)
		GROUP BY rm.module_id
	`

	permRows, err := r.db.Query(permQuery, userID)
	if err != nil {
		return permissions, fmt.Errorf("failed to get basic permissions: %w", err)
	}
	defer permRows.Close()

	for permRows.Next() {
		var moduleID int64
		var canRead, canWrite, canDelete, canApprove int

		if err := permRows.Scan(&moduleID, &canRead, &canWrite, &canDelete, &canApprove); err != nil {
			continue
		}

		permissions.Modules[moduleID] = ModulePermission{
			ModuleID:   moduleID,
			CanRead:    canRead == 1,
			CanWrite:   canWrite == 1,
			CanDelete:  canDelete == 1,
			CanApprove: canApprove == 1,
		}
	}

	return permissions, nil
}

// HasPermission checks if user has specific permission for a module
func (r *RBACService) HasPermission(userID int64, moduleID int64, permission string) (bool, error) {
	// Special case: CONSOLE ADMIN role has full access to API Documentation modules (139-143)
	if moduleID >= 139 && moduleID <= 143 {
		hasConsoleAdminRole, err := r.hasConsoleAdminRole(userID)
		if err != nil {
			return false, err
		}
		if hasConsoleAdminRole {
			return true, nil // CONSOLE ADMIN has full access to all API Documentation modules
		}
	}

	permissions, err := r.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	modulePerm, exists := permissions.Modules[moduleID]
	if !exists {
		return false, nil
	}

	switch permission {
	case "read":
		return modulePerm.CanRead, nil
	case "write":
		return modulePerm.CanWrite, nil
	case "delete":
		return modulePerm.CanDelete, nil
	case "approve":
		return modulePerm.CanApprove, nil
	default:
		return false, fmt.Errorf("invalid permission type: %s", permission)
	}
}

// hasConsoleAdminRole checks if user has CONSOLE ADMIN role (role_id=13)
func (r *RBACService) hasConsoleAdminRole(userID int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.id = 13
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check console admin role: %w", err)
	}

	return count > 0, nil
}

// HasRole checks if user has specific role
func (r *RBACService) HasRole(userID int64, roleName string) (bool, error) {
	permissions, err := r.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, role := range permissions.Roles {
		if role == roleName {
			return true, nil
		}
	}
	return false, nil
}

// GetAccessibleModules returns list of module IDs user can access with specific permission
func (r *RBACService) GetAccessibleModules(userID int64, permission string) ([]int64, error) {
	permissions, err := r.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	var moduleIDs []int64
	for moduleID, modulePerm := range permissions.Modules {
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
		}

		if hasPermission {
			moduleIDs = append(moduleIDs, moduleID)
		}
	}

	return moduleIDs, nil
}

// IsSuperAdmin checks if user is super admin
func (r *RBACService) IsSuperAdmin(userID int64) (bool, error) {
	return r.HasRole(userID, "SUPER_ADMIN")
}

// GetFilteredModules returns modules that user can access with specific permission
func (r *RBACService) GetFilteredModules(userID int64, permission string, limit, offset int, category, search string, isActive *bool) ([]*ModuleInfo, error) {
	permissions, err := r.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// Build accessible module IDs based on permission
	var accessibleModuleIDs []int64
	for moduleID, modulePerm := range permissions.Modules {
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
			hasPermission = modulePerm.CanRead // Default to read permission
		}

		if hasPermission {
			accessibleModuleIDs = append(accessibleModuleIDs, moduleID)
		}
	}

	if len(accessibleModuleIDs) == 0 {
		return []*ModuleInfo{}, nil
	}

	// Build IN clause for module IDs
	placeholders := make([]string, len(accessibleModuleIDs))
	args := make([]interface{}, len(accessibleModuleIDs))
	for i, id := range accessibleModuleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// Build query to get module details
	query := fmt.Sprintf(`
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active
		FROM modules 
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	argIndex := len(args) + 1

	// Add filters
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Add ordering and pagination
	query += " ORDER BY category, name"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered modules: %w", err)
	}
	defer rows.Close()

	var modules []*ModuleInfo
	for rows.Next() {
		module := &ModuleInfo{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID,
			&module.SubscriptionTier, &module.IsActive,
		)
		if err != nil {
			continue
		}
		modules = append(modules, module)
	}

	return modules, nil
}

type ModuleInfo struct {
	ID               int64  `json:"id"`
	Category         string `json:"category"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier"`
	IsActive         bool   `json:"is_active"`
}
