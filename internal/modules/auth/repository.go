package auth

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/pkg/model"
	"time"
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
		// If no company found, return empty modules (no fallback)
		return make(map[string][][]string), nil
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
			AND (s.billing_cycle = 'lifetime' OR s.end_date >= CURRENT_DATE)
		ORDER BY m.category, sort_order, m.name
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		// If subscription query fails, return empty modules (no fallback)
		return make(map[string][][]string), nil
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

	// If no modules found (expired subscription), return empty instead of fallback
	if len(modules) == 0 {
		return make(map[string][][]string), nil
	}

	return modules, nil
}

// GetUserApplicationsWithModules retrieves user applications with modules grouped by category
func (r *Repository) GetUserApplicationsWithModules(userID int64) (map[string]interface{}, error) {
	// Get user's company ID
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		return make(map[string]interface{}), nil
	}

	// Query applications with modules and subscription filtering
	query := `
		SELECT DISTINCT 
			a.id as app_id, a.name as app_name, a.code as app_code, 
			a.icon as app_icon, a.url as app_url, a.sort_order as app_sort_order,
			m.name as module_name, m.url as module_url, m.icon as module_icon, 
			m.description as module_description, m.category as module_category
		FROM user_roles ur
		JOIN role_modules rm ON ur.role_id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		JOIN applications a ON m.application_id = a.id
		JOIN plan_modules pm ON m.id = pm.module_id AND pm.is_included = true
		JOIN plan_applications pa ON a.id = pa.application_id AND pa.is_included = true
		JOIN subscriptions s ON pm.plan_id = s.plan_id AND pa.plan_id = s.plan_id
		WHERE ur.user_id = $1 
			AND rm.can_read = true
			AND m.is_active = true
			AND a.is_active = true
			AND s.company_id = $2
			AND s.status = 'active'
			AND (s.billing_cycle = 'lifetime' OR s.end_date >= CURRENT_DATE)
		ORDER BY a.sort_order, a.name, m.category, m.name
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		return make(map[string]interface{}), nil
	}
	defer rows.Close()

	applications := make(map[string]interface{})

	for rows.Next() {
		var appID int64
		var appName, appCode, appIcon, appURL string
		var appSortOrder int
		var moduleName, moduleURL, moduleIcon, moduleDescription, moduleCategory string

		if err := rows.Scan(&appID, &appName, &appCode, &appIcon, &appURL, &appSortOrder,
			&moduleName, &moduleURL, &moduleIcon, &moduleDescription, &moduleCategory); err != nil {
			continue
		}

		// Initialize application if not exists
		if applications[appName] == nil {
			applications[appName] = map[string]interface{}{
				"id":         appID,
				"name":       appName,
				"code":       appCode,
				"icon":       appIcon,
				"url":        appURL,
				"sort_order": appSortOrder,
				"modules":    make(map[string][][]string),
			}
		}

		// Add module to application
		app := applications[appName].(map[string]interface{})
		modules := app["modules"].(map[string][][]string)
		modules[moduleCategory] = append(modules[moduleCategory],
			[]string{moduleName, moduleURL, moduleIcon, moduleDescription})
		app["modules"] = modules
	}

	return applications, nil
}

// GetUserModulesWithSubscription retrieves user module URLs
func (r *Repository) GetUserModulesWithSubscription(userID int64) ([]string, error) {
	// Get user's company ID
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		// If no company found, return empty modules (no fallback)
		return []string{}, nil
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
			AND (s.billing_cycle = 'lifetime' OR s.end_date >= CURRENT_DATE)
		ORDER BY m.url
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		// If subscription query fails, return empty modules (no fallback)
		return []string{}, nil
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

	// If no modules found (expired subscription), return empty instead of fallback
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

// GetUserSubscriptionInfo retrieves user's company subscription information
func (r *Repository) GetUserSubscriptionInfo(userID int64) (map[string]interface{}, error) {
	// Get user's company ID first
	var companyID int64
	err := r.db.QueryRow("SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1", userID).Scan(&companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user company: %w", err)
	}

	// Get subscription information with plan details
	query := `
		SELECT 
			s.id as subscription_id,
			s.company_id,
			c.name as company_name,
			s.plan_id,
			sp.name as plan_name,
			sp.description as plan_description,
			s.price,
			s.billing_cycle,
			sp.max_users,
			sp.max_branches,
			s.status,
			s.start_date,
			s.end_date,
			s.created_at as subscription_created_at,
			s.updated_at as subscription_updated_at,
			CASE 
				WHEN s.billing_cycle = 'lifetime' THEN 'lifetime'
				WHEN s.end_date < CURRENT_DATE THEN 'expired'
				WHEN s.end_date = CURRENT_DATE THEN 'expiring_today'
				WHEN s.end_date <= CURRENT_DATE + INTERVAL '7 days' THEN 'expiring_soon'
				ELSE 'active'
			END as computed_status,
			CASE 
				WHEN s.billing_cycle = 'lifetime' THEN NULL
				ELSE (s.end_date - CURRENT_DATE)
			END as days_remaining
		FROM subscriptions s
		JOIN subscription_plans sp ON s.plan_id = sp.id
		JOIN companies c ON s.company_id = c.id
		WHERE s.company_id = $1 
			AND s.status = 'active'
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	var subscriptionID, planID int64
	var companyName, planName, planDescription, billingCycle, status, computedStatus string
	var price float64
	var maxUsers, maxBranches, daysRemaining *int64
	var startDate, endDate, subscriptionCreatedAt, subscriptionUpdatedAt string

	err = r.db.QueryRow(query, companyID).Scan(
		&subscriptionID, &companyID, &companyName, &planID, &planName, &planDescription,
		&price, &billingCycle, &maxUsers, &maxBranches,
		&status, &startDate, &endDate, &subscriptionCreatedAt, &subscriptionUpdatedAt,
		&computedStatus, &daysRemaining,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Get company name even if no subscription
			var companyNameOnly string
			companyQuery := "SELECT name FROM companies WHERE id = $1"
			if err := r.db.QueryRow(companyQuery, companyID).Scan(&companyNameOnly); err != nil {
				companyNameOnly = "Unknown Company"
			}

			return map[string]interface{}{
				"has_subscription": false,
				"company_id":       companyID,
				"company_name":     companyNameOnly,
				"message":          "No active subscription found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get subscription info: %w", err)
	}

	subscriptionInfo := map[string]interface{}{
		"has_subscription": true,
		"subscription": map[string]interface{}{
			"id":           subscriptionID,
			"company_id":   companyID,
			"company_name": companyName,
			"plan": map[string]interface{}{
				"id":            planID,
				"name":          planName,
				"description":   planDescription,
				"price":         price,
				"billing_cycle": billingCycle,
			},
			"limits": map[string]interface{}{
				"max_users":    maxUsers,
				"max_branches": maxBranches,
			},
			"status":          status,
			"computed_status": computedStatus,
			"start_date":      startDate,
			"end_date":        endDate,
			"created_at":      subscriptionCreatedAt,
			"updated_at":      subscriptionUpdatedAt,
		},
	}

	// Add days remaining if not null
	if daysRemaining != nil {
		subscriptionInfo["subscription"].(map[string]interface{})["days_remaining"] = *daysRemaining
	}

	return subscriptionInfo, nil
}

