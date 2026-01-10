package mapper

import (
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/models"
	"time"
)

// ModuleMapper handles conversion between module models and DTOs
type ModuleMapper struct{}

// NewModuleMapper creates a new module mapper
func NewModuleMapper() *ModuleMapper {
	return &ModuleMapper{}
}

// ToResponse converts model to response DTO
func (m *ModuleMapper) ToResponse(module *models.Module) *dto.ModuleResponse {
	if module == nil {
		return nil
	}

	return &dto.ModuleResponse{
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

// ToResponseList converts model slice to response DTO slice
func (m *ModuleMapper) ToResponseList(modules []*models.Module) []*dto.ModuleResponse {
	if modules == nil {
		return nil
	}

	responses := make([]*dto.ModuleResponse, len(modules))
	for i, module := range modules {
		responses[i] = m.ToResponse(module)
	}
	return responses
}

// ToNestedResponse converts model to nested response DTO
func (m *ModuleMapper) ToNestedResponse(module *models.Module) *dto.NestedModuleResponse {
	if module == nil {
		return nil
	}

	return &dto.NestedModuleResponse{
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
		Children:         []*dto.NestedModuleResponse{},
	}
}

// ToUserModuleResponse converts user module model to response DTO
func (m *ModuleMapper) ToUserModuleResponse(userModule *models.UserModule) *dto.UserModuleResponse {
	if userModule == nil {
		return nil
	}

	return &dto.UserModuleResponse{
		ModuleResponse: *m.ToResponse(&userModule.Module),
		CanRead:        userModule.CanRead,
		CanWrite:       userModule.CanWrite,
		CanDelete:      userModule.CanDelete,
	}
}

// ToUserModuleResponseList converts user module slice to response DTO slice
func (m *ModuleMapper) ToUserModuleResponseList(userModules []*models.UserModule) []*dto.UserModuleResponse {
	if userModules == nil {
		return nil
	}

	responses := make([]*dto.UserModuleResponse, len(userModules))
	for i, userModule := range userModules {
		responses[i] = m.ToUserModuleResponse(userModule)
	}
	return responses
}

// ToModel converts create request DTO to model
func (m *ModuleMapper) ToModel(req *dto.CreateModuleRequest) *models.Module {
	if req == nil {
		return nil
	}

	return &models.Module{
		Category:         req.Category,
		Name:             req.Name,
		URL:              req.URL,
		Icon:             req.Icon,
		Description:      req.Description,
		ParentID:         req.ParentID,
		SubscriptionTier: req.SubscriptionTier,
		IsActive:         true, // Default to active
	}
}

// UpdateModel updates model with update request DTO
func (m *ModuleMapper) UpdateModel(module *models.Module, req *dto.UpdateModuleRequest) {
	if module == nil || req == nil {
		return
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
}

// ToListResponse creates paginated list response
func (m *ModuleMapper) ToListResponse(modules []*models.Module, total int64, limit, offset int) *dto.ModuleListResponse {
	return &dto.ModuleListResponse{
		Data:    m.ToResponseList(modules),
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(modules)) < total,
	}
}
