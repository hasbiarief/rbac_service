package application

import "github.com/go-playground/validator/v10"

// CreateApplicationRequest DTO
type CreateApplicationRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateApplicationRequest DTO
type UpdateApplicationRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Code        string `json:"code" validate:"omitempty,min=2,max=50"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   *int   `json:"sort_order"`
}

// PlanApplicationRequest DTO
type PlanApplicationRequest struct {
	ApplicationIDs []int64 `json:"application_ids" validate:"required,min=1"`
	IsIncluded     bool    `json:"is_included"`
}

// ApplicationListRequest DTO
type ApplicationListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// ApplicationResponse DTO
type ApplicationResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	IsActive    bool   `json:"is_active"`
	SortOrder   int    `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ApplicationWithModulesResponse DTO
type ApplicationWithModulesResponse struct {
	ApplicationResponse
	Modules map[string][][]string `json:"modules"`
}

// PlanApplicationResponse DTO
type PlanApplicationResponse struct {
	ID              int64  `json:"id"`
	PlanID          int64  `json:"plan_id"`
	PlanName        string `json:"plan_name"`
	ApplicationID   int64  `json:"application_id"`
	ApplicationName string `json:"application_name"`
	ApplicationCode string `json:"application_code"`
	IsIncluded      bool   `json:"is_included"`
	CreatedAt       string `json:"created_at"`
}

// ApplicationListResponse DTO
type ApplicationListResponse struct {
	Data    []*ApplicationResponse `json:"data"`
	Total   int64                  `json:"total"`
	Limit   int                    `json:"limit"`
	Offset  int                    `json:"offset"`
	HasMore bool                   `json:"has_more"`
}

// Validation functions
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCreateApplicationRequest validates create application request
func ValidateCreateApplicationRequest(req *CreateApplicationRequest) error {
	return validate.Struct(req)
}

// ValidateUpdateApplicationRequest validates update application request
func ValidateUpdateApplicationRequest(req *UpdateApplicationRequest) error {
	return validate.Struct(req)
}

// ValidatePlanApplicationRequest validates plan application request
func ValidatePlanApplicationRequest(req *PlanApplicationRequest) error {
	return validate.Struct(req)
}

// ValidateApplicationListRequest validates application list request
func ValidateApplicationListRequest(req *ApplicationListRequest) error {
	return validate.Struct(req)
}
