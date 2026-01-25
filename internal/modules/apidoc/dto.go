package apidoc

import (
	"database/sql"
	"gin-scalable-api/pkg/model"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ===== REQUEST DTOs =====

// Collection Request DTOs
type CreateCollectionRequest struct {
	Name          string  `json:"name" validate:"required,min=2,max=100"`
	Description   *string `json:"description" validate:"omitempty,max=500"`
	Version       string  `json:"version" validate:"required,version"`
	BaseURL       *string `json:"base_url" validate:"omitempty,url"`
	SchemaVersion string  `json:"schema_version" validate:"required,version"`
}

type UpdateCollectionRequest struct {
	Name          *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description   *string `json:"description" validate:"omitempty,max=500"`
	Version       *string `json:"version" validate:"omitempty,version"`
	BaseURL       *string `json:"base_url" validate:"omitempty,url"`
	SchemaVersion *string `json:"schema_version" validate:"omitempty,version"`
	IsActive      *bool   `json:"is_active"`
}

// Folder Request DTOs
type CreateFolderRequest struct {
	ParentID    *int64  `json:"parent_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`
}

type UpdateFolderRequest struct {
	ParentID    *int64  `json:"parent_id" validate:"omitempty,min=1"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`
}

// Endpoint Request DTOs
type CreateEndpointRequest struct {
	FolderID    *int64  `json:"folder_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Method      string  `json:"method" validate:"required,http_method"`
	URL         string  `json:"url" validate:"required,min=1,max=500,api_path"`
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateEndpointRequest struct {
	FolderID    *int64  `json:"folder_id" validate:"omitempty,min=1"`
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Method      *string `json:"method" validate:"omitempty,http_method"`
	URL         *string `json:"url" validate:"omitempty,min=1,max=500,api_path"`
	SortOrder   *int    `json:"sort_order" validate:"omitempty,min=0"`
	IsActive    *bool   `json:"is_active"`
}

// Header Request DTOs
type CreateHeaderRequest struct {
	KeyName     string  `json:"key_name" validate:"required,min=1,max=100"`
	Value       *string `json:"value" validate:"omitempty,max=1000"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsRequired  *bool   `json:"is_required"`
	HeaderType  string  `json:"header_type" validate:"required,oneof=request response"`
}

type UpdateHeaderRequest struct {
	KeyName     *string `json:"key_name" validate:"omitempty,min=1,max=100"`
	Value       *string `json:"value" validate:"omitempty,max=1000"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsRequired  *bool   `json:"is_required"`
	HeaderType  *string `json:"header_type" validate:"omitempty,oneof=request response"`
}

// Parameter Request DTOs
type CreateParameterRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=100"`
	Type         string  `json:"type" validate:"required,oneof=query path form header"`
	DataType     string  `json:"data_type" validate:"required,oneof=string integer number boolean array object"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
	DefaultValue *string `json:"default_value" validate:"omitempty,max=1000"`
	ExampleValue *string `json:"example_value" validate:"omitempty,max=1000"`
	IsRequired   *bool   `json:"is_required"`
}

type UpdateParameterRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=1,max=100"`
	Type         *string `json:"type" validate:"omitempty,oneof=query path form header"`
	DataType     *string `json:"data_type" validate:"omitempty,oneof=string integer number boolean array object"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
	DefaultValue *string `json:"default_value" validate:"omitempty,max=1000"`
	ExampleValue *string `json:"example_value" validate:"omitempty,max=1000"`
	IsRequired   *bool   `json:"is_required"`
}

// Request Body Request DTOs
type CreateRequestBodyRequest struct {
	ContentType      string                 `json:"content_type" validate:"required,content_type"`
	BodyContent      *string                `json:"body_content" validate:"omitempty"`
	Description      *string                `json:"description" validate:"omitempty,max=500"`
	SchemaDefinition map[string]interface{} `json:"schema_definition" validate:"omitempty"`
}

type UpdateRequestBodyRequest struct {
	ContentType      *string                `json:"content_type" validate:"omitempty,content_type"`
	BodyContent      *string                `json:"body_content" validate:"omitempty"`
	Description      *string                `json:"description" validate:"omitempty,max=500"`
	SchemaDefinition map[string]interface{} `json:"schema_definition" validate:"omitempty"`
}

// Response Request DTOs
type CreateResponseRequest struct {
	StatusCode   int     `json:"status_code" validate:"required,min=100,max=599"`
	StatusText   *string `json:"status_text" validate:"omitempty,max=100"`
	ContentType  string  `json:"content_type" validate:"required,content_type"`
	ResponseBody *string `json:"response_body" validate:"omitempty"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
	IsDefault    *bool   `json:"is_default"`
}

type UpdateResponseRequest struct {
	StatusCode   *int    `json:"status_code" validate:"omitempty,min=100,max=599"`
	StatusText   *string `json:"status_text" validate:"omitempty,max=100"`
	ContentType  *string `json:"content_type" validate:"omitempty,content_type"`
	ResponseBody *string `json:"response_body" validate:"omitempty"`
	Description  *string `json:"description" validate:"omitempty,max=500"`
	IsDefault    *bool   `json:"is_default"`
}

// Environment Request DTOs
type CreateEnvironmentRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsDefault   *bool   `json:"is_default"`
}

type UpdateEnvironmentRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsDefault   *bool   `json:"is_default"`
}

// Environment Variable Request DTOs
type CreateEnvironmentVariableRequest struct {
	KeyName     string  `json:"key_name" validate:"required,min=1,max=100,env_var_name"`
	Value       *string `json:"value" validate:"omitempty,max=1000"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsSecret    *bool   `json:"is_secret"`
}

type UpdateEnvironmentVariableRequest struct {
	KeyName     *string `json:"key_name" validate:"omitempty,min=1,max=100,env_var_name"`
	Value       *string `json:"value" validate:"omitempty,max=1000"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IsSecret    *bool   `json:"is_secret"`
}

// Tag Request DTOs
type CreateTagRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=50"`
	Color       string  `json:"color" validate:"required,hex_color"`
	Description *string `json:"description" validate:"omitempty,max=200"`
}

type UpdateTagRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=50"`
	Color       *string `json:"color" validate:"omitempty,hex_color"`
	Description *string `json:"description" validate:"omitempty,max=200"`
}

// List Request DTOs
type CollectionListRequest struct {
	Limit    int     `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int     `form:"offset" validate:"omitempty,min=0"`
	Search   *string `form:"search" validate:"omitempty,max=100"`
	IsActive *bool   `form:"is_active"`
}

type EndpointListRequest struct {
	Limit    int     `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int     `form:"offset" validate:"omitempty,min=0"`
	Search   *string `form:"search" validate:"omitempty,max=100"`
	Method   *string `form:"method" validate:"omitempty,http_method"`
	FolderID *int64  `form:"folder_id" validate:"omitempty,min=1"`
	IsActive *bool   `form:"is_active"`
}

// Export Request DTOs
type ExportRequest struct {
	Format        string  `form:"format" validate:"required,oneof=postman openapi insomnia swagger apidog"`
	EnvironmentID *int64  `form:"environment_id" validate:"omitempty,min=1"`
	IncludeTests  *bool   `form:"include_tests"`
	OutputFormat  *string `form:"output_format" validate:"omitempty,oneof=json yaml"`
}

// Folder Order Request DTO for reordering
type FolderOrderRequest struct {
	FolderID  int64 `json:"folder_id" validate:"required,min=1"`
	SortOrder int   `json:"sort_order" validate:"min=0"`
}

// Endpoint Bulk Operations DTOs
type EndpointBulkUpdateRequest struct {
	EndpointID    int64                 `json:"endpoint_id" validate:"required,min=1"`
	UpdateRequest UpdateEndpointRequest `json:"update_request"`
}

type BulkCreateEndpointsRequest struct {
	Endpoints []*CreateEndpointRequest `json:"endpoints" validate:"required,min=1,dive"`
}

type BulkUpdateEndpointsRequest struct {
	Updates []EndpointBulkUpdateRequest `json:"updates" validate:"required,min=1,dive"`
}

type BulkDeleteEndpointsRequest struct {
	EndpointIDs []int64 `json:"endpoint_ids" validate:"required,min=1,dive,min=1"`
}

