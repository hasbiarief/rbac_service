package user

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateUserRequest validates create user request
func ValidateCreateUserRequest(req *CreateUserRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateUserRequest validates update user request
func ValidateUpdateUserRequest(req *UpdateUserRequest) error {
	return validate.Struct(req)
}

// ValidateChangePasswordRequest validates change password request
func ValidateChangePasswordRequest(req *ChangePasswordRequest) error {
	return validate.Struct(req)
}

// ValidateAccessCheckRequest validates access check request
func ValidateAccessCheckRequest(req *AccessCheckRequest) error {
	return validate.Struct(req)
}

// ValidateUserListRequest validates user list request
func ValidateUserListRequest(req *UserListRequest) error {
	return validate.Struct(req)
}
