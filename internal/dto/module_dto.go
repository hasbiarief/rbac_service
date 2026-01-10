package dto

// Module Request DTO
type CreateModuleRequest struct {
	Category         string `json:"category" validate:"required"`
	Name             string `json:"name" validate:"required,min=2,max=100"`
	URL              string `json:"url" validate:"required"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier" validate:"required"`
}

// Update Module Request DTO
type UpdateModuleRequest struct {
	Category         string `json:"category"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier"`
	IsActive         *bool  `json:"is_active"`
}

// Module List Request DTO
type ModuleListRequest struct {
	Limit            int    `form:"limit"`
	Offset           int    `form:"offset"`
	Search           string `form:"search"`
	Category         string `form:"category"`
	SubscriptionTier string `form:"subscription_tier"`
	ParentID         *int64 `form:"parent_id"`
	IsActive         *bool  `form:"is_active"`
	UserID           *int64 `form:"user_id"` // For user-specific modules
	CompanyID        *int64 `form:"company_id"`
}

// Module Response DTO
type ModuleResponse struct {
	ID               int64  `json:"id"`
	Category         string `json:"category"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier"`
	IsActive         bool   `json:"is_active"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Nested Module Response DTO
type NestedModuleResponse struct {
	ID               int64                   `json:"id"`
	Category         string                  `json:"category"`
	Name             string                  `json:"name"`
	URL              string                  `json:"url"`
	Icon             string                  `json:"icon"`
	Description      string                  `json:"description"`
	ParentID         *int64                  `json:"parent_id"`
	SubscriptionTier string                  `json:"subscription_tier"`
	IsActive         bool                    `json:"is_active"`
	CreatedAt        string                  `json:"created_at"`
	UpdatedAt        string                  `json:"updated_at"`
	Children         []*NestedModuleResponse `json:"children"`
}

// User Module Response DTO
type UserModuleResponse struct {
	ModuleResponse
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanDelete bool `json:"can_delete"`
}

// Module List Response DTO
type ModuleListResponse struct {
	Data    []*ModuleResponse `json:"data"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}

// Module Tree Response DTO
type ModuleTreeResponse struct {
	ID               int64                 `json:"id"`
	Category         string                `json:"category"`
	Name             string                `json:"name"`
	URL              string                `json:"url"`
	Icon             string                `json:"icon"`
	Description      string                `json:"description"`
	ParentID         *int64                `json:"parent_id"`
	SubscriptionTier string                `json:"subscription_tier"`
	IsActive         bool                  `json:"is_active"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
	Children         []*ModuleTreeResponse `json:"children"`
	Level            int                   `json:"level"`
	Path             string                `json:"path"`
}
