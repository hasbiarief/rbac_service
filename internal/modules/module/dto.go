package module

type CreateModuleRequest struct {
	Category         string `json:"category" validate:"required"`
	Name             string `json:"name" validate:"required,min=2,max=100"`
	URL              string `json:"url" validate:"required"`
	Icon             string `json:"icon" validate:"omitempty"`
	Description      string `json:"description" validate:"omitempty"`
	ParentID         *int64 `json:"parent_id" validate:"omitempty"`
	SubscriptionTier string `json:"subscription_tier" validate:"required"`
	IsActive         *bool  `json:"is_active" validate:"omitempty"`
}

type UpdateModuleRequest struct {
	Category         string `json:"category" validate:"omitempty"`
	Name             string `json:"name" validate:"omitempty,min=2,max=100"`
	URL              string `json:"url" validate:"omitempty"`
	Icon             string `json:"icon" validate:"omitempty"`
	Description      string `json:"description" validate:"omitempty"`
	ParentID         *int64 `json:"parent_id" validate:"omitempty"`
	SubscriptionTier string `json:"subscription_tier" validate:"omitempty"`
	IsActive         *bool  `json:"is_active" validate:"omitempty"`
}

type ModuleListRequest struct {
	Limit            int    `form:"limit"`
	Offset           int    `form:"offset"`
	Search           string `form:"search"`
	Category         string `form:"category"`
	SubscriptionTier string `form:"subscription_tier"`
	ParentID         *int64 `form:"parent_id"`
	IsActive         *bool  `form:"is_active"`
	UserID           *int64 `form:"user_id"`
	CompanyID        *int64 `form:"company_id"`
}

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

type UserModuleResponse struct {
	ModuleResponse
	CanRead   bool `json:"can_read"`
	CanWrite  bool `json:"can_write"`
	CanDelete bool `json:"can_delete"`
}

type ModuleListResponse struct {
	Data    []*ModuleResponse `json:"data"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}

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
