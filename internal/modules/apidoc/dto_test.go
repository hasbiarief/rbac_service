package apidoc

import (
	"database/sql"
	"gin-scalable-api/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateCollectionRequest(t *testing.T) {
	// Valid request
	validReq := &CreateCollectionRequest{
		Name:          "Test API",
		Description:   stringPtr("Test API description"),
		Version:       "1.0.0",
		BaseURL:       stringPtr("https://api.example.com"),
		SchemaVersion: "1.0",
	}

	err := ValidateCreateCollectionRequest(validReq)
	assert.NoError(t, err, "Valid request should not return error")

	// Invalid request - missing required fields
	invalidReq := &CreateCollectionRequest{
		Description: stringPtr("Test API description"),
	}

	err = ValidateCreateCollectionRequest(invalidReq)
	assert.Error(t, err, "Invalid request should return error")
}

func TestValidateCreateEndpointRequest(t *testing.T) {
	// Valid request
	validReq := &CreateEndpointRequest{
		Name:        "Get Users",
		Description: stringPtr("Get all users"),
		Method:      "GET",
		URL:         "/api/users",
		SortOrder:   intPtr(1),
		IsActive:    boolPtr(true),
	}

	err := ValidateCreateEndpointRequest(validReq)
	assert.NoError(t, err, "Valid request should not return error")

	// Invalid request - invalid HTTP method
	invalidReq := &CreateEndpointRequest{
		Name:   "Invalid Endpoint",
		Method: "INVALID",
		URL:    "/api/test",
	}

	err = ValidateCreateEndpointRequest(invalidReq)
	assert.Error(t, err, "Invalid HTTP method should return error")

	// Invalid request - invalid URL path
	invalidURLReq := &CreateEndpointRequest{
		Name:   "Invalid URL",
		Method: "GET",
		URL:    "invalid-url",
	}

	err = ValidateCreateEndpointRequest(invalidURLReq)
	assert.Error(t, err, "Invalid URL path should return error")
}

func TestValidateCreateTagRequest(t *testing.T) {
	// Valid request
	validReq := &CreateTagRequest{
		Name:        "Authentication",
		Color:       "#FF5733",
		Description: stringPtr("Authentication related endpoints"),
	}

	err := ValidateCreateTagRequest(validReq)
	assert.NoError(t, err, "Valid request should not return error")

	// Invalid request - invalid hex color
	invalidReq := &CreateTagRequest{
		Name:  "Invalid Color",
		Color: "not-a-color",
	}

	err = ValidateCreateTagRequest(invalidReq)
	assert.Error(t, err, "Invalid hex color should return error")
}

func TestValidateCreateEnvironmentVariableRequest(t *testing.T) {
	// Valid request
	validReq := &CreateEnvironmentVariableRequest{
		KeyName:     "API_KEY",
		Value:       stringPtr("secret-key"),
		Description: stringPtr("API key for authentication"),
		IsSecret:    boolPtr(true),
	}

	err := ValidateCreateEnvironmentVariableRequest(validReq)
	assert.NoError(t, err, "Valid request should not return error")

	// Invalid request - invalid environment variable name
	invalidReq := &CreateEnvironmentVariableRequest{
		KeyName: "invalid-name",
		Value:   stringPtr("value"),
	}

	err = ValidateCreateEnvironmentVariableRequest(invalidReq)
	assert.Error(t, err, "Invalid environment variable name should return error")
}

func TestCustomValidators(t *testing.T) {
	// Test with actual validation requests instead of mocking FieldLevel

	// Test version validator through actual validation
	validVersionReq := &CreateCollectionRequest{
		Name:          "Test",
		Version:       "1.0.0",
		SchemaVersion: "1.0",
	}
	err := ValidateCreateCollectionRequest(validVersionReq)
	assert.NoError(t, err, "Valid version should pass")

	invalidVersionReq := &CreateCollectionRequest{
		Name:          "Test",
		Version:       "invalid-version",
		SchemaVersion: "1.0",
	}
	err = ValidateCreateCollectionRequest(invalidVersionReq)
	assert.Error(t, err, "Invalid version should fail")

	// Test HTTP method validator through endpoint validation
	validMethodReq := &CreateEndpointRequest{
		Name:   "Test Endpoint",
		Method: "GET",
		URL:    "/api/test",
	}
	err = ValidateCreateEndpointRequest(validMethodReq)
	assert.NoError(t, err, "Valid HTTP method should pass")

	invalidMethodReq := &CreateEndpointRequest{
		Name:   "Test Endpoint",
		Method: "INVALID",
		URL:    "/api/test",
	}
	err = ValidateCreateEndpointRequest(invalidMethodReq)
	assert.Error(t, err, "Invalid HTTP method should fail")

	// Test hex color validator through tag validation
	validColorReq := &CreateTagRequest{
		Name:  "Test Tag",
		Color: "#FF5733",
	}
	err = ValidateCreateTagRequest(validColorReq)
	assert.NoError(t, err, "Valid hex color should pass")

	invalidColorReq := &CreateTagRequest{
		Name:  "Test Tag",
		Color: "invalid-color",
	}
	err = ValidateCreateTagRequest(invalidColorReq)
	assert.Error(t, err, "Invalid hex color should fail")
}

