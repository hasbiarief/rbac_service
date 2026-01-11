package rbac

import (
	"database/sql"
	"fmt"
	"strings"
)

// UnitUserPermissions represents user permissions with unit context
type UnitUserPermissions struct {
	UserID         int64
	CompanyID      int64
	BranchID       *int64
	UnitID         *int64
	Roles          []string
	UnitRoles      []UnitRoleInfo
	Modules        map[int64]UnitModulePermission
	EffectiveUnits []int64 // Units user has access to (including parent units)
	IsUnitAdmin    bool    // Can manage unit-level permissions
	IsBranchAdmin  bool    // Can manage branch-level permissions
	IsCompanyAdmin bool    // Can manage company-level permissions
}

// UnitRoleInfo represents role information with unit context
type UnitRoleInfo struct {
	RoleID      int64  `json:"role_id"`
	RoleName    string `json:"role_name"`
	UnitID      *int64 `json:"unit_id,omitempty"`
	UnitName    string `json:"unit_name,omitempty"`
	BranchID    *int64 `json:"branch_id,omitempty"`
	BranchName  string `json:"branch_name,omitempty"`
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`
	Level       string `json:"level"` // "company", "branch", "unit"
}

// UnitModulePermission represents module permissions with unit context
type UnitModulePermission struct {
	ModuleID     int64
	CanRead      bool
	CanWrite     bool
	CanDelete    bool
	CanApprove   bool
	GrantedBy    []PermissionSource // Track where permission comes from
	HighestLevel string             // "company", "branch", "unit"
}

// PermissionSource tracks the source of a permission
type PermissionSource struct {
	Type     string `json:"type"` // "role", "unit_role"
	RoleID   int64  `json:"role_id"`
	RoleName string `json:"role_name"`
	UnitID   *int64 `json:"unit_id,omitempty"`
	UnitName string `json:"unit_name,omitempty"`
	Level    string `json:"level"` // "company", "branch", "unit"
}

// UnitRBACService provides unit-aware RBAC functionality
type UnitRBACService struct {
	db *sql.DB
}

// NewUnitRBACService creates a new unit-aware RBAC service
func NewUnitRBACService(db *sql.DB) *UnitRBACService {
	return &UnitRBACService{db: db}
}

// GetUserUnitPermissions retrieves comprehensive unit-aware permissions for a user
func (r *UnitRBACService) GetUserUnitPermissions(userID int64) (*UnitUserPermissions, error) {
	permissions := &UnitUserPermissions{
		UserID:         userID,
		Roles:          []string{},
		UnitRoles:      []UnitRoleInfo{},
		Modules:        make(map[int64]UnitModulePermission),
		EffectiveUnits: []int64{},
	}

	// Get user's role assignments with unit context
	if err := r.loadUserRoleAssignments(permissions); err != nil {
		return nil, fmt.Errorf("failed to load user role assignments: %w", err)
	}

	// Get effective units (including parent units in hierarchy)
	if err := r.loadEffectiveUnits(permissions); err != nil {
		return nil, fmt.Errorf("failed to load effective units: %w", err)
	}

	// Get module permissions from all sources
	if err := r.loadModulePermissions(permissions); err != nil {
		return nil, fmt.Errorf("failed to load module permissions: %w", err)
	}

	// Determine admin levels
	r.determineAdminLevels(permissions)

	return permissions, nil
}

// loadUserRoleAssignments loads user role assignments with unit context
func (r *UnitRBACService) loadUserRoleAssignments(permissions *UnitUserPermissions) error {
	query := `
		SELECT 
			ur.id, ur.user_id, ur.role_id, ur.company_id, ur.branch_id, ur.unit_id,
			r.name as role_name,
			c.name as company_name,
			b.name as branch_name,
			u.name as unit_name,
			CASE 
				WHEN ur.unit_id IS NOT NULL THEN 'unit'
				WHEN ur.branch_id IS NOT NULL THEN 'branch'
				ELSE 'company'
			END as level
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN companies c ON ur.company_id = c.id
		LEFT JOIN branches b ON ur.branch_id = b.id
		LEFT JOIN units u ON ur.unit_id = u.id
		WHERE ur.user_id = $1 AND r.is_active = true
		ORDER BY 
			CASE 
				WHEN ur.unit_id IS NOT NULL THEN 3
				WHEN ur.branch_id IS NOT NULL THEN 2
				ELSE 1
			END DESC
	`

	rows, err := r.db.Query(query, permissions.UserID)
	if err != nil {
		return fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var roleInfo UnitRoleInfo
		var branchName, unitName *string

		err := rows.Scan(
			&roleInfo.RoleID, &permissions.UserID, &roleInfo.RoleID, &roleInfo.CompanyID,
			&roleInfo.BranchID, &roleInfo.UnitID, &roleInfo.RoleName, &roleInfo.CompanyName,
			&branchName, &unitName, &roleInfo.Level,
		)
		if err != nil {
			continue
		}

		if branchName != nil {
			roleInfo.BranchName = *branchName
		}
		if unitName != nil {
			roleInfo.UnitName = *unitName
		}

		// Set primary context from highest level assignment
		if permissions.CompanyID == 0 {
			permissions.CompanyID = roleInfo.CompanyID
		}
		if permissions.BranchID == nil && roleInfo.BranchID != nil {
			permissions.BranchID = roleInfo.BranchID
		}
		if permissions.UnitID == nil && roleInfo.UnitID != nil {
			permissions.UnitID = roleInfo.UnitID
		}

		permissions.UnitRoles = append(permissions.UnitRoles, roleInfo)
		permissions.Roles = append(permissions.Roles, roleInfo.RoleName)
	}

	return nil
}

// loadEffectiveUnits loads all units user has access to (including parent units)
func (r *UnitRBACService) loadEffectiveUnits(permissions *UnitUserPermissions) error {
	if len(permissions.UnitRoles) == 0 {
		return nil
	}

	// Collect all unit IDs from role assignments
	unitIDs := make(map[int64]bool)
	for _, roleInfo := range permissions.UnitRoles {
		if roleInfo.UnitID != nil {
			unitIDs[*roleInfo.UnitID] = true
		}
	}

	if len(unitIDs) == 0 {
		return nil
	}

	// Build IN clause for unit IDs
	placeholders := make([]string, 0, len(unitIDs))
	args := make([]interface{}, 0, len(unitIDs))
	for unitID := range unitIDs {
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
		args = append(args, unitID)
	}

	// Get unit hierarchy (including parent units)
	query := fmt.Sprintf(`
		WITH RECURSIVE unit_hierarchy AS (
			-- Base case: direct unit assignments
			SELECT id, parent_id, branch_id, 0 as level
			FROM units 
			WHERE id IN (%s)
			
			UNION ALL
			
			-- Recursive case: parent units
			SELECT u.id, u.parent_id, u.branch_id, uh.level + 1
			FROM units u
			JOIN unit_hierarchy uh ON u.id = uh.parent_id
			WHERE uh.level < 10 -- Prevent infinite recursion
		)
		SELECT DISTINCT id FROM unit_hierarchy
		ORDER BY id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to get effective units: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var unitID int64
		if err := rows.Scan(&unitID); err != nil {
			continue
		}
		permissions.EffectiveUnits = append(permissions.EffectiveUnits, unitID)
	}

	return nil
}

