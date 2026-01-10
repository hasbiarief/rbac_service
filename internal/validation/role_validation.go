package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Aturan validasi peran
var CreateRoleValidation = middleware.ValidationRules{
	Body: &dto.CreateRoleRequest{},
}

var UpdateRoleValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateRoleRequest{},
}

var AssignUserRoleValidation = middleware.ValidationRules{
	Body: &dto.AssignRoleRequest{},
}

var UpdateRoleModulesValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
	},
	Body: &dto.UpdateRolePermissionsRequest{},
}

var RoleUserValidation = middleware.ValidationRules{
	Params: []middleware.ParamValidation{
		{Name: "userId", Type: "int", Required: true, Min: IntPtr(1)},
		{Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
	},
}
