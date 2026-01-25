package export

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	FormatPostman  ExportFormat = "postman"
	FormatOpenAPI  ExportFormat = "openapi"
	FormatInsomnia ExportFormat = "insomnia"
	FormatSwagger  ExportFormat = "swagger"
	FormatApidog   ExportFormat = "apidog"
)

// ExportOptions contains configuration for export operations
type ExportOptions struct {
	Format            ExportFormat `json:"format"`
	EnvironmentID     *int64       `json:"environment_id,omitempty"`
	IncludeTests      bool         `json:"include_tests"`
	IncludeExamples   bool         `json:"include_examples"`
	IncludePreRequest bool         `json:"include_prerequest"`
	OutputFormat      string       `json:"output_format"` // json, yaml (for OpenAPI)
	SpecVersion       string       `json:"spec_version"`  // OpenAPI version
	IncludeServers    bool         `json:"include_servers"`
	IncludeSecurity   bool         `json:"include_security"`
}

// DefaultExportOptions returns default export options
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		IncludeTests:      true,
		IncludeExamples:   true,
		IncludePreRequest: true,
		OutputFormat:      "json",
		SpecVersion:       "3.0.0",
		IncludeServers:    true,
		IncludeSecurity:   true,
	}
}

// ExportResult contains the result of an export operation
type ExportResult struct {
	Content     interface{} `json:"content"`
	ContentType string      `json:"content_type"`
	Filename    string      `json:"filename"`
	Size        int64       `json:"size"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// ExportService interface defines the contract for export services
type ExportService interface {
	Export(collection *CollectionWithDetails, options *ExportOptions) (*ExportResult, error)
	GetSupportedFormats() []ExportFormat
	ValidateOptions(options *ExportOptions) error
}

// BaseExporter provides common functionality for all exporters
type BaseExporter struct {
	format ExportFormat
}

// NewBaseExporter creates a new base exporter
func NewBaseExporter(format ExportFormat) *BaseExporter {
	return &BaseExporter{
		format: format,
	}
}

// GetFormat returns the exporter format
func (e *BaseExporter) GetFormat() ExportFormat {
	return e.format
}

// ValidateOptions validates common export options
func (e *BaseExporter) ValidateOptions(options *ExportOptions) error {
	if options == nil {
		return fmt.Errorf("export options cannot be nil")
	}

	if options.Format != e.format {
		return fmt.Errorf("invalid format: expected %s, got %s", e.format, options.Format)
	}

	return nil
}

// GenerateFilename generates a filename for the export
func (e *BaseExporter) GenerateFilename(collectionName string, options *ExportOptions) string {
	// Sanitize collection name
	sanitized := strings.ReplaceAll(collectionName, " ", "_")
	sanitized = strings.ToLower(sanitized)

	// Remove special characters
	var result strings.Builder
	for _, r := range sanitized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}

	filename := result.String()
	if filename == "" {
		filename = "collection"
	}

	// Add appropriate extension based on format
	switch e.format {
	case FormatPostman:
		return fmt.Sprintf("%s.postman_collection.json", filename)
	case FormatOpenAPI:
		if options.OutputFormat == "yaml" {
			return fmt.Sprintf("%s.openapi.yaml", filename)
		}
		return fmt.Sprintf("%s.openapi.json", filename)
	case FormatInsomnia:
		return fmt.Sprintf("%s.insomnia_collection.json", filename)
	case FormatSwagger:
		return fmt.Sprintf("%s.swagger.json", filename)
	case FormatApidog:
		return fmt.Sprintf("%s.apidog_collection.json", filename)
	default:
		return fmt.Sprintf("%s.%s.json", filename, string(e.format))
	}
}

// SubstituteVariables replaces environment variables in text
func (e *BaseExporter) SubstituteVariables(text string, variables map[string]string) string {
	if text == "" || len(variables) == 0 {
		return text
	}

	result := text
	for key, value := range variables {
		// Replace {{variable}} format
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", key), value)
		// Replace ${variable} format
		result = strings.ReplaceAll(result, fmt.Sprintf("${%s}", key), value)
	}

	return result
}

// BuildVariableMap creates a map of environment variables
func (e *BaseExporter) BuildVariableMap(environment *EnvironmentWithVariables) map[string]string {
	if environment == nil || len(environment.Variables) == 0 {
		return make(map[string]string)
	}

	variables := make(map[string]string)
	for _, variable := range environment.Variables {
		if !variable.IsSecret || variable.Value != "" {
			variables[variable.KeyName] = variable.Value
		}
	}

	return variables
}

// SerializeJSON serializes content to JSON with proper formatting
func (e *BaseExporter) SerializeJSON(content interface{}) ([]byte, error) {
	return json.MarshalIndent(content, "", "  ")
}

// GetContentType returns the appropriate content type for the format
func (e *BaseExporter) GetContentType(options *ExportOptions) string {
	switch e.format {
	case FormatOpenAPI:
		if options.OutputFormat == "yaml" {
			return "application/x-yaml"
		}
		return "application/json"
	default:
		return "application/json"
	}
}

// ExportFactory creates appropriate exporter based on format
type ExportFactory struct{}

// NewExportFactory creates a new export factory
func NewExportFactory() *ExportFactory {
	return &ExportFactory{}
}

// CreateExporter creates an exporter for the specified format
func (f *ExportFactory) CreateExporter(format ExportFormat) (ExportService, error) {
	switch format {
	case FormatPostman:
		return NewPostmanExporter(), nil
	case FormatOpenAPI:
		return NewOpenAPIExporter(), nil
	case FormatInsomnia:
		return NewInsomniaExporter(), nil
	case FormatSwagger:
		return NewSwaggerExporter(), nil
	case FormatApidog:
		return NewApidogExporter(), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// GetSupportedFormats returns all supported export formats
func (f *ExportFactory) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{
		FormatPostman,
		FormatOpenAPI,
		FormatInsomnia,
		FormatSwagger,
		FormatApidog,
	}
}

// ValidateFormat checks if the format is supported
func (f *ExportFactory) ValidateFormat(format ExportFormat) error {
	supported := f.GetSupportedFormats()
	for _, supportedFormat := range supported {
		if format == supportedFormat {
			return nil
		}
	}
	return fmt.Errorf("unsupported format: %s", format)
}

// ExportManager manages export operations with caching and performance optimization
type ExportManager struct {
	factory   *ExportFactory
	cache     *ExportCache
	optimizer *PerformanceOptimizer
}

// NewExportManager creates a new export manager
func NewExportManager() *ExportManager {
	return &ExportManager{
		factory:   NewExportFactory(),
		cache:     nil, // Cache will be set separately
		optimizer: NewPerformanceOptimizer(DefaultPerformanceConfig()),
	}
}

// NewExportManagerWithCache creates a new export manager with caching
func NewExportManagerWithCache(cache *ExportCache) *ExportManager {
	return &ExportManager{
		factory:   NewExportFactory(),
		cache:     cache,
		optimizer: NewPerformanceOptimizer(DefaultPerformanceConfig()),
	}
}

// NewExportManagerWithOptimizations creates a new export manager with full optimizations
func NewExportManagerWithOptimizations(cache *ExportCache, perfConfig *PerformanceConfig) *ExportManager {
	return &ExportManager{
		factory:   NewExportFactory(),
		cache:     cache,
		optimizer: NewPerformanceOptimizer(perfConfig),
	}
}

// SetCache sets the cache instance
func (m *ExportManager) SetCache(cache *ExportCache) {
	m.cache = cache
}

// SetPerformanceConfig updates performance configuration
func (m *ExportManager) SetPerformanceConfig(config *PerformanceConfig) {
	if m.optimizer != nil {
		m.optimizer.UpdateConfig(config)
	}
}

// Export performs export operation with caching and performance optimization
func (m *ExportManager) Export(collection *CollectionWithDetails, format ExportFormat, options *ExportOptions) (*ExportResult, error) {
	return m.ExportWithContext(context.Background(), collection, format, options)
}

// ExportWithContext performs export operation with context support
func (m *ExportManager) ExportWithContext(ctx context.Context, collection *CollectionWithDetails, format ExportFormat, options *ExportOptions) (*ExportResult, error) {
	// Validate format
	if err := m.factory.ValidateFormat(format); err != nil {
		return nil, err
	}

	// Set format in options if not set
	if options == nil {
		options = DefaultExportOptions()
	}
	options.Format = format

	// Try to get from cache first
	if m.cache != nil && m.cache.IsEnabled() {
		cacheKey := m.cache.GenerateCacheKey(collection.ID, format, options)
		if cacheKey != "" {
			if cachedResult, err := m.cache.Get(cacheKey); err == nil && cachedResult != nil {
				// Cache hit - return cached result
				return cachedResult, nil
			}
		}
	}

	// Cache miss or caching disabled - perform export with optimization
	exporter, err := m.factory.CreateExporter(format)
	if err != nil {
		return nil, err
	}

	var result *ExportResult
	if m.optimizer != nil {
		// Use performance optimizer
		result, err = m.optimizer.OptimizedExport(ctx, collection, format, options, exporter)
	} else {
		// Fallback to direct export
		result, err = exporter.Export(collection, options)
	}

	if err != nil {
		return nil, err
	}

	// Store result in cache
	if m.cache != nil && m.cache.IsEnabled() && result != nil {
		cacheKey := m.cache.GenerateCacheKey(collection.ID, format, options)
		if cacheKey != "" {
			// Store in cache (ignore cache errors to not affect export functionality)
			_ = m.cache.Set(cacheKey, result)
		}
	}

	return result, nil
}

// InvalidateCache invalidates cache for a collection
func (m *ExportManager) InvalidateCache(collectionID int64) error {
	if m.cache != nil && m.cache.IsEnabled() {
		return m.cache.InvalidateCollection(collectionID)
	}
	return nil
}

// InvalidateCacheFormat invalidates cache for a collection and format
func (m *ExportManager) InvalidateCacheFormat(collectionID int64, format ExportFormat) error {
	if m.cache != nil && m.cache.IsEnabled() {
		return m.cache.InvalidateCollectionFormat(collectionID, format)
	}
	return nil
}

// GetCacheStats returns cache statistics
func (m *ExportManager) GetCacheStats() (*CacheStats, error) {
	if m.cache != nil {
		return m.cache.GetCacheStats()
	}
	return &CacheStats{Enabled: false}, nil
}

// GetPerformanceStats returns performance statistics
func (m *ExportManager) GetPerformanceStats() *PerformanceStats {
	if m.optimizer != nil {
		return m.optimizer.GetStats()
	}
	return &PerformanceStats{}
}

// GetExportProgress returns progress for a specific export operation
func (m *ExportManager) GetExportProgress(progressID string) (*ExportProgress, bool) {
	if m.optimizer != nil {
		return m.optimizer.GetProgress(progressID)
	}
	return nil, false
}

// GetAllExportProgress returns all active export progress
func (m *ExportManager) GetAllExportProgress() map[string]*ExportProgress {
	if m.optimizer != nil {
		return m.optimizer.GetAllProgress()
	}
	return make(map[string]*ExportProgress)
}

// GetPerformanceConfig returns current performance configuration
func (m *ExportManager) GetPerformanceConfig() *PerformanceConfig {
	if m.optimizer != nil {
		return m.optimizer.GetConfig()
	}
	return DefaultPerformanceConfig()
}

// GetSupportedFormats returns all supported formats
func (m *ExportManager) GetSupportedFormats() []ExportFormat {
	return m.factory.GetSupportedFormats()
}