type BulkMoveEndpointsRequest struct {
	EndpointIDs    []int64 `json:"endpoint_ids" validate:"required,min=1,dive,min=1"`
	TargetFolderID *int64  `json:"target_folder_id" validate:"omitempty,min=1"`
}

// ===== RESPONSE DTOs =====

// Collection Response DTOs
type CollectionResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	Version       string  `json:"version"`
	BaseURL       *string `json:"base_url"`
	SchemaVersion string  `json:"schema_version"`
	CreatedBy     int64   `json:"created_by"`
	CompanyID     int64   `json:"company_id"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CollectionWithStatsResponse struct {
	CollectionResponse
	TotalFolders      int `json:"total_folders"`
	TotalEndpoints    int `json:"total_endpoints"`
	TotalEnvironments int `json:"total_environments"`
	GetEndpoints      int `json:"get_endpoints"`
	PostEndpoints     int `json:"post_endpoints"`
	PutEndpoints      int `json:"put_endpoints"`
	DeleteEndpoints   int `json:"delete_endpoints"`
}

type CollectionListResponse struct {
	Data    []*CollectionResponse `json:"data"`
	Total   int64                 `json:"total"`
	Limit   int                   `json:"limit"`
	Offset  int                   `json:"offset"`
	HasMore bool                  `json:"has_more"`
}

// Folder Response DTOs
type FolderResponse struct {
	ID           int64   `json:"id"`
	CollectionID int64   `json:"collection_id"`
	ParentID     *int64  `json:"parent_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	SortOrder    int     `json:"sort_order"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type FolderWithChildrenResponse struct {
	FolderResponse
	Children []*FolderWithChildrenResponse `json:"children,omitempty"`
}

