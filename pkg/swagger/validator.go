package swagger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Validator validates Swagger annotations
type Validator interface {
	// ValidateAnnotations validates all annotation files
	ValidateAnnotations() ([]ValidationError, error)

	// ValidateFile validates a single annotation file
	ValidateFile(filePath string) ([]ValidationError, error)

	// ValidateSyntax validates annotation syntax
	ValidateSyntax(filePath string) ([]ValidationError, error)

	// ValidateRequiredFields validates required fields in annotations
	ValidateRequiredFields(filePath string) ([]ValidationError, error)

	// DetectConflicts detects conflicts like duplicate paths and operation IDs
	DetectConflicts() ([]ValidationError, error)

	// ValidateExampleSchema validates that examples conform to their schemas
	ValidateExampleSchema(filePath string) ([]ValidationError, error)
}

// validator implements Validator
type validator struct {
	searchDir string
	// Track all endpoints for conflict detection
	endpoints    map[string][]*endpointInfo // key: "METHOD /path"
	operationIDs map[string][]*endpointInfo // key: operationID
}

// endpointInfo stores information about an endpoint for validation
type endpointInfo struct {
	file        string
	line        int
	method      string
	path        string
	operationID string
}

// NewValidator creates a new annotation validator
func NewValidator(searchDir string) Validator {
	if searchDir == "" {
		searchDir = "./"
	}
	return &validator{
		searchDir:    searchDir,
		endpoints:    make(map[string][]*endpointInfo),
		operationIDs: make(map[string][]*endpointInfo),
	}
}

// ValidateAnnotations validates all annotation files
func (v *validator) ValidateAnnotations() ([]ValidationError, error) {
	var allErrors []ValidationError

	// Find all annotation files
	files, err := v.findAnnotationFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find annotation files: %w", err)
	}

	// Validate each file
	for _, file := range files {
		errors, err := v.ValidateFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to validate file %s: %w", file, err)
		}
		allErrors = append(allErrors, errors...)
	}

	// Detect conflicts across all files
	conflictErrors, err := v.DetectConflicts()
	if err != nil {
		return nil, fmt.Errorf("failed to detect conflicts: %w", err)
	}
	allErrors = append(allErrors, conflictErrors...)

	return allErrors, nil
}

// ValidateFile validates a single annotation file
func (v *validator) ValidateFile(filePath string) ([]ValidationError, error) {
	var allErrors []ValidationError

	// Validate syntax
	syntaxErrors, err := v.ValidateSyntax(filePath)
	if err != nil {
		return nil, err
	}
	allErrors = append(allErrors, syntaxErrors...)

	// Validate required fields
	requiredErrors, err := v.ValidateRequiredFields(filePath)
	if err != nil {
		return nil, err
	}
	allErrors = append(allErrors, requiredErrors...)

	// Validate example schemas
	schemaErrors, err := v.ValidateExampleSchema(filePath)
	if err != nil {
		return nil, err
	}
	allErrors = append(allErrors, schemaErrors...)

	return allErrors, nil
}

// ValidateSyntax validates annotation syntax
func (v *validator) ValidateSyntax(filePath string) ([]ValidationError, error) {
	var errors []ValidationError

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Regex patterns for validation
	routerPattern := regexp.MustCompile(`^//\s*@Router\s+(/[^\s]*)\s+\[(\w+)\]`)
	paramPattern := regexp.MustCompile(`^//\s*@Param\s+(\w+)\s+(\w+)\s+(\S+)\s+(true|false)\s+"([^"]*)"`)
	successPattern := regexp.MustCompile(`^//\s*@Success\s+(\d{3})\s+\{(\w+)\}\s+(.+)`)
	failurePattern := regexp.MustCompile(`^//\s*@Failure\s+(\d{3})\s+\{(\w+)\}\s+(.+)`)

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Skip non-annotation lines
		if !strings.HasPrefix(trimmed, "// @") {
			continue
		}

		annotation := strings.TrimPrefix(trimmed, "// @")
		parts := strings.Fields(annotation)
		if len(parts) == 0 {
			continue
		}

		annotationType := parts[0]

		switch annotationType {
		case "Router":
			if !routerPattern.MatchString(trimmed) {
				errors = append(errors, ValidationError{
					File:    filePath,
					Line:    lineNum,
					Message: "invalid @Router syntax: expected format '@Router /path [method]'",
				})
			}

		case "Param":
			if !paramPattern.MatchString(trimmed) {
				errors = append(errors, ValidationError{
					File:    filePath,
					Line:    lineNum,
					Message: "invalid @Param syntax: expected format '@Param name location type required \"description\"'",
				})
			}

		case "Success":
			if !successPattern.MatchString(trimmed) {
				errors = append(errors, ValidationError{
					File:    filePath,
					Line:    lineNum,
					Message: "invalid @Success syntax: expected format '@Success statusCode {type} schema'",
				})
			}

		case "Failure":
			if !failurePattern.MatchString(trimmed) {
				errors = append(errors, ValidationError{
					File:    filePath,
					Line:    lineNum,
					Message: "invalid @Failure syntax: expected format '@Failure statusCode {type} schema'",
				})
			}

		case "Summary", "Description", "Tags", "Accept", "Produce", "Security":
			// These annotations should have content after the annotation type
			if len(parts) < 2 {
				errors = append(errors, ValidationError{
					File:    filePath,
					Line:    lineNum,
					Message: fmt.Sprintf("@%s annotation requires a value", annotationType),
				})
			}
		}
	}

	return errors, nil
}

