package validation

import "gin-scalable-api/middleware"

// Auth validation rules
var LoginValidation = middleware.ValidationRules{
	Body: &struct {
		UserIdentity string `json:"user_identity" validate:"required"`
		Password     string `json:"password" validate:"required,min=6"`
	}{},
}

var LoginEmailValidation = middleware.ValidationRules{
	Body: &struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}{},
}

var RefreshValidation = middleware.ValidationRules{
	Body: &struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}{},
}

var LogoutValidation = middleware.ValidationRules{
	Body: &struct {
		Token  string `json:"token"`
		UserID int64  `json:"user_id"`
	}{},
}
