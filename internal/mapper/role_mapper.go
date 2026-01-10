package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// RoleMapper handles conversion between role models and DTOs
type RoleMapper struct{}

// NewRoleMapper creates a new role mapper
func NewRoleMapper() *RoleMapper {
	return &RoleMapper{}
}

// ToResponse converts model to response DTO
func (m *RoleMapper) ToResponse(role *models.Role) *dto.RoleResponse {
	if role == nil {
		return nil
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}
}

// ToResponseList converts model slice to response DTO slice
func (m *RoleMapper) ToResponseList(roles []*models.Role) []*dto.RoleResponse {
	if roles == nil {
		return nil
	}

	responses := make([]*dto.RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = m.ToResponse(role)
	}
	return responses
}

// ToWithPermissionsResponse converts model with permissions to response DTO
func (m *RoleMapper) ToWithPermissionsResponse(role *models.RoleWithPermissions) *dto.RoleWithPermissionsResponse {
	if role == nil {
		return nil
	}

	modules := make([]dto.RoleModulePermissionResponse, len(role.Modules))
	for i, module := range role.Modules {
		modules[i] = dto.RoleModulePermissionResponse{
			ModuleID:   module.ModuleID,
			ModuleName: module.ModuleName,
			ModuleURL:  module.ModuleURL,
			CanRead:    module.CanRead,
			CanWrite:   module.CanWrite,
			CanDelete:  module.CanDelete,
		}
	}

	return &dto.RoleWithPermissionsResponse{
		RoleResponse: *m.ToResponse(&role.Role),
		Modules:      modules,
	}
}

// ToModel converts create request DTO to model
func (m *RoleMapper) ToModel(req *dto.CreateRoleRequest) *models.Role {
	if req == nil {
		return nil
	}

	return &models.Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true, // Default to active
	}
}

// UpdateModel updates model with update request DTO
func (m *RoleMapper) UpdateModel(role *models.Role, req *dto.UpdateRoleRequest) {
	if role == nil || req == nil {
		return
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}
}

// ToListResponse creates paginated list response
func (m *RoleMapper) ToListResponse(roles []*models.Role, total int64, limit, offset int) *dto.RoleListResponse {
	return &dto.RoleListResponse{
		Data:    m.ToResponseList(roles),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(roles)) < total,
	}
}

// ToUserRoleModel converts assign role request to user role model
func (m *RoleMapper) ToUserRoleModel(req *dto.AssignRoleRequest) *models.UserRole {
	if req == nil {
		return nil
	}

	return &models.UserRole{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
		BranchID:  req.BranchID,
	}
}
