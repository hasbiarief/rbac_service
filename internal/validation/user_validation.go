package validation

import "gin-scalable-api/middleware"

// User validation rules
var UserListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{
		{Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
		{Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
		{Name: "search", Type: "string"},
		{Name: "is_active", Type: "bool"},
	},
}

var CreateUserValidation = middleware.ValidationRules{
	Body: &struct {
		Name         string  `json:"name" validate:"required,min=2,max=100"`
		Email        string  `json:"email" validate:"required,email,max=255"`
		UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
		Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
	}{},
}

var UpdateUserValidation = middleware.ValidationRules{
	Body: &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		IsActive *bool  `json:"is_active"`
	}{},
}

var AccessCheckValidation = middleware.ValidationRules{
	Body: &struct {
		UserIdentity string `json:"user_identity" validate:"required"`
		ModuleURL    string `json:"module_url" validate:"required"`
	}{},
}

var PasswordChangeValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
	},
	Body: &struct {
		CurrentPassword string `json:"current_password" validate:"required,min=6"`
		NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
	}{},
}

var MyPasswordChangeValidation = middleware.ValidationRules{
	Body: &struct {
		CurrentPassword string `json:"current_password" validate:"required,min=6"`
		NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
	}{},
}
