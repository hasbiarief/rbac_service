package export

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PostmanCollection represents a Postman collection v2.1
type PostmanCollection struct {
	Info     PostmanInfo       `json:"info"`
	Item     []PostmanItem     `json:"item"`
	Auth     *PostmanAuth      `json:"auth,omitempty"`
	Event    []PostmanEvent    `json:"event,omitempty"`
	Variable []PostmanVariable `json:"variable,omitempty"`
}

// PostmanInfo contains collection metadata
type PostmanInfo struct {
	PostmanID   string `json:"_postman_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
	ExporterID  string `json:"_exporter_id,omitempty"`
	Version     string `json:"version,omitempty"`
}

// PostmanItem represents a folder or request
type PostmanItem struct {
	Name               string            `json:"name"`
	Description        string            `json:"description,omitempty"`
	Item               []PostmanItem     `json:"item,omitempty"`
	Request            *PostmanRequest   `json:"request,omitempty"`
	Response           []PostmanResponse `json:"response,omitempty"`
	Event              []PostmanEvent    `json:"event,omitempty"`
	PostmanIsSubFolder bool              `json:"_postman_isSubFolder,omitempty"`
	PostmanFolderState string            `json:"_postman_folderState,omitempty"`
}

// PostmanRequest represents an HTTP request
type PostmanRequest struct {
	Method      string          `json:"method"`
	Header      []PostmanHeader `json:"header,omitempty"`
	Body        *PostmanBody    `json:"body,omitempty"`
	URL         PostmanURL      `json:"url"`
	Description string          `json:"description,omitempty"`
	Auth        *PostmanAuth    `json:"auth,omitempty"`
}

// PostmanURL represents a request URL
type PostmanURL struct {
	Raw      string               `json:"raw"`
	Protocol string               `json:"protocol,omitempty"`
	Host     []string             `json:"host,omitempty"`
	Path     []string             `json:"path,omitempty"`
	Port     string               `json:"port,omitempty"`
	Query    []PostmanQueryParam  `json:"query,omitempty"`
	Variable []PostmanURLVariable `json:"variable,omitempty"`
}

// PostmanHeader represents an HTTP header
type PostmanHeader struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// PostmanBody represents request body
type PostmanBody struct {
	Mode       string                 `json:"mode"`
	Raw        string                 `json:"raw,omitempty"`
	URLEncoded []PostmanFormParameter `json:"urlencoded,omitempty"`
	FormData   []PostmanFormParameter `json:"formdata,omitempty"`
	File       *PostmanFile           `json:"file,omitempty"`
	Options    *PostmanBodyOptions    `json:"options,omitempty"`
}

// PostmanBodyOptions contains body options
type PostmanBodyOptions struct {
	Raw *PostmanRawOptions `json:"raw,omitempty"`
}

// PostmanRawOptions contains raw body options
type PostmanRawOptions struct {
	Language string `json:"language,omitempty"`
}

// PostmanFormParameter represents form parameter
type PostmanFormParameter struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// PostmanFile represents file upload
type PostmanFile struct {
	Src string `json:"src"`
}

// PostmanQueryParam represents query parameter
type PostmanQueryParam struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// PostmanURLVariable represents URL variable
type PostmanURLVariable struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// PostmanResponse represents example response
type PostmanResponse struct {
	Name            string          `json:"name"`
	OriginalRequest PostmanRequest  `json:"originalRequest"`
	Status          string          `json:"status"`
	Code            int             `json:"code"`
	Header          []PostmanHeader `json:"header,omitempty"`
	Cookie          []PostmanCookie `json:"cookie,omitempty"`
	Body            string          `json:"body,omitempty"`
}

// PostmanCookie represents HTTP cookie
type PostmanCookie struct {
	Domain   string `json:"domain"`
	Expires  string `json:"expires,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Secure   bool   `json:"secure,omitempty"`
	Value    string `json:"value"`
}

// PostmanEvent represents pre-request or test script
type PostmanEvent struct {
	Listen string        `json:"listen"`
	Script PostmanScript `json:"script"`
}

