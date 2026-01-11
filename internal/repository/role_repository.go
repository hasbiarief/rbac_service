package repository

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/models"
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
func (r *RoleRepository) GetAll(limit, offset int, search string, isActive *bool) ([]*models.Role, error) {
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

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
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
func (r *RoleRepository) GetByID(id int64) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, id).Scan(
		&role.ID, &role.Name, &role.Description,
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
func (r *RoleRepository) Create(role *models.Role) error {
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
func (r *RoleRepository) Update(role *models.Role) error {
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

// AssignUserRole assigns a role to a user
func (r *RoleRepository) AssignUserRole(userRole *models.UserRole) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, company_id, branch_id, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, role_id, company_id, branch_id) DO NOTHING
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, userRole.UserID, userRole.RoleID, userRole.CompanyID, userRole.BranchID).Scan(
		&userRole.ID, &userRole.CreatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
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
func (r *RoleRepository) GetUserRoles(userID int64) ([]*models.UserRole, error) {
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

	var userRoles []*models.UserRole
	for rows.Next() {
		userRole := &models.UserRole{}
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
func (r *RoleRepository) GetUsersByRole(roleID int64, limit int) ([]*models.User, error) {
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
			u.id, u.name, u.email, u.user_identity, u.password_hash, u.is_active, u.created_at, u.updated_at,
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

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var unitID, branchID, companyID *int64
		var roleNameFromQuery string

		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.UserIdentity,
			&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
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
		return []*models.User{}, nil
	}

	return users, nil
}

// GetAllUserRoleAssignments - Debug method to see all user role assignments
func (r *RoleRepository) GetAllUserRoleAssignments() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			ur.id as assignment_id,
			u.id as user_id, 
			u.name as user_name,
			u.email as user_email,
			u.is_active as user_active,
			r.id as role_id,
			r.name as role_name,
			ur.company_id,
			ur.branch_id,
			ur.unit_id,
			c.name as company_name,
			b.name as branch_name,
			un.name as unit_name
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		JOIN roles r ON ur.role_id = r.id
		LEFT JOIN companies c ON ur.company_id = c.id
		LEFT JOIN branches b ON ur.branch_id = b.id
		LEFT JOIN units un ON ur.unit_id = un.id
		ORDER BY ur.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role assignments: %w", err)
	}
	defer rows.Close()

	var assignments []map[string]interface{}
	for rows.Next() {
		var assignmentID, userID, roleID, companyID int64
		var branchID, unitID *int64
		var userName, userEmail, roleName, companyName string
		var branchName, unitName *string
		var userActive bool

		err := rows.Scan(
			&assignmentID, &userID, &userName, &userEmail, &userActive,
			&roleID, &roleName, &companyID, &branchID, &unitID,
			&companyName, &branchName, &unitName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}

		assignment := map[string]interface{}{
			"assignment_id": assignmentID,
			"user_id":       userID,
			"user_name":     userName,
			"user_email":    userEmail,
			"user_active":   userActive,
			"role_id":       roleID,
			"role_name":     roleName,
			"company_id":    companyID,
			"company_name":  companyName,
			"branch_id":     branchID,
			"unit_id":       unitID,
		}

		if branchName != nil {
			assignment["branch_name"] = *branchName
		}
		if unitName != nil {
			assignment["unit_name"] = *unitName
		}

		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// UpdateRoleModules updates module permissions for a role
func (r *RoleRepository) UpdateRoleModules(roleID int64, modules []*models.RoleModule) error {
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
func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &models.Role{}
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
func (r *RoleRepository) GetWithPermissions(id int64) (*models.RoleWithPermissions, error) {
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

	var modulePermissions []models.RoleModulePermission
	for rows.Next() {
		perm := models.RoleModulePermission{}
		err := rows.Scan(
			&perm.ModuleID, &perm.ModuleName, &perm.ModuleURL,
			&perm.CanRead, &perm.CanWrite, &perm.CanDelete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role module permission: %w", err)
		}
		modulePermissions = append(modulePermissions, perm)
	}

	return &models.RoleWithPermissions{
		Role:    *role,
		Modules: modulePermissions,
	}, nil
}

// GetRoleModules retrieves module permissions for a role
func (r *RoleRepository) GetRoleModules(roleID int64) ([]*models.RoleModule, error) {
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

	var roleModules []*models.RoleModule
	for rows.Next() {
		roleModule := &models.RoleModule{}
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

// GetUserRolesByUserID - Debug method to check specific user's role assignments
func (r *RoleRepository) GetUserRolesByUserID(userID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			ur.id as assignment_id,
			ur.user_id,
			u.name as user_name,
			u.email as user_email,
			u.is_active as user_active,
			ur.role_id,
			r.name as role_name,
			ur.company_id,
			ur.branch_id,
			ur.unit_id,
			c.name as company_name,
			b.name as branch_name,
			un.name as unit_name
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		JOIN roles r ON ur.role_id = r.id
		LEFT JOIN companies c ON ur.company_id = c.id
		LEFT JOIN branches b ON ur.branch_id = b.id
		LEFT JOIN units un ON ur.unit_id = un.id
		WHERE ur.user_id = $1
		ORDER BY ur.id
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user role assignments for user %d: %w", userID, err)
	}
	defer rows.Close()

	var assignments []map[string]interface{}
	for rows.Next() {
		var assignmentID, roleID, companyID int64
		var branchID, unitID *int64
		var userName, userEmail, roleName, companyName string
		var branchName, unitName *string
		var userActive bool

		err := rows.Scan(
			&assignmentID, &userID, &userName, &userEmail, &userActive,
			&roleID, &roleName, &companyID, &branchID, &unitID,
			&companyName, &branchName, &unitName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}

		assignment := map[string]interface{}{
			"assignment_id": assignmentID,
			"user_id":       userID,
			"user_name":     userName,
			"user_email":    userEmail,
			"user_active":   userActive,
			"role_id":       roleID,
			"role_name":     roleName,
			"company_id":    companyID,
			"company_name":  companyName,
			"branch_id":     branchID,
			"unit_id":       unitID,
		}

		if branchName != nil {
			assignment["branch_name"] = *branchName
		}
		if unitName != nil {
			assignment["unit_name"] = *unitName
		}

		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

// GetRoleUsersMapping - Debug method to show role-user mapping
func (r *RoleRepository) GetRoleUsersMapping() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			r.id as role_id,
			r.name as role_name,
			COUNT(ur.user_id) as user_count,
			STRING_AGG(DISTINCT u.id::text || ':' || u.name, ', ') as users_list
		FROM roles r
		LEFT JOIN user_roles ur ON r.id = ur.role_id
		LEFT JOIN users u ON ur.user_id = u.id AND u.is_active = true
		WHERE r.is_active = true
		GROUP BY r.id, r.name
		ORDER BY r.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get role users mapping: %w", err)
	}
	defer rows.Close()

	var mappings []map[string]interface{}
	for rows.Next() {
		var roleID int64
		var roleName string
		var userCount int
		var usersList *string

		err := rows.Scan(&roleID, &roleName, &userCount, &usersList)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mapping: %w", err)
		}

		mapping := map[string]interface{}{
			"role_id":    roleID,
			"role_name":  roleName,
			"user_count": userCount,
		}

		if usersList != nil {
			mapping["users_list"] = *usersList
		} else {
			mapping["users_list"] = "No users assigned"
		}

		mappings = append(mappings, mapping)
	}

	return mappings, nil
}
