package unit

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	GetAll(branchID *int64, limit, offset int, search string, isActive *bool) ([]*UnitWithBranch, error)
	Count(branchID *int64, search string, isActive *bool) (int64, error)
	GetByID(id int64) (*UnitWithBranch, error)
	GetHierarchy(branchID int64) ([]*Unit, error)
	GetWithStats(id int64) (*UnitWithStats, error)
	Create(unit *Unit) error
	Update(unit *Unit) error
	Delete(id int64) error

	// Unit Role methods
	AssignRole(unitID int64, roleID int64) error
	RemoveRole(unitID int64, roleID int64) error
	GetUnitRoles(unitID int64) ([]*UnitRole, error)

	// Permission methods
	GetUnitPermissions(unitID int64, roleID int64) ([]*UnitRoleModule, error)
	UpdatePermissions(unitRoleID int64, modules []UpdateUnitRoleModulePermission) error
	CopyPermissions(sourceUnitID int64, targetUnitID int64, roleID int64, overwrite bool) error
	GetUserEffectivePermissions(userID int64) ([]*UnitRoleModule, error)
}

type repository struct {
	db *sql.DB
}

type UnitWithStats struct {
	Unit
	TotalUsers    int `json:"total_users" db:"total_users"`
	TotalSubUnits int `json:"total_sub_units" db:"total_sub_units"`
	TotalRoles    int `json:"total_roles" db:"total_roles"`
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(branchID *int64, limit, offset int, search string, isActive *bool) ([]*UnitWithBranch, error) {
	query := `
		SELECT 
			u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, 
			u.level, u.path, u.is_active, u.created_at, u.updated_at,
			b.name as branch_name, b.code as branch_code,
			c.name as company_name, c.code as company_code
		FROM units u
		JOIN branches b ON u.branch_id = b.id
		JOIN companies c ON b.company_id = c.id
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	if branchID != nil {
		query += fmt.Sprintf(` AND u.branch_id = $%d`, argCount)
		args = append(args, *branchID)
		argCount++
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query += fmt.Sprintf(` AND (LOWER(u.name) LIKE $%d OR LOWER(u.code) LIKE $%d)`, argCount, argCount+1)
		args = append(args, searchPattern, searchPattern)
		argCount += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(` AND u.is_active = $%d`, argCount)
		args = append(args, *isActive)
		argCount++
	}

	query += ` ORDER BY u.branch_id, u.level, u.name`

	if limit > 0 {
		query += fmt.Sprintf(` LIMIT $%d`, argCount)
		args = append(args, limit)
		argCount++
	}

	if offset > 0 {
		query += fmt.Sprintf(` OFFSET $%d`, argCount)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*UnitWithBranch
	for rows.Next() {
		unit := &UnitWithBranch{}
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

func (r *repository) Count(branchID *int64, search string, isActive *bool) (int64, error) {
	query := `SELECT COUNT(*) FROM units WHERE 1=1`
	var args []interface{}
	argCount := 1

	if branchID != nil {
		query += fmt.Sprintf(` AND branch_id = $%d`, argCount)
		args = append(args, *branchID)
		argCount++
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query += fmt.Sprintf(` AND (LOWER(name) LIKE $%d OR LOWER(code) LIKE $%d)`, argCount, argCount+1)
		args = append(args, searchPattern, searchPattern)
		argCount += 2
	}

	if isActive != nil {
		query += fmt.Sprintf(` AND is_active = $%d`, argCount)
		args = append(args, *isActive)
	}

	var count int64
	err := r.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (r *repository) GetByID(id int64) (*UnitWithBranch, error) {
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

	unit := &UnitWithBranch{}
	err := r.db.QueryRow(query, id).Scan(
		&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
		&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		&unit.BranchName, &unit.BranchCode, &unit.CompanyName, &unit.CompanyCode,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return unit, err
}

func (r *repository) GetHierarchy(branchID int64) ([]*Unit, error) {
	query := `
		WITH RECURSIVE unit_tree AS (
			SELECT id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at
			FROM units 
			WHERE branch_id = $1 AND parent_id IS NULL AND is_active = true
			
			UNION ALL
			
			SELECT u.id, u.branch_id, u.parent_id, u.name, u.code, u.description, u.level, u.path, u.is_active, u.created_at, u.updated_at
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

	var units []*Unit
	for rows.Next() {
		unit := &Unit{}
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

func (r *repository) GetWithStats(id int64) (*UnitWithStats, error) {
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

	unit := &UnitWithStats{}
	err := r.db.QueryRow(query, id).Scan(
		&unit.ID, &unit.BranchID, &unit.ParentID, &unit.Name, &unit.Code, &unit.Description,
		&unit.Level, &unit.Path, &unit.IsActive, &unit.CreatedAt, &unit.UpdatedAt,
		&unit.TotalUsers, &unit.TotalSubUnits, &unit.TotalRoles,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return unit, err
}

func (r *repository) Create(unit *Unit) error {
	query := `
		INSERT INTO units (branch_id, parent_id, name, code, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, level, path, created_at, updated_at
	`

	return r.db.QueryRow(query,
		unit.BranchID, unit.ParentID, unit.Name, unit.Code, unit.Description, unit.IsActive,
	).Scan(&unit.ID, &unit.Level, &unit.Path, &unit.CreatedAt, &unit.UpdatedAt)
}

func (r *repository) Update(unit *Unit) error {
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

func (r *repository) Delete(id int64) error {
	query := `UPDATE units SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *repository) AssignRole(unitID int64, roleID int64) error {
	query := `INSERT INTO unit_roles (unit_id, role_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, unitID, roleID)
	return err
}

func (r *repository) RemoveRole(unitID int64, roleID int64) error {
	query := `DELETE FROM unit_roles WHERE unit_id = $1 AND role_id = $2`
	_, err := r.db.Exec(query, unitID, roleID)
	return err
}

func (r *repository) GetUnitRoles(unitID int64) ([]*UnitRole, error) {
	query := `
		SELECT ur.id, ur.unit_id, ur.role_id, ur.created_at, ur.updated_at,
			u.name as unit_name, r.name as role_name
		FROM unit_roles ur
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.unit_id = $1
	`

	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*UnitRole
	for rows.Next() {
		role := &UnitRole{}
		err := rows.Scan(&role.ID, &role.UnitID, &role.RoleID, &role.CreatedAt, &role.UpdatedAt,
			&role.UnitName, &role.RoleName)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *repository) GetUnitPermissions(unitID int64, roleID int64) ([]*UnitRoleModule, error) {
	query := `
		SELECT urm.id, urm.unit_role_id, urm.module_id, urm.can_read, urm.can_write, 
			urm.can_delete, urm.can_approve, urm.created_at, urm.updated_at
		FROM unit_role_modules urm
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		WHERE ur.unit_id = $1 AND ur.role_id = $2
	`

	rows, err := r.db.Query(query, unitID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*UnitRoleModule
	for rows.Next() {
		perm := &UnitRoleModule{}
		err := rows.Scan(&perm.ID, &perm.UnitRoleID, &perm.ModuleID, &perm.CanRead, &perm.CanWrite,
			&perm.CanDelete, &perm.CanApprove, &perm.CreatedAt, &perm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *repository) UpdatePermissions(unitRoleID int64, modules []UpdateUnitRoleModulePermission) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, module := range modules {
		query := `
			INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (unit_role_id, module_id) 
			DO UPDATE SET can_read = $3, can_write = $4, can_delete = $5, can_approve = $6, updated_at = CURRENT_TIMESTAMP
		`
		_, err := tx.Exec(query, unitRoleID, module.ModuleID, module.CanRead, module.CanWrite, module.CanDelete, module.CanApprove)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) CopyPermissions(sourceUnitID int64, targetUnitID int64, roleID int64, overwrite bool) error {
	query := `
		INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
		SELECT 
			(SELECT id FROM unit_roles WHERE unit_id = $2 AND role_id = $3),
			urm.module_id, urm.can_read, urm.can_write, urm.can_delete, urm.can_approve
		FROM unit_role_modules urm
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		WHERE ur.unit_id = $1 AND ur.role_id = $3
	`

	if overwrite {
		query += ` ON CONFLICT (unit_role_id, module_id) DO UPDATE SET 
			can_read = EXCLUDED.can_read, can_write = EXCLUDED.can_write, 
			can_delete = EXCLUDED.can_delete, can_approve = EXCLUDED.can_approve`
	} else {
		query += ` ON CONFLICT (unit_role_id, module_id) DO NOTHING`
	}

	_, err := r.db.Exec(query, sourceUnitID, targetUnitID, roleID)
	return err
}

func (r *repository) GetUserEffectivePermissions(userID int64) ([]*UnitRoleModule, error) {
	query := `
		SELECT DISTINCT urm.id, urm.unit_role_id, urm.module_id, urm.can_read, urm.can_write, 
			urm.can_delete, urm.can_approve, urm.created_at, urm.updated_at
		FROM unit_role_modules urm
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		JOIN user_roles usr ON usr.unit_id = ur.unit_id AND usr.role_id = ur.role_id
		WHERE usr.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*UnitRoleModule
	for rows.Next() {
		perm := &UnitRoleModule{}
		err := rows.Scan(&perm.ID, &perm.UnitRoleID, &perm.ModuleID, &perm.CanRead, &perm.CanWrite,
			&perm.CanDelete, &perm.CanApprove, &perm.CreatedAt, &perm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}