// loadModulePermissions loads module permissions from all sources (roles and unit roles)
func (r *UnitRBACService) loadModulePermissions(permissions *UnitUserPermissions) error {
	// Load permissions from traditional role assignments
	if err := r.loadRoleModulePermissions(permissions); err != nil {
		return err
	}

	// Load permissions from unit role assignments
	if err := r.loadUnitRoleModulePermissions(permissions); err != nil {
		return err
	}

	return nil
}

// loadRoleModulePermissions loads permissions from traditional role-module assignments
func (r *UnitRBACService) loadRoleModulePermissions(permissions *UnitUserPermissions) error {
	if len(permissions.UnitRoles) == 0 {
		return nil
	}

	// Get company ID for subscription filtering
	companyID := permissions.CompanyID

	// Build role IDs for query
	roleIDs := make([]interface{}, len(permissions.UnitRoles))
	placeholders := make([]string, len(permissions.UnitRoles))
	for i, roleInfo := range permissions.UnitRoles {
		roleIDs[i] = roleInfo.RoleID
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Query with subscription filtering
	query := fmt.Sprintf(`
		SELECT 
			rm.module_id,
			MAX(CASE WHEN rm.can_read THEN 1 ELSE 0 END) as can_read,
			MAX(CASE WHEN rm.can_write THEN 1 ELSE 0 END) as can_write,
			MAX(CASE WHEN rm.can_delete THEN 1 ELSE 0 END) as can_delete,
			MAX(CASE WHEN rm.can_approve THEN 1 ELSE 0 END) as can_approve,
			r.name as role_name
		FROM role_modules rm
		JOIN roles r ON rm.role_id = r.id
		JOIN modules m ON rm.module_id = m.id
		LEFT JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
		LEFT JOIN subscriptions s ON pm.plan_id = s.plan_id AND s.company_id = $%d
		WHERE rm.role_id IN (%s)
			AND m.is_active = true
			AND (
				s.status = 'active' AND s.end_date > CURRENT_DATE
				OR m.subscription_tier = 'basic' 
				OR m.subscription_tier IS NULL
			)
		GROUP BY rm.module_id, r.name
	`, len(roleIDs)+1, strings.Join(placeholders, ","))

	args := append(roleIDs, companyID)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to get role module permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var moduleID int64
		var canRead, canWrite, canDelete, canApprove int
		var roleName string

		if err := rows.Scan(&moduleID, &canRead, &canWrite, &canDelete, &canApprove, &roleName); err != nil {
			continue
		}

		// Get or create module permission
		modulePerm, exists := permissions.Modules[moduleID]
		if !exists {
			modulePerm = UnitModulePermission{
				ModuleID:     moduleID,
				GrantedBy:    []PermissionSource{},
				HighestLevel: "company",
			}
		}

		// Merge permissions (OR logic - any true makes it true)
		modulePerm.CanRead = modulePerm.CanRead || (canRead == 1)
		modulePerm.CanWrite = modulePerm.CanWrite || (canWrite == 1)
		modulePerm.CanDelete = modulePerm.CanDelete || (canDelete == 1)
		modulePerm.CanApprove = modulePerm.CanApprove || (canApprove == 1)

		// Add permission source
		source := PermissionSource{
			Type:     "role",
			RoleName: roleName,
			Level:    "company", // Traditional roles are company-level
		}
		modulePerm.GrantedBy = append(modulePerm.GrantedBy, source)

		permissions.Modules[moduleID] = modulePerm
	}

	return nil
}