// PostmanScript represents JavaScript code
type PostmanScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// PostmanAuth represents authentication
type PostmanAuth struct {
	Type   string                `json:"type"`
	Bearer []PostmanAuthProperty `json:"bearer,omitempty"`
	Basic  []PostmanAuthProperty `json:"basic,omitempty"`
	APIKey []PostmanAuthProperty `json:"apikey,omitempty"`
}

// PostmanAuthProperty represents auth property
type PostmanAuthProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// PostmanVariable represents collection variable
type PostmanVariable struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// PostmanExporter handles Postman collection export
type PostmanExporter struct {
	*BaseExporter
}

// NewPostmanExporter creates a new Postman exporter
func NewPostmanExporter() *PostmanExporter {
	return &PostmanExporter{
		BaseExporter: NewBaseExporter(FormatPostman),
	}
}

// Export exports collection to Postman format
func (e *PostmanExporter) Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error) {
	if err := e.ValidateOptions(options); err != nil {
		return nil, err
	}

	// Build variable map from environment
	var variables map[string]string
	if options.EnvironmentID != nil && collection.Environment != nil {
		variables = e.BuildVariableMap(collection.Environment)
	}

	// Create Postman collection
	postmanCollection := &PostmanCollection{
		Info: PostmanInfo{
			PostmanID:   uuid.New().String(),
			Name:        collection.Name,
			Description: collection.Description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
			ExporterID:  "huminor-api-doc-system",
			Version:     collection.Version,
		},
		Item:     []PostmanItem{},
		Variable: []PostmanVariable{},
	}

	// Add environment variables
	if collection.Environment != nil {
		for _, variable := range collection.Environment.Variables {
			postmanVar := PostmanVariable{
				Key:         variable.KeyName,
				Value:       variable.Value,
				Type:        "string",
				Description: variable.Description,
			}
			postmanCollection.Variable = append(postmanCollection.Variable, postmanVar)
		}
	}

	// Process folders and endpoints
	folderMap := make(map[int64]*PostmanItem)

	// Create folder items
	for _, folder := range collection.Folders {
		folderItem := &PostmanItem{
			Name:               folder.Name,
			Description:        folder.Description,
			Item:               []PostmanItem{},
			PostmanIsSubFolder: true,
			PostmanFolderState: "collapsed",
		}
		folderMap[folder.ID] = folderItem

		// Add to root if no parent, otherwise will be added to parent later
		if folder.ParentID == nil {
			postmanCollection.Item = append(postmanCollection.Item, *folderItem)
		}
	}

	// Handle nested folders
	for _, folder := range collection.Folders {
		if folder.ParentID != nil {
			if parentFolder, exists := folderMap[*folder.ParentID]; exists {
				if folderItem, exists := folderMap[folder.ID]; exists {
					parentFolder.Item = append(parentFolder.Item, *folderItem)
				}
			}
		}
	}

	// Process endpoints
	for _, endpoint := range collection.Endpoints {
		requestItem := e.convertEndpointToPostmanItem(&endpoint, variables, options)

		// Add to appropriate folder or root
		if endpoint.FolderID != nil {
			if folder, exists := folderMap[*endpoint.FolderID]; exists {
				folder.Item = append(folder.Item, requestItem)
			}
		} else {
			postmanCollection.Item = append(postmanCollection.Item, requestItem)
		}
	}

	// Serialize to JSON
	content, err := e.SerializeJSON(postmanCollection)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Postman collection: %w", err)
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

// convertEndpointToPostmanItem converts an endpoint to Postman item
func (e *PostmanExporter) convertEndpointToPostmanItem(endpoint *EndpointWithDetails, variables map[string]string, options *ExportOptions) PostmanItem {
	// Build URL
	rawURL := e.SubstituteVariables(endpoint.URL, variables)
	postmanURL := e.buildPostmanURL(rawURL, endpoint.Parameters)

	// Build headers
	headers := []PostmanHeader{}
	for _, header := range endpoint.Headers {
		if header.HeaderType == "request" {
			postmanHeader := PostmanHeader{
				Key:         header.KeyName,
				Value:       e.SubstituteVariables(header.Value, variables),
				Description: header.Description,
				Disabled:    !header.IsRequired,
			}
			headers = append(headers, postmanHeader)
		}
	}

	// Build request body
	var body *PostmanBody
	if endpoint.RequestBody != nil {
		body = &PostmanBody{
			Mode: "raw",
			Raw:  e.SubstituteVariables(endpoint.RequestBody.BodyContent, variables),
			Options: &PostmanBodyOptions{
				Raw: &PostmanRawOptions{
					Language: e.getLanguageFromContentType(endpoint.RequestBody.ContentType),
				},
			},
		}
	}

	// Build request
	request := &PostmanRequest{
		Method:      endpoint.Method,
		Header:      headers,
		Body:        body,
		URL:         postmanURL,
		Description: endpoint.Description,
	}

	// Build responses
	responses := []PostmanResponse{}
	for _, response := range endpoint.Responses {
		postmanResponse := PostmanResponse{
			Name:   fmt.Sprintf("%d %s", response.StatusCode, response.StatusText),
			Status: response.StatusText,
			Code:   response.StatusCode,
			Body:   response.ResponseBody,
			Header: []PostmanHeader{},
		}

		// Add response headers
		for _, header := range endpoint.Headers {
			if header.HeaderType == "response" {
				responseHeader := PostmanHeader{
					Key:   header.KeyName,
					Value: header.Value,
				}
				postmanResponse.Header = append(postmanResponse.Header, responseHeader)
			}
		}

		responses = append(responses, postmanResponse)
	}

	// Build events (scripts)
	events := []PostmanEvent{}
	if options.IncludeTests || options.IncludePreRequest {
		for _, test := range endpoint.Tests {
			if test.TestType == "pre_request" && options.IncludePreRequest {
				event := PostmanEvent{
					Listen: "prerequest",
					Script: PostmanScript{
						Type: "text/javascript",
						Exec: strings.Split(test.ScriptContent, "\n"),
					},
				}
				events = append(events, event)
			} else if test.TestType == "test" && options.IncludeTests {
				event := PostmanEvent{
					Listen: "test",
					Script: PostmanScript{
						Type: "text/javascript",
						Exec: strings.Split(test.ScriptContent, "\n"),
					},
				}
				events = append(events, event)
			}
		}
	}

	return PostmanItem{
		Name:        endpoint.Name,
		Description: endpoint.Description,
		Request:     request,
		Response:    responses,
		Event:       events,
	}
}

// buildPostmanURL builds Postman URL structure
func (e *PostmanExporter) buildPostmanURL(rawURL string, parameters []Parameter) PostmanURL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// If parsing fails, return simple structure
		return PostmanURL{
			Raw: rawURL,
		}
	}

	// Build host array
	host := []string{}
	if parsedURL.Host != "" {
		host = strings.Split(parsedURL.Host, ".")
	}

	// Build path array
	path := []string{}
	if parsedURL.Path != "" {
		pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
		for _, part := range pathParts {
			if part != "" {
				path = append(path, part)
			}
		}
	}

	// Build query parameters
	query := []PostmanQueryParam{}
	for _, param := range parameters {
		if param.Type == "query" {
			queryParam := PostmanQueryParam{
				Key:         param.Name,
				Value:       param.ExampleValue,
				Description: param.Description,
				Disabled:    !param.IsRequired,
			}
			query = append(query, queryParam)
		}
	}

	return PostmanURL{
		Raw:      rawURL,
		Protocol: parsedURL.Scheme,
		Host:     host,
		Path:     path,
		Port:     parsedURL.Port(),
		Query:    query,
	}
}

// getLanguageFromContentType returns language for syntax highlighting
func (e *PostmanExporter) getLanguageFromContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "json"):
		return "json"
	case strings.Contains(contentType, "xml"):
		return "xml"
	case strings.Contains(contentType, "html"):
		return "html"
	case strings.Contains(contentType, "javascript"):
		return "javascript"
	default:
		return "text"
	}
}

// GetSupportedFormats returns supported formats
func (e *PostmanExporter) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatPostman}
}