// ValidateRequiredFields validates required fields in annotations
func (v *validator) ValidateRequiredFields(filePath string) ([]ValidationError, error) {
	var errors []ValidationError

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Track current endpoint annotations
	var currentEndpoint *endpointAnnotations
	var inAnnotationBlock bool
	var lastAnnotationLine int

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Check if this is an annotation line
		if strings.HasPrefix(trimmed, "// @") {
			lastAnnotationLine = lineNum
			annotation := strings.TrimPrefix(trimmed, "// @")
			parts := strings.Fields(annotation)
			if len(parts) == 0 {
				continue
			}

			annotationType := parts[0]

			// Router annotation marks the end of an endpoint definition
			if annotationType == "Router" {
				// Extract method and path for conflict detection
				if len(parts) >= 3 {
					path := parts[1]
					method := strings.Trim(parts[2], "[]")
					key := fmt.Sprintf("%s %s", method, path)
					v.endpoints[key] = append(v.endpoints[key], &endpointInfo{
						file:   filePath,
						line:   lineNum,
						method: method,
						path:   path,
					})
				}

				if currentEndpoint != nil {
					currentEndpoint.hasRouter = true
					// Validate and close this endpoint block
					endpointErrors := v.validateEndpointRequiredFields(currentEndpoint, filePath)
					errors = append(errors, endpointErrors...)
					currentEndpoint = nil
					inAnnotationBlock = false
				}
				continue
			}

			// Security annotation after Router is part of the previous endpoint, not a new one
			// Skip it if we don't have a current endpoint
			if annotationType == "Security" && currentEndpoint == nil {
				continue
			}

			// Start tracking annotations - start new endpoint if we don't have one
			if currentEndpoint == nil {
				inAnnotationBlock = true
				currentEndpoint = &endpointAnnotations{
					startLine: lineNum,
				}
			}

			switch annotationType {
			case "Summary":
				currentEndpoint.hasSummary = true
			case "Description":
				currentEndpoint.hasDescription = true
			case "Tags":
				currentEndpoint.hasTags = true
			case "Accept":
				currentEndpoint.hasAccept = true
			case "Produce":
				currentEndpoint.hasProduce = true
			case "Success":
				currentEndpoint.hasSuccess = true
			case "Router":
				currentEndpoint.hasRouter = true
			case "ID":
				// Track operation ID for conflict detection
				if len(parts) >= 2 {
					operationID := parts[1]
					v.operationIDs[operationID] = append(v.operationIDs[operationID], &endpointInfo{
						file:        filePath,
						line:        lineNum,
						operationID: operationID,
					})
				}
			}
		} else if inAnnotationBlock && trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			// End of annotation block - validate current endpoint
			inAnnotationBlock = false
			if currentEndpoint != nil {
				endpointErrors := v.validateEndpointRequiredFields(currentEndpoint, filePath)
				errors = append(errors, endpointErrors...)
				currentEndpoint = nil
			}
		} else if inAnnotationBlock && (trimmed == "" || strings.HasPrefix(trimmed, "//")) && lineNum > lastAnnotationLine+1 {
			// Empty line or regular comment after annotations - end of block
			if currentEndpoint != nil {
				endpointErrors := v.validateEndpointRequiredFields(currentEndpoint, filePath)
				errors = append(errors, endpointErrors...)
				currentEndpoint = nil
			}
			inAnnotationBlock = false
		}
	}

	// Validate last endpoint if exists (in case file doesn't end with proper closure)
	if currentEndpoint != nil {
		endpointErrors := v.validateEndpointRequiredFields(currentEndpoint, filePath)
		errors = append(errors, endpointErrors...)
	}

	return errors, nil
}

