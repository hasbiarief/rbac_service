package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Company validation rules
var CreateCompanyValidation = middleware.ValidationRules{
	Body: &dto.CreateCompanyRequest{},
}

var UpdateCompanyValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateCompanyRequest{},
}
