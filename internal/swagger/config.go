package swagger

import "time"

// Config holds the configuration for Swagger documentation
type Config struct {
	// General Info
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string

	// Security
	SecurityDefinitions map[string]SecurityScheme

	// Paths
	OutputDir   string
	SearchDir   string
	ExcludeDirs []string

	// Features
	EnableUI bool
	UIPath   string
	SpecPath string

	// Cache
	EnableCache bool
	CacheTTL    time.Duration
}

// SecurityScheme defines a security scheme for the API
type SecurityScheme struct {
	Type         string
	Description  string
	Name         string
	In           string
	Scheme       string
	BearerFormat string
	Flows        *OAuthFlows
}

// OAuthFlows defines OAuth 2.0 flows
type OAuthFlows struct {
	Implicit          *OAuthFlow
	Password          *OAuthFlow
	ClientCredentials *OAuthFlow
	AuthorizationCode *OAuthFlow
}

// OAuthFlow defines a single OAuth 2.0 flow
type OAuthFlow struct {
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

// DefaultConfig returns the default Swagger configuration
func DefaultConfig() *Config {
	return &Config{
		Title:       "Huminor Console API",
		Description: "Complete API for ERP with RBAC and API Documentation System",
		Version:     "1.0",
		Host:        "localhost:8081",
		BasePath:    "/",
		Schemes:     []string{"http", "https"},
		SecurityDefinitions: map[string]SecurityScheme{
			"BearerAuth": {
				Type:         "apiKey",
				Description:  "JWT Authorization header using the Bearer scheme",
				Name:         "Authorization",
				In:           "header",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
		OutputDir:   "docs",
		SearchDir:   "./",
		ExcludeDirs: []string{"vendor", "tmp", "bin"},
		EnableUI:    true,
		UIPath:      "/swagger",
		SpecPath:    "/api/swagger.json",
		EnableCache: true,
		CacheTTL:    5 * time.Minute,
	}
}
