package module

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateModuleRequest validates create module request
func ValidateCreateModuleRequest(req *CreateModuleRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateModuleRequest validates update module request
func ValidateUpdateModuleRequest(req *UpdateModuleRequest) error {
	return validate.Struct(req)
}

// ValidateModuleListRequest validates module list request
func ValidateModuleListRequest(req *ModuleListRequest) error {
	return validate.Struct(req)
}