// GetUserProfileByApplication retrieves user profile with modules for specific application
func (r *Repository) GetUserProfileByApplication(userIdentity, applicationCode string) (map[string]interface{}, error) {
	// First get user by user_identity
	user, err := r.GetByUserIdentity(userIdentity)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get application by code
	var applicationID int64
	var applicationName, applicationIcon, applicationURL string
	var applicationSortOrder int

	appQuery := `
		SELECT id, name, icon, url, sort_order
		FROM applications 
		WHERE code = $1 AND is_active = true
	`
	err = r.db.QueryRow(appQuery, applicationCode).Scan(
		&applicationID, &applicationName, &applicationIcon, &applicationURL, &applicationSortOrder,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("application not found")
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Get user's company ID for subscription check
	var companyID int64
	companyQuery := `
		SELECT COALESCE(ur.company_id, 0) as company_id
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		WHERE u.id = $1
		LIMIT 1
	`
	err = r.db.QueryRow(companyQuery, user.ID).Scan(&companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user company: %w", err)
	}

	// Get user role assignments for this application
	// Filter roles based on their assigned application_id
	roleQuery := `
		SELECT DISTINCT
			ur.id as assignment_id,
			ur.role_id,
			ro.name as role_name,
			ro.description as role_description,
			CASE 
				WHEN ur.unit_id IS NOT NULL THEN 'unit'
				WHEN ur.branch_id IS NOT NULL THEN 'branch'
				ELSE 'company'
			END as assignment_level,
			ur.company_id,
			c.name as company_name,
			ur.branch_id,
			b.name as branch_name,
			ur.unit_id,
			u.name as unit_name
		FROM user_roles ur
		JOIN roles ro ON ur.role_id = ro.id
		LEFT JOIN companies c ON ur.company_id = c.id
		LEFT JOIN branches b ON ur.branch_id = b.id
		LEFT JOIN units u ON ur.unit_id = u.id
		WHERE ur.user_id = $1 
		AND ro.is_active = true
		AND ro.application_id = $2
		ORDER BY ur.id
	`

	rows, err := r.db.Query(roleQuery, user.ID, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roleAssignments []map[string]interface{}
	roleMap := make(map[int64]bool) // To avoid duplicates

	for rows.Next() {
		var assignmentID, roleID, companyID, branchID, unitID sql.NullInt64
		var roleName, roleDescription, assignmentLevel, companyName, branchName, unitName sql.NullString

		err := rows.Scan(
			&assignmentID, &roleID, &roleName, &roleDescription, &assignmentLevel,
			&companyID, &companyName, &branchID, &branchName, &unitID, &unitName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role assignment: %w", err)
		}

		// Skip if we already processed this role assignment
		if assignmentID.Valid && roleMap[assignmentID.Int64] {
			continue
		}
		if assignmentID.Valid {
			roleMap[assignmentID.Int64] = true
		}

		assignment := map[string]interface{}{
			"assignment_id":    assignmentID.Int64,
			"role_id":          roleID.Int64,
			"role_name":        roleName.String,
			"role_description": roleDescription.String,
			"assignment_level": assignmentLevel.String,
		}

		if companyID.Valid {
			assignment["company_id"] = companyID.Int64
			assignment["company_name"] = companyName.String
		}
		if branchID.Valid {
			assignment["branch_id"] = branchID.Int64
			assignment["branch_name"] = branchName.String
		}
		if unitID.Valid {
			assignment["unit_id"] = unitID.Int64
			assignment["unit_name"] = unitName.String
		}

		roleAssignments = append(roleAssignments, assignment)
	}

	// Get modules for this application with subscription filtering
	modulesQuery := `
		SELECT DISTINCT
			m.id, m.name, m.url, m.icon, m.description, m.category
		FROM modules m
		JOIN role_modules rm ON m.id = rm.module_id
		JOIN user_roles ur ON rm.role_id = ur.role_id
		JOIN plan_modules pm ON m.id = pm.module_id
		JOIN subscriptions s ON pm.plan_id = s.plan_id
		WHERE ur.user_id = $1 
		AND m.application_id = $2
		AND s.company_id = $3
		AND s.status = 'active'
		AND (s.end_date >= CURRENT_DATE OR s.end_date IS NULL)
		AND pm.is_included = true
		AND m.is_active = true
		AND rm.can_read = true
		ORDER BY m.category, m.name
	`

	moduleRows, err := r.db.Query(modulesQuery, user.ID, applicationID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user modules: %w", err)
	}
	defer moduleRows.Close()

	// Group modules by category
	modulesByCategory := make(map[string][][]string)

	for moduleRows.Next() {
		var moduleID int64
		var moduleName, moduleURL, moduleIcon, moduleDescription, moduleCategory string

		err := moduleRows.Scan(
			&moduleID, &moduleName, &moduleURL, &moduleIcon, &moduleDescription, &moduleCategory,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}

		moduleInfo := []string{moduleName, moduleURL, moduleIcon, moduleDescription}
		modulesByCategory[moduleCategory] = append(modulesByCategory[moduleCategory], moduleInfo)
	}

	// Build application structure
	applicationData := map[string]interface{}{
		"id":               applicationID,
		"name":             applicationName,
		"code":             applicationCode,
		"icon":             applicationIcon,
		"url":              applicationURL,
		"sort_order":       applicationSortOrder,
		"role_assignments": roleAssignments,
		"modules":          modulesByCategory,
	}

	// Build user profile response
	userProfile := map[string]interface{}{
		"name":          user.Name,
		"user_identity": user.UserIdentity,
		"email":         user.Email,
		"is_active":     user.IsActive,
		"total_roles":   len(roleAssignments),
		"created_at":    user.CreatedAt.Format(time.RFC3339),
		"updated_at":    user.UpdatedAt.Format(time.RFC3339),
		"applications": map[string]interface{}{
			applicationCode: applicationData,
		},
	}

	return userProfile, nil
}
