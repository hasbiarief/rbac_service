package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Unit validation rules
var CreateUnitValidation = middleware.ValidationRules{
	Body: &dto.CreateUnitRequest{},
}

var UpdateUnitValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateUnitRequest{},
}

var UnitListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{},
}

// Unit Role validation rules
var CreateUnitRoleValidation = middleware.ValidationRules{
	Body: &dto.CreateUnitRoleRequest{},
}

var UnitRoleListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{},
}

// Unit Role Module validation rules
var CreateUnitRoleModuleValidation = middleware.ValidationRules{
	Body: &dto.CreateUnitRoleModuleRequest{},
}

var UpdateUnitRoleModuleValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateUnitRoleModuleRequest{},
}

var BulkUpdateUnitRoleModulesValidation = middleware.ValidationRules{
	Body: &dto.BulkUpdateUnitRoleModulesRequest{},
}

var CopyUnitPermissionsValidation = middleware.ValidationRules{
	Body: &dto.CopyUnitPermissionsRequest{},
}

// Enhanced User Role validation with unit support
var CreateUserRoleWithUnitValidation = middleware.ValidationRules{
	Body: &dto.CreateUserRoleRequest{},
}

var UpdateUserRoleWithUnitValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateUserRoleRequest{},
}
