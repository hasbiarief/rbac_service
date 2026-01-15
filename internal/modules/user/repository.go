package user

import (
	"database/sql"
	"fmt"
	// removed - using local model
	"gin-scalable-api/pkg/model"
	"strings"
	"time"
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
func (r *UserRepository) Create(user *User) error {
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
func (r *UserRepository) GetByID(id int64) (*User, error) {
	user := &User{}
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
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
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
func (r *UserRepository) GetByUserIdentity(userIdentity string) (*User, error) {
	user := &User{}
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
func (r *UserRepository) Update(user *User) error {
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
func (r *UserRepository) GetAll(limit, offset int, search string, isActive *bool) ([]*User, error) {
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

	var users []*User
	for rows.Next() {
		user := &User{}
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

// Count returns total count of users with filtering (excluding soft deleted)
func (r *UserRepository) Count(search string, isActive *bool) (int64, error) {
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

// GetWithRoles retrieves a user with their roles
func (r *UserRepository) GetWithRoles(id int64) (*User, error) {
	// For now, just return the user without roles
	// This would need to be implemented properly with role data
	return r.GetByID(id)
}
func (r *UserRepository) GetUserRoles(userID int64) ([]*UserRole, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, ur.company_id, ur.created_at, ur.updated_at,
		       r.name as role_name, c.name as company_name
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		LEFT JOIN companies c ON ur.company_id = c.id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var userRoles []*UserRole
	for rows.Next() {
		userRole := &UserRole{}
		var roleName, companyName *string

		err := rows.Scan(
			&userRole.ID, &userRole.UserID, &userRole.RoleID, &userRole.CompanyID,
			&userRole.CreatedAt, &userRole.UpdatedAt, &roleName, &companyName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}

		if roleName != nil {
			userRole.RoleName = *roleName
		}
		if companyName != nil {
			userRole.CompanyName = *companyName
		}

		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
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

// GetByIDWithRoles retrieves a user by ID with complete role assignments
func (r *UserRepository) GetByIDWithRoles(id int64) (map[string]interface{}, error) {
	// Get basic user info
	userQuery := `
		SELECT id, name, email, user_identity, is_active, created_at, updated_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user map[string]interface{}
	var userID int64
	var name, email string
	var userIdentity *string
	var isActive bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(userQuery, id).Scan(
		&userID, &name, &email, &userIdentity, &isActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user = map[string]interface{}{
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

	// Get role assignments with company, branch, unit info
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

	var roleAssignments []map[string]interface{}
	// Initialize as empty slice to avoid nil
	roleAssignments = make([]map[string]interface{}, 0)
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

	// Ensure role_assignments is never nil
	if user["role_assignments"] == nil {
		user["role_assignments"] = []map[string]interface{}{}
	}

	return user, nil
}

// GetAllWithRoles retrieves all users with complete role assignments
func (r *UserRepository) GetAllWithRoles(limit, offset int, search string, isActive *bool) ([]map[string]interface{}, error) {
	// Build user query with filters
	userQuery := `
		SELECT id, name, email, user_identity, is_active, created_at, updated_at
		FROM users 
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	argIndex := 1

	// Add search filter
	if search != "" {
		userQuery += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Add active filter
	if isActive != nil {
		userQuery += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	// Add ordering and pagination
	userQuery += " ORDER BY created_at DESC"
	if limit > 0 {
		userQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++
	}
	if offset > 0 {
		userQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	rows, err := r.db.Query(userQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []map[string]interface{}
	var userIDs []int64

	// Collect users and their IDs
	for rows.Next() {
		var userID int64
		var name, email string
		var userIdentity *string
		var isActive bool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&userID, &name, &email, &userIdentity, &isActive, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
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

		users = append(users, user)
		userIDs = append(userIDs, userID)
	}

	// If no users found, return empty slice
	if len(users) == 0 {
		return users, nil
	}

	// Get role assignments for all users
	placeholders := make([]string, len(userIDs))
	roleArgs := make([]interface{}, len(userIDs))
	for i, userID := range userIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		roleArgs[i] = userID
	}

	rolesQuery := fmt.Sprintf(`
		SELECT 
			ur.user_id,
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
		WHERE ur.user_id IN (%s)
		ORDER BY ur.user_id, ur.id
	`, strings.Join(placeholders, ","))

	roleRows, err := r.db.Query(rolesQuery, roleArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer roleRows.Close()

	// Group role assignments by user ID
	userRoles := make(map[int64][]map[string]interface{})
	for roleRows.Next() {
		var userID, assignmentID, roleID, companyID int64
		var branchID, unitID *int64
		var roleName, roleDescription, companyName string
		var branchName, unitName *string
		var assignmentLevel string

		err := roleRows.Scan(
			&userID, &assignmentID, &roleID, &roleName, &roleDescription,
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

		userRoles[userID] = append(userRoles[userID], assignment)
	}

	// Add role assignments to users
	for i, user := range users {
		userID := user["id"].(int64)
		if roles, exists := userRoles[userID]; exists {
			users[i]["role_assignments"] = roles
			users[i]["total_roles"] = len(roles)
		} else {
			users[i]["role_assignments"] = []map[string]interface{}{}
			users[i]["total_roles"] = 0
		}
	}

	return users, nil
}
