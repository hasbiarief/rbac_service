package user

import (
	"fmt"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/rbac"
	"time"
)

type Service struct {
	userRepo    *UserRepository
	rbacService *rbac.RBACService
}

func NewService(userRepo *UserRepository, rbacService *rbac.RBACService) *Service {
	return &Service{
		userRepo:    userRepo,
		rbacService: rbacService,
	}
}

func (s *Service) GetUsers(req *UserListRequest) (*UserListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.userRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	return &UserListResponse{
		Data:    userResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(userResponses)) < total,
	}, nil
}

func (s *Service) GetUsersFiltered(requestingUserID int64, req *UserListRequest) (*UserListResponse, error) {
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		return s.GetUsers(req)
	}

	isHRAdmin, err := s.rbacService.HasRole(requestingUserID, "HR_ADMIN")
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa peran HR_ADMIN: %w", err)
	}

	isHRManager, err := s.rbacService.HasRole(requestingUserID, "HR_MANAGER")
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa peran HR_MANAGER: %w", err)
	}

	if isHRAdmin || isHRManager {
		return s.GetUsers(req)
	}

	return &UserListResponse{
		Data:    []*UserResponse{},
		Total:   0,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: false,
	}, nil
}

func (s *Service) GetUserByID(id int64) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *Service) GetUserWithRoles(id int64) (*UserWithRolesResponse, error) {
	user, err := s.userRepo.GetWithRoles(id)
	if err != nil {
		return nil, err
	}

	return &UserWithRolesResponse{
		UserResponse: *toUserResponse(user),
		Roles:        []UserRoleResponse{},
	}, nil
}

func (s *Service) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	hashedPassword := ""
	if req.Password != "" {
		hash, err := password.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = hash
	} else {
		hash, err := password.HashPassword("password123")
		if err != nil {
			return nil, err
		}
		hashedPassword = hash
	}

	user := &User{
		Name:         req.Name,
		Email:        req.Email,
		UserIdentity: req.UserIdentity,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *Service) UpdateUser(id int64, req *UpdateUserRequest) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.UserIdentity != nil {
		user.UserIdentity = req.UserIdentity
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *Service) DeleteUser(id int64) error {
	return s.userRepo.Delete(id)
}

func (s *Service) ChangePassword(userID int64, req *ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("password baru dan konfirmasi password tidak cocok")
	}

	if req.CurrentPassword == req.NewPassword {
		return fmt.Errorf("password baru harus berbeda dari password saat ini")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := password.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		return fmt.Errorf("password saat ini salah")
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *Service) ChangeUserPassword(userID int64, req *ChangePasswordRequest) error {
	return s.ChangePassword(userID, req)
}

func (s *Service) GetUserByIDWithRoles(id int64) (interface{}, error) {
	return s.userRepo.GetByIDWithRoles(id)
}

func (s *Service) GetUsersWithRoles(req *UserListRequest) (interface{}, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	users, err := s.userRepo.GetAllWithRoles(req.Limit, req.Offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.userRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"users":  users,
		"total":  total,
		"limit":  req.Limit,
		"offset": req.Offset,
	}, nil
}

// Helper function to convert model to DTO
func toUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		UserIdentity:    user.UserIdentity,
		IsActive:        user.IsActive,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       user.UpdatedAt.Format(time.RFC3339),
		RoleAssignments: []map[string]interface{}{},
		TotalRoles:      0,
	}
}
