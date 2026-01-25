# API Documentation System - Authentication Integration

## Overview

The API Documentation System integrates with the existing token-based authentication system that uses Redis for token storage instead of JWT tokens.

## Authentication System Details

### Current Authentication Flow

1. **Login Process:**
   - User provides credentials (email/user_identity + password)
   - System generates random access token and refresh token
   - Tokens are stored in Redis with metadata
   - Access token expires in 15 minutes
   - Refresh token expires in 7 days

2. **Token Storage in Redis:**
   ```
   Key: access:user:{user_id}
   Value: {
     "token": "generated_token_string",
     "metadata": {
       "user_id": 123,
       "user_agent": "...",
       "ip": "...",
       "abilities": ["module1", "module2"],
       "expires_at": 1234567890
     }
   }
   ```

3. **Authentication Middleware:**
   - Extracts token from `Authorization: Bearer {token}` header
   - Looks up token in Redis
   - Validates expiration
   - Sets user context in Gin context

### User Context Available

From the authentication middleware, the following context is available:

- `user_id` (int64) - The authenticated user's ID
- `abilities` ([]string) - User's module permissions
- `company_id` (int64) - Available if using UnitAwareAuthMiddleware
- `unit_permissions` - Comprehensive unit permissions (if using UnitAwareAuthMiddleware)

### Company ID Resolution

Since the basic auth middleware doesn't provide `company_id`, the API Documentation Service includes a helper method to resolve it:

```go
func (s *Service) getUserCompanyID(userID int64) (int64, error) {
    var companyID int64
    query := `SELECT company_id FROM user_roles WHERE user_id = $1 LIMIT 1`
    err := s.db.QueryRow(query, userID).Scan(&companyID)
    return companyID, err
}
```

## Service Layer Integration

### Wrapper Methods

The service provides wrapper methods that automatically resolve company_id from user_id:

```go
// Instead of: CreateCollection(req, userID, companyID)
// Use: CreateCollectionByUser(req, userID)

func (s *Service) CreateCollectionByUser(req *CreateCollectionRequest, userID int64) (*CollectionResponse, error) {
    companyID, err := s.getUserCompanyID(userID)
    if err != nil {
        return nil, err
    }
    return s.CreateCollection(req, userID, companyID)
}
```

### Available Wrapper Methods

All service methods have corresponding `ByUser` wrapper methods:

**Collections:**
- `CreateCollectionByUser`
- `GetCollectionsByUser`
- `GetCollectionByIDByUser`
- `GetCollectionWithStatsByUser`
- `UpdateCollectionByUser`
- `DeleteCollectionByUser`

**Folders:**
- `CreateFolderByUser`
- `GetFoldersByUser`
- `GetFolderByIDByUser`
- `UpdateFolderByUser`
- `DeleteFolderByUser`
- `ReorderFoldersByUser`

**Endpoints:**
- `CreateEndpointByUser`
- `GetEndpointsByUser`
- `GetEndpointByIDByUser`
- `GetEndpointWithDetailsByUser`
- `UpdateEndpointByUser`
- `DeleteEndpointByUser`
- `BulkCreateEndpointsByUser`
- `BulkUpdateEndpointsByUser`
- `BulkDeleteEndpointsByUser`
- `BulkMoveEndpointsByUser`

**Environments:**
- `CreateEnvironmentByUser`
- `GetEnvironmentsByUser`
- `GetEnvironmentByIDByUser`
- `GetEnvironmentWithVariablesByUser`
- `UpdateEnvironmentByUser`
- `DeleteEnvironmentByUser`

**Environment Variables:**
- `CreateEnvironmentVariableByUser`
- `GetEnvironmentVariablesByUser`
- `UpdateEnvironmentVariableByUser`
- `DeleteEnvironmentVariableByUser`

## Route Handler Integration

Route handlers should extract `user_id` from context and use the `ByUser` wrapper methods:

```go
func (h *Handler) CreateCollection(c *gin.Context) {
    // Extract user ID from context
    userID, exists := c.Get("user_id")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
        return
    }
    
    userIDInt64, ok := userID.(int64)
    if !ok {
        response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid user ID type")
        return
    }
    
    // Parse request
    var req CreateCollectionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
        return
    }
    
    // Use wrapper method that automatically resolves company_id
    result, err := h.service.CreateCollectionByUser(&req, userIDInt64)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create collection", err.Error())
        return
    }
    
    response.Success(c, "Collection created successfully", result)
}
```

## Security Considerations

1. **Company Isolation:** All operations are automatically isolated by company through the user_roles table lookup
2. **Token Validation:** Tokens are validated against Redis storage and expiration times
3. **User Context:** User context is securely extracted from validated tokens
4. **Authorization:** Future RBAC integration will use the existing rbac.RBACService

## Migration Notes

- The original service methods (with explicit companyID parameter) are still available for internal use
- The new `ByUser` wrapper methods should be used in route handlers
- No changes needed to the authentication middleware
- Database connection is required for the service to resolve company_id from user_id

## Testing

Unit tests use `nil` for the database connection since they mock the repository layer:

```go
service := NewService(mockRepo, mockRBAC, nil) // nil DB for unit tests
```

Integration tests should provide a real database connection for company_id resolution.