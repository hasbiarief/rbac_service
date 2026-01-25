package export

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI      string                 `json:"openapi" yaml:"openapi"`
	Info         OpenAPIInfo            `json:"info" yaml:"info"`
	Servers      []OpenAPIServer        `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths        map[string]OpenAPIPath `json:"paths" yaml:"paths"`
	Components   *OpenAPIComponents     `json:"components,omitempty" yaml:"components,omitempty"`
	Security     []map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []OpenAPITag           `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *OpenAPIExternalDocs   `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// OpenAPIInfo contains API information
type OpenAPIInfo struct {
	Title          string          `json:"title" yaml:"title"`
	Description    string          `json:"description,omitempty" yaml:"description,omitempty"`
	Version        string          `json:"version" yaml:"version"`
	TermsOfService string          `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *OpenAPIContact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *OpenAPILicense `json:"license,omitempty" yaml:"license,omitempty"`
}

// OpenAPIContact contains contact information
type OpenAPIContact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// OpenAPILicense contains license information
type OpenAPILicense struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// OpenAPIServer represents a server
type OpenAPIServer struct {
	URL         string                     `json:"url" yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]OpenAPIVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// OpenAPIVariable represents server variable
type OpenAPIVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// OpenAPIPath represents path operations
type OpenAPIPath map[string]OpenAPIOperation

// OpenAPIOperation represents an operation
type OpenAPIOperation struct {
	Tags        []string                   `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary     string                     `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	OperationID string                     `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []OpenAPIParameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]OpenAPIResponse `json:"responses" yaml:"responses"`
	Callbacks   map[string]interface{}     `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	Deprecated  bool                       `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security    []map[string][]string      `json:"security,omitempty" yaml:"security,omitempty"`
	Servers     []OpenAPIServer            `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// OpenAPIParameter represents a parameter
type OpenAPIParameter struct {
	Name            string                 `json:"name" yaml:"name"`
	In              string                 `json:"in" yaml:"in"`
	Description     string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool                   `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool                   `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Style           string                 `json:"style,omitempty" yaml:"style,omitempty"`
	Explode         bool                   `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved   bool                   `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
	Schema          *OpenAPISchema         `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples        map[string]interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// OpenAPIRequestBody represents request body
type OpenAPIRequestBody struct {
	Description string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content" yaml:"content"`
	Required    bool                        `json:"required,omitempty" yaml:"required,omitempty"`
}

// OpenAPIResponse represents a response
type OpenAPIResponse struct {
	Description string                      `json:"description" yaml:"description"`
	Headers     map[string]OpenAPIHeader    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]OpenAPIMediaType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]interface{}      `json:"links,omitempty" yaml:"links,omitempty"`
}

