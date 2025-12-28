package models

import (
	"gin-scalable-api/pkg/model"
	"time"
)

type SubscriptionPlan struct {
	model.BaseModel
	Name         string      `json:"name" db:"name"`
	DisplayName  string      `json:"display_name" db:"display_name"`
	Description  string      `json:"description" db:"description"`
	PriceMonthly float64     `json:"price_monthly" db:"price_monthly"`
	PriceYearly  float64     `json:"price_yearly" db:"price_yearly"`
	MaxUsers     *int        `json:"max_users" db:"max_users"`
	MaxBranches  *int        `json:"max_branches" db:"max_branches"`
	Features     model.JSONB `json:"features" db:"features"`
	IsActive     bool        `json:"is_active" db:"is_active"`
}

type Subscription struct {
	model.BaseModel
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

	// Additional fields for responses
	CompanyName     string `json:"company_name,omitempty" db:"-"`
	PlanDisplayName string `json:"plan_display_name,omitempty" db:"-"`
}

type PlanModule struct {
	ID         int64 `json:"id" db:"id"`
	PlanID     int64 `json:"plan_id" db:"plan_id"`
	ModuleID   int64 `json:"module_id" db:"module_id"`
	IsIncluded bool  `json:"is_included" db:"is_included"`
}
