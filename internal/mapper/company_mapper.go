package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// CompanyMapper handles conversion between company models and DTOs
type CompanyMapper struct{}

// NewCompanyMapper creates a new company mapper
func NewCompanyMapper() *CompanyMapper {
	return &CompanyMapper{}
}

// ToResponse converts model to response DTO
func (m *CompanyMapper) ToResponse(company *models.Company) *dto.CompanyResponse {
	if company == nil {
		return nil
	}

	return &dto.CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Code:      company.Code,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt.Format(time.RFC3339),
		UpdatedAt: company.UpdatedAt.Format(time.RFC3339),
	}
}

// ToResponseList converts model slice to response DTO slice
func (m *CompanyMapper) ToResponseList(companies []*models.Company) []*dto.CompanyResponse {
	if companies == nil {
		return nil
	}

	responses := make([]*dto.CompanyResponse, len(companies))
	for i, company := range companies {
		responses[i] = m.ToResponse(company)
	}
	return responses
}

// ToStatsResponse converts model with stats to response DTO
func (m *CompanyMapper) ToStatsResponse(company *models.CompanyWithStats) *dto.CompanyWithStatsResponse {
	if company == nil {
		return nil
	}

	return &dto.CompanyWithStatsResponse{
		CompanyResponse: *m.ToResponse(&company.Company),
		TotalUsers:      company.TotalUsers,
		TotalBranches:   company.TotalBranches,
	}
}

// ToModel converts create request DTO to model
func (m *CompanyMapper) ToModel(req *dto.CreateCompanyRequest) *models.Company {
	if req == nil {
		return nil
	}

	return &models.Company{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: true, // Default to active
	}
}

// UpdateModel updates model with update request DTO
func (m *CompanyMapper) UpdateModel(company *models.Company, req *dto.UpdateCompanyRequest) {
	if company == nil || req == nil {
		return
	}

	if req.Name != "" {
		company.Name = req.Name
	}
	if req.Code != "" {
		company.Code = req.Code
	}
	if req.IsActive != nil {
		company.IsActive = *req.IsActive
	}
}

// ToListResponse creates paginated list response
func (m *CompanyMapper) ToListResponse(companies []*models.Company, total int64, limit, offset int) *dto.CompanyListResponse {
	return &dto.CompanyListResponse{
		Data:    m.ToResponseList(companies),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(companies)) < total,
	}
}