// OpenAPIMediaType represents media type
type OpenAPIMediaType struct {
	Schema   *OpenAPISchema         `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]interface{} `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// OpenAPIHeader represents header
type OpenAPIHeader struct {
	Description     string         `json:"description,omitempty" yaml:"description,omitempty"`
	Required        bool           `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated      bool           `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	AllowEmptyValue bool           `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	Schema          *OpenAPISchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example         interface{}    `json:"example,omitempty" yaml:"example,omitempty"`
}

// OpenAPISchema represents JSON schema
type OpenAPISchema struct {
	Type                 string                   `json:"type,omitempty" yaml:"type,omitempty"`
	AllOf                []OpenAPISchema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf                []OpenAPISchema          `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf                []OpenAPISchema          `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	Not                  *OpenAPISchema           `json:"not,omitempty" yaml:"not,omitempty"`
	Items                *OpenAPISchema           `json:"items,omitempty" yaml:"items,omitempty"`
	Properties           map[string]OpenAPISchema `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties interface{}              `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Description          string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Format               string                   `json:"format,omitempty" yaml:"format,omitempty"`
	Default              interface{}              `json:"default,omitempty" yaml:"default,omitempty"`
	Title                string                   `json:"title,omitempty" yaml:"title,omitempty"`
	MultipleOf           float64                  `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum              float64                  `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum     bool                     `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum              float64                  `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum     bool                     `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength            int                      `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength            int                      `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern              string                   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems             int                      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems             int                      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems          bool                     `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	MaxProperties        int                      `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        int                      `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	Enum                 []interface{}            `json:"enum,omitempty" yaml:"enum,omitempty"`
	Example              interface{}              `json:"example,omitempty" yaml:"example,omitempty"`
}

// OpenAPIComponents contains reusable components
type OpenAPIComponents struct {
	Schemas         map[string]OpenAPISchema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]OpenAPIResponse       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]OpenAPIParameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Examples        map[string]interface{}           `json:"examples,omitempty" yaml:"examples,omitempty"`
	RequestBodies   map[string]OpenAPIRequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	Headers         map[string]OpenAPIHeader         `json:"headers,omitempty" yaml:"headers,omitempty"`
	SecuritySchemes map[string]OpenAPISecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Links           map[string]interface{}           `json:"links,omitempty" yaml:"links,omitempty"`
	Callbacks       map[string]interface{}           `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

// OpenAPISecurityScheme represents security scheme
type OpenAPISecurityScheme struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string      `json:"name,omitempty" yaml:"name,omitempty"`
	In               string      `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            interface{} `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OpenAPITag represents a tag
type OpenAPITag struct {
	Name         string               `json:"name" yaml:"name"`
	Description  string               `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *OpenAPIExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// OpenAPIExternalDocs represents external documentation
type OpenAPIExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// OpenAPIExporter handles OpenAPI specification export
type OpenAPIExporter struct {
	*BaseExporter
}

// NewOpenAPIExporter creates a new OpenAPI exporter
func NewOpenAPIExporter() *OpenAPIExporter {
	return &OpenAPIExporter{
		BaseExporter: NewBaseExporter(FormatOpenAPI),
	}
}

// Export exports collection to OpenAPI format
func (e *OpenAPIExporter) Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error) {
	if err := e.ValidateOptions(options); err != nil {
		return nil, err
	}

	// Create OpenAPI spec
	spec := &OpenAPISpec{
		OpenAPI: options.SpecVersion,
		Info: OpenAPIInfo{
			Title:       collection.Name,
			Description: collection.Description,
			Version:     collection.Version,
		},
		Paths:      make(map[string]OpenAPIPath),
		Components: &OpenAPIComponents{},
	}

	// Add servers if enabled
	if options.IncludeServers && collection.Environment != nil {
		servers := e.buildServers(collection.Environment)
		spec.Servers = servers
	}

	// Add security schemes if enabled
	if options.IncludeSecurity {
		securitySchemes := e.buildSecuritySchemes()
		if len(securitySchemes) > 0 {
			spec.Components.SecuritySchemes = securitySchemes
			spec.Security = []map[string][]string{
				{"BearerAuth": {}},
			}
		}
	}

	// Process endpoints
	tags := make(map[string]bool)
	for _, endpoint := range collection.Endpoints {
		path := e.normalizePath(endpoint.URL)
		method := strings.ToLower(endpoint.Method)

		// Initialize path if not exists
		if spec.Paths[path] == nil {
			spec.Paths[path] = make(OpenAPIPath)
		}

		// Convert endpoint to operation
		operation := e.convertEndpointToOperation(&endpoint, options)
		spec.Paths[path][method] = operation

		// Collect tags
		for _, tag := range operation.Tags {
			tags[tag] = true
		}
	}

	// Add tags
	for tagName := range tags {
		spec.Tags = append(spec.Tags, OpenAPITag{
			Name: tagName,
		})
	}

	// Serialize based on output format
	var content []byte
	var err error

	if options.OutputFormat == "yaml" {
		content, err = yaml.Marshal(spec)
	} else {
		content, err = e.SerializeJSON(spec)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to serialize OpenAPI spec: %w", err)
	}

	result := &ExportResult{
		Content:     string(content),
		ContentType: e.GetContentType(options),
		Filename:    e.GenerateFilename(collection.Name, options),
		Size:        int64(len(content)),
		GeneratedAt: time.Now(),
	}

	return result, nil
}

// buildServers creates server definitions from environment
func (e *OpenAPIExporter) buildServers(environment *EnvironmentWithVariables) []OpenAPIServer {
	servers := []OpenAPIServer{}

	// Find base_url variable
	baseURL := ""
	for _, variable := range environment.Variables {
		if variable.KeyName == "base_url" {
			baseURL = variable.Value
			break
		}
	}

	if baseURL != "" {
		server := OpenAPIServer{
			URL:         baseURL,
			Description: fmt.Sprintf("%s server", environment.Name),
		}
		servers = append(servers, server)
	}

	return servers
}

// buildSecuritySchemes creates security scheme definitions
func (e *OpenAPIExporter) buildSecuritySchemes() map[string]OpenAPISecurityScheme {
	schemes := make(map[string]OpenAPISecurityScheme)

	schemes["BearerAuth"] = OpenAPISecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT Authorization header using the Bearer scheme",
	}

	return schemes
}

// convertEndpointToOperation converts endpoint to OpenAPI operation
func (e *OpenAPIExporter) convertEndpointToOperation(endpoint *EndpointWithDetails, options *ExportOptions) OpenAPIOperation {
	operation := OpenAPIOperation{
		Summary:     endpoint.Name,
		Description: endpoint.Description,
		OperationID: e.generateOperationID(endpoint),
		Parameters:  []OpenAPIParameter{},
		Responses:   make(map[string]OpenAPIResponse),
	}

	// Add tags based on folder
	if endpoint.FolderID != nil {
		// This would need folder information to be passed or looked up
		// For now, we'll use a generic tag
		operation.Tags = []string{"API"}
	}

	// Process parameters
	for _, param := range endpoint.Parameters {
		openAPIParam := e.convertParameter(&param)
		operation.Parameters = append(operation.Parameters, openAPIParam)
	}

	// Process request body
	if endpoint.RequestBody != nil {
		requestBody := e.convertRequestBody(endpoint.RequestBody)
		operation.RequestBody = requestBody
	}

	// Process responses
	if len(endpoint.Responses) > 0 {
		for _, response := range endpoint.Responses {
			statusCode := strconv.Itoa(response.StatusCode)
			openAPIResponse := e.convertResponse(&response)
			operation.Responses[statusCode] = openAPIResponse
		}
	} else {
		// Add default response
		operation.Responses["200"] = OpenAPIResponse{
			Description: "Successful response",
		}
	}

	return operation
}

// convertParameter converts parameter to OpenAPI parameter
func (e *OpenAPIExporter) convertParameter(param *Parameter) OpenAPIParameter {
	openAPIParam := OpenAPIParameter{
		Name:        param.Name,
		In:          param.Type,
		Description: param.Description,
		Required:    param.IsRequired,
		Schema: &OpenAPISchema{
			Type: e.convertDataType(param.DataType),
		},
	}

	if param.ExampleValue != "" {
		openAPIParam.Example = param.ExampleValue
	}

	if param.DefaultValue != "" {
		openAPIParam.Schema.Default = param.DefaultValue
	}

	return openAPIParam
}

// convertRequestBody converts request body to OpenAPI request body
func (e *OpenAPIExporter) convertRequestBody(requestBody *RequestBody) *OpenAPIRequestBody {
	content := make(map[string]OpenAPIMediaType)

	mediaType := OpenAPIMediaType{}

	// Try to parse schema from JSON if available
	if requestBody.SchemaDefinition != "" {
		var schema OpenAPISchema
		if err := json.Unmarshal([]byte(requestBody.SchemaDefinition), &schema); err == nil {
			mediaType.Schema = &schema
		}
	}

	// Add example if available
	if requestBody.BodyContent != "" {
		var example interface{}
		if err := json.Unmarshal([]byte(requestBody.BodyContent), &example); err == nil {
			mediaType.Example = example
		} else {
			mediaType.Example = requestBody.BodyContent
		}
	}

	content[requestBody.ContentType] = mediaType

	return &OpenAPIRequestBody{
		Description: requestBody.Description,
		Content:     content,
		Required:    true,
	}
}

// convertResponse converts response to OpenAPI response
func (e *OpenAPIExporter) convertResponse(response *Response) OpenAPIResponse {
	openAPIResponse := OpenAPIResponse{
		Description: response.Description,
	}

	if response.Description == "" {
		openAPIResponse.Description = http.StatusText(response.StatusCode)
	}

	// Add content if available
	if response.ResponseBody != "" {
		content := make(map[string]OpenAPIMediaType)
		mediaType := OpenAPIMediaType{}

		// Try to parse as JSON for example
		var example interface{}
		if err := json.Unmarshal([]byte(response.ResponseBody), &example); err == nil {
			mediaType.Example = example
		} else {
			mediaType.Example = response.ResponseBody
		}

		content[response.ContentType] = mediaType
		openAPIResponse.Content = content
	}

	return openAPIResponse
}

// normalizePath normalizes URL path for OpenAPI
func (e *OpenAPIExporter) normalizePath(url string) string {
	// Remove base URL if present
	if strings.Contains(url, "://") {
		parts := strings.SplitN(url, "/", 4)
		if len(parts) >= 4 {
			url = "/" + parts[3]
		}
	}

	// Convert :param to {param} format
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + part[1:] + "}"
		}
	}

	return strings.Join(parts, "/")
}

// generateOperationID generates operation ID
func (e *OpenAPIExporter) generateOperationID(endpoint *EndpointWithDetails) string {
	method := strings.ToLower(endpoint.Method)
	name := strings.ReplaceAll(endpoint.Name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	return method + name
}

// convertDataType converts data type to OpenAPI type
func (e *OpenAPIExporter) convertDataType(dataType string) string {
	switch strings.ToLower(dataType) {
	case "int", "integer", "int64", "int32":
		return "integer"
	case "float", "float64", "float32", "double":
		return "number"
	case "bool", "boolean":
		return "boolean"
	case "array", "slice":
		return "array"
	case "object", "map":
		return "object"
	default:
		return "string"
	}
}

// GetSupportedFormats returns supported formats
func (e *OpenAPIExporter) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatOpenAPI}
}
