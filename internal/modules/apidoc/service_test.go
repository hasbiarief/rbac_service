package apidoc

import (
	"database/sql"
	"errors"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/rbac"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

// MockRBACService is a mock implementation of the RBAC service for testing
type MockRBACService struct {
	mock.Mock
}

func (m *MockRBACService) GetUserPermissions(userID int64) (*rbac.UserPermissions, error) {
	args := m.Called(userID)
	return args.Get(0).(*rbac.UserPermissions), args.Error(1)
}

func (m *MockRBACService) HasPermission(userID int64, moduleID int64, permission string) (bool, error) {
	args := m.Called(userID, moduleID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) HasRole(userID int64, roleName string) (bool, error) {
	args := m.Called(userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) GetAccessibleModules(userID int64, permission string) ([]int64, error) {
	args := m.Called(userID, permission)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *MockRBACService) IsSuperAdmin(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRBACService) GetFilteredModules(userID int64, permission string, limit, offset int, category, search string, isActive *bool) ([]*rbac.ModuleInfo, error) {
	args := m.Called(userID, permission, limit, offset, category, search, isActive)
	return args.Get(0).([]*rbac.ModuleInfo), args.Error(1)
}

// Helper function to set up common RBAC mocks for tests
func setupRBACMocks(mockRBAC *MockRBACService, userID int64, hasPermission bool) {
	// Collection permissions
	mockRBAC.On("HasPermission", userID, int64(140), "read").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(140), "write").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(140), "delete").Return(hasPermission, nil).Maybe()

	// Endpoint permissions
	mockRBAC.On("HasPermission", userID, int64(141), "read").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(141), "write").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(141), "delete").Return(hasPermission, nil).Maybe()

	// Environment permissions
	mockRBAC.On("HasPermission", userID, int64(142), "read").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(142), "write").Return(hasPermission, nil).Maybe()
	mockRBAC.On("HasPermission", userID, int64(142), "delete").Return(hasPermission, nil).Maybe()

	// Export permissions
	mockRBAC.On("HasPermission", userID, int64(143), "read").Return(hasPermission, nil).Maybe()
}

func (m *MockRepository) CreateCollection(collection *Collection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *MockRepository) GetCollectionByID(id int64, companyID int64) (*Collection, error) {
	args := m.Called(id, companyID)
	return args.Get(0).(*Collection), args.Error(1)
}

func (m *MockRepository) GetCollections(companyID int64, limit, offset int) ([]*Collection, int64, error) {
	args := m.Called(companyID, limit, offset)
	return args.Get(0).([]*Collection), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetCollectionWithStats(id int64, companyID int64) (*CollectionWithStats, error) {
	args := m.Called(id, companyID)
	return args.Get(0).(*CollectionWithStats), args.Error(1)
}

func (m *MockRepository) UpdateCollection(collection *Collection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *MockRepository) DeleteCollection(id int64, companyID int64) error {
	args := m.Called(id, companyID)
	return args.Error(0)
}

func (m *MockRepository) CreateFolder(folder *Folder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockRepository) GetFoldersByCollectionID(collectionID int64) ([]*Folder, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*Folder), args.Error(1)
}

func (m *MockRepository) GetFolderByID(id int64) (*Folder, error) {
	args := m.Called(id)
	return args.Get(0).(*Folder), args.Error(1)
}

func (m *MockRepository) GetFoldersHierarchy(collectionID int64) ([]*FolderWithChildren, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*FolderWithChildren), args.Error(1)
}

