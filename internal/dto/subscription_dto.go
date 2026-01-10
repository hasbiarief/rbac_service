package dto

// Subscription Plan Request DTO
type CreateSubscriptionPlanRequest struct {
	Name         string                 `json:"name" validate:"required,min=2,max=100"`
	DisplayName  string                 `json:"display_name" validate:"required,min=2,max=100"`
	Description  string                 `json:"description"`
	PriceMonthly float64                `json:"price_monthly" validate:"min=0"`
	PriceYearly  float64                `json:"price_yearly" validate:"min=0"`
	MaxUsers     *int                   `json:"max_users"`
	MaxBranches  *int                   `json:"max_branches"`
	Features     map[string]interface{} `json:"features"`
}

// Update Subscription Plan Request DTO
type UpdateSubscriptionPlanRequest struct {
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	PriceMonthly *float64               `json:"price_monthly" validate:"omitempty,min=0"`
	PriceYearly  *float64               `json:"price_yearly" validate:"omitempty,min=0"`
	MaxUsers     *int                   `json:"max_users"`
	MaxBranches  *int                   `json:"max_branches"`
	Features     map[string]interface{} `json:"features"`
	IsActive     *bool                  `json:"is_active"`
}

// Subscription Request DTO
type CreateSubscriptionRequest struct {
	CompanyID    int64   `json:"company_id" validate:"required"`
	PlanID       int64   `json:"plan_id" validate:"required"`
	BillingCycle string  `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	StartDate    string  `json:"start_date" validate:"required"` // Format: 2006-01-02
	EndDate      string  `json:"end_date" validate:"required"`   // Format: 2006-01-02
	Price        float64 `json:"price" validate:"min=0"`
	Currency     string  `json:"currency" validate:"required,len=3"`
	AutoRenew    bool    `json:"auto_renew"`
}

// Simple Create Subscription Request DTO (for basic creation)
type CreateSubscriptionBasicRequest struct {
	CompanyID    int64  `json:"company_id" validate:"required"`
	PlanID       int64  `json:"plan_id" validate:"required"`
	BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
}

// Update Subscription Request DTO
type UpdateSubscriptionRequest struct {
	PlanID          *int64   `json:"plan_id"`
	Status          string   `json:"status" validate:"omitempty,oneof=active inactive cancelled expired"`
	BillingCycle    string   `json:"billing_cycle" validate:"omitempty,oneof=monthly yearly"`
	EndDate         string   `json:"end_date"` // Format: 2006-01-02
	Price           *float64 `json:"price" validate:"omitempty,min=0"`
	PaymentStatus   string   `json:"payment_status" validate:"omitempty,oneof=pending paid failed"`
	NextPaymentDate string   `json:"next_payment_date"` // Format: 2006-01-02
	AutoRenew       *bool    `json:"auto_renew"`
}

// Renew Subscription Request DTO
type RenewSubscriptionRequest struct {
	BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
	PlanID       *int64 `json:"plan_id"`
}

// Subscription List Request DTO
type SubscriptionListRequest struct {
	Limit         int    `form:"limit"`
	Offset        int    `form:"offset"`
	Search        string `form:"search"`
	CompanyID     *int64 `form:"company_id"`
	PlanID        *int64 `form:"plan_id"`
	Status        string `form:"status"`
	BillingCycle  string `form:"billing_cycle"`
	PaymentStatus string `form:"payment_status"`
}

// Subscription Plan Response DTO
type SubscriptionPlanResponse struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	PriceMonthly float64                `json:"price_monthly"`
	PriceYearly  float64                `json:"price_yearly"`
	MaxUsers     *int                   `json:"max_users"`
	MaxBranches  *int                   `json:"max_branches"`
	Features     map[string]interface{} `json:"features"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

// Subscription Response DTO
type SubscriptionResponse struct {
	ID              int64   `json:"id"`
	CompanyID       int64   `json:"company_id"`
	PlanID          int64   `json:"plan_id"`
	Status          string  `json:"status"`
	BillingCycle    string  `json:"billing_cycle"`
	StartDate       string  `json:"start_date"`
	EndDate         string  `json:"end_date"`
	Price           float64 `json:"price"`
	Currency        string  `json:"currency"`
	PaymentStatus   string  `json:"payment_status"`
	LastPaymentDate *string `json:"last_payment_date"`
	NextPaymentDate *string `json:"next_payment_date"`
	AutoRenew       bool    `json:"auto_renew"`
	CompanyName     string  `json:"company_name,omitempty"`
	PlanDisplayName string  `json:"plan_display_name,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// Subscription Plan List Response DTO
type SubscriptionPlanListResponse struct {
	Data    []*SubscriptionPlanResponse `json:"data"`
	Total   int64                       `json:"total"`
	Limit   int                         `json:"limit"`
	Offset  int                         `json:"offset"`
	HasMore bool                        `json:"has_more"`
}

// Subscription List Response DTO
type SubscriptionListResponse struct {
	Data    []*SubscriptionResponse `json:"data"`
	Total   int64                   `json:"total"`
	Limit   int                     `json:"limit"`
	Offset  int                     `json:"offset"`
	HasMore bool                    `json:"has_more"`
}
