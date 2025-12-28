package models

import "time"

type Module struct {
	ID               int64     `json:"id" db:"id"`
	Category         string    `json:"category" db:"category"`
	Name             string    `json:"name" db:"name"`
	URL              string    `json:"url" db:"url"`
	Icon             string    `json:"icon" db:"icon"`
	Description      string    `json:"description" db:"description"`
	ParentID         *int64    `json:"parent_id" db:"parent_id"`
	SubscriptionTier string    `json:"subscription_tier" db:"subscription_tier"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (Module) TableName() string {
	return "modules"
}

type ModuleWithChildren struct {
	Module
	Children []Module `json:"children,omitempty"`
}

type UserModule struct {
	Module
	CanRead   bool `json:"can_read" db:"can_read"`
	CanWrite  bool `json:"can_write" db:"can_write"`
	CanDelete bool `json:"can_delete" db:"can_delete"`
}
