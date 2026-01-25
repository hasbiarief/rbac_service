package export

import (
	"fmt"
	"strings"
	"time"
)

// ApidogCollection represents Apidog collection format
type ApidogCollection struct {
	ApidogVersion string              `json:"apidogVersion"`
	Info          ApidogInfo          `json:"info"`
	Servers       []ApidogServer      `json:"servers,omitempty"`
	Folders       []ApidogFolder      `json:"folders,omitempty"`
	APIs          []ApidogAPI         `json:"apis"`
	Environments  []ApidogEnvironment `json:"environments,omitempty"`
	Schemas       []ApidogSchema      `json:"schemas,omitempty"`
}

// ApidogInfo contains collection information
type ApidogInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// ApidogServer represents a server
type ApidogServer struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// ApidogFolder represents a folder
type ApidogFolder struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
	Sort        int    `json:"sort"`
}

// ApidogAPI represents an API endpoint
type ApidogAPI struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	FolderID    string             `json:"folderId,omitempty"`
	Method      string             `json:"method"`
	Path        string             `json:"path"`
	Parameters  ApidogParameters   `json:"parameters,omitempty"`
	RequestBody *ApidogRequestBody `json:"requestBody,omitempty"`
	Responses   []ApidogResponse   `json:"responses,omitempty"`
	Sort        int                `json:"sort"`
	Tags        []string           `json:"tags,omitempty"`
}

// ApidogParameters contains all parameter types
type ApidogParameters struct {
	Query  []ApidogParameter `json:"query,omitempty"`
	Header []ApidogParameter `json:"header,omitempty"`
	Path   []ApidogParameter `json:"path,omitempty"`
	Cookie []ApidogParameter `json:"cookie,omitempty"`
}

// ApidogParameter represents a parameter
type ApidogParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// ApidogRequestBody represents request body
type ApidogRequestBody struct {
	Type       string                `json:"type"`
	JSONSchema *ApidogJSONSchema     `json:"jsonSchema,omitempty"`
	Example    interface{}           `json:"example,omitempty"`
	FormData   []ApidogFormParameter `json:"formData,omitempty"`
}

// ApidogJSONSchema represents JSON schema
type ApidogJSONSchema struct {
	Type                 string                      `json:"type,omitempty"`
	Properties           map[string]ApidogJSONSchema `json:"properties,omitempty"`
	Required             []string                    `json:"required,omitempty"`
	Items                *ApidogJSONSchema           `json:"items,omitempty"`
	AdditionalProperties interface{}                 `json:"additionalProperties,omitempty"`
	Description          string                      `json:"description,omitempty"`
	Example              interface{}                 `json:"example,omitempty"`
	Format               string                      `json:"format,omitempty"`
	Enum                 []interface{}               `json:"enum,omitempty"`
	Default              interface{}                 `json:"default,omitempty"`
}

