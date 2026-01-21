package role

import (
	"database/sql"
	"fmt"

	// removed
	"gin-scalable-api/pkg/model"
)

type RoleRepository struct {
	*model.Repository
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAll retrieves all roles with pagination and filtering
func (r *RoleRepository) GetAll(limit, offset int, search string, isActive *bool) ([]*Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM roles
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

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

	query += " ORDER BY name"

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
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	var roles []*Role
	for rows.Next() {
		role := &Role{}
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.IsActive,
			&role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(id int64) (*Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &Role{}
	err := r.db.QueryRow(query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.IsActive,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// Create creates a new role
func (r *RoleRepository) Create(role *Role) error {
	query, values := r.BuildInsertQuery(role)

	err := r.db.QueryRow(query, values...).Scan(
		&role.ID, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// Update updates a role
func (r *RoleRepository) Update(role *Role) error {
	query, values := r.BuildUpdateQuery(role, role.ID)

	err := r.db.QueryRow(query, values...).Scan(&role.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

// Delete deletes a role
func (r *RoleRepository) Delete(id int64) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// CheckUserExists verifies if a user exists (query minimal field, no cross-module import)
func (r *RoleRepository) CheckUserExists(userID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)"
	err := r.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}

// AssignUserRole assigns a role to a user
func (r *RoleRepository) AssignUserRole(userRole *UserRole) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, company_id, branch_id, unit_id, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, userRole.UserID, userRole.RoleID, userRole.CompanyID, userRole.BranchID, userRole.UnitID).Scan(
		&userRole.ID, &userRole.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to assign user role: %w", err)
	}

	return nil
}

// RemoveUserRole removes a role from a user
func (r *RoleRepository) RemoveUserRole(userID, roleID, companyID int64) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2 AND company_id = $3`

	result, err := r.db.Exec(query, userID, roleID, companyID)
	if err != nil {
		return fmt.Errorf("failed to remove user role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user role assignment not found")
	}

	return nil
}

// GetUserRoles retrieves user role assignments for a specific user
func (r *RoleRepository) GetUserRoles(userID int64) ([]*UserRole, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.company_id, ur.branch_id, ur.created_at
		FROM user_roles ur
		WHERE ur.user_id = $1
		ORDER BY ur.created_at
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var userRoles []*UserRole
	for rows.Next() {
		userRole := &UserRole{}
		err := rows.Scan(
			&userRole.ID, &userRole.UserID, &userRole.RoleID,
			&userRole.CompanyID, &userRole.BranchID, &userRole.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// GetUsersByRole retrieves users assigned to a specific role
func (r *RoleRepository) GetUsersByRole(roleID int64, limit int) ([]*User, error) {
	// First, let's check what role we're looking for
	var roleName string
	roleCheckQuery := "SELECT name FROM roles WHERE id = $1"
	err := r.db.QueryRow(roleCheckQuery, roleID).Scan(&roleName)
	if err != nil {
		return nil, fmt.Errorf("role with ID %d not found: %w", roleID, err)
	}

	// Enhanced query that includes unit context and provides more debugging info
	query := `
		SELECT DISTINCT 
			u.id, u.name, u.email, u.user_identity, u.is_active,
			ur.unit_id, ur.branch_id, ur.company_id,
			r.name as role_name
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.role_id = $1 AND u.is_active = true
		ORDER BY u.name
	`
	args := []interface{}{roleID}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role %s (ID: %d): %w", roleName, roleID, err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		var unitID, branchID, companyID *int64
		var roleNameFromQuery string

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.UserIdentity,
			&user.IsActive,
			&unitID, &branchID, &companyID, &roleNameFromQuery,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user for role %s: %w", roleName, err)
		}
		users = append(users, user)
	}

	// If no users found, let's provide detailed debugging info
	if len(users) == 0 {
		// Check total user_roles for this role (including inactive users)
		var totalAssignments int
		countQuery := `
			SELECT COUNT(*) 
			FROM user_roles ur 
			JOIN users u ON ur.user_id = u.id 
			WHERE ur.role_id = $1
		`
		err = r.db.QueryRow(countQuery, roleID).Scan(&totalAssignments)
		if err == nil && totalAssignments > 0 {
			// There are assignments but no active users
			var inactiveCount int
			inactiveQuery := `
				SELECT COUNT(*) 
				FROM user_roles ur 
				JOIN users u ON ur.user_id = u.id 
				WHERE ur.role_id = $1 AND u.is_active = false
			`
			r.db.QueryRow(inactiveQuery, roleID).Scan(&inactiveCount)

			return nil, fmt.Errorf("role %s (ID: %d) has %d total assignments but %d are inactive users",
				roleName, roleID, totalAssignments, inactiveCount)
		}

		// Return empty slice for no assignments (this is valid)
		return []*User{}, nil
	}

	return users, nil
}

// UpdateRoleModules updates module permissions for a role
func (r *RoleRepository) UpdateRoleModules(roleID int64, modules []*RoleModule) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing role modules
	_, err = tx.Exec("DELETE FROM role_modules WHERE role_id = $1", roleID)
	if err != nil {
		return fmt.Errorf("failed to delete existing role modules: %w", err)
	}

	// Insert new role modules
	for _, module := range modules {
		_, err = tx.Exec(`
			INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
			VALUES ($1, $2, $3, $4, $5)
		`, roleID, module.ModuleID, module.CanRead, module.CanWrite, module.CanDelete)
		if err != nil {
			return fmt.Errorf("failed to insert role module: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(name string) (*Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &Role{}
	err := r.db.QueryRow(query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.IsActive,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// GetWithPermissions retrieves a role with its module permissions
func (r *RoleRepository) GetWithPermissions(id int64) (*RoleWithPermissions, error) {
	// First get the role
	role, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get role module permissions with module details
	query := `
		SELECT rm.module_id, m.name as module_name, m.url as module_url,
		       rm.can_read, rm.can_write, rm.can_delete
		FROM role_modules rm
		JOIN modules m ON rm.module_id = m.id
		WHERE rm.role_id = $1
		ORDER BY m.name
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role module permissions: %w", err)
	}
	defer rows.Close()

	var modulePermissions []RoleModulePermission
	for rows.Next() {
		perm := RoleModulePermission{}
		err := rows.Scan(
			&perm.ModuleID, &perm.ModuleName, &perm.ModuleURL,
			&perm.CanRead, &perm.CanWrite, &perm.CanDelete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role module permission: %w", err)
		}
		modulePermissions = append(modulePermissions, perm)
	}

	return &RoleWithPermissions{
		Role:    *role,
		Modules: modulePermissions,
	}, nil
}

// GetRoleModules retrieves module permissions for a role
func (r *RoleRepository) GetRoleModules(roleID int64) ([]*RoleModule, error) {
	query := `
		SELECT rm.id, rm.role_id, rm.module_id, rm.can_read, rm.can_write, rm.can_delete, rm.created_at
		FROM role_modules rm
		WHERE rm.role_id = $1
		ORDER BY rm.module_id
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role modules: %w", err)
	}
	defer rows.Close()

	var roleModules []*RoleModule
	for rows.Next() {
		roleModule := &RoleModule{}
		err := rows.Scan(
			&roleModule.ID, &roleModule.RoleID, &roleModule.ModuleID,
			&roleModule.CanRead, &roleModule.CanWrite, &roleModule.CanDelete,
			&roleModule.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role module: %w", err)
		}
		roleModules = append(roleModules, roleModule)
	}

	return roleModules, nil
}

// Count returns total count of roles with filtering
func (r *RoleRepository) Count(search string, isActive *bool) (int64, error) {
	query := "SELECT COUNT(*) FROM roles WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get role count: %w", err)
	}

	return count, nil
}
