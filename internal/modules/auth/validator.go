package auth

import (
	"github.com/go-playground/validator/v10"
)

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
