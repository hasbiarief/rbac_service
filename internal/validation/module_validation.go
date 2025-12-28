package validation

import "gin-scalable-api/middleware"

// Module validation rules
var ModuleListValidation = middleware.ValidationRules{
	Query: []middleware.QueryValidation{
		{Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
		{Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
		{Name: "search", Type: "string"},
		{Name: "category", Type: "string"},
	},
}

var CreateModuleValidation = middleware.ValidationRules{
	Body: &struct {
		Category         string `json:"category" validate:"required,min=2,max=50"`
		Name             string `json:"name" validate:"required,min=2,max=100"`
		URL              string `json:"url" validate:"required,min=1,max=255"`
		Icon             string `json:"icon" validate:"max=50"`
		Description      string `json:"description" validate:"max=500"`
		ParentID         *int64 `json:"parent_id"`
		SubscriptionTier string `json:"subscription_tier" validate:"required,oneof=basic pro enterprise"`
	}{},
}

var UpdateModuleValidation = middleware.ValidationRules{
	Params: IDValidation.Params,
	Body: &struct {
		Category         string `json:"category"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		Icon             string `json:"icon"`
		Description      string `json:"description"`
		ParentID         *int64 `json:"parent_id"`
		SubscriptionTier string `json:"subscription_tier"`
	}{},
}
