package validation

import "gin-scalable-api/middleware"

// Branch validation rules
var CreateBranchValidation = middleware.ValidationRules{
	Body: &struct {
		CompanyID int64  `json:"company_id" validate:"required,min=1"`
		Name      string `json:"name" validate:"required,min=2,max=100"`
		Code      string `json:"code" validate:"required,min=2,max=20"`
		ParentID  *int64 `json:"parent_id"`
	}{},
}

var UpdateBranchValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Name     string `json:"name"`
		Code     string `json:"code"`
		ParentID *int64 `json:"parent_id"`
		IsActive *bool  `json:"is_active"`
	}{},
}
