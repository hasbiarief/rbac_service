package repository

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
)

type UserRepository struct {
	*model.Repository
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetDB returns the database connection
func (r *UserRepository) GetDB() *sql.DB {
	return r.db
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (name, email, user_identity, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, user.Name, user.Email, user.UserIdentity, user.PasswordHash, user.IsActive).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID (excluding soft deleted)
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.UserIdentity,
		&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email (excluding soft deleted)
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE email = $1 AND is_active = true AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.UserIdentity,
		&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUserIdentity retrieves a user by user_identity (excluding soft deleted)
func (r *UserRepository) GetByUserIdentity(userIdentity string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE user_identity = $1 AND is_active = true AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, userIdentity).Scan(
		&user.ID, &user.Name, &user.Email, &user.UserIdentity,
		&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query, values := r.BuildUpdateQuery(user, user.ID)

	err := r.db.QueryRow(query, values...).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete soft deletes a user by setting deleted_at timestamp
func (r *UserRepository) Delete(id int64) error {
	query := `
		UPDATE users 
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// GetAll retrieves all users with pagination and filtering (excluding soft deleted)
func (r *UserRepository) GetAll(limit, offset int, search string, isActive *bool) ([]*models.User, error) {
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at
		FROM users 
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Add active filter
	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Add ordering and pagination
	query += " ORDER BY created_at DESC"
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
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.UserIdentity,
			&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// GetCount returns total count of users with filtering (excluding soft deleted)
func (r *UserRepository) GetCount(search string, isActive *bool) (int64, error) {
	query := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Add active filter
	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	return count, nil
}
func (r *UserRepository) GetUserRoles(userID int64) ([]string, error) {
	query := `
		SELECT r.name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetUserModulesWithSubscription retrieves user modules filtered by subscription
func (r *UserRepository) GetUserModulesWithSubscription(userID int64) ([]string, error) {
	// First, get user's company ID
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		return r.getUserBasicModules(userID)
	}

	// Query modules with subscription filtering
	query := `
		SELECT DISTINCT m.url
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
		JOIN subscriptions s ON pm.plan_id = s.plan_id
		WHERE ur.user_id = $1 
			AND rm.can_read = true
			AND m.is_active = true
			AND s.company_id = $2
			AND s.status = 'active'
			AND s.end_date > CURRENT_DATE
		ORDER BY m.url
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		return r.getUserBasicModules(userID)
	}
	defer rows.Close()

	var modules []string
	for rows.Next() {
		var moduleURL string
		if err := rows.Scan(&moduleURL); err != nil {
			continue
		}
		modules = append(modules, moduleURL)
	}

	if len(modules) == 0 {
		return r.getUserBasicModules(userID)
	}

	return modules, nil
}

// GetUserModulesGroupedWithSubscription retrieves user modules grouped by category with subscription filtering
func (r *UserRepository) GetUserModulesGroupedWithSubscription(userID int64) (map[string][][]string, error) {
	// First, get user's company ID
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		return r.getUserBasicModulesGrouped(userID)
	}

	// Query modules with subscription filtering - include icon, description, and parent_id for sorting
	query := `
		SELECT DISTINCT m.name, m.url, m.icon, m.description, m.category, m.parent_id,
			CASE WHEN m.parent_id IS NULL THEN 0 ELSE 1 END as sort_order
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
		JOIN subscriptions s ON pm.plan_id = s.plan_id
		WHERE ur.user_id = $1 
			AND rm.can_read = true
			AND m.is_active = true
			AND s.company_id = $2
			AND s.status = 'active'
			AND s.end_date > CURRENT_DATE
		ORDER BY m.category, sort_order, m.name
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		return r.getUserBasicModulesGrouped(userID)
	}
	defer rows.Close()

	modules := make(map[string][][]string)
	for rows.Next() {
		var moduleName, moduleURL, moduleIcon, moduleDescription, category string
		var parentID *int64
		var sortOrder int
		if err := rows.Scan(&moduleName, &moduleURL, &moduleIcon, &moduleDescription, &category, &parentID, &sortOrder); err != nil {
			continue
		}

		// Add module to category with [name, url, icon, description]
		modules[category] = append(modules[category], []string{moduleName, moduleURL, moduleIcon, moduleDescription})
	}

	if len(modules) == 0 {
		return r.getUserBasicModulesGrouped(userID)
	}

	return modules, nil
}

// getUserBasicModules returns basic tier modules only
func (r *UserRepository) getUserBasicModules(userID int64) ([]string, error) {
	query := `
		SELECT DISTINCT m.url
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		WHERE ur.user_id = $1 
			AND rm.can_read = true
			AND m.is_active = true
			AND (m.subscription_tier = 'basic' OR m.subscription_tier IS NULL)
		ORDER BY m.url
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic modules: %w", err)
	}
	defer rows.Close()

	var modules []string
	for rows.Next() {
		var moduleURL string
		if err := rows.Scan(&moduleURL); err != nil {
			continue
		}
		modules = append(modules, moduleURL)
	}

	return modules, nil
}

// getUserBasicModulesGrouped returns basic tier modules grouped by category
func (r *UserRepository) getUserBasicModulesGrouped(userID int64) (map[string][][]string, error) {
	query := `
		SELECT DISTINCT m.name, m.url, m.icon, m.description, m.category, m.parent_id,
			CASE WHEN m.parent_id IS NULL THEN 0 ELSE 1 END as sort_order
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		WHERE ur.user_id = $1 
			AND rm.can_read = true
			AND m.is_active = true
			AND (m.subscription_tier = 'basic' OR m.subscription_tier IS NULL)
		ORDER BY m.category, sort_order, m.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic modules grouped: %w", err)
	}
	defer rows.Close()

	modules := make(map[string][][]string)
	for rows.Next() {
		var moduleName, moduleURL, moduleIcon, moduleDescription, category string
		var parentID *int64
		var sortOrder int
		if err := rows.Scan(&moduleName, &moduleURL, &moduleIcon, &moduleDescription, &category, &parentID, &sortOrder); err != nil {
			continue
		}

		// Add module to category with [name, url, icon, description]
		modules[category] = append(modules[category], []string{moduleName, moduleURL, moduleIcon, moduleDescription})
	}

	return modules, nil
}
