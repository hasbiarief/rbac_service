package export

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// SwaggerSpec represents Swagger 2.0 specification
type SwaggerSpec struct {
	Swagger             string                       `json:"swagger"`
	Info                SwaggerInfo                  `json:"info"`
	Host                string                       `json:"host,omitempty"`
	BasePath            string                       `json:"basePath,omitempty"`
	Schemes             []string                     `json:"schemes,omitempty"`
	Consumes            []string                     `json:"consumes,omitempty"`
	Produces            []string                     `json:"produces,omitempty"`
	Paths               map[string]SwaggerPath       `json:"paths"`
	Definitions         map[string]SwaggerDefinition `json:"definitions,omitempty"`
	Parameters          map[string]SwaggerParameter  `json:"parameters,omitempty"`
	Responses           map[string]SwaggerResponse   `json:"responses,omitempty"`
	SecurityDefinitions map[string]SwaggerSecurity   `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string        `json:"security,omitempty"`
	Tags                []SwaggerTag                 `json:"tags,omitempty"`
	ExternalDocs        *SwaggerExternalDocs         `json:"externalDocs,omitempty"`
}

// SwaggerInfo contains API information
type SwaggerInfo struct {
	Title          string          `json:"title"`
	Description    string          `json:"description,omitempty"`
	Version        string          `json:"version"`
	TermsOfService string          `json:"termsOfService,omitempty"`
	Contact        *SwaggerContact `json:"contact,omitempty"`
	License        *SwaggerLicense `json:"license,omitempty"`
}

// SwaggerContact contains contact information
type SwaggerContact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// SwaggerLicense contains license information
type SwaggerLicense struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// SwaggerPath represents path operations
type SwaggerPath map[string]SwaggerOperation

// SwaggerOperation represents an operation
type SwaggerOperation struct {
	Tags        []string                   `json:"tags,omitempty"`
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	OperationID string                     `json:"operationId,omitempty"`
	Consumes    []string                   `json:"consumes,omitempty"`
	Produces    []string                   `json:"produces,omitempty"`
	Parameters  []SwaggerParameter         `json:"parameters,omitempty"`
	Responses   map[string]SwaggerResponse `json:"responses"`
	Schemes     []string                   `json:"schemes,omitempty"`
	Deprecated  bool                       `json:"deprecated,omitempty"`
	Security    []map[string][]string      `json:"security,omitempty"`
}

// SwaggerParameter represents a parameter
type SwaggerParameter struct {
	Name             string             `json:"name"`
	In               string             `json:"in"`
	Description      string             `json:"description,omitempty"`
	Required         bool               `json:"required,omitempty"`
	Type             string             `json:"type,omitempty"`
	Format           string             `json:"format,omitempty"`
	Items            *SwaggerItems      `json:"items,omitempty"`
	CollectionFormat string             `json:"collectionFormat,omitempty"`
	Default          interface{}        `json:"default,omitempty"`
	Maximum          float64            `json:"maximum,omitempty"`
	ExclusiveMaximum bool               `json:"exclusiveMaximum,omitempty"`
	Minimum          float64            `json:"minimum,omitempty"`
	ExclusiveMinimum bool               `json:"exclusiveMinimum,omitempty"`
	MaxLength        int                `json:"maxLength,omitempty"`
	MinLength        int                `json:"minLength,omitempty"`
	Pattern          string             `json:"pattern,omitempty"`
	MaxItems         int                `json:"maxItems,omitempty"`
	MinItems         int                `json:"minItems,omitempty"`
	UniqueItems      bool               `json:"uniqueItems,omitempty"`
	Enum             []interface{}      `json:"enum,omitempty"`
	MultipleOf       float64            `json:"multipleOf,omitempty"`
	Schema           *SwaggerDefinition `json:"schema,omitempty"`
}

// SwaggerItems represents array items
type SwaggerItems struct {
	Type             string        `json:"type,omitempty"`
	Format           string        `json:"format,omitempty"`
	Items            *SwaggerItems `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Maximum          float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        int           `json:"maxLength,omitempty"`
	MinLength        int           `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         int           `json:"maxItems,omitempty"`
	MinItems         int           `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	MultipleOf       float64       `json:"multipleOf,omitempty"`
}

// SwaggerResponse represents a response
type SwaggerResponse struct {
	Description string                   `json:"description"`
	Schema      *SwaggerDefinition       `json:"schema,omitempty"`
	Headers     map[string]SwaggerHeader `json:"headers,omitempty"`
	Examples    map[string]interface{}   `json:"examples,omitempty"`
}

