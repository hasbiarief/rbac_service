package service

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/errors"
)

type UnitService struct {
	unitRepo           interfaces.UnitRepositoryInterface
	unitRoleRepo       interfaces.UnitRoleRepositoryInterface
	unitRoleModuleRepo interfaces.UnitRoleModuleRepositoryInterface
	branchRepo         interfaces.BranchRepositoryInterface
	roleRepo           interfaces.RoleRepositoryInterface
	moduleRepo         interfaces.ModuleRepositoryInterface
	unitMapper         *mapper.UnitMapper
}

func NewUnitService(
	unitRepo interfaces.UnitRepositoryInterface,
	unitRoleRepo interfaces.UnitRoleRepositoryInterface,
	unitRoleModuleRepo interfaces.UnitRoleModuleRepositoryInterface,
	branchRepo interfaces.BranchRepositoryInterface,
	roleRepo interfaces.RoleRepositoryInterface,
	moduleRepo interfaces.ModuleRepositoryInterface,
	unitMapper *mapper.UnitMapper,
) interfaces.UnitServiceInterface {
	return &UnitService{
		unitRepo:           unitRepo,
		unitRoleRepo:       unitRoleRepo,
		unitRoleModuleRepo: unitRoleModuleRepo,
		branchRepo:         branchRepo,
		roleRepo:           roleRepo,
		moduleRepo:         moduleRepo,
		unitMapper:         unitMapper,
	}
}

// GetUnits retrieves units with filtering and pagination
func (s *UnitService) GetUnits(req *dto.UnitListRequest) (*dto.UnitListResponse, error) {
	// Set default values
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get units from repository
	units, err := s.unitRepo.GetAll(req.BranchID, req.Limit, req.Offset, req.Search, req.IsActive)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get units: " + err.Error())
	}

	// Count total units for pagination
	// Note: We need to implement Count method in repository
	total := int64(len(units)) // Temporary, should be actual count

	// Convert to response DTOs
	unitResponses := make([]*dto.UnitResponse, len(units))
	for i, unit := range units {
		unitResponses[i] = s.unitMapper.ToResponse(unit)
	}

	return &dto.UnitListResponse{
		Data:    unitResponses,
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+req.Limit < int(total),
	}, nil
}

// GetUnitByID retrieves a unit by ID
func (s *UnitService) GetUnitByID(id int64) (*dto.UnitResponse, error) {
	unit, err := s.unitRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit: " + err.Error())
	}

	if unit == nil {
		return nil, errors.NewNotFoundError("Unit")
	}

	return s.unitMapper.ToResponse(unit), nil
}

// GetUnitHierarchy retrieves unit hierarchy for a branch
func (s *UnitService) GetUnitHierarchy(branchID int64) ([]*dto.UnitHierarchyResponse, error) {
	// Verify branch exists
	branch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to verify branch: " + err.Error())
	}
	if branch == nil {
		return nil, errors.NewNotFoundError("Branch")
	}

	// Get hierarchy
	hierarchy, err := s.unitRepo.GetHierarchy(branchID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit hierarchy: " + err.Error())
	}

	return s.unitMapper.ToHierarchyResponse(hierarchy), nil
}

// GetUnitWithStats retrieves unit with statistics
func (s *UnitService) GetUnitWithStats(id int64) (*dto.UnitWithStatsResponse, error) {
	unit, err := s.unitRepo.GetWithStats(id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit with stats: " + err.Error())
	}

	if unit == nil {
		return nil, errors.NewNotFoundError("Unit")
	}

	return s.unitMapper.ToStatsResponse(unit), nil
}

