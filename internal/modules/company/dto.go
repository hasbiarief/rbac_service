package company

type CreateCompanyRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Code string `json:"code" validate:"required,min=2,max=20"`
}

type UpdateCompanyRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

type CompanyListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

type CompanyResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CompanyWithStatsResponse struct {
	CompanyResponse
	TotalUsers    int `json:"total_users"`
	TotalBranches int `json:"total_branches"`
}

type CompanyListResponse struct {
	Data    []*CompanyResponse `json:"data"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
	HasMore bool               `json:"has_more"`
}
