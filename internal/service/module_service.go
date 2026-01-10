package service

import (
	"fmt"
	"gin-scalable-api/internal/dto"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/mapper"
	"gin-scalable-api/internal/models"
	"gin-scalable-api/pkg/rbac"
)

type ModuleService struct {
	moduleRepo   interfaces.ModuleRepositoryInterface
	userRepo     interfaces.UserRepositoryInterface
	rbacService  *rbac.RBACService
	moduleMapper *mapper.ModuleMapper
}

func NewModuleService(moduleRepo interfaces.ModuleRepositoryInterface, userRepo interfaces.UserRepositoryInterface, rbacService *rbac.RBACService) *ModuleService {
	return &ModuleService{
		moduleRepo:   moduleRepo,
		userRepo:     userRepo,
		rbacService:  rbacService,
		moduleMapper: mapper.NewModuleMapper(),
	}
}

func (s *ModuleService) GetModules(req *dto.ModuleListRequest) (*dto.ModuleListResponse, error) {
	// Set default values jika tidak disediakan
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	modules, err := s.moduleRepo.GetAll(limit, offset, req.Search, req.Category, req.SubscriptionTier, req.ParentID, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Dapatkan total count untuk pagination
	total, err := s.moduleRepo.Count(req.Search, req.Category, req.SubscriptionTier, req.ParentID, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Konversi ke DTO menggunakan mapper
	var moduleResponses []*dto.ModuleResponse
	for _, module := range modules {
		moduleResponses = append(moduleResponses, s.moduleMapper.ToResponse(module))
	}

	return &dto.ModuleListResponse{
		Data:    moduleResponses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(moduleResponses)) < total,
	}, nil
}

// GetModulesFiltered mengembalikan modul yang difilter berdasarkan izin pengguna
func (s *ModuleService) GetModulesFiltered(requestingUserID int64, req *dto.ModuleListRequest) (*dto.ModuleListResponse, error) {
	// Periksa apakah user adalah super admin - jika ya, kembalikan semua modul
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		return s.GetModules(req)
	}

	// Dapatkan modul yang difilter berdasarkan izin pengguna
	moduleInfos, err := s.rbacService.GetFilteredModules(requestingUserID, "read", req.Limit, req.Offset, req.Category, req.Search, req.IsActive)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan modul yang difilter: %w", err)
	}

	var moduleResponses []*dto.ModuleResponse
	for _, moduleInfo := range moduleInfos {
		moduleResponses = append(moduleResponses, &dto.ModuleResponse{
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

	// Hitung total untuk pagination (perkiraan berdasarkan hasil yang difilter)
	total := int64(len(moduleResponses))

	return &dto.ModuleListResponse{
		Data:    moduleResponses,
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: int64(req.Offset+len(moduleResponses)) < total,
	}, nil
}

func (s *ModuleService) GetModulesNested(req *dto.ModuleListRequest) ([]*dto.NestedModuleResponse, error) {
	// Set default values jika tidak disediakan
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	modules, err := s.moduleRepo.GetAll(limit, offset, req.Search, req.Category, req.SubscriptionTier, req.ParentID, req.IsActive)
	if err != nil {
		return nil, err
	}

	// Konversi ke nested response menggunakan mapper
	var nestedResponses []*dto.NestedModuleResponse
	for _, module := range modules {
		nestedResponses = append(nestedResponses, s.moduleMapper.ToNestedResponse(module))
	}

	return nestedResponses, nil
}

func (s *ModuleService) GetModuleByID(id int64) (*dto.ModuleResponse, error) {
	module, err := s.moduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.moduleMapper.ToResponse(module), nil
}

func (s *ModuleService) GetUserModules(userID int64, category string, limit int) ([]*dto.ModuleResponse, error) {
	// Untuk saat ini, gunakan companyID = 0 sebagai default, ini harus diteruskan dari handler
	userModules, err := s.moduleRepo.GetUserModules(userID, 0)
	if err != nil {
		return nil, err
	}

	var response []*dto.ModuleResponse
	for _, userModule := range userModules {
		// Konversi UserModule ke Module untuk mapper
		module := &models.Module{
			ID:               userModule.ID,
			Category:         userModule.Category,
			Name:             userModule.Name,
			URL:              userModule.URL,
			Icon:             userModule.Icon,
			Description:      userModule.Description,
			ParentID:         userModule.ParentID,
			SubscriptionTier: userModule.SubscriptionTier,
			IsActive:         userModule.IsActive,
		}
		response = append(response, s.moduleMapper.ToResponse(module))
	}

	return response, nil
}

func (s *ModuleService) CreateModule(req *dto.CreateModuleRequest) (*dto.ModuleResponse, error) {
	// Konversi DTO ke model menggunakan mapper
	module := s.moduleMapper.ToModel(req)

	if err := s.moduleRepo.Create(module); err != nil {
		return nil, err
	}

	return s.moduleMapper.ToResponse(module), nil
}

func (s *ModuleService) UpdateModule(id int64, req *dto.UpdateModuleRequest) (*dto.ModuleResponse, error) {
	module, err := s.moduleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields menggunakan mapper
	s.moduleMapper.UpdateModel(module, req)

	if err := s.moduleRepo.Update(module); err != nil {
		return nil, err
	}

	return s.moduleMapper.ToResponse(module), nil
}

func (s *ModuleService) DeleteModule(id int64) error {
	return s.moduleRepo.Delete(id)
}

// filterModulesByPermissions memfilter modul berdasarkan izin secara rekursif
func (s *ModuleService) filterModulesByPermissions(modules []*models.Module, permissions map[int64]rbac.ModulePermission) []*models.Module {
	var filtered []*models.Module

	for _, module := range modules {
		// Periksa apakah user memiliki izin baca untuk modul ini
		if perm, exists := permissions[module.ID]; exists && perm.CanRead {
			filtered = append(filtered, module)
		}
	}

	return filtered
}

// convertModulesToTreeResponse mengkonversi slice modul ke tree response
func (s *ModuleService) convertModulesToTreeResponse(modules []*models.Module) []*dto.ModuleTreeResponse {
	var response []*dto.ModuleTreeResponse

	for _, module := range modules {
		treeResponse := &dto.ModuleTreeResponse{
			ID:               module.ID,
			Category:         module.Category,
			Name:             module.Name,
			URL:              module.URL,
			Icon:             module.Icon,
			Description:      module.Description,
			ParentID:         module.ParentID,
			SubscriptionTier: module.SubscriptionTier,
			IsActive:         module.IsActive,
			CreatedAt:        module.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:        module.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Children:         []*dto.ModuleTreeResponse{}, // Empty children for now
			Level:            0,                           // Calculate level if needed
			Path:             module.Name,                 // Simple path for now
		}
		response = append(response, treeResponse)
	}

	return response
}

func (s *ModuleService) GetModuleTreeByParentFiltered(userID int64, parentName string) ([]*dto.ModuleTreeResponse, error) {
	// Periksa apakah user adalah super admin - jika ya, kembalikan semua modul
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		modules, err := s.moduleRepo.GetTreeByParentName(parentName, 0) // userID 0 untuk akses admin
		if err != nil {
			return nil, fmt.Errorf("gagal mendapatkan tree berdasarkan nama parent '%s': %w", parentName, err)
		}
		return s.convertModulesToTreeResponse(modules), nil
	}

	// Dapatkan izin pengguna untuk memfilter tree
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan izin pengguna: %w", err)
	}

	// Dapatkan semua modul terlebih dahulu
	modules, err := s.moduleRepo.GetTreeByParentName(parentName, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan tree berdasarkan nama parent '%s': %w", parentName, err)
	}

	// Filter modul berdasarkan izin
	filteredModules := s.filterModulesByPermissions(modules, permissions.Modules)
	return s.convertModulesToTreeResponse(filteredModules), nil
}

