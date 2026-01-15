package branch

type CreateBranchRequest struct {
	CompanyID int64  `json:"company_id" validate:"required,min=1"`
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Code      string `json:"code" validate:"required,min=2,max=20"`
	ParentID  *int64 `json:"parent_id"`
}

type UpdateBranchRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	ParentID *int64 `json:"parent_id"`
	IsActive *bool  `json:"is_active"`
}

type BranchListRequest struct {
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
	Search    string `form:"search"`
	CompanyID *int64 `form:"company_id"`
	IsActive  *bool  `form:"is_active"`
}

type BranchResponse struct {
	ID        int64  `json:"id"`
	CompanyID int64  `json:"company_id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	ParentID  *int64 `json:"parent_id"`
	Level     int    `json:"level"`
	Path      string `json:"path"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type NestedBranchResponse struct {
	ID        int64                   `json:"id"`
	CompanyID int64                   `json:"company_id"`
	Name      string                  `json:"name"`
	Code      string                  `json:"code"`
	ParentID  *int64                  `json:"parent_id"`
	Level     int                     `json:"level"`
	Path      string                  `json:"path"`
	IsActive  bool                    `json:"is_active"`
	CreatedAt string                  `json:"created_at"`
	UpdatedAt string                  `json:"updated_at"`
	Children  []*NestedBranchResponse `json:"children"`
}

type BranchListResponse struct {
	Data    []*BranchResponse `json:"data"`
	Total   int64             `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}
