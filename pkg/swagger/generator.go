package swagger

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Generator generates Swagger documentation from annotations
type Generator interface {
	// Generate generates swagger.json and swagger.yaml files
	Generate(ctx context.Context, opts GenerateOptions) error

	// Validate validates annotations before generation
	Validate(ctx context.Context) ([]ValidationError, error)

	// Watch monitors changes and regenerates automatically
	Watch(ctx context.Context, opts WatchOptions) error
}

// GenerateOptions contains options for documentation generation
type GenerateOptions struct {
	OutputDir       string   // Output directory (default: docs/)
	SearchDir       string   // Search directory (default: ./)
	Exclude         []string // Directories to exclude
	ParseVendor     bool     // Parse vendor directory
	ParseDependency bool     // Parse dependency
	MarkdownFiles   []string // Markdown files to include
	GeneralInfo     GeneralInfo
}

// GeneralInfo contains general API information
type GeneralInfo struct {
	Title        string
	Description  string
	Version      string
	Host         string
	BasePath     string
	Schemes      []string
	ContactName  string
	ContactEmail string
	ContactURL   string
	LicenseName  string
	LicenseURL   string
}

// WatchOptions contains options for watch mode
type WatchOptions struct {
	GenerateOptions
	Interval      time.Duration // Check interval
	DebounceDelay time.Duration // Delay before regeneration
}

// ValidationError represents a validation error
type ValidationError struct {
	File    string
	Line    int
	Message string
}

// swagGenerator implements Generator using swaggo/swag
type swagGenerator struct {
	config *Config
}

// Config holds generator configuration
type Config struct {
	SwagPath string // Path to swag binary
}

// NewGenerator creates a new Swagger generator
func NewGenerator(config *Config) Generator {
	if config == nil {
		config = &Config{
			SwagPath: "swag",
		}
	}
	return &swagGenerator{config: config}
}

// Generate generates Swagger documentation
func (g *swagGenerator) Generate(ctx context.Context, opts GenerateOptions) error {
	args := []string{"init"}

	// Set general info file (main.go location)
	args = append(args, "--generalInfo", "cmd/api/main.go")

	// Set output directory
	outputDir := opts.OutputDir
	if outputDir == "" {
		outputDir = "docs"
	}
	args = append(args, "--output", outputDir)

	// Set search directory
	searchDir := opts.SearchDir
	if searchDir == "" {
		searchDir = "./"
	}
	args = append(args, "--dir", searchDir)

	// Exclude directories
	if len(opts.Exclude) > 0 {
		args = append(args, "--exclude", strings.Join(opts.Exclude, ","))
	}

	// Parse vendor
	if opts.ParseVendor {
		args = append(args, "--parseVendor")
	}

	// Parse dependency
	if opts.ParseDependency {
		args = append(args, "--parseDependency")
	}

	// Parse internal packages
	args = append(args, "--parseInternal")

	// General info overrides (if provided)
	if opts.GeneralInfo.Title != "" {
		args = append(args, "--generalInfo.title", opts.GeneralInfo.Title)
	}

	// Execute swag command
	cmd := exec.CommandContext(ctx, g.config.SwagPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Parse swag output for better error messages
		outputStr := string(output)
		if strings.Contains(outputStr, "cannot find type definition") {
			return fmt.Errorf("swag generation failed: missing type definitions\nOutput: %s", outputStr)
		}
		if strings.Contains(outputStr, "ParseComment error") {
			return fmt.Errorf("swag generation failed: annotation parsing error\nOutput: %s", outputStr)
		}
		return fmt.Errorf("swag generation failed: %w\nOutput: %s", err, outputStr)
	}

	return nil
}

// Validate validates Swagger annotations
func (g *swagGenerator) Validate(ctx context.Context) ([]ValidationError, error) {
	var errors []ValidationError

	// Find all swagger annotation files
	annotationFiles, err := g.findAnnotationFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find annotation files: %w", err)
	}

	// Validate each file
	for _, file := range annotationFiles {
		fileErrors := g.validateFile(file)
		errors = append(errors, fileErrors...)
	}

	return errors, nil
}

// findAnnotationFiles finds all swagger annotation files
func (g *swagGenerator) findAnnotationFiles() ([]string, error) {
	var files []string

	// Search in internal/modules/*/docs/swagger.go
	modulesPath := "internal/modules"
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		return files, nil
	}

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

	// Also check for main docs.go
	if _, err := os.Stat("docs/docs.go"); err == nil {
		files = append(files, "docs/docs.go")
	}

	return files, nil
}

