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
	tokenService *token.SimpleTokenService
	jwtSecret    string
}

func NewAuthService(userRepo *repository.UserRepository, tokenService *token.SimpleTokenService, jwtSecret string) *AuthService {
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
	ID           int64                 `json:"id"`
	Name         string                `json:"name"`
	Email        string                `json:"email"`
	UserIdentity string                `json:"user_identity"`
	IsActive     bool                  `json:"is_active"`
	Roles        []string              `json:"roles"`
	Modules      map[string][][]string `json:"modules"`
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

	// Get user modules with subscription filtering (grouped by category)
	modules, err := s.userRepo.GetUserModulesGroupedWithSubscription(user.ID)
	if err != nil {
		modules = make(map[string][][]string) // Default to empty if error
	}

	// Get user modules as URLs for token abilities
	moduleURLs, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{} // Default to empty if error
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

	// Store access token (overwrites any existing token for this user)
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: userAgent,
		IP:        ip,
		Abilities: moduleURLs, // Use URLs for abilities
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Store refresh token (overwrites any existing token for this user)
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

	modules, err := s.userRepo.GetUserModulesGroupedWithSubscription(user.ID)
	if err != nil {
		modules = make(map[string][][]string)
	}

	// Get user modules as URLs for token abilities
	moduleURLs, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{}
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

	// Store new access token (overwrites existing)
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: "",         // We don't have user agent in refresh
		IP:        "",         // We don't have IP in refresh
		Abilities: moduleURLs, // Use URLs for abilities
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Store new refresh token (overwrites existing) with same family ID
	newRefreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: refreshMetadata.FamilyID,
	}

	if err := s.tokenService.StoreRefreshToken(newRefreshToken, newRefreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

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
	// Try to get token metadata to find user ID
	metadata, err := s.tokenService.GetAccessToken(accessToken)
	if err != nil {
		// If token is not found (expired/invalid), we can't determine user ID
		// In this case, we'll return success since the token is already gone
		return nil
	}

	// Delete ALL tokens for this user (access and refresh tokens)
	if err := s.tokenService.RevokeAllUserTokens(metadata.UserID); err != nil {
		return err
	}

	return nil
}

// LogoutByUserID logs out user by user ID (alternative when token is expired)
func (s *AuthService) LogoutByUserID(userID int64) error {
	// Delete ALL tokens for this user (access and refresh tokens)
	if err := s.tokenService.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	return nil
}

// CheckUserTokens checks if user has valid tokens in Redis
func (s *AuthService) CheckUserTokens(userID int64) (*token.UserTokensResponse, error) {
	return s.tokenService.GetUserTokens(userID)
}

// CleanupExpiredTokens removes expired tokens from Redis
func (s *AuthService) CleanupExpiredTokens() error {
	return s.tokenService.CleanupExpiredTokens()
}

// GetUserSessionCount returns the number of active sessions for a user
func (s *AuthService) GetUserSessionCount(userID int64) (int, error) {
	return s.tokenService.GetUserSessionCount(userID)
}

// GetUserRefreshTokenCount returns the number of active refresh tokens for a user
func (s *AuthService) GetUserRefreshTokenCount(userID int64) (int, error) {
	// Get all refresh tokens for user
	tokensResponse, err := s.tokenService.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}

	return len(tokensResponse.RefreshTokens), nil
}
