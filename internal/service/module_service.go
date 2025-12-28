package service

import (
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
)

type ModuleService struct {
	moduleRepo *repository.ModuleRepository
}

func NewModuleService(moduleRepo *repository.ModuleRepository) *ModuleService {
	return &ModuleService{
		moduleRepo: moduleRepo,
	}
}

type ModuleListRequest struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Category string `form:"category"`
	Search   string `form:"search"`
	IsActive *bool  `form:"is_active"`
}

type ModuleResponse struct {
	ID               int64  `json:"id"`
	Category         string `json:"category"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier"`
	IsActive         bool   `json:"is_active"`
}

type CreateModuleRequest struct {
	Category         string `json:"category" binding:"required"`
	Name             string `json:"name" binding:"required"`
	URL              string `json:"url" binding:"required"`
	Icon             string `json:"icon"`
	Description      string `json:"description"`
	ParentID         *int64 `json:"parent_id"`
	SubscriptionTier string `json:"subscription_tier"`
	IsActive         bool   `json:"is_active"`
}

type UpdateModuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type CheckAccessRequest struct {
	UserIdentity string `json:"user_identity" binding:"required"`
	ModuleURL    string `json:"module_url" binding:"required"`
}

func (s *ModuleService) GetModules(req *ModuleListRequest) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.GetAll(req.Limit, req.Offset, req.Category, req.Search, req.IsActive)
	if err != nil {
		return nil, err
	}

	var response []*ModuleResponse
	for _, module := range modules {
		response = append(response, &ModuleResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
		})
	}

	return response, nil
}

func (s *ModuleService) GetModuleByID(id int64) (*ModuleResponse, error) {
	module, err := s.moduleRepo.GetByID(id)
	if err != nil {
		return nil, err
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
	}, nil
}

func (s *ModuleService) CreateModule(req *CreateModuleRequest) (*ModuleResponse, error) {
	module := &models.Module{
		Category:         req.Category,
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Description:      req.Description,
		ParentID:         req.ParentID,
		SubscriptionTier: req.SubscriptionTier,
		IsActive:         req.IsActive,
	}

	if err := s.moduleRepo.Create(module); err != nil {
		return nil, err
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
	}, nil
}

func (s *ModuleService) UpdateModule(id int64, req *UpdateModuleRequest) (*ModuleResponse, error) {
	module, err := s.moduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		module.Name = req.Name
	}
	if req.Description != "" {
		module.Description = req.Description
	}
	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	if err := s.moduleRepo.Update(module); err != nil {
		return nil, err
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
	}, nil
}

func (s *ModuleService) DeleteModule(id int64) error {
	return s.moduleRepo.Delete(id)
}

func (s *ModuleService) GetModuleTree(category string) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.GetTree(category)
	if err != nil {
		return nil, err
	}

	var response []*ModuleResponse
	for _, module := range modules {
		response = append(response, &ModuleResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
		})
	}

	return response, nil
}

func (s *ModuleService) GetModuleChildren(parentID int64) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	var response []*ModuleResponse
	for _, module := range modules {
		response = append(response, &ModuleResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
		})
	}

	return response, nil
}

func (s *ModuleService) GetModuleAncestors(moduleID int64) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.GetAncestors(moduleID)
	if err != nil {
		return nil, err
	}

	var response []*ModuleResponse
	for _, module := range modules {
		response = append(response, &ModuleResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
		})
	}

	return response, nil
}

func (s *ModuleService) GetUserModules(userID int64, category string, limit int) ([]*ModuleResponse, error) {
	modules, err := s.moduleRepo.GetUserModules(userID, category, limit)
	if err != nil {
		return nil, err
	}

	var response []*ModuleResponse
	for _, module := range modules {
		response = append(response, &ModuleResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
		})
	}

	return response, nil
}

func (s *ModuleService) CheckUserAccess(req *CheckAccessRequest) (bool, error) {
	return s.moduleRepo.CheckUserAccess(req.UserIdentity, req.ModuleURL)
}
