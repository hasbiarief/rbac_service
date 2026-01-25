package apidoc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gin-scalable-api/internal/constants"
	"gin-scalable-api/internal/modules/apidoc/export"
	"gin-scalable-api/pkg/rbac"
)

// RBACServiceInterface defines the interface for RBAC operations
type RBACServiceInterface interface {
	HasPermission(userID int64, moduleID int64, permission string) (bool, error)
	GetUserPermissions(userID int64) (*rbac.UserPermissions, error)
	HasRole(userID int64, roleName string) (bool, error)
	GetAccessibleModules(userID int64, permission string) ([]int64, error)
	IsSuperAdmin(userID int64) (bool, error)
	GetFilteredModules(userID int64, permission string, limit, offset int, category, search string, isActive *bool) ([]*rbac.ModuleInfo, error)
}

// Service provides business logic for API documentation operations
type Service struct {
	repo          Repository
	rbacService   RBACServiceInterface
	db            *sql.DB
	exportManager *export.ExportManager
}

// NewService creates a new API documentation service instance
func NewService(repo Repository, rbacService RBACServiceInterface, db *sql.DB) *Service {
	return &Service{
		repo:          repo,
		rbacService:   rbacService,
		db:            db,
		exportManager: export.NewExportManager(),
	}
}

// NewServiceWithCache creates a new API documentation service instance with caching
func NewServiceWithCache(repo Repository, rbacService RBACServiceInterface, db *sql.DB, cache *export.ExportCache) *Service {
	return &Service{
		repo:          repo,
		rbacService:   rbacService,
		db:            db,
		exportManager: export.NewExportManagerWithCache(cache),
	}
}

// SetExportCache sets the export cache
func (s *Service) SetExportCache(cache *export.ExportCache) {
	s.exportManager.SetCache(cache)
}

// CollectionService methods

// CreateCollection creates a new API documentation collection with company isolation
func (s *Service) CreateCollection(req *CreateCollectionRequest, userID, companyID int64) (*CollectionResponse, error) {
	// Check user permissions for collection creation
	if err := s.checkCollectionPermission(userID, "write"); err != nil {
		return nil, err
	}

	// Validate request
	if err := ValidateCreateCollectionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert request to model
	collection := CreateCollectionRequestToModel(req, userID, companyID)

	// Create collection in database
	if err := s.repo.CreateCollection(collection); err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	// Convert to response DTO
	return CollectionToResponse(collection), nil
}

