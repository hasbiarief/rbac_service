package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"strings"
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

// NestedBranchResponse represents a branch with its children in nested structure
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

// GetBranchesNested returns all branches in nested hierarchy structure
func (s *BranchService) GetBranchesNested(req *BranchListRequest) ([]*NestedBranchResponse, error) {
	branches, err := s.branchRepo.GetAll(0, 0, req.Search, req.CompanyID, req.IsActive) // Get all for hierarchy
	if err != nil {
		return nil, err
	}

	// Convert to map for easier lookup
	branchMap := make(map[int64]*NestedBranchResponse)
	var rootBranches []*NestedBranchResponse

	// First pass: create all branch nodes
	for _, branch := range branches {
		nestedBranch := &NestedBranchResponse{
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
			Children:  []*NestedBranchResponse{},
		}
		branchMap[branch.ID] = nestedBranch
	}

	// Second pass: build hierarchy
	for _, branch := range branches {
		nestedBranch := branchMap[branch.ID]
		if branch.ParentID == nil {
			// Root branch
			rootBranches = append(rootBranches, nestedBranch)
		} else {
			// Child branch - add to parent's children
			if parent, exists := branchMap[*branch.ParentID]; exists {
				parent.Children = append(parent.Children, nestedBranch)
			}
		}
	}

	return rootBranches, nil
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

// GetCompanyBranchesNested returns branches in nested hierarchy structure
func (s *BranchService) GetCompanyBranchesNested(companyID int64) ([]*NestedBranchResponse, error) {
	branches, err := s.branchRepo.GetByCompany(companyID, true)
	if err != nil {
		return nil, err
	}

	// Convert to map for easier lookup
	branchMap := make(map[int64]*NestedBranchResponse)
	var rootBranches []*NestedBranchResponse

	// First pass: create all branch nodes
	for _, branch := range branches {
		nestedBranch := &NestedBranchResponse{
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
			Children:  []*NestedBranchResponse{},
		}
		branchMap[branch.ID] = nestedBranch
	}

	// Second pass: build hierarchy
	for _, branch := range branches {
		nestedBranch := branchMap[branch.ID]
		if branch.ParentID == nil {
			// Root branch
			rootBranches = append(rootBranches, nestedBranch)
		} else {
			// Child branch - add to parent's children
			if parent, exists := branchMap[*branch.ParentID]; exists {
				parent.Children = append(parent.Children, nestedBranch)
			}
		}
	}

	return rootBranches, nil
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
		var descendants []*BranchResponse
		targetPath := targetBranch.Path

		for _, branch := range allBranches {
			// Include the target branch itself and all its descendants
			if branch.Path == targetPath || strings.HasPrefix(branch.Path, targetPath+".") {
				descendants = append(descendants, &BranchResponse{
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
		}
		return descendants, nil
	}

	// Return nested structure - build tree starting from target branch
	branchMap := make(map[int64]*NestedBranchResponse)
	targetPath := targetBranch.Path

	// First pass: create branch nodes for target and its descendants
	for _, branch := range allBranches {
		if branch.Path == targetPath || strings.HasPrefix(branch.Path, targetPath+".") {
			nestedBranch := &NestedBranchResponse{
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
				Children:  []*NestedBranchResponse{},
			}
			branchMap[branch.ID] = nestedBranch
		}
	}

	// Second pass: build hierarchy relationships
	var rootBranch *NestedBranchResponse
	for _, branch := range allBranches {
		if branch.Path == targetPath || strings.HasPrefix(branch.Path, targetPath+".") {
			nestedBranch := branchMap[branch.ID]

			if branch.ID == branchID {
				// This is our root branch
				rootBranch = nestedBranch
			} else if branch.ParentID != nil {
				// Add to parent's children if parent exists in our filtered set
				if parent, exists := branchMap[*branch.ParentID]; exists {
					parent.Children = append(parent.Children, nestedBranch)
				}
			}
		}
	}

	return rootBranch, nil
}
