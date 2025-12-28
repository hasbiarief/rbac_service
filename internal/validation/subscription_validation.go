package validation

import "gin-scalable-api/middleware"

// Subscription validation rules
var CreateSubscriptionValidation = middleware.ValidationRules{
	Body: &struct {
		CompanyID    int64  `json:"company_id" validate:"required,min=1"`
		PlanID       int64  `json:"plan_id" validate:"required,min=1"`
		BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	}{},
}

var UpdateSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		PlanID    *int64 `json:"plan_id"`
		AutoRenew *bool  `json:"auto_renew"`
	}{},
}

var RenewSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
		PlanID       *int64 `json:"plan_id"`
	}{},
}

var CancelSubscriptionValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Reason            string `json:"reason"`
		CancelImmediately bool   `json:"cancel_immediately"`
	}{},
}
