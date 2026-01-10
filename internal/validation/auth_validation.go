package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Aturan validasi autentikasi
var LoginValidation = middleware.ValidationRules{
	Body: &dto.LoginRequest{},
}

var LoginEmailValidation = middleware.ValidationRules{
	Body: &dto.LoginRequest{}, // Sama dengan LoginValidation karena DTO menangani keduanya
}

var RegisterValidation = middleware.ValidationRules{
	Body: &dto.RegisterRequest{},
}

var RefreshValidation = middleware.ValidationRules{
	Body: &dto.RefreshTokenRequest{},
}

var LogoutValidation = middleware.ValidationRules{
	Body: &dto.LogoutRequest{},
}

var ForgotPasswordValidation = middleware.ValidationRules{
	Body: &dto.ForgotPasswordRequest{},
}

var ResetPasswordValidation = middleware.ValidationRules{
	Body: &dto.ResetPasswordRequest{},
}
