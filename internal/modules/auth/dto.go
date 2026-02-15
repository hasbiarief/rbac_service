package auth

import (
	"github.com/go-playground/validator/v10"
)

// Auth Request DTO
type LoginRequest struct {
	UserIdentity string `json:"user_identity" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
}

// Login with Email Request DTO (specific for email login)
type LoginEmailRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Register Request DTO
type RegisterRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Email        string  `json:"email" validate:"required,email"`
	UserIdentity *string `json:"user_identity"`
	Password     string  `json:"password" validate:"required,min=6"`
}

// Refresh Token Request DTO
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Logout Request DTO
type LogoutRequest struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

// Forgot Password Request DTO
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Reset Password Request DTO
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// Profile Request DTO
type ProfileRequest struct {
	UserIdentity    string `json:"user_identity" validate:"required"`
	ApplicationCode string `json:"application_code" validate:"required"`
}

// Profile Response DTO
type ProfileResponse struct {
	User interface{} `json:"user"`
}

// Auth Response DTO
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	User         interface{} `json:"user"`
}

// Refresh Token Response DTO
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Register Response DTO
type RegisterResponse struct {
	User    interface{} `json:"user"`
	Message string      `json:"message"`
}

// Validation functions
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateLoginRequest validates login request
func ValidateLoginRequest(req *LoginRequest) error {
	return validate.Struct(req)
}

// ValidateLoginEmailRequest validates login with email request
func ValidateLoginEmailRequest(req *LoginEmailRequest) error {
	return validate.Struct(req)
}

// ValidateRegisterRequest validates register request
func ValidateRegisterRequest(req *RegisterRequest) error {
	return validate.Struct(req)
}

// ValidateRefreshTokenRequest validates refresh token request
func ValidateRefreshTokenRequest(req *RefreshTokenRequest) error {
	return validate.Struct(req)
}

// ValidateLogoutRequest validates logout request
func ValidateLogoutRequest(req *LogoutRequest) error {
	return validate.Struct(req)
}

// ValidateForgotPasswordRequest validates forgot password request
func ValidateForgotPasswordRequest(req *ForgotPasswordRequest) error {
	return validate.Struct(req)
}

// ValidateResetPasswordRequest validates reset password request
func ValidateResetPasswordRequest(req *ResetPasswordRequest) error {
	return validate.Struct(req)
}

// ValidateProfileRequest validates profile request
func ValidateProfileRequest(req *ProfileRequest) error {
	return validate.Struct(req)
}
