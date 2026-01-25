package export

import "time"

// Collection represents an API documentation collection for export
type Collection struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Version       string    `json:"version"`
	BaseURL       string    `json:"base_url"`
	SchemaVersion string    `json:"schema_version"`
	CreatedBy     int64     `json:"created_by"`
	CompanyID     int64     `json:"company_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Folder represents a hierarchical folder structure for export
type Folder struct {
	ID           int64     `json:"id"`
	CollectionID int64     `json:"collection_id"`
	ParentID     *int64    `json:"parent_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Endpoint represents an API endpoint for export
type Endpoint struct {
	ID           int64     `json:"id"`
	CollectionID int64     `json:"collection_id"`
	FolderID     *int64    `json:"folder_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	SortOrder    int       `json:"sort_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Header represents an HTTP header for export
type Header struct {
	ID          int64     `json:"id"`
	EndpointID  int64     `json:"endpoint_id"`
	KeyName     string    `json:"key_name"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	IsRequired  bool      `json:"is_required"`
	HeaderType  string    `json:"header_type"`
	CreatedAt   time.Time `json:"created_at"`
}

// Parameter represents a parameter for export
type Parameter struct {
	ID           int64     `json:"id"`
	EndpointID   int64     `json:"endpoint_id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	DataType     string    `json:"data_type"`
	Description  string    `json:"description"`
	DefaultValue string    `json:"default_value"`
	ExampleValue string    `json:"example_value"`
	IsRequired   bool      `json:"is_required"`
	CreatedAt    time.Time `json:"created_at"`
}

// RequestBody represents request body for export
type RequestBody struct {
	ID               int64     `json:"id"`
	EndpointID       int64     `json:"endpoint_id"`
	ContentType      string    `json:"content_type"`
	BodyContent      string    `json:"body_content"`
	Description      string    `json:"description"`
	SchemaDefinition string    `json:"schema_definition"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Response represents a response for export
type Response struct {
	ID           int64     `json:"id"`
	EndpointID   int64     `json:"endpoint_id"`
	StatusCode   int       `json:"status_code"`
	StatusText   string    `json:"status_text"`
	ContentType  string    `json:"content_type"`
	ResponseBody string    `json:"response_body"`
	Description  string    `json:"description"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
}

// Test represents test script for export
type Test struct {
	ID            int64     `json:"id"`
	EndpointID    int64     `json:"endpoint_id"`
	TestType      string    `json:"test_type"`
	ScriptContent string    `json:"script_content"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Environment represents an environment for export
type Environment struct {
	ID           int64     `json:"id"`
	CollectionID int64     `json:"collection_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EnvironmentVariable represents environment variable for export
type EnvironmentVariable struct {
	ID            int64     `json:"id"`
	EnvironmentID int64     `json:"environment_id"`
	KeyName       string    `json:"key_name"`
	Value         string    `json:"value"`
	Description   string    `json:"description"`
	IsSecret      bool      `json:"is_secret"`
	CreatedAt     time.Time `json:"created_at"`
}

// Tag represents a tag for export
type Tag struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// EndpointWithDetails represents an endpoint with all its related data for export
type EndpointWithDetails struct {
	Endpoint
	Headers     []Header     `json:"headers,omitempty"`
	Parameters  []Parameter  `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"request_body,omitempty"`
	Responses   []Response   `json:"responses,omitempty"`
	Tags        []Tag        `json:"tags,omitempty"`
	Tests       []Test       `json:"tests,omitempty"`
}

// EnvironmentWithVariables represents an environment with its variables for export
type EnvironmentWithVariables struct {
	Environment
	Variables []EnvironmentVariable `json:"variables,omitempty"`
}

// CollectionWithDetails represents a collection with all details for export
type CollectionWithDetails struct {
	Collection
	Folders     []Folder                  `json:"folders"`
	Endpoints   []EndpointWithDetails     `json:"endpoints"`
	Environment *EnvironmentWithVariables `json:"environment,omitempty"`
}
