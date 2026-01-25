package apidoc

import (
	"gin-scalable-api/pkg/model"
	"time"
)

// Collection represents an API documentation collection
type Collection struct {
	model.BaseModel
	Name          string           `json:"name" db:"name"`
	Description   model.NullString `json:"description" db:"description"`
	Version       string           `json:"version" db:"version"`
	BaseURL       model.NullString `json:"base_url" db:"base_url"`
	SchemaVersion string           `json:"schema_version" db:"schema_version"`
	CreatedBy     int64            `json:"created_by" db:"created_by"`
	CompanyID     int64            `json:"company_id" db:"company_id"`
	IsActive      bool             `json:"is_active" db:"is_active"`
}

func (Collection) TableName() string {
	return "api_collections"
}

// Folder represents a hierarchical folder structure
type Folder struct {
	model.BaseModel
	CollectionID int64            `json:"collection_id" db:"collection_id"`
	ParentID     model.NullInt64  `json:"parent_id" db:"parent_id"`
	Name         string           `json:"name" db:"name"`
	Description  model.NullString `json:"description" db:"description"`
	SortOrder    int              `json:"sort_order" db:"sort_order"`
}

func (Folder) TableName() string {
	return "api_folders"
}

// Endpoint represents an API endpoint
type Endpoint struct {
	model.BaseModel
	CollectionID int64            `json:"collection_id" db:"collection_id"`
	FolderID     model.NullInt64  `json:"folder_id" db:"folder_id"`
	Name         string           `json:"name" db:"name"`
	Description  model.NullString `json:"description" db:"description"`
	Method       string           `json:"method" db:"method"`
	URL          string           `json:"url" db:"url"`
	SortOrder    int              `json:"sort_order" db:"sort_order"`
	IsActive     bool             `json:"is_active" db:"is_active"`
}

func (Endpoint) TableName() string {
	return "api_endpoints"
}

// Header represents request/response headers
type Header struct {
	ID          int64            `json:"id" db:"id"`
	EndpointID  int64            `json:"endpoint_id" db:"endpoint_id"`
	KeyName     string           `json:"key_name" db:"key_name"`
	Value       model.NullString `json:"value" db:"value"`
	Description model.NullString `json:"description" db:"description"`
	IsRequired  bool             `json:"is_required" db:"is_required"`
	HeaderType  string           `json:"header_type" db:"header_type"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}

func (Header) TableName() string {
	return "api_headers"
}

// Parameter represents query, path, or form parameters
type Parameter struct {
	ID           int64            `json:"id" db:"id"`
	EndpointID   int64            `json:"endpoint_id" db:"endpoint_id"`
	Name         string           `json:"name" db:"name"`
	Type         string           `json:"type" db:"type"`
	DataType     string           `json:"data_type" db:"data_type"`
	Description  model.NullString `json:"description" db:"description"`
	DefaultValue model.NullString `json:"default_value" db:"default_value"`
	ExampleValue model.NullString `json:"example_value" db:"example_value"`
	IsRequired   bool             `json:"is_required" db:"is_required"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
}

func (Parameter) TableName() string {
	return "api_parameters"
}

// RequestBody represents request body documentation
type RequestBody struct {
	model.BaseModel
	EndpointID       int64            `json:"endpoint_id" db:"endpoint_id"`
	ContentType      string           `json:"content_type" db:"content_type"`
	BodyContent      model.NullString `json:"body_content" db:"body_content"`
	Description      model.NullString `json:"description" db:"description"`
	SchemaDefinition model.JSONB      `json:"schema_definition" db:"schema_definition"`
}

func (RequestBody) TableName() string {
	return "api_request_bodies"
}

