package unit

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateUnitRequest validates create unit request
func ValidateCreateUnitRequest(req *CreateUnitRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateUnitRequest validates update unit request
func ValidateUpdateUnitRequest(req *UpdateUnitRequest) error {
	return validate.Struct(req)
}

// ValidateUnitListRequest validates unit list request
func ValidateUnitListRequest(req *UnitListRequest) error {
	return validate.Struct(req)
}

// ValidateBulkUpdateUnitRoleModulesRequest validates bulk update unit role modules request
func ValidateBulkUpdateUnitRoleModulesRequest(req *BulkUpdateUnitRoleModulesRequest) error {
	return validate.Struct(req)
}

// ValidateCopyUnitPermissionsRequest validates copy unit permissions request
func ValidateCopyUnitPermissionsRequest(req *CopyUnitPermissionsRequest) error {
	return validate.Struct(req)
}
