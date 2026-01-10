package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// UserMapper handles conversion between user models and DTOs
type UserMapper struct{}

// NewUserMapper creates a new user mapper
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToResponse converts model to response DTO
func (m *UserMapper) ToResponse(user *models.User) *dto.UserResponse {
	if user == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		UserIdentity: user.UserIdentity,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
	}
}

// ToResponseList converts model slice to response DTO slice
func (m *UserMapper) ToResponseList(users []*models.User) []*dto.UserResponse {
	if users == nil {
		return nil
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = m.ToResponse(user)
	}
	return responses
}

// ToModel converts create request DTO to model
func (m *UserMapper) ToModel(req *dto.CreateUserRequest) *models.User {
	if req == nil {
		return nil
	}

	return &models.User{
		Name:         req.Name,
		Email:        req.Email,
		UserIdentity: req.UserIdentity,
		IsActive:     true, // Default to active
	}
}

// UpdateModel updates model with update request DTO
func (m *UserMapper) UpdateModel(user *models.User, req *dto.UpdateUserRequest) {
	if user == nil || req == nil {
		return
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
}

// ToListResponse creates paginated list response
func (m *UserMapper) ToListResponse(users []*models.User, total int64, limit, offset int) *dto.UserListResponse {
	return &dto.UserListResponse{
		Data:    m.ToResponseList(users),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(users)) < total,
	}
}

// ToWithRolesResponse converts user with roles model to response DTO
func (m *UserMapper) ToWithRolesResponse(user *models.User) *dto.UserWithRolesResponse {
	if user == nil {
		return nil
	}

	// For now, return basic user response with empty roles
	// This would need to be implemented properly with actual role data
	return &dto.UserWithRolesResponse{
		UserResponse: *m.ToResponse(user),
		Roles:        []dto.UserRoleResponse{}, // Empty for now
	}
}
