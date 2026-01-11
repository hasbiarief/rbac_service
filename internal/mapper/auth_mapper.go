package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// AuthMapper handles conversion between auth models and DTOs
type AuthMapper struct {
	userMapper *UserMapper
}

// NewAuthMapper creates a new auth mapper
func NewAuthMapper() *AuthMapper {
	return &AuthMapper{
		userMapper: NewUserMapper(),
	}
}

// ToLoginResponse creates login response DTO with roles and modules
func (m *AuthMapper) ToLoginResponse(user *models.User, accessToken, refreshToken string, expiresIn int64) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         *m.userMapper.ToResponse(user),
	}
}

// ToLoginResponseWithRolesAndModules creates login response DTO with roles and modules
func (m *AuthMapper) ToLoginResponseWithRolesAndModules(user *models.User, accessToken, refreshToken string, expiresIn int64, roles []string, modules map[string][][]string) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	userResponse := m.userMapper.ToResponse(user)

	// Create enhanced user response with roles and modules
	enhancedUser := dto.UserResponse{
		ID:           userResponse.ID,
		Name:         userResponse.Name,
		Email:        userResponse.Email,
		UserIdentity: userResponse.UserIdentity,
		IsActive:     userResponse.IsActive,
		CreatedAt:    userResponse.CreatedAt,
		UpdatedAt:    userResponse.UpdatedAt,
		Roles:        roles,   // Add roles
		Modules:      modules, // Add modules
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         enhancedUser,
	}
}

// ToRefreshTokenResponse creates refresh token response DTO
func (m *AuthMapper) ToRefreshTokenResponse(accessToken, refreshToken string, expiresIn int64) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}

// ToRegisterResponse creates register response DTO
func (m *AuthMapper) ToRegisterResponse(user *models.User, message string) *dto.RegisterResponse {
	if user == nil {
		return nil
	}

	return &dto.RegisterResponse{
		User:    *m.userMapper.ToResponse(user),
		Message: message,
	}
}

// ToUserModel converts register request DTO to user model
func (m *AuthMapper) ToUserModel(req *dto.RegisterRequest) *models.User {
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

// ToLoginResponseWithUserRoles creates login response DTO with complete user role assignments
func (m *AuthMapper) ToLoginResponseWithUserRoles(user *models.User, userWithRoles map[string]interface{}, accessToken, refreshToken string, expiresIn int64, modules map[string][][]string) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	// Handle nil userWithRoles
	if userWithRoles == nil {
		userWithRoles = map[string]interface{}{
			"role_assignments": []map[string]interface{}{},
			"total_roles":      0,
		}
	}

	// Extract role assignments from userWithRoles
	roleAssignments, ok := userWithRoles["role_assignments"].([]map[string]interface{})
	if !ok || roleAssignments == nil {
		roleAssignments = []map[string]interface{}{}
	}
	totalRoles, ok := userWithRoles["total_roles"].(int)
	if !ok {
		totalRoles = 0
	}

	// Create enhanced user response with complete role assignments
	enhancedUser := dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		UserIdentity: user.UserIdentity,
		IsActive:     user.IsActive,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
		Modules:      modules, // Add modules
	}

	// Add role assignments to user response
	enhancedUser.RoleAssignments = roleAssignments
	enhancedUser.TotalRoles = totalRoles

	// Ensure role_assignments is never nil in JSON output
	if enhancedUser.RoleAssignments == nil {
		enhancedUser.RoleAssignments = []map[string]interface{}{}
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         enhancedUser,
	}
}
