# Swagger Documentation System

This directory contains the Swagger/OpenAPI documentation infrastructure for the Huminor Console API.

## Overview

The Swagger system provides automatic API documentation generation from code annotations. Documentation is generated using [swaggo/swag](https://github.com/swaggo/swag) and served via Swagger UI.

## Directory Structure

```
internal/swagger/
├── config.go      # Swagger configuration
├── handler.go     # HTTP handler for Swagger UI
└── README.md      # This file

pkg/swagger/
└── generator.go   # Wrapper for swag generation

docs/
├── docs.go        # Generated Go code (DO NOT EDIT)
├── swagger.json   # Generated OpenAPI spec (JSON)
└── swagger.yaml   # Generated OpenAPI spec (YAML)
```

## Quick Start

### 1. Generate Documentation

```bash
make swagger-gen
```

This will scan your code for Swagger annotations and generate:
- `docs/swagger.json` - OpenAPI specification in JSON format
- `docs/swagger.yaml` - OpenAPI specification in YAML format
- `docs/docs.go` - Go code for embedding the spec

### 2. Validate Annotations

```bash
make swagger-validate
```

### 3. Clean Generated Files

```bash
make swagger-clean
```

## Adding Swagger Annotations

### Main API Information

The main API information is defined in `cmd/api/main.go`:

```go
// @title           Huminor Console API
// @version         1.0
// @description     Complete API for ERP with RBAC and API Documentation System
// @host            localhost:8081
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Endpoint Annotations

For each module, create a `docs/swagger.go` file with endpoint annotations:

```go
// internal/modules/auth/docs/swagger.go
package docs

// @Summary      Login with user identity
// @Description  Authenticate user and return access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginRequest  true  "Login credentials"
// @Success      200          {object}  response.Response{data=auth.LoginResponse}
// @Failure      400          {object}  response.Response
// @Failure      401          {object}  response.Response
// @Router       /api/v1/auth/login [post]
```

### Common Annotations

- `@Summary` - Short description of the endpoint
- `@Description` - Detailed description
- `@Tags` - Group endpoints by tags
- `@Accept` - Request content type (json, xml, etc.)
- `@Produce` - Response content type
- `@Param` - Request parameters (query, path, body, header)
- `@Success` - Success response
- `@Failure` - Error responses
- `@Router` - Route path and HTTP method
- `@Security` - Security requirements (e.g., BearerAuth)

## Configuration

The default configuration is defined in `internal/swagger/config.go`:

```go
config := swagger.DefaultConfig()
// Customize as needed
config.Host = "api.example.com"
config.EnableUI = true
config.UIPath = "/swagger"
```

## Integration with Application

To integrate Swagger UI with your application:

```go
import (
    "gin-scalable-api/internal/swagger"
    _ "gin-scalable-api/docs" // Import generated docs
)

// Create Swagger handler
swaggerConfig := swagger.DefaultConfig()
swaggerHandler := swagger.NewHandler(swaggerConfig)

// Register routes
swaggerHandler.RegisterRoutes(router)
```

Then access Swagger UI at: `http://localhost:8081/swagger/index.html`

## Best Practices

1. **Separate Annotations from Route Files**: Keep annotations in `docs/swagger.go` files within each module
2. **Keep Annotations Up-to-Date**: Regenerate documentation after any API changes
3. **Use Descriptive Names**: Make summaries and descriptions clear and concise
4. **Document All Parameters**: Include all query, path, body, and header parameters
5. **Provide Examples**: Include request/response examples for better understanding
6. **Group by Tags**: Use tags to organize endpoints logically
7. **Document Error Cases**: Include all possible error responses

## Troubleshooting

### Documentation Not Updating

1. Clean and regenerate:
   ```bash
   make swagger-clean
   make swagger-gen
   ```

2. Clear browser cache or use incognito mode

### Annotations Not Found

- Ensure annotations are in files that are scanned by swag
- Check that the package is imported (even with `_` blank import)
- Verify annotation syntax is correct

### Build Errors

- Run `go mod tidy` to ensure all dependencies are installed
- Check that `github.com/swaggo/swag` is in go.mod

## References

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
