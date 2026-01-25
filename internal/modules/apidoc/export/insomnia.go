package export

import (
	"fmt"
	"strings"
	"time"
)

// InsomniaExport represents Insomnia export format
type InsomniaExport struct {
	Type         string             `json:"_type"`
	ExportFormat int                `json:"__export_format"`
	ExportDate   string             `json:"__export_date"`
	ExportSource string             `json:"__export_source"`
	Resources    []InsomniaResource `json:"resources"`
}

// InsomniaResource represents a resource in Insomnia
type InsomniaResource struct {
	ID                              string                 `json:"_id"`
	Type                            string                 `json:"_type"`
	ParentID                        string                 `json:"parentId,omitempty"`
	Modified                        int64                  `json:"modified"`
	Created                         int64                  `json:"created"`
	URL                             string                 `json:"url,omitempty"`
	Name                            string                 `json:"name"`
	Description                     string                 `json:"description,omitempty"`
	Method                          string                 `json:"method,omitempty"`
	Body                            *InsomniaBody          `json:"body,omitempty"`
	Parameters                      []InsomniaParameter    `json:"parameters,omitempty"`
	Headers                         []InsomniaHeader       `json:"headers,omitempty"`
	Authentication                  map[string]interface{} `json:"authentication,omitempty"`
	MetaSortKey                     int64                  `json:"metaSortKey,omitempty"`
	IsPrivate                       bool                   `json:"isPrivate,omitempty"`
	SettingStoreCookies             bool                   `json:"settingStoreCookies,omitempty"`
	SettingSendCookies              bool                   `json:"settingSendCookies,omitempty"`
	SettingDisableRenderRequestBody bool                   `json:"settingDisableRenderRequestBody,omitempty"`
	SettingEncodeUrl                bool                   `json:"settingEncodeUrl,omitempty"`
	SettingRebuildPath              bool                   `json:"settingRebuildPath,omitempty"`
	SettingFollowRedirects          string                 `json:"settingFollowRedirects,omitempty"`
	Scope                           string                 `json:"scope,omitempty"`
	Data                            map[string]interface{} `json:"data,omitempty"`
	DataPropertyOrder               map[string][]string    `json:"dataPropertyOrder,omitempty"`
	Color                           string                 `json:"color,omitempty"`
	Environment                     map[string]interface{} `json:"environment,omitempty"`
	EnvironmentPropertyOrder        interface{}            `json:"environmentPropertyOrder,omitempty"`
	Collapsed                       bool                   `json:"collapsed,omitempty"`
}

// InsomniaBody represents request body
type InsomniaBody struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
}

// InsomniaParameter represents query parameter
type InsomniaParameter struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// InsomniaHeader represents HTTP header
type InsomniaHeader struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// InsomniaExporter handles Insomnia collection export
type InsomniaExporter struct {
	*BaseExporter
}

// NewInsomniaExporter creates a new Insomnia exporter
func NewInsomniaExporter() *InsomniaExporter {
	return &InsomniaExporter{
		BaseExporter: NewBaseExporter(FormatInsomnia),
	}
}

