package export

import (
	"testing"
	"time"
)

func TestExportManager(t *testing.T) {
	manager := NewExportManager()

	// Test supported formats
	formats := manager.GetSupportedFormats()
	if len(formats) == 0 {
		t.Error("Expected supported formats, got none")
	}

	expectedFormats := []ExportFormat{
		FormatPostman,
		FormatOpenAPI,
		FormatInsomnia,
		FormatSwagger,
		FormatApidog,
	}

	if len(formats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(formats))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected format %s not found in supported formats", expected)
		}
	}
}

func TestPostmanExporter(t *testing.T) {
	exporter := NewPostmanExporter()

	// Test sample collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: []EndpointWithDetails{},
	}

	options := DefaultExportOptions()
	options.Format = FormatPostman

	result, err := exporter.Export(collection, options)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	if result == nil {
		t.Error("Expected export result, got nil")
	}

	if result.ContentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", result.ContentType)
	}

	if result.Filename == "" {
		t.Error("Expected filename, got empty string")
	}
}

func TestApidogExporter(t *testing.T) {
	exporter := NewApidogExporter()

	// Test sample collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Test Collection",
			Description: "Test Description",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders:   []Folder{},
		Endpoints: []EndpointWithDetails{},
	}

	options := DefaultExportOptions()
	options.Format = FormatApidog

	result, err := exporter.Export(collection, options)
	if err != nil {
		t.Errorf("Export failed: %v", err)
	}

	if result == nil {
		t.Error("Expected export result, got nil")
	}

	if result.ContentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", result.ContentType)
	}
}

func TestBaseExporter(t *testing.T) {
	exporter := NewBaseExporter(FormatPostman)

	// Test filename generation
	filename := exporter.GenerateFilename("Test Collection", &ExportOptions{Format: FormatPostman})
	expected := "test_collection.postman_collection.json"
	if filename != expected {
		t.Errorf("Expected filename %s, got %s", expected, filename)
	}

	// Test variable substitution
	variables := map[string]string{
		"base_url": "https://api.example.com",
		"version":  "v1",
	}

	text := "{{base_url}}/{{version}}/users"
	result := exporter.SubstituteVariables(text, variables)
	expected = "https://api.example.com/v1/users"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
