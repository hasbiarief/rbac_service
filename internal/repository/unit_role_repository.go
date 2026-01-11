package repository

import (
	"database/sql"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/query"
)

type UnitRoleRepository struct {
	*model.Repository
	db *sql.DB
}

func NewUnitRoleRepository(db *sql.DB) interfaces.UnitRoleRepositoryInterface {
	return &UnitRoleRepository{
		Repository: &model.Repository{},
		db:         db,
	}
}

// GetAll retrieves unit roles with optional filtering
func (r *UnitRoleRepository) GetAll(unitID, roleID *int64, limit, offset int) ([]*models.UnitRole, error) {
	qb := query.NewQueryBuilder(`
		SELECT 
			ur.id, ur.unit_id, ur.role_id, ur.created_at, ur.updated_at,
			u.name as unit_name, r.name as role_name
		FROM unit_roles ur
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE 1=1
	`)

	if unitID != nil {
		qb.AddCondition("AND ur.unit_id = $%d", *unitID)
	}

	if roleID != nil {
		qb.AddCondition("AND ur.role_id = $%d", *roleID)
	}

	qb.AddCondition("ORDER BY u.name, r.name LIMIT $%d", limit)
	qb.AddCondition("OFFSET $%d", offset)

	query, args := qb.Build()
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitRoles []*models.UnitRole
	for rows.Next() {
		unitRole := &models.UnitRole{}
		err := rows.Scan(
			&unitRole.ID, &unitRole.UnitID, &unitRole.RoleID, &unitRole.CreatedAt, &unitRole.UpdatedAt,
			&unitRole.UnitName, &unitRole.RoleName,
		)
		if err != nil {
			return nil, err
		}
		unitRoles = append(unitRoles, unitRole)
	}

	return unitRoles, nil
}

// GetByID retrieves a unit role by ID
func (r *UnitRoleRepository) GetByID(id int64) (*models.UnitRole, error) {
	query := `
		SELECT 
			ur.id, ur.unit_id, ur.role_id, ur.created_at, ur.updated_at,
			u.name as unit_name, r.name as role_name
		FROM unit_roles ur
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.id = $1
	`

	unitRole := &models.UnitRole{}
	err := r.db.QueryRow(query, id).Scan(
		&unitRole.ID, &unitRole.UnitID, &unitRole.RoleID, &unitRole.CreatedAt, &unitRole.UpdatedAt,
		&unitRole.UnitName, &unitRole.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return unitRole, nil
}

// Create creates a new unit role
func (r *UnitRoleRepository) Create(unitRole *models.UnitRole) error {
	query := `
		INSERT INTO unit_roles (unit_id, role_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(query, unitRole.UnitID, unitRole.RoleID).Scan(
		&unitRole.ID, &unitRole.CreatedAt, &unitRole.UpdatedAt,
	)
}

// Delete deletes a unit role
func (r *UnitRoleRepository) Delete(id int64) error {
	query := `DELETE FROM unit_roles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetRolesByUnit retrieves all roles for a unit
func (r *UnitRoleRepository) GetRolesByUnit(unitID int64) ([]*models.UnitRole, error) {
	query := `
		SELECT 
			ur.id, ur.unit_id, ur.role_id, ur.created_at, ur.updated_at,
			u.name as unit_name, r.name as role_name
		FROM unit_roles ur
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.unit_id = $1
		ORDER BY r.name
	`

	rows, err := r.db.Query(query, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitRoles []*models.UnitRole
	for rows.Next() {
		unitRole := &models.UnitRole{}
		err := rows.Scan(
			&unitRole.ID, &unitRole.UnitID, &unitRole.RoleID, &unitRole.CreatedAt, &unitRole.UpdatedAt,
			&unitRole.UnitName, &unitRole.RoleName,
		)
		if err != nil {
			return nil, err
		}
		unitRoles = append(unitRoles, unitRole)
	}

	return unitRoles, nil
}

// GetUnitsByRole retrieves all units for a role
func (r *UnitRoleRepository) GetUnitsByRole(roleID int64) ([]*models.UnitRole, error) {
	query := `
		SELECT 
			ur.id, ur.unit_id, ur.role_id, ur.created_at, ur.updated_at,
			u.name as unit_name, r.name as role_name
		FROM unit_roles ur
		JOIN units u ON ur.unit_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.role_id = $1
		ORDER BY u.name
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unitRoles []*models.UnitRole
	for rows.Next() {
		unitRole := &models.UnitRole{}
		err := rows.Scan(
			&unitRole.ID, &unitRole.UnitID, &unitRole.RoleID, &unitRole.CreatedAt, &unitRole.UpdatedAt,
			&unitRole.UnitName, &unitRole.RoleName,
		)
		if err != nil {
			return nil, err
		}
		unitRoles = append(unitRoles, unitRole)
	}

	return unitRoles, nil
}

// ExistsByUnitAndRole checks if a unit-role combination exists
func (r *UnitRoleRepository) ExistsByUnitAndRole(unitID, roleID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM unit_roles WHERE unit_id = $1 AND role_id = $2`

	var count int64
	err := r.db.QueryRow(query, unitID, roleID).Scan(&count)
	return count > 0, err
}

// GetWithPermissions retrieves unit role with its permissions
func (r *UnitRoleRepository) GetWithPermissions(unitRoleID int64) (*models.UnitRoleWithPermissions, error) {
	// First get the unit role
	unitRole, err := r.GetByID(unitRoleID)
	if err != nil || unitRole == nil {
		return nil, err
	}

	// Then get the permissions
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

	var permissions []models.UnitRoleModulePermission
	for rows.Next() {
		var perm models.UnitRoleModule
		err := rows.Scan(
			&perm.ID, &perm.UnitRoleID, &perm.ModuleID, &perm.CanRead, &perm.CanWrite,
			&perm.CanDelete, &perm.CanApprove, &perm.CreatedAt, &perm.UpdatedAt,
			&perm.ModuleName, &perm.ModuleCategory, &perm.UnitName, &perm.RoleName,
		)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, models.UnitRoleModulePermission{
			ModuleID:       perm.ModuleID,
			ModuleName:     perm.ModuleName,
			ModuleCategory: perm.ModuleCategory,
			ModuleURL:      "", // Will be filled if needed
			CanRead:        perm.CanRead,
			CanWrite:       perm.CanWrite,
			CanDelete:      perm.CanDelete,
			CanApprove:     perm.CanApprove,
			IsCustomized:   true, // Since these are unit-specific permissions
		})
	}

	return &models.UnitRoleWithPermissions{
		UnitRole: *unitRole,
		Modules:  permissions,
	}, nil
}
