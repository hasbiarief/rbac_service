package apidoc

import (
	"database/sql"
	"gin-scalable-api/internal/modules/apidoc/export"
	"gin-scalable-api/pkg/model"
	"gin-scalable-api/pkg/rbac"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestPhase5Comprehensive runs comprehensive tests for Phase 5
func TestPhase5Comprehensive(t *testing.T) {
	t.Run("RepositoryComprehensive", func(t *testing.T) {
		testRepositoryComprehensive(t)
	})

	t.Run("ServiceComprehensive", func(t *testing.T) {
		testServiceComprehensive(t)
	})

	t.Run("DTOValidation", func(t *testing.T) {
		testDTOValidation(t)
	})

	t.Run("Integration", func(t *testing.T) {
		testIntegration(t)
	})
}

func testRepositoryComprehensive(t *testing.T) {
	// Skip if no test database available
	db := setupPhase5TestDatabase(t)
	if db == nil {
		t.Skip("Test database not available")
		return
	}
	defer db.Close()

	repo := NewRepository(db)

	// Test Collection CRUD operations
	t.Run("CollectionCRUD", func(t *testing.T) {
		collection := &Collection{
			Name:          "Test Collection",
			Description:   model.NullString{NullString: sql.NullString{String: "Test Description", Valid: true}},
			Version:       "1.0.0",
			BaseURL:       model.NullString{NullString: sql.NullString{String: "https://api.test.com", Valid: true}},
			SchemaVersion: "2.1.0",
			CreatedBy:     1,
			CompanyID:     1,
			IsActive:      true,
		}

		// Test Create
		err := repo.CreateCollection(collection)
		assert.NoError(t, err)
		assert.NotZero(t, collection.ID)

		// Test GetByID
		retrieved, err := repo.GetCollectionByID(collection.ID, collection.CompanyID)
		assert.NoError(t, err)
		assert.Equal(t, collection.Name, retrieved.Name)

		// Test Update
		retrieved.Name = "Updated Collection"
		err = repo.UpdateCollection(retrieved)
		assert.NoError(t, err)

		// Test Delete
		err = repo.DeleteCollection(collection.ID, collection.CompanyID)
		assert.NoError(t, err)
	})
}

func testServiceComprehensive(t *testing.T) {
	mockRepo := &Phase5MockRepository{}
	mockRBAC := &Phase5MockRBACService{}

	service := &Service{
		repo:          mockRepo,
		rbacService:   mockRBAC,
		exportManager: export.NewExportManager(),
	}

	// Test GetSupportedExportFormats
	formats := service.GetSupportedExportFormats()
	assert.NotEmpty(t, formats)
	assert.Contains(t, formats, "postman")
	assert.Contains(t, formats, "openapi")
	assert.Contains(t, formats, "insomnia")
	assert.Contains(t, formats, "swagger")
	assert.Contains(t, formats, "apidog")
}

func testDTOValidation(t *testing.T) {
	// Test CreateCollectionRequest validation
	t.Run("CreateCollectionRequest", func(t *testing.T) {
		// Valid request
		validReq := &CreateCollectionRequest{
			Name:          "Test Collection",
			Description:   phase5StringPtr("Test Description"),
			Version:       "1.0.0",
			BaseURL:       phase5StringPtr("https://api.test.com"),
			SchemaVersion: "2.1.0",
		}

		err := ValidateCreateCollectionRequest(validReq)
		assert.NoError(t, err)

		// Invalid request - missing name
		invalidReq := &CreateCollectionRequest{
			Version:       "1.0.0",
			SchemaVersion: "2.1.0",
		}

		err = ValidateCreateCollectionRequest(invalidReq)
		assert.Error(t, err)
	})

	// Test CreateEndpointRequest validation
	t.Run("CreateEndpointRequest", func(t *testing.T) {
		// Valid request
		validReq := &CreateEndpointRequest{
			Name:   "Test Endpoint",
			Method: "GET",
			URL:    "/api/test",
		}

		err := ValidateCreateEndpointRequest(validReq)
		assert.NoError(t, err)

		// Invalid request - invalid method
		invalidReq := &CreateEndpointRequest{
			Name:   "Test Endpoint",
			Method: "INVALID",
			URL:    "/api/test",
		}

		err = ValidateCreateEndpointRequest(invalidReq)
		assert.Error(t, err)
	})
}

func testIntegration(t *testing.T) {
	// Skip integration tests for now as they require full setup
	t.Skip("Integration tests require full application setup")
}

func setupPhase5TestDatabase(t *testing.T) *sql.DB {
	// Try to connect to test database
	dbURL := "postgres://localhost/test_db?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Logf("Cannot connect to test database: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		t.Logf("Cannot ping test database: %v", err)
		db.Close()
		return nil
	}

	return db
}