// validateFile validates a single annotation file
func (g *swagGenerator) validateFile(filePath string) []ValidationError {
	var errors []ValidationError

	content, err := os.ReadFile(filePath)
	if err != nil {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    0,
			Message: fmt.Sprintf("failed to read file: %v", err),
		})
		return errors
	}

	lines := strings.Split(string(content), "\n")

	// Track required annotations for endpoints
	var currentEndpoint *endpointAnnotations
	var inAnnotationBlock bool

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Check if this is an annotation line
		if strings.HasPrefix(trimmed, "// @") {
			inAnnotationBlock = true
			annotation := strings.TrimPrefix(trimmed, "// @")
			parts := strings.Fields(annotation)
			if len(parts) == 0 {
				continue
			}

			annotationType := parts[0]

			// Start tracking endpoint if we don't have one
			if currentEndpoint == nil {
				currentEndpoint = &endpointAnnotations{
					startLine: lineNum,
				}
			}

			// Track annotations
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
				// Validate router format: /path [method]
				if len(parts) < 3 {
					errors = append(errors, ValidationError{
						File:    filePath,
						Line:    lineNum,
						Message: "Router annotation must have format: @Router /path [method]",
					})
				}
				// Router marks the end of endpoint definition - validate and close
				endpointErrors := g.validateEndpoint(currentEndpoint, filePath)
				errors = append(errors, endpointErrors...)
				currentEndpoint = nil
				inAnnotationBlock = false
			case "Param":
				// Validate param format
				if len(parts) < 5 {
					errors = append(errors, ValidationError{
						File:    filePath,
						Line:    lineNum,
						Message: "Param annotation must have format: @Param name location type required \"description\"",
					})
				}
			case "Security":
				// Security can appear after Router, skip it
				if currentEndpoint == nil {
					continue
				}
			}
		} else if inAnnotationBlock && trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			// End of annotation block
			inAnnotationBlock = false
			if currentEndpoint != nil {
				endpointErrors := g.validateEndpoint(currentEndpoint, filePath)
				errors = append(errors, endpointErrors...)
				currentEndpoint = nil
			}
		}
	}

	// No need to validate last endpoint - we validate when we see @Router

	return errors
}

// endpointAnnotations tracks required annotations for an endpoint
type endpointAnnotations struct {
	startLine      int
	hasSummary     bool
	hasDescription bool
	hasTags        bool
	hasAccept      bool
	hasProduce     bool
	hasSuccess     bool
	hasRouter      bool
}

// validateEndpoint validates that an endpoint has all required annotations
func (g *swagGenerator) validateEndpoint(endpoint *endpointAnnotations, filePath string) []ValidationError {
	var errors []ValidationError

	if !endpoint.hasSummary {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing @Summary annotation",
		})
	}

	if !endpoint.hasTags {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing @Tags annotation",
		})
	}

	if !endpoint.hasSuccess {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing @Success annotation",
		})
	}

	if !endpoint.hasRouter {
		errors = append(errors, ValidationError{
			File:    filePath,
			Line:    endpoint.startLine,
			Message: "endpoint missing @Router annotation",
		})
	}

	return errors
}

// Watch monitors changes and regenerates documentation
func (g *swagGenerator) Watch(ctx context.Context, opts WatchOptions) error {
	// Set default values
	if opts.Interval == 0 {
		opts.Interval = 2 * time.Second
	}
	if opts.DebounceDelay == 0 {
		opts.DebounceDelay = 500 * time.Millisecond
	}

	fmt.Println("Starting Swagger watch mode...")
	fmt.Printf("Watching for changes every %v\n", opts.Interval)

	// Track file modification times
	fileModTimes := make(map[string]time.Time)

	// Initial scan
	files, err := g.findAnnotationFiles()
	if err != nil {
		return fmt.Errorf("failed to find annotation files: %w", err)
	}

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fileModTimes[file] = info.ModTime()
	}

	// Generate initial documentation
	fmt.Println("Generating initial documentation...")
	if err := g.Generate(ctx, opts.GenerateOptions); err != nil {
		fmt.Printf("Initial generation failed: %v\n", err)
	} else {
		fmt.Println("Initial documentation generated successfully")
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	var pendingRegeneration bool
	var lastChangeTime time.Time

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Watch mode stopped")
			return ctx.Err()

		case <-ticker.C:
			// Check for file changes
			files, err := g.findAnnotationFiles()
			if err != nil {
				fmt.Printf("Error finding files: %v\n", err)
				continue
			}

			for _, file := range files {
				info, err := os.Stat(file)
				if err != nil {
					continue
				}

				lastMod, exists := fileModTimes[file]
				if !exists || info.ModTime().After(lastMod) {
					fileModTimes[file] = info.ModTime()
					if exists {
						fmt.Printf("Detected change in: %s\n", file)
						lastChangeTime = time.Now()
						pendingRegeneration = true
					}
				}
			}

			// Check for new files
			if len(files) != len(fileModTimes) {
				lastChangeTime = time.Now()
				pendingRegeneration = true
				fmt.Println("Detected new annotation files")
			}

			// Regenerate if changes detected and debounce period passed
			if pendingRegeneration && time.Since(lastChangeTime) >= opts.DebounceDelay {
				fmt.Println("Regenerating documentation...")
				if err := g.Generate(ctx, opts.GenerateOptions); err != nil {
					fmt.Printf("Generation failed: %v\n", err)
				} else {
					fmt.Println("Documentation regenerated successfully")
				}
				pendingRegeneration = false
			}
		}
	}
}
