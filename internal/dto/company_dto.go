package dto

// Company Request DTO
type CreateCompanyRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	Code string `json:"code" validate:"required,min=2,max=20"`
}

// Update Company Request DTO
type UpdateCompanyRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

// Company List Request DTO
type CompanyListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

// Company Response DTO
type CompanyResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Company With Stats Response DTO
type CompanyWithStatsResponse struct {
	CompanyResponse
	TotalUsers    int `json:"total_users"`
	TotalBranches int `json:"total_branches"`
}

// Company List Response DTO
type CompanyListResponse struct {
	Data    []*CompanyResponse `json:"data"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
	HasMore bool               `json:"has_more"`
}