// Phase5MockRepository is a mock implementation for Phase 5 tests
type Phase5MockRepository struct {
	mock.Mock
}

func (m *Phase5MockRepository) CreateCollection(collection *Collection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetCollectionByID(id int64, companyID int64) (*Collection, error) {
	args := m.Called(id, companyID)
	return args.Get(0).(*Collection), args.Error(1)
}

func (m *Phase5MockRepository) UpdateCollection(collection *Collection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteCollection(id int64, companyID int64) error {
	args := m.Called(id, companyID)
	return args.Error(0)
}

// Add other required repository methods as stubs
func (m *Phase5MockRepository) GetCollections(companyID int64, limit, offset int) ([]*Collection, int64, error) {
	args := m.Called(companyID, limit, offset)
	return args.Get(0).([]*Collection), args.Get(1).(int64), args.Error(2)
}

func (m *Phase5MockRepository) GetCollectionWithStats(id int64, companyID int64) (*CollectionWithStats, error) {
	args := m.Called(id, companyID)
	return args.Get(0).(*CollectionWithStats), args.Error(1)
}

func (m *Phase5MockRepository) CreateFolder(folder *Folder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetFoldersByCollectionID(collectionID int64) ([]*Folder, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*Folder), args.Error(1)
}

func (m *Phase5MockRepository) GetFolderByID(id int64) (*Folder, error) {
	args := m.Called(id)
	return args.Get(0).(*Folder), args.Error(1)
}

func (m *Phase5MockRepository) GetFoldersHierarchy(collectionID int64) ([]*FolderWithChildren, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*FolderWithChildren), args.Error(1)
}

func (m *Phase5MockRepository) UpdateFolder(folder *Folder) error {
	args := m.Called(folder)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteFolder(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateEndpoint(endpoint *Endpoint) error {
	args := m.Called(endpoint)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetEndpointsByCollectionID(collectionID int64, limit, offset int) ([]*Endpoint, int64, error) {
	args := m.Called(collectionID, limit, offset)
	return args.Get(0).([]*Endpoint), args.Get(1).(int64), args.Error(2)
}

func (m *Phase5MockRepository) GetEndpointByID(id int64) (*Endpoint, error) {
	args := m.Called(id)
	return args.Get(0).(*Endpoint), args.Error(1)
}

func (m *Phase5MockRepository) GetEndpointWithDetails(id int64) (*EndpointWithDetails, error) {
	args := m.Called(id)
	return args.Get(0).(*EndpointWithDetails), args.Error(1)
}

func (m *Phase5MockRepository) UpdateEndpoint(endpoint *Endpoint) error {
	args := m.Called(endpoint)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteEndpoint(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateEnvironment(environment *Environment) error {
	args := m.Called(environment)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetEnvironmentsByCollectionID(collectionID int64) ([]*Environment, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*Environment), args.Error(1)
}

func (m *Phase5MockRepository) GetEnvironmentByID(id int64) (*Environment, error) {
	args := m.Called(id)
	return args.Get(0).(*Environment), args.Error(1)
}

func (m *Phase5MockRepository) GetEnvironmentWithVariables(id int64) (*EnvironmentWithVariables, error) {
	args := m.Called(id)
	return args.Get(0).(*EnvironmentWithVariables), args.Error(1)
}

func (m *Phase5MockRepository) UpdateEnvironment(environment *Environment) error {
	args := m.Called(environment)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteEnvironment(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateEnvironmentVariable(variable *EnvironmentVariable) error {
	args := m.Called(variable)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetEnvironmentVariables(environmentID int64) ([]*EnvironmentVariable, error) {
	args := m.Called(environmentID)
	return args.Get(0).([]*EnvironmentVariable), args.Error(1)
}

func (m *Phase5MockRepository) GetEnvironmentVariableByID(id int64) (*EnvironmentVariable, error) {
	args := m.Called(id)
	return args.Get(0).(*EnvironmentVariable), args.Error(1)
}

func (m *Phase5MockRepository) UpdateEnvironmentVariable(variable *EnvironmentVariable) error {
	args := m.Called(variable)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteEnvironmentVariable(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateTag(tag *Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetTagsByCollectionID(collectionID int64) ([]*Tag, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *Phase5MockRepository) GetTagByID(id int64) (*Tag, error) {
	args := m.Called(id)
	return args.Get(0).(*Tag), args.Error(1)
}

func (m *Phase5MockRepository) UpdateTag(tag *Tag) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteTag(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) AddEndpointTag(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *Phase5MockRepository) RemoveEndpointTag(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetEndpointTags(endpointID int64) ([]*Tag, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *Phase5MockRepository) GetAllTags() ([]*Tag, error) {
	args := m.Called()
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *Phase5MockRepository) GetTagsByEndpointID(endpointID int64) ([]*Tag, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Tag), args.Error(1)
}

func (m *Phase5MockRepository) AddTagToEndpoint(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *Phase5MockRepository) RemoveTagFromEndpoint(endpointID, tagID int64) error {
	args := m.Called(endpointID, tagID)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateHeader(header *Header) error {
	args := m.Called(header)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetHeadersByEndpointID(endpointID int64) ([]*Header, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Header), args.Error(1)
}

func (m *Phase5MockRepository) UpdateHeader(header *Header) error {
	args := m.Called(header)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteHeader(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateParameter(parameter *Parameter) error {
	args := m.Called(parameter)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetParametersByEndpointID(endpointID int64) ([]*Parameter, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Parameter), args.Error(1)
}

func (m *Phase5MockRepository) UpdateParameter(parameter *Parameter) error {
	args := m.Called(parameter)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteParameter(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateRequestBody(requestBody *RequestBody) error {
	args := m.Called(requestBody)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetRequestBodyByEndpointID(endpointID int64) (*RequestBody, error) {
	args := m.Called(endpointID)
	return args.Get(0).(*RequestBody), args.Error(1)
}

func (m *Phase5MockRepository) UpdateRequestBody(requestBody *RequestBody) error {
	args := m.Called(requestBody)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteRequestBody(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) CreateResponse(response *Response) error {
	args := m.Called(response)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetResponsesByEndpointID(endpointID int64) ([]*Response, error) {
	args := m.Called(endpointID)
	return args.Get(0).([]*Response), args.Error(1)
}

func (m *Phase5MockRepository) UpdateResponse(response *Response) error {
	args := m.Called(response)
	return args.Error(0)
}

func (m *Phase5MockRepository) DeleteResponse(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Phase5MockRepository) GetCollectionForExport(collectionID int64, companyID int64) (*CollectionExport, error) {
	args := m.Called(collectionID, companyID)
	return args.Get(0).(*CollectionExport), args.Error(1)
}

// Phase5MockRBACService is a mock implementation for Phase 5 tests
type Phase5MockRBACService struct {
	mock.Mock
}

func (m *Phase5MockRBACService) GetUserPermissions(userID int64) (*rbac.UserPermissions, error) {
	args := m.Called(userID)
	return args.Get(0).(*rbac.UserPermissions), args.Error(1)
}

func (m *Phase5MockRBACService) HasPermission(userID int64, moduleID int64, permission string) (bool, error) {
	args := m.Called(userID, moduleID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *Phase5MockRBACService) HasRole(userID int64, roleName string) (bool, error) {
	args := m.Called(userID, roleName)
	return args.Bool(0), args.Error(1)
}

func (m *Phase5MockRBACService) GetAccessibleModules(userID int64, permission string) ([]int64, error) {
	args := m.Called(userID, permission)
	return args.Get(0).([]int64), args.Error(1)
}

func (m *Phase5MockRBACService) IsSuperAdmin(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *Phase5MockRBACService) GetFilteredModules(userID int64, permission string, limit, offset int, category, search string, isActive *bool) ([]*rbac.ModuleInfo, error) {
	args := m.Called(userID, permission, limit, offset, category, search, isActive)
	return args.Get(0).([]*rbac.ModuleInfo), args.Error(1)
}

// Helper function for Phase 5 tests
func phase5StringPtr(s string) *string {
	return &s
}