// GetCollections retrieves collections with pagination and filtering, enforcing company isolation
func (s *Service) GetCollections(req *CollectionListRequest, companyID int64) (*CollectionListResponse, error) {
	// Note: Permission check is done in the ByUser wrapper method
	// This method is called internally with pre-validated company isolation

	// Validate request
	if err := ValidateCollectionListRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Set default values
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get collections from repository with company isolation
	collections, total, err := s.repo.GetCollections(companyID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	// Convert to response DTOs
	var responses []*CollectionResponse
	for _, collection := range collections {
		responses = append(responses, CollectionToResponse(collection))
	}

	return &CollectionListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

// GetCollectionByID retrieves a collection by ID with authorization check
func (s *Service) GetCollectionByID(id, companyID int64) (*CollectionResponse, error) {
	// Get collection with company isolation
	collection, err := s.repo.GetCollectionByID(id, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	return CollectionToResponse(collection), nil
}

// GetCollectionWithStats retrieves a collection with statistics
func (s *Service) GetCollectionWithStats(id, companyID int64) (*CollectionWithStatsResponse, error) {
	// Get collection with stats and company isolation
	stats, err := s.repo.GetCollectionWithStats(id, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection with stats: %w", err)
	}

	return CollectionWithStatsToResponse(stats), nil
}

// UpdateCollection updates a collection with ownership validation
func (s *Service) UpdateCollection(id int64, req *UpdateCollectionRequest, userID, companyID int64) (*CollectionResponse, error) {
	// Check user permissions for collection updates
	if err := s.checkCollectionPermission(userID, "write"); err != nil {
		return nil, err
	}

	// Validate request
	if err := ValidateUpdateCollectionRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing collection to verify ownership and company isolation
	collection, err := s.repo.GetCollectionByID(id, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		collection.Name = *req.Name
	}
	if req.Description != nil {
		collection.Description = pointerToNullString(req.Description)
	}
	if req.Version != nil {
		collection.Version = *req.Version
	}
	if req.BaseURL != nil {
		collection.BaseURL = pointerToNullString(req.BaseURL)
	}
	if req.SchemaVersion != nil {
		collection.SchemaVersion = *req.SchemaVersion
	}
	if req.IsActive != nil {
		collection.IsActive = *req.IsActive
	}

	// Update in database
	if err := s.repo.UpdateCollection(collection); err != nil {
		return nil, fmt.Errorf("failed to update collection: %w", err)
	}

	// Invalidate export cache for this collection
	if err := s.exportManager.InvalidateCache(collection.ID); err != nil {
		// Log cache invalidation error but don't fail the update
		// In production, you might want to log this error
	}

	return CollectionToResponse(collection), nil
}

// DeleteCollection deletes a collection with cascade handling
func (s *Service) DeleteCollection(id, companyID int64) error {
	// Note: Permission check is done in the ByUser wrapper method
	// This method is called internally with pre-validated company isolation

	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(id, companyID)
	if err != nil {
		return fmt.Errorf("failed to get collection: %w", err)
	}

	// Delete collection (cascade will handle related records)
	if err := s.repo.DeleteCollection(id, companyID); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	// Invalidate export cache for this collection
	if err := s.exportManager.InvalidateCache(id); err != nil {
		// Log cache invalidation error but don't fail the delete
		// In production, you might want to log this error
	}

	return nil
}

// FolderService methods

// CreateFolder creates a new folder with hierarchy validation
func (s *Service) CreateFolder(collectionID int64, req *CreateFolderRequest, companyID int64) (*FolderResponse, error) {
	// Validate request
	if err := ValidateCreateFolderRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// If parent folder is specified, verify it exists and belongs to the same collection
	if req.ParentID != nil {
		parentFolder, err := s.repo.GetFolderByID(*req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found: %w", err)
		}
		if parentFolder.CollectionID != collectionID {
			return nil, fmt.Errorf("parent folder does not belong to the same collection")
		}
	}

	// Convert request to model
	folder := CreateFolderRequestToModel(req, collectionID)

	// Create folder in database
	if err := s.repo.CreateFolder(folder); err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return FolderToResponse(folder), nil
}

// GetFolders retrieves folders for a collection with nested structure
func (s *Service) GetFolders(collectionID, companyID int64) ([]*FolderWithChildrenResponse, error) {
	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// Get folders hierarchy
	folders, err := s.repo.GetFoldersHierarchy(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	// Convert to response DTOs
	var responses []*FolderWithChildrenResponse
	for _, folder := range folders {
		responses = append(responses, FolderWithChildrenToResponse(folder))
	}

	return responses, nil
}

// GetFolderByID retrieves a folder by ID
func (s *Service) GetFolderByID(id, companyID int64) (*FolderResponse, error) {
	folder, err := s.repo.GetFolderByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	// Verify the folder's collection belongs to the company
	_, err = s.repo.GetCollectionByID(folder.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	return FolderToResponse(folder), nil
}

// UpdateFolder updates a folder with parent validation
func (s *Service) UpdateFolder(id int64, req *UpdateFolderRequest, companyID int64) (*FolderResponse, error) {
	// Validate request
	if err := ValidateUpdateFolderRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing folder
	folder, err := s.repo.GetFolderByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get folder: %w", err)
	}

	// Verify the folder's collection belongs to the company
	_, err = s.repo.GetCollectionByID(folder.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// If parent folder is being changed, validate it
	if req.ParentID != nil {
		if *req.ParentID != 0 { // 0 means root level
			parentFolder, err := s.repo.GetFolderByID(*req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("parent folder not found: %w", err)
			}
			if parentFolder.CollectionID != folder.CollectionID {
				return nil, fmt.Errorf("parent folder must belong to the same collection")
			}
			// Prevent circular references
			if *req.ParentID == id {
				return nil, fmt.Errorf("folder cannot be its own parent")
			}
		}
		folder.ParentID = pointerToNullInt64(req.ParentID)
	}

	// Update other fields if provided
	if req.Name != nil {
		folder.Name = *req.Name
	}
	if req.Description != nil {
		folder.Description = pointerToNullString(req.Description)
	}
	if req.SortOrder != nil {
		folder.SortOrder = *req.SortOrder
	}

	// Update in database
	if err := s.repo.UpdateFolder(folder); err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	return FolderToResponse(folder), nil
}

// DeleteFolder deletes a folder with children handling
func (s *Service) DeleteFolder(id, companyID int64) error {
	// Get folder to verify access
	folder, err := s.repo.GetFolderByID(id)
	if err != nil {
		return fmt.Errorf("failed to get folder: %w", err)
	}

	// Verify the folder's collection belongs to the company
	_, err = s.repo.GetCollectionByID(folder.CollectionID, companyID)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	// Delete folder (cascade will handle children and endpoints)
	if err := s.repo.DeleteFolder(id); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}

// ReorderFolders updates the sort order of multiple folders
func (s *Service) ReorderFolders(collectionID, companyID int64, folderOrders []FolderOrderRequest) error {
	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return fmt.Errorf("collection not found or access denied: %w", err)
	}

	// Validate all folders belong to the collection
	for _, order := range folderOrders {
		folder, err := s.repo.GetFolderByID(order.FolderID)
		if err != nil {
			return fmt.Errorf("folder %d not found: %w", order.FolderID, err)
		}
		if folder.CollectionID != collectionID {
			return fmt.Errorf("folder %d does not belong to collection %d", order.FolderID, collectionID)
		}
	}

	// Update sort orders
	for _, order := range folderOrders {
		folder, err := s.repo.GetFolderByID(order.FolderID)
		if err != nil {
			return fmt.Errorf("failed to get folder %d: %w", order.FolderID, err)
		}

		folder.SortOrder = order.SortOrder
		if err := s.repo.UpdateFolder(folder); err != nil {
			return fmt.Errorf("failed to update folder %d sort order: %w", order.FolderID, err)
		}
	}

	return nil
}

// EndpointService methods

// CreateEndpoint creates a new endpoint with validation
func (s *Service) CreateEndpoint(collectionID int64, req *CreateEndpointRequest, companyID int64) (*EndpointResponse, error) {
	// Validate request
	if err := ValidateCreateEndpointRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// If folder is specified, verify it exists and belongs to the same collection
	if req.FolderID != nil {
		folder, err := s.repo.GetFolderByID(*req.FolderID)
		if err != nil {
			return nil, fmt.Errorf("folder not found: %w", err)
		}
		if folder.CollectionID != collectionID {
			return nil, fmt.Errorf("folder does not belong to the same collection")
		}
	}

	// Convert request to model
	endpoint := CreateEndpointRequestToModel(req, collectionID)

	// Create endpoint in database
	if err := s.repo.CreateEndpoint(endpoint); err != nil {
		return nil, fmt.Errorf("failed to create endpoint: %w", err)
	}

	return EndpointToResponse(endpoint), nil
}

// GetEndpoints retrieves endpoints with filtering and pagination
func (s *Service) GetEndpoints(collectionID, companyID int64, req *EndpointListRequest) (*EndpointListResponse, error) {
	// Validate request
	if err := ValidateEndpointListRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// Set default values
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get endpoints from repository
	endpoints, total, err := s.repo.GetEndpointsByCollectionID(collectionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	// Convert to response DTOs
	var responses []*EndpointResponse
	for _, endpoint := range endpoints {
		responses = append(responses, EndpointToResponse(endpoint))
	}

	return &EndpointListResponse{
		Data:    responses,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+len(responses)) < total,
	}, nil
}

// GetEndpointByID retrieves an endpoint by ID
func (s *Service) GetEndpointByID(id, companyID int64) (*EndpointResponse, error) {
	endpoint, err := s.repo.GetEndpointByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Verify the endpoint's collection belongs to the company
	_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	return EndpointToResponse(endpoint), nil
}

// GetEndpointWithDetails retrieves an endpoint with all its related data
func (s *Service) GetEndpointWithDetails(id, companyID int64) (*EndpointWithDetailsResponse, error) {
	// First verify access
	endpoint, err := s.repo.GetEndpointByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Verify the endpoint's collection belongs to the company
	_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Get endpoint with all details
	details, err := s.repo.GetEndpointWithDetails(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint details: %w", err)
	}

	return EndpointWithDetailsToResponse(details), nil
}

// UpdateEndpoint updates an endpoint with partial updates
func (s *Service) UpdateEndpoint(id int64, req *UpdateEndpointRequest, companyID int64) (*EndpointResponse, error) {
	// Validate request
	if err := ValidateUpdateEndpointRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing endpoint
	endpoint, err := s.repo.GetEndpointByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Verify the endpoint's collection belongs to the company
	_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// If folder is being changed, validate it
	if req.FolderID != nil {
		if *req.FolderID != 0 { // 0 means no folder
			folder, err := s.repo.GetFolderByID(*req.FolderID)
			if err != nil {
				return nil, fmt.Errorf("folder not found: %w", err)
			}
			if folder.CollectionID != endpoint.CollectionID {
				return nil, fmt.Errorf("folder must belong to the same collection")
			}
		}
		endpoint.FolderID = pointerToNullInt64(req.FolderID)
	}

	// Update other fields if provided
	if req.Name != nil {
		endpoint.Name = *req.Name
	}
	if req.Description != nil {
		endpoint.Description = pointerToNullString(req.Description)
	}
	if req.Method != nil {
		endpoint.Method = *req.Method
	}
	if req.URL != nil {
		endpoint.URL = *req.URL
	}
	if req.SortOrder != nil {
		endpoint.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		endpoint.IsActive = *req.IsActive
	}

	// Update in database
	if err := s.repo.UpdateEndpoint(endpoint); err != nil {
		return nil, fmt.Errorf("failed to update endpoint: %w", err)
	}

	return EndpointToResponse(endpoint), nil
}

// DeleteEndpoint deletes an endpoint with cleanup
func (s *Service) DeleteEndpoint(id, companyID int64) error {
	// Get endpoint to verify access
	endpoint, err := s.repo.GetEndpointByID(id)
	if err != nil {
		return fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Verify the endpoint's collection belongs to the company
	_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	// Delete endpoint (cascade will handle related records)
	if err := s.repo.DeleteEndpoint(id); err != nil {
		return fmt.Errorf("failed to delete endpoint: %w", err)
	}

	return nil
}

// BulkCreateEndpoints creates multiple endpoints in a single operation
func (s *Service) BulkCreateEndpoints(collectionID, companyID int64, requests []*CreateEndpointRequest) ([]*EndpointResponse, error) {
	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	var responses []*EndpointResponse

	// Process each endpoint request
	for i, req := range requests {
		// Validate request
		if err := ValidateCreateEndpointRequest(req); err != nil {
			return nil, fmt.Errorf("validation failed for endpoint %d: %w", i+1, err)
		}

		// If folder is specified, verify it exists and belongs to the same collection
		if req.FolderID != nil {
			folder, err := s.repo.GetFolderByID(*req.FolderID)
			if err != nil {
				return nil, fmt.Errorf("folder not found for endpoint %d: %w", i+1, err)
			}
			if folder.CollectionID != collectionID {
				return nil, fmt.Errorf("folder does not belong to the same collection for endpoint %d", i+1)
			}
		}

		// Convert request to model
		endpoint := CreateEndpointRequestToModel(req, collectionID)

		// Create endpoint in database
		if err := s.repo.CreateEndpoint(endpoint); err != nil {
			return nil, fmt.Errorf("failed to create endpoint %d: %w", i+1, err)
		}

		responses = append(responses, EndpointToResponse(endpoint))
	}

	return responses, nil
}

// BulkUpdateEndpoints updates multiple endpoints in a single operation
func (s *Service) BulkUpdateEndpoints(companyID int64, updates []EndpointBulkUpdateRequest) ([]*EndpointResponse, error) {
	var responses []*EndpointResponse

	// Process each update request
	for i, update := range updates {
		// Validate request
		if err := ValidateUpdateEndpointRequest(&update.UpdateRequest); err != nil {
			return nil, fmt.Errorf("validation failed for endpoint %d: %w", i+1, err)
		}

		// Get existing endpoint
		endpoint, err := s.repo.GetEndpointByID(update.EndpointID)
		if err != nil {
			return nil, fmt.Errorf("failed to get endpoint %d: %w", update.EndpointID, err)
		}

		// Verify the endpoint's collection belongs to the company
		_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
		if err != nil {
			return nil, fmt.Errorf("access denied for endpoint %d: %w", update.EndpointID, err)
		}

		// Apply updates
		req := &update.UpdateRequest
		if req.FolderID != nil {
			if *req.FolderID != 0 { // 0 means no folder
				folder, err := s.repo.GetFolderByID(*req.FolderID)
				if err != nil {
					return nil, fmt.Errorf("folder not found for endpoint %d: %w", update.EndpointID, err)
				}
				if folder.CollectionID != endpoint.CollectionID {
					return nil, fmt.Errorf("folder must belong to the same collection for endpoint %d", update.EndpointID)
				}
			}
			endpoint.FolderID = pointerToNullInt64(req.FolderID)
		}

		if req.Name != nil {
			endpoint.Name = *req.Name
		}
		if req.Description != nil {
			endpoint.Description = pointerToNullString(req.Description)
		}
		if req.Method != nil {
			endpoint.Method = *req.Method
		}
		if req.URL != nil {
			endpoint.URL = *req.URL
		}
		if req.SortOrder != nil {
			endpoint.SortOrder = *req.SortOrder
		}
		if req.IsActive != nil {
			endpoint.IsActive = *req.IsActive
		}

		// Update in database
		if err := s.repo.UpdateEndpoint(endpoint); err != nil {
			return nil, fmt.Errorf("failed to update endpoint %d: %w", update.EndpointID, err)
		}

		responses = append(responses, EndpointToResponse(endpoint))
	}

	return responses, nil
}

// BulkDeleteEndpoints deletes multiple endpoints in a single operation
func (s *Service) BulkDeleteEndpoints(companyID int64, endpointIDs []int64) error {
	// Verify all endpoints exist and belong to company
	for _, id := range endpointIDs {
		endpoint, err := s.repo.GetEndpointByID(id)
		if err != nil {
			return fmt.Errorf("failed to get endpoint %d: %w", id, err)
		}

		// Verify the endpoint's collection belongs to the company
		_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
		if err != nil {
			return fmt.Errorf("access denied for endpoint %d: %w", id, err)
		}
	}

	// Delete all endpoints
	for _, id := range endpointIDs {
		if err := s.repo.DeleteEndpoint(id); err != nil {
			return fmt.Errorf("failed to delete endpoint %d: %w", id, err)
		}
	}

	return nil
}

// BulkMoveEndpoints moves multiple endpoints to a different folder
func (s *Service) BulkMoveEndpoints(companyID int64, endpointIDs []int64, targetFolderID *int64) error {
	// If target folder is specified, verify it exists
	var targetCollectionID int64
	if targetFolderID != nil && *targetFolderID != 0 {
		folder, err := s.repo.GetFolderByID(*targetFolderID)
		if err != nil {
			return fmt.Errorf("target folder not found: %w", err)
		}
		targetCollectionID = folder.CollectionID

		// Verify folder's collection belongs to company
		_, err = s.repo.GetCollectionByID(targetCollectionID, companyID)
		if err != nil {
			return fmt.Errorf("access denied to target folder's collection: %w", err)
		}
	}

	// Process each endpoint
	for _, id := range endpointIDs {
		endpoint, err := s.repo.GetEndpointByID(id)
		if err != nil {
			return fmt.Errorf("failed to get endpoint %d: %w", id, err)
		}

		// Verify the endpoint's collection belongs to the company
		_, err = s.repo.GetCollectionByID(endpoint.CollectionID, companyID)
		if err != nil {
			return fmt.Errorf("access denied for endpoint %d: %w", id, err)
		}

		// If target folder is specified, ensure it belongs to the same collection
		if targetFolderID != nil && *targetFolderID != 0 {
			if targetCollectionID != endpoint.CollectionID {
				return fmt.Errorf("target folder must belong to the same collection as endpoint %d", id)
			}
		}

		// Update folder assignment
		endpoint.FolderID = pointerToNullInt64(targetFolderID)
		if err := s.repo.UpdateEndpoint(endpoint); err != nil {
			return fmt.Errorf("failed to move endpoint %d: %w", id, err)
		}
	}

	return nil
}

// EnvironmentService methods

// CreateEnvironment creates a new environment with uniqueness check
func (s *Service) CreateEnvironment(collectionID, companyID int64, req *CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Validate request
	if err := ValidateCreateEnvironmentRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// Check for name uniqueness within the collection
	existingEnvironments, err := s.repo.GetEnvironmentsByCollectionID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing environments: %w", err)
	}

	for _, env := range existingEnvironments {
		if env.Name == req.Name {
			return nil, fmt.Errorf("environment with name '%s' already exists in this collection", req.Name)
		}
	}

	// If this is set as default, unset other defaults
	isDefault := false
	if req.IsDefault != nil {
		isDefault = *req.IsDefault
	}

	if isDefault {
		for _, env := range existingEnvironments {
			if env.IsDefault {
				env.IsDefault = false
				if err := s.repo.UpdateEnvironment(env); err != nil {
					return nil, fmt.Errorf("failed to unset previous default environment: %w", err)
				}
			}
		}
	}

	// Convert request to model
	environment := &Environment{
		CollectionID: collectionID,
		Name:         req.Name,
		Description:  pointerToNullString(req.Description),
		IsDefault:    isDefault,
	}

	// Create environment in database
	if err := s.repo.CreateEnvironment(environment); err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return EnvironmentToResponse(environment), nil
}

// GetEnvironments retrieves environments for a collection
func (s *Service) GetEnvironments(collectionID, companyID int64) ([]*EnvironmentResponse, error) {
	// Verify collection exists and belongs to company
	_, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found or access denied: %w", err)
	}

	// Get environments from repository
	environments, err := s.repo.GetEnvironmentsByCollectionID(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments: %w", err)
	}

	// Convert to response DTOs
	var responses []*EnvironmentResponse
	for _, environment := range environments {
		responses = append(responses, EnvironmentToResponse(environment))
	}

	return responses, nil
}

// GetEnvironmentByID retrieves an environment by ID
func (s *Service) GetEnvironmentByID(id, companyID int64) (*EnvironmentResponse, error) {
	environment, err := s.repo.GetEnvironmentByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	return EnvironmentToResponse(environment), nil
}

// GetEnvironmentWithVariables retrieves an environment with its variables
func (s *Service) GetEnvironmentWithVariables(id, companyID int64) (*EnvironmentWithVariablesResponse, error) {
	environment, err := s.repo.GetEnvironmentByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Get environment with variables
	envWithVars, err := s.repo.GetEnvironmentWithVariables(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment with variables: %w", err)
	}

	return EnvironmentWithVariablesToResponse(envWithVars), nil
}

// UpdateEnvironment updates an environment with variable management
func (s *Service) UpdateEnvironment(id int64, req *UpdateEnvironmentRequest, companyID int64) (*EnvironmentResponse, error) {
	// Validate request
	if err := ValidateUpdateEnvironmentRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing environment
	environment, err := s.repo.GetEnvironmentByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Check for name uniqueness if name is being changed
	if req.Name != nil && *req.Name != environment.Name {
		existingEnvironments, err := s.repo.GetEnvironmentsByCollectionID(environment.CollectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing environments: %w", err)
		}

		for _, env := range existingEnvironments {
			if env.ID != id && env.Name == *req.Name {
				return nil, fmt.Errorf("environment with name '%s' already exists in this collection", *req.Name)
			}
		}
		environment.Name = *req.Name
	}

	// Update other fields if provided
	if req.Description != nil {
		environment.Description = pointerToNullString(req.Description)
	}

	// Handle default flag changes
	if req.IsDefault != nil {
		if *req.IsDefault && !environment.IsDefault {
			// Setting as default, unset other defaults
			existingEnvironments, err := s.repo.GetEnvironmentsByCollectionID(environment.CollectionID)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing environments: %w", err)
			}

			for _, env := range existingEnvironments {
				if env.ID != id && env.IsDefault {
					env.IsDefault = false
					if err := s.repo.UpdateEnvironment(env); err != nil {
						return nil, fmt.Errorf("failed to unset previous default environment: %w", err)
					}
				}
			}
		}
		environment.IsDefault = *req.IsDefault
	}

	// Update in database
	if err := s.repo.UpdateEnvironment(environment); err != nil {
		return nil, fmt.Errorf("failed to update environment: %w", err)
	}

	return EnvironmentToResponse(environment), nil
}

// DeleteEnvironment deletes an environment with dependency check
func (s *Service) DeleteEnvironment(id, companyID int64) error {
	// Get environment to verify access
	environment, err := s.repo.GetEnvironmentByID(id)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	// Check if this is the only environment in the collection
	environments, err := s.repo.GetEnvironmentsByCollectionID(environment.CollectionID)
	if err != nil {
		return fmt.Errorf("failed to check existing environments: %w", err)
	}

	if len(environments) == 1 {
		return fmt.Errorf("cannot delete the last environment in a collection")
	}

	// If this is the default environment, set another one as default
	if environment.IsDefault {
		for _, env := range environments {
			if env.ID != id {
				env.IsDefault = true
				if err := s.repo.UpdateEnvironment(env); err != nil {
					return fmt.Errorf("failed to set new default environment: %w", err)
				}
				break
			}
		}
	}

	// Delete environment (cascade will handle variables)
	if err := s.repo.DeleteEnvironment(id); err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	return nil
}

// ManageEnvironmentVariables CRUD operations

// CreateEnvironmentVariable creates a new environment variable
func (s *Service) CreateEnvironmentVariable(environmentID, companyID int64, req *CreateEnvironmentVariableRequest) (*EnvironmentVariableResponse, error) {
	// Validate request
	if err := ValidateCreateEnvironmentVariableRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get environment and verify access
	environment, err := s.repo.GetEnvironmentByID(environmentID)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Check for key name uniqueness within the environment
	existingVariables, err := s.repo.GetEnvironmentVariables(environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing variables: %w", err)
	}

	for _, variable := range existingVariables {
		if variable.KeyName == req.KeyName {
			return nil, fmt.Errorf("variable with key '%s' already exists in this environment", req.KeyName)
		}
	}

	// Convert request to model
	isSecret := false
	if req.IsSecret != nil {
		isSecret = *req.IsSecret
	}

	variable := &EnvironmentVariable{
		EnvironmentID: environmentID,
		KeyName:       req.KeyName,
		Value:         pointerToNullString(req.Value),
		Description:   pointerToNullString(req.Description),
		IsSecret:      isSecret,
	}

	// Create variable in database
	if err := s.repo.CreateEnvironmentVariable(variable); err != nil {
		return nil, fmt.Errorf("failed to create environment variable: %w", err)
	}

	return EnvironmentVariableToResponse(variable), nil
}

// GetEnvironmentVariables retrieves variables for an environment
func (s *Service) GetEnvironmentVariables(environmentID, companyID int64) ([]*EnvironmentVariableResponse, error) {
	// Get environment and verify access
	environment, err := s.repo.GetEnvironmentByID(environmentID)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Get variables from repository
	variables, err := s.repo.GetEnvironmentVariables(environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %w", err)
	}

	// Convert to response DTOs
	var responses []*EnvironmentVariableResponse
	for _, variable := range variables {
		responses = append(responses, EnvironmentVariableToResponse(variable))
	}

	return responses, nil
}

// UpdateEnvironmentVariable updates an environment variable
func (s *Service) UpdateEnvironmentVariable(id int64, req *UpdateEnvironmentVariableRequest, companyID int64) (*EnvironmentVariableResponse, error) {
	// Validate request
	if err := ValidateUpdateEnvironmentVariableRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing variable
	variable, err := s.repo.GetEnvironmentVariableByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variable: %w", err)
	}

	// Get environment and verify access
	environment, err := s.repo.GetEnvironmentByID(variable.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	// Check for key name uniqueness if key name is being changed
	if req.KeyName != nil && *req.KeyName != variable.KeyName {
		existingVariables, err := s.repo.GetEnvironmentVariables(variable.EnvironmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing variables: %w", err)
		}

		for _, existingVar := range existingVariables {
			if existingVar.ID != id && existingVar.KeyName == *req.KeyName {
				return nil, fmt.Errorf("variable with key '%s' already exists in this environment", *req.KeyName)
			}
		}
		variable.KeyName = *req.KeyName
	}

	// Update other fields if provided
	if req.Value != nil {
		variable.Value = pointerToNullString(req.Value)
	}
	if req.Description != nil {
		variable.Description = pointerToNullString(req.Description)
	}
	if req.IsSecret != nil {
		variable.IsSecret = *req.IsSecret
	}

	// Update in database
	if err := s.repo.UpdateEnvironmentVariable(variable); err != nil {
		return nil, fmt.Errorf("failed to update environment variable: %w", err)
	}

	return EnvironmentVariableToResponse(variable), nil
}

// DeleteEnvironmentVariable deletes an environment variable
func (s *Service) DeleteEnvironmentVariable(id, companyID int64) error {
	// Get variable to verify access
	variable, err := s.repo.GetEnvironmentVariableByID(id)
	if err != nil {
		return fmt.Errorf("failed to get environment variable: %w", err)
	}

	// Get environment and verify access
	environment, err := s.repo.GetEnvironmentByID(variable.EnvironmentID)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	// Verify the environment's collection belongs to the company
	_, err = s.repo.GetCollectionByID(environment.CollectionID, companyID)
	if err != nil {
		return fmt.Errorf("access denied: %w", err)
	}

	// Delete variable
	if err := s.repo.DeleteEnvironmentVariable(id); err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	return nil
}

// Helper method to get user's company ID
func (s *Service) getUserCompanyID(userID int64) (int64, error) {
	var companyID int64
	query := `SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1`
	err := s.db.QueryRow(query, userID).Scan(&companyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user company: %w", err)
	}
	return companyID, nil
}

// RBAC permission checking methods

// checkCollectionPermission checks if user has permission for collection operations
func (s *Service) checkCollectionPermission(userID int64, permission string) error {
	// Skip permission check if rbacService is nil (for unit tests)
	if s.rbacService == nil {
		return nil
	}

	hasPermission, err := s.rbacService.HasPermission(userID, constants.ModuleAPIDocCollections, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// checkEndpointPermission checks if user has permission for endpoint operations
func (s *Service) checkEndpointPermission(userID int64, permission string) error {
	// Skip permission check if rbacService is nil (for unit tests)
	if s.rbacService == nil {
		return nil
	}

	hasPermission, err := s.rbacService.HasPermission(userID, constants.ModuleAPIDocEndpoints, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// checkEnvironmentPermission checks if user has permission for environment operations
func (s *Service) checkEnvironmentPermission(userID int64, permission string) error {
	// Skip permission check if rbacService is nil (for unit tests)
	if s.rbacService == nil {
		return nil
	}

	hasPermission, err := s.rbacService.HasPermission(userID, constants.ModuleAPIDocEnvironments, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// checkExportPermission checks if user has permission for export operations
func (s *Service) checkExportPermission(userID int64, permission string) error {
	// Skip permission check if rbacService is nil (for unit tests)
	if s.rbacService == nil {
		return nil
	}

	hasPermission, err := s.rbacService.HasPermission(userID, constants.ModuleAPIDocExport, permission)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf(constants.MsgInsufficientPermissions)
	}
	return nil
}

// checkResourceOwnership checks if user has access to a collection (through company isolation)
func (s *Service) checkResourceOwnership(userID, collectionID int64) error {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}

	// Verify collection belongs to user's company
	collection, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return fmt.Errorf("access denied or collection not found: %w", err)
	}

	// Additional check: verify collection actually belongs to the company
	if collection.CompanyID != companyID {
		return fmt.Errorf("access denied: collection does not belong to your company")
	}

	return nil
}

// Wrapper methods that automatically get company ID from user ID

// CreateCollectionByUser creates a new collection using user ID to get company ID
func (s *Service) CreateCollectionByUser(req *CreateCollectionRequest, userID int64) (*CollectionResponse, error) {
	// Check user permissions for collection creation
	if err := s.checkCollectionPermission(userID, "write"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.CreateCollection(req, userID, companyID)
}

// GetCollectionsByUser retrieves collections using user ID to get company ID
func (s *Service) GetCollectionsByUser(req *CollectionListRequest, userID int64) (*CollectionListResponse, error) {
	// Check user permissions for collection reading
	if err := s.checkCollectionPermission(userID, "read"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetCollections(req, companyID)
}

// GetCollectionByIDByUser retrieves a collection using user ID to get company ID
func (s *Service) GetCollectionByIDByUser(id, userID int64) (*CollectionResponse, error) {
	// Check user permissions for collection reading
	if err := s.checkCollectionPermission(userID, "read"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetCollectionByID(id, companyID)
}

// GetCollectionWithStatsByUser retrieves a collection with stats using user ID to get company ID
func (s *Service) GetCollectionWithStatsByUser(id, userID int64) (*CollectionWithStatsResponse, error) {
	// Check user permissions for collection reading
	if err := s.checkCollectionPermission(userID, "read"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetCollectionWithStats(id, companyID)
}

// UpdateCollectionByUser updates a collection using user ID to get company ID
func (s *Service) UpdateCollectionByUser(id int64, req *UpdateCollectionRequest, userID int64) (*CollectionResponse, error) {
	// Check user permissions for collection updates
	if err := s.checkCollectionPermission(userID, "write"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateCollection(id, req, userID, companyID)
}

// DeleteCollectionByUser deletes a collection using user ID to get company ID
func (s *Service) DeleteCollectionByUser(id, userID int64) error {
	// Check user permissions for collection deletion
	if err := s.checkCollectionPermission(userID, "delete"); err != nil {
		return err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.DeleteCollection(id, companyID)
}

// CreateFolderByUser creates a folder using user ID to get company ID
func (s *Service) CreateFolderByUser(collectionID int64, req *CreateFolderRequest, userID int64) (*FolderResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.CreateFolder(collectionID, req, companyID)
}

// GetFoldersByUser retrieves folders using user ID to get company ID
func (s *Service) GetFoldersByUser(collectionID, userID int64) ([]*FolderWithChildrenResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetFolders(collectionID, companyID)
}

// GetFolderByIDByUser retrieves a folder using user ID to get company ID
func (s *Service) GetFolderByIDByUser(id, userID int64) (*FolderResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetFolderByID(id, companyID)
}

// UpdateFolderByUser updates a folder using user ID to get company ID
func (s *Service) UpdateFolderByUser(id int64, req *UpdateFolderRequest, userID int64) (*FolderResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateFolder(id, req, companyID)
}

// DeleteFolderByUser deletes a folder using user ID to get company ID
func (s *Service) DeleteFolderByUser(id, userID int64) error {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.DeleteFolder(id, companyID)
}

// ReorderFoldersByUser reorders folders using user ID to get company ID
func (s *Service) ReorderFoldersByUser(collectionID, userID int64, folderOrders []FolderOrderRequest) error {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.ReorderFolders(collectionID, companyID, folderOrders)
}

// CreateEndpointByUser creates an endpoint using user ID to get company ID
func (s *Service) CreateEndpointByUser(collectionID int64, req *CreateEndpointRequest, userID int64) (*EndpointResponse, error) {
	// Check user permissions for endpoint creation
	if err := s.checkEndpointPermission(userID, "write"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.CreateEndpoint(collectionID, req, companyID)
}

// GetEndpointsByUser retrieves endpoints using user ID to get company ID
func (s *Service) GetEndpointsByUser(collectionID, userID int64, req *EndpointListRequest) (*EndpointListResponse, error) {
	// Check user permissions for endpoint reading
	if err := s.checkEndpointPermission(userID, "read"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEndpoints(collectionID, companyID, req)
}

// GetEndpointByIDByUser retrieves an endpoint using user ID to get company ID
func (s *Service) GetEndpointByIDByUser(id, userID int64) (*EndpointResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEndpointByID(id, companyID)
}

// GetEndpointWithDetailsByUser retrieves an endpoint with details using user ID to get company ID
func (s *Service) GetEndpointWithDetailsByUser(id, userID int64) (*EndpointWithDetailsResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEndpointWithDetails(id, companyID)
}

// UpdateEndpointByUser updates an endpoint using user ID to get company ID
func (s *Service) UpdateEndpointByUser(id int64, req *UpdateEndpointRequest, userID int64) (*EndpointResponse, error) {
	// Check user permissions for endpoint updates
	if err := s.checkEndpointPermission(userID, "write"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateEndpoint(id, req, companyID)
}

// DeleteEndpointByUser deletes an endpoint using user ID to get company ID
func (s *Service) DeleteEndpointByUser(id, userID int64) error {
	// Check user permissions for endpoint deletion
	if err := s.checkEndpointPermission(userID, "delete"); err != nil {
		return err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.DeleteEndpoint(id, companyID)
}

// BulkCreateEndpointsByUser creates multiple endpoints using user ID to get company ID
func (s *Service) BulkCreateEndpointsByUser(collectionID, userID int64, requests []*CreateEndpointRequest) ([]*EndpointResponse, error) {
	// Check user permissions for bulk endpoint creation
	if err := s.checkBulkOperationPermission(userID, "write"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.BulkCreateEndpoints(collectionID, companyID, requests)
}

// BulkUpdateEndpointsByUser updates multiple endpoints using user ID to get company ID
func (s *Service) BulkUpdateEndpointsByUser(userID int64, updates []EndpointBulkUpdateRequest) ([]*EndpointResponse, error) {
	// Check user permissions for bulk endpoint updates
	if err := s.checkBulkOperationPermission(userID, "write"); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.BulkUpdateEndpoints(companyID, updates)
}

// BulkDeleteEndpointsByUser deletes multiple endpoints using user ID to get company ID
func (s *Service) BulkDeleteEndpointsByUser(userID int64, endpointIDs []int64) error {
	// Check user permissions for bulk endpoint deletion
	if err := s.checkBulkOperationPermission(userID, "delete"); err != nil {
		return err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.BulkDeleteEndpoints(companyID, endpointIDs)
}

// BulkMoveEndpointsByUser moves multiple endpoints using user ID to get company ID
func (s *Service) BulkMoveEndpointsByUser(userID int64, endpointIDs []int64, targetFolderID *int64) error {
	// Check user permissions for bulk endpoint moves
	if err := s.checkBulkOperationPermission(userID, "write"); err != nil {
		return err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.BulkMoveEndpoints(companyID, endpointIDs, targetFolderID)
}

// CreateEnvironmentByUser creates an environment using user ID to get company ID
func (s *Service) CreateEnvironmentByUser(collectionID, userID int64, req *CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Check user permissions for environment creation
	if err := s.checkEnvironmentPermission(userID, "write"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.CreateEnvironment(collectionID, companyID, req)
}

// GetEnvironmentsByUser retrieves environments using user ID to get company ID
func (s *Service) GetEnvironmentsByUser(collectionID, userID int64) ([]*EnvironmentResponse, error) {
	// Check user permissions for environment reading
	if err := s.checkEnvironmentPermission(userID, "read"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEnvironments(collectionID, companyID)
}

// GetEnvironmentByIDByUser retrieves an environment using user ID to get company ID
func (s *Service) GetEnvironmentByIDByUser(id, userID int64) (*EnvironmentResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEnvironmentByID(id, companyID)
}

// GetEnvironmentWithVariablesByUser retrieves an environment with variables using user ID to get company ID
func (s *Service) GetEnvironmentWithVariablesByUser(id, userID int64) (*EnvironmentWithVariablesResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEnvironmentWithVariables(id, companyID)
}

// UpdateEnvironmentByUser updates an environment using user ID to get company ID
func (s *Service) UpdateEnvironmentByUser(id int64, req *UpdateEnvironmentRequest, userID int64) (*EnvironmentResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateEnvironment(id, req, companyID)
}

// DeleteEnvironmentByUser deletes an environment using user ID to get company ID
func (s *Service) DeleteEnvironmentByUser(id, userID int64) error {
	// Check user permissions for environment deletion
	if err := s.checkEnvironmentPermission(userID, "delete"); err != nil {
		return err
	}

	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.DeleteEnvironment(id, companyID)
}

// CreateEnvironmentVariableByUser creates an environment variable using user ID to get company ID
func (s *Service) CreateEnvironmentVariableByUser(environmentID, userID int64, req *CreateEnvironmentVariableRequest) (*EnvironmentVariableResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.CreateEnvironmentVariable(environmentID, companyID, req)
}

// GetEnvironmentVariablesByUser retrieves environment variables using user ID to get company ID
func (s *Service) GetEnvironmentVariablesByUser(environmentID, userID int64) ([]*EnvironmentVariableResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.GetEnvironmentVariables(environmentID, companyID)
}

// UpdateEnvironmentVariableByUser updates an environment variable using user ID to get company ID
func (s *Service) UpdateEnvironmentVariableByUser(id int64, req *UpdateEnvironmentVariableRequest, userID int64) (*EnvironmentVariableResponse, error) {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}
	return s.UpdateEnvironmentVariable(id, req, companyID)
}

// DeleteEnvironmentVariableByUser deletes an environment variable using user ID to get company ID
func (s *Service) DeleteEnvironmentVariableByUser(id, userID int64) error {
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return err
	}
	return s.DeleteEnvironmentVariable(id, companyID)
}

// Export-related methods with RBAC

// ExportCollectionByUser exports a collection in specified format with permission checks
func (s *Service) ExportCollectionByUser(collectionID, userID int64, format string, environmentID *int64) (*export.ExportResult, error) {
	// Check user permissions for export operations
	if err := s.checkExportPermission(userID, "read"); err != nil {
		return nil, err
	}

	// Check resource ownership
	if err := s.checkResourceOwnership(userID, collectionID); err != nil {
		return nil, err
	}

	// Get user's company ID
	companyID, err := s.getUserCompanyID(userID)
	if err != nil {
		return nil, err
	}

	// Get collection with full details
	collectionWithDetails, err := s.getCollectionWithFullDetails(collectionID, companyID, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection details: %w", err)
	}

	// Validate export format
	exportFormat := export.ExportFormat(format)
	supportedFormats := s.exportManager.GetSupportedFormats()
	supported := false
	for _, supportedFormat := range supportedFormats {
		if exportFormat == supportedFormat {
			supported = true
			break
		}
	}
	if !supported {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}

	// Create export options
	options := export.DefaultExportOptions()
	options.Format = exportFormat
	options.EnvironmentID = environmentID

	// Perform export
	result, err := s.exportManager.Export(collectionWithDetails, exportFormat, options)
	if err != nil {
		return nil, fmt.Errorf("export failed: %w", err)
	}

	return result, nil
}

// Helper method to check if user can perform bulk operations
func (s *Service) checkBulkOperationPermission(userID int64, operation string) error {
	// Bulk operations require write permissions on endpoints
	return s.checkEndpointPermission(userID, operation)
}

// getCollectionWithFullDetails retrieves collection with all related data for export
func (s *Service) getCollectionWithFullDetails(collectionID, companyID int64, environmentID *int64) (*export.CollectionWithDetails, error) {
	// Get collection
	collection, err := s.repo.GetCollectionByID(collectionID, companyID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}

	// Get folders
	folders, err := s.repo.GetFoldersHierarchy(collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders: %w", err)
	}

	// Get endpoints with details
	endpointsWithDetails := []export.EndpointWithDetails{}
	endpoints, _, err := s.repo.GetEndpointsByCollectionID(collectionID, 1000, 0) // Get all endpoints
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %w", err)
	}

	for _, endpoint := range endpoints {
		endpointDetails, err := s.repo.GetEndpointWithDetails(endpoint.ID)
		if err != nil {
			continue // Skip endpoints that can't be loaded
		}

		// Convert to export types
		exportEndpoint := s.convertEndpointToExportType(endpointDetails)
		endpointsWithDetails = append(endpointsWithDetails, exportEndpoint)
	}

	// Get environment if specified
	var environment *export.EnvironmentWithVariables
	if environmentID != nil {
		env, err := s.repo.GetEnvironmentWithVariables(*environmentID)
		if err != nil {
			return nil, fmt.Errorf("environment not found: %w", err)
		}

		// Verify environment belongs to collection
		if env.CollectionID != collectionID {
			return nil, fmt.Errorf("environment does not belong to collection")
		}

		environment = s.convertEnvironmentToExportType(env)
	}

	// Convert collection to export type
	exportCollection := s.convertCollectionToExportType(collection)
	exportFolders := s.convertFoldersToExportType(folders)

	return &export.CollectionWithDetails{
		Collection:  exportCollection,
		Folders:     exportFolders,
		Endpoints:   endpointsWithDetails,
		Environment: environment,
	}, nil
}

// Helper methods to convert between apidoc types and export types

func (s *Service) convertCollectionToExportType(collection *Collection) export.Collection {
	return export.Collection{
		ID:            collection.ID,
		Name:          collection.Name,
		Description:   collection.Description.String,
		Version:       collection.Version,
		BaseURL:       collection.BaseURL.String,
		SchemaVersion: collection.SchemaVersion,
		CreatedBy:     collection.CreatedBy,
		CompanyID:     collection.CompanyID,
		IsActive:      collection.IsActive,
		CreatedAt:     collection.CreatedAt,
		UpdatedAt:     collection.UpdatedAt,
	}
}

func (s *Service) convertFoldersToExportType(folders []*FolderWithChildren) []export.Folder {
	var exportFolders []export.Folder
	for _, folder := range folders {
		exportFolder := export.Folder{
			ID:           folder.ID,
			CollectionID: folder.CollectionID,
			ParentID:     nil,
			Name:         folder.Name,
			Description:  folder.Description.String,
			SortOrder:    folder.SortOrder,
			CreatedAt:    folder.CreatedAt,
			UpdatedAt:    folder.UpdatedAt,
		}
		if folder.ParentID.Valid {
			exportFolder.ParentID = &folder.ParentID.Int64
		}
		exportFolders = append(exportFolders, exportFolder)
	}
	return exportFolders
}

func (s *Service) convertEndpointToExportType(endpoint *EndpointWithDetails) export.EndpointWithDetails {
	exportEndpoint := export.EndpointWithDetails{
		Endpoint: export.Endpoint{
			ID:           endpoint.ID,
			CollectionID: endpoint.CollectionID,
			FolderID:     nil,
			Name:         endpoint.Name,
			Description:  endpoint.Description.String,
			Method:       endpoint.Method,
			URL:          endpoint.URL,
			SortOrder:    endpoint.SortOrder,
			IsActive:     endpoint.IsActive,
			CreatedAt:    endpoint.CreatedAt,
			UpdatedAt:    endpoint.UpdatedAt,
		},
	}

	if endpoint.FolderID.Valid {
		exportEndpoint.FolderID = &endpoint.FolderID.Int64
	}

	// Convert headers
	for _, header := range endpoint.Headers {
		exportHeader := export.Header{
			ID:          header.ID,
			EndpointID:  header.EndpointID,
			KeyName:     header.KeyName,
			Value:       header.Value.String,
			Description: header.Description.String,
			IsRequired:  header.IsRequired,
			HeaderType:  header.HeaderType,
			CreatedAt:   header.CreatedAt,
		}
		exportEndpoint.Headers = append(exportEndpoint.Headers, exportHeader)
	}

	// Convert parameters
	for _, param := range endpoint.Parameters {
		exportParam := export.Parameter{
			ID:           param.ID,
			EndpointID:   param.EndpointID,
			Name:         param.Name,
			Type:         param.Type,
			DataType:     param.DataType,
			Description:  param.Description.String,
			DefaultValue: param.DefaultValue.String,
			ExampleValue: param.ExampleValue.String,
			IsRequired:   param.IsRequired,
			CreatedAt:    param.CreatedAt,
		}
		exportEndpoint.Parameters = append(exportEndpoint.Parameters, exportParam)
	}

	// Convert request body
	if endpoint.RequestBody != nil {
		// Convert JSONB to string
		var schemaDefinition string
		if endpoint.RequestBody.SchemaDefinition != nil {
			// Convert JSONB to JSON string
			schemaBytes, err := json.Marshal(endpoint.RequestBody.SchemaDefinition)
			if err == nil {
				schemaDefinition = string(schemaBytes)
			}
		}

		exportRequestBody := &export.RequestBody{
			ID:               endpoint.RequestBody.ID,
			EndpointID:       endpoint.RequestBody.EndpointID,
			ContentType:      endpoint.RequestBody.ContentType,
			BodyContent:      endpoint.RequestBody.BodyContent.String,
			Description:      endpoint.RequestBody.Description.String,
			SchemaDefinition: schemaDefinition,
			CreatedAt:        endpoint.RequestBody.CreatedAt,
			UpdatedAt:        endpoint.RequestBody.UpdatedAt,
		}
		exportEndpoint.RequestBody = exportRequestBody
	}

	// Convert responses
	for _, response := range endpoint.Responses {
		exportResponse := export.Response{
			ID:           response.ID,
			EndpointID:   response.EndpointID,
			StatusCode:   response.StatusCode,
			StatusText:   response.StatusText.String,
			ContentType:  response.ContentType,
			ResponseBody: response.ResponseBody.String,
			Description:  response.Description.String,
			IsDefault:    response.IsDefault,
			CreatedAt:    response.CreatedAt,
		}
		exportEndpoint.Responses = append(exportEndpoint.Responses, exportResponse)
	}

	// Convert tests
	for _, test := range endpoint.Tests {
		exportTest := export.Test{
			ID:            test.ID,
			EndpointID:    test.EndpointID,
			TestType:      test.TestType,
			ScriptContent: test.ScriptContent.String,
			Description:   test.Description.String,
			CreatedAt:     test.CreatedAt,
			UpdatedAt:     test.UpdatedAt,
		}
		exportEndpoint.Tests = append(exportEndpoint.Tests, exportTest)
	}

	return exportEndpoint
}

func (s *Service) convertEnvironmentToExportType(env *EnvironmentWithVariables) *export.EnvironmentWithVariables {
	exportEnv := &export.EnvironmentWithVariables{
		Environment: export.Environment{
			ID:           env.ID,
			CollectionID: env.CollectionID,
			Name:         env.Name,
			Description:  env.Description.String,
			IsDefault:    env.IsDefault,
			CreatedAt:    env.CreatedAt,
			UpdatedAt:    env.UpdatedAt,
		},
	}

	// Convert variables
	for _, variable := range env.Variables {
		exportVar := export.EnvironmentVariable{
			ID:            variable.ID,
			EnvironmentID: variable.EnvironmentID,
			KeyName:       variable.KeyName,
			Value:         variable.Value.String,
			Description:   variable.Description.String,
			IsSecret:      variable.IsSecret,
			CreatedAt:     variable.CreatedAt,
		}
		exportEnv.Variables = append(exportEnv.Variables, exportVar)
	}

	return exportEnv
}

// GetSupportedExportFormats returns list of supported export formats
func (s *Service) GetSupportedExportFormats() []string {
	formats := s.exportManager.GetSupportedFormats()
	result := make([]string, len(formats))
	for i, format := range formats {
		result[i] = string(format)
	}
	return result
}

// GetExportCacheStats returns export cache statistics
func (s *Service) GetExportCacheStats() (*export.CacheStats, error) {
	return s.exportManager.GetCacheStats()
}

// GetExportPerformanceStats returns export performance statistics
func (s *Service) GetExportPerformanceStats() *export.PerformanceStats {
	return s.exportManager.GetPerformanceStats()
}

// GetExportProgress returns progress for a specific export operation
func (s *Service) GetExportProgress(progressID string) (*export.ExportProgress, bool) {
	return s.exportManager.GetExportProgress(progressID)
}

// GetAllExportProgress returns all active export progress
func (s *Service) GetAllExportProgress() map[string]*export.ExportProgress {
	return s.exportManager.GetAllExportProgress()
}

// InvalidateExportCache invalidates export cache for a collection
func (s *Service) InvalidateExportCache(collectionID int64) error {
	return s.exportManager.InvalidateCache(collectionID)
}

// InvalidateExportCacheFormat invalidates export cache for a collection and format
func (s *Service) InvalidateExportCacheFormat(collectionID int64, format string) error {
	exportFormat := export.ExportFormat(format)
	return s.exportManager.InvalidateCacheFormat(collectionID, exportFormat)
}