func TestTransformers(t *testing.T) {
	// Test Collection transformer
	collection := &Collection{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "Test Collection",
		Version:   "1.0.0",
		CompanyID: 1,
		IsActive:  true,
	}

	response := CollectionToResponse(collection)
	assert.Equal(t, collection.ID, response.ID)
	assert.Equal(t, collection.Name, response.Name)
	assert.Equal(t, collection.Version, response.Version)
	assert.Equal(t, collection.CompanyID, response.CompanyID)
	assert.Equal(t, collection.IsActive, response.IsActive)

	// Test request to model transformer
	createReq := &CreateCollectionRequest{
		Name:          "New Collection",
		Description:   stringPtr("Description"),
		Version:       "1.0.0",
		BaseURL:       stringPtr("https://api.example.com"),
		SchemaVersion: "1.0",
	}

	model := CreateCollectionRequestToModel(createReq, 1, 1)
	assert.Equal(t, createReq.Name, model.Name)
	assert.Equal(t, *createReq.Description, model.Description.String)
	assert.Equal(t, createReq.Version, model.Version)
	assert.Equal(t, *createReq.BaseURL, model.BaseURL.String)
	assert.Equal(t, createReq.SchemaVersion, model.SchemaVersion)
	assert.Equal(t, int64(1), model.CreatedBy)
	assert.Equal(t, int64(1), model.CompanyID)
	assert.True(t, model.IsActive)
}

func TestPaginationUtility(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	total := int64(10)
	limit := 3
	offset := 0

	result := CreatePaginatedResponse(data, total, limit, offset)

	assert.Equal(t, data, result["data"])
	assert.Equal(t, total, result["total"])
	assert.Equal(t, limit, result["limit"])
	assert.Equal(t, offset, result["offset"])
	assert.Equal(t, true, result["has_more"])

	// Test last page
	offset = 9
	limit = 3
	result = CreatePaginatedResponse(data, total, limit, offset)
	assert.Equal(t, false, result["has_more"])
}

func TestNullTypeConversions(t *testing.T) {
	// Test string pointer to null string
	str := "test"
	nullStr := pointerToNullString(&str)
	assert.True(t, nullStr.Valid)
	assert.Equal(t, str, nullStr.String)

	nullStr = pointerToNullString(nil)
	assert.False(t, nullStr.Valid)

	// Test null string to pointer
	nullStr = model.NullString{NullString: sql.NullString{String: "test", Valid: true}}
	ptr := nullStringToPointer(nullStr)
	require.NotNil(t, ptr)
	assert.Equal(t, "test", *ptr)

	nullStr = model.NullString{NullString: sql.NullString{Valid: false}}
	ptr = nullStringToPointer(nullStr)
	assert.Nil(t, ptr)

	// Test int64 pointer to null int64
	num := int64(123)
	nullInt := pointerToNullInt64(&num)
	assert.True(t, nullInt.Valid)
	assert.Equal(t, num, nullInt.Int64)

	nullInt = pointerToNullInt64(nil)
	assert.False(t, nullInt.Valid)

	// Test null int64 to pointer
	nullInt = model.NullInt64{NullInt64: sql.NullInt64{Int64: 123, Valid: true}}
	intPtr := nullInt64ToPointer(nullInt)
	require.NotNil(t, intPtr)
	assert.Equal(t, int64(123), *intPtr)

	nullInt = model.NullInt64{NullInt64: sql.NullInt64{Valid: false}}
	intPtr = nullInt64ToPointer(nullInt)
	assert.Nil(t, intPtr)
}

// Helper functions for tests
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
