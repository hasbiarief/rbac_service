package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

type UnitMapper struct{}

func NewUnitMapper() *UnitMapper {
	return &UnitMapper{}
}

// ToResponse converts UnitWithBranch model to UnitResponse DTO
func (m *UnitMapper) ToResponse(unit *models.UnitWithBranch) *dto.UnitResponse {
	if unit == nil {
		return nil
	}

	return &dto.UnitResponse{
		ID:          unit.ID,
		BranchID:    unit.BranchID,
		ParentID:    unit.ParentID,
		Name:        unit.Name,
		Code:        unit.Code,
		Description: unit.Description,
		Level:       unit.Level,
		Path:        unit.Path,
		IsActive:    unit.IsActive,
		CreatedAt:   unit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   unit.UpdatedAt.Format(time.RFC3339),
		BranchName:  unit.BranchName,
		BranchCode:  unit.BranchCode,
		CompanyName: unit.CompanyName,
		CompanyCode: unit.CompanyCode,
	}
}

// ToModel converts CreateUnitRequest DTO to Unit model
func (m *UnitMapper) ToModel(req *dto.CreateUnitRequest) *models.Unit {
	if req == nil {
		return nil
	}

	return &models.Unit{
		BranchID:    req.BranchID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}
}

// ToHierarchyResponse converts UnitHierarchy slice to UnitHierarchyResponse slice
func (m *UnitMapper) ToHierarchyResponse(hierarchy []*models.UnitHierarchy) []*dto.UnitHierarchyResponse {
	if hierarchy == nil {
		return nil
	}

	responses := make([]*dto.UnitHierarchyResponse, len(hierarchy))
	for i, unit := range hierarchy {
		responses[i] = m.toHierarchyResponseSingle(unit)
	}

	return responses
}

// toHierarchyResponseSingle converts single UnitHierarchy to UnitHierarchyResponse
func (m *UnitMapper) toHierarchyResponseSingle(unit *models.UnitHierarchy) *dto.UnitHierarchyResponse {
	if unit == nil {
		return nil
	}

	response := &dto.UnitHierarchyResponse{
		UnitResponse: dto.UnitResponse{
			ID:          unit.ID,
			BranchID:    unit.BranchID,
			ParentID:    unit.ParentID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
			Level:       unit.Level,
			Path:        unit.Path,
			IsActive:    unit.IsActive,
			CreatedAt:   unit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   unit.UpdatedAt.Format(time.RFC3339),
		},
	}

	// Convert children recursively
	if len(unit.Children) > 0 {
		response.Children = make([]dto.UnitHierarchyResponse, len(unit.Children))
		for i, child := range unit.Children {
			childResponse := m.toHierarchyResponseSingle(&child)
			if childResponse != nil {
				response.Children[i] = *childResponse
			}
		}
	}

	return response
}

// ToStatsResponse converts UnitWithStats model to UnitWithStatsResponse DTO
func (m *UnitMapper) ToStatsResponse(unit *models.UnitWithStats) *dto.UnitWithStatsResponse {
	if unit == nil {
		return nil
	}

	return &dto.UnitWithStatsResponse{
		UnitResponse: dto.UnitResponse{
			ID:          unit.ID,
			BranchID:    unit.BranchID,
			ParentID:    unit.ParentID,
			Name:        unit.Name,
			Code:        unit.Code,
			Description: unit.Description,
			Level:       unit.Level,
			Path:        unit.Path,
			IsActive:    unit.IsActive,
			CreatedAt:   unit.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   unit.UpdatedAt.Format(time.RFC3339),
		},
		TotalUsers:    unit.TotalUsers,
		TotalSubUnits: unit.TotalSubUnits,
		TotalRoles:    unit.TotalRoles,
	}
}

