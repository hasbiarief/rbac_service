package branch

import (
	"time"
)

type Service struct {
	repo *BranchRepository
}

func NewService(repo *BranchRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetBranches(req *BranchListRequest) (*BranchListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	branches, err := s.repo.GetAll(limit, offset, req.Search, req.CompanyID, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Count not implemented in old repository, return 0 for now
	total := int64(0)

	var responses []*BranchResponse
	for _, branch := range branches {
		responses = append(responses, toBranchResponse(branch))
	}

	return &BranchListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) GetBranchesNested(req *BranchListRequest) ([]*NestedBranchResponse, error) {
	branches, err := s.repo.GetAll(0, 0, req.Search, req.CompanyID, req.IsActive)
	if err != nil {
		return nil, err
	}

	return buildNestedBranches(branches, nil), nil
}

func (s *Service) GetBranchByID(id int64) (*BranchResponse, error) {
	branch, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toBranchResponse(branch), nil
}

func (s *Service) CreateBranch(req *CreateBranchRequest) (*BranchResponse, error) {
	branch := &Branch{
		CompanyID: req.CompanyID,
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		IsActive:  true,
		Level:     1,
		Path:      "/",
	}

	if req.ParentID != nil {
		parent, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			return nil, err
		}
		branch.Level = parent.Level + 1
		branch.Path = parent.Path + "/" + parent.Code
	}

	if err := s.repo.Create(branch); err != nil {
		return nil, err
	}

	return toBranchResponse(branch), nil
}

func (s *Service) UpdateBranch(id int64, req *UpdateBranchRequest) (*BranchResponse, error) {
	branch, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		branch.Name = req.Name
	}
	if req.Code != "" {
		branch.Code = req.Code
	}
	if req.ParentID != nil {
		branch.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}

	if err := s.repo.Update(branch); err != nil {
		return nil, err
	}

	return toBranchResponse(branch), nil
}

func (s *Service) DeleteBranch(id int64) error {
	return s.repo.Delete(id)
}

func (s *Service) GetCompanyBranches(companyID int64, includeHierarchy bool) (interface{}, error) {
	branches, err := s.repo.GetAll(1000, 0, "", &companyID, nil)
	if err != nil {
		return nil, err
	}

	var responses []*BranchResponse
	for _, branch := range branches {
		responses = append(responses, toBranchResponse(branch))
	}

	return responses, nil
}

func (s *Service) GetCompanyBranchesNested(companyID int64) ([]*NestedBranchResponse, error) {
	branches, err := s.repo.GetAll(1000, 0, "", &companyID, nil)
	if err != nil {
		return nil, err
	}

	return buildNestedBranches(branches, nil), nil
}

func (s *Service) GetBranchChildren(id int64) ([]*BranchResponse, error) {
	branches, err := s.repo.GetChildren(id)
	if err != nil {
		return nil, err
	}

	var responses []*BranchResponse
	for _, branch := range branches {
		responses = append(responses, toBranchResponse(branch))
	}

	return responses, nil
}

func (s *Service) GetBranchHierarchyByID(id int64, nested bool) (interface{}, error) {
	// GetHierarchy not available in old repository, just get the branch itself
	branch, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	branches := []*Branch{branch}

	if nested {
		return buildNestedBranches(branches, nil), nil
	}

	var responses []*BranchResponse
	for _, branch := range branches {
		responses = append(responses, toBranchResponse(branch))
	}

	return responses, nil
}

// Helper functions
func toBranchResponse(branch *Branch) *BranchResponse {
	if branch == nil {
		return nil
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
	}
}

func buildNestedBranches(branches []*Branch, parentID *int64) []*NestedBranchResponse {
	var result []*NestedBranchResponse

	for _, branch := range branches {
		if (parentID == nil && branch.ParentID == nil) || (parentID != nil && branch.ParentID != nil && *branch.ParentID == *parentID) {
			nested := &NestedBranchResponse{
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
				Children:  buildNestedBranches(branches, &branch.ID),
			}
			result = append(result, nested)
		}
	}

	return result
}