// CreateUnit creates a new unit
func (s *UnitService) CreateUnit(req *dto.CreateUnitRequest) (*dto.UnitResponse, error) {
	// Verify branch exists
	branch, err := s.branchRepo.GetByID(req.BranchID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to verify branch: " + err.Error())
	}
	if branch == nil {
		return nil, errors.NewNotFoundError("Branch")
	}

	// Verify parent unit exists if specified
	if req.ParentID != nil {
		parent, err := s.unitRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to verify parent unit: " + err.Error())
		}
		if parent == nil {
			return nil, errors.NewNotFoundError("Parent unit")
		}
		// Ensure parent is in the same branch
		if parent.BranchID != req.BranchID {
			return nil, errors.NewBadRequestError("Parent unit must be in the same branch")
		}
	}

	// Check if code already exists in the branch
	exists, err := s.unitRepo.ExistsByCode(req.BranchID, req.Code, nil)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to check unit code: " + err.Error())
	}
	if exists {
		return nil, errors.NewBadRequestError("Unit code already exists in this branch")
	}

	// Convert to model
	unit := s.unitMapper.ToModel(req)
	if req.IsActive == nil {
		unit.IsActive = true
	} else {
		unit.IsActive = *req.IsActive
	}

	// Create unit
	if err := s.unitRepo.Create(unit); err != nil {
		return nil, errors.NewInternalServerError("Failed to create unit: " + err.Error())
	}

	// Get created unit with full details
	createdUnit, err := s.unitRepo.GetByID(unit.ID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get created unit: " + err.Error())
	}

	return s.unitMapper.ToResponse(createdUnit), nil
}

// UpdateUnit updates an existing unit
func (s *UnitService) UpdateUnit(id int64, req *dto.UpdateUnitRequest) (*dto.UnitResponse, error) {
	// Get existing unit
	existingUnit, err := s.unitRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit: " + err.Error())
	}
	if existingUnit == nil {
		return nil, errors.NewNotFoundError("Unit")
	}

	// Verify parent unit exists if specified and changed
	if req.ParentID != nil && (existingUnit.ParentID == nil || *req.ParentID != *existingUnit.ParentID) {
		parent, err := s.unitRepo.GetByID(*req.ParentID)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to verify parent unit: " + err.Error())
		}
		if parent == nil {
			return nil, errors.NewNotFoundError("Parent unit")
		}
		// Ensure parent is in the same branch
		if parent.BranchID != existingUnit.BranchID {
			return nil, errors.NewBadRequestError("Parent unit must be in the same branch")
		}
		// Prevent circular reference
		if *req.ParentID == id {
			return nil, errors.NewBadRequestError("Unit cannot be its own parent")
		}
	}

	// Check if code already exists in the branch (excluding current unit)
	if req.Code != existingUnit.Code {
		exists, err := s.unitRepo.ExistsByCode(existingUnit.BranchID, req.Code, &id)
		if err != nil {
			return nil, errors.NewInternalServerError("Failed to check unit code: " + err.Error())
		}
		if exists {
			return nil, errors.NewBadRequestError("Unit code already exists in this branch")
		}
	}

	// Update unit fields
	unit := &models.Unit{
		ID:          id,
		BranchID:    existingUnit.BranchID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    existingUnit.IsActive,
	}

	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
	}

	// Update unit
	if err := s.unitRepo.Update(unit); err != nil {
		return nil, errors.NewInternalServerError("Failed to update unit: " + err.Error())
	}

	// Get updated unit with full details
	updatedUnit, err := s.unitRepo.GetByID(id)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get updated unit: " + err.Error())
	}

	return s.unitMapper.ToResponse(updatedUnit), nil
}

// DeleteUnit soft deletes a unit
func (s *UnitService) DeleteUnit(id int64) error {
	// Check if unit exists
	unit, err := s.unitRepo.GetByID(id)
	if err != nil {
		return errors.NewInternalServerError("Failed to get unit: " + err.Error())
	}
	if unit == nil {
		return errors.NewNotFoundError("Unit")
	}

	// Check if unit has children
	hasChildren, err := s.unitRepo.HasChildren(id)
	if err != nil {
		return errors.NewInternalServerError("Failed to check unit children: " + err.Error())
	}
	if hasChildren {
		return errors.NewBadRequestError("Cannot delete unit with children")
	}

	// Check if unit has assigned users
	hasUsers, err := s.unitRepo.HasUsers(id)
	if err != nil {
		return errors.NewInternalServerError("Failed to check unit users: " + err.Error())
	}
	if hasUsers {
		return errors.NewBadRequestError("Cannot delete unit with assigned users")
	}

	// Soft delete unit
	if err := s.unitRepo.Delete(id); err != nil {
		return errors.NewInternalServerError("Failed to delete unit: " + err.Error())
	}

	return nil
}

