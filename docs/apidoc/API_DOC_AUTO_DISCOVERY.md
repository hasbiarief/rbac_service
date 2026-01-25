# API Documentation Auto-Discovery

## Overview

Sistem auto-discovery akan menganalisis struktur routes yang ada di project untuk secara otomatis mengisi dokumentasi API. Ini akan mengurangi effort manual dalam mendokumentasikan setiap endpoint.

## Current Route Structure Analysis

Berdasarkan analisis `internal/app/routes.go`, project ini menggunakan struktur modular dengan pattern:

```go
// Auth module routes (public)
authModule.RegisterRoutes(api, h.Auth)

// Protected routes
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
{
    userModule.RegisterRoutes(protected, h.User)
    roleModule.RegisterRoutes(protected, h.Role)
    companyModule.RegisterRoutes(protected, h.Company)
    // ... other modules
}
```

## Auto-Discovery Strategy

### 1. Route Registration Scanning

#### Gin Route Extraction
```go
type RouteInfo struct {
    Method      string            `json:"method"`
    Path        string            `json:"path"`
    HandlerName string            `json:"handler_name"`
    Module      string            `json:"module"`
    IsProtected bool              `json:"is_protected"`
    Middleware  []string          `json:"middleware"`
    Parameters  []ParameterInfo   `json:"parameters"`
}

type ParameterInfo struct {
    Name        string `json:"name"`
    Type        string `json:"type"` // path, query, body
    Required    bool   `json:"required"`
    DataType    string `json:"data_type"`
    Description string `json:"description"`
}
```

#### Route Scanner Implementation
```go
func ScanGinRoutes(engine *gin.Engine) ([]RouteInfo, error) {
    var routes []RouteInfo
    
    for _, route := range engine.Routes() {
        routeInfo := RouteInfo{
            Method:      route.Method,
            Path:        route.Path,
            HandlerName: extractHandlerName(route.Handler),
            Module:      extractModuleFromPath(route.Path),
            IsProtected: isProtectedRoute(route.Path),
            Middleware:  extractMiddleware(route),
            Parameters:  extractParameters(route.Path),
        }
        routes = append(routes, routeInfo)
    }
    
    return routes, nil
}
```

### 2. Module-Based Organization

#### Automatic Folder Creation
Berdasarkan struktur module yang ada:
- Authentication → `/api/v1/auth/*`
- User Management → `/api/v1/users/*`
- Role Management → `/api/v1/role-management/*`
- Company Management → `/api/v1/companies/*`
- Branch Management → `/api/v1/branches/*`
- Unit Management → `/api/v1/units/*`
- Module System → `/api/v1/modules/*`
- Audit Logging → `/api/v1/audit/*`
- Subscription → `/api/v1/subscriptions/*`

#### Module Mapping
```go
var moduleMapping = map[string]string{
    "/api/v1/auth":             "Authentication",
    "/api/v1/users":            "User Management",
    "/api/v1/role-management":  "Role Management System",
    "/api/v1/companies":        "Company Management",
    "/api/v1/branches":         "Branch Management",
    "/api/v1/units":            "Unit Management",
    "/api/v1/modules":          "Module System",
    "/api/v1/audit":            "Audit Logging",
    "/api/v1/subscriptions":    "Subscription Management",
}
```

### 3. Parameter Extraction

#### Path Parameters
```go
// From route: /api/v1/users/:id
// Extract: {name: "id", type: "path", required: true, data_type: "integer"}

func extractPathParameters(path string) []ParameterInfo {
    var params []ParameterInfo
    
    // Find :param patterns
    re := regexp.MustCompile(`:(\w+)`)
    matches := re.FindAllStringSubmatch(path, -1)
    
    for _, match := range matches {
        param := ParameterInfo{
            Name:     match[1],
            Type:     "path",
            Required: true,
            DataType: inferDataType(match[1]), // id -> integer, slug -> string
        }
        params = append(params, param)
    }
    
    return params
}
```

#### Query Parameters (dari existing Postman collection)
```go
// Analyze existing collection untuk common query parameters
var commonQueryParams = map[string]ParameterInfo{
    "limit": {
        Name:        "limit",
        Type:        "query",
        Required:    false,
        DataType:    "integer",
        Description: "Number of items to return",
        DefaultValue: "20",
    },
    "offset": {
        Name:        "offset", 
        Type:        "query",
        Required:    false,
        DataType:    "integer",
        Description: "Number of items to skip",
        DefaultValue: "0",
    },
    "search": {
        Name:        "search",
        Type:        "query", 
        Required:    false,
        DataType:    "string",
        Description: "Search term",
    },
    // ... other common params
}
```

### 4. Authentication Detection

#### Middleware Analysis
```go
func detectAuthentication(route gin.RouteInfo) []HeaderInfo {
    var headers []HeaderInfo
    
    // Check if route uses AuthMiddleware
    if isProtectedRoute(route.Path) {
        headers = append(headers, HeaderInfo{
            KeyName:     "Authorization",
            Value:       "Bearer {{accessToken}}",
            Description: "JWT Access Token",
            IsRequired:  true,
            HeaderType:  "request",
        })
    }
    
    return headers
}
```

### 5. Response Schema Inference