func (s *ModuleService) GetModuleTreeFiltered(userID int64, category string) ([]*dto.ModuleTreeResponse, error) {
	// Periksa apakah user adalah super admin - jika ya, kembalikan semua modul
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		modules, err := s.moduleRepo.GetTreeStructure(category, 0) // userID 0 untuk akses admin
		if err != nil {
			return nil, err
		}
		return s.convertModulesToTreeResponse(modules), nil
	}

	// Dapatkan izin pengguna untuk memfilter tree
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan izin pengguna: %w", err)
	}

	// Dapatkan semua modul terlebih dahulu
	modules, err := s.moduleRepo.GetTreeStructure(category, userID)
	if err != nil {
		return nil, err
	}

	// Filter modul berdasarkan izin
	filteredModules := s.filterModulesByPermissions(modules, permissions.Modules)
	return s.convertModulesToTreeResponse(filteredModules), nil
}

func (s *ModuleService) GetModuleChildrenFiltered(userID int64, id int64) ([]*dto.ModuleResponse, error) {
	// Periksa apakah user adalah super admin - jika ya, kembalikan semua modul
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		modules, err := s.moduleRepo.GetChildren(id)
		if err != nil {
			return nil, err
		}

		var response []*dto.ModuleResponse
		for _, module := range modules {
			response = append(response, s.moduleMapper.ToResponse(module))
		}
		return response, nil
	}

	// Dapatkan izin pengguna
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan izin pengguna: %w", err)
	}

	// Dapatkan semua children terlebih dahulu
	modules, err := s.moduleRepo.GetChildren(id)
	if err != nil {
		return nil, err
	}

	// Filter berdasarkan izin
	var response []*dto.ModuleResponse
	for _, module := range modules {
		// Periksa apakah user memiliki izin baca untuk modul ini
		if perm, exists := permissions.Modules[module.ID]; exists && perm.CanRead {
			response = append(response, s.moduleMapper.ToResponse(module))
		}
	}

	return response, nil
}

func (s *ModuleService) GetModuleAncestorsFiltered(userID int64, id int64) ([]*dto.ModuleResponse, error) {
	// Periksa apakah user adalah super admin - jika ya, kembalikan semua modul
	isSuperAdmin, err := s.rbacService.IsSuperAdmin(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa status super admin: %w", err)
	}

	if isSuperAdmin {
		modules, err := s.moduleRepo.GetAncestors(id, 0) // userID 0 untuk akses admin
		if err != nil {
			return nil, err
		}

		var response []*dto.ModuleResponse
		for _, module := range modules {
			response = append(response, s.moduleMapper.ToResponse(module))
		}
		return response, nil
	}

	// Dapatkan izin pengguna
	permissions, err := s.rbacService.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan izin pengguna: %w", err)
	}

	// Dapatkan semua ancestors terlebih dahulu
	modules, err := s.moduleRepo.GetAncestors(id, userID)
	if err != nil {
		return nil, err
	}

	// Filter berdasarkan izin
	var response []*dto.ModuleResponse
	for _, module := range modules {
		// Periksa apakah user memiliki izin baca untuk modul ini
		if perm, exists := permissions.Modules[module.ID]; exists && perm.CanRead {
			response = append(response, s.moduleMapper.ToResponse(module))
		}
	}

	return response, nil
}
