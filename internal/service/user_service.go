package service

import (
	"fmt"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/rbac"
)

type UserService struct {
	userRepo    interfaces.UserRepositoryInterface
	rbacService *rbac.RBACService
	userMapper  *mapper.UserMapper
}

func NewUserService(userRepo interfaces.UserRepositoryInterface, rbacService *rbac.RBACService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		rbacService: rbacService,
		userMapper:  mapper.NewUserMapper(),
	}
}

func (s *UserService) GetUsers(req *dto.UserListRequest) (*dto.UserListResponse, error) {
	// Set default values
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Dapatkan users dari repository
	users, err := s.userRepo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Dapatkan total count untuk pagination
	total, err := s.userRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Konversi ke response format menggunakan mapper
	var userResponses []*dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, s.userMapper.ToResponse(user))
	}

	return &dto.UserListResponse{
		Data:    userResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(userResponses)) < total,
	}, nil
}

// GetUsersFiltered mengembalikan users yang difilter berdasarkan izin requesting user
func (s *UserService) GetUsersFiltered(requestingUserID int64, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	// Periksa apakah requesting user adalah super admin - jika ya, kembalikan semua users
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		return s.GetUsers(req)
	}

	// Untuk non-super admin users, implementasikan filtering berbasis peran
	// Untuk saat ini, mari periksa apakah mereka memiliki peran HR_ADMIN atau HR_MANAGER
	isHRAdmin, err := s.rbacService.HasRole(requestingUserID, "HR_ADMIN")
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa peran HR_ADMIN: %w", err)
	}

	isHRManager, err := s.rbacService.HasRole(requestingUserID, "HR_MANAGER")
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa peran HR_MANAGER: %w", err)
	}

	// HR_ADMIN dan HR_MANAGER dapat melihat semua users
	if isHRAdmin || isHRManager {
		return s.GetUsers(req)
	}

	// Peran lain hanya dapat melihat profil mereka sendiri
	// Untuk saat ini, kembalikan list kosong untuk peran lain
	// Dalam implementasi nyata, Anda mungkin ingin mengembalikan hanya data requesting user
	return &dto.UserListResponse{
		Data:    []*dto.UserResponse{},
		Total:   0,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: false,
	}, nil
}

func (s *UserService) GetUserByID(id int64) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) GetUserWithRoles(id int64) (*dto.UserWithRolesResponse, error) {
	user, err := s.userRepo.GetWithRoles(id)
	if err != nil {
		return nil, err
	}

	return s.userMapper.ToWithRolesResponse(user), nil
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
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

	// Konversi DTO ke model menggunakan mapper
	user := s.userMapper.ToModel(req)
	user.PasswordHash = hashedPassword

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) UpdateUser(id int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields menggunakan mapper
	s.userMapper.UpdateModel(user, req)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.userMapper.ToResponse(user), nil
}

func (s *UserService) DeleteUser(id int64) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ChangePassword(userID int64, req *dto.ChangePasswordRequest) error {
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

	// Verifikasi password saat ini
	if err := password.VerifyPassword(user.PasswordHash, req.CurrentPassword); err != nil {
		return fmt.Errorf("password saat ini salah")
	}

	// Hash password baru
	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *UserService) ChangeUserPassword(userID int64, req *dto.ChangePasswordRequest) error {
	return s.ChangePassword(userID, req)
}