func (m *MockRepository) UpdateFolder(folder *Folder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *MockRepository) DeleteFolder(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateEndpoint(endpoint *Endpoint) error {
	args := m.Called(endpoint)
	return args.Error(0)
}

func (m *MockRepository) GetEndpointsByCollectionID(collectionID int64, limit, offset int) ([]*Endpoint, int64, error) {
	args := m.Called(collectionID, limit, offset)
	return args.Get(0).([]*Endpoint), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) GetEndpointByID(id int64) (*Endpoint, error) {
	args := m.Called(id)
	return args.Get(0).(*Endpoint), args.Error(1)
}

func (m *MockRepository) GetEndpointWithDetails(id int64) (*EndpointWithDetails, error) {
	args := m.Called(id)
	return args.Get(0).(*EndpointWithDetails), args.Error(1)
}

func (m *MockRepository) UpdateEndpoint(endpoint *Endpoint) error {
	args := m.Called(endpoint)
	return args.Error(0)
}

func (m *MockRepository) DeleteEndpoint(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateEnvironment(environment *Environment) error {
	args := m.Called(environment)
	return args.Error(0)
}

func (m *MockRepository) GetEnvironmentsByCollectionID(collectionID int64) ([]*Environment, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*Environment), args.Error(1)
}

func (m *MockRepository) GetEnvironmentByID(id int64) (*Environment, error) {
	args := m.Called(id)
	return args.Get(0).(*Environment), args.Error(1)
}

func (m *MockRepository) GetEnvironmentWithVariables(id int64) (*EnvironmentWithVariables, error) {
	args := m.Called(id)
	return args.Get(0).(*EnvironmentWithVariables), args.Error(1)
}

func (m *MockRepository) UpdateEnvironment(environment *Environment) error {
	args := m.Called(environment)
	return args.Error(0)
}

func (m *MockRepository) DeleteEnvironment(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateEnvironmentVariable(variable *EnvironmentVariable) error {
	args := m.Called(variable)
	return args.Error(0)
}

func (m *MockRepository) GetEnvironmentVariables(environmentID int64) ([]*EnvironmentVariable, error) {
	args := m.Called(environmentID)
	return args.Get(0).([]*EnvironmentVariable), args.Error(1)
}

func (m *MockRepository) GetEnvironmentVariableByID(id int64) (*EnvironmentVariable, error) {
	args := m.Called(id)
	return args.Get(0).(*EnvironmentVariable), args.Error(1)
}

func (m *MockRepository) UpdateEnvironmentVariable(variable *EnvironmentVariable) error {
	args := m.Called(variable)
	return args.Error(0)
}

func (m *MockRepository) DeleteEnvironmentVariable(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateHeader(header *Header) error {
	args := m.Called(header)
	return args.Error(0)
}

func (m *MockRepository) GetHeadersByEndpointID(endpointID int64) ([]*Header, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Header), args.Error(1)
}

func (m *MockRepository) UpdateHeader(header *Header) error {
	args := m.Called(header)
	return args.Error(0)
}

func (m *MockRepository) DeleteHeader(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateParameter(parameter *Parameter) error {
	args := m.Called(parameter)
	return args.Error(0)
}

func (m *MockRepository) GetParametersByEndpointID(endpointID int64) ([]*Parameter, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Parameter), args.Error(1)
}

func (m *MockRepository) UpdateParameter(parameter *Parameter) error {
	args := m.Called(parameter)
	return args.Error(0)
}

func (m *MockRepository) DeleteParameter(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateRequestBody(requestBody *RequestBody) error {
	args := m.Called(requestBody)
	return args.Error(0)
}

func (m *MockRepository) GetRequestBodyByEndpointID(endpointID int64) (*RequestBody, error) {
	args := m.Called(endpointID)
	return args.Get(0).(*RequestBody), args.Error(1)
}

func (m *MockRepository) UpdateRequestBody(requestBody *RequestBody) error {
	args := m.Called(requestBody)
	return args.Error(0)
}

func (m *MockRepository) DeleteRequestBody(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) CreateResponse(response *Response) error {
	args := m.Called(response)
	return args.Error(0)
}

func (m *MockRepository) GetResponsesByEndpointID(endpointID int64) ([]*Response, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Response), args.Error(1)
}

func (m *MockRepository) UpdateResponse(response *Response) error {
	args := m.Called(response)
	return args.Error(0)
}

func (m *MockRepository) DeleteResponse(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) GetAllTags() ([]*Tag, error) {
	args := m.Called()
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *MockRepository) CreateTag(tag *Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *MockRepository) GetTagsByEndpointID(endpointID int64) ([]*Tag, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *MockRepository) AddTagToEndpoint(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *MockRepository) RemoveTagFromEndpoint(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *MockRepository) GetCollectionForExport(collectionID int64, companyID int64) (*CollectionExport, error) {
	args := m.Called(collectionID, companyID)
	return args.Get(0).(*CollectionExport), args.Error(1)
}

// Test helper functions
func createTestCollection() *Collection {
	return &Collection{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:          "Test Collection",
		Description:   model.NullString{NullString: sql.NullString{String: "Test Description", Valid: true}},
		Version:       "1.0.0",
		BaseURL:       model.NullString{NullString: sql.NullString{String: "https://api.example.com", Valid: true}},
		SchemaVersion: "1.0",
		CreatedBy:     1,
		CompanyID:     1,
		IsActive:      true,
	}
}

func createTestFolder() *Folder {
	return &Folder{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CollectionID: 1,
		ParentID:     model.NullInt64{NullInt64: sql.NullInt64{Valid: false}},
		Name:         "Test Folder",
		Description:  model.NullString{NullString: sql.NullString{String: "Test Description", Valid: true}},
		SortOrder:    0,
	}
}

func createTestEndpoint() *Endpoint {
	return &Endpoint{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CollectionID: 1,
		FolderID:     model.NullInt64{NullInt64: sql.NullInt64{Valid: false}},
		Name:         "Test Endpoint",
		Description:  model.NullString{NullString: sql.NullString{String: "Test Description", Valid: true}},
		Method:       "GET",
		URL:          "/api/test",
		SortOrder:    0,
		IsActive:     true,
	}
}

func createTestEnvironment() *Environment {
	return &Environment{
		BaseModel: model.BaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CollectionID: 1,
		Name:         "Test Environment",
		Description:  model.NullString{NullString: sql.NullString{String: "Test Description", Valid: true}},
		IsDefault:    true,
	}
}

func createTestEnvironmentVariable() *EnvironmentVariable {
	return &EnvironmentVariable{
		ID:            1,
		EnvironmentID: 1,
		KeyName:       "API_KEY",
		Value:         model.NullString{NullString: sql.NullString{String: "test-key", Valid: true}},
		Description:   model.NullString{NullString: sql.NullString{String: "Test API Key", Valid: true}},
		IsSecret:      true,
		CreatedAt:     time.Now(),
	}
}

// Collection Service Tests

func TestService_CreateCollection(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	req := &CreateCollectionRequest{
		Name:          "Test Collection",
		Description:   stringPtr("Test Description"),
		Version:       "1.0.0",
		BaseURL:       stringPtr("https://api.example.com"),
		SchemaVersion: "1.0",
	}

	collection := createTestCollection()

	mockRepo.On("CreateCollection", mock.AnythingOfType("*apidoc.Collection")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Collection)
		arg.ID = collection.ID
		arg.CreatedAt = collection.CreatedAt
		arg.UpdatedAt = collection.UpdatedAt
	})

	result, err := service.CreateCollection(req, 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Collection", result.Name)
	assert.Equal(t, "Test Description", *result.Description)
	assert.Equal(t, "1.0.0", result.Version)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateCollection_ValidationError(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	req := &CreateCollectionRequest{
		Name:          "", // Invalid: empty name
		Version:       "1.0.0",
		SchemaVersion: "1.0",
	}

	result, err := service.CreateCollection(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestService_CreateCollection_InsufficientPermissions(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRBAC := new(MockRBACService)
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	req := &CreateCollectionRequest{
		Name:          "Test Collection",
		Description:   stringPtr("Test Description"),
		Version:       "1.0.0",
		BaseURL:       stringPtr("https://api.example.com"),
		SchemaVersion: "1.0",
	}

	// Mock RBAC permission check - user doesn't have write permission
	mockRBAC.On("HasPermission", int64(1), int64(140), "write").Return(false, nil)

	result, err := service.CreateCollection(req, 1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Insufficient permissions")
	mockRBAC.AssertExpectations(t)
}

func TestService_GetCollections(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	req := &CollectionListRequest{
		Limit:  10,
		Offset: 0,
	}

	collections := []*Collection{createTestCollection()}
	mockRepo.On("GetCollections", int64(1), 10, 0).Return(collections, int64(1), nil)

	result, err := service.GetCollections(req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 0, result.Offset)
	assert.False(t, result.HasMore)
	mockRepo.AssertExpectations(t)
}

func TestService_GetCollectionByID(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)

	result, err := service.GetCollectionByID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Collection", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetCollectionByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return((*Collection)(nil), errors.New("collection not found"))

	result, err := service.GetCollectionByID(1, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get collection")
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateCollection(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	req := &UpdateCollectionRequest{
		Name:        stringPtr("Updated Collection"),
		Description: stringPtr("Updated Description"),
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("UpdateCollection", mock.AnythingOfType("*apidoc.Collection")).Return(nil)

	result, err := service.UpdateCollection(1, req, 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Collection", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteCollection(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("DeleteCollection", int64(1), int64(1)).Return(nil)

	err := service.DeleteCollection(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Folder Service Tests

func TestService_CreateFolder(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	req := &CreateFolderRequest{
		Name:        "Test Folder",
		Description: stringPtr("Test Description"),
		SortOrder:   intPtr(0),
	}

	folder := createTestFolder()
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("CreateFolder", mock.AnythingOfType("*apidoc.Folder")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Folder)
		arg.ID = folder.ID
		arg.CreatedAt = folder.CreatedAt
		arg.UpdatedAt = folder.UpdatedAt
	})

	result, err := service.CreateFolder(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Folder", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateFolder_WithParent(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	parentFolder := createTestFolder()
	parentFolder.ID = 2
	req := &CreateFolderRequest{
		ParentID:    int64Ptr(2),
		Name:        "Child Folder",
		Description: stringPtr("Child Description"),
	}

	folder := createTestFolder()
	folder.Name = "Child Folder"
	folder.ParentID = model.NullInt64{NullInt64: sql.NullInt64{Int64: 2, Valid: true}}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetFolderByID", int64(2)).Return(parentFolder, nil)
	mockRepo.On("CreateFolder", mock.AnythingOfType("*apidoc.Folder")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Folder)
		arg.ID = folder.ID
		arg.CreatedAt = folder.CreatedAt
		arg.UpdatedAt = folder.UpdatedAt
	})

	result, err := service.CreateFolder(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Child Folder", result.Name)
	assert.Equal(t, int64(2), *result.ParentID)
	mockRepo.AssertExpectations(t)
}

func TestService_GetFolders(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder := createTestFolder()
	folderWithChildren := &FolderWithChildren{
		Folder:   *folder,
		Children: []*FolderWithChildren{},
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetFoldersHierarchy", int64(1)).Return([]*FolderWithChildren{folderWithChildren}, nil)

	result, err := service.GetFolders(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Folder", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateFolder(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder := createTestFolder()
	req := &UpdateFolderRequest{
		Name:        stringPtr("Updated Folder"),
		Description: stringPtr("Updated Description"),
	}

	mockRepo.On("GetFolderByID", int64(1)).Return(folder, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("UpdateFolder", mock.AnythingOfType("*apidoc.Folder")).Return(nil)

	result, err := service.UpdateFolder(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Folder", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteFolder(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder := createTestFolder()

	mockRepo.On("GetFolderByID", int64(1)).Return(folder, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("DeleteFolder", int64(1)).Return(nil)

	err := service.DeleteFolder(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_ReorderFolders(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder1 := createTestFolder()
	folder1.ID = 1
	folder2 := createTestFolder()
	folder2.ID = 2

	folderOrders := []FolderOrderRequest{
		{FolderID: 1, SortOrder: 1},
		{FolderID: 2, SortOrder: 0},
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetFolderByID", int64(1)).Return(folder1, nil).Twice()
	mockRepo.On("GetFolderByID", int64(2)).Return(folder2, nil).Twice()
	mockRepo.On("UpdateFolder", mock.AnythingOfType("*apidoc.Folder")).Return(nil).Twice()

	err := service.ReorderFolders(1, 1, folderOrders)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_ReorderFolders_InvalidFolder(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder := createTestFolder()
	folder.CollectionID = 2 // Different collection

	folderOrders := []FolderOrderRequest{
		{FolderID: 1, SortOrder: 1},
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetFolderByID", int64(1)).Return(folder, nil)

	err := service.ReorderFolders(1, 1, folderOrders)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong to collection")
	mockRepo.AssertExpectations(t)
}

// Endpoint Service Tests

func TestService_CreateEndpoint(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	req := &CreateEndpointRequest{
		Name:        "Test Endpoint",
		Description: stringPtr("Test Description"),
		Method:      "GET",
		URL:         "/api/test",
		SortOrder:   intPtr(0),
		IsActive:    boolPtr(true),
	}

	endpoint := createTestEndpoint()
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("CreateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Endpoint)
		arg.ID = endpoint.ID
		arg.CreatedAt = endpoint.CreatedAt
		arg.UpdatedAt = endpoint.UpdatedAt
	})

	result, err := service.CreateEndpoint(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Endpoint", result.Name)
	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, "/api/test", result.URL)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEndpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	req := &EndpointListRequest{
		Limit:  10,
		Offset: 0,
	}

	endpoints := []*Endpoint{createTestEndpoint()}
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEndpointsByCollectionID", int64(1), 10, 0).Return(endpoints, int64(1), nil)

	result, err := service.GetEndpoints(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(1), result.Total)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEndpointByID(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint := createTestEndpoint()

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)

	result, err := service.GetEndpointByID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Endpoint", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEndpointWithDetails(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint := createTestEndpoint()
	details := &EndpointWithDetails{
		Endpoint:    *endpoint,
		Headers:     []*Header{},
		Parameters:  []*Parameter{},
		RequestBody: nil,
		Responses:   []*Response{},
		Tags:        []*Tag{},
		Tests:       []*Test{},
	}

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEndpointWithDetails", int64(1)).Return(details, nil)

	result, err := service.GetEndpointWithDetails(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Endpoint", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateEndpoint(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint := createTestEndpoint()
	req := &UpdateEndpointRequest{
		Name:        stringPtr("Updated Endpoint"),
		Description: stringPtr("Updated Description"),
	}

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("UpdateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil)

	result, err := service.UpdateEndpoint(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Endpoint", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEndpoint(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint := createTestEndpoint()

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("DeleteEndpoint", int64(1)).Return(nil)

	err := service.DeleteEndpoint(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_BulkCreateEndpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	requests := []*CreateEndpointRequest{
		{
			Name:        "Endpoint 1",
			Description: stringPtr("Description 1"),
			Method:      "GET",
			URL:         "/api/test1",
			IsActive:    boolPtr(true),
		},
		{
			Name:        "Endpoint 2",
			Description: stringPtr("Description 2"),
			Method:      "POST",
			URL:         "/api/test2",
			IsActive:    boolPtr(true),
		},
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("CreateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil).Twice().Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Endpoint)
		arg.ID = 1
		arg.CreatedAt = time.Now()
		arg.UpdatedAt = time.Now()
	})

	result, err := service.BulkCreateEndpoints(1, 1, requests)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "Endpoint 1", result[0].Name)
	assert.Equal(t, "Endpoint 2", result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestService_BulkUpdateEndpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint1 := createTestEndpoint()
	endpoint1.ID = 1
	endpoint2 := createTestEndpoint()
	endpoint2.ID = 2

	updates := []EndpointBulkUpdateRequest{
		{
			EndpointID: 1,
			UpdateRequest: UpdateEndpointRequest{
				Name: stringPtr("Updated Endpoint 1"),
			},
		},
		{
			EndpointID: 2,
			UpdateRequest: UpdateEndpointRequest{
				Name: stringPtr("Updated Endpoint 2"),
			},
		},
	}

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint1, nil)
	mockRepo.On("GetEndpointByID", int64(2)).Return(endpoint2, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil).Twice()
	mockRepo.On("UpdateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil).Twice()

	result, err := service.BulkUpdateEndpoints(1, updates)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestService_BulkDeleteEndpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	endpoint1 := createTestEndpoint()
	endpoint1.ID = 1
	endpoint2 := createTestEndpoint()
	endpoint2.ID = 2

	endpointIDs := []int64{1, 2}

	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint1, nil)
	mockRepo.On("GetEndpointByID", int64(2)).Return(endpoint2, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil).Twice()
	mockRepo.On("DeleteEndpoint", int64(1)).Return(nil)
	mockRepo.On("DeleteEndpoint", int64(2)).Return(nil)

	err := service.BulkDeleteEndpoints(1, endpointIDs)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_BulkMoveEndpoints(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	folder := createTestFolder()
	folder.ID = 2
	endpoint1 := createTestEndpoint()
	endpoint1.ID = 1
	endpoint2 := createTestEndpoint()
	endpoint2.ID = 2

	endpointIDs := []int64{1, 2}
	targetFolderID := int64Ptr(2)

	mockRepo.On("GetFolderByID", int64(2)).Return(folder, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil).Times(3)
	mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint1, nil)
	mockRepo.On("GetEndpointByID", int64(2)).Return(endpoint2, nil)
	mockRepo.On("UpdateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil).Twice()

	err := service.BulkMoveEndpoints(1, endpointIDs, targetFolderID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Helper functions
func int64Ptr(i int64) *int64 {
	return &i
}

// Environment Service Tests

func TestService_CreateEnvironment(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	req := &CreateEnvironmentRequest{
		Name:        "Test Environment",
		Description: stringPtr("Test Description"),
		IsDefault:   boolPtr(true),
	}

	environment := createTestEnvironment()
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{}, nil)
	mockRepo.On("CreateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Environment)
		arg.ID = environment.ID
		arg.CreatedAt = environment.CreatedAt
		arg.UpdatedAt = environment.UpdatedAt
	})

	result, err := service.CreateEnvironment(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Environment", result.Name)
	assert.Equal(t, "Test Description", *result.Description)
	assert.True(t, result.IsDefault)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateEnvironment_DuplicateName(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	existingEnv := createTestEnvironment()
	req := &CreateEnvironmentRequest{
		Name:        "Test Environment", // Same name as existing
		Description: stringPtr("Test Description"),
	}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{existingEnv}, nil)

	result, err := service.CreateEnvironment(1, 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestService_CreateEnvironment_SetAsDefault(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	existingEnv := createTestEnvironment()
	existingEnv.Name = "Existing Environment"
	req := &CreateEnvironmentRequest{
		Name:      "New Environment",
		IsDefault: boolPtr(true),
	}

	environment := createTestEnvironment()
	environment.Name = "New Environment"

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{existingEnv}, nil)
	mockRepo.On("UpdateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil) // Unset previous default
	mockRepo.On("CreateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*Environment)
		arg.ID = environment.ID
		arg.CreatedAt = environment.CreatedAt
		arg.UpdatedAt = environment.UpdatedAt
	})

	result, err := service.CreateEnvironment(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Environment", result.Name)
	assert.True(t, result.IsDefault)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEnvironments(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environments := []*Environment{createTestEnvironment()}

	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return(environments, nil)

	result, err := service.GetEnvironments(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Environment", result[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEnvironmentByID(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)

	result, err := service.GetEnvironmentByID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Environment", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_GetEnvironmentWithVariables(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	variable := createTestEnvironmentVariable()
	envWithVars := &EnvironmentWithVariables{
		Environment: *environment,
		Variables:   []*EnvironmentVariable{variable},
	}

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentWithVariables", int64(1)).Return(envWithVars, nil)

	result, err := service.GetEnvironmentWithVariables(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Environment", result.Name)
	assert.Len(t, result.Variables, 1)
	assert.Equal(t, "API_KEY", result.Variables[0].KeyName)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateEnvironment(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	req := &UpdateEnvironmentRequest{
		Name:        stringPtr("Updated Environment"),
		Description: stringPtr("Updated Description"),
	}

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment}, nil)
	mockRepo.On("UpdateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil)

	result, err := service.UpdateEnvironment(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Environment", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateEnvironment_SetAsDefault(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	environment.IsDefault = false
	otherEnv := createTestEnvironment()
	otherEnv.ID = 2
	otherEnv.Name = "Other Environment"
	req := &UpdateEnvironmentRequest{
		IsDefault: boolPtr(true),
	}

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment, otherEnv}, nil)
	mockRepo.On("UpdateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil).Twice()

	result, err := service.UpdateEnvironment(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsDefault)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEnvironment(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	environment.IsDefault = false
	otherEnv := createTestEnvironment()
	otherEnv.ID = 2
	otherEnv.Name = "Other Environment"

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment, otherEnv}, nil)
	mockRepo.On("DeleteEnvironment", int64(1)).Return(nil)

	err := service.DeleteEnvironment(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEnvironment_LastEnvironment(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment}, nil)

	err := service.DeleteEnvironment(1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete the last environment")
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEnvironment_DefaultEnvironment(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	environment.IsDefault = true
	otherEnv := createTestEnvironment()
	otherEnv.ID = 2
	otherEnv.Name = "Other Environment"
	otherEnv.IsDefault = false

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment, otherEnv}, nil)
	mockRepo.On("UpdateEnvironment", mock.AnythingOfType("*apidoc.Environment")).Return(nil) // Set new default
	mockRepo.On("DeleteEnvironment", int64(1)).Return(nil)

	err := service.DeleteEnvironment(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Environment Variable Service Tests

func TestService_CreateEnvironmentVariable(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	req := &CreateEnvironmentVariableRequest{
		KeyName:     "API_KEY",
		Value:       stringPtr("test-key"),
		Description: stringPtr("Test API Key"),
		IsSecret:    boolPtr(true),
	}

	variable := createTestEnvironmentVariable()
	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentVariables", int64(1)).Return([]*EnvironmentVariable{}, nil)
	mockRepo.On("CreateEnvironmentVariable", mock.AnythingOfType("*apidoc.EnvironmentVariable")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*EnvironmentVariable)
		arg.ID = variable.ID
		arg.CreatedAt = variable.CreatedAt
	})

	result, err := service.CreateEnvironmentVariable(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "API_KEY", result.KeyName)
	assert.Equal(t, "test-key", *result.Value)
	assert.True(t, result.IsSecret)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateEnvironmentVariable_DuplicateKey(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	existingVariable := createTestEnvironmentVariable()
	req := &CreateEnvironmentVariableRequest{
		KeyName: "API_KEY", // Same key as existing
		Value:   stringPtr("new-key"),
	}

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentVariables", int64(1)).Return([]*EnvironmentVariable{existingVariable}, nil)

	result, err := service.CreateEnvironmentVariable(1, 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestService_GetEnvironmentVariables(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	variables := []*EnvironmentVariable{createTestEnvironmentVariable()}

	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentVariables", int64(1)).Return(variables, nil)

	result, err := service.GetEnvironmentVariables(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "API_KEY", result[0].KeyName)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateEnvironmentVariable(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	variable := createTestEnvironmentVariable()
	req := &UpdateEnvironmentVariableRequest{
		KeyName:     stringPtr("UPDATED_KEY"),
		Value:       stringPtr("updated-value"),
		Description: stringPtr("Updated Description"),
		IsSecret:    boolPtr(false),
	}

	mockRepo.On("GetEnvironmentVariableByID", int64(1)).Return(variable, nil)
	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentVariables", int64(1)).Return([]*EnvironmentVariable{variable}, nil)
	mockRepo.On("UpdateEnvironmentVariable", mock.AnythingOfType("*apidoc.EnvironmentVariable")).Return(nil)

	result, err := service.UpdateEnvironmentVariable(1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "UPDATED_KEY", result.KeyName)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateEnvironmentVariable_DuplicateKey(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	variable := createTestEnvironmentVariable()
	otherVariable := createTestEnvironmentVariable()
	otherVariable.ID = 2
	otherVariable.KeyName = "OTHER_KEY"
	req := &UpdateEnvironmentVariableRequest{
		KeyName: stringPtr("OTHER_KEY"), // Same as existing variable
	}

	mockRepo.On("GetEnvironmentVariableByID", int64(1)).Return(variable, nil)
	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("GetEnvironmentVariables", int64(1)).Return([]*EnvironmentVariable{variable, otherVariable}, nil)

	result, err := service.UpdateEnvironmentVariable(1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEnvironmentVariable(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()
	environment := createTestEnvironment()
	variable := createTestEnvironmentVariable()

	mockRepo.On("GetEnvironmentVariableByID", int64(1)).Return(variable, nil)
	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
	mockRepo.On("DeleteEnvironmentVariable", int64(1)).Return(nil)

	err := service.DeleteEnvironmentVariable(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteEnvironmentVariable_AccessDenied(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	environment := createTestEnvironment()
	variable := createTestEnvironmentVariable()

	mockRepo.On("GetEnvironmentVariableByID", int64(1)).Return(variable, nil)
	mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
	mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return((*Collection)(nil), errors.New("collection not found"))

	err := service.DeleteEnvironmentVariable(1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
	mockRepo.AssertExpectations(t)
}

// Additional Business Rule and Authorization Tests

func TestService_CreateCollection_ValidationEdgeCases(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	tests := []struct {
		name        string
		req         *CreateCollectionRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Empty name",
			req: &CreateCollectionRequest{
				Name:          "",
				Version:       "1.0.0",
				SchemaVersion: "1.0",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
		{
			name: "Invalid version format",
			req: &CreateCollectionRequest{
				Name:          "Test Collection",
				Version:       "invalid-version",
				SchemaVersion: "1.0",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
		{
			name: "Invalid base URL",
			req: &CreateCollectionRequest{
				Name:          "Test Collection",
				Version:       "1.0.0",
				BaseURL:       stringPtr("not-a-url"),
				SchemaVersion: "1.0",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateCollection(tt.req, 1, 1)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_CreateEndpoint_ValidationEdgeCases(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	collection := createTestCollection()

	tests := []struct {
		name        string
		req         *CreateEndpointRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Invalid HTTP method",
			req: &CreateEndpointRequest{
				Name:   "Test Endpoint",
				Method: "INVALID",
				URL:    "/api/test",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
		{
			name: "Invalid URL format",
			req: &CreateEndpointRequest{
				Name:   "Test Endpoint",
				Method: "GET",
				URL:    "invalid-url",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
		{
			name: "Empty name",
			req: &CreateEndpointRequest{
				Name:   "",
				Method: "GET",
				URL:    "/api/test",
			},
			expectError: true,
			errorMsg:    "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
				mockRepo.On("CreateEndpoint", mock.AnythingOfType("*apidoc.Endpoint")).Return(nil)
			}

			result, err := service.CreateEndpoint(1, tt.req, 1)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_AuthorizationChecks(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	tests := []struct {
		name        string
		testFunc    func() error
		setupMocks  func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Collection access denied",
			testFunc: func() error {
				_, err := service.GetCollectionByID(1, 999) // Wrong company ID
				return err
			},
			setupMocks: func() {
				mockRepo.On("GetCollectionByID", int64(1), int64(999)).Return((*Collection)(nil), errors.New("collection not found"))
			},
			expectError: true,
			errorMsg:    "failed to get collection",
		},
		{
			name: "Folder access through collection denied",
			testFunc: func() error {
				_, err := service.GetFolderByID(1, 999) // Wrong company ID
				return err
			},
			setupMocks: func() {
				folder := createTestFolder()
				mockRepo.On("GetFolderByID", int64(1)).Return(folder, nil)
				mockRepo.On("GetCollectionByID", int64(1), int64(999)).Return((*Collection)(nil), errors.New("collection not found"))
			},
			expectError: true,
			errorMsg:    "access denied",
		},
		{
			name: "Endpoint access through collection denied",
			testFunc: func() error {
				_, err := service.GetEndpointByID(1, 999) // Wrong company ID
				return err
			},
			setupMocks: func() {
				endpoint := createTestEndpoint()
				mockRepo.On("GetEndpointByID", int64(1)).Return(endpoint, nil)
				mockRepo.On("GetCollectionByID", int64(1), int64(999)).Return((*Collection)(nil), errors.New("collection not found"))
			},
			expectError: true,
			errorMsg:    "access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMocks()
			err := tt.testFunc()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_BusinessRuleEnforcement(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	t.Run("Folder parent must belong to same collection", func(t *testing.T) {
		collection := createTestCollection()
		parentFolder := createTestFolder()
		parentFolder.CollectionID = 2 // Different collection

		req := &CreateFolderRequest{
			ParentID: int64Ptr(1),
			Name:     "Child Folder",
		}

		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
		mockRepo.On("GetFolderByID", int64(1)).Return(parentFolder, nil)

		result, err := service.CreateFolder(1, req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not belong to the same collection")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Endpoint folder must belong to same collection", func(t *testing.T) {
		collection := createTestCollection()
		folder := createTestFolder()
		folder.CollectionID = 2 // Different collection

		req := &CreateEndpointRequest{
			FolderID: int64Ptr(1),
			Name:     "Test Endpoint",
			Method:   "GET",
			URL:      "/api/test",
		}

		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
		mockRepo.On("GetFolderByID", int64(1)).Return(folder, nil)

		result, err := service.CreateEndpoint(1, req, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "does not belong to the same collection")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Cannot delete last environment", func(t *testing.T) {
		collection := createTestCollection()
		environment := createTestEnvironment()

		mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
		mockRepo.On("GetEnvironmentsByCollectionID", int64(1)).Return([]*Environment{environment}, nil)

		err := service.DeleteEnvironment(1, 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete the last environment")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Environment variable key uniqueness within environment", func(t *testing.T) {
		collection := createTestCollection()
		environment := createTestEnvironment()
		existingVariable := createTestEnvironmentVariable()

		req := &CreateEnvironmentVariableRequest{
			KeyName: "API_KEY", // Same as existing
			Value:   stringPtr("new-value"),
		}

		mockRepo.On("GetEnvironmentByID", int64(1)).Return(environment, nil)
		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
		mockRepo.On("GetEnvironmentVariables", int64(1)).Return([]*EnvironmentVariable{existingVariable}, nil)

		result, err := service.CreateEnvironmentVariable(1, 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_ErrorHandlingScenarios(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	t.Run("Repository error propagation", func(t *testing.T) {
		// Reset mock for this test
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return((*Collection)(nil), errors.New("database connection failed"))

		result, err := service.GetCollectionByID(1, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get collection")
		assert.Contains(t, err.Error(), "database connection failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Cascade delete handling", func(t *testing.T) {
		// Reset mock for this test
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		collection := createTestCollection()
		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)
		mockRepo.On("DeleteCollection", int64(1), int64(1)).Return(errors.New("foreign key constraint violation"))

		err := service.DeleteCollection(1, 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete collection")
		assert.Contains(t, err.Error(), "foreign key constraint violation")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Bulk operation partial failure", func(t *testing.T) {
		// Reset mock for this test
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		collection := createTestCollection()
		requests := []*CreateEndpointRequest{
			{
				Name:   "", // This will cause validation failure
				Method: "GET",
				URL:    "/api/test1",
			},
		}

		mockRepo.On("GetCollectionByID", int64(1), int64(1)).Return(collection, nil)

		result, err := service.BulkCreateEndpoints(1, 1, requests)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "validation failed for endpoint 1")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CompanyIsolationEnforcement(t *testing.T) {
	mockRepo := new(MockRepository)
	var mockRBAC RBACServiceInterface = nil
	service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests

	t.Run("Collections are isolated by company", func(t *testing.T) {
		// Reset mock for this test
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		req := &CollectionListRequest{Limit: 10, Offset: 0}

		// Company 1 collections
		company1Collections := []*Collection{createTestCollection()}
		mockRepo.On("GetCollections", int64(1), 10, 0).Return(company1Collections, int64(1), nil)

		result1, err1 := service.GetCollections(req, 1)
		assert.NoError(t, err1)
		assert.Len(t, result1.Data, 1)

		// Reset mock for second call
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		// Company 2 collections (different set)
		company2Collections := []*Collection{}
		mockRepo.On("GetCollections", int64(2), 10, 0).Return(company2Collections, int64(0), nil)

		result2, err2 := service.GetCollections(req, 2)
		assert.NoError(t, err2)
		assert.Len(t, result2.Data, 0)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Cross-company access denied", func(t *testing.T) {
		// Reset mock for this test
		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil

		// Try to access company 1's collection with company 2's credentials
		mockRepo.On("GetCollectionByID", int64(1), int64(2)).Return((*Collection)(nil), errors.New("collection not found"))

		result, err := service.GetCollectionByID(1, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
