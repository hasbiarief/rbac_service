package subscription

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateSubscriptionPlanRequest validates create subscription plan request
func ValidateCreateSubscriptionPlanRequest(req *CreateSubscriptionPlanRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateSubscriptionPlanRequest validates update subscription plan request
func ValidateUpdateSubscriptionPlanRequest(req *UpdateSubscriptionPlanRequest) error {
	return validate.Struct(req)
}

// ValidateCreateSubscriptionRequest validates create subscription request
func ValidateCreateSubscriptionRequest(req *CreateSubscriptionRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateSubscriptionRequest validates update subscription request
func ValidateUpdateSubscriptionRequest(req *UpdateSubscriptionRequest) error {
	return validate.Struct(req)
}

// ValidateSubscriptionListRequest validates subscription list request
func ValidateSubscriptionListRequest(req *SubscriptionListRequest) error {
	return validate.Struct(req)
}
