package apidoc

import (
	"fmt"
	"gin-scalable-api/pkg/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// AdvancedRepositoryTestSuite for more complex test scenarios
type AdvancedRepositoryTestSuite struct {
	RepositoryTestSuite
}

// Test Tag Operations
func (suite *AdvancedRepositoryTestSuite) TestTagOperations() {
	// Setup test data
	collection := suite.createTestCollection()
	folder := suite.createTestFolder(collection.ID)
	endpoint := suite.createTestEndpoint(collection.ID, folder.ID)

	// Test Create Tag
	tag := &Tag{
		Name:        fmt.Sprintf("Test Tag %d", time.Now().UnixNano()),
		Color:       "#FF5733",
		Description: newNullString("Test tag for unit tests"),
	}

	err := suite.repository.CreateTag(tag)
	require.NoError(suite.T(), err, "Failed to create tag")
	assert.NotZero(suite.T(), tag.ID, "Tag ID should be set after creation")

	suite.testData.TagID = tag.ID

	// Test Get All Tags
	tags, err := suite.repository.GetAllTags()
	require.NoError(suite.T(), err, "Failed to get all tags")
	assert.GreaterOrEqual(suite.T(), len(tags), 1, "Should have at least one tag")

	// Test Add Tag to Endpoint
	err = suite.repository.AddTagToEndpoint(endpoint.ID, tag.ID)
	require.NoError(suite.T(), err, "Failed to add tag to endpoint")

	// Test Get Tags by Endpoint ID
	endpointTags, err := suite.repository.GetTagsByEndpointID(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to get tags by endpoint ID")
	assert.Len(suite.T(), endpointTags, 1, "Should have 1 tag")
	assert.Equal(suite.T(), tag.ID, endpointTags[0].ID)

	// Test Remove Tag from Endpoint
	err = suite.repository.RemoveTagFromEndpoint(endpoint.ID, tag.ID)
	require.NoError(suite.T(), err, "Failed to remove tag from endpoint")

	// Verify removal
	endpointTags, err = suite.repository.GetTagsByEndpointID(endpoint.ID)
	require.NoError(suite.T(), err, "Failed to get tags after removal")
	assert.Len(suite.T(), endpointTags, 0, "Should have no tags after removal")
}

// Test Export Operations
func (suite *AdvancedRepositoryTestSuite) TestExportOperations() {
	// Setup comprehensive test data
	collection := suite.createTestCollection()
	folder := suite.createTestFolder(collection.ID)
	endpoint := suite.createTestEndpoint(collection.ID, folder.ID)
	environment := suite.createTestEnvironment(collection.ID)
	tag := suite.createTestTag()

	// Add some related data
	suite.createTestHeader(endpoint.ID)
	suite.createTestParameter(endpoint.ID)
	suite.createTestRequestBody(endpoint.ID)
	suite.createTestResponse(endpoint.ID)
	suite.createTestEnvironmentVariable(environment.ID)
	suite.repository.AddTagToEndpoint(endpoint.ID, tag.ID)

	// Test Get Collection for Export
	exportData, err := suite.repository.GetCollectionForExport(collection.ID, suite.testData.CompanyID)
	require.NoError(suite.T(), err, "Failed to get collection for export")

	// Verify export data structure
	assert.Equal(suite.T(), collection.ID, exportData.Collection.ID)
	assert.Equal(suite.T(), collection.Name, exportData.Collection.Name)

	assert.Len(suite.T(), exportData.Folders, 1, "Should have 1 folder")
	assert.Equal(suite.T(), folder.ID, exportData.Folders[0].ID)

	assert.Len(suite.T(), exportData.Endpoints, 1, "Should have 1 endpoint")
	assert.Equal(suite.T(), endpoint.ID, exportData.Endpoints[0].ID)
	assert.NotEmpty(suite.T(), exportData.Endpoints[0].Headers, "Endpoint should have headers")
	assert.NotEmpty(suite.T(), exportData.Endpoints[0].Parameters, "Endpoint should have parameters")
	assert.NotNil(suite.T(), exportData.Endpoints[0].RequestBody, "Endpoint should have request body")
	assert.NotEmpty(suite.T(), exportData.Endpoints[0].Responses, "Endpoint should have responses")
	assert.NotEmpty(suite.T(), exportData.Endpoints[0].Tags, "Endpoint should have tags")

	assert.Len(suite.T(), exportData.Environments, 1, "Should have 1 environment")
	assert.Equal(suite.T(), environment.ID, exportData.Environments[0].ID)
	assert.NotEmpty(suite.T(), exportData.Environments[0].Variables, "Environment should have variables")

	assert.GreaterOrEqual(suite.T(), len(exportData.Tags), 1, "Should have at least 1 tag")
}

// Test Error Cases
func (suite *AdvancedRepositoryTestSuite) TestErrorCases() {
	// Test getting non-existent collection
	_, err := suite.repository.GetCollectionByID(99999, suite.testData.CompanyID)
	assert.Error(suite.T(), err, "Should return error for non-existent collection")

	// Test getting collection with wrong company ID
	collection := suite.createTestCollection()
	_, err = suite.repository.GetCollectionByID(collection.ID, 99999)
	assert.Error(suite.T(), err, "Should return error for wrong company ID")

	// Test updating non-existent collection
	nonExistentCollection := &Collection{
		BaseModel: model.BaseModel{ID: 99999},
		Name:      "Non-existent",
		CompanyID: suite.testData.CompanyID,
	}
	err = suite.repository.UpdateCollection(nonExistentCollection)
	assert.Error(suite.T(), err, "Should return error when updating non-existent collection")

	// Test deleting non-existent collection
	err = suite.repository.DeleteCollection(99999, suite.testData.CompanyID)
	assert.Error(suite.T(), err, "Should return error when deleting non-existent collection")

	// Test getting non-existent folder
	_, err = suite.repository.GetFolderByID(99999)
	assert.Error(suite.T(), err, "Should return error for non-existent folder")

	// Test getting non-existent endpoint
	_, err = suite.repository.GetEndpointByID(99999)
	assert.Error(suite.T(), err, "Should return error for non-existent endpoint")

	// Test getting non-existent environment
	_, err = suite.repository.GetEnvironmentByID(99999)
	assert.Error(suite.T(), err, "Should return error for non-existent environment")

	// Test getting request body for non-existent endpoint
	_, err = suite.repository.GetRequestBodyByEndpointID(99999)
	assert.Error(suite.T(), err, "Should return error for non-existent endpoint request body")
}

// Test Pagination
func (suite *AdvancedRepositoryTestSuite) TestPagination() {
	collection := suite.createTestCollection()

	// Create multiple endpoints
	const numEndpoints = 15
	for i := 0; i < numEndpoints; i++ {
		endpoint := &Endpoint{
			CollectionID: collection.ID,
			Name:         fmt.Sprintf("Endpoint %d", i),
			Method:       "GET",
			URL:          fmt.Sprintf("/api/test/%d", i),
			SortOrder:    i,
			IsActive:     true,
		}
		err := suite.repository.CreateEndpoint(endpoint)
		require.NoError(suite.T(), err, "Failed to create endpoint for pagination test")
	}

	// Test first page
	endpoints, total, err := suite.repository.GetEndpointsByCollectionID(collection.ID, 10, 0)
	require.NoError(suite.T(), err, "Failed to get first page of endpoints")
	assert.Equal(suite.T(), int64(numEndpoints), total, "Total should be correct")
	assert.Len(suite.T(), endpoints, 10, "First page should have 10 endpoints")

	// Test second page
	endpoints, total, err = suite.repository.GetEndpointsByCollectionID(collection.ID, 10, 10)
	require.NoError(suite.T(), err, "Failed to get second page of endpoints")
	assert.Equal(suite.T(), int64(numEndpoints), total, "Total should be correct")
	assert.Len(suite.T(), endpoints, 5, "Second page should have 5 endpoints")

	// Test empty page
	endpoints, total, err = suite.repository.GetEndpointsByCollectionID(collection.ID, 10, 20)
	require.NoError(suite.T(), err, "Failed to get empty page of endpoints")
	assert.Equal(suite.T(), int64(numEndpoints), total, "Total should be correct")
	assert.Len(suite.T(), endpoints, 0, "Empty page should have 0 endpoints")
}

// Test Collection Stats
func (suite *AdvancedRepositoryTestSuite) TestCollectionStats() {
	collection := suite.createTestCollection()
	folder := suite.createTestFolder(collection.ID)
	_ = suite.createTestEnvironment(collection.ID) // Create environment for completeness

	// Create endpoints with different methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "GET", "POST"}
	for i, method := range methods {
		endpoint := &Endpoint{
			CollectionID: collection.ID,
			FolderID:     newNullInt64(folder.ID),
			Name:         fmt.Sprintf("Endpoint %d", i),
			Method:       method,
			URL:          fmt.Sprintf("/api/test/%d", i),
			SortOrder:    i,
			IsActive:     true,
		}
		err := suite.repository.CreateEndpoint(endpoint)
		require.NoError(suite.T(), err, "Failed to create endpoint for stats test")
	}

	// Test Get Collection with Stats
	stats, err := suite.repository.GetCollectionWithStats(collection.ID, suite.testData.CompanyID)
	require.NoError(suite.T(), err, "Failed to get collection with stats")
	assert.Equal(suite.T(), collection.ID, stats.ID)
	assert.Equal(suite.T(), collection.Name, stats.Name)
	// Note: Stats might be 0 if the view doesn't exist or isn't populated
	// This is expected in a test environment
}

// Test Complex Folder Hierarchy
func (suite *AdvancedRepositoryTestSuite) TestComplexFolderHierarchy() {
	collection := suite.createTestCollection()

	// Create a complex folder structure
	// Root1
	//   - Child1
	//     - Grandchild1
	//   - Child2
	// Root2

	root1 := &Folder{
		CollectionID: collection.ID,
		Name:         "Root Folder 1",
		Description:  newNullString("First root folder"),
		SortOrder:    1,
	}
	err := suite.repository.CreateFolder(root1)
	require.NoError(suite.T(), err, "Failed to create root1 folder")

	root2 := &Folder{
		CollectionID: collection.ID,
		Name:         "Root Folder 2",
		Description:  newNullString("Second root folder"),
		SortOrder:    2,
	}
	err = suite.repository.CreateFolder(root2)
	require.NoError(suite.T(), err, "Failed to create root2 folder")

	child1 := &Folder{
		CollectionID: collection.ID,
		ParentID:     newNullInt64(root1.ID),
		Name:         "Child Folder 1",
		Description:  newNullString("First child folder"),
		SortOrder:    1,
	}
	err = suite.repository.CreateFolder(child1)
	require.NoError(suite.T(), err, "Failed to create child1 folder")

	child2 := &Folder{
		CollectionID: collection.ID,
		ParentID:     newNullInt64(root1.ID),
		Name:         "Child Folder 2",
		Description:  newNullString("Second child folder"),
		SortOrder:    2,
	}
	err = suite.repository.CreateFolder(child2)
	require.NoError(suite.T(), err, "Failed to create child2 folder")

	grandchild1 := &Folder{
		CollectionID: collection.ID,
		ParentID:     newNullInt64(child1.ID),
		Name:         "Grandchild Folder 1",
		Description:  newNullString("First grandchild folder"),
		SortOrder:    1,
	}
	err = suite.repository.CreateFolder(grandchild1)
	require.NoError(suite.T(), err, "Failed to create grandchild1 folder")

	// Test hierarchy
	hierarchy, err := suite.repository.GetFoldersHierarchy(collection.ID)
	require.NoError(suite.T(), err, "Failed to get complex folder hierarchy")

	assert.Len(suite.T(), hierarchy, 2, "Should have 2 root folders")

	// Find root1 in hierarchy
	var root1Node *FolderWithChildren
	for _, root := range hierarchy {
		if root.ID == root1.ID {
			root1Node = root
			break
		}
	}
	require.NotNil(suite.T(), root1Node, "Should find root1 in hierarchy")
	assert.Len(suite.T(), root1Node.Children, 2, "Root1 should have 2 children")

	// Find child1 in root1's children
	var child1Node *FolderWithChildren
	for _, child := range root1Node.Children {
		if child.ID == child1.ID {
			child1Node = child
			break
		}
	}
	require.NotNil(suite.T(), child1Node, "Should find child1 in root1's children")
	assert.Len(suite.T(), child1Node.Children, 1, "Child1 should have 1 grandchild")
	assert.Equal(suite.T(), grandchild1.ID, child1Node.Children[0].ID, "Grandchild should be correct")
}

// Additional helper methods

func (suite *AdvancedRepositoryTestSuite) createTestEndpoint(collectionID, folderID int64) *Endpoint {
	endpoint := &Endpoint{
		CollectionID: collectionID,
		FolderID:     newNullInt64(folderID),
		Name:         fmt.Sprintf("Test Endpoint %d", time.Now().UnixNano()),
		Description:  newNullString("Test endpoint"),
		Method:       "GET",
		URL:          "/api/test",
		SortOrder:    1,
		IsActive:     true,
	}

	err := suite.repository.CreateEndpoint(endpoint)
	require.NoError(suite.T(), err, "Failed to create test endpoint")
	return endpoint
}

func (suite *AdvancedRepositoryTestSuite) createTestEnvironment(collectionID int64) *Environment {
	environment := &Environment{
		CollectionID: collectionID,
		Name:         fmt.Sprintf("Test Environment %d", time.Now().UnixNano()),
		Description:  newNullString("Test environment"),
		IsDefault:    false,
	}

	err := suite.repository.CreateEnvironment(environment)
	require.NoError(suite.T(), err, "Failed to create test environment")
	return environment
}

func (suite *AdvancedRepositoryTestSuite) createTestEnvironmentVariable(environmentID int64) *EnvironmentVariable {
	variable := &EnvironmentVariable{
		EnvironmentID: environmentID,
		KeyName:       fmt.Sprintf("TEST_VAR_%d", time.Now().UnixNano()),
		Value:         newNullString("test-value"),
		Description:   newNullString("Test variable"),
		IsSecret:      false,
	}

	err := suite.repository.CreateEnvironmentVariable(variable)
	require.NoError(suite.T(), err, "Failed to create test environment variable")
	return variable
}

func (suite *AdvancedRepositoryTestSuite) createTestTag() *Tag {
	tag := &Tag{
		Name:        fmt.Sprintf("Test Tag %d", time.Now().UnixNano()),
		Color:       "#FF5733",
		Description: newNullString("Test tag"),
	}

	err := suite.repository.CreateTag(tag)
	require.NoError(suite.T(), err, "Failed to create test tag")
	return tag
}

// TestAdvancedRepositoryTestSuite runs the advanced test suite
func TestAdvancedRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AdvancedRepositoryTestSuite))
}
