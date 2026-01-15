package unit

import (
	"errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUnits(req *UnitListRequest) (*UnitListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	units, err := s.repo.GetAll(req.BranchID, limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(req.BranchID, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var responses []*UnitResponse
	for _, unit := range units {
		responses = append(responses, toUnitResponse(unit))
	}

	return &UnitListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) GetUnitByID(id int64) (*UnitResponse, error) {
	unit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, errors.New("unit not found")
	}

	return toUnitResponse(unit), nil
}

func (s *Service) GetUnitHierarchy(branchID int64) ([]*UnitHierarchyResponse, error) {
	units, err := s.repo.GetHierarchy(branchID)
	if err != nil {
		return nil, err
	}

	return buildHierarchy(units), nil
}

func (s *Service) GetUnitWithStats(id int64) (*UnitWithStatsResponse, error) {
	unit, err := s.repo.GetWithStats(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, errors.New("unit not found")
	}

	return &UnitWithStatsResponse{
		UnitResponse: UnitResponse{
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
	}, nil
}

func (s *Service) CreateUnit(req *CreateUnitRequest) (*UnitResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	unit := &Unit{
		BranchID:    req.BranchID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    isActive,
		Level:       1,
		Path:        "/",
	}

	if req.ParentID != nil {
		parent, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("parent unit not found")
		}
		unit.Level = parent.Level + 1
		unit.Path = parent.Path + "/" + parent.Code
	}

	if err := s.repo.Create(unit); err != nil {
		return nil, err
	}

	unitWithBranch, _ := s.repo.GetByID(unit.ID)
	return toUnitResponse(unitWithBranch), nil
}

func (s *Service) UpdateUnit(id int64, req *UpdateUnitRequest) (*UnitResponse, error) {
	unit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, errors.New("unit not found")
	}

	unit.Name = req.Name
	unit.Code = req.Code
	unit.Description = req.Description
	if req.ParentID != nil {
		unit.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		unit.IsActive = *req.IsActive
	}

	if err := s.repo.Update(&unit.Unit); err != nil {
		return nil, err
	}

	unitWithBranch, _ := s.repo.GetByID(id)
	return toUnitResponse(unitWithBranch), nil
}

func (s *Service) DeleteUnit(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) AssignRoleToUnit(unitID int64, roleID int64) error {
	return s.repo.AssignRole(unitID, roleID)
}

func (s *Service) RemoveRoleFromUnit(unitID int64, roleID int64) error {
	return s.repo.RemoveRole(unitID, roleID)
}

func (s *Service) GetUnitRoles(unitID int64) ([]*UnitRoleResponse, error) {
	roles, err := s.repo.GetUnitRoles(unitID)
	if err != nil {
		return nil, err
	}

	var responses []*UnitRoleResponse
	for _, role := range roles {
		responses = append(responses, toUnitRoleResponse(role))
	}

	return responses, nil
}

func (s *Service) GetUnitPermissions(unitID int64, roleID int64) ([]*UnitRoleModuleResponse, error) {
	permissions, err := s.repo.GetUnitPermissions(unitID, roleID)
	if err != nil {
		return nil, err
	}

	var responses []*UnitRoleModuleResponse
	for _, perm := range permissions {
		responses = append(responses, toUnitRoleModuleResponse(perm))
	}

	return responses, nil
}

func (s *Service) UpdateUnitPermissions(unitRoleID int64, req *BulkUpdateUnitRoleModulesRequest) error {
	return s.repo.UpdatePermissions(unitRoleID, req.Modules)
}

func (s *Service) CopyPermissions(req *CopyUnitPermissionsRequest) error {
	return s.repo.CopyPermissions(req.SourceUnitID, req.TargetUnitID, req.RoleID, req.OverwriteExisting)
}

func (s *Service) GetUserEffectivePermissions(userID int64) ([]*UnitRoleModuleResponse, error) {
	permissions, err := s.repo.GetUserEffectivePermissions(userID)
	if err != nil {
		return nil, err
	}

	var responses []*UnitRoleModuleResponse
	for _, perm := range permissions {
		responses = append(responses, toUnitRoleModuleResponse(perm))
	}

	return responses, nil
}

func toUnitResponse(unit *UnitWithBranch) *UnitResponse {
	if unit == nil {
		return nil
	}

	return &UnitResponse{
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

func toUnitRoleResponse(role *UnitRole) *UnitRoleResponse {
	if role == nil {
		return nil
	}

	return &UnitRoleResponse{
		ID:        role.ID,
		UnitID:    role.UnitID,
		RoleID:    role.RoleID,
		CreatedAt: role.CreatedAt.Format(time.RFC3339),
		UpdatedAt: role.UpdatedAt.Format(time.RFC3339),
		UnitName:  role.UnitName,
		RoleName:  role.RoleName,
	}
}

func toUnitRoleModuleResponse(perm *UnitRoleModule) *UnitRoleModuleResponse {
	if perm == nil {
		return nil
	}

	return &UnitRoleModuleResponse{
		ID:         perm.ID,
		UnitRoleID: perm.UnitRoleID,
		ModuleID:   perm.ModuleID,
		CanRead:    perm.CanRead,
		CanWrite:   perm.CanWrite,
		CanDelete:  perm.CanDelete,
		CanApprove: perm.CanApprove,
		CreatedAt:  perm.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  perm.UpdatedAt.Format(time.RFC3339),
	}
}

func buildHierarchy(units []*Unit) []*UnitHierarchyResponse {
	unitMap := make(map[int64]*UnitHierarchyResponse)
	var roots []*UnitHierarchyResponse

	for _, unit := range units {
		unitMap[unit.ID] = &UnitHierarchyResponse{
			UnitResponse: UnitResponse{
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
			Children: []UnitHierarchyResponse{},
		}
	}

	for _, unit := range units {
		if unit.ParentID == nil {
			roots = append(roots, unitMap[unit.ID])
		} else {
			parent := unitMap[*unit.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, *unitMap[unit.ID])
			}
		}
	}

	return roots
}