// loadUnitRoleModulePermissions loads permissions from unit-specific role assignments
func (r *UnitRBACService) loadUnitRoleModulePermissions(permissions *UnitUserPermissions) error {
	if len(permissions.EffectiveUnits) == 0 {
		return nil
	}

	// Build unit IDs for query
	unitIDs := make([]interface{}, len(permissions.EffectiveUnits))
	placeholders := make([]string, len(permissions.EffectiveUnits))
	for i, unitID := range permissions.EffectiveUnits {
		unitIDs[i] = unitID
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Query unit role module permissions
	query := fmt.Sprintf(`
		SELECT 
			urm.module_id,
			MAX(CASE WHEN urm.can_read THEN 1 ELSE 0 END) as can_read,
			MAX(CASE WHEN urm.can_write THEN 1 ELSE 0 END) as can_write,
			MAX(CASE WHEN urm.can_delete THEN 1 ELSE 0 END) as can_delete,
			MAX(CASE WHEN urm.can_approve THEN 1 ELSE 0 END) as can_approve,
			r.name as role_name,
			u.name as unit_name,
			u.id as unit_id
		FROM unit_role_modules urm
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		JOIN roles r ON ur.role_id = r.id
		JOIN units u ON ur.unit_id = u.id
		JOIN modules m ON urm.module_id = m.id
		WHERE ur.unit_id IN (%s)
			AND m.is_active = true
		GROUP BY urm.module_id, r.name, u.name, u.id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, unitIDs...)
	if err != nil {
		return fmt.Errorf("failed to get unit role module permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var moduleID, unitID int64
		var canRead, canWrite, canDelete, canApprove int
		var roleName, unitName string

		if err := rows.Scan(&moduleID, &canRead, &canWrite, &canDelete, &canApprove, &roleName, &unitName, &unitID); err != nil {
			continue
		}

		// Get or create module permission
		modulePerm, exists := permissions.Modules[moduleID]
		if !exists {
			modulePerm = UnitModulePermission{
				ModuleID:     moduleID,
				GrantedBy:    []PermissionSource{},
				HighestLevel: "unit",
			}
		}

		// Merge permissions (OR logic - any true makes it true)
		modulePerm.CanRead = modulePerm.CanRead || (canRead == 1)
		modulePerm.CanWrite = modulePerm.CanWrite || (canWrite == 1)
		modulePerm.CanDelete = modulePerm.CanDelete || (canDelete == 1)
		modulePerm.CanApprove = modulePerm.CanApprove || (canApprove == 1)

		// Update highest level if unit-level is more specific
		if modulePerm.HighestLevel == "company" {
			modulePerm.HighestLevel = "unit"
		}

		// Add permission source
		source := PermissionSource{
			Type:     "unit_role",
			RoleName: roleName,
			UnitID:   &unitID,
			UnitName: unitName,
			Level:    "unit",
		}
		modulePerm.GrantedBy = append(modulePerm.GrantedBy, source)

		permissions.Modules[moduleID] = modulePerm
	}

	return nil
}

// determineAdminLevels determines user's administrative levels
func (r *UnitRBACService) determineAdminLevels(permissions *UnitUserPermissions) {
	for _, role := range permissions.Roles {
		switch role {
		case "SUPER_ADMIN", "COMPANY_ADMIN":
			permissions.IsCompanyAdmin = true
			permissions.IsBranchAdmin = true
			permissions.IsUnitAdmin = true
		case "BRANCH_ADMIN":
			permissions.IsBranchAdmin = true
			permissions.IsUnitAdmin = true
		case "UNIT_ADMIN":
			permissions.IsUnitAdmin = true
		}
	}
}

// HasUnitPermission checks if user has specific permission for a module in a unit context
func (r *UnitRBACService) HasUnitPermission(userID int64, moduleID int64, permission string, unitID *int64) (bool, error) {
	permissions, err := r.GetUserUnitPermissions(userID)
	if err != nil {
		return false, err
	}

	// Check if user has access to the unit
	if unitID != nil {
		hasUnitAccess := false
		for _, effectiveUnitID := range permissions.EffectiveUnits {
			if effectiveUnitID == *unitID {
				hasUnitAccess = true
				break
			}
		}
		if !hasUnitAccess {
			return false, nil
		}
	}

	// Check module permission
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

// CanAccessUnit checks if user can access a specific unit
func (r *UnitRBACService) CanAccessUnit(userID int64, unitID int64) (bool, error) {
	permissions, err := r.GetUserUnitPermissions(userID)
	if err != nil {
		return false, err
	}

	// Company/Branch admins can access all units in their scope
	if permissions.IsCompanyAdmin || permissions.IsBranchAdmin {
		return true, nil
	}

	// Check if unit is in user's effective units
	for _, effectiveUnitID := range permissions.EffectiveUnits {
		if effectiveUnitID == unitID {
			return true, nil
		}
	}

	return false, nil
}

// GetAccessibleUnits returns all units user can access
func (r *UnitRBACService) GetAccessibleUnits(userID int64) ([]int64, error) {
	permissions, err := r.GetUserUnitPermissions(userID)
	if err != nil {
		return nil, err
	}

	// Company/Branch admins can access all units in their scope
	if permissions.IsCompanyAdmin {
		return r.getAllCompanyUnits(permissions.CompanyID)
	}

	if permissions.IsBranchAdmin && permissions.BranchID != nil {
		return r.getAllBranchUnits(*permissions.BranchID)
	}

	// Return user's effective units
	return permissions.EffectiveUnits, nil
}

// getAllCompanyUnits returns all unit IDs for a company
func (r *UnitRBACService) getAllCompanyUnits(companyID int64) ([]int64, error) {
	query := `
		SELECT u.id 
		FROM units u
		JOIN branches b ON u.branch_id = b.id
		WHERE b.company_id = $1 AND u.is_active = true
		ORDER BY u.id
	`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company units: %w", err)
	}
	defer rows.Close()

	var unitIDs []int64
	for rows.Next() {
		var unitID int64
		if err := rows.Scan(&unitID); err != nil {
			continue
		}
		unitIDs = append(unitIDs, unitID)
	}

	return unitIDs, nil
}

// getAllBranchUnits returns all unit IDs for a branch
func (r *UnitRBACService) getAllBranchUnits(branchID int64) ([]int64, error) {
	query := `
		SELECT id 
		FROM units 
		WHERE branch_id = $1 AND is_active = true
		ORDER BY id
	`

	rows, err := r.db.Query(query, branchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branch units: %w", err)
	}
	defer rows.Close()

	var unitIDs []int64
	for rows.Next() {
		var unitID int64
		if err := rows.Scan(&unitID); err != nil {
			continue
		}
		unitIDs = append(unitIDs, unitID)
	}

	return unitIDs, nil
}
