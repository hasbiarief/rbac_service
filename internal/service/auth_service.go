package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/password"
	"gin-scalable-api/pkg/token"
	"time"
)

type AuthService struct {
	userRepo     interfaces.UserRepositoryInterface
	tokenService *token.SimpleTokenService
	jwtSecret    string
	authMapper   *mapper.AuthMapper
	userMapper   *mapper.UserMapper
}

func NewAuthService(userRepo interfaces.UserRepositoryInterface, tokenService *token.SimpleTokenService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		jwtSecret:    jwtSecret,
		authMapper:   mapper.NewAuthMapper(),
		userMapper:   mapper.NewUserMapper(),
	}
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Tentukan user identity dari request
	userIdentity := ""
	if req.UserIdentity != "" {
		userIdentity = req.UserIdentity
	} else if req.Email != "" {
		userIdentity = req.Email
	} else {
		return nil, errors.New("user identity atau email harus diisi")
	}

	// Dapatkan user berdasarkan user_identity (metode utama sesuai dokumentasi Postman)
	user, err := s.userRepo.GetByUserIdentity(userIdentity)
	if err != nil {
		// Fallback: coba berdasarkan email untuk kompatibilitas mundur
		user, err = s.userRepo.GetByEmail(userIdentity)
		if err != nil {
			return nil, errors.New("kredensial tidak valid")
		}
	}

	// Periksa apakah user aktif
	if !user.IsActive {
		return nil, errors.New("akun pengguna tidak aktif")
	}

	// Verifikasi password menggunakan bcrypt
	if err := password.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("kredensial tidak valid")
	}

	// Dapatkan user agent dan IP dari request
	userAgent := ""
	ip := ""
	if req.UserAgent != nil {
		userAgent = *req.UserAgent
	}
	if req.IP != nil {
		ip = *req.IP
	}

	// Dapatkan user dengan role assignments lengkap
	userWithRoles, err := s.userRepo.GetByIDWithRoles(user.ID)
	if err != nil || userWithRoles == nil {
		// If error getting role assignments or nil result, create empty structure
		userWithRoles = map[string]interface{}{
			"role_assignments": []map[string]interface{}{},
			"total_roles":      0,
		}
	}

	// Ensure role_assignments is never nil
	if userWithRoles["role_assignments"] == nil {
		userWithRoles["role_assignments"] = []map[string]interface{}{}
	}
	if userWithRoles["total_roles"] == nil {
		userWithRoles["total_roles"] = 0
	}

	// Dapatkan modul user dengan unit-aware enhancement (fallback to traditional for now)
	modules, err := s.userRepo.GetUserModulesGroupedWithSubscription(user.ID)
	if err != nil {
		modules = make(map[string][][]string) // Default ke kosong jika error
	}

	moduleURLs, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{} // Default ke kosong jika error
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

	// Generate family ID untuk refresh token
	familyID, err := s.generateFamilyID()
	if err != nil {
		return nil, err
	}

	// Set token expiration
	expiresAt := time.Now().Add(15 * time.Minute) // 15 menit untuk access token

	// Simpan access token (menimpa token yang ada untuk user ini)
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: userAgent,
		IP:        ip,
		Abilities: moduleURLs, // Gunakan URL untuk kemampuan
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Simpan refresh token (menimpa token yang ada untuk user ini)
	refreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: familyID,
	}

	if err := s.tokenService.StoreRefreshToken(refreshToken, refreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	// Hitung expires in dalam detik
	expiresIn := int64(15 * 60) // 15 menit dalam detik

	// Gunakan mapper untuk membuat response dengan user data lengkap termasuk role assignments
	return s.authMapper.ToLoginResponseWithUserRoles(user, userWithRoles, accessToken, refreshToken, expiresIn, modules), nil
}

