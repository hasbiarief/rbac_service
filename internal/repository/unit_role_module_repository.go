package repository

import (
	"database/sql"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/query"
)

type UnitRoleModuleRepository struct {
	*model.Repository
	db *sql.DB
}

func NewUnitRoleModuleRepository(db *sql.DB) interfaces.UnitRoleModuleRepositoryInterface {
	return &UnitRoleModuleRepository{
		Repository: &model.Repository{},
		db:         db,
	}
}

// GetAll retrieves unit role modules with optional filtering
func (r *UnitRoleModuleRepository) GetAll(unitRoleID, moduleID *int64, limit, offset int) ([]*models.UnitRoleModule, error) {
	qb := query.NewQueryBuilder(`
		SELECT 
			urm.id, urm.unit_role_id, urm.module_id, urm.can_read, urm.can_write, 
			urm.can_delete, urm.can_approve, urm.created_at, urm.updated_at,
			m.name as module_name, m.category as module_category,
			u.name as unit_name, r.name as role_name
		FROM unit_role_modules urm
		JOIN modules m ON urm.module_id = m.id
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE 1=1
	`)

	if unitRoleID != nil {
		qb.AddCondition("AND urm.unit_role_id = $%d", *unitRoleID)
	}

	if moduleID != nil {
		qb.AddCondition("AND urm.module_id = $%d", *moduleID)
	}

	qb.AddCondition("ORDER BY m.category, m.name LIMIT $%d", limit)
	qb.AddCondition("OFFSET $%d", offset)

	query, args := qb.Build()
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitRoleModules []*models.UnitRoleModule
	for rows.Next() {
		urm := &models.UnitRoleModule{}
		err := rows.Scan(
			&urm.ID, &urm.UnitRoleID, &urm.ModuleID, &urm.CanRead, &urm.CanWrite,
			&urm.CanDelete, &urm.CanApprove, &urm.CreatedAt, &urm.UpdatedAt,
			&urm.ModuleName, &urm.ModuleCategory, &urm.UnitName, &urm.RoleName,
		)
		if err != nil {
			return nil, err
		}
		unitRoleModules = append(unitRoleModules, urm)
	}

	return unitRoleModules, nil
}

