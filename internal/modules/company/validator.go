package company

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateCompanyRequest validates create company request
func ValidateCreateCompanyRequest(req *CreateCompanyRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateCompanyRequest validates update company request
func ValidateUpdateCompanyRequest(req *UpdateCompanyRequest) error {
	return validate.Struct(req)
}

// ValidateCompanyListRequest validates company list request
func ValidateCompanyListRequest(req *CompanyListRequest) error {
	return validate.Struct(req)
}
