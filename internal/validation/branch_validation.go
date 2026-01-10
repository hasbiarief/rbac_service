package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Branch validation rules
var CreateBranchValidation = middleware.ValidationRules{
	Body: &dto.CreateBranchRequest{},
}

var UpdateBranchValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateBranchRequest{},
}
