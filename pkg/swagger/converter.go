package swagger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Converter converts Postman collections to Swagger annotations
type Converter interface {
	// Convert converts a Postman collection to Swagger annotations
	Convert(postmanFile string, outputDir string) (*ConversionReport, error)

	// ConvertItem converts a single Postman item to Swagger annotation
	ConvertItem(item PostmanItem) (string, error)
}

// PostmanCollection represents a Postman collection
type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanInfo contains collection metadata
type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

// PostmanItem represents a Postman request or folder
type PostmanItem struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Request     *PostmanRequest `json:"request,omitempty"`
	Item        []PostmanItem   `json:"item,omitempty"`
}

// PostmanRequest represents a Postman HTTP request
type PostmanRequest struct {
	Method      string          `json:"method"`
	Header      []PostmanHeader `json:"header"`
	Body        *PostmanBody    `json:"body,omitempty"`
	URL         PostmanURL      `json:"url"`
	Description string          `json:"description,omitempty"`
}

// PostmanHeader represents a request header
type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostmanBody represents a request body
type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw,omitempty"`
}

// PostmanURL represents a request URL
type PostmanURL struct {
	Raw  string   `json:"raw"`
	Path []string `json:"path"`
}

// ConversionReport contains the results of a conversion
type ConversionReport struct {
	TotalEndpoints     int
	ConvertedEndpoints int
	FailedEndpoints    []FailedEndpoint
	GeneratedFiles     []string
}

// FailedEndpoint represents an endpoint that failed to convert
type FailedEndpoint struct {
	Name   string
	Reason string
}

// converter implements Converter
type converter struct{}

// NewConverter creates a new Postman to Swagger converter
func NewConverter() Converter {
	return &converter{}
}

// Convert converts a Postman collection to Swagger annotations
func (c *converter) Convert(postmanFile string, outputDir string) (*ConversionReport, error) {
	report := &ConversionReport{
		FailedEndpoints: []FailedEndpoint{},
		GeneratedFiles:  []string{},
	}

	// Read Postman collection file
	data, err := os.ReadFile(postmanFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read Postman file: %w", err)
	}

	// Parse Postman collection
	var collection PostmanCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse Postman collection: %w", err)
	}

	// Group endpoints by module
	moduleEndpoints := c.groupByModule(collection.Item)

	// Generate swagger.go files for each module
	for moduleName, items := range moduleEndpoints {
		annotations := []string{}

		for _, item := range items {
			report.TotalEndpoints++

			annotation, err := c.ConvertItem(item)
			if err != nil {
				report.FailedEndpoints = append(report.FailedEndpoints, FailedEndpoint{
					Name:   item.Name,
					Reason: err.Error(),
				})
				continue
			}

			annotations = append(annotations, annotation)
			report.ConvertedEndpoints++
		}

		// Generate file for this module
		if len(annotations) > 0 {
			filePath, err := c.generateAnnotationFile(moduleName, annotations, outputDir)
			if err != nil {
				return report, fmt.Errorf("failed to generate file for module %s: %w", moduleName, err)
			}
			report.GeneratedFiles = append(report.GeneratedFiles, filePath)
		}
	}

	return report, nil
}

// groupByModule groups Postman items by module based on path
func (c *converter) groupByModule(items []PostmanItem) map[string][]PostmanItem {
	moduleEndpoints := make(map[string][]PostmanItem)

	var processItems func([]PostmanItem, string)
	processItems = func(items []PostmanItem, parentPath string) {
		for _, item := range items {
			// If item has sub-items (folder), process recursively
			if len(item.Item) > 0 {
				// Use folder name as potential module name
				processItems(item.Item, item.Name)
				continue
			}

			// If item has a request, it's an endpoint
			if item.Request != nil {
				moduleName := c.extractModuleName(item.Request.URL.Path, parentPath)
				moduleEndpoints[moduleName] = append(moduleEndpoints[moduleName], item)
			}
		}
	}

	processItems(items, "")
	return moduleEndpoints
}

