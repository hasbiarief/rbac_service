package application

import (
	"gin-scalable-api/pkg/model"
	"time"
)

// Application represents an application in the system
type Application struct {
	model.BaseModel
	Name        string    `json:"name" db:"name"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"`
	URL         string    `json:"url" db:"url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// PlanApplication represents the relationship between subscription plans and applications
type PlanApplication struct {
	ID            int64     `json:"id" db:"id"`
	PlanID        int64     `json:"plan_id" db:"plan_id"`
	ApplicationID int64     `json:"application_id" db:"application_id"`
	IsIncluded    bool      `json:"is_included" db:"is_included"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ApplicationWithModules represents an application with its modules grouped by category
type ApplicationWithModules struct {
	Application
	Modules map[string][][]string `json:"modules"`
}

// TableName returns the table name for Application
func (Application) TableName() string {
	return "applications"
}

// TableName returns the table name for PlanApplication
func (PlanApplication) TableName() string {
	return "plan_applications"
}