// Export exports collection to Insomnia format
func (e *InsomniaExporter) Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error) {
	if err := e.ValidateOptions(options); err != nil {
		return nil, err
	}

	now := time.Now()
	nowUnix := now.UnixMilli()

	// Build variable map from environment
	var variables map[string]string
	if options.EnvironmentID != nil && collection.Environment != nil {
		variables = e.BuildVariableMap(collection.Environment)
	}

	// Create Insomnia export
	insomniaExport := &InsomniaExport{
		Type:         "export",
		ExportFormat: 4,
		ExportDate:   now.Format(time.RFC3339),
		ExportSource: "huminor.api.doc.system:v1.0.0",
		Resources:    []InsomniaResource{},
	}

	// Create workspace
	workspaceID := e.generateID("wrk")
	workspace := InsomniaResource{
		ID:          workspaceID,
		Type:        "workspace",
		Modified:    nowUnix,
		Created:     nowUnix,
		Name:        collection.Name,
		Description: collection.Description,
		Scope:       "collection",
	}
	insomniaExport.Resources = append(insomniaExport.Resources, workspace)

	// Create environment if available
	if collection.Environment != nil {
		envID := e.generateID("env")
		envData := make(map[string]interface{})
		envOrder := []string{}

		for _, variable := range collection.Environment.Variables {
			envData[variable.KeyName] = variable.Value
			envOrder = append(envOrder, variable.KeyName)
		}

		environment := InsomniaResource{
			ID:       envID,
			Type:     "environment",
			ParentID: workspaceID,
			Modified: nowUnix,
			Created:  nowUnix,
			Name:     collection.Environment.Name,
			Data:     envData,
			DataPropertyOrder: map[string][]string{
				"&": envOrder,
			},
			Color:       "",
			IsPrivate:   false,
			MetaSortKey: nowUnix,
		}
		insomniaExport.Resources = append(insomniaExport.Resources, environment)
	}

	// Create folder map
	folderMap := make(map[int64]string)

	// Create folders
	for _, folder := range collection.Folders {
		folderID := e.generateID("fld")
		folderMap[folder.ID] = folderID

		parentID := workspaceID
		if folder.ParentID != nil {
			if parentFolderID, exists := folderMap[*folder.ParentID]; exists {
				parentID = parentFolderID
			}
		}

		folderResource := InsomniaResource{
			ID:                       folderID,
			Type:                     "request_group",
			ParentID:                 parentID,
			Modified:                 nowUnix,
			Created:                  nowUnix,
			Name:                     folder.Name,
			Description:              folder.Description,
			Environment:              make(map[string]interface{}),
			EnvironmentPropertyOrder: nil,
			MetaSortKey:              nowUnix - int64(folder.SortOrder),
			Collapsed:                true, // Set folder as collapsed by default
		}
		insomniaExport.Resources = append(insomniaExport.Resources, folderResource)
	}

	// Create requests
	for _, endpoint := range collection.Endpoints {
		requestID := e.generateID("req")

		parentID := workspaceID
		if endpoint.FolderID != nil {
			if folderID, exists := folderMap[*endpoint.FolderID]; exists {
				parentID = folderID
			}
		}

		request := e.convertEndpointToInsomniaRequest(&endpoint, requestID, parentID, nowUnix, variables)
		insomniaExport.Resources = append(insomniaExport.Resources, request)
	}

	// Serialize to JSON
	content, err := e.SerializeJSON(insomniaExport)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize Insomnia collection: %w", err)
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

// convertEndpointToInsomniaRequest converts endpoint to Insomnia request
func (e *InsomniaExporter) convertEndpointToInsomniaRequest(endpoint *EndpointWithDetails, requestID, parentID string, timestamp int64, variables map[string]string) InsomniaResource {
	// Build URL
	url := e.SubstituteVariables(endpoint.URL, variables)

	// Convert variable format from {{var}} to {{ _.var }}
	for key := range variables {
		url = strings.ReplaceAll(url, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("{{ _.%s }}", key))
	}

	// Build parameters
	parameters := []InsomniaParameter{}
	for _, param := range endpoint.Parameters {
		if param.Type == "query" {
			insomniaParam := InsomniaParameter{
				Name:        param.Name,
				Value:       param.ExampleValue,
				Description: param.Description,
				Disabled:    !param.IsRequired,
			}
			parameters = append(parameters, insomniaParam)
		}
	}

	// Build headers
	headers := []InsomniaHeader{}
	for _, header := range endpoint.Headers {
		if header.HeaderType == "request" {
			value := e.SubstituteVariables(header.Value, variables)
			// Convert variable format
			for key := range variables {
				value = strings.ReplaceAll(value, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("{{ _.%s }}", key))
			}

			insomniaHeader := InsomniaHeader{
				Name:        header.KeyName,
				Value:       value,
				Description: header.Description,
				Disabled:    !header.IsRequired,
			}
			headers = append(headers, insomniaHeader)
		}
	}

	// Build body
	var body *InsomniaBody
	if endpoint.RequestBody != nil {
		bodyContent := e.SubstituteVariables(endpoint.RequestBody.BodyContent, variables)
		body = &InsomniaBody{
			MimeType: endpoint.RequestBody.ContentType,
			Text:     bodyContent,
		}
	}

	return InsomniaResource{
		ID:                              requestID,
		Type:                            "request",
		ParentID:                        parentID,
		Modified:                        timestamp,
		Created:                         timestamp,
		URL:                             url,
		Name:                            endpoint.Name,
		Description:                     endpoint.Description,
		Method:                          endpoint.Method,
		Body:                            body,
		Parameters:                      parameters,
		Headers:                         headers,
		Authentication:                  make(map[string]interface{}),
		MetaSortKey:                     timestamp - int64(endpoint.SortOrder),
		IsPrivate:                       false,
		SettingStoreCookies:             true,
		SettingSendCookies:              true,
		SettingDisableRenderRequestBody: false,
		SettingEncodeUrl:                true,
		SettingRebuildPath:              true,
		SettingFollowRedirects:          "global",
	}
}

// generateID generates Insomnia-style ID
func (e *InsomniaExporter) generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

// GetSupportedFormats returns supported formats
func (e *InsomniaExporter) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatInsomnia}
}
