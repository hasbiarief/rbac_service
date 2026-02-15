# Swagger/OpenAPI Documentation System - Complete Guide

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Adding Swagger Annotations](#adding-swagger-annotations)
4. [Generating Documentation](#generating-documentation)
5. [Accessing Swagger UI](#accessing-swagger-ui)
6. [Common Patterns and Examples](#common-patterns-and-examples)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)
9. [Advanced Topics](#advanced-topics)

## Overview

The Huminor Console API uses an automated Swagger/OpenAPI documentation system that generates API documentation from code annotations. This approach ensures documentation stays synchronized with the actual API implementation.

### Key Features

- ✅ **Automatic Generation**: Documentation is generated from code annotations
- ✅ **Modular Structure**: Each module has its own annotation file
- ✅ **Interactive UI**: Test endpoints directly from Swagger UI
- ✅ **Multiple Formats**: Generates JSON, YAML, and embedded Go code
- ✅ **Validation**: Built-in validation ensures documentation quality
- ✅ **Watch Mode**: Auto-regeneration during development

### Architecture

```
Application Code
    ↓
Swagger Annotations (internal/modules/*/docs/swagger.go)
    ↓
Swag Generator
    ↓
OpenAPI Specification (docs/swagger.json, docs/swagger.yaml)
    ↓
Swagger UI (http://localhost:8081/swagger/index.html)
```

## Quick Start

### 1. Install Prerequisites

```bash
# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
swag --version
```

### 2. Generate Documentation

```bash
# Using Makefile (recommended)
make swagger-gen

# Or using the script
./scripts/generate-swagger.sh

# Or using the CLI tool
./bin/swagger -generate
```

### 3. Start the Application

```bash
make run
```

### 4. Access Swagger UI

Open your browser and navigate to:
```
http://localhost:8081/swagger/index.html
```

## Adding Swagger Annotations

### Module Structure

Each module should have a `docs/swagger.go` file containing its endpoint annotations:

```
internal/modules/
└── mymodule/
    ├── route.go          # Route definitions (no annotations here)
    ├── service.go        # Business logic
    ├── repository.go     # Data access
    └── docs/
        └── swagger.go    # Swagger annotations
```

### Basic Annotation Template

```go
// internal/modules/mymodule/docs/swagger.go
package docs

// Package docs berisi anotasi Swagger untuk modul mymodule.
// File ini terpisah dari route.go untuk menjaga kebersihan kode.

// @Summary      Short description of the endpoint
// @Description  Detailed description of what this endpoint does
// @Tags         ModuleName
// @Accept       json
// @Produce      json
// @Param        paramName  paramType  dataType  required  "Parameter description"
// @Success      200        {object}   response.Response{data=YourResponseType}  "Success message"
// @Failure      400        {object}   response.Response  "Bad request message"
// @Failure      401        {object}   response.Response  "Unauthorized message"
// @Security     BearerAuth
// @Router       /api/v1/mymodule/endpoint [post]
```

### Annotation Reference

#### Required Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description (1 line) | `@Summary Get user by ID` |
| `@Description` | Detailed description | `@Description Retrieves user information by user ID` |
| `@Tags` | Group endpoints by category | `@Tags Users` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Router` | Route path and HTTP method | `@Router /api/v1/users/{id} [get]` |

#### Optional Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Param` | Request parameter | `@Param id path int true "User ID"` |
| `@Success` | Success response | `@Success 200 {object} User` |
| `@Failure` | Error response | `@Failure 404 {object} response.Response` |
| `@Security` | Security requirement | `@Security BearerAuth` |
| `@Header` | Response header | `@Header 200 {string} Token "qwerty"` |

## Generating Documentation

### Using Makefile (Recommended)

```bash
# Generate documentation
make swagger-gen

# Validate annotations
make swagger-validate

# Watch for changes and auto-regenerate
make swagger-watch

# Clean generated files
make swagger-clean
```

### Using Scripts

```bash
# Generate with default settings
./scripts/generate-swagger.sh

# Validate only
./scripts/generate-swagger.sh -v

# Clean and regenerate
./scripts/generate-swagger.sh -c

# Verbose output
./scripts/generate-swagger.sh --verbose

# Custom output directory
./scripts/generate-swagger.sh -o ./api-docs
```

### Using CLI Tool

```bash
# Build the CLI tool first
make build-swagger

# Generate documentation
./bin/swagger -generate

# Validate annotations
./bin/swagger -validate

# Watch mode
./bin/swagger -watch

# Custom directories
./bin/swagger -generate -output ./api-docs -dir ./internal
```

### Generated Files

After generation, you'll find these files in the `docs/` directory:

- **swagger.json** - OpenAPI specification in JSON format
- **swagger.yaml** - OpenAPI specification in YAML format
- **docs.go** - Go code for embedding the specification

## Accessing Swagger UI

### Local Development

Start your application and access Swagger UI at:
```
http://localhost:8081/swagger/index.html
```

### Available Endpoints

- **Swagger UI**: `/swagger/index.html`
- **OpenAPI Spec (JSON)**: `/swagger/doc.json`
- **OpenAPI Spec (YAML)**: `/swagger/doc.yaml`

### Using Swagger UI

1. **Browse Endpoints**: Expand tags to see grouped endpoints
2. **View Details**: Click on an endpoint to see parameters and responses
3. **Try It Out**: Click "Try it out" to test the endpoint
4. **Authenticate**: Click "Authorize" to add your Bearer token
5. **Execute**: Fill in parameters and click "Execute"

## Common Patterns and Examples

### 1. Simple GET Endpoint

```go
// @Summary      Get user by ID
// @Description  Retrieves a user's information by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=user.User}  "User found"
// @Failure      404  {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [get]
```

### 2. POST Endpoint with Request Body

```go
// @Summary      Create new user
// @Description  Creates a new user with the provided information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      user.CreateUserRequest  true  "User information"
// @Success      201   {object}  response.Response{data=user.User}  "User created"
// @Failure      400   {object}  response.Response  "Invalid request"
// @Failure      409   {object}  response.Response  "User already exists"
// @Security     BearerAuth
// @Router       /api/v1/users [post]
```

### 3. PUT Endpoint for Updates

```go
// @Summary      Update user
// @Description  Updates an existing user's information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int                    true  "User ID"
// @Param        user  body      user.UpdateUserRequest true  "Updated user information"
// @Success      200   {object}  response.Response{data=user.User}  "User updated"
// @Failure      400   {object}  response.Response  "Invalid request"
// @Failure      404   {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [put]
```

### 4. DELETE Endpoint

```go
// @Summary      Delete user
// @Description  Deletes a user by their ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response  "User deleted"
// @Failure      404  {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id} [delete]
```

### 5. Query Parameters

```go
// @Summary      List users
// @Description  Retrieves a paginated list of users with optional filters
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"  default(1)
// @Param        limit     query     int     false  "Items per page"  default(10)
// @Param        search    query     string  false  "Search term"
// @Param        status    query     string  false  "User status"  Enums(active, inactive)
// @Success      200       {object}  response.Response{data=[]user.User}  "Users retrieved"
// @Failure      400       {object}  response.Response  "Invalid parameters"
// @Security     BearerAuth
// @Router       /api/v1/users [get]
```

### 6. Multiple Response Types

```go
// @Summary      Get user profile
// @Description  Retrieves user profile with different detail levels
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id      path      int     true   "User ID"
// @Param        detail  query     string  false  "Detail level"  Enums(basic, full)
// @Success      200     {object}  response.Response{data=user.BasicProfile}  "Basic profile"
// @Success      200     {object}  response.Response{data=user.FullProfile}   "Full profile"
// @Failure      404     {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/profile [get]
```

### 7. File Upload

```go
// @Summary      Upload user avatar
// @Description  Uploads a new avatar image for the user
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      int    true  "User ID"
// @Param        avatar  formData  file   true  "Avatar image file"
// @Success      200     {object}  response.Response{data=user.AvatarResponse}  "Avatar uploaded"
// @Failure      400     {object}  response.Response  "Invalid file"
// @Failure      413     {object}  response.Response  "File too large"
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/avatar [post]
```

### 8. Authentication Endpoint (No Security)

```go
// @Summary      Login with credentials
// @Description  Authenticates user and returns access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      auth.LoginRequest  true  "Login credentials"
// @Success      200          {object}  response.Response{data=auth.LoginResponse}  "Login successful"
// @Failure      400          {object}  response.Response  "Invalid request"
// @Failure      401          {object}  response.Response  "Invalid credentials"
// @Router       /api/v1/auth/login [post]
```

### 9. Array Response

```go
// @Summary      Get user roles
// @Description  Retrieves all roles assigned to a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=[]role.Role}  "Roles retrieved"
// @Failure      404  {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/roles [get]
```

### 10. Nested Object Response

```go
// @Summary      Get user with details
// @Description  Retrieves user with company and branch information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=user.UserWithDetails}  "User with details"
// @Failure      404  {object}  response.Response  "User not found"
// @Security     BearerAuth
// @Router       /api/v1/users/{id}/details [get]
```

## Best Practices

### 1. Annotation Organization

✅ **DO**: Keep annotations in separate `docs/swagger.go` files
```go
// internal/modules/user/docs/swagger.go
package docs

// All user endpoint annotations here
```

❌ **DON'T**: Mix annotations with route definitions
```go
// internal/modules/user/route.go - AVOID THIS
// @Summary Get user
func (r *UserRoute) GetUser(c *gin.Context) {
    // ...
}
```

### 2. Descriptive Documentation

✅ **DO**: Provide clear, detailed descriptions
```go
// @Summary      Create new user
// @Description  Creates a new user account with the provided information. Email must be unique. Password must be at least 8 characters.
```

❌ **DON'T**: Use vague descriptions
```go
// @Summary      User creation
// @Description  Creates user
```

### 3. Consistent Naming

✅ **DO**: Use consistent tag names across modules
```go
// @Tags         Users
// @Tags         Authentication
// @Tags         Companies
```

❌ **DON'T**: Use inconsistent or unclear tags
```go
// @Tags         user
// @Tags         Auth
// @Tags         company-management
```

### 4. Complete Parameter Documentation

✅ **DO**: Document all parameters with types and descriptions
```go
// @Param        id      path      int     true   "User ID"
// @Param        page    query     int     false  "Page number"  default(1)
// @Param        user    body      User    true   "User object"
```

❌ **DON'T**: Skip parameter documentation
```go
// @Param        id      path      int     true
// @Param        page    query     int     false
```

### 5. Error Response Documentation

✅ **DO**: Document all possible error responses
```go
// @Failure      400  {object}  response.Response  "Invalid request format"
// @Failure      401  {object}  response.Response  "Unauthorized - invalid token"
// @Failure      404  {object}  response.Response  "User not found"
// @Failure      500  {object}  response.Response  "Internal server error"
```

❌ **DON'T**: Only document success cases
```go
// @Success      200  {object}  User
```

### 6. Security Annotations

✅ **DO**: Add security annotations to protected endpoints
```go
// @Security     BearerAuth
// @Router       /api/v1/users [get]
```

❌ **DON'T**: Forget security annotations on protected routes

### 7. Type References

✅ **DO**: Use proper type references
```go
// @Success      200  {object}  response.Response{data=user.User}
// @Param        user  body     user.CreateUserRequest  true  "User data"
```

❌ **DON'T**: Use generic types
```go
// @Success      200  {object}  object
// @Param        user  body     object  true  "User data"
```

### 8. Keep Documentation Updated

✅ **DO**: Regenerate documentation after API changes
```bash
# After making changes
make swagger-gen
```

❌ **DON'T**: Let documentation become stale

### 9. Use Watch Mode During Development

✅ **DO**: Enable watch mode while developing
```bash
make swagger-watch
```

### 10. Validate Before Committing

✅ **DO**: Always validate annotations before committing
```bash
make swagger-validate
```

## Troubleshooting

### Problem: Documentation Not Updating

**Symptoms**: Changes to annotations don't appear in Swagger UI

**Solutions**:
1. Clean and regenerate:
   ```bash
   make swagger-clean
   make swagger-gen
   ```

2. Clear browser cache or use incognito mode

3. Restart the application:
   ```bash
   make run
   ```

### Problem: Annotations Not Found

**Symptoms**: Endpoints missing from generated documentation

**Solutions**:
1. Ensure the `docs` package is imported in `main.go`:
   ```go
   import _ "gin-scalable-api/docs"
   ```

2. Check that annotation files are in the correct location:
   ```
   internal/modules/{module}/docs/swagger.go
   ```

3. Verify annotation syntax is correct

4. Run with verbose mode to see errors:
   ```bash
   ./scripts/generate-swagger.sh --verbose
   ```

### Problem: Invalid Annotation Syntax

**Symptoms**: Generation fails with syntax errors

**Solutions**:
1. Check for common syntax errors:
   - Missing spaces after `@` symbol
   - Incorrect parameter format
   - Missing required fields

2. Validate annotations:
   ```bash
   make swagger-validate
   ```

3. Compare with working examples in other modules

### Problem: Type Not Found

**Symptoms**: Error about undefined types in annotations

**Solutions**:
1. Ensure the type is exported (starts with capital letter):
   ```go
   type User struct { ... }  // ✅ Correct
   type user struct { ... }  // ❌ Wrong
   ```

2. Use fully qualified type names:
   ```go
   // @Success 200 {object} user.User  // ✅ Correct
   // @Success 200 {object} User       // ❌ May fail
   ```

3. Ensure the package containing the type is in the search path

### Problem: Swagger UI Not Loading

**Symptoms**: 404 error when accessing `/swagger/index.html`

**Solutions**:
1. Verify Swagger handler is registered in `main.go`:
   ```go
   swaggerHandler := swagger.NewHandler(swagger.DefaultConfig())
   swaggerHandler.RegisterRoutes(router)
   ```

2. Check that Swagger UI is enabled in config:
   ```go
   config.EnableUI = true
   ```

3. Ensure the application is running on the correct port

### Problem: Authentication Not Working in Swagger UI

**Symptoms**: 401 errors when testing endpoints from Swagger UI

**Solutions**:
1. Click the "Authorize" button in Swagger UI

2. Enter your token in the format:
   ```
   Bearer your-token-here
   ```

3. Ensure the endpoint has the `@Security BearerAuth` annotation

### Problem: Slow Generation

**Symptoms**: Documentation generation takes too long

**Solutions**:
1. Exclude unnecessary directories:
   ```bash
   ./bin/swagger -generate -exclude "vendor,tmp,bin"
   ```

2. Use caching (enabled by default)

3. Only regenerate when needed (use watch mode)

## Advanced Topics

### Custom Configuration

Create a custom configuration in your application:

```go
config := &swagger.Config{
    Title:       "My Custom API",
    Description: "Custom API description",
    Version:     "2.0",
    Host:        "api.example.com",
    BasePath:    "/v2",
    Schemes:     []string{"https"},
    EnableUI:    true,
    UIPath:      "/docs",
    SpecPath:    "/docs/spec.json",
}

handler := swagger.NewHandler(config)
handler.RegisterRoutes(router)
```

### Multiple Environments

Define multiple servers in your main documentation:

```go
// @host      localhost:8081
// @BasePath  /

// @x-servers [{"url": "http://localhost:8081", "description": "Development"}, {"url": "https://staging-api.example.com", "description": "Staging"}, {"url": "https://api.example.com", "description": "Production"}]
```

### Custom Security Schemes

Add custom security schemes in `cmd/api/main.go`:

```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants admin access
```

### Migrating from Postman

If you have existing Postman collections, use the migration script:

```bash
# Migrate Postman collection to Swagger annotations
./scripts/migrate-postman.sh -f postman_collection.json

# Review generated annotations
ls internal/modules/*/docs/swagger.go

# Regenerate documentation
make swagger-gen
```

See [Migration Guide](../scripts/README-migrate-postman.md) for details.

### CI/CD Integration

Add documentation generation to your CI/CD pipeline:

```yaml
# .github/workflows/ci.yml
- name: Validate Swagger Annotations
  run: make swagger-validate

- name: Generate Documentation
  run: make swagger-gen

- name: Upload Documentation
  uses: actions/upload-artifact@v2
  with:
    name: swagger-docs
    path: docs/
```

### Versioning

Maintain multiple API versions:

```go
// v1 endpoints
// @Router /api/v1/users [get]

// v2 endpoints
// @Router /api/v2/users [get]
```

## Related Documentation

- [Swagger System README](../internal/swagger/README.md) - Technical details
- [Generation Scripts](../scripts/README-generate-swagger.md) - Script documentation
- [Migration Guide](../scripts/README-migrate-postman.md) - Postman migration
- [Swaggo Documentation](https://github.com/swaggo/swag) - Official swaggo docs
- [OpenAPI Specification](https://swagger.io/specification/) - OpenAPI standard

## Requirements Validation

This documentation satisfies the following requirements:

- **Requirement 1.2**: Documents annotation file structure and location
- **Requirement 2.1**: Documents generation process and commands
- **Requirement 5.1**: Documents Swagger UI access and usage
- **Requirement 6.1**: Documents make commands and build workflow
- **Requirement 8.4**: Provides templates and examples for new modules

---

**Need Help?** Check the troubleshooting section or refer to the related documentation above.
