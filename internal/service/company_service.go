package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
	}
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

type CreateCompanyRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateCompanyRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

func (s *CompanyService) GetCompanies(req *CompanyListRequest) ([]*CompanyResponse, error) {
	companies, err := s.companyRepo.GetAll(req.Limit, req.Offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var response []*CompanyResponse
	for _, company := range companies {
		response = append(response, &CompanyResponse{
			ID:        company.ID,
			Name:      company.Name,
			Code:      company.Code,
			IsActive:  company.IsActive,
			CreatedAt: company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return response, nil
}

func (s *CompanyService) GetCompanyByID(id int64) (*CompanyResponse, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Code:      company.Code,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *CompanyService) CreateCompany(req *CreateCompanyRequest) (*CompanyResponse, error) {
	company := &models.Company{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: true,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return &CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Code:      company.Code,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *CompanyService) UpdateCompany(id int64, req *UpdateCompanyRequest) (*CompanyResponse, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
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

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}

	return &CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Code:      company.Code,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: company.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *CompanyService) DeleteCompany(id int64) error {
	return s.companyRepo.Delete(id)
}