// validateEndpointRequiredFields validates that an endpoint has all required fields
func (v *validator) validateEndpointRequiredFields(endpoint *endpointAnnotations, filePath string) []ValidationError {
	var errors []ValidationError

	// Debug: print what we're validating
	// fmt.Printf("DEBUG: Validating endpoint at line %d: hasSummary=%v, hasTags=%v, hasSuccess=%v, hasRouter=%v\n",
	// 	endpoint.startLine, endpoint.hasSummary, endpoint.hasTags, endpoint.hasSuccess, endpoint.hasRouter)

	if !endpoint.hasSummary {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing required @Summary annotation",
		})
	}

	if !endpoint.hasTags {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing required @Tags annotation",
		})
	}

	if !endpoint.hasSuccess {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing required @Success annotation",
		})
	}

	if !endpoint.hasRouter {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing required @Router annotation",
		})
	}

	return errors
}

// DetectConflicts detects conflicts like duplicate paths and operation IDs
func (v *validator) DetectConflicts() ([]ValidationError, error) {
	var errors []ValidationError

	// Check for duplicate endpoints (same method + path)
	for key, infos := range v.endpoints {
		if len(infos) > 1 {
			// Found duplicate
			for _, info := range infos {
				errors = append(errors, ValidationError{
					File:    info.file,
					Line:    info.line,
					Message: fmt.Sprintf("duplicate endpoint: %s already defined in another location", key),
				})
			}
		}
	}

	// Check for duplicate operation IDs
	for id, infos := range v.operationIDs {
		if len(infos) > 1 {
			// Found duplicate operation ID
			for _, info := range infos {
				errors = append(errors, ValidationError{
					File:    info.file,
					Line:    info.line,
					Message: fmt.Sprintf("duplicate operation ID: %s already defined in another location", id),
				})
			}
		}
	}

	return errors, nil
}

// ValidateExampleSchema validates that examples conform to their schemas
func (v *validator) ValidateExampleSchema(filePath string) ([]ValidationError, error) {
	var errors []ValidationError

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Look for example annotations and validate JSON structure
	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Check for example in Success/Failure annotations
		if strings.Contains(trimmed, "@Success") || strings.Contains(trimmed, "@Failure") {
			// Extract schema reference if present
			// Format: @Success 200 {object} response.Response{data=auth.LoginResponse}
			// We need to find the {type} part
			startIdx := strings.Index(trimmed, "{")
			endIdx := strings.Index(trimmed, "}")
			if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
				schemaType := trimmed[startIdx+1 : endIdx]
				if schemaType != "object" && schemaType != "array" && schemaType != "string" &&
					schemaType != "integer" && schemaType != "number" && schemaType != "boolean" {
					errors = append(errors, ValidationError{
						File:    filePath,
						Line:    lineNum,
						Message: fmt.Sprintf("invalid schema type: %s (must be object, array, string, integer, number, or boolean)", schemaType),
					})
				}
			}
		}

		// Check for inline JSON examples (if any)
		if strings.Contains(trimmed, "example:") {
			// Extract JSON after "example:"
			exampleStart := strings.Index(trimmed, "example:")
			if exampleStart != -1 {
				jsonStr := strings.TrimSpace(trimmed[exampleStart+8:])
				if jsonStr != "" && (strings.HasPrefix(jsonStr, "{") || strings.HasPrefix(jsonStr, "[")) {
					// Try to parse as JSON
					var data interface{}
					if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
						errors = append(errors, ValidationError{
							File:    filePath,
							Line:    lineNum,
							Message: fmt.Sprintf("invalid JSON in example: %v", err),
						})
					}
				}
			}
		}
	}

	return errors, nil
}

// findAnnotationFiles finds all swagger annotation files
func (v *validator) findAnnotationFiles() ([]string, error) {
	var files []string

	// Search in internal/modules/*/docs/swagger.go
	modulesPath := filepath.Join(v.searchDir, "internal/modules")
	if _, err := os.Stat(modulesPath); !os.IsNotExist(err) {
		err := filepath.Walk(modulesPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Look for docs/swagger.go files
			if !info.IsDir() && strings.HasSuffix(path, "docs/swagger.go") {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	// Also check for main docs.go
	docsPath := filepath.Join(v.searchDir, "docs/docs.go")
	if _, err := os.Stat(docsPath); err == nil {
		files = append(files, docsPath)
	}

	return files, nil
}
