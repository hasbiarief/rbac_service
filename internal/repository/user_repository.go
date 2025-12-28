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

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.UserIdentity,
		&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUserIdentity retrieves a user by user_identity
func (r *UserRepository) GetByUserIdentity(userIdentity string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE user_identity = $1 AND is_active = true
	`

	err := r.db.QueryRow(query, userIdentity).Scan(
		&user.ID, &user.Name, &user.Email, &user.UserIdentity,
		&user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
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

// Delete soft deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetUserRoles retrieves user roles
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
