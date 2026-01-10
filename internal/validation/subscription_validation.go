package validation

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/middleware"
)

// Subscription validation rules
var CreateSubscriptionValidation = middleware.ValidationRules{
	Body: &dto.CreateSubscriptionBasicRequest{},
}

var UpdateSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateSubscriptionRequest{},
}

var CreateSubscriptionPlanValidation = middleware.ValidationRules{
	Body: &dto.CreateSubscriptionPlanRequest{},
}

var UpdateSubscriptionPlanValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.UpdateSubscriptionPlanRequest{},
}

var RenewSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body:   &dto.RenewSubscriptionRequest{},
}

var CancelSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Reason            string `json:"reason"`
		CancelImmediately bool   `json:"cancel_immediately"`
	}{},
}
