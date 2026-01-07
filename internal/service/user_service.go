package service

import (
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/rbac"
)

type UserService struct {
	userRepo    *repository.UserRepository
	rbacService *rbac.RBACService
}

func NewUserService(userRepo *repository.UserRepository, rbacService *rbac.RBACService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		rbacService: rbacService,
	}
}

type UserListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

type UserResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	UserIdentity string `json:"user_identity"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserListResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

type CreateUserRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	UserIdentity string `json:"user_identity" binding:"required"`
	Password     string `json:"password"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive *bool  `json:"is_active"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (s *UserService) GetUsers(req *UserListRequest) (*UserListResponse, error) {
	// Set default values
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Get users from repository
	users, err := s.userRepo.GetAll(req.Limit, req.Offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	total, err := s.userRepo.GetCount(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	var userResponses []*UserResponse
	for _, user := range users {
		userIdentity := ""
		if user.UserIdentity != nil {
			userIdentity = *user.UserIdentity
		}

		userResponses = append(userResponses, &UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			UserIdentity: userIdentity,
			IsActive:     user.IsActive,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Calculate pagination info
	page := (req.Offset / req.Limit) + 1
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetUsersFiltered returns users filtered by requesting user's permissions
func (s *UserService) GetUsersFiltered(requestingUserID int64, req *UserListRequest) (*UserListResponse, error) {
	// Check if requesting user is super admin - if so, return all users
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetUsers(req)
	}

	// For non-super admin users, implement role-based filtering
	// For now, let's check if they have HR_ADMIN or HR_MANAGER role
	isHRAdmin, err := s.rbacService.HasRole(requestingUserID, "HR_ADMIN")
	if err != nil {
		return nil, fmt.Errorf("failed to check HR_ADMIN role: %w", err)
	}

	isHRManager, err := s.rbacService.HasRole(requestingUserID, "HR_MANAGER")
	if err != nil {
		return nil, fmt.Errorf("failed to check HR_MANAGER role: %w", err)
	}

	// HR_ADMIN and HR_MANAGER can see all users
	if isHRAdmin || isHRManager {
		return s.GetUsers(req)
	}

	// Other roles can only see their own profile
	// For now, return empty list for other roles
	// In a real implementation, you might want to return only the requesting user's data
	return &UserListResponse{
		Users:      []*UserResponse{},
		Total:      0,
		Page:       1,
		Limit:      req.Limit,
		TotalPages: 0,
	}, nil
}

func (s *UserService) GetUserByID(id int64) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	userIdentity := ""
	if user.UserIdentity != nil {
		userIdentity = *user.UserIdentity
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		UserIdentity: userIdentity,
		IsActive:     user.IsActive,
	}, nil
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	// Hash password
	hashedPassword := ""
	if req.Password != "" {
		hash, err := password.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		hashedPassword = hash
	} else {
		// Default password
		hash, err := password.HashPassword("password123")
		if err != nil {
			return nil, err
		}
		hashedPassword = hash
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		UserIdentity: &req.UserIdentity,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		UserIdentity: req.UserIdentity,
		IsActive:     user.IsActive,
	}, nil
}

func (s *UserService) UpdateUser(id int64, req *UpdateUserRequest) (*UserResponse, error) {
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
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	userIdentity := ""
	if user.UserIdentity != nil {
		userIdentity = *user.UserIdentity
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		UserIdentity: userIdentity,
		IsActive:     user.IsActive,
	}, nil
}

func (s *UserService) DeleteUser(id int64) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ChangePassword(userID int64, req *ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("new password and confirm password do not match")
	}

	if req.CurrentPassword == req.NewPassword {
		return fmt.Errorf("new password must be different from current password")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := password.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserService) ChangeUserPassword(userID int64, req *ChangePasswordRequest) error {
	return s.ChangePassword(userID, req)
}
