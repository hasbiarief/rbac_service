package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/token"
	"time"
)

type Service struct {
	repo         *Repository
	tokenService *token.SimpleTokenService
	jwtSecret    string
}

func NewService(repo *Repository, tokenService *token.SimpleTokenService, jwtSecret string) *Service {
	return &Service{
		repo:         repo,
		tokenService: tokenService,
		jwtSecret:    jwtSecret,
	}
}

func (s *Service) Login(req *LoginRequest, userAgent, ip string) (*LoginResponse, error) {
	// Get user by user_identity
	user, err := s.repo.GetByUserIdentity(req.UserIdentity)
	if err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	if !user.IsActive {
		return nil, errors.New("akun pengguna tidak aktif")
	}

	if err := password.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	// UserAgent and IP are now passed as parameters
	userWithRoles, err := s.repo.GetByIDWithRoles(user.ID)
	if err != nil || userWithRoles == nil {
		userWithRoles = map[string]interface{}{
			"role_assignments": []map[string]interface{}{},
			"total_roles":      0,
		}
	}

	if userWithRoles["role_assignments"] == nil {
		userWithRoles["role_assignments"] = []map[string]interface{}{}
	}
	if userWithRoles["total_roles"] == nil {
		userWithRoles["total_roles"] = 0
	}

	// Get applications with modules (new hierarchical structure)
	applications, err := s.repo.GetUserApplicationsWithModules(user.ID)
	if err != nil {
		applications = make(map[string]interface{})
	}

	// Extract application codes for simplified login response
	var applicationCodes []string
	for appCode := range applications {
		applicationCodes = append(applicationCodes, appCode)
	}

	moduleURLs, err := s.repo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{}
	}

	// Get subscription information
	subscriptionInfo, err := s.repo.GetUserSubscriptionInfo(user.ID)
	if err != nil {
		// If subscription info fails, provide basic info
		subscriptionInfo = map[string]interface{}{
			"has_subscription": false,
			"message":          "Unable to retrieve subscription information",
		}
	}

	accessToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	familyID, err := s.generateFamilyID()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: userAgent,
		IP:        ip,
		Abilities: moduleURLs,
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	refreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: familyID,
	}

	if err := s.tokenService.StoreRefreshToken(refreshToken, refreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	expiresIn := int64(15 * 60)

	// Build simplified login response
	loginData := map[string]interface{}{
		"user_identity": user.UserIdentity,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"applications":  applicationCodes, // Simple array of application codes
		"subscription":  subscriptionInfo,
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         loginData,
	}, nil
}

func (s *Service) LoginWithEmail(req *LoginEmailRequest, userAgent, ip string) (*LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	if !user.IsActive {
		return nil, errors.New("akun pengguna tidak aktif")
	}

	if err := password.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	// Convert to LoginRequest and call Login
	loginReq := &LoginRequest{
		UserIdentity: *user.UserIdentity,
		Password:     req.Password,
	}

	return s.Login(loginReq, userAgent, ip)
}

func (s *Service) generateFamilyID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *Service) RefreshToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	refreshMetadata, err := s.tokenService.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau kedaluwarsa")
	}

	user, err := s.repo.GetByID(refreshMetadata.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	if !user.IsActive {
		return nil, errors.New("akun pengguna tidak aktif")
	}

	userRoles, err := s.repo.GetUserRoles(user.ID)
	if err != nil {
		userRoles = []*UserRole{}
	}

	var roles []string
	for _, userRole := range userRoles {
		roles = append(roles, fmt.Sprintf("role_%d", userRole.RoleID))
	}

	moduleURLs, err := s.repo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{}
	}

	accessToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: "",
		IP:        "",
		Abilities: moduleURLs,
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	newRefreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: refreshMetadata.FamilyID,
	}

	if err := s.tokenService.StoreRefreshToken(newRefreshToken, newRefreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	expiresIn := int64(15 * 60)

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *Service) Logout(accessToken string) error {
	metadata, err := s.tokenService.GetAccessToken(accessToken)
	if err != nil {
		return nil
	}

	if err := s.tokenService.RevokeAllUserTokens(metadata.UserID); err != nil {
		return err
	}

	return nil
}

func (s *Service) LogoutByUserID(userID int64) error {
	if err := s.tokenService.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckUserTokens(userID int64) (interface{}, error) {
	return s.tokenService.GetUserTokens(userID)
}

func (s *Service) CleanupExpiredTokens() error {
	return s.tokenService.CleanupExpiredTokens()
}

func (s *Service) GetUserSessionCount(userID int64) (int64, error) {
	count, err := s.tokenService.GetUserSessionCount(userID)
	return int64(count), err
}

func (s *Service) GetUserRefreshTokenCount(userID int64) (int64, error) {
	tokensResponse, err := s.tokenService.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}

	return int64(len(tokensResponse.RefreshTokens)), nil
}
func (s *Service) GetUserProfile(req *ProfileRequest) (*ProfileResponse, error) {
	userProfile, err := s.repo.GetUserProfileByApplication(req.UserIdentity, req.ApplicationCode)
	if err != nil {
		return nil, err
	}

	return &ProfileResponse{
		User: userProfile,
	}, nil
}
