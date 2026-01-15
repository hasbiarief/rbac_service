package auth

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/pkg/model"
)

type Repository struct {
	*model.Repository
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetByUserIdentity retrieves a user by user_identity for authentication
func (r *Repository) GetByUserIdentity(userIdentity string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE user_identity = $1 AND is_active = true AND deleted_at IS NULL
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

// GetByEmail retrieves a user by email for authentication
func (r *Repository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE email = $1 AND is_active = true AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, email).Scan(
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

// GetByID retrieves a user by ID for authentication
func (r *Repository) GetByID(id int64) (*User, error) {
	user := &User{}
	query := `
		SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
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

// GetByIDWithRoles retrieves a user with role assignments
func (r *Repository) GetByIDWithRoles(id int64) (map[string]interface{}, error) {
	// Get basic user info
	userQuery := `
		SELECT id, name, email, user_identity, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var userID int64
	var name, email string
	var userIdentity *string
	var isActive bool
	var createdAt, updatedAt string

	err := r.db.QueryRow(userQuery, id).Scan(
		&userID, &name, &email, &userIdentity, &isActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := map[string]interface{}{
		"id":         userID,
		"name":       name,
		"email":      email,
		"is_active":  isActive,
		"created_at": createdAt,
		"updated_at": updatedAt,
	}

	if userIdentity != nil {
		user["user_identity"] = *userIdentity
	}

	// Get role assignments
	rolesQuery := `
		SELECT 
			ur.id as assignment_id,
			r.id as role_id,
			r.name as role_name,
			r.description as role_description,
			ur.company_id,
			c.name as company_name,
			ur.branch_id,
			b.name as branch_name,
			ur.unit_id,
			u.name as unit_name,
			CASE 
				WHEN ur.unit_id IS NOT NULL THEN 'unit'
				WHEN ur.branch_id IS NOT NULL THEN 'branch'
				ELSE 'company'
			END as assignment_level
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN companies c ON ur.company_id = c.id
		LEFT JOIN branches b ON ur.branch_id = b.id
		LEFT JOIN units u ON ur.unit_id = u.id
		WHERE ur.user_id = $1
		ORDER BY ur.id
	`

	rows, err := r.db.Query(rolesQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	roleAssignments := make([]map[string]interface{}, 0)
	for rows.Next() {
		var assignmentID, roleID, companyID int64
		var branchID, unitID *int64
		var roleName, roleDescription, companyName string
		var branchName, unitName *string
		var assignmentLevel string

		err := rows.Scan(
			&assignmentID, &roleID, &roleName, &roleDescription,
			&companyID, &companyName, &branchID, &branchName,
			&unitID, &unitName, &assignmentLevel,
		)
		if err != nil {
			continue
		}

		assignment := map[string]interface{}{
			"assignment_id":    assignmentID,
			"role_id":          roleID,
			"role_name":        roleName,
			"role_description": roleDescription,
			"company_id":       companyID,
			"company_name":     companyName,
			"assignment_level": assignmentLevel,
		}

		if branchID != nil {
			assignment["branch_id"] = *branchID
			if branchName != nil {
				assignment["branch_name"] = *branchName
			}
		}

		if unitID != nil {
			assignment["unit_id"] = *unitID
			if unitName != nil {
				assignment["unit_name"] = *unitName
			}
		}

		roleAssignments = append(roleAssignments, assignment)
	}

	user["role_assignments"] = roleAssignments
	user["total_roles"] = len(roleAssignments)

	return user, nil
}

// GetUserModulesGroupedWithSubscription retrieves user modules grouped by category
func (r *Repository) GetUserModulesGroupedWithSubscription(userID int64) (map[string][][]string, error) {
	// Get user's company ID
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		return r.getUserBasicModulesGrouped(userID)
	}

	// Query modules with subscription filtering
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

		modules[category] = append(modules[category], []string{moduleName, moduleURL, moduleIcon, moduleDescription})
	}

	if len(modules) == 0 {
		return r.getUserBasicModulesGrouped(userID)
	}

	return modules, nil
}

// GetUserModulesWithSubscription retrieves user module URLs
func (r *Repository) GetUserModulesWithSubscription(userID int64) ([]string, error) {
	// Get user's company ID
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

// GetUserRoles retrieves user role assignments
func (r *Repository) GetUserRoles(userID int64) ([]*UserRole, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.company_id, ur.branch_id, ur.unit_id
		FROM user_roles ur
		WHERE ur.user_id = $1
		ORDER BY ur.id
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
			&userRole.CompanyID, &userRole.BranchID, &userRole.UnitID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// Create creates a new user (for registration)
func (r *Repository) Create(user *User) error {
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

// Helper methods for basic tier modules
func (r *Repository) getUserBasicModules(userID int64) ([]string, error) {
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

func (r *Repository) getUserBasicModulesGrouped(userID int64) (map[string][][]string, error) {
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

		modules[category] = append(modules[category], []string{moduleName, moduleURL, moduleIcon, moduleDescription})
	}

	return modules, nil
}
