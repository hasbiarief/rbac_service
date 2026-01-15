package module

import (
	"database/sql"
	"fmt"
	// removed
	"gin-scalable-api/pkg/model"
)

type ModuleRepository struct {
	*model.Repository
	db *sql.DB
}

func NewModuleRepository(db *sql.DB) *ModuleRepository {
	return &ModuleRepository{
		Repository: model.NewRepository(db),
		db:         db,
	}
}

// GetAll retrieves all modules with pagination and filtering
func (r *ModuleRepository) GetAll(limit, offset int, search string, category, subscriptionTier string, parentID *int64, isActive *bool) ([]*Module, error) {
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
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

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if subscriptionTier != "" {
		query += fmt.Sprintf(" AND subscription_tier = $%d", argIndex)
		args = append(args, subscriptionTier)
		argIndex++
	}

	if parentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argIndex)
		args = append(args, *parentID)
		argIndex++
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	query += " ORDER BY category, subscription_tier, name"

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
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetByID retrieves a module by ID
func (r *ModuleRepository) GetByID(id int64) (*Module, error) {
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
		WHERE id = $1
	`

	module := &Module{}
	err := r.db.QueryRow(query, id).Scan(
		&module.ID, &module.Category, &module.Name, &module.URL,
		&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
		&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module not found")
		}
		return nil, fmt.Errorf("failed to get module: %w", err)
	}

	return module, nil
}

// Create creates a new module
func (r *ModuleRepository) Create(module *Module) error {
	query, values := r.BuildInsertQuery(module)

	err := r.db.QueryRow(query, values...).Scan(
		&module.ID, &module.CreatedAt, &module.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create module: %w", err)
	}

	return nil
}

// Update updates a module
func (r *ModuleRepository) Update(module *Module) error {
	query, values := r.BuildUpdateQuery(module, module.ID)

	err := r.db.QueryRow(query, values...).Scan(&module.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update module: %w", err)
	}

	return nil
}

// Delete soft deletes a module
func (r *ModuleRepository) Delete(id int64) error {
	query := `UPDATE modules SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("module not found")
	}

	return nil
}

// GetChildren retrieves child modules of a parent module
func (r *ModuleRepository) GetChildren(parentID int64) ([]*Module, error) {
	// print parentID
	println(parentID)
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
		WHERE parent_id = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get child modules: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetAncestors retrieves ancestor modules of a module with user filtering
func (r *ModuleRepository) GetAncestors(moduleID int64, userID int64) ([]*Module, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
			FROM modules
			WHERE id = $1
			
			UNION ALL
			
			SELECT m.id, m.category, m.name, m.url, m.icon, m.description, m.parent_id, m.subscription_tier, m.is_active, m.created_at, m.updated_at
			FROM modules m
			INNER JOIN ancestors a ON m.id = a.parent_id
		)
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM ancestors
		WHERE id != $1
		ORDER BY subscription_tier
	`

	rows, err := r.db.Query(query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ancestor modules: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetTree retrieves modules in tree structure
func (r *ModuleRepository) GetTree(category string) ([]*Module, error) {
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
		WHERE is_active = true
	`
	args := []interface{}{}

	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}

	query += " ORDER BY subscription_tier, name"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get module tree: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetTreeStructure retrieves modules in hierarchical tree structure
