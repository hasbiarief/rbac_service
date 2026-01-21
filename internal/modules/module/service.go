package module

import (
	"time"
)

type Service struct {
	repo *ModuleRepository
}

func NewService(repo *ModuleRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetModules(req *ModuleListRequest) (*ModuleListResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	modules, err := s.repo.GetAll(limit, offset, req.Search, req.Category, req.SubscriptionTier, req.ParentID, req.IsActive)
	if err != nil {
		return nil, err
	}

	total := int64(0)

	var responses []*ModuleResponse
	for _, module := range modules {
		responses = append(responses, toModuleResponse(module))
	}

	return &ModuleListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

func (s *Service) GetModulesFiltered(requestingUserID int64, req *ModuleListRequest) (*ModuleListResponse, error) {
	return s.GetModules(req)
}

func (s *Service) GetModulesNested(req *ModuleListRequest) ([]*NestedModuleResponse, error) {
	modules, err := s.repo.GetAll(0, 0, req.Search, req.Category, req.SubscriptionTier, req.ParentID, req.IsActive)
	if err != nil {
		return nil, err
	}

	return buildNestedModules(modules, nil), nil
}

func (s *Service) GetModuleByID(id int64) (*ModuleResponse, error) {
	module, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toModuleResponse(module), nil
}

func (s *Service) GetUserModules(userID int64, category string, limit int) ([]*ModuleResponse, error) {
	userModules, err := s.repo.GetUserModules(userID, 0)
	if err != nil {
		return nil, err
	}

	var responses []*ModuleResponse
	for _, um := range userModules {
		if category == "" || um.Category == category {
			// Convert UserModule to ModuleResponse
			moduleResp := &ModuleResponse{
				ID:               um.ID,
				Category:         um.Category,
				Name:             um.Name,
				URL:              um.URL,
				Icon:             um.Icon,
				Description:      um.Description,
				ParentID:         um.ParentID,
				SubscriptionTier: um.SubscriptionTier,
				IsActive:         um.IsActive,
				CreatedAt:        um.CreatedAt.Format(time.RFC3339),
				UpdatedAt:        um.UpdatedAt.Format(time.RFC3339),
			}
			responses = append(responses, moduleResp)
			if limit > 0 && len(responses) >= limit {
				break
			}
		}
	}

	return responses, nil
}

func (s *Service) GetModuleTreeByParentFiltered(userID int64, parentName string) ([]*ModuleTreeResponse, error) {
	modules, err := s.repo.GetTreeByParentName(parentName, userID)
	if err != nil {
		return nil, err
	}

	return buildModuleTree(modules, nil, 0, ""), nil
}

func (s *Service) GetModuleTreeFiltered(userID int64, category string) ([]*ModuleTreeResponse, error) {
	modules, err := s.repo.GetTreeStructure(category, userID)
	if err != nil {
		return nil, err
	}

	return buildModuleTree(modules, nil, 0, ""), nil
}

func (s *Service) GetModuleChildrenFiltered(userID int64, id int64) ([]*ModuleResponse, error) {
	modules, err := s.repo.GetChildren(id)
	if err != nil {
		return nil, err
	}

	var responses []*ModuleResponse
	for _, module := range modules {
		responses = append(responses, toModuleResponse(module))
	}

	return responses, nil
}

func (s *Service) GetModuleAncestorsFiltered(userID int64, id int64) ([]*ModuleResponse, error) {
	modules, err := s.repo.GetAncestors(id, userID)
	if err != nil {
		return nil, err
	}

	var responses []*ModuleResponse
	for _, module := range modules {
		responses = append(responses, toModuleResponse(module))
	}

	return responses, nil
}

func (s *Service) CreateModule(req *CreateModuleRequest) (*ModuleResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	module := &Module{
		Category:         req.Category,
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Description:      req.Description,
		ParentID:         req.ParentID,
		SubscriptionTier: req.SubscriptionTier,
		IsActive:         isActive,
	}

	if err := s.repo.Create(module); err != nil {
		return nil, err
	}

	return toModuleResponse(module), nil
}

func (s *Service) UpdateModule(id int64, req *UpdateModuleRequest) (*ModuleResponse, error) {
	module, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Category != "" {
		module.Category = req.Category
	}
	if req.Name != "" {
		module.Name = req.Name
	}
	if req.URL != "" {
		module.URL = req.URL
	}
	if req.Icon != "" {
		module.Icon = req.Icon
	}
	if req.Description != "" {
		module.Description = req.Description
	}
	if req.ParentID != nil {
		module.ParentID = req.ParentID
	}
	if req.SubscriptionTier != "" {
		module.SubscriptionTier = req.SubscriptionTier
	}
	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	if err := s.repo.Update(module); err != nil {
		return nil, err
	}

	return toModuleResponse(module), nil
}

func (s *Service) DeleteModule(id int64) error {
	return s.repo.Delete(id)
}

func toModuleResponse(module *Module) *ModuleResponse {
	if module == nil {
		return nil
	}

	return &ModuleResponse{
		ID:               module.ID,
		Category:         module.Category,
		Name:             module.Name,
		URL:              module.URL,
		Icon:             module.Icon,
		Description:      module.Description,
		ParentID:         module.ParentID,
		SubscriptionTier: module.SubscriptionTier,
		IsActive:         module.IsActive,
		CreatedAt:        module.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        module.UpdatedAt.Format(time.RFC3339),
	}
}

func buildNestedModules(modules []*Module, parentID *int64) []*NestedModuleResponse {
	var result []*NestedModuleResponse

	for _, module := range modules {
		if (parentID == nil && module.ParentID == nil) || (parentID != nil && module.ParentID != nil && *module.ParentID == *parentID) {
			nested := &NestedModuleResponse{
				ID:               module.ID,
				Category:         module.Category,
				Name:             module.Name,
				URL:              module.URL,
				Icon:             module.Icon,
				Description:      module.Description,
				ParentID:         module.ParentID,
				SubscriptionTier: module.SubscriptionTier,
				IsActive:         module.IsActive,
				CreatedAt:        module.CreatedAt.Format(time.RFC3339),
				UpdatedAt:        module.UpdatedAt.Format(time.RFC3339),
				Children:         buildNestedModules(modules, &module.ID),
			}
			result = append(result, nested)
		}
	}

	return result
}

func buildModuleTree(modules []*Module, parentID *int64, level int, path string) []*ModuleTreeResponse {
	var result []*ModuleTreeResponse

	for _, module := range modules {
		if (parentID == nil && module.ParentID == nil) || (parentID != nil && module.ParentID != nil && *module.ParentID == *parentID) {
			currentPath := path
			if currentPath != "" {
				currentPath += "/"
			}
			currentPath += module.URL

			tree := &ModuleTreeResponse{
				ID:               module.ID,
				Category:         module.Category,
				Name:             module.Name,
				URL:              module.URL,
				Icon:             module.Icon,
				Description:      module.Description,
				ParentID:         module.ParentID,
				SubscriptionTier: module.SubscriptionTier,
				IsActive:         module.IsActive,
				CreatedAt:        module.CreatedAt.Format(time.RFC3339),
				UpdatedAt:        module.UpdatedAt.Format(time.RFC3339),
				Level:            level,
				Path:             currentPath,
				Children:         buildModuleTree(modules, &module.ID, level+1, currentPath),
			}
			result = append(result, tree)
		}
	}

	return result
}
