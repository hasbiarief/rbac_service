package service

import (
	"fmt"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
)

type RoleService struct {
	roleRepo   interfaces.RoleRepositoryInterface
	userRepo   interfaces.UserRepositoryInterface
	roleMapper *mapper.RoleMapper
}

func NewRoleService(roleRepo interfaces.RoleRepositoryInterface, userRepo interfaces.UserRepositoryInterface) *RoleService {
	return &RoleService{
		roleRepo:   roleRepo,
		userRepo:   userRepo,
		roleMapper: mapper.NewRoleMapper(),
	}
}

func (s *RoleService) GetRoles(req *dto.RoleListRequest) (*dto.RoleListResponse, error) {
	// Set default values jika tidak disediakan
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

	// Dapatkan total count untuk pagination
	total, err := s.roleRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Konversi ke DTO menggunakan mapper
	var roleResponses []*dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, s.roleMapper.ToResponse(role))
	}

	return &dto.RoleListResponse{
		Data:    roleResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(roleResponses)) < total,
	}, nil
}

func (s *RoleService) GetRoleByID(id int64) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.roleMapper.ToResponse(role), nil
}

func (s *RoleService) GetRoleWithPermissions(id int64) (*dto.RoleWithPermissionsResponse, error) {
	roleWithPermissions, err := s.roleRepo.GetWithPermissions(id)
	if err != nil {
		return nil, err
	}

	return s.roleMapper.ToWithPermissionsResponse(roleWithPermissions), nil
}

func (s *RoleService) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// Konversi DTO ke model menggunakan mapper
	role := s.roleMapper.ToModel(req)

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return s.roleMapper.ToResponse(role), nil
}

func (s *RoleService) UpdateRole(id int64, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields menggunakan mapper
	s.roleMapper.UpdateModel(role, req)

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	return s.roleMapper.ToResponse(role), nil
}

func (s *RoleService) DeleteRole(id int64) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) UpdateRolePermissions(roleID int64, req *dto.UpdateRolePermissionsRequest) error {
	var modules []*models.RoleModule
	for _, perm := range req.Permissions {
		modules = append(modules, &models.RoleModule{
			RoleID:    roleID,
			ModuleID:  perm.ModuleID,
			CanRead:   perm.CanRead,
			CanWrite:  perm.CanWrite,
			CanDelete: perm.CanDelete,
		})
	}
	return s.roleRepo.UpdateRoleModules(roleID, modules)
}

func (s *RoleService) AssignRoleToUser(req *dto.AssignRoleRequest) (*dto.UserRoleAssignmentResponse, error) {
	// Pertama, periksa apakah user ada
	user, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("pengguna dengan ID %d tidak ditemukan", req.UserID)
	}

	// Periksa apakah role ada
	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("peran dengan ID %d tidak ditemukan", req.RoleID)
	}

	userRole := &models.UserRole{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		CompanyID: req.CompanyID,
		BranchID:  req.BranchID,
	}

	if err := s.roleRepo.AssignUserRole(userRole); err != nil {
		return nil, err
	}

	// Buat response
	return &dto.UserRoleAssignmentResponse{
		ID:          0, // Akan diisi oleh database
		UserID:      user.ID,
		RoleID:      role.ID,
		CompanyID:   req.CompanyID,
		BranchID:    req.BranchID,
		RoleName:    role.Name,
		CompanyName: "",    // Perlu diambil dari company service jika diperlukan
		BranchName:  nil,   // Perlu diambil dari branch service jika diperlukan
		CreatedAt:   "now", // Atau gunakan timestamp yang sesuai
	}, nil
}

func (s *RoleService) RemoveRoleFromUser(userID, roleID, companyID int64) error {
	return s.roleRepo.RemoveUserRole(userID, roleID, companyID)
}

func (s *RoleService) GetUsersByRole(roleID int64, limit int) (interface{}, error) {
	users, err := s.roleRepo.GetUsersByRole(roleID, limit)
	if err != nil {
		return nil, err
	}

	// Always return an array, even if empty
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

	// Return the response with metadata
	return map[string]interface{}{
		"users":   response,
		"total":   len(response),
		"role_id": roleID,
	}, nil
}

func (s *RoleService) GetUserRoles(userID int64) (interface{}, error) {
	userRoles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, userRole := range userRoles {
		// Dapatkan detail role dari role repository
		role, err := s.roleRepo.GetByID(userRole.RoleID)
		if err != nil {
			continue // Skip jika role tidak ditemukan
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

func (s *RoleService) GetUserAccessSummary(userID int64) (interface{}, error) {
	// Dapatkan user roles
	userRolesInterface, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	userRoles, ok := userRolesInterface.([]map[string]interface{})
	if !ok {
		userRoles = []map[string]interface{}{}
	}

	// Dapatkan user modules (ini perlu diimplementasikan di user service)
	// Untuk saat ini, kembalikan ringkasan dasar
	summary := map[string]interface{}{
		"user_id":     userID,
		"roles":       userRoles,
		"total_roles": len(userRoles),
		"permissions": map[string]interface{}{
			"can_read":   true,  // Ini akan dihitung berdasarkan role modules
			"can_write":  false, // Ini akan dihitung berdasarkan role modules
			"can_delete": false, // Ini akan dihitung berdasarkan role modules
		},
	}

	return summary, nil
}

// GetAllUserRoleAssignments - Debug method to see all user role assignments
func (s *RoleService) GetAllUserRoleAssignments() (interface{}, error) {
	assignments, err := s.roleRepo.GetAllUserRoleAssignments()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"assignments": assignments,
		"total":       len(assignments),
	}, nil
}

// GetUserRolesByUserID - Debug method to check specific user's role assignments
func (s *RoleService) GetUserRolesByUserID(userID int64) (interface{}, error) {
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

// GetRoleUsersMapping - Debug method to show role-user mapping
func (s *RoleService) GetRoleUsersMapping() (interface{}, error) {
	mappings, err := s.roleRepo.GetRoleUsersMapping()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"role_mappings": mappings,
		"total_roles":   len(mappings),
	}, nil
}