// extractModuleName extracts module name from URL path
func (c *converter) extractModuleName(pathSegments []string, folderName string) string {
	// Try to extract from path first
	// Expected format: /api/v1/{module}/...
	if len(pathSegments) >= 3 {
		// Skip "api" and version, use next segment as module
		for i, segment := range pathSegments {
			if segment == "api" && i+2 < len(pathSegments) {
				return pathSegments[i+2]
			}
		}
	}

	// If path doesn't match expected format, use folder name
	if folderName != "" {
		return strings.ToLower(folderName)
	}

	// Default to "general" if we can't determine module
	return "general"
}

// ConvertItem converts a single Postman item to Swagger annotation
func (c *converter) ConvertItem(item PostmanItem) (string, error) {
	if item.Request == nil {
		return "", fmt.Errorf("item has no request")
	}

	req := item.Request
	var builder strings.Builder

	// Add summary
	builder.WriteString(fmt.Sprintf("// @Summary      %s\n", item.Name))

	// Add description if available
	description := item.Description
	if description == "" && req.Description != "" {
		description = req.Description
	}
	if description != "" {
		// Clean up description (remove newlines, limit length)
		description = strings.ReplaceAll(description, "\n", " ")
		if len(description) > 200 {
			description = description[:197] + "..."
		}
		builder.WriteString(fmt.Sprintf("// @Description  %s\n", description))
	}

	// Add tags (use module name from path)
	tag := c.extractTag(req.URL.Path)
	builder.WriteString(fmt.Sprintf("// @Tags         %s\n", tag))

	// Add Accept and Produce based on Content-Type header
	contentType := c.extractContentType(req.Header)
	if contentType != "" {
		builder.WriteString(fmt.Sprintf("// @Accept       %s\n", c.convertContentType(contentType)))
	} else {
		builder.WriteString("// @Accept       json\n")
	}
	builder.WriteString("// @Produce      json\n")

	// Add parameters
	params := c.extractParameters(req)
	for _, param := range params {
		builder.WriteString(param + "\n")
	}

	// Add request body if present
	if req.Body != nil && req.Body.Raw != "" {
		bodyType := c.inferBodyType(req.Body.Raw)
		builder.WriteString(fmt.Sprintf("// @Param        request  body      %s  true  \"Request body\"\n", bodyType))
	}

	// Add success response
	builder.WriteString("// @Success      200  {object}  response.Response\n")

	// Add common failure responses
	builder.WriteString("// @Failure      400  {object}  response.Response\n")
	builder.WriteString("// @Failure      401  {object}  response.Response\n")
	builder.WriteString("// @Failure      500  {object}  response.Response\n")

	// Add security if Authorization header is present
	if c.hasAuthHeader(req.Header) {
		builder.WriteString("// @Security     BearerAuth\n")
	}

	// Add router
	path := c.buildPath(req.URL.Path)
	method := strings.ToLower(req.Method)
	builder.WriteString(fmt.Sprintf("// @Router       %s [%s]\n", path, method))

	return builder.String(), nil
}

// extractTag extracts tag name from URL path
func (c *converter) extractTag(pathSegments []string) string {
	// Try to extract module name from path
	for i, segment := range pathSegments {
		if segment == "api" && i+2 < len(pathSegments) {
			module := pathSegments[i+2]
			// Capitalize first letter
			return strings.ToUpper(module[:1]) + module[1:]
		}
	}
	return "General"
}

// extractContentType extracts Content-Type from headers
func (c *converter) extractContentType(headers []PostmanHeader) string {
	for _, header := range headers {
		if strings.ToLower(header.Key) == "content-type" {
			return header.Value
		}
	}
	return ""
}

// convertContentType converts MIME type to Swagger format
func (c *converter) convertContentType(mimeType string) string {
	// Extract base type (remove charset, etc.)
	parts := strings.Split(mimeType, ";")
	baseType := strings.TrimSpace(parts[0])

	switch baseType {
	case "application/json":
		return "json"
	case "application/xml":
		return "xml"
	case "application/x-www-form-urlencoded":
		return "x-www-form-urlencoded"
	case "multipart/form-data":
		return "mpfd"
	case "text/plain":
		return "plain"
	default:
		return "json"
	}
}

