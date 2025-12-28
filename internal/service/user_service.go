package service

import (
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/pkg/password"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type UserListRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
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

func (s *UserService) GetUsers(req *UserListRequest) ([]*UserResponse, error) {
	// For now, we'll implement a simple query. In a real system, you'd want pagination
	// This is a placeholder implementation
	users := []*UserResponse{} // We'll need to implement GetAll in UserRepository
	return users, nil
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