// GetTreeStructure retrieves modules in hierarchical tree structure with user filtering
func (r *ModuleRepository) GetTreeStructure(category string, userID int64) ([]*Module, error) {
	// For now, just return flat modules - tree structure building is done in service
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
		WHERE is_active = true
	`
	args := []interface{}{}

	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}

	query += " ORDER BY subscription_tier, name"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get module tree: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetTreeByParentName retrieves modules tree structure by parent module name with user filtering
func (r *ModuleRepository) GetTreeByParentName(parentName string, userID int64) ([]*Module, error) {
	// First find the parent module
	parentQuery := `
		SELECT id FROM modules 
		WHERE name = $1 AND is_active = true
		LIMIT 1
	`

	var parentID int64
	err := r.db.QueryRow(parentQuery, parentName).Scan(&parentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*Module{}, nil // Return empty array if parent not found
		}
		return nil, fmt.Errorf("failed to find parent module: %w", err)
	}

	// Get all descendants of this parent using recursive CTE
	query := `
		WITH RECURSIVE module_tree AS (
			-- Base case: the parent module itself
			SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at, 0 as level
			FROM modules 
			WHERE id = $1 AND is_active = true
			
			UNION ALL
			
			-- Recursive case: all descendants
			SELECT m.id, m.category, m.name, m.url, m.icon, m.description, m.parent_id, m.subscription_tier, m.is_active, m.created_at, m.updated_at, mt.level + 1
			FROM modules m
			INNER JOIN module_tree mt ON m.parent_id = mt.id
			WHERE m.is_active = true
		)
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM module_tree 
		ORDER BY level, name
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module tree by parent name: %w", err)
	}
	defer rows.Close()

	var modules []*Module
	for rows.Next() {
		module := &Module{}
		err := rows.Scan(
			&module.ID, &module.Category, &module.Name, &module.URL,
			&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
			&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// GetUserModules retrieves modules accessible by a user with company filtering
func (r *ModuleRepository) GetUserModules(userID, companyID int64) ([]*UserModule, error) {
	query := `
		SELECT DISTINCT m.id, m.category, m.name, m.url, m.icon, m.description, m.parent_id, m.subscription_tier, m.is_active, m.created_at, m.updated_at,
			rm.can_read, rm.can_write, rm.can_delete
		FROM modules m
		JOIN role_modules rm ON m.id = rm.module_id
		JOIN user_roles ur ON rm.role_id = ur.role_id
		WHERE ur.user_id = $1 
			AND (ur.company_id = $2 OR $2 = 0)
			AND rm.can_read = true
			AND m.is_active = true
		ORDER BY m.category, m.subscription_tier, m.name
	`

	rows, err := r.db.Query(query, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user modules: %w", err)
	}
	defer rows.Close()

	var userModules []*UserModule
	for rows.Next() {
		userModule := &UserModule{}
		err := rows.Scan(
			&userModule.ID, &userModule.Category, &userModule.Name, &userModule.URL,
			&userModule.Icon, &userModule.Description, &userModule.ParentID, &userModule.SubscriptionTier,
			&userModule.IsActive, &userModule.CreatedAt, &userModule.UpdatedAt,
			&userModule.CanRead, &userModule.CanWrite, &userModule.CanDelete,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user module: %w", err)
		}
		userModules = append(userModules, userModule)
	}

	return userModules, nil
}

// CheckUserAccess checks if user has access to a specific module
func (r *ModuleRepository) CheckUserAccess(userID int64, moduleURL string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM modules m
			JOIN role_modules rm ON m.id = rm.module_id
			JOIN user_roles ur ON rm.role_id = ur.role_id
			WHERE ur.user_id = $1 
				AND m.url = $2
				AND rm.can_read = true
				AND m.is_active = true
		)
	`

	var hasAccess bool
	err := r.db.QueryRow(query, userID, moduleURL).Scan(&hasAccess)
	if err != nil {
		return false, fmt.Errorf("failed to check user access: %w", err)
	}

	return hasAccess, nil
}

// GetByURL retrieves a module by URL
func (r *ModuleRepository) GetByURL(url string) (*Module, error) {
	query := `
		SELECT id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at
		FROM modules
		WHERE url = $1
	`

	module := &Module{}
	err := r.db.QueryRow(query, url).Scan(
		&module.ID, &module.Category, &module.Name, &module.URL,
		&module.Icon, &module.Description, &module.ParentID, &module.SubscriptionTier,
		&module.IsActive, &module.CreatedAt, &module.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module not found")
		}
		return nil, fmt.Errorf("failed to get module: %w", err)
	}

	return module, nil
}

// Count returns total count of modules with filtering
func (r *ModuleRepository) Count(search string, category, subscriptionTier string, parentID *int64, isActive *bool) (int64, error) {
	query := "SELECT COUNT(*) FROM modules WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if subscriptionTier != "" {
		query += fmt.Sprintf(" AND subscription_tier = $%d", argIndex)
		args = append(args, subscriptionTier)
		argIndex++
	}

	if parentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argIndex)
		args = append(args, *parentID)
		argIndex++
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get module count: %w", err)
	}

	return count, nil
}
