package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/token"
	"time"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	tokenService *token.TokenService
	jwtSecret    string
}

func NewAuthService(userRepo *repository.UserRepository, tokenService *token.TokenService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		jwtSecret:    jwtSecret,
	}
}

type LoginRequest struct {
	UserIdentity string `json:"user_identity" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type UserInfo struct {
	ID           int64    `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	UserIdentity string   `json:"user_identity"`
	IsActive     bool     `json:"is_active"`
	Roles        []string `json:"roles"`
	Modules      []string `json:"modules"`
}

func (s *AuthService) Login(req *LoginRequest, userAgent, ip string) (*LoginResponse, error) {
	// Get user by user_identity (primary method as per Postman documentation)
	user, err := s.userRepo.GetByUserIdentity(req.UserIdentity)
	if err != nil {
		// Fallback: try by email for backward compatibility
		user, err = s.userRepo.GetByEmail(req.UserIdentity)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password using bcrypt
	if err := password.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Get user roles and modules
	roles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{} // Default to empty if error
	}

	// Get user modules with subscription filtering
	modules, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		modules = []string{} // Default to empty if error
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Generate family ID for refresh token
	familyID, err := s.generateFamilyID()
	if err != nil {
		return nil, err
	}

	// Set token expiration
	expiresAt := time.Now().Add(15 * time.Minute) // 15 minutes for access token

	// Store access token
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: userAgent,
		IP:        ip,
		Abilities: modules, // Use modules as abilities
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Store refresh token
	refreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: familyID,
	}

	if err := s.tokenService.StoreRefreshToken(refreshToken, refreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	// Handle UserIdentity pointer
	userIdentity := ""
	if user.UserIdentity != nil {
		userIdentity = *user.UserIdentity
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserInfo{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			UserIdentity: userIdentity,
			IsActive:     user.IsActive,
			Roles:        roles,
			Modules:      modules,
		},
	}, nil
}

func (s *AuthService) generateFamilyID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// Get refresh token metadata
	refreshMetadata, err := s.tokenService.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user by ID
	user, err := s.userRepo.GetByID(refreshMetadata.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Get user roles and modules
	roles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{}
	}

	modules, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		modules = []string{}
	}

	// Generate new access token
	accessToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Set token expiration
	expiresAt := time.Now().Add(15 * time.Minute)

	// Store new access token
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: "", // We don't have user agent in refresh
		IP:        "", // We don't have IP in refresh
		Abilities: modules,
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Store new refresh token with same family ID
	newRefreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: refreshMetadata.FamilyID,
	}

	if err := s.tokenService.StoreRefreshToken(newRefreshToken, newRefreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	// Delete old refresh token
	s.tokenService.RevokeToken(refreshToken, "refresh")

	// Handle UserIdentity pointer
	userIdentity := ""
	if user.UserIdentity != nil {
		userIdentity = *user.UserIdentity
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		User: UserInfo{
			ID:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			UserIdentity: userIdentity,
			IsActive:     user.IsActive,
			Roles:        roles,
			Modules:      modules,
		},
	}, nil
}

func (s *AuthService) Logout(accessToken string) error {
	// Validate that the token exists before revoking it
	_, err := s.tokenService.GetAccessToken(accessToken)
	if err != nil {
		return errors.New("invalid token")
	}

	// Delete access token
	if err := s.tokenService.RevokeToken(accessToken, "access"); err != nil {
		return err
	}

	// Optionally, you could also delete all refresh tokens for this user
	// This would log them out from all devices
	// For now, we'll just delete the access token

	return nil
}
