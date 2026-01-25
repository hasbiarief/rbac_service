package apidoc

import (
	"database/sql"
	"fmt"
	"gin-scalable-api/pkg/database"
	"gin-scalable-api/pkg/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Helper functions for creating nullable types
func newNullString(s string) model.NullString {
	return model.NullString{NullString: sql.NullString{String: s, Valid: true}}
}

func newNullInt64(i int64) model.NullInt64 {
	return model.NullInt64{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}

func emptyNullString() model.NullString {
	return model.NullString{NullString: sql.NullString{Valid: false}}
}

func emptyNullInt64() model.NullInt64 {
	return model.NullInt64{NullInt64: sql.NullInt64{Valid: false}}
}

// RepositoryTestSuite defines the test suite for repository operations
type RepositoryTestSuite struct {
	suite.Suite
	db         *sql.DB
	repository Repository
	testData   *TestData
}

// TestData holds test data for various entities
type TestData struct {
	CompanyID     int64
	UserID        int64
	CollectionID  int64
	FolderID      int64
	EndpointID    int64
	EnvironmentID int64
	TagID         int64
}

// SetupSuite runs once before all tests
func (suite *RepositoryTestSuite) SetupSuite() {
	// Setup test database connection
	config := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "hasbi"),
		Password: getEnv("DB_PASS", "hasbi"),
		DBName:   getEnv("DB_NAME", "huminor_rbac"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbConn, err := database.NewConnection(config)
	require.NoError(suite.T(), err, "Failed to connect to test database")

	suite.db = dbConn.DB
	suite.repository = NewRepository(suite.db)

	// Initialize test data
	suite.testData = &TestData{
		CompanyID: 1, // Use existing company from seed data
		UserID:    1, // Use existing user from seed data
	}
}

// TearDownSuite runs once after all tests
func (suite *RepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest runs before each test
func (suite *RepositoryTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.cleanupTestData()
}

// TearDownTest runs after each test
func (suite *RepositoryTestSuite) TearDownTest() {
	// Clean up test data after each test
	suite.cleanupTestData()
}

// cleanupTestData removes all test data
func (suite *RepositoryTestSuite) cleanupTestData() {
	// Delete in reverse order of dependencies
	suite.db.Exec("DELETE FROM api_endpoint_tags WHERE endpoint_id IN (SELECT id FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_responses WHERE endpoint_id IN (SELECT id FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_request_bodies WHERE endpoint_id IN (SELECT id FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_parameters WHERE endpoint_id IN (SELECT id FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_headers WHERE endpoint_id IN (SELECT id FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_environment_variables WHERE environment_id IN (SELECT id FROM api_environments WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%'))")
	suite.db.Exec("DELETE FROM api_environments WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%')")
	suite.db.Exec("DELETE FROM api_endpoints WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%')")
	suite.db.Exec("DELETE FROM api_folders WHERE collection_id IN (SELECT id FROM api_collections WHERE name LIKE 'Test%')")
	suite.db.Exec("DELETE FROM api_collections WHERE name LIKE 'Test%'")
	suite.db.Exec("DELETE FROM api_tags WHERE name LIKE 'Test%'")
}

// Test Collection CRUD Operations
func (suite *RepositoryTestSuite) TestCollectionCRUD() {
	// Test Create Collection
	collection := &Collection{
		Name:          "Test API Collection",
		Description:   newNullString("Test collection for unit tests"),
		Version:       "1.0.0",
		BaseURL:       newNullString("https://api.test.com"),
		SchemaVersion: "1.0",
		CreatedBy:     suite.testData.UserID,
		CompanyID:     suite.testData.CompanyID,
		IsActive:      true,
	}

	err := suite.repository.CreateCollection(collection)
	require.NoError(suite.T(), err, "Failed to create collection")
	assert.NotZero(suite.T(), collection.ID, "Collection ID should be set after creation")
	assert.NotZero(suite.T(), collection.CreatedAt, "CreatedAt should be set")
	assert.NotZero(suite.T(), collection.UpdatedAt, "UpdatedAt should be set")

	suite.testData.CollectionID = collection.ID

	// Test Get Collection by ID
	retrievedCollection, err := suite.repository.GetCollectionByID(collection.ID, suite.testData.CompanyID)
	require.NoError(suite.T(), err, "Failed to get collection by ID")
	assert.Equal(suite.T(), collection.Name, retrievedCollection.Name)
	assert.Equal(suite.T(), collection.Description.String, retrievedCollection.Description.String)
	assert.Equal(suite.T(), collection.Version, retrievedCollection.Version)
	assert.Equal(suite.T(), collection.BaseURL.String, retrievedCollection.BaseURL.String)

	// Test Get Collections with pagination
	collections, total, err := suite.repository.GetCollections(suite.testData.CompanyID, 10, 0)
	require.NoError(suite.T(), err, "Failed to get collections")
	assert.GreaterOrEqual(suite.T(), total, int64(1), "Should have at least one collection")
	assert.GreaterOrEqual(suite.T(), len(collections), 1, "Should return at least one collection")

	// Test Update Collection
	collection.Name = "Updated Test API Collection"
	collection.Version = "1.1.0"
	err = suite.repository.UpdateCollection(collection)
	require.NoError(suite.T(), err, "Failed to update collection")

	updatedCollection, err := suite.repository.GetCollectionByID(collection.ID, suite.testData.CompanyID)
	require.NoError(suite.T(), err, "Failed to get updated collection")
	assert.Equal(suite.T(), "Updated Test API Collection", updatedCollection.Name)
	assert.Equal(suite.T(), "1.1.0", updatedCollection.Version)

	// Test Delete Collection
	err = suite.repository.DeleteCollection(collection.ID, suite.testData.CompanyID)
	require.NoError(suite.T(), err, "Failed to delete collection")

	// Verify deletion
	_, err = suite.repository.GetCollectionByID(collection.ID, suite.testData.CompanyID)
	assert.Error(suite.T(), err, "Should return error when getting deleted collection")
}

// Test Folder CRUD Operations and Hierarchy
func (suite *RepositoryTestSuite) TestFolderCRUDAndHierarchy() {
	// First create a collection
	collection := suite.createTestCollection()
	suite.testData.CollectionID = collection.ID

	// Test Create Root Folder
	rootFolder := &Folder{
		CollectionID: collection.ID,
		Name:         "Test Root Folder",
		Description:  newNullString("Root folder for testing"),
		SortOrder:    1,
	}

	err := suite.repository.CreateFolder(rootFolder)
	require.NoError(suite.T(), err, "Failed to create root folder")
	assert.NotZero(suite.T(), rootFolder.ID, "Folder ID should be set after creation")

	// Test Create Child Folder
	childFolder := &Folder{
		CollectionID: collection.ID,
		ParentID:     newNullInt64(rootFolder.ID),
		Name:         "Test Child Folder",
		Description:  newNullString("Child folder for testing"),
		SortOrder:    1,
	}

	err = suite.repository.CreateFolder(childFolder)
	require.NoError(suite.T(), err, "Failed to create child folder")
	assert.NotZero(suite.T(), childFolder.ID, "Child folder ID should be set after creation")

	// Test Get Folder by ID
	retrievedFolder, err := suite.repository.GetFolderByID(rootFolder.ID)
	require.NoError(suite.T(), err, "Failed to get folder by ID")
	assert.Equal(suite.T(), rootFolder.Name, retrievedFolder.Name)
	assert.Equal(suite.T(), rootFolder.CollectionID, retrievedFolder.CollectionID)

	// Test Get Folders by Collection ID
	folders, err := suite.repository.GetFoldersByCollectionID(collection.ID)
	require.NoError(suite.T(), err, "Failed to get folders by collection ID")
	assert.Len(suite.T(), folders, 2, "Should have 2 folders")

	// Test Get Folders Hierarchy
	hierarchy, err := suite.repository.GetFoldersHierarchy(collection.ID)
	require.NoError(suite.T(), err, "Failed to get folders hierarchy")
	assert.Len(suite.T(), hierarchy, 1, "Should have 1 root folder")
	assert.Len(suite.T(), hierarchy[0].Children, 1, "Root folder should have 1 child")
	assert.Equal(suite.T(), childFolder.ID, hierarchy[0].Children[0].ID)

	// Test Update Folder
	rootFolder.Name = "Updated Root Folder"
	err = suite.repository.UpdateFolder(rootFolder)
	require.NoError(suite.T(), err, "Failed to update folder")

	updatedFolder, err := suite.repository.GetFolderByID(rootFolder.ID)
	require.NoError(suite.T(), err, "Failed to get updated folder")
	assert.Equal(suite.T(), "Updated Root Folder", updatedFolder.Name)

	// Test Delete Folder
	err = suite.repository.DeleteFolder(childFolder.ID)
	require.NoError(suite.T(), err, "Failed to delete child folder")

	err = suite.repository.DeleteFolder(rootFolder.ID)
	require.NoError(suite.T(), err, "Failed to delete root folder")

	// Verify deletion
	folders, err = suite.repository.GetFoldersByCollectionID(collection.ID)
	require.NoError(suite.T(), err, "Failed to get folders after deletion")
	assert.Len(suite.T(), folders, 0, "Should have no folders after deletion")
}

// Test Endpoint CRUD Operations with Joins
func (suite *RepositoryTestSuite) TestEndpointCRUDWithJoins() {
	// Setup test data
	collection := suite.createTestCollection()
	folder := suite.createTestFolder(collection.ID)
	suite.testData.CollectionID = collection.ID
	suite.testData.FolderID = folder.ID

	// Test Create Endpoint
	endpoint := &Endpoint{
		CollectionID: collection.ID,
		FolderID:     newNullInt64(folder.ID),
		Name:         "Test Endpoint",
		Description:  newNullString("Test endpoint for unit tests"),
		Method:       "GET",
		URL:          "/api/test",
		SortOrder:    1,
		IsActive:     true,
	}

	err := suite.repository.CreateEndpoint(endpoint)
	require.NoError(suite.T(), err, "Failed to create endpoint")
	assert.NotZero(suite.T(), endpoint.ID, "Endpoint ID should be set after creation")

	suite.testData.EndpointID = endpoint.ID

	// Create related data for joins testing
	suite.createTestHeader(endpoint.ID)
	suite.createTestParameter(endpoint.ID)
	suite.createTestRequestBody(endpoint.ID)
	suite.createTestResponse(endpoint.ID)

	// Test Get Endpoint by ID
	retrievedEndpoint, err := suite.repository.GetEndpointByID(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to get endpoint by ID")
	assert.Equal(suite.T(), endpoint.Name, retrievedEndpoint.Name)
	assert.Equal(suite.T(), endpoint.Method, retrievedEndpoint.Method)
	assert.Equal(suite.T(), endpoint.URL, retrievedEndpoint.URL)

	// Test Get Endpoints by Collection ID
	endpoints, total, err := suite.repository.GetEndpointsByCollectionID(collection.ID, 10, 0)
	require.NoError(suite.T(), err, "Failed to get endpoints by collection ID")
	assert.Equal(suite.T(), int64(1), total, "Should have 1 endpoint")
	assert.Len(suite.T(), endpoints, 1, "Should return 1 endpoint")

	// Test Get Endpoint with Details (with joins)
	details, err := suite.repository.GetEndpointWithDetails(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to get endpoint with details")
	assert.Equal(suite.T(), endpoint.ID, details.ID)
	assert.NotEmpty(suite.T(), details.Headers, "Should have headers")
	assert.NotEmpty(suite.T(), details.Parameters, "Should have parameters")
	assert.NotNil(suite.T(), details.RequestBody, "Should have request body")
	assert.NotEmpty(suite.T(), details.Responses, "Should have responses")

	// Test Update Endpoint
	endpoint.Name = "Updated Test Endpoint"
	endpoint.Method = "POST"
	err = suite.repository.UpdateEndpoint(endpoint)
	require.NoError(suite.T(), err, "Failed to update endpoint")

	updatedEndpoint, err := suite.repository.GetEndpointByID(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to get updated endpoint")
	assert.Equal(suite.T(), "Updated Test Endpoint", updatedEndpoint.Name)
	assert.Equal(suite.T(), "POST", updatedEndpoint.Method)

	// Test Delete Endpoint
	err = suite.repository.DeleteEndpoint(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to delete endpoint")

	// Verify deletion
	_, err = suite.repository.GetEndpointByID(endpoint.ID)
	assert.Error(suite.T(), err, "Should return error when getting deleted endpoint")
}

// Test Environment Operations
func (suite *RepositoryTestSuite) TestEnvironmentOperations() {
	// Setup test data
	collection := suite.createTestCollection()
	suite.testData.CollectionID = collection.ID

	// Test Create Environment
	environment := &Environment{
		CollectionID: collection.ID,
		Name:         "Test Environment",
		Description:  newNullString("Test environment for unit tests"),
		IsDefault:    true,
	}

	err := suite.repository.CreateEnvironment(environment)
	require.NoError(suite.T(), err, "Failed to create environment")
	assert.NotZero(suite.T(), environment.ID, "Environment ID should be set after creation")

	suite.testData.EnvironmentID = environment.ID

	// Create environment variables
	variable := &EnvironmentVariable{
		EnvironmentID: environment.ID,
		KeyName:       "API_KEY",
		Value:         newNullString("test-api-key"),
		Description:   newNullString("Test API key"),
		IsSecret:      true,
	}

	err = suite.repository.CreateEnvironmentVariable(variable)
	require.NoError(suite.T(), err, "Failed to create environment variable")
	assert.NotZero(suite.T(), variable.ID, "Variable ID should be set after creation")

	// Test Get Environment by ID
	retrievedEnv, err := suite.repository.GetEnvironmentByID(environment.ID)
	require.NoError(suite.T(), err, "Failed to get environment by ID")
	assert.Equal(suite.T(), environment.Name, retrievedEnv.Name)
	assert.Equal(suite.T(), environment.IsDefault, retrievedEnv.IsDefault)

	// Test Get Environments by Collection ID
	environments, err := suite.repository.GetEnvironmentsByCollectionID(collection.ID)
	require.NoError(suite.T(), err, "Failed to get environments by collection ID")
	assert.Len(suite.T(), environments, 1, "Should have 1 environment")

	// Test Get Environment with Variables
	envWithVars, err := suite.repository.GetEnvironmentWithVariables(environment.ID)
	require.NoError(suite.T(), err, "Failed to get environment with variables")
	assert.Equal(suite.T(), environment.ID, envWithVars.ID)
	assert.Len(suite.T(), envWithVars.Variables, 1, "Should have 1 variable")
	assert.Equal(suite.T(), "API_KEY", envWithVars.Variables[0].KeyName)

	// Test Update Environment
	environment.Name = "Updated Test Environment"
	environment.IsDefault = false
	err = suite.repository.UpdateEnvironment(environment)
	require.NoError(suite.T(), err, "Failed to update environment")

	updatedEnv, err := suite.repository.GetEnvironmentByID(environment.ID)
	require.NoError(suite.T(), err, "Failed to get updated environment")
	assert.Equal(suite.T(), "Updated Test Environment", updatedEnv.Name)
	assert.False(suite.T(), updatedEnv.IsDefault)

	// Test Delete Environment
	err = suite.repository.DeleteEnvironment(environment.ID)
	require.NoError(suite.T(), err, "Failed to delete environment")

	// Verify deletion
	_, err = suite.repository.GetEnvironmentByID(environment.ID)
	assert.Error(suite.T(), err, "Should return error when getting deleted environment")
}

// Helper methods for creating test data

func (suite *RepositoryTestSuite) createTestCollection() *Collection {
	collection := &Collection{
		Name:          fmt.Sprintf("Test Collection %d", time.Now().UnixNano()),
		Description:   newNullString("Test collection"),
		Version:       "1.0.0",
		BaseURL:       newNullString("https://api.test.com"),
		SchemaVersion: "1.0",
		CreatedBy:     suite.testData.UserID,
		CompanyID:     suite.testData.CompanyID,
		IsActive:      true,
	}

	err := suite.repository.CreateCollection(collection)
	require.NoError(suite.T(), err, "Failed to create test collection")
	return collection
}

func (suite *RepositoryTestSuite) createTestFolder(collectionID int64) *Folder {
	folder := &Folder{
		CollectionID: collectionID,
		Name:         fmt.Sprintf("Test Folder %d", time.Now().UnixNano()),
		Description:  newNullString("Test folder"),
		SortOrder:    1,
	}

	err := suite.repository.CreateFolder(folder)
	require.NoError(suite.T(), err, "Failed to create test folder")
	return folder
}

func (suite *RepositoryTestSuite) createTestHeader(endpointID int64) *Header {
	header := &Header{
		EndpointID:  endpointID,
		KeyName:     "Content-Type",
		Value:       newNullString("application/json"),
		Description: newNullString("Content type header"),
		IsRequired:  true,
		HeaderType:  "request",
	}

	err := suite.repository.CreateHeader(header)
	require.NoError(suite.T(), err, "Failed to create test header")
	return header
}

func (suite *RepositoryTestSuite) createTestParameter(endpointID int64) *Parameter {
	parameter := &Parameter{
		EndpointID:   endpointID,
		Name:         "id",
		Type:         "path",
		DataType:     "integer",
		Description:  newNullString("Resource ID"),
		DefaultValue: emptyNullString(),
		ExampleValue: newNullString("123"),
		IsRequired:   true,
	}

	err := suite.repository.CreateParameter(parameter)
	require.NoError(suite.T(), err, "Failed to create test parameter")
	return parameter
}

func (suite *RepositoryTestSuite) createTestRequestBody(endpointID int64) *RequestBody {
	requestBody := &RequestBody{
		EndpointID:       endpointID,
		ContentType:      "application/json",
		BodyContent:      newNullString(`{"name": "test"}`),
		Description:      newNullString("Test request body"),
		SchemaDefinition: model.JSONB{"type": "object", "properties": map[string]interface{}{"name": map[string]interface{}{"type": "string"}}},
	}

	err := suite.repository.CreateRequestBody(requestBody)
	require.NoError(suite.T(), err, "Failed to create test request body")
	return requestBody
}

func (suite *RepositoryTestSuite) createTestResponse(endpointID int64) *Response {
	response := &Response{
		EndpointID:   endpointID,
		StatusCode:   200,
		StatusText:   newNullString("OK"),
		ContentType:  "application/json",
		ResponseBody: newNullString(`{"success": true}`),
		Description:  newNullString("Success response"),
		IsDefault:    true,
	}

	err := suite.repository.CreateResponse(response)
	require.NoError(suite.T(), err, "Failed to create test response")
	return response
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// TestRepositoryTestSuite runs the test suite
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
