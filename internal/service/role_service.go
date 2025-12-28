package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

type RoleListRequest struct {
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
	Search string `form:"search"`
}

type RoleResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssignUserRoleRequest struct {
	UserID    int64  `json:"user_id" binding:"required"`
	RoleID    int64  `json:"role_id" binding:"required"`
	CompanyID int64  `json:"company_id" binding:"required"`
	BranchID  *int64 `json:"branch_id"`
}

type BulkAssignRolesRequest struct {
	Assignments []AssignUserRoleRequest `json:"assignments" binding:"required"`
}

type UpdateRoleModulesRequest struct {
	Modules []RoleModulePermission `json:"modules" binding:"required"`
}

type RoleModulePermission struct {
	ModuleID  int64 `json:"module_id" binding:"required"`
	CanRead   bool  `json:"can_read"`
	CanWrite  bool  `json:"can_write"`
	CanDelete bool  `json:"can_delete"`
}

type UserResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	UserIdentity string `json:"user_identity"`
	IsActive     bool   `json:"is_active"`
}

func (s *RoleService) GetRoles(req *RoleListRequest) ([]*RoleResponse, error) {
	roles, err := s.roleRepo.GetAll(req.Limit, req.Offset, req.Search)
	if err != nil {
		return nil, err
	}

	var response []*RoleResponse
	for _, role := range roles {
		response = append(response, &RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *RoleService) GetRoleByID(id int64) (*RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *RoleService) CreateRole(req *CreateRoleRequest) (*RoleResponse, error) {
	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *RoleService) UpdateRole(id int64, req *UpdateRoleRequest) (*RoleResponse, error) {
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

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *RoleService) DeleteRole(id int64) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) AssignUserRole(req *AssignUserRoleRequest) error {
	return s.roleRepo.AssignUserRole(req.UserID, req.RoleID, req.CompanyID, req.BranchID)
}

func (s *RoleService) RemoveUserRole(userID, roleID int64) error {
	return s.roleRepo.RemoveUserRole(userID, roleID)
}

func (s *RoleService) GetUsersByRole(roleID int64, limit int) ([]*UserResponse, error) {
	users, err := s.roleRepo.GetUsersByRole(roleID, limit)
	if err != nil {
		return nil, err
	}

	var response []*UserResponse
	for _, user := range users {
		userIdentity := ""
		if user.UserIdentity != nil {
			userIdentity = *user.UserIdentity
		}

		response = append(response, &UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			UserIdentity: userIdentity,
			IsActive:     user.IsActive,
		})
	}

	return response, nil
}

func (s *RoleService) GetUserRoles(userID int64) ([]*RoleResponse, error) {
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	var response []*RoleResponse
	for _, role := range roles {
		response = append(response, &RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			IsActive:    role.IsActive,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *RoleService) GetUserAccessSummary(userID int64) (map[string]interface{}, error) {
	// Get user roles
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// Get user modules (this would need to be implemented in user service)
	// For now, return basic summary
	summary := map[string]interface{}{
		"user_id":     userID,
		"roles":       roles,
		"total_roles": len(roles),
		"permissions": map[string]interface{}{
			"can_read":   true,  // This would be calculated based on role modules
			"can_write":  false, // This would be calculated based on role modules
			"can_delete": false, // This would be calculated based on role modules
		},
	}

	return summary, nil
}
func (s *RoleService) BulkAssignRoles(req *BulkAssignRolesRequest) error {
	for _, assignment := range req.Assignments {
		if err := s.roleRepo.AssignUserRole(assignment.UserID, assignment.RoleID, assignment.CompanyID, assignment.BranchID); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoleService) UpdateRoleModules(roleID int64, req *UpdateRoleModulesRequest) error {
	var modules []models.RoleModule
	for _, perm := range req.Modules {
		modules = append(modules, models.RoleModule{
			RoleID:    roleID,
			ModuleID:  perm.ModuleID,
			CanRead:   perm.CanRead,
			CanWrite:  perm.CanWrite,
			CanDelete: perm.CanDelete,
		})
	}
	return s.roleRepo.UpdateRoleModules(roleID, modules)
}