// Endpoint Response DTOs
type EndpointResponse struct {
	ID           int64   `json:"id"`
	CollectionID int64   `json:"collection_id"`
	FolderID     *int64  `json:"folder_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Method       string  `json:"method"`
	URL          string  `json:"url"`
	SortOrder    int     `json:"sort_order"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type EndpointWithDetailsResponse struct {
	EndpointResponse
	Headers     []*HeaderResponse    `json:"headers,omitempty"`
	Parameters  []*ParameterResponse `json:"parameters,omitempty"`
	RequestBody *RequestBodyResponse `json:"request_body,omitempty"`
	Responses   []*ResponseResponse  `json:"responses,omitempty"`
	Tags        []*TagResponse       `json:"tags,omitempty"`
	Tests       []*TestResponse      `json:"tests,omitempty"`
}

type EndpointListResponse struct {
	Data    []*EndpointResponse `json:"data"`
	Total   int64               `json:"total"`
	Limit   int                 `json:"limit"`
	Offset  int                 `json:"offset"`
	HasMore bool                `json:"has_more"`
}

// Header Response DTOs
type HeaderResponse struct {
	ID          int64   `json:"id"`
	EndpointID  int64   `json:"endpoint_id"`
	KeyName     string  `json:"key_name"`
	Value       *string `json:"value"`
	Description *string `json:"description"`
	IsRequired  bool    `json:"is_required"`
	HeaderType  string  `json:"header_type"`
	CreatedAt   string  `json:"created_at"`
}

// Parameter Response DTOs
type ParameterResponse struct {
	ID           int64   `json:"id"`
	EndpointID   int64   `json:"endpoint_id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	DataType     string  `json:"data_type"`
	Description  *string `json:"description"`
	DefaultValue *string `json:"default_value"`
	ExampleValue *string `json:"example_value"`
	IsRequired   bool    `json:"is_required"`
	CreatedAt    string  `json:"created_at"`
}

// Request Body Response DTOs
type RequestBodyResponse struct {
	ID               int64                  `json:"id"`
	EndpointID       int64                  `json:"endpoint_id"`
	ContentType      string                 `json:"content_type"`
	BodyContent      *string                `json:"body_content"`
	Description      *string                `json:"description"`
	SchemaDefinition map[string]interface{} `json:"schema_definition"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

// Response Response DTOs
type ResponseResponse struct {
	ID           int64   `json:"id"`
	EndpointID   int64   `json:"endpoint_id"`
	StatusCode   int     `json:"status_code"`
	StatusText   *string `json:"status_text"`
	ContentType  string  `json:"content_type"`
	ResponseBody *string `json:"response_body"`
	Description  *string `json:"description"`
	IsDefault    bool    `json:"is_default"`
	CreatedAt    string  `json:"created_at"`
}

// Environment Response DTOs
type EnvironmentResponse struct {
	ID           int64   `json:"id"`
	CollectionID int64   `json:"collection_id"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	IsDefault    bool    `json:"is_default"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type EnvironmentWithVariablesResponse struct {
	EnvironmentResponse
	Variables []*EnvironmentVariableResponse `json:"variables,omitempty"`
}

// Environment Variable Response DTOs
type EnvironmentVariableResponse struct {
	ID            int64   `json:"id"`
	EnvironmentID int64   `json:"environment_id"`
	KeyName       string  `json:"key_name"`
	Value         *string `json:"value"`
	Description   *string `json:"description"`
	IsSecret      bool    `json:"is_secret"`
	CreatedAt     string  `json:"created_at"`
}

// Tag Response DTOs
type TagResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// Test Response DTOs
type TestResponse struct {
	ID            int64   `json:"id"`
	EndpointID    int64   `json:"endpoint_id"`
	TestType      string  `json:"test_type"`
	ScriptContent *string `json:"script_content"`
	Description   *string `json:"description"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// Export Response DTOs
type ExportResponse struct {
	Format      string `json:"format"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	GeneratedAt string `json:"generated_at"`
}

// ===== VALIDATION =====

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	validate.RegisterValidation("version", validateVersion)
	validate.RegisterValidation("http_method", validateHTTPMethod)
	validate.RegisterValidation("api_path", validateAPIPath)
	validate.RegisterValidation("content_type", validateContentType)
	validate.RegisterValidation("hex_color", validateHexColor)
	validate.RegisterValidation("env_var_name", validateEnvVarName)
}

// Custom validation functions
func validateVersion(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	// Simple version validation (e.g., 1.0.0, 2.1, v1.0)
	versionRegex := regexp.MustCompile(`^v?(\d+)(\.\d+)*(\.\d+)*(-[a-zA-Z0-9]+)*$`)
	return versionRegex.MatchString(version)
}

func validateHTTPMethod(fl validator.FieldLevel) bool {
	method := strings.ToUpper(fl.Field().String())
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

func validateAPIPath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	// API path should start with / and contain valid URL characters
	pathRegex := regexp.MustCompile(`^/[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=%{}]*$`)
	return pathRegex.MatchString(path)
}

func validateContentType(fl validator.FieldLevel) bool {
	contentType := fl.Field().String()
	// Basic content type validation
	contentTypeRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9!#$&\-\^_]*\/[a-zA-Z0-9][a-zA-Z0-9!#$&\-\^_]*(\s*;\s*[a-zA-Z0-9\-]+=([a-zA-Z0-9\-]+|"[^"]*"))*$`)
	return contentTypeRegex.MatchString(contentType)
}

func validateHexColor(fl validator.FieldLevel) bool {
	color := fl.Field().String()
	// Hex color validation (#RRGGBB or #RGB)
	hexColorRegex := regexp.MustCompile(`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`)
	return hexColorRegex.MatchString(color)
}

func validateEnvVarName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	// Environment variable name validation (uppercase letters, numbers, underscores)
	envVarRegex := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	return envVarRegex.MatchString(name)
}

// ===== VALIDATION FUNCTIONS =====

// Collection validation functions
func ValidateCreateCollectionRequest(req *CreateCollectionRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateCollectionRequest(req *UpdateCollectionRequest) error {
	return validate.Struct(req)
}

func ValidateCollectionListRequest(req *CollectionListRequest) error {
	return validate.Struct(req)
}

// Folder validation functions
func ValidateCreateFolderRequest(req *CreateFolderRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateFolderRequest(req *UpdateFolderRequest) error {
	return validate.Struct(req)
}

// Endpoint validation functions
func ValidateCreateEndpointRequest(req *CreateEndpointRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateEndpointRequest(req *UpdateEndpointRequest) error {
	return validate.Struct(req)
}

func ValidateEndpointListRequest(req *EndpointListRequest) error {
	return validate.Struct(req)
}

// Header validation functions
func ValidateCreateHeaderRequest(req *CreateHeaderRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateHeaderRequest(req *UpdateHeaderRequest) error {
	return validate.Struct(req)
}

// Parameter validation functions
func ValidateCreateParameterRequest(req *CreateParameterRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateParameterRequest(req *UpdateParameterRequest) error {
	return validate.Struct(req)
}

// Request Body validation functions
func ValidateCreateRequestBodyRequest(req *CreateRequestBodyRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateRequestBodyRequest(req *UpdateRequestBodyRequest) error {
	return validate.Struct(req)
}

// Response validation functions
func ValidateCreateResponseRequest(req *CreateResponseRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateResponseRequest(req *UpdateResponseRequest) error {
	return validate.Struct(req)
}

// Environment validation functions
func ValidateCreateEnvironmentRequest(req *CreateEnvironmentRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateEnvironmentRequest(req *UpdateEnvironmentRequest) error {
	return validate.Struct(req)
}

// Environment Variable validation functions
func ValidateCreateEnvironmentVariableRequest(req *CreateEnvironmentVariableRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateEnvironmentVariableRequest(req *UpdateEnvironmentVariableRequest) error {
	return validate.Struct(req)
}

// Tag validation functions
func ValidateCreateTagRequest(req *CreateTagRequest) error {
	return validate.Struct(req)
}

func ValidateUpdateTagRequest(req *UpdateTagRequest) error {
	return validate.Struct(req)
}

// Export validation functions
func ValidateExportRequest(req *ExportRequest) error {
	return validate.Struct(req)
}

// Folder order validation functions
func ValidateFolderOrderRequest(req *FolderOrderRequest) error {
	return validate.Struct(req)
}

// Bulk operations validation functions
func ValidateBulkCreateEndpointsRequest(req *BulkCreateEndpointsRequest) error {
	return validate.Struct(req)
}

func ValidateBulkUpdateEndpointsRequest(req *BulkUpdateEndpointsRequest) error {
	return validate.Struct(req)
}

func ValidateBulkDeleteEndpointsRequest(req *BulkDeleteEndpointsRequest) error {
	return validate.Struct(req)
}

func ValidateBulkMoveEndpointsRequest(req *BulkMoveEndpointsRequest) error {
	return validate.Struct(req)
}

// ===== TRANSFORMATION UTILITIES =====

// Collection transformers
func CollectionToResponse(collection *Collection) *CollectionResponse {
	return &CollectionResponse{
		ID:            collection.ID,
		Name:          collection.Name,
		Description:   nullStringToPointer(collection.Description),
		Version:       collection.Version,
		BaseURL:       nullStringToPointer(collection.BaseURL),
		SchemaVersion: collection.SchemaVersion,
		CreatedBy:     collection.CreatedBy,
		CompanyID:     collection.CompanyID,
		IsActive:      collection.IsActive,
		CreatedAt:     collection.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     collection.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func CollectionWithStatsToResponse(stats *CollectionWithStats) *CollectionWithStatsResponse {
	return &CollectionWithStatsResponse{
		CollectionResponse: *CollectionToResponse(&stats.Collection),
		TotalFolders:       stats.TotalFolders,
		TotalEndpoints:     stats.TotalEndpoints,
		TotalEnvironments:  stats.TotalEnvironments,
		GetEndpoints:       stats.GetEndpoints,
		PostEndpoints:      stats.PostEndpoints,
		PutEndpoints:       stats.PutEndpoints,
		DeleteEndpoints:    stats.DeleteEndpoints,
	}
}

// Folder transformers
func FolderToResponse(folder *Folder) *FolderResponse {
	return &FolderResponse{
		ID:           folder.ID,
		CollectionID: folder.CollectionID,
		ParentID:     nullInt64ToPointer(folder.ParentID),
		Name:         folder.Name,
		Description:  nullStringToPointer(folder.Description),
		SortOrder:    folder.SortOrder,
		CreatedAt:    folder.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    folder.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func FolderWithChildrenToResponse(folder *FolderWithChildren) *FolderWithChildrenResponse {
	response := &FolderWithChildrenResponse{
		FolderResponse: *FolderToResponse(&folder.Folder),
		Children:       make([]*FolderWithChildrenResponse, len(folder.Children)),
	}

	for i, child := range folder.Children {
		response.Children[i] = FolderWithChildrenToResponse(child)
	}

	return response
}

// Endpoint transformers
func EndpointToResponse(endpoint *Endpoint) *EndpointResponse {
	return &EndpointResponse{
		ID:           endpoint.ID,
		CollectionID: endpoint.CollectionID,
		FolderID:     nullInt64ToPointer(endpoint.FolderID),
		Name:         endpoint.Name,
		Description:  nullStringToPointer(endpoint.Description),
		Method:       endpoint.Method,
		URL:          endpoint.URL,
		SortOrder:    endpoint.SortOrder,
		IsActive:     endpoint.IsActive,
		CreatedAt:    endpoint.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    endpoint.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func EndpointWithDetailsToResponse(details *EndpointWithDetails) *EndpointWithDetailsResponse {
	response := &EndpointWithDetailsResponse{
		EndpointResponse: *EndpointToResponse(&details.Endpoint),
		Headers:          make([]*HeaderResponse, len(details.Headers)),
		Parameters:       make([]*ParameterResponse, len(details.Parameters)),
		Responses:        make([]*ResponseResponse, len(details.Responses)),
		Tags:             make([]*TagResponse, len(details.Tags)),
		Tests:            make([]*TestResponse, len(details.Tests)),
	}

	for i, header := range details.Headers {
		response.Headers[i] = HeaderToResponse(header)
	}

	for i, param := range details.Parameters {
		response.Parameters[i] = ParameterToResponse(param)
	}

	if details.RequestBody != nil {
		response.RequestBody = RequestBodyToResponse(details.RequestBody)
	}

	for i, resp := range details.Responses {
		response.Responses[i] = ResponseToResponse(resp)
	}

	for i, tag := range details.Tags {
		response.Tags[i] = TagToResponse(tag)
	}

	for i, test := range details.Tests {
		response.Tests[i] = TestToResponse(test)
	}

	return response
}

// Header transformers
func HeaderToResponse(header *Header) *HeaderResponse {
	return &HeaderResponse{
		ID:          header.ID,
		EndpointID:  header.EndpointID,
		KeyName:     header.KeyName,
		Value:       nullStringToPointer(header.Value),
		Description: nullStringToPointer(header.Description),
		IsRequired:  header.IsRequired,
		HeaderType:  header.HeaderType,
		CreatedAt:   header.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Parameter transformers
func ParameterToResponse(param *Parameter) *ParameterResponse {
	return &ParameterResponse{
		ID:           param.ID,
		EndpointID:   param.EndpointID,
		Name:         param.Name,
		Type:         param.Type,
		DataType:     param.DataType,
		Description:  nullStringToPointer(param.Description),
		DefaultValue: nullStringToPointer(param.DefaultValue),
		ExampleValue: nullStringToPointer(param.ExampleValue),
		IsRequired:   param.IsRequired,
		CreatedAt:    param.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Request Body transformers
func RequestBodyToResponse(requestBody *RequestBody) *RequestBodyResponse {
	return &RequestBodyResponse{
		ID:               requestBody.ID,
		EndpointID:       requestBody.EndpointID,
		ContentType:      requestBody.ContentType,
		BodyContent:      nullStringToPointer(requestBody.BodyContent),
		Description:      nullStringToPointer(requestBody.Description),
		SchemaDefinition: map[string]interface{}(requestBody.SchemaDefinition),
		CreatedAt:        requestBody.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        requestBody.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Response transformers
func ResponseToResponse(response *Response) *ResponseResponse {
	return &ResponseResponse{
		ID:           response.ID,
		EndpointID:   response.EndpointID,
		StatusCode:   response.StatusCode,
		StatusText:   nullStringToPointer(response.StatusText),
		ContentType:  response.ContentType,
		ResponseBody: nullStringToPointer(response.ResponseBody),
		Description:  nullStringToPointer(response.Description),
		IsDefault:    response.IsDefault,
		CreatedAt:    response.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Environment transformers
func EnvironmentToResponse(env *Environment) *EnvironmentResponse {
	return &EnvironmentResponse{
		ID:           env.ID,
		CollectionID: env.CollectionID,
		Name:         env.Name,
		Description:  nullStringToPointer(env.Description),
		IsDefault:    env.IsDefault,
		CreatedAt:    env.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    env.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func EnvironmentWithVariablesToResponse(envWithVars *EnvironmentWithVariables) *EnvironmentWithVariablesResponse {
	response := &EnvironmentWithVariablesResponse{
		EnvironmentResponse: *EnvironmentToResponse(&envWithVars.Environment),
		Variables:           make([]*EnvironmentVariableResponse, len(envWithVars.Variables)),
	}

	for i, variable := range envWithVars.Variables {
		response.Variables[i] = EnvironmentVariableToResponse(variable)
	}

	return response
}

// Environment Variable transformers
func EnvironmentVariableToResponse(variable *EnvironmentVariable) *EnvironmentVariableResponse {
	return &EnvironmentVariableResponse{
		ID:            variable.ID,
		EnvironmentID: variable.EnvironmentID,
		KeyName:       variable.KeyName,
		Value:         nullStringToPointer(variable.Value),
		Description:   nullStringToPointer(variable.Description),
		IsSecret:      variable.IsSecret,
		CreatedAt:     variable.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Tag transformers
func TagToResponse(tag *Tag) *TagResponse {
	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: nullStringToPointer(tag.Description),
		CreatedAt:   tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Test transformers
func TestToResponse(test *Test) *TestResponse {
	return &TestResponse{
		ID:            test.ID,
		EndpointID:    test.EndpointID,
		TestType:      test.TestType,
		ScriptContent: nullStringToPointer(test.ScriptContent),
		Description:   nullStringToPointer(test.Description),
		CreatedAt:     test.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     test.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ===== UTILITY FUNCTIONS =====

// Helper functions for null type conversions
func nullStringToPointer(ns model.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullInt64ToPointer(ni model.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

func pointerToNullString(s *string) model.NullString {
	if s != nil {
		return model.NullString{NullString: sql.NullString{String: *s, Valid: true}}
	}
	return model.NullString{NullString: sql.NullString{Valid: false}}
}

func pointerToNullInt64(i *int64) model.NullInt64 {
	if i != nil {
		return model.NullInt64{NullInt64: sql.NullInt64{Int64: *i, Valid: true}}
	}
	return model.NullInt64{NullInt64: sql.NullInt64{Valid: false}}
}

// Pagination utility
func CreatePaginatedResponse(data interface{}, total int64, limit, offset int) map[string]interface{} {
	hasMore := int64(offset+limit) < total
	return map[string]interface{}{
		"data":     data,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"has_more": hasMore,
	}
}

// Request to Model transformers
func CreateCollectionRequestToModel(req *CreateCollectionRequest, createdBy, companyID int64) *Collection {
	return &Collection{
		Name:          req.Name,
		Description:   pointerToNullString(req.Description),
		Version:       req.Version,
		BaseURL:       pointerToNullString(req.BaseURL),
		SchemaVersion: req.SchemaVersion,
		CreatedBy:     createdBy,
		CompanyID:     companyID,
		IsActive:      true,
	}
}

func CreateFolderRequestToModel(req *CreateFolderRequest, collectionID int64) *Folder {
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	return &Folder{
		CollectionID: collectionID,
		ParentID:     pointerToNullInt64(req.ParentID),
		Name:         req.Name,
		Description:  pointerToNullString(req.Description),
		SortOrder:    sortOrder,
	}
}

func CreateEndpointRequestToModel(req *CreateEndpointRequest, collectionID int64) *Endpoint {
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &Endpoint{
		CollectionID: collectionID,
		FolderID:     pointerToNullInt64(req.FolderID),
		Name:         req.Name,
		Description:  pointerToNullString(req.Description),
		Method:       strings.ToUpper(req.Method),
		URL:          req.URL,
		SortOrder:    sortOrder,
		IsActive:     isActive,
	}
}