// AssignRoleToUnit assigns a role to a unit
func (s *UnitService) AssignRoleToUnit(unitID, roleID int64) error {
	// Verify unit exists
	unit, err := s.unitRepo.GetByID(unitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify unit: " + err.Error())
	}
	if unit == nil {
		return errors.NewNotFoundError("Unit")
	}

	// Verify role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify role: " + err.Error())
	}
	if role == nil {
		return errors.NewNotFoundError("Role")
	}

	// Check if assignment already exists
	exists, err := s.unitRoleRepo.ExistsByUnitAndRole(unitID, roleID)
	if err != nil {
		return errors.NewInternalServerError("Failed to check existing assignment: " + err.Error())
	}
	if exists {
		return errors.NewBadRequestError("Role already assigned to unit")
	}

	// Create assignment
	unitRole := &models.UnitRole{
		UnitID: unitID,
		RoleID: roleID,
	}

	if err := s.unitRoleRepo.Create(unitRole); err != nil {
		return errors.NewInternalServerError("Failed to assign role to unit: " + err.Error())
	}

	return nil
}

// RemoveRoleFromUnit removes a role from a unit
func (s *UnitService) RemoveRoleFromUnit(unitID, roleID int64) error {
	// Find the unit role assignment
	unitRoles, err := s.unitRoleRepo.GetRolesByUnit(unitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get unit roles: " + err.Error())
	}

	var unitRoleID int64
	found := false
	for _, ur := range unitRoles {
		if ur.RoleID == roleID {
			unitRoleID = ur.ID
			found = true
			break
		}
	}

	if !found {
		return errors.NewNotFoundError("Unit role assignment")
	}

	// Delete the assignment
	if err := s.unitRoleRepo.Delete(unitRoleID); err != nil {
		return errors.NewInternalServerError("Failed to remove role from unit: " + err.Error())
	}

	return nil
}

// GetUnitRoles retrieves all roles for a unit
func (s *UnitService) GetUnitRoles(unitID int64) ([]*dto.UnitRoleResponse, error) {
	// Verify unit exists
	unit, err := s.unitRepo.GetByID(unitID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to verify unit: " + err.Error())
	}
	if unit == nil {
		return nil, errors.NewNotFoundError("Unit")
	}

	// Get unit roles
	unitRoles, err := s.unitRoleRepo.GetRolesByUnit(unitID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit roles: " + err.Error())
	}

	// Convert to response DTOs
	responses := make([]*dto.UnitRoleResponse, len(unitRoles))
	for i, ur := range unitRoles {
		responses[i] = s.unitMapper.ToUnitRoleResponse(ur)
	}

	return responses, nil
}

// GetUnitPermissions retrieves permissions for a unit-role combination
func (s *UnitService) GetUnitPermissions(unitID, roleID int64) ([]*dto.UnitRoleModuleResponse, error) {
	// Find unit role
	unitRoles, err := s.unitRoleRepo.GetRolesByUnit(unitID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit roles: " + err.Error())
	}

	var unitRoleID int64
	found := false
	for _, ur := range unitRoles {
		if ur.RoleID == roleID {
			unitRoleID = ur.ID
			found = true
			break
		}
	}

	if !found {
		return nil, errors.NewNotFoundError("Unit role assignment")
	}

	// Get permissions
	permissions, err := s.unitRoleModuleRepo.GetByUnitRole(unitRoleID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get unit permissions: " + err.Error())
	}

	// Convert to response DTOs
	responses := make([]*dto.UnitRoleModuleResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = s.unitMapper.ToUnitRoleModuleResponse(perm)
	}

	return responses, nil
}

// UpdateUnitPermissions updates permissions for a unit role
func (s *UnitService) UpdateUnitPermissions(unitRoleID int64, req *dto.BulkUpdateUnitRoleModulesRequest) error {
	// Verify unit role exists
	unitRole, err := s.unitRoleRepo.GetByID(unitRoleID)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify unit role: " + err.Error())
	}
	if unitRole == nil {
		return errors.NewNotFoundError("Unit role")
	}

	// Convert request to models
	permissions := make([]*models.UnitRoleModule, len(req.Modules))
	for i, module := range req.Modules {
		permissions[i] = &models.UnitRoleModule{
			UnitRoleID: unitRoleID,
			ModuleID:   module.ModuleID,
			CanRead:    module.CanRead,
			CanWrite:   module.CanWrite,
			CanDelete:  module.CanDelete,
			CanApprove: module.CanApprove,
		}
	}

	// Delete existing permissions
	if err := s.unitRoleModuleRepo.DeleteByUnitRole(unitRoleID); err != nil {
		return errors.NewInternalServerError("Failed to delete existing permissions: " + err.Error())
	}

	// Create new permissions
	if err := s.unitRoleModuleRepo.BulkCreate(permissions); err != nil {
		return errors.NewInternalServerError("Failed to create permissions: " + err.Error())
	}

	return nil
}

