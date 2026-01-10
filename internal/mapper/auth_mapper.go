package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
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
