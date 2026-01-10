package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// BranchMapper menangani konversi antara model branch dan DTO
type BranchMapper struct{}

// NewBranchMapper membuat mapper cabang baru
func NewBranchMapper() *BranchMapper {
	return &BranchMapper{}
}

// ToResponse mengkonversi model ke DTO respons
func (m *BranchMapper) ToResponse(branch *models.Branch) *dto.BranchResponse {
	if branch == nil {
		return nil
	}

	return &dto.BranchResponse{
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

// ToResponseList mengkonversi slice model ke slice DTO respons
func (m *BranchMapper) ToResponseList(branches []*models.Branch) []*dto.BranchResponse {
	if branches == nil {
		return nil
	}

	responses := make([]*dto.BranchResponse, len(branches))
	for i, branch := range branches {
		responses[i] = m.ToResponse(branch)
	}
	return responses
}

// ToNestedResponse mengkonversi model ke DTO respons bersarang
func (m *BranchMapper) ToNestedResponse(branch *models.Branch) *dto.NestedBranchResponse {
	if branch == nil {
		return nil
	}

	return &dto.NestedBranchResponse{
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
		Children:  []*dto.NestedBranchResponse{},
	}
}

// ToModel mengkonversi DTO permintaan buat ke model
func (m *BranchMapper) ToModel(req *dto.CreateBranchRequest) *models.Branch {
	if req == nil {
		return nil
	}

	return &models.Branch{
		CompanyID: req.CompanyID,
		Name:      req.Name,
		Code:      req.Code,
		ParentID:  req.ParentID,
		IsActive:  true, // Default aktif
	}
}

// UpdateModel memperbarui model dengan DTO permintaan update
func (m *BranchMapper) UpdateModel(branch *models.Branch, req *dto.UpdateBranchRequest) {
	if branch == nil || req == nil {
		return
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
}

// ToListResponse membuat respons daftar dengan paginasi
func (m *BranchMapper) ToListResponse(branches []*models.Branch, total int64, limit, offset int) *dto.BranchListResponse {
	return &dto.BranchListResponse{
		Data:    m.ToResponseList(branches),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(branches)) < total,
	}
}

// ToNestedResponseList mengkonversi slice model ke struktur hierarki bersarang
func (m *BranchMapper) ToNestedResponseList(branches []*models.Branch) []*dto.NestedBranchResponse {
	if branches == nil {
		return nil
	}

	// Convert to map for easier lookup
	branchMap := make(map[int64]*dto.NestedBranchResponse)
	var rootBranches []*dto.NestedBranchResponse

	// First pass: create all branch nodes
	for _, branch := range branches {
		nestedBranch := m.ToNestedResponse(branch)
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

	return rootBranches
}