// ApidogFormParameter represents form parameter
type ApidogFormParameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// ApidogResponse represents a response
type ApidogResponse struct {
	StatusCode  int               `json:"statusCode"`
	Description string            `json:"description"`
	ContentType string            `json:"contentType"`
	JSONSchema  *ApidogJSONSchema `json:"jsonSchema,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Headers     []ApidogParameter `json:"headers,omitempty"`
}

// ApidogEnvironment represents an environment
type ApidogEnvironment struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Variables []ApidogEnvironmentVar `json:"variables"`
}

// ApidogEnvironmentVar represents environment variable
type ApidogEnvironmentVar struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// ApidogSchema represents reusable schema
type ApidogSchema struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	JSONSchema *ApidogJSONSchema `json:"jsonSchema"`
}

// ApidogExporter handles Apidog collection export
type ApidogExporter struct {
	*BaseExporter
}

// NewApidogExporter creates a new Apidog exporter
func NewApidogExporter() *ApidogExporter {
	return &ApidogExporter{
		BaseExporter: NewBaseExporter(FormatApidog),
	}
}

// Export exports collection to Apidog format
func (e *ApidogExporter) Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error) {
	if err := e.ValidateOptions(options); err != nil {
		return nil, err
	}

	// Build variable map from environment
	var variables map[string]string
	if options.EnvironmentID != nil && collection.Environment != nil {
		variables = e.BuildVariableMap(collection.Environment)
	}

	// Create Apidog collection
	apidogCollection := &ApidogCollection{
		ApidogVersion: "2.0.0",
		Info: ApidogInfo{
			Name:        collection.Name,
			Description: collection.Description,
			Version:     collection.Version,
		},
		Servers:      []ApidogServer{},
		Folders:      []ApidogFolder{},
		APIs:         []ApidogAPI{},
		Environments: []ApidogEnvironment{},
		Schemas:      []ApidogSchema{},
	}

	// Add servers from environment
	if collection.Environment != nil {
		servers := e.buildServers(collection.Environment)
		apidogCollection.Servers = servers
	}

	// Add environment
	if collection.Environment != nil {
		environment := e.convertEnvironment(collection.Environment)
		apidogCollection.Environments = append(apidogCollection.Environments, environment)
	}

	// Create folder map
	folderMap := make(map[int64]string)

	// Create folders
	for i, folder := range collection.Folders {
		folderID := fmt.Sprintf("folder_%d", folder.ID)
		folderMap[folder.ID] = folderID

		apidogFolder := ApidogFolder{
			ID:          folderID,
			Name:        folder.Name,
			Description: folder.Description,
			Sort:        i + 1,
		}

		if folder.ParentID != nil {
			if parentID, exists := folderMap[*folder.ParentID]; exists {
				apidogFolder.ParentID = parentID
			}
		}

		apidogCollection.Folders = append(apidogCollection.Folders, apidogFolder)
	}

	// Create APIs
	for i, endpoint := range collection.Endpoints {
		api := e.convertEndpointToAPI(&endpoint, folderMap, variables, i+1)
		apidogCollection.APIs = append(apidogCollection.APIs, api)
	}

	// Serialize to JSON
	content, err := e.SerializeJSON(apidogCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Apidog collection: %w", err)
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
func (e *ApidogExporter) buildServers(environment *EnvironmentWithVariables) []ApidogServer {
	servers := []ApidogServer{}

	// Find base_url variable
	baseURL := ""
	for _, variable := range environment.Variables {
		if variable.KeyName == "base_url" {
			baseURL = variable.Value
			break
		}
	}

	if baseURL != "" {
		server := ApidogServer{
			URL:         baseURL,
			Description: fmt.Sprintf("%s Server", environment.Name),
		}
		servers = append(servers, server)
	}

	return servers
}

// convertEnvironment converts environment to Apidog format
func (e *ApidogExporter) convertEnvironment(environment *EnvironmentWithVariables) ApidogEnvironment {
	apidogEnv := ApidogEnvironment{
		ID:        fmt.Sprintf("env_%d", environment.ID),
		Name:      environment.Name,
		Variables: []ApidogEnvironmentVar{},
	}

	for _, variable := range environment.Variables {
		envVar := ApidogEnvironmentVar{
			Key:         variable.KeyName,
			Value:       variable.Value,
			Description: variable.Description,
		}
		apidogEnv.Variables = append(apidogEnv.Variables, envVar)
	}

	return apidogEnv
}

// convertEndpointToAPI converts endpoint to Apidog API
func (e *ApidogExporter) convertEndpointToAPI(endpoint *EndpointWithDetails, folderMap map[int64]string, variables map[string]string, sort int) ApidogAPI {
	api := ApidogAPI{
		ID:          fmt.Sprintf("api_%d", endpoint.ID),
		Name:        endpoint.Name,
		Description: endpoint.Description,
		Method:      endpoint.Method,
		Path:        e.SubstituteVariables(endpoint.URL, variables),
		Sort:        sort,
		Parameters:  ApidogParameters{},
		Responses:   []ApidogResponse{},
		Tags:        []string{},
	}

	// Set folder ID
	if endpoint.FolderID != nil {
		if folderID, exists := folderMap[*endpoint.FolderID]; exists {
			api.FolderID = folderID
		}
	}

	// Process parameters
	queryParams := []ApidogParameter{}
	headerParams := []ApidogParameter{}
	pathParams := []ApidogParameter{}

	for _, param := range endpoint.Parameters {
		apidogParam := ApidogParameter{
			Name:        param.Name,
			Type:        e.convertDataType(param.DataType),
			Required:    param.IsRequired,
			Description: param.Description,
		}

		if param.ExampleValue != "" {
			apidogParam.Example = param.ExampleValue
		}
		if param.DefaultValue != "" {
			apidogParam.Default = param.DefaultValue
		}

		switch param.Type {
		case "query":
			queryParams = append(queryParams, apidogParam)
		case "path":
			pathParams = append(pathParams, apidogParam)
		}
	}

	// Process headers
	for _, header := range endpoint.Headers {
		if header.HeaderType == "request" {
			headerParam := ApidogParameter{
				Name:        header.KeyName,
				Type:        "string",
				Required:    header.IsRequired,
				Description: header.Description,
				Example:     e.SubstituteVariables(header.Value, variables),
			}
			headerParams = append(headerParams, headerParam)
		}
	}

	api.Parameters.Query = queryParams
	api.Parameters.Header = headerParams
	api.Parameters.Path = pathParams

	// Process request body
	if endpoint.RequestBody != nil {
		requestBody := e.convertRequestBody(endpoint.RequestBody, variables)
		api.RequestBody = requestBody
	}

	// Process responses
	for _, response := range endpoint.Responses {
		apidogResponse := e.convertResponse(&response)
		api.Responses = append(api.Responses, apidogResponse)
	}

	return api
}

// convertRequestBody converts request body to Apidog format
func (e *ApidogExporter) convertRequestBody(requestBody *RequestBody, variables map[string]string) *ApidogRequestBody {
	apidogBody := &ApidogRequestBody{
		Type: requestBody.ContentType,
	}

	// Handle JSON content
	if strings.Contains(requestBody.ContentType, "json") {
		// Try to create schema from schema definition
		if requestBody.SchemaDefinition != "" {
			// This would need proper JSON schema parsing
			apidogBody.JSONSchema = &ApidogJSONSchema{
				Type: "object",
			}
		}

		// Add example
		if requestBody.BodyContent != "" {
			bodyContent := e.SubstituteVariables(requestBody.BodyContent, variables)
			apidogBody.Example = bodyContent
		}
	}

	return apidogBody
}

// convertResponse converts response to Apidog format
func (e *ApidogExporter) convertResponse(response *Response) ApidogResponse {
	apidogResponse := ApidogResponse{
		StatusCode:  response.StatusCode,
		Description: response.Description,
		ContentType: response.ContentType,
		Headers:     []ApidogParameter{},
	}

	// Add example if available
	if response.ResponseBody != "" {
		apidogResponse.Example = response.ResponseBody
	}

	// Create basic schema for JSON responses
	if strings.Contains(response.ContentType, "json") {
		apidogResponse.JSONSchema = &ApidogJSONSchema{
			Type: "object",
		}
	}

	return apidogResponse
}

// convertDataType converts data type to Apidog type
func (e *ApidogExporter) convertDataType(dataType string) string {
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
func (e *ApidogExporter) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatApidog}
}