// extractParameters extracts query and path parameters from URL
func (c *converter) extractParameters(req *PostmanRequest) []string {
	var params []string

	// Extract path parameters (e.g., :id, {id})
	pathParams := c.extractPathParams(req.URL.Path)
	for _, param := range pathParams {
		params = append(params, fmt.Sprintf("// @Param        %s  path      string  true  \"%s parameter\"", param, param))
	}

	// Extract query parameters from URL
	// Note: Postman URL structure may include query in Raw field
	if strings.Contains(req.URL.Raw, "?") {
		queryParams := c.extractQueryParams(req.URL.Raw)
		for _, param := range queryParams {
			params = append(params, fmt.Sprintf("// @Param        %s  query     string  false  \"%s parameter\"", param, param))
		}
	}

	return params
}

// extractPathParams extracts path parameters from URL path
func (c *converter) extractPathParams(pathSegments []string) []string {
	var params []string
	paramRegex := regexp.MustCompile(`^:(.+)$|^\{(.+)\}$`)

	for _, segment := range pathSegments {
		matches := paramRegex.FindStringSubmatch(segment)
		if len(matches) > 0 {
			// Extract parameter name (either from :param or {param})
			paramName := matches[1]
			if paramName == "" {
				paramName = matches[2]
			}
			params = append(params, paramName)
		}
	}

	return params
}

// extractQueryParams extracts query parameter names from URL
func (c *converter) extractQueryParams(rawURL string) []string {
	var params []string

	// Find query string
	parts := strings.Split(rawURL, "?")
	if len(parts) < 2 {
		return params
	}

	queryString := parts[1]
	// Split by & to get individual parameters
	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		// Split by = to get parameter name
		paramParts := strings.Split(pair, "=")
		if len(paramParts) > 0 {
			paramName := paramParts[0]
			if paramName != "" {
				params = append(params, paramName)
			}
		}
	}

	return params
}

// hasAuthHeader checks if Authorization header is present
func (c *converter) hasAuthHeader(headers []PostmanHeader) bool {
	for _, header := range headers {
		if strings.ToLower(header.Key) == "authorization" {
			return true
		}
	}
	return false
}

// buildPath builds the path string from path segments
func (c *converter) buildPath(pathSegments []string) string {
	if len(pathSegments) == 0 {
		return "/"
	}

	var parts []string
	for _, segment := range pathSegments {
		if segment == "" {
			continue
		}
		// Convert :param to {param} format for Swagger
		if strings.HasPrefix(segment, ":") {
			segment = "{" + segment[1:] + "}"
		}
		parts = append(parts, segment)
	}

	return "/" + strings.Join(parts, "/")
}

// inferBodyType infers the body type from raw JSON
func (c *converter) inferBodyType(rawBody string) string {
	// Try to parse as JSON to determine structure
	var data map[string]any
	if err := json.Unmarshal([]byte(rawBody), &data); err == nil {
		// Successfully parsed as JSON object
		// Try to infer type from common patterns
		if _, ok := data["email"]; ok {
			if _, ok := data["password"]; ok {
				return "auth.LoginRequest"
			}
		}
		// Default to generic object
		return "object"
	}

	// Try array
	var arr []any
	if err := json.Unmarshal([]byte(rawBody), &arr); err == nil {
		return "array"
	}

	// Default to object
	return "object"
}

// generateAnnotationFile generates a swagger.go file for a module
func (c *converter) generateAnnotationFile(moduleName string, annotations []string, outputDir string) (string, error) {
	// Create module directory structure
	moduleDir := filepath.Join(outputDir, "internal", "modules", moduleName, "docs")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file path
	filePath := filepath.Join(moduleDir, "swagger.go")

	// Build file content
	var content strings.Builder
	content.WriteString("package docs\n\n")
	content.WriteString(fmt.Sprintf("// Package docs contains Swagger annotations for the %s module.\n", moduleName))
	content.WriteString("// This file was generated from Postman collection.\n\n")

	// Add all annotations
	for i, annotation := range annotations {
		if i > 0 {
			content.WriteString("\n")
		}
		content.WriteString(annotation)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}
