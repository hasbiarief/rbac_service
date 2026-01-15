package role

import (
	"fmt"
	"time"
)

type Service struct {
	roleRepo *RoleRepository
}

func NewService(roleRepo *RoleRepository) *Service {
	return &Service{
		roleRepo: roleRepo,
	}
}

func (s *Service) GetRoles(req *RoleListRequest) (*RoleListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	roles, err := s.roleRepo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.roleRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var roleResponses []*RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, toRoleResponse(role))
	}

	return &RoleListResponse{
		Data:    roleResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(roleResponses)) < total,
	}, nil
}

func (s *Service) GetRoleByID(id int64) (*RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

func (s *Service) GetRoleWithPermissions(id int64) (*RoleWithPermissionsResponse, error) {
	roleWithPermissions, err := s.roleRepo.GetWithPermissions(id)
	if err != nil {
		return nil, err
	}

	response := &RoleWithPermissionsResponse{
		RoleResponse: *toRoleResponse(&roleWithPermissions.Role),
		Modules:      []RoleModulePermissionResponse{},
	}

	for _, module := range roleWithPermissions.Modules {
		response.Modules = append(response.Modules, RoleModulePermissionResponse{
			ModuleID:   module.ModuleID,
			ModuleName: module.ModuleName,
			ModuleURL:  module.ModuleURL,
			CanRead:    module.CanRead,
			CanWrite:   module.CanWrite,
			CanDelete:  module.CanDelete,
		})
	}

	return response, nil
}

func (s *Service) CreateRole(req *CreateRoleRequest) (*RoleResponse, error) {
	role := &Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

func (s *Service) UpdateRole(id int64, req *UpdateRoleRequest) (*RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
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

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return toRoleResponse(role), nil
}

func (s *Service) DeleteRole(id int64) error {
	return s.roleRepo.Delete(id)
}

func (s *Service) UpdateRolePermissions(roleID int64, req *UpdateRolePermissionsRequest) error {
	var modules []*RoleModule
	for _, perm := range req.Permissions {
		modules = append(modules, &RoleModule{
			RoleID:    roleID,
			ModuleID:  perm.ModuleID,
			CanRead:   perm.CanRead,
			CanWrite:  perm.CanWrite,
			CanDelete: perm.CanDelete,
		})
	}
	return s.roleRepo.UpdateRoleModules(roleID, modules)
}

func (s *Service) AssignRoleToUser(req *AssignRoleRequest) (*UserRoleAssignmentResponse, error) {
	// Verify user exists (query via repository, no cross-module import)
	userExists, err := s.roleRepo.CheckUserExists(req.UserID)
	if err != nil || !userExists {
		return nil, fmt.Errorf("pengguna dengan ID %d tidak ditemukan", req.UserID)
	}

	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("peran dengan ID %d tidak ditemukan", req.RoleID)
	}

	userRole := &UserRole{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
		BranchID:  req.BranchID,
	}

	if err := s.roleRepo.AssignUserRole(userRole); err != nil {
		return nil, err
	}

	return &UserRoleAssignmentResponse{
		ID:          userRole.ID,
		UserID:      req.UserID,
		RoleID:      role.ID,
		CompanyID:   req.CompanyID,
		BranchID:    req.BranchID,
		RoleName:    role.Name,
		CompanyName: "",
		BranchName:  nil,
		CreatedAt:   userRole.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Service) RemoveRoleFromUser(userID, roleID, companyID int64) error {
	return s.roleRepo.RemoveUserRole(userID, roleID, companyID)
}

func (s *Service) GetUsersByRole(roleID int64, limit int) (interface{}, error) {
	users, err := s.roleRepo.GetUsersByRole(roleID, limit)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, user := range users {
		userIdentity := ""
		if user.UserIdentity != nil {
			userIdentity = *user.UserIdentity
		}

		response = append(response, map[string]interface{}{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"user_identity": userIdentity,
			"is_active":     user.IsActive,
		})
	}

	return map[string]interface{}{
		"users":   response,
		"total":   len(response),
		"role_id": roleID,
	}, nil
}

func (s *Service) GetUserRoles(userID int64) (interface{}, error) {
	userRoles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, userRole := range userRoles {
		role, err := s.roleRepo.GetByID(userRole.RoleID)
		if err != nil {
			continue
		}

		response = append(response, map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
			"is_active":   role.IsActive,
			"created_at":  role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"updated_at":  role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *Service) GetUserAccessSummary(userID int64) (interface{}, error) {
	userRolesInterface, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	userRoles, ok := userRolesInterface.([]map[string]interface{})
	if !ok {
		userRoles = []map[string]interface{}{}
	}

	summary := map[string]interface{}{
		"user_id":     userID,
		"roles":       userRoles,
		"total_roles": len(userRoles),
		"permissions": map[string]interface{}{
			"can_read":   true,
			"can_write":  false,
			"can_delete": false,
		},
	}

	return summary, nil
}

func (s *Service) GetAllUserRoleAssignments() (interface{}, error) {
	assignments, err := s.roleRepo.GetAllUserRoleAssignments()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"assignments": assignments,
		"total":       len(assignments),
	}, nil
}

func (s *Service) GetUserRolesByUserID(userID int64) (interface{}, error) {
	assignments, err := s.roleRepo.GetUserRolesByUserID(userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"user_id":     userID,
		"assignments": assignments,
		"total":       len(assignments),
	}, nil
}

func (s *Service) GetRoleUsersMapping() (interface{}, error) {
	mappings, err := s.roleRepo.GetRoleUsersMapping()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"role_mappings": mappings,
		"total_roles":   len(mappings),
	}, nil
}

// Helper function
func toRoleResponse(role *Role) *RoleResponse {
	if role == nil {
		return nil
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   role.UpdatedAt.Format(time.RFC3339),
	}
}