// CopyPermissions copies permissions from one unit to another
func (s *UnitService) CopyPermissions(req *dto.CopyUnitPermissionsRequest) error {
	// Verify source and target units exist
	sourceUnit, err := s.unitRepo.GetByID(req.SourceUnitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify source unit: " + err.Error())
	}
	if sourceUnit == nil {
		return errors.NewNotFoundError("Source unit")
	}

	targetUnit, err := s.unitRepo.GetByID(req.TargetUnitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify target unit: " + err.Error())
	}
	if targetUnit == nil {
		return errors.NewNotFoundError("Target unit")
	}

	// Find source unit role
	sourceUnitRoles, err := s.unitRoleRepo.GetRolesByUnit(req.SourceUnitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get source unit roles: " + err.Error())
	}

	var sourceUnitRoleID int64
	found := false
	for _, ur := range sourceUnitRoles {
		if ur.RoleID == req.RoleID {
			sourceUnitRoleID = ur.ID
			found = true
			break
		}
	}

	if !found {
		return errors.NewNotFoundError("Source unit role assignment")
	}

	// Find or create target unit role
	targetUnitRoles, err := s.unitRoleRepo.GetRolesByUnit(req.TargetUnitID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get target unit roles: " + err.Error())
	}

	var targetUnitRoleID int64
	found = false
	for _, ur := range targetUnitRoles {
		if ur.RoleID == req.RoleID {
			targetUnitRoleID = ur.ID
			found = true
			break
		}
	}

	if !found {
		// Create target unit role
		targetUnitRole := &models.UnitRole{
			UnitID: req.TargetUnitID,
			RoleID: req.RoleID,
		}
		if err := s.unitRoleRepo.Create(targetUnitRole); err != nil {
			return errors.NewInternalServerError("Failed to create target unit role: " + err.Error())
		}
		targetUnitRoleID = targetUnitRole.ID
	}

	// Get source permissions
	sourcePermissions, err := s.unitRoleModuleRepo.GetByUnitRole(sourceUnitRoleID)
	if err != nil {
		return errors.NewInternalServerError("Failed to get source permissions: " + err.Error())
	}

	// Delete existing target permissions if overwrite is enabled
	if req.OverwriteExisting {
		if err := s.unitRoleModuleRepo.DeleteByUnitRole(targetUnitRoleID); err != nil {
			return errors.NewInternalServerError("Failed to delete existing target permissions: " + err.Error())
		}
	}

	// Copy permissions
	targetPermissions := make([]*models.UnitRoleModule, len(sourcePermissions))
	for i, perm := range sourcePermissions {
		targetPermissions[i] = &models.UnitRoleModule{
			UnitRoleID: targetUnitRoleID,
			ModuleID:   perm.ModuleID,
			CanRead:    perm.CanRead,
			CanWrite:   perm.CanWrite,
			CanDelete:  perm.CanDelete,
			CanApprove: perm.CanApprove,
		}
	}

	if err := s.unitRoleModuleRepo.BulkCreate(targetPermissions); err != nil {
		return errors.NewInternalServerError("Failed to copy permissions: " + err.Error())
	}

	return nil
}

// GetUserEffectivePermissions retrieves effective permissions for a user
func (s *UnitService) GetUserEffectivePermissions(userID int64) ([]*dto.UnitRoleModuleResponse, error) {
	permissions, err := s.unitRoleModuleRepo.GetEffectivePermissions(userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to get user permissions: " + err.Error())
	}

	// Convert to response DTOs
	responses := make([]*dto.UnitRoleModuleResponse, len(permissions))
	for i, perm := range permissions {
		responses[i] = &dto.UnitRoleModuleResponse{
			ModuleID:       perm.ModuleID,
			ModuleName:     perm.ModuleName,
			ModuleCategory: perm.ModuleCategory,
			CanRead:        perm.CanRead,
			CanWrite:       perm.CanWrite,
			CanDelete:      perm.CanDelete,
			CanApprove:     perm.CanApprove,
			IsCustomized:   perm.IsCustomized,
		}
	}

	return responses, nil
}
