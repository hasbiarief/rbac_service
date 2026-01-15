package subscription

import "time"

type SubscriptionPlan struct {
	ID           int64                  `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name"`
	DisplayName  string                 `json:"display_name" db:"display_name"`
	Description  string                 `json:"description" db:"description"`
	PriceMonthly float64                `json:"price_monthly" db:"price_monthly"`
	PriceYearly  float64                `json:"price_yearly" db:"price_yearly"`
	MaxUsers     *int                   `json:"max_users" db:"max_users"`
	MaxBranches  *int                   `json:"max_branches" db:"max_branches"`
	Features     map[string]interface{} `json:"features" db:"features"`
	IsActive     bool                   `json:"is_active" db:"is_active"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

type Subscription struct {
	ID              int64      `json:"id" db:"id"`
	CompanyID       int64      `json:"company_id" db:"company_id"`
	PlanID          int64      `json:"plan_id" db:"plan_id"`
	Status          string     `json:"status" db:"status"`
	BillingCycle    string     `json:"billing_cycle" db:"billing_cycle"`
	StartDate       time.Time  `json:"start_date" db:"start_date"`
	EndDate         time.Time  `json:"end_date" db:"end_date"`
	Price           float64    `json:"price" db:"price"`
	Currency        string     `json:"currency" db:"currency"`
	PaymentStatus   string     `json:"payment_status" db:"payment_status"`
	LastPaymentDate *time.Time `json:"last_payment_date" db:"last_payment_date"`
	NextPaymentDate *time.Time `json:"next_payment_date" db:"next_payment_date"`
	AutoRenew       bool       `json:"auto_renew" db:"auto_renew"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	CompanyName     string     `json:"company_name,omitempty" db:"company_name"`
	PlanDisplayName string     `json:"plan_display_name,omitempty" db:"plan_display_name"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
