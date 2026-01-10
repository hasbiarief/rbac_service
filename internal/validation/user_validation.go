package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

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
	Body: &dto.CreateUserRequest{},
}

var UpdateUserValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateUserRequest{},
}

var ChangePasswordValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.ChangePasswordRequest{},
}

// Access check validation
var AccessCheckValidation = middleware.ValidationRules{
	Body: &dto.AccessCheckRequest{},
}