func (s *AuthService) generateFamilyID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	// Dapatkan metadata refresh token
	refreshMetadata, err := s.tokenService.GetRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token tidak valid atau kedaluwarsa")
	}

	// Dapatkan user berdasarkan ID
	user, err := s.userRepo.GetByID(refreshMetadata.UserID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	// Periksa apakah user masih aktif
	if !user.IsActive {
		return nil, errors.New("akun pengguna tidak aktif")
	}

	// Dapatkan peran dan modul user
	userRoles, err := s.userRepo.GetUserRoles(user.ID)
	if err != nil {
		userRoles = []*models.UserRole{}
	}

	// Konversi peran ke slice string untuk token
	var roles []string
	for _, userRole := range userRoles {
		// Untuk saat ini, gunakan role ID sebagai string karena tidak ada akses langsung ke nama peran
		roles = append(roles, fmt.Sprintf("role_%d", userRole.RoleID))
	}

	// Dapatkan modul user dengan unit-aware enhancement (fallback to traditional for now)
	moduleURLs, err := s.userRepo.GetUserModulesWithSubscription(user.ID)
	if err != nil {
		moduleURLs = []string{}
	}

	// Generate access token baru
	accessToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Generate refresh token baru
	newRefreshToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Set token expiration
	expiresAt := time.Now().Add(15 * time.Minute)

	// Simpan access token baru (menimpa yang ada)
	accessMetadata := token.TokenMetadata{
		UserID:    user.ID,
		UserAgent: "",         // Tidak ada user agent dalam refresh
		IP:        "",         // Tidak ada IP dalam refresh
		Abilities: moduleURLs, // Gunakan URL untuk kemampuan
		ExpiresAt: expiresAt.Unix(),
	}

	if err := s.tokenService.StoreAccessToken(accessToken, accessMetadata, 15*time.Minute); err != nil {
		return nil, err
	}

	// Simpan refresh token baru (menimpa yang ada) dengan family ID yang sama
	newRefreshMetadata := token.RefreshTokenMetadata{
		UserID:   user.ID,
		FamilyID: refreshMetadata.FamilyID,
	}

	if err := s.tokenService.StoreRefreshToken(newRefreshToken, newRefreshMetadata, 7*24*time.Hour); err != nil {
		return nil, err
	}

	// Hitung expires in dalam detik
	expiresIn := int64(15 * 60) // 15 menit dalam detik

	// Gunakan mapper untuk membuat response
	return s.authMapper.ToRefreshTokenResponse(accessToken, newRefreshToken, expiresIn), nil
}

// Logout melakukan logout menggunakan access token
func (s *AuthService) Logout(accessToken string) error {
	// Coba dapatkan metadata token untuk menemukan user ID
	metadata, err := s.tokenService.GetAccessToken(accessToken)
	if err != nil {
		// Jika token tidak ditemukan (kedaluwarsa/tidak valid), kita tidak bisa menentukan user ID
		// Dalam kasus ini, kita akan mengembalikan sukses karena token sudah hilang
		return nil
	}

	// Hapus SEMUA token untuk user ini (access dan refresh token)
	if err := s.tokenService.RevokeAllUserTokens(metadata.UserID); err != nil {
		return err
	}

	return nil
}

// LogoutByUserID melakukan logout user berdasarkan user ID (alternatif ketika token kedaluwarsa)
func (s *AuthService) LogoutByUserID(userID int64) error {
	// Hapus SEMUA token untuk user ini (access dan refresh token)
	if err := s.tokenService.RevokeAllUserTokens(userID); err != nil {
		return err
	}

	return nil
}

// CheckUserTokens memeriksa apakah user memiliki token valid di Redis
func (s *AuthService) CheckUserTokens(userID int64) (interface{}, error) {
	return s.tokenService.GetUserTokens(userID)
}

// CleanupExpiredTokens menghapus token kedaluwarsa dari Redis
func (s *AuthService) CleanupExpiredTokens() error {
	return s.tokenService.CleanupExpiredTokens()
}

// Register mendaftarkan pengguna baru
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Periksa apakah email sudah ada
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Periksa apakah user identity sudah ada (jika disediakan)
	if req.UserIdentity != nil && *req.UserIdentity != "" {
		existingUser, _ := s.userRepo.GetByUserIdentity(*req.UserIdentity)
		if existingUser != nil {
			return nil, errors.New("user identity sudah terdaftar")
		}
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	// Konversi DTO ke model menggunakan mapper
	user := s.authMapper.ToUserModel(req)
	user.PasswordHash = hashedPassword

	// Simpan user ke database
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("gagal membuat pengguna")
	}

	// Gunakan mapper untuk membuat response
	return s.authMapper.ToRegisterResponse(user, "Pengguna berhasil didaftarkan"), nil
}

// ForgotPassword mengirim email reset password
func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	// Periksa apakah user ada
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Jangan beri tahu bahwa email tidak ditemukan untuk keamanan
		return nil
	}

	// Generate reset token
	resetToken, err := s.tokenService.GenerateToken()
	if err != nil {
		return errors.New("gagal membuat token reset")
	}

	// Simpan reset token dengan expiry 1 jam
	// TODO: Implementasi penyimpanan reset token
	// TODO: Kirim email dengan reset token

	_ = user
	_ = resetToken

	return nil
}

// ResetPassword mereset password menggunakan token
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// TODO: Validasi reset token
	// TODO: Dapatkan user dari token
	// TODO: Update password user

	// Hash password baru
	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	_ = hashedPassword

	return errors.New("fitur reset password belum diimplementasi")
}

// GetUserSessionCount mengembalikan jumlah sesi aktif untuk user
func (s *AuthService) GetUserSessionCount(userID int64) (int64, error) {
	count, err := s.tokenService.GetUserSessionCount(userID)
	return int64(count), err
}

// GetUserRefreshTokenCount mengembalikan jumlah refresh token aktif untuk user
func (s *AuthService) GetUserRefreshTokenCount(userID int64) (int64, error) {
	// Dapatkan semua refresh token untuk user
	tokensResponse, err := s.tokenService.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}

	return int64(len(tokensResponse.RefreshTokens)), nil
}
