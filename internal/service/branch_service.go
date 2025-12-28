package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"time"
)

type BranchService struct {
	branchRepo *repository.BranchRepository
}

func NewBranchService(branchRepo *repository.BranchRepository) *BranchService {
	return &BranchService{
		branchRepo: branchRepo,
	}
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

type BranchListRequest struct {
	Limit            int    `form:"limit"`
	Offset           int    `form:"offset"`
	Search           string `form:"search"`
	CompanyID        *int64 `form:"company_id"`
	IsActive         *bool  `form:"is_active"`
	IncludeHierarchy bool   `form:"include_hierarchy"`
}

type CreateBranchRequest struct {
	CompanyID int64  `json:"company_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code" binding:"required"`
	ParentID  *int64 `json:"parent_id"`
}

type UpdateBranchRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	ParentID *int64 `json:"parent_id"`
	IsActive *bool  `json:"is_active"`
}

func (s *BranchService) GetBranches(req *BranchListRequest) ([]*BranchResponse, error) {
	branches, err := s.branchRepo.GetAll(req.Limit, req.Offset, req.Search, req.CompanyID, req.IsActive)
	if err != nil {
		return nil, err
	}

	var response []*BranchResponse
	for _, branch := range branches {
		response = append(response, &BranchResponse{
			ID:        branch.ID,
			CompanyID: branch.CompanyID,
			Name:      branch.Name,
			Code:      branch.Code,
			ParentID:  branch.ParentID,
			Level:     branch.Level,
			Path:      branch.Path,
			IsActive:  branch.IsActive,
			CreatedAt: branch.CreatedAt.Format(time.RFC3339),
			UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *BranchService) GetBranchByID(id int64) (*BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &BranchResponse{
		ID:        branch.ID,
		CompanyID: branch.CompanyID,
		Name:      branch.Name,
		Code:      branch.Code,
		ParentID:  branch.ParentID,
		Level:     branch.Level,
		Path:      branch.Path,
		IsActive:  branch.IsActive,
		CreatedAt: branch.CreatedAt.Format(time.RFC3339),
		UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *BranchService) CreateBranch(req *CreateBranchRequest) (*BranchResponse, error) {
	branch := &models.Branch{
		CompanyID: req.CompanyID,
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		IsActive:  true,
	}

	err := s.branchRepo.Create(branch)
	if err != nil {
		return nil, err
	}

	return &BranchResponse{
		ID:        branch.ID,
		CompanyID: branch.CompanyID,
		Name:      branch.Name,
		Code:      branch.Code,
		ParentID:  branch.ParentID,
		Level:     branch.Level,
		Path:      branch.Path,
		IsActive:  branch.IsActive,
		CreatedAt: branch.CreatedAt.Format(time.RFC3339),
		UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *BranchService) UpdateBranch(id int64, req *UpdateBranchRequest) (*BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	branch.Name = req.Name
	branch.Code = req.Code
	branch.ParentID = req.ParentID
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}

	err = s.branchRepo.Update(branch)
	if err != nil {
		return nil, err
	}

	return &BranchResponse{
		ID:        branch.ID,
		CompanyID: branch.CompanyID,
		Name:      branch.Name,
		Code:      branch.Code,
		ParentID:  branch.ParentID,
		Level:     branch.Level,
		Path:      branch.Path,
		IsActive:  branch.IsActive,
		CreatedAt: branch.CreatedAt.Format(time.RFC3339),
		UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *BranchService) DeleteBranch(id int64) error {
	return s.branchRepo.Delete(id)
}

func (s *BranchService) GetCompanyBranches(companyID int64, includeHierarchy bool) ([]*BranchResponse, error) {
	branches, err := s.branchRepo.GetByCompany(companyID, includeHierarchy)
	if err != nil {
		return nil, err
	}

	var response []*BranchResponse
	for _, branch := range branches {
		response = append(response, &BranchResponse{
			ID:        branch.ID,
			CompanyID: branch.CompanyID,
			Name:      branch.Name,
			Code:      branch.Code,
			ParentID:  branch.ParentID,
			Level:     branch.Level,
			Path:      branch.Path,
			IsActive:  branch.IsActive,
			CreatedAt: branch.CreatedAt.Format(time.RFC3339),
			UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *BranchService) GetBranchChildren(parentID int64) ([]*BranchResponse, error) {
	branches, err := s.branchRepo.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	var response []*BranchResponse
	for _, branch := range branches {
		response = append(response, &BranchResponse{
			ID:        branch.ID,
			CompanyID: branch.CompanyID,
			Name:      branch.Name,
			Code:      branch.Code,
			ParentID:  branch.ParentID,
			Level:     branch.Level,
			Path:      branch.Path,
			IsActive:  branch.IsActive,
			CreatedAt: branch.CreatedAt.Format(time.RFC3339),
			UpdatedAt: branch.UpdatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}
