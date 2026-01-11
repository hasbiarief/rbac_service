package repository

import (
	"database/sql"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/query"
	"strings"
)

type UnitRepository struct {
	*model.Repository
	db *sql.DB
}

func NewUnitRepository(db *sql.DB) interfaces.UnitRepositoryInterface {
	return &UnitRepository{
		Repository: &model.Repository{},
		db:         db,
	}
}

// GetAll retrieves units with optional filtering
func (r *UnitRepository) GetAll(branchID *int64, limit, offset int, search string, isActive *bool) ([]*models.UnitWithBranch, error) {
	qb := query.NewQueryBuilder(`
		SELECT 
			u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, 
			u.level, u.path, u.is_active, u.created_at, u.updated_at,
			b.name as branch_name, b.code as branch_code,
			c.name as company_name, c.code as company_code
		FROM units u
		JOIN branches b ON u.branch_id = b.id
		JOIN companies c ON b.company_id = c.id
	`)

	if branchID != nil {
		qb.AddCondition("u.branch_id = $%d", *branchID)
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		qb.AddCondition("(LOWER(u.name) LIKE $%d OR LOWER(u.code) LIKE $%d)", searchPattern)
		qb.AddCondition("", searchPattern) // Add the second search parameter
	}

	if isActive != nil {
		qb.AddCondition("u.is_active = $%d", *isActive)
	}

	qb.AddOrderBy("u.branch_id, u.level, u.name")
	qb.AddLimit(limit)
	qb.AddOffset(offset)

	query, args := qb.Build()
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*models.UnitWithBranch
	for rows.Next() {
		unit := &models.UnitWithBranch{}
		err := rows.Scan(
			&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
			&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
			&unit.BranchName, &unit.BranchCode, &unit.CompanyName, &unit.CompanyCode,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}

// GetByID retrieves a unit by ID with branch and company information
func (r *UnitRepository) GetByID(id int64) (*models.UnitWithBranch, error) {
	query := `
		SELECT 
			u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, 
			u.level, u.path, u.is_active, u.created_at, u.updated_at,
			b.name as branch_name, b.code as branch_code,
			c.name as company_name, c.code as company_code
		FROM units u
		JOIN branches b ON u.branch_id = b.id
		JOIN companies c ON b.company_id = c.id
		WHERE u.id = $1
	`

	unit := &models.UnitWithBranch{}
	err := r.db.QueryRow(query, id).Scan(
		&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
		&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		&unit.BranchName, &unit.BranchCode, &unit.CompanyName, &unit.CompanyCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return unit, nil
}

// GetByCode retrieves a unit by branch ID and code
func (r *UnitRepository) GetByCode(branchID int64, code string) (*models.Unit, error) {
	query := `
		SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
		FROM units 
		WHERE branch_id = $1 AND code = $2
	`

	unit := &models.Unit{}
	err := r.db.QueryRow(query, branchID, code).Scan(
		&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
		&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return unit, nil
}

// Create creates a new unit
func (r *UnitRepository) Create(unit *models.Unit) error {
	query := `
		INSERT INTO units (branch_id, parent_id, name, code, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, level, path, created_at, updated_at
	`

	return r.db.QueryRow(query,
		unit.BranchID, unit.ParentID, unit.Name, unit.Code, unit.Description, unit.IsActive,
	).Scan(&unit.ID, &unit.Level, &unit.Path, &unit.CreatedAt, &unit.UpdatedAt)
}

// Update updates an existing unit
func (r *UnitRepository) Update(unit *models.Unit) error {
	query := `
		UPDATE units 
		SET parent_id = $2, name = $3, code = $4, description = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING level, path, updated_at
	`

	return r.db.QueryRow(query,
		unit.ID, unit.ParentID, unit.Name, unit.Code, unit.Description, unit.IsActive,
	).Scan(&unit.Level, &unit.Path, &unit.UpdatedAt)
}

// Delete soft deletes a unit
func (r *UnitRepository) Delete(id int64) error {
	query := `UPDATE units SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetHierarchy retrieves unit hierarchy for a branch
func (r *UnitRepository) GetHierarchy(branchID int64) ([]*models.UnitHierarchy, error) {
	query := `
		WITH RECURSIVE unit_tree AS (
			-- Root units (no parent)
			SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at, 0 as depth
			FROM units 
			WHERE branch_id = $1 AND parent_id IS NULL AND is_active = true
			
			UNION ALL
			
			-- Child units
			SELECT u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, u.level, u.path, u.is_active, u.created_at, u.updated_at, ut.depth + 1
			FROM units u
			INNER JOIN unit_tree ut ON u.parent_id = ut.id
			WHERE u.is_active = true
		)
		SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
		FROM unit_tree
		ORDER BY level, name
	`

	rows, err := r.db.Query(query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*models.UnitHierarchy
	for rows.Next() {
		unit := &models.UnitHierarchy{}
		err := rows.Scan(
			&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
			&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	// Build hierarchy structure
	return r.buildHierarchy(units), nil
}

// buildHierarchy builds the hierarchical structure from flat list
func (r *UnitRepository) buildHierarchy(units []*models.UnitHierarchy) []*models.UnitHierarchy {
	unitMap := make(map[int64]*models.UnitHierarchy)
	var roots []*models.UnitHierarchy

	// Create map for quick lookup
	for _, unit := range units {
		unitMap[unit.ID] = unit
	}

	// Build hierarchy
	for _, unit := range units {
		if unit.ParentID == nil {
			roots = append(roots, unit)
		} else {
			parent := unitMap[*unit.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, *unit)
			}
		}
	}

	return roots
}

// GetChildren retrieves direct children of a unit
func (r *UnitRepository) GetChildren(parentID int64) ([]*models.Unit, error) {
	query := `
		SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
		FROM units 
		WHERE parent_id = $1 AND is_active = true
		ORDER BY name
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*models.Unit
	for rows.Next() {
		unit := &models.Unit{}
		err := rows.Scan(
			&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
			&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}

// GetParentPath retrieves the path from root to unit
func (r *UnitRepository) GetParentPath(unitID int64) ([]*models.Unit, error) {
	query := `
		WITH RECURSIVE parent_path AS (
			SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
			FROM units 
			WHERE id = $1
			
			UNION ALL
			
			SELECT u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, u.level, u.path, u.is_active, u.created_at, u.updated_at
			FROM units u
			INNER JOIN parent_path pp ON u.id = pp.parent_id
		)
		SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
		FROM parent_path
		ORDER BY level
	`

	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*models.Unit
	for rows.Next() {
		unit := &models.Unit{}
		err := rows.Scan(
			&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
			&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}

// GetWithStats retrieves unit with statistics
func (r *UnitRepository) GetWithStats(id int64) (*models.UnitWithStats, error) {
	query := `
		SELECT 
			u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, 
			u.level, u.path, u.is_active, u.created_at, u.updated_at,
			COALESCE(user_count.total, 0) as total_users,
			COALESCE(sub_unit_count.total, 0) as total_sub_units,
			COALESCE(role_count.total, 0) as total_roles
		FROM units u
		LEFT JOIN (
			SELECT unit_id, COUNT(*) as total 
			FROM user_roles 
			WHERE unit_id IS NOT NULL 
			GROUP BY unit_id
		) user_count ON u.id = user_count.unit_id
		LEFT JOIN (
			SELECT parent_id, COUNT(*) as total 
			FROM units 
			WHERE parent_id IS NOT NULL AND is_active = true 
			GROUP BY parent_id
		) sub_unit_count ON u.id = sub_unit_count.parent_id
		LEFT JOIN (
			SELECT unit_id, COUNT(*) as total 
			FROM unit_roles 
			GROUP BY unit_id
		) role_count ON u.id = role_count.unit_id
		WHERE u.id = $1
	`

	unit := &models.UnitWithStats{}
	err := r.db.QueryRow(query, id).Scan(
		&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
		&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		&unit.TotalUsers, &unit.TotalSubUnits, &unit.TotalRoles,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return unit, nil
}

// GetByBranch retrieves all units in a branch
func (r *UnitRepository) GetByBranch(branchID int64) ([]*models.Unit, error) {
	query := `
		SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
		FROM units 
		WHERE branch_id = $1 AND is_active = true
		ORDER BY level, name
	`

	rows, err := r.db.Query(query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*models.Unit
	for rows.Next() {
		unit := &models.Unit{}
		err := rows.Scan(
			&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
			&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}

// CountByBranch counts units in a branch
func (r *UnitRepository) CountByBranch(branchID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM units WHERE branch_id = $1 AND is_active = true`

	var count int64
	err := r.db.QueryRow(query, branchID).Scan(&count)
	return count, err
}

// ExistsByCode checks if a unit code exists in a branch
func (r *UnitRepository) ExistsByCode(branchID int64, code string, excludeID *int64) (bool, error) {
	qb := query.NewQueryBuilder("SELECT COUNT(*) FROM units")
	qb.AddCondition("branch_id = $%d", branchID)
	qb.AddCondition("code = $%d", code)

	if excludeID != nil {
		qb.AddCondition("id != $%d", *excludeID)
	}

	query, args := qb.Build()
	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count > 0, err
}

// HasChildren checks if a unit has children
func (r *UnitRepository) HasChildren(unitID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM units WHERE parent_id = $1 AND is_active = true`

	var count int64
	err := r.db.QueryRow(query, unitID).Scan(&count)
	return count > 0, err
}

// HasUsers checks if a unit has assigned users
func (r *UnitRepository) HasUsers(unitID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM user_roles WHERE unit_id = $1`

	var count int64
	err := r.db.QueryRow(query, unitID).Scan(&count)
	return count > 0, err
}
