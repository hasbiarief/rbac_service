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
func (r *RoleRepository) GetAll(limit, offset int, search string) ([]*models.Role, error) {
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
func (r *RoleRepository) AssignUserRole(userID, roleID, companyID int64, branchID *int64) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, company_id, branch_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role_id, company_id, branch_id) DO NOTHING
	`

	_, err := r.db.Exec(query, userID, roleID, companyID, branchID)
	if err != nil {
		return fmt.Errorf("failed to assign user role: %w", err)
	}

	return nil
}

// RemoveUserRole removes a role from a user
func (r *RoleRepository) RemoveUserRole(userID, roleID int64) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.db.Exec(query, userID, roleID)
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

// GetUsersByRole retrieves users assigned to a specific role
func (r *RoleRepository) GetUsersByRole(roleID int64, limit int) ([]*models.User, error) {
	query := `
		SELECT DISTINCT u.id, u.name, u.email, u.user_identity, u.password_hash, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN user_roles ur ON u.id = ur.user_id
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
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.UserIdentity,
			&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserRoles retrieves roles assigned to a specific user
func (r *RoleRepository) GetUserRoles(userID int64) ([]*models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description,
			&role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateRoleModules updates module permissions for a role
func (r *RoleRepository) UpdateRoleModules(roleID int64, modules []models.RoleModule) error {
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
