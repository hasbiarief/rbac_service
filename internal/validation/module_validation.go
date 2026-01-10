package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Module validation rules
var ModuleListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{
		{Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
		{Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
		{Name: "search", Type: "string"},
		{Name: "category", Type: "string"},
	},
}

var CreateModuleValidation = middleware.ValidationRules{
	Body: &dto.CreateModuleRequest{},
}

var UpdateModuleValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateModuleRequest{},
}
