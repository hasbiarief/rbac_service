package service

import (
	"fmt"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/pkg/rbac"
)

type ModuleService struct {
	moduleRepo  *repository.ModuleRepository
	rbacService *rbac.RBACService
}

func NewModuleService(moduleRepo *repository.ModuleRepository, rbacService *rbac.RBACService) *ModuleService {
	return &ModuleService{
		moduleRepo:  moduleRepo,
		rbacService: rbacService,
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

type ModuleTreeResponse struct {
	ID               int64                 `json:"id"`
	Category         string                `json:"category"`
	Name             string                `json:"name"`
	URL              string                `json:"url"`
	Icon             string                `json:"icon"`
	Description      string                `json:"description"`
	ParentID         *int64                `json:"parent_id"`
	SubscriptionTier string                `json:"subscription_tier"`
	IsActive         bool                  `json:"is_active"`
	Children         []*ModuleTreeResponse `json:"children"`
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

// GetModulesFiltered returns modules filtered by user permissions
func (s *ModuleService) GetModulesFiltered(userID int64, req *ModuleListRequest) ([]*ModuleResponse, error) {
	// Check if user is super admin - if so, return all modules
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetModules(req)
	}

	// Get filtered modules based on user permissions
	moduleInfos, err := s.rbacService.GetFilteredModules(userID, "read", req.Limit, req.Offset, req.Category, req.Search, req.IsActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered modules: %w", err)
	}

	var response []*ModuleResponse
	for _, moduleInfo := range moduleInfos {
		response = append(response, &ModuleResponse{
			ID:               moduleInfo.ID,
			Category:         moduleInfo.Category,
			Name:             moduleInfo.Name,
			URL:              moduleInfo.URL,
			Icon:             moduleInfo.Icon,
			Description:      moduleInfo.Description,
			ParentID:         moduleInfo.ParentID,
			SubscriptionTier: moduleInfo.SubscriptionTier,
			IsActive:         moduleInfo.IsActive,
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

type CheckAccessRequest struct {
	UserIdentity string `json:"user_identity" binding:"required"`
	ModuleURL    string `json:"module_url" binding:"required"`
}

func (s *ModuleService) GetModuleTree(category string) ([]*ModuleTreeResponse, error) {
	modules, err := s.moduleRepo.GetTreeStructure(category)
	if err != nil {
		return nil, err
	}

	return s.convertToTreeResponse(modules), nil
}

// GetModuleTreeFiltered returns module tree filtered by user permissions
func (s *ModuleService) GetModuleTreeFiltered(userID int64, category string) ([]*ModuleTreeResponse, error) {
	// Check if user is super admin - if so, return all modules
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetModuleTree(category)
	}

	// Get user permissions to filter tree
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get all modules first
	modules, err := s.moduleRepo.GetTreeStructure(category)
	if err != nil {
		return nil, err
	}

	// Filter modules based on permissions
	filteredModules := s.filterModuleTree(modules, permissions.Modules)
	return s.convertToTreeResponse(filteredModules), nil
}

// filterModuleTree recursively filters module tree based on permissions
func (s *ModuleService) filterModuleTree(modules []*models.ModuleWithChildren, permissions map[int64]rbac.ModulePermission) []*models.ModuleWithChildren {
	var filtered []*models.ModuleWithChildren

	for _, module := range modules {
		// Check if user has read permission for this module
		if perm, exists := permissions[module.ID]; exists && perm.CanRead {
			// Create a copy of the module
			filteredModule := &models.ModuleWithChildren{
				ID:               module.ID,
				Category:         module.Category,
				Name:             module.Name,
				URL:              module.URL,
				Icon:             module.Icon,
				Description:      module.Description,
				ParentID:         module.ParentID,
				SubscriptionTier: module.SubscriptionTier,
				IsActive:         module.IsActive,
				CreatedAt:        module.CreatedAt,
				UpdatedAt:        module.UpdatedAt,
			}

			// Recursively filter children
			if len(module.Children) > 0 {
				filteredModule.Children = s.filterModuleTree(module.Children, permissions)
			}

			filtered = append(filtered, filteredModule)
		}
	}

	return filtered
}

func (s *ModuleService) GetModuleTreeByParent(parentName string) ([]*ModuleTreeResponse, error) {
	modules, err := s.moduleRepo.GetTreeByParentName(parentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by parent name '%s': %w", parentName, err)
	}

	return s.convertToTreeResponse(modules), nil
}

// GetModuleTreeByParentFiltered returns module tree by parent filtered by user permissions
func (s *ModuleService) GetModuleTreeByParentFiltered(userID int64, parentName string) ([]*ModuleTreeResponse, error) {
	// Check if user is super admin - if so, return all modules
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetModuleTreeByParent(parentName)
	}

	// Get user permissions to filter tree
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get all modules first
	modules, err := s.moduleRepo.GetTreeByParentName(parentName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by parent name '%s': %w", parentName, err)
	}

	// Filter modules based on permissions
	filteredModules := s.filterModuleTree(modules, permissions.Modules)
	return s.convertToTreeResponse(filteredModules), nil
}

func (s *ModuleService) convertToTreeResponse(modules []*models.ModuleWithChildren) []*ModuleTreeResponse {
	var response []*ModuleTreeResponse

	for _, module := range modules {
		treeResponse := &ModuleTreeResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
			Children:         s.convertToTreeResponse(module.Children),
		}
		response = append(response, treeResponse)
	}

	return response
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

// GetModuleChildrenFiltered returns module children filtered by user permissions
func (s *ModuleService) GetModuleChildrenFiltered(userID int64, parentID int64) ([]*ModuleResponse, error) {
	// Check if user is super admin - if so, return all modules
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetModuleChildren(parentID)
	}

	// Get user permissions
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get all children first
	modules, err := s.moduleRepo.GetChildren(parentID)
	if err != nil {
		return nil, err
	}

	// Filter based on permissions
	var response []*ModuleResponse
	for _, module := range modules {
		// Check if user has read permission for this module
		if perm, exists := permissions.Modules[module.ID]; exists && perm.CanRead {
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

// GetModuleAncestorsFiltered returns module ancestors filtered by user permissions
func (s *ModuleService) GetModuleAncestorsFiltered(userID int64, moduleID int64) ([]*ModuleResponse, error) {
	// Check if user is super admin - if so, return all modules
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}

	if isSuperAdmin {
		return s.GetModuleAncestors(moduleID)
	}

	// Get user permissions
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get all ancestors first
	modules, err := s.moduleRepo.GetAncestors(moduleID)
	if err != nil {
		return nil, err
	}

	// Filter based on permissions
	var response []*ModuleResponse
	for _, module := range modules {
		// Check if user has read permission for this module
		if perm, exists := permissions.Modules[module.ID]; exists && perm.CanRead {
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

func (s *ModuleService) GetUserModulesGrouped(userID int64) (map[string][][]string, error) {
	// Use the user repository method directly since it already handles subscription filtering
	userRepo := &repository.UserRepository{} // This is not ideal, but for now...
	return userRepo.GetUserModulesGroupedWithSubscription(userID)
}

func (s *ModuleService) CheckUserAccess(req *CheckAccessRequest) (bool, error) {
	return s.moduleRepo.CheckUserAccess(req.UserIdentity, req.ModuleURL)
}
