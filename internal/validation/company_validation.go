package validation

import "gin-scalable-api/middleware"

// Company validation rules
var CreateCompanyValidation = middleware.ValidationRules{
	Body: &struct {
		Name string `json:"name" validate:"required,min=2,max=100"`
		Code string `json:"code" validate:"required,min=2,max=20"`
	}{},
}

var UpdateCompanyValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}{},
}
