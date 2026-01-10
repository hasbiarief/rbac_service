package service

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
	"strings"
)

type BranchService struct {
	branchRepo   interfaces.BranchRepositoryInterface
	branchMapper *mapper.BranchMapper
}

func NewBranchService(branchRepo interfaces.BranchRepositoryInterface) *BranchService {
	return &BranchService{
		branchRepo:   branchRepo,
		branchMapper: mapper.NewBranchMapper(),
	}
}

func (s *BranchService) GetBranches(req *dto.BranchListRequest) (*dto.BranchListResponse, error) {
	branches, err := s.branchRepo.GetAll(req.Limit, req.Offset, req.Search, req.CompanyID, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Convert to DTO responses
	var branchResponses []*dto.BranchResponse
	for _, branch := range branches {
		branchResponses = append(branchResponses, s.branchMapper.ToResponse(branch))
	}

	// For now, return without total count (would need Count method in repository)
	return &dto.BranchListResponse{
		Data:    branchResponses,
		Total:   int64(len(branchResponses)),
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: false,
	}, nil
}

// GetBranchesNested returns all branches in nested hierarchy structure
func (s *BranchService) GetBranchesNested(req *dto.BranchListRequest) ([]*dto.NestedBranchResponse, error) {
	branches, err := s.branchRepo.GetAll(0, 0, req.Search, req.CompanyID, req.IsActive) // Get all for hierarchy
	if err != nil {
		return nil, err
	}

	return s.branchMapper.ToNestedResponseList(branches), nil
}

func (s *BranchService) GetBranchByID(id int64) (*dto.BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.branchMapper.ToResponse(branch), nil
}

func (s *BranchService) CreateBranch(req *dto.CreateBranchRequest) (*dto.BranchResponse, error) {
	branch := s.branchMapper.ToModel(req)
	branch.IsActive = true

	err := s.branchRepo.Create(branch)
	if err != nil {
		return nil, err
	}

	return s.branchMapper.ToResponse(branch), nil
}

func (s *BranchService) UpdateBranch(id int64, req *dto.UpdateBranchRequest) (*dto.BranchResponse, error) {
	branch, err := s.branchRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields from request
	if req.Name != "" {
		branch.Name = req.Name
	}
	if req.Code != "" {
		branch.Code = req.Code
	}
	branch.ParentID = req.ParentID
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}

	err = s.branchRepo.Update(branch)
	if err != nil {
		return nil, err
	}

	return s.branchMapper.ToResponse(branch), nil
}

func (s *BranchService) DeleteBranch(id int64) error {
	return s.branchRepo.Delete(id)
}

func (s *BranchService) GetCompanyBranches(companyID int64, includeHierarchy bool) ([]*dto.BranchResponse, error) {
	branches, err := s.branchRepo.GetByCompany(companyID, includeHierarchy)
	if err != nil {
		return nil, err
	}

	var response []*dto.BranchResponse
	for _, branch := range branches {
		response = append(response, s.branchMapper.ToResponse(branch))
	}

	return response, nil
}

// GetCompanyBranchesNested returns branches in nested hierarchy structure
func (s *BranchService) GetCompanyBranchesNested(companyID int64) ([]*dto.NestedBranchResponse, error) {
	branches, err := s.branchRepo.GetByCompany(companyID, true)
	if err != nil {
		return nil, err
	}

	return s.branchMapper.ToNestedResponseList(branches), nil
}

func (s *BranchService) GetBranchChildren(parentID int64) ([]*dto.BranchResponse, error) {
	branches, err := s.branchRepo.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	var response []*dto.BranchResponse
	for _, branch := range branches {
		response = append(response, s.branchMapper.ToResponse(branch))
	}

	return response, nil
}

// GetBranchHierarchyByID returns branch hierarchy starting from specific branch ID
func (s *BranchService) GetBranchHierarchyByID(branchID int64, nested bool) (interface{}, error) {
	// First, get the specific branch to ensure it exists
	targetBranch, err := s.branchRepo.GetByID(branchID)
	if err != nil {
		return nil, err
	}

	// Get all branches to build the hierarchy
	allBranches, err := s.branchRepo.GetByCompany(targetBranch.CompanyID, true)
	if err != nil {
		return nil, err
	}

	if !nested {
		// Return flat structure - filter branches that are descendants of target branch
		var descendants []*dto.BranchResponse
		targetPath := targetBranch.Path

		for _, branch := range allBranches {
			// Include the target branch itself and all its descendants
			if branch.Path == targetPath || strings.HasPrefix(branch.Path, targetPath+".") {
				descendants = append(descendants, s.branchMapper.ToResponse(branch))
			}
		}
		return descendants, nil
	}

	// Return nested structure - build tree starting from target branch
	// Filter branches that are descendants of target branch
	var filteredBranches []*models.Branch
	targetPath := targetBranch.Path

	for _, branch := range allBranches {
		if branch.Path == targetPath || strings.HasPrefix(branch.Path, targetPath+".") {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	nestedBranches := s.branchMapper.ToNestedResponseList(filteredBranches)

	// Find the root branch (target branch)
	for _, nestedBranch := range nestedBranches {
		if nestedBranch.ID == branchID {
			return nestedBranch, nil
		}
	}

	return nil, nil
}
