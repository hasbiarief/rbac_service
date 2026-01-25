package export

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCompleteExportWorkflow(t *testing.T) {
	// Create a comprehensive test collection
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Complete API Collection",
			Description: "A comprehensive test collection",
			Version:     "2.1.0",
			BaseURL:     "https://api.example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Folders: []Folder{
			{
				ID:          1,
				Name:        "Authentication",
				Description: "Auth endpoints",
				SortOrder:   1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				Name:        "Users",
				Description: "User management",
				SortOrder:   2,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		Endpoints: []EndpointWithDetails{
			{
				Endpoint: Endpoint{
					ID:          1,
					Name:        "Login",
					Description: "User login endpoint",
					Method:      "POST",
					URL:         "{{base_url}}/auth/login",
					SortOrder:   1,
					IsActive:    true,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Headers: []Header{
					{
						ID:          1,
						KeyName:     "Content-Type",
						Value:       "application/json",
						Description: "Request content type",
						IsRequired:  true,
						HeaderType:  "request",
						CreatedAt:   time.Now(),
					},
				},
				Parameters: []Parameter{
					{
						ID:           1,
						Name:         "email",
						Type:         "body",
						DataType:     "string",
						Description:  "User email",
						IsRequired:   true,
						ExampleValue: "user@example.com",
						CreatedAt:    time.Now(),
					},
				},
				RequestBody: &RequestBody{
					ID:          1,
					ContentType: "application/json",
					BodyContent: `{"email": "{{user_email}}", "password": "{{user_password}}"}`,
					Description: "Login credentials",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Responses: []Response{
					{
						ID:           1,
						StatusCode:   200,
						StatusText:   "OK",
						ContentType:  "application/json",
						ResponseBody: `{"token": "jwt_token_here", "user": {"id": 1, "email": "user@example.com"}}`,
						Description:  "Successful login",
						IsDefault:    true,
						CreatedAt:    time.Now(),
					},
					{
						ID:           2,
						StatusCode:   401,
						StatusText:   "Unauthorized",
						ContentType:  "application/json",
						ResponseBody: `{"error": "Invalid credentials"}`,
						Description:  "Invalid login",
						IsDefault:    false,
						CreatedAt:    time.Now(),
					},
				},
			},
		},
		Environment: &EnvironmentWithVariables{
			Environment: Environment{
				ID:          1,
				Name:        "Production",
				Description: "Production environment",
				IsDefault:   true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			Variables: []EnvironmentVariable{
				{
					ID:          1,
					KeyName:     "base_url",
					Value:       "https://api.example.com",
					Description: "Base API URL",
					IsSecret:    false,
					CreatedAt:   time.Now(),
				},
				{
					ID:          2,
					KeyName:     "user_email",
					Value:       "test@example.com",
					Description: "Test user email",
					IsSecret:    false,
					CreatedAt:   time.Now(),
				},
				{
					ID:          3,
					KeyName:     "user_password",
					Value:       "secret123",
					Description: "Test user password",
					IsSecret:    true,
					CreatedAt:   time.Now(),
				},
			},
		},
	}

	// Test all export formats
	formats := []ExportFormat{
		FormatPostman,
		FormatOpenAPI,
		FormatInsomnia,
		FormatSwagger,
		FormatApidog,
	}

	manager := NewExportManager()

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			options := DefaultExportOptions()
			options.Format = format
			envID := int64(1)
			options.EnvironmentID = &envID

			result, err := manager.Export(collection, format, options)
			if err != nil {
				t.Errorf("Export failed for format %s: %v", format, err)
				return
			}

			if result == nil {
				t.Errorf("Expected export result for format %s, got nil", format)
				return
			}

			// Validate result structure
			if result.Content == "" {
				t.Errorf("Expected content for format %s, got empty", format)
			}

			if result.ContentType == "" {
				t.Errorf("Expected content type for format %s, got empty", format)
			}

			if result.Filename == "" {
				t.Errorf("Expected filename for format %s, got empty", format)
			}

			if result.Size <= 0 {
				t.Errorf("Expected positive size for format %s, got %d", format, result.Size)
			}

			// Try to parse JSON for JSON formats
			if result.ContentType == "application/json" {
				var parsed interface{}
				if err := json.Unmarshal([]byte(result.Content.(string)), &parsed); err != nil {
					t.Errorf("Invalid JSON for format %s: %v", format, err)
				}
			}

			t.Logf("Successfully exported %s format: %s (%d bytes)",
				format, result.Filename, result.Size)
		})
	}
}

func TestVariableSubstitution(t *testing.T) {
	collection := &CollectionWithDetails{
		Collection: Collection{
			ID:          1,
			Name:        "Variable Test",
			Description: "Test variable substitution",
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Endpoints: []EndpointWithDetails{
			{
				Endpoint: Endpoint{
					ID:        1,
					Name:      "Test Endpoint",
					Method:    "GET",
					URL:       "{{base_url}}/api/{{version}}/test",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				RequestBody: &RequestBody{
					ID:          1,
					ContentType: "application/json",
					BodyContent: `{"api_key": "{{api_key}}", "user_id": "{{user_id}}"}`,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		},
		Environment: &EnvironmentWithVariables{
			Environment: Environment{
				ID:        1,
				Name:      "Test Env",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Variables: []EnvironmentVariable{
				{
					ID:        1,
					KeyName:   "base_url",
					Value:     "https://test.api.com",
					CreatedAt: time.Now(),
				},
				{
					ID:        2,
					KeyName:   "version",
					Value:     "v2",
					CreatedAt: time.Now(),
				},
				{
					ID:        3,
					KeyName:   "api_key",
					Value:     "test_key_123",
					CreatedAt: time.Now(),
				},
				{
					ID:        4,
					KeyName:   "user_id",
					Value:     "12345",
					CreatedAt: time.Now(),
				},
			},
		},
	}

	exporter := NewPostmanExporter()
	options := DefaultExportOptions()
	options.Format = FormatPostman
	envID := int64(1)
	options.EnvironmentID = &envID

	result, err := exporter.Export(collection, options)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	content := result.Content.(string)

	// Check if variables were substituted
	expectedSubstitutions := []string{
		"https://test.api.com/api/v2/test",
		"test_key_123",
		"12345",
	}

	for _, expected := range expectedSubstitutions {
		if !contains(content, expected) {
			t.Errorf("Expected substitution '%s' not found in exported content", expected)
		}
	}

	// Check that variable placeholders were replaced
	unexpectedPlaceholders := []string{
		"{{base_url}}",
		"{{version}}",
		"{{api_key}}",
		"{{user_id}}",
	}

	for _, placeholder := range unexpectedPlaceholders {
		if contains(content, placeholder) {
			t.Errorf("Variable placeholder '%s' was not substituted", placeholder)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