// GetByID retrieves a unit role module by ID
func (r *UnitRoleModuleRepository) GetByID(id int64) (*models.UnitRoleModule, error) {
	query := `
		SELECT 
			urm.id, urm.unit_role_id, urm.module_id, urm.can_read, urm.can_write, 
			urm.can_delete, urm.can_approve, urm.created_at, urm.updated_at,
			m.name as module_name, m.category as module_category,
			u.name as unit_name, r.name as role_name
		FROM unit_role_modules urm
		JOIN modules m ON urm.module_id = m.id
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE urm.id = $1
	`

	urm := &models.UnitRoleModule{}
	err := r.db.QueryRow(query, id).Scan(
		&urm.ID, &urm.UnitRoleID, &urm.ModuleID, &urm.CanRead, &urm.CanWrite,
		&urm.CanDelete, &urm.CanApprove, &urm.CreatedAt, &urm.UpdatedAt,
		&urm.ModuleName, &urm.ModuleCategory, &urm.UnitName, &urm.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return urm, nil
}

// Create creates a new unit role module
func (r *UnitRoleModuleRepository) Create(unitRoleModule *models.UnitRoleModule) error {
	query := `
		INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(query,
		unitRoleModule.UnitRoleID, unitRoleModule.ModuleID, unitRoleModule.CanRead,
		unitRoleModule.CanWrite, unitRoleModule.CanDelete, unitRoleModule.CanApprove,
	).Scan(&unitRoleModule.ID, &unitRoleModule.CreatedAt, &unitRoleModule.UpdatedAt)
}

// Update updates an existing unit role module
func (r *UnitRoleModuleRepository) Update(unitRoleModule *models.UnitRoleModule) error {
	query := `
		UPDATE unit_role_modules 
		SET can_read = $2, can_write = $3, can_delete = $4, can_approve = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	return r.db.QueryRow(query,
		unitRoleModule.ID, unitRoleModule.CanRead, unitRoleModule.CanWrite,
		unitRoleModule.CanDelete, unitRoleModule.CanApprove,
	).Scan(&unitRoleModule.UpdatedAt)
}

// Delete deletes a unit role module
func (r *UnitRoleModuleRepository) Delete(id int64) error {
	query := `DELETE FROM unit_role_modules WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// BulkCreate creates multiple unit role modules
func (r *UnitRoleModuleRepository) BulkCreate(unitRoleModules []*models.UnitRoleModule) error {
	if len(unitRoleModules) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, urm := range unitRoleModules {
		err = stmt.QueryRow(
			urm.UnitRoleID, urm.ModuleID, urm.CanRead,
			urm.CanWrite, urm.CanDelete, urm.CanApprove,
		).Scan(&urm.ID, &urm.CreatedAt, &urm.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// BulkUpdate updates multiple unit role modules
func (r *UnitRoleModuleRepository) BulkUpdate(unitRoleModules []*models.UnitRoleModule) error {
	if len(unitRoleModules) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		UPDATE unit_role_modules 
		SET can_read = $2, can_write = $3, can_delete = $4, can_approve = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, urm := range unitRoleModules {
		err = stmt.QueryRow(
			urm.ID, urm.CanRead, urm.CanWrite, urm.CanDelete, urm.CanApprove,
		).Scan(&urm.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteByUnitRole deletes all permissions for a unit role
func (r *UnitRoleModuleRepository) DeleteByUnitRole(unitRoleID int64) error {
	query := `DELETE FROM unit_role_modules WHERE unit_role_id = $1`
	_, err := r.db.Exec(query, unitRoleID)
	return err
}

// GetByUnitRole retrieves all permissions for a unit role
func (r *UnitRoleModuleRepository) GetByUnitRole(unitRoleID int64) ([]*models.UnitRoleModule, error) {
	query := `
		SELECT 
			urm.id, urm.unit_role_id, urm.module_id, urm.can_read, urm.can_write, 
			urm.can_delete, urm.can_approve, urm.created_at, urm.updated_at,
			m.name as module_name, m.category as module_category,
			u.name as unit_name, r.name as role_name
		FROM unit_role_modules urm
		JOIN modules m ON urm.module_id = m.id
		JOIN unit_roles ur ON urm.unit_role_id = ur.id
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE urm.unit_role_id = $1
		ORDER BY m.category, m.name
	`

	rows, err := r.db.Query(query, unitRoleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitRoleModules []*models.UnitRoleModule
	for rows.Next() {
		urm := &models.UnitRoleModule{}
		err := rows.Scan(
			&urm.ID, &urm.UnitRoleID, &urm.ModuleID, &urm.CanRead, &urm.CanWrite,
			&urm.CanDelete, &urm.CanApprove, &urm.CreatedAt, &urm.UpdatedAt,
			&urm.ModuleName, &urm.ModuleCategory, &urm.UnitName, &urm.RoleName,
		)
		if err != nil {
			return nil, err
		}
		unitRoleModules = append(unitRoleModules, urm)
	}

	return unitRoleModules, nil
}

// GetEffectivePermissions retrieves effective permissions for a user
func (r *UnitRoleModuleRepository) GetEffectivePermissions(userID int64) ([]*models.UnitRoleModulePermission, error) {
	query := `
		SELECT DISTINCT
			m.id as module_id,
			m.name as module_name,
			m.category as module_category,
			m.url as module_url,
			COALESCE(urm.can_read, rm.can_read, false) as can_read,
			COALESCE(urm.can_write, rm.can_write, false) as can_write,
			COALESCE(urm.can_delete, rm.can_delete, false) as can_delete,
			COALESCE(urm.can_approve, rm.can_approve, false) as can_approve,
			CASE WHEN urm.id IS NOT NULL THEN true ELSE false END as is_customized
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		JOIN role_modules rm ON r.id = rm.role_id
		JOIN modules m ON rm.module_id = m.id
		LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
		LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id AND m.id = urm.module_id
		WHERE ur.user_id = $1
		AND r.is_active = true
		AND m.is_active = true
		ORDER BY m.category, m.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.UnitRoleModulePermission
	for rows.Next() {
		perm := &models.UnitRoleModulePermission{}
		err := rows.Scan(
			&perm.ModuleID, &perm.ModuleName, &perm.ModuleCategory, &perm.ModuleURL,
			&perm.CanRead, &perm.CanWrite, &perm.CanDelete, &perm.CanApprove,
			&perm.IsCustomized,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// ExistsByUnitRoleAndModule checks if a unit role module combination exists
func (r *UnitRoleModuleRepository) ExistsByUnitRoleAndModule(unitRoleID, moduleID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM unit_role_modules WHERE unit_role_id = $1 AND module_id = $2`

	var count int64
	err := r.db.QueryRow(query, unitRoleID, moduleID).Scan(&count)
	return count > 0, err
}
