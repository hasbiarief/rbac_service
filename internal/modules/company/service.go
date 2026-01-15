package company

import (
	"time"
)

type Service struct {
	repo *CompanyRepository
}

func NewService(repo *CompanyRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCompanies(req *CompanyListRequest) (*CompanyListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	companies, err := s.repo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var responses []*CompanyResponse
	for _, company := range companies {
		responses = append(responses, toCompanyResponse(company))
	}

	return &CompanyListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) GetCompanyByID(id int64) (*CompanyResponse, error) {
	company, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toCompanyResponse(company), nil
}

func (s *Service) GetCompanyWithStats(id int64) (*CompanyWithStatsResponse, error) {
	companyWithStats, err := s.repo.GetWithStats(id)
	if err != nil {
		return nil, err
	}

	return &CompanyWithStatsResponse{
		CompanyResponse: *toCompanyResponse(&companyWithStats.Company),
		TotalUsers:      companyWithStats.TotalUsers,
		TotalBranches:   companyWithStats.TotalBranches,
	}, nil
}

func (s *Service) CreateCompany(req *CreateCompanyRequest) (*CompanyResponse, error) {
	company := &Company{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: true,
	}

	if err := s.repo.Create(company); err != nil {
		return nil, err
	}

	return toCompanyResponse(company), nil
}

func (s *Service) UpdateCompany(id int64, req *UpdateCompanyRequest) (*CompanyResponse, error) {
	company, err := s.repo.GetByID(id)
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

	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	return toCompanyResponse(company), nil
}

func (s *Service) DeleteCompany(id int64) error {
	return s.repo.Delete(id)
}

func toCompanyResponse(company *Company) *CompanyResponse {
	if company == nil {
		return nil
	}

	return &CompanyResponse{
		ID:        company.ID,
		Name:      company.Name,
		Code:      company.Code,
		IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt.Format(time.RFC3339),
		UpdatedAt: company.UpdatedAt.Format(time.RFC3339),
	}
}
