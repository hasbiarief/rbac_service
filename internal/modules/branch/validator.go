package branch

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateBranchRequest validates create branch request
func ValidateCreateBranchRequest(req *CreateBranchRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateBranchRequest validates update branch request
func ValidateUpdateBranchRequest(req *UpdateBranchRequest) error {
	return validate.Struct(req)
}

// ValidateBranchListRequest validates branch list request
func ValidateBranchListRequest(req *BranchListRequest) error {
	return validate.Struct(req)
}