// ToUnitRoleResponse converts UnitRole model to UnitRoleResponse DTO
func (m *UnitMapper) ToUnitRoleResponse(unitRole *models.UnitRole) *dto.UnitRoleResponse {
	if unitRole == nil {
		return nil
	}

	return &dto.UnitRoleResponse{
		ID:        unitRole.ID,
		UnitID:    unitRole.UnitID,
		RoleID:    unitRole.RoleID,
		CreatedAt: unitRole.CreatedAt.Format(time.RFC3339),
		UpdatedAt: unitRole.UpdatedAt.Format(time.RFC3339),
		UnitName:  unitRole.UnitName,
		RoleName:  unitRole.RoleName,
	}
}

// ToUnitRoleModuleResponse converts UnitRoleModule model to UnitRoleModuleResponse DTO
func (m *UnitMapper) ToUnitRoleModuleResponse(urm *models.UnitRoleModule) *dto.UnitRoleModuleResponse {
	if urm == nil {
		return nil
	}

	return &dto.UnitRoleModuleResponse{
		ID:             urm.ID,
		UnitRoleID:     urm.UnitRoleID,
		ModuleID:       urm.ModuleID,
		CanRead:        urm.CanRead,
		CanWrite:       urm.CanWrite,
		CanDelete:      urm.CanDelete,
		CanApprove:     urm.CanApprove,
		CreatedAt:      urm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      urm.UpdatedAt.Format(time.RFC3339),
		ModuleName:     urm.ModuleName,
		ModuleCategory: urm.ModuleCategory,
		UnitName:       urm.UnitName,
		RoleName:       urm.RoleName,
		IsCustomized:   true, // Unit-specific permissions are always customized
	}
}

// ToListResponse converts slice of UnitWithBranch to UnitListResponse
func (m *UnitMapper) ToListResponse(units []*models.UnitWithBranch, total int64, limit, offset int) *dto.UnitListResponse {
	responses := make([]*dto.UnitResponse, len(units))
	for i, unit := range units {
		responses[i] = m.ToResponse(unit)
	}

	return &dto.UnitListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+limit < int(total),
	}
}

// ToUnitRoleListResponse converts slice of UnitRole to UnitRoleListResponse
func (m *UnitMapper) ToUnitRoleListResponse(unitRoles []*models.UnitRole, total int64, limit, offset int) *dto.UnitRoleListResponse {
	responses := make([]*dto.UnitRoleResponse, len(unitRoles))
	for i, ur := range unitRoles {
		responses[i] = m.ToUnitRoleResponse(ur)
	}

	return &dto.UnitRoleListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+limit < int(total),
	}
}

// ToUnitRoleWithPermissionsResponse converts UnitRoleWithPermissions to response DTO
func (m *UnitMapper) ToUnitRoleWithPermissionsResponse(urwp *models.UnitRoleWithPermissions) *dto.UnitRoleWithPermissionsResponse {
	if urwp == nil {
		return nil
	}

	modules := make([]dto.UnitRoleModuleResponse, len(urwp.Modules))
	for i, module := range urwp.Modules {
		modules[i] = dto.UnitRoleModuleResponse{
			ModuleID:       module.ModuleID,
			ModuleName:     module.ModuleName,
			ModuleCategory: module.ModuleCategory,
			ModuleURL:      module.ModuleURL,
			CanRead:        module.CanRead,
			CanWrite:       module.CanWrite,
			CanDelete:      module.CanDelete,
			CanApprove:     module.CanApprove,
			IsCustomized:   module.IsCustomized,
		}
	}

	return &dto.UnitRoleWithPermissionsResponse{
		UnitRoleResponse: *m.ToUnitRoleResponse(&urwp.UnitRole),
		Modules:          modules,
	}
}

// ToPermissionSummaryResponse converts unit permission data to summary response
func (m *UnitMapper) ToPermissionSummaryResponse(
	unit *models.UnitWithBranch,
	roles []dto.UnitRolePermissionSummary,
) *dto.UnitPermissionSummaryResponse {
	if unit == nil {
		return nil
	}

	return &dto.UnitPermissionSummaryResponse{
		UnitID:      unit.ID,
		UnitName:    unit.Name,
		UnitCode:    unit.Code,
		BranchName:  unit.BranchName,
		CompanyName: unit.CompanyName,
		Roles:       roles,
	}
}
