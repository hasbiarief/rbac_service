package role

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateRoleRequest validates create role request
func ValidateCreateRoleRequest(req *CreateRoleRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateRoleRequest validates update role request
func ValidateUpdateRoleRequest(req *UpdateRoleRequest) error {
	return validate.Struct(req)
}

// ValidateAssignRoleRequest validates assign role request
func ValidateAssignRoleRequest(req *AssignRoleRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateRolePermissionsRequest validates update role permissions request
func ValidateUpdateRolePermissionsRequest(req *UpdateRolePermissionsRequest) error {
	return validate.Struct(req)
}

// ValidateRoleListRequest validates role list request
func ValidateRoleListRequest(req *RoleListRequest) error {
	return validate.Struct(req)
}