// Response represents response examples
type Response struct {
	ID           int64            `json:"id" db:"id"`
	EndpointID   int64            `json:"endpoint_id" db:"endpoint_id"`
	StatusCode   int              `json:"status_code" db:"status_code"`
	StatusText   model.NullString `json:"status_text" db:"status_text"`
	ContentType  string           `json:"content_type" db:"content_type"`
	ResponseBody model.NullString `json:"response_body" db:"response_body"`
	Description  model.NullString `json:"description" db:"description"`
	IsDefault    bool             `json:"is_default" db:"is_default"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
}

func (Response) TableName() string {
	return "api_responses"
}

// Environment represents different API environments
type Environment struct {
	model.BaseModel
	CollectionID int64            `json:"collection_id" db:"collection_id"`
	Name         string           `json:"name" db:"name"`
	Description  model.NullString `json:"description" db:"description"`
	IsDefault    bool             `json:"is_default" db:"is_default"`
}

func (Environment) TableName() string {
	return "api_environments"
}

// EnvironmentVariable represents environment-specific variables
type EnvironmentVariable struct {
	ID            int64            `json:"id" db:"id"`
	EnvironmentID int64            `json:"environment_id" db:"environment_id"`
	KeyName       string           `json:"key_name" db:"key_name"`
	Value         model.NullString `json:"value" db:"value"`
	Description   model.NullString `json:"description" db:"description"`
	IsSecret      bool             `json:"is_secret" db:"is_secret"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
}

func (EnvironmentVariable) TableName() string {
	return "api_environment_variables"
}

// Tag represents endpoint categorization tags
type Tag struct {
	ID          int64            `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Color       string           `json:"color" db:"color"`
	Description model.NullString `json:"description" db:"description"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}

func (Tag) TableName() string {
	return "api_tags"
}

// Test represents test scripts for endpoints
type Test struct {
	model.BaseModel
	EndpointID    int64            `json:"endpoint_id" db:"endpoint_id"`
	TestType      string           `json:"test_type" db:"test_type"`
	ScriptContent model.NullString `json:"script_content" db:"script_content"`
	Description   model.NullString `json:"description" db:"description"`
}

func (Test) TableName() string {
	return "api_tests"
}

// EndpointTag represents the many-to-many relationship between endpoints and tags
type EndpointTag struct {
	EndpointID int64     `json:"endpoint_id" db:"endpoint_id"`
	TagID      int64     `json:"tag_id" db:"tag_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func (EndpointTag) TableName() string {
	return "api_endpoint_tags"
}

// CollectionWithStats represents a collection with statistics
type CollectionWithStats struct {
	Collection
	TotalFolders      int `json:"total_folders" db:"total_folders"`
	TotalEndpoints    int `json:"total_endpoints" db:"total_endpoints"`
	TotalEnvironments int `json:"total_environments" db:"total_environments"`
	GetEndpoints      int `json:"get_endpoints" db:"get_endpoints"`
	PostEndpoints     int `json:"post_endpoints" db:"post_endpoints"`
	PutEndpoints      int `json:"put_endpoints" db:"put_endpoints"`
	DeleteEndpoints   int `json:"delete_endpoints" db:"delete_endpoints"`
}

// FolderWithChildren represents a folder with its children for hierarchical display
type FolderWithChildren struct {
	Folder
	Children []*FolderWithChildren `json:"children,omitempty"`
}

// EndpointWithDetails represents an endpoint with all its related data
type EndpointWithDetails struct {
	Endpoint
	Headers     []*Header    `json:"headers,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"request_body,omitempty"`
	Responses   []*Response  `json:"responses,omitempty"`
	Tags        []*Tag       `json:"tags,omitempty"`
	Tests       []*Test      `json:"tests,omitempty"`
}

// EnvironmentWithVariables represents an environment with its variables
type EnvironmentWithVariables struct {
	Environment
	Variables []*EnvironmentVariable `json:"variables,omitempty"`
}

// CollectionExport represents a complete collection for export
type CollectionExport struct {
	Collection   Collection                  `json:"collection"`
	Folders      []*FolderWithChildren       `json:"folders"`
	Endpoints    []*EndpointWithDetails      `json:"endpoints"`
	Environments []*EnvironmentWithVariables `json:"environments"`
	Tags         []*Tag                      `json:"tags"`
}

// CollectionWithDetails represents a collection with all details for export
type CollectionWithDetails struct {
	Collection
	Folders     []Folder                  `json:"folders"`
	Endpoints   []EndpointWithDetails     `json:"endpoints"`
	Environment *EnvironmentWithVariables `json:"environment,omitempty"`
}