// SwaggerHeader represents header
type SwaggerHeader struct {
	Description      string        `json:"description,omitempty"`
	Type             string        `json:"type"`
	Format           string        `json:"format,omitempty"`
	Items            *SwaggerItems `json:"items,omitempty"`
	CollectionFormat string        `json:"collectionFormat,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Maximum          float64       `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          float64       `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        int           `json:"maxLength,omitempty"`
	MinLength        int           `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         int           `json:"maxItems,omitempty"`
	MinItems         int           `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
	MultipleOf       float64       `json:"multipleOf,omitempty"`
}

// SwaggerDefinition represents schema definition
type SwaggerDefinition struct {
	Type                 string                       `json:"type,omitempty"`
	AllOf                []SwaggerDefinition          `json:"allOf,omitempty"`
	Properties           map[string]SwaggerDefinition `json:"properties,omitempty"`
	AdditionalProperties interface{}                  `json:"additionalProperties,omitempty"`
	Description          string                       `json:"description,omitempty"`
	Format               string                       `json:"format,omitempty"`
	Default              interface{}                  `json:"default,omitempty"`
	Title                string                       `json:"title,omitempty"`
	MultipleOf           float64                      `json:"multipleOf,omitempty"`
	Maximum              float64                      `json:"maximum,omitempty"`
	ExclusiveMaximum     bool                         `json:"exclusiveMaximum,omitempty"`
	Minimum              float64                      `json:"minimum,omitempty"`
	ExclusiveMinimum     bool                         `json:"exclusiveMinimum,omitempty"`
	MaxLength            int                          `json:"maxLength,omitempty"`
	MinLength            int                          `json:"minLength,omitempty"`
	Pattern              string                       `json:"pattern,omitempty"`
	MaxItems             int                          `json:"maxItems,omitempty"`
	MinItems             int                          `json:"minItems,omitempty"`
	UniqueItems          bool                         `json:"uniqueItems,omitempty"`
	MaxProperties        int                          `json:"maxProperties,omitempty"`
	MinProperties        int                          `json:"minProperties,omitempty"`
	Required             []string                     `json:"required,omitempty"`
	Enum                 []interface{}                `json:"enum,omitempty"`
	Items                *SwaggerDefinition           `json:"items,omitempty"`
	Example              interface{}                  `json:"example,omitempty"`
	Ref                  string                       `json:"$ref,omitempty"`
}

// SwaggerSecurity represents security scheme
type SwaggerSecurity struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// SwaggerTag represents a tag
type SwaggerTag struct {
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	ExternalDocs *SwaggerExternalDocs `json:"externalDocs,omitempty"`
}

// SwaggerExternalDocs represents external documentation
type SwaggerExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// SwaggerExporter handles Swagger 2.0 export
type SwaggerExporter struct {
	*BaseExporter
}

// NewSwaggerExporter creates a new Swagger exporter
func NewSwaggerExporter() *SwaggerExporter {
	return &SwaggerExporter{
		BaseExporter: NewBaseExporter(FormatSwagger),
	}
}

// Export exports collection to Swagger format
func (e *SwaggerExporter) Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error) {
	if err := e.ValidateOptions(options); err != nil {
		return nil, err
	}

	// Create Swagger spec
	spec := &SwaggerSpec{
		Swagger: "2.0",
		Info: SwaggerInfo{
			Title:       collection.Name,
			Description: collection.Description,
			Version:     collection.Version,
		},
		Paths:       make(map[string]SwaggerPath),
		Definitions: make(map[string]SwaggerDefinition),
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
	}

	// Add host and base path from environment
	if collection.Environment != nil {
		host, basePath := e.extractHostAndBasePath(collection.Environment)
		if host != "" {
			spec.Host = host
		}
		if basePath != "" {
			spec.BasePath = basePath
		}
		spec.Schemes = []string{"http", "https"}
	}

	// Add security definitions if enabled
	if options.IncludeSecurity {
		securityDefs := e.buildSecurityDefinitions()
		if len(securityDefs) > 0 {
			spec.SecurityDefinitions = securityDefs
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
			spec.Paths[path] = make(SwaggerPath)
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
		spec.Tags = append(spec.Tags, SwaggerTag{
			Name: tagName,
		})
	}

	// Add common error definition
	spec.Definitions["Error"] = SwaggerDefinition{
		Type: "object",
		Properties: map[string]SwaggerDefinition{
			"success": {
				Type:    "boolean",
				Example: false,
			},
			"message": {
				Type:    "string",
				Example: "Error message",
			},
		},
	}

	// Serialize to JSON
	content, err := e.SerializeJSON(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Swagger spec: %w", err)
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

// extractHostAndBasePath extracts host and base path from environment
func (e *SwaggerExporter) extractHostAndBasePath(environment *EnvironmentWithVariables) (string, string) {
	baseURL := ""
	for _, variable := range environment.Variables {
		if variable.KeyName == "base_url" {
			baseURL = variable.Value
			break
		}
	}

	if baseURL == "" {
		return "", ""
	}

	// Parse URL to extract host and path
	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.TrimPrefix(baseURL, "http://")
	} else if strings.HasPrefix(baseURL, "https://") {
		baseURL = strings.TrimPrefix(baseURL, "https://")
	}

	parts := strings.SplitN(baseURL, "/", 2)
	host := parts[0]
	basePath := ""

	if len(parts) > 1 {
		basePath = "/" + parts[1]
	}

	return host, basePath
}

// buildSecurityDefinitions creates security definitions
func (e *SwaggerExporter) buildSecurityDefinitions() map[string]SwaggerSecurity {
	definitions := make(map[string]SwaggerSecurity)

	definitions["BearerAuth"] = SwaggerSecurity{
		Type:        "apiKey",
		Name:        "Authorization",
		In:          "header",
		Description: "JWT Authorization header using the Bearer scheme",
	}

	return definitions
}

// convertEndpointToOperation converts endpoint to Swagger operation
func (e *SwaggerExporter) convertEndpointToOperation(endpoint *EndpointWithDetails, options *ExportOptions) SwaggerOperation {
	operation := SwaggerOperation{
		Summary:     endpoint.Name,
		Description: endpoint.Description,
		OperationID: e.generateOperationID(endpoint),
		Parameters:  []SwaggerParameter{},
		Responses:   make(map[string]SwaggerResponse),
	}

	// Add tags based on folder
	if endpoint.FolderID != nil {
		// This would need folder information to be passed or looked up
		// For now, we'll use a generic tag
		operation.Tags = []string{"API"}
	}

	// Process parameters
	for _, param := range endpoint.Parameters {
		swaggerParam := e.convertParameter(&param)
		operation.Parameters = append(operation.Parameters, swaggerParam)
	}

	// Process request body
	if endpoint.RequestBody != nil {
		bodyParam := e.convertRequestBody(endpoint.RequestBody)
		operation.Parameters = append(operation.Parameters, bodyParam)
	}

	// Process responses
	if len(endpoint.Responses) > 0 {
		for _, response := range endpoint.Responses {
			statusCode := strconv.Itoa(response.StatusCode)
			swaggerResponse := e.convertResponse(&response)
			operation.Responses[statusCode] = swaggerResponse
		}
	} else {
		// Add default response
		operation.Responses["200"] = SwaggerResponse{
			Description: "Successful response",
		}
	}

	// Add error responses
	operation.Responses["400"] = SwaggerResponse{
		Description: "Bad request",
		Schema: &SwaggerDefinition{
			Ref: "#/definitions/Error",
		},
	}

	return operation
}

// convertParameter converts parameter to Swagger parameter
func (e *SwaggerExporter) convertParameter(param *Parameter) SwaggerParameter {
	swaggerParam := SwaggerParameter{
		Name:        param.Name,
		In:          param.Type,
		Description: param.Description,
		Required:    param.IsRequired,
		Type:        e.convertDataType(param.DataType),
	}

	if param.DefaultValue != "" {
		swaggerParam.Default = param.DefaultValue
	}

	return swaggerParam
}

// convertRequestBody converts request body to Swagger parameter
func (e *SwaggerExporter) convertRequestBody(requestBody *RequestBody) SwaggerParameter {
	param := SwaggerParameter{
		Name:        "body",
		In:          "body",
		Description: requestBody.Description,
		Required:    true,
	}

	// Try to parse schema from JSON if available
	if requestBody.SchemaDefinition != "" {
		var schema SwaggerDefinition
		if err := json.Unmarshal([]byte(requestBody.SchemaDefinition), &schema); err == nil {
			param.Schema = &schema
		}
	}

	if param.Schema == nil {
		// Create basic object schema
		param.Schema = &SwaggerDefinition{
			Type: "object",
		}
	}

	return param
}

// convertResponse converts response to Swagger response
func (e *SwaggerExporter) convertResponse(response *Response) SwaggerResponse {
	swaggerResponse := SwaggerResponse{
		Description: response.Description,
	}

	if response.Description == "" {
		swaggerResponse.Description = http.StatusText(response.StatusCode)
	}

	// Add schema and example if available
	if response.ResponseBody != "" {
		swaggerResponse.Schema = &SwaggerDefinition{
			Type: "object",
		}

		// Try to parse as JSON for example
		var example interface{}
		if err := json.Unmarshal([]byte(response.ResponseBody), &example); err == nil {
			swaggerResponse.Examples = map[string]interface{}{
				response.ContentType: example,
			}
		}
	}

	return swaggerResponse
}

// normalizePath normalizes URL path for Swagger
func (e *SwaggerExporter) normalizePath(url string) string {
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
func (e *SwaggerExporter) generateOperationID(endpoint *EndpointWithDetails) string {
	method := strings.ToLower(endpoint.Method)
	name := strings.ReplaceAll(endpoint.Name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	return method + name
}

// convertDataType converts data type to Swagger type
func (e *SwaggerExporter) convertDataType(dataType string) string {
	switch strings.ToLower(dataType) {
	case "int", "integer", "int64", "int32":
		return "integer"
	case "float", "float64", "float32", "double":
		return "number"
	case "bool", "boolean":
		return "boolean"
	case "array", "slice":
		return "array"
	default:
		return "string"
	}
}

// GetSupportedFormats returns supported formats
func (e *SwaggerExporter) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatSwagger}
}
