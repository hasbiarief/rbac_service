package service

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
)

type CompanyService struct {
	companyRepo   interfaces.CompanyRepositoryInterface
	companyMapper *mapper.CompanyMapper
}

func NewCompanyService(companyRepo interfaces.CompanyRepositoryInterface) *CompanyService {
	return &CompanyService{
		companyRepo:   companyRepo,
		companyMapper: mapper.NewCompanyMapper(),
	}
}

func (s *CompanyService) GetCompanies(req *dto.CompanyListRequest) (*dto.CompanyListResponse, error) {
	// Set default values jika tidak disediakan
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	companies, err := s.companyRepo.GetAll(limit, offset, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Dapatkan total count untuk pagination
	total, err := s.companyRepo.Count(req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Konversi ke DTO menggunakan mapper
	var companyResponses []*dto.CompanyResponse
	for _, company := range companies {
		companyResponses = append(companyResponses, s.companyMapper.ToResponse(company))
	}

	return &dto.CompanyListResponse{
		Data:    companyResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(companyResponses)) < total,
	}, nil
}

func (s *CompanyService) GetCompanyByID(id int64) (*dto.CompanyResponse, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.companyMapper.ToResponse(company), nil
}

func (s *CompanyService) GetCompanyWithStats(id int64) (*dto.CompanyWithStatsResponse, error) {
	companyWithStats, err := s.companyRepo.GetWithStats(id)
	if err != nil {
		return nil, err
	}

	return s.companyMapper.ToStatsResponse(companyWithStats), nil
}

func (s *CompanyService) CreateCompany(req *dto.CreateCompanyRequest) (*dto.CompanyResponse, error) {
	// Konversi DTO ke model menggunakan mapper
	company := s.companyMapper.ToModel(req)

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return s.companyMapper.ToResponse(company), nil
}

func (s *CompanyService) UpdateCompany(id int64, req *dto.UpdateCompanyRequest) (*dto.CompanyResponse, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields menggunakan mapper
	s.companyMapper.UpdateModel(company, req)

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}

	return s.companyMapper.ToResponse(company), nil
}

func (s *CompanyService) DeleteCompany(id int64) error {
	return s.companyRepo.Delete(id)
}