#### From Existing Postman Collection
```go
type ResponseTemplate struct {
    StatusCode   int                    `json:"status_code"`
    ContentType  string                 `json:"content_type"`
    Schema       map[string]interface{} `json:"schema"`
    Example      map[string]interface{} `json:"example"`
    Description  string                 `json:"description"`
}

// Common response patterns
var responseTemplates = map[string][]ResponseTemplate{
    "GET_LIST": {
        {
            StatusCode:  200,
            ContentType: "application/json",
            Schema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "success": map[string]interface{}{"type": "boolean"},
                    "data": map[string]interface{}{
                        "type": "array",
                        "items": map[string]interface{}{"type": "object"},
                    },
                    "pagination": map[string]interface{}{"type": "object"},
                },
            },
            Example: map[string]interface{}{
                "success": true,
                "data": []interface{}{},
                "pagination": map[string]interface{}{
                    "total": 0,
                    "limit": 20,
                    "offset": 0,
                },
            },
        },
    },
    "GET_SINGLE": {
        {
            StatusCode:  200,
            ContentType: "application/json",
            Example: map[string]interface{}{
                "success": true,
                "data": map[string]interface{}{},
            },
        },
    },
    // ... other patterns
}
```

## Implementation Plan

### Phase 1: Basic Route Scanning
```go
// Service method
func (s *ApiDocService) ScanAndPopulateRoutes(collectionID int64, ginEngine *gin.Engine) error {
    // 1. Scan all routes from Gin engine
    routes, err := ScanGinRoutes(ginEngine)
    if err != nil {
        return err
    }
    
    // 2. Create folders based on modules
    folderMap, err := s.createModuleFolders(collectionID, routes)
    if err != nil {
        return err
    }
    
    // 3. Create endpoints
    for _, route := range routes {
        endpoint := &models.ApiEndpoint{
            CollectionID: collectionID,
            FolderID:     folderMap[route.Module],
            Name:         generateEndpointName(route),
            Method:       route.Method,
            URL:          route.Path,
            Description:  generateDescription(route),
        }
        
        // Save endpoint
        endpointID, err := s.repo.CreateEndpoint(endpoint)
        if err != nil {
            continue
        }
        
        // Add parameters
        s.addParameters(endpointID, route.Parameters)
        
        // Add headers
        s.addHeaders(endpointID, route)
        
        // Add response templates
        s.addResponseTemplates(endpointID, route)
    }
    
    return nil
}
```

### Phase 2: Enhanced Analysis
```go
// Analyze handler functions untuk lebih detail
func analyzeHandlerFunction(handlerName string) (*HandlerInfo, error) {
    // Use reflection atau static analysis
    // Extract parameter types, return types, comments
    return &HandlerInfo{
        Parameters:   extractHandlerParams(handlerName),
        ReturnType:   extractReturnType(handlerName),
        Comments:     extractComments(handlerName),
        RequestBody:  inferRequestBody(handlerName),
    }, nil
}
```

## API Endpoints untuk Auto-Discovery

### 1. Scan Routes
```
POST /api/v1/api-docs/collections/{id}/scan-routes
```

Request Body:
```json
{
    "overwrite_existing": false,
    "modules_to_scan": ["auth", "users", "companies"],
    "include_middleware_info": true
}
```

Response:
```json
{
    "success": true,
    "data": {
        "scanned_routes": 45,
        "created_endpoints": 42,
        "created_folders": 8,
        "skipped_routes": 3,
        "summary": {
            "Authentication": 6,
            "User Management": 12,
            "Company Management": 8
        }
    }
}
```

### 2. Preview Scan Results
```
GET /api/v1/api-docs/routes/scan-preview
```

Response:
```json
{
    "success": true,
    "data": {
        "total_routes": 45,
        "routes_by_module": {
            "Authentication": [
                {
                    "method": "POST",
                    "path": "/api/v1/auth/login",
                    "handler": "Login",
                    "estimated_name": "Login with credentials"
                }
            ]
        }
    }
}
```

### 3. Get Available Routes
```
GET /api/v1/api-docs/routes/available
```

Response:
```json
{
    "success": true,
    "data": [
        {
            "method": "GET",
            "path": "/health",
            "module": "System",
            "is_protected": false
        },
        {
            "method": "POST", 
            "path": "/api/v1/auth/login",
            "module": "Authentication",
            "is_protected": false
        }
    ]
}
```

## Benefits

1. **Rapid Documentation**: Instantly populate basic documentation
2. **Consistency**: Standardized naming dan structure
3. **Maintenance**: Easy to keep documentation in sync
4. **Coverage**: Ensure all endpoints are documented
5. **Efficiency**: Reduce manual documentation effort

## Limitations

1. **Basic Information Only**: Auto-discovery provides skeleton, manual enhancement needed
2. **No Business Logic**: Cannot infer complex business rules
3. **Limited Schema Detection**: Basic type inference only
4. **No Examples**: Real examples need manual input
5. **Comment Dependency**: Better results with good code comments

## Manual Enhancement Areas

After auto-discovery, manual enhancement needed for:
- Detailed descriptions
- Request/response examples
- Business logic documentation
- Error scenarios
- Authentication specifics
- Validation rules