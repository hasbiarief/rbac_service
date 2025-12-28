# Validation System Documentation

## Gambaran Umum

Project ini menggunakan sistem validasi yang modular dan terpusat untuk memastikan konsistensi dan kemudahan maintenance. Sistem validasi telah dipisahkan dari routes untuk meningkatkan readability dan reusability.

## Arsitektur Validation System

### 1. Struktur Modular

```
internal/validation/
├── common_validation.go      # Common rules (ID, pagination, etc.)
├── auth_validation.go        # Authentication validation
├── user_validation.go        # User management validation
├── company_validation.go     # Company validation
├── role_validation.go        # Role & RBAC validation
├── module_validation.go      # Module validation
├── subscription_validation.go # Subscription validation
├── audit_validation.go       # Audit logging validation
└── branch_validation.go      # Branch validation
```

### 2. Middleware Integration

**Lokasi:** `middleware/validation.go`

**Flow:**
1. Request masuk ke endpoint
2. Validation middleware memproses request
3. Validated data disimpan di context
4. Handler mengambil validated data dari context

## Validation Rules

### Common Validations (`common_validation.go`)

**ID Validation:**
```go
var IDValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
    },
}
```

**User ID Validation:**
```go
var UserIDValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "userId", Type: "int", Required: true, Min: IntPtr(1)},
    },
}
```

**Pagination Validation:**
```go
var ListValidation = middleware.ValidationRules{
    Query: []middleware.QueryValidation{
        {Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
        {Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
        {Name: "search", Type: "string"},
    },
}
```

**Identity Validation:**
```go
var IdentityValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "identity", Type: "string", Required: true, Min: IntPtr(3), Max: IntPtr(50)},
    },
}
```

### Authentication Validations (`auth_validation.go`)

**Login dengan User Identity:**
```go
var LoginValidation = middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
    }{},
}
```

**Login dengan Email:**
```go
var LoginEmailValidation = middleware.ValidationRules{
    Body: &struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=6"`
    }{},
}
```

**Refresh Token:**
```go
var RefreshValidation = middleware.ValidationRules{
    Body: &struct {
        RefreshToken string `json:"refresh_token" validate:"required"`
    }{},
}
```

**Logout:**
```go
var LogoutValidation = middleware.ValidationRules{
    Body: &struct {
        Token string `json:"token" validate:"required"`
    }{},
}
```

### User Management Validations (`user_validation.go`)

**Create User:**
```go
var CreateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name         string  `json:"name" validate:"required,min=2,max=100"`
        Email        string  `json:"email" validate:"required,email,max=255"`
        UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
        Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
    }{},
}
```

**Update User:**
```go
var UpdateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        IsActive *bool  `json:"is_active"`
    }{},
}
```

**Password Change:**
```go
var PasswordChangeValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
    },
    Body: &struct {
        CurrentPassword string `json:"current_password" validate:"required,min=6"`
        NewPassword     string `json:"new_password" validate:"required,min=6,max=100"`
        ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
    }{},
}
```

**Access Check:**
```go
var AccessCheckValidation = middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        ModuleURL    string `json:"module_url" validate:"required"`
    }{},
}
```

### Role Management Validations (`role_validation.go`)

**Create Role:**
```go
var CreateRoleValidation = middleware.ValidationRules{
    Body: &struct {
        Name        string `json:"name" validate:"required,min=2,max=100"`
        Description string `json:"description"`
    }{},
}
```

**Assign User Role:**
```go
var AssignUserRoleValidation = middleware.ValidationRules{
    Body: &struct {
        UserID    int64  `json:"user_id" validate:"required,min=1"`
        RoleID    int64  `json:"role_id" validate:"required,min=1"`
        CompanyID int64  `json:"company_id" validate:"required,min=1"`
        BranchID  *int64 `json:"branch_id"`
    }{},
}
```

**Bulk Assign Roles:**
```go
var BulkAssignRolesValidation = middleware.ValidationRules{
    Body: &struct {
        UserIDs   []int64 `json:"user_ids" validate:"required,min=1"`
        RoleID    int64   `json:"role_id" validate:"required,min=1"`
        CompanyID int64   `json:"company_id" validate:"required,min=1"`
        BranchID  *int64  `json:"branch_id"`
    }{},
}
```

**Update Role Modules:**
```go
var UpdateRoleModulesValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "roleId", Type: "int", Required: true, Min: IntPtr(1)},
    },
    Body: &struct {
        Modules []struct {
            ModuleID  int64 `json:"module_id" validate:"required,min=1"`
            CanRead   bool  `json:"can_read"`
            CanWrite  bool  `json:"can_write"`
            CanDelete bool  `json:"can_delete"`
        } `json:"modules" validate:"required,min=1"`
    }{},
}
```

### Module Management Validations (`module_validation.go`)

**Create Module:**
```go
var CreateModuleValidation = middleware.ValidationRules{
    Body: &struct {
        Category         string `json:"category" validate:"required,min=2,max=50"`
        Name             string `json:"name" validate:"required,min=2,max=100"`
        URL              string `json:"url" validate:"required,min=1,max=255"`
        Icon             string `json:"icon" validate:"max=50"`
        Description      string `json:"description" validate:"max=500"`
        ParentID         *int64 `json:"parent_id"`
        SubscriptionTier string `json:"subscription_tier" validate:"required,oneof=basic pro enterprise"`
    }{},
}
```

### Subscription Validations (`subscription_validation.go`)

**Create Subscription:**
```go
var CreateSubscriptionValidation = middleware.ValidationRules{
    Body: &struct {
        CompanyID    int64  `json:"company_id" validate:"required,min=1"`
        PlanID       int64  `json:"plan_id" validate:"required,min=1"`
        BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
    }{},
}
```

**Renew Subscription:**
```go
var RenewSubscriptionValidation = middleware.ValidationRules{
    Params: IDValidation.Params,
    Body: &struct {
        BillingCycle string `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
        PlanID       *int64 `json:"plan_id"`
    }{},
}
```

## Validation Middleware

### Middleware Configuration

**Lokasi:** `middleware/validation.go`

**Fungsi Utama:**
- `ValidateRequest()` - Main validation function
- `validateParams()` - Path parameter validation
- `validateQuery()` - Query parameter validation  
- `validateBody()` - Request body validation

### Validation Flow

1. **Parameter Validation:**
   ```go
   // Validate path parameters (/users/:id)
   if err := validateParams(c, rules.Params); err != nil {
       response.Error(c, http.StatusBadRequest, "Invalid path parameter", err.Error())
       c.Abort()
       return
   }
   ```

2. **Query Validation:**
   ```go
   // Validate query parameters (?page=1&limit=10)
   if err := validateQuery(c, rules.Query); err != nil {
       response.Error(c, http.StatusBadRequest, "Invalid query parameter", err.Error())
       c.Abort()
       return
   }
   ```

3. **Body Validation:**
   ```go
   // Validate request body (JSON)
   if rules.Body != nil {
       if err := validateBody(c, rules.Body); err != nil {
           response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
           c.Abort()
           return
       }
   }
   ```

### Handler Integration

**Pattern dalam Handler:**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Get validated body from context (set by validation middleware)
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }

    // Type assert to the expected struct
    req, ok := validatedBody.(*struct {
        Name         string  `json:"name" validate:"required,min=2,max=100"`
        Email        string  `json:"email" validate:"required,email,max=255"`
        UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
        Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
    })
    if !ok {
        response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
        return
    }

    // Convert to service request
    createReq := &service.CreateUserRequest{
        Name:         req.Name,
        Email:        req.Email,
        UserIdentity: req.UserIdentity,
        Password:     req.Password,
    }

    // Call service
    result, err := h.userService.CreateUser(createReq)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "User created successfully", result)
}
```

## Route Integration

### Usage dalam Routes

**Before (Inline Validation):**
```go
// Validation rules defined inline di routes (1400+ lines)
loginValidation := middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
    }{},
}
auth.POST("/login", middleware.ValidateRequest(loginValidation), authHandler.Login)
```

**After (Modular Validation):**
```go
// Import validation rules
import "gin-scalable-api/internal/validation"

// Clean route definition (~200 lines)
auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)
```

### Route Examples

**Authentication Routes:**
```go
func setupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
    auth := api.Group("/auth")
    {
        auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)
        auth.POST("/login-email", middleware.ValidateRequest(validation.LoginEmailValidation), authHandler.LoginWithEmail)
        auth.POST("/refresh", middleware.ValidateRequest(validation.RefreshValidation), authHandler.RefreshToken)
        auth.POST("/logout", middleware.ValidateRequest(validation.LogoutValidation), authHandler.Logout)
    }
}
```

**User Routes:**
```go
func setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler) {
    users := api.Group("/users")
    {
        users.GET("", middleware.ValidateRequest(validation.UserListValidation), userHandler.GetUsers)
        users.GET("/:id", middleware.ValidateRequest(validation.IDValidation), userHandler.GetUserByID)
        users.POST("", middleware.ValidateRequest(validation.CreateUserValidation), userHandler.CreateUser)
        users.PUT("/:id", middleware.ValidateRequest(middleware.ValidationRules{
            Params: validation.IDValidation.Params,
            Body:   validation.UpdateUserValidation.Body,
        }), userHandler.UpdateUser)
        users.POST("/check-access", middleware.ValidateRequest(validation.AccessCheckValidation), userHandler.CheckAccess)
    }
}
```

## Validation Tags

### Supported Tags

**Basic Validation:**
- `required` - Field harus ada
- `omitempty` - Field boleh kosong
- `min=N` - Minimum length/value
- `max=N` - Maximum length/value
- `len=N` - Exact length

**String Validation:**
- `email` - Valid email format
- `oneof=val1 val2` - Must be one of specified values

**Numeric Validation:**
- `min=N` - Minimum value
- `max=N` - Maximum value

**Cross-field Validation:**
- `eqfield=FieldName` - Must equal another field (untuk confirm password)

### Custom Error Messages

**Error Format:**
```go
func formatValidationError(err validator.FieldError) string {
    field := err.Field()
    tag := err.Tag()

    switch tag {
    case "required":
        return fmt.Sprintf("field '%s' is required", field)
    case "email":
        return fmt.Sprintf("field '%s' must be a valid email", field)
    case "min":
        return fmt.Sprintf("field '%s' must be at least %s characters/value", field, err.Param())
    case "max":
        return fmt.Sprintf("field '%s' must be at most %s characters/value", field, err.Param())
    case "oneof":
        return fmt.Sprintf("field '%s' must be one of: %s", field, err.Param())
    default:
        return fmt.Sprintf("field '%s' failed validation for '%s'", field, tag)
    }
}
```

## Error Responses

### Validation Error Format

**Single Error:**
```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "field 'email' must be a valid email",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

**Multiple Errors:**
```json
{
  "success": false,
  "message": "Invalid request body", 
  "error": "field 'name' is required; field 'email' must be a valid email; field 'password' must be at least 6 characters",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

### Parameter Validation Errors

**Invalid Path Parameter:**
```json
{
  "success": false,
  "message": "Invalid path parameter",
  "error": "parameter 'id' must be a valid integer",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

**Invalid Query Parameter:**
```json
{
  "success": false,
  "message": "Invalid query parameter",
  "error": "query parameter 'limit' must be at most 100",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

## Best Practices

### 1. Validation Rules Organization

**Group by Domain:**
- Pisahkan validation rules berdasarkan domain/module
- Gunakan common validation untuk rules yang sama
- Reuse validation rules sebisa mungkin

**Naming Convention:**
- `CreateXXXValidation` - Untuk create operations
- `UpdateXXXValidation` - Untuk update operations
- `XXXListValidation` - Untuk list operations dengan pagination

### 2. Handler Implementation

**Always Use Validated Data:**
```go
// ✅ Good - Use validated data from context
validatedBody, exists := c.Get("validated_body")
if !exists {
    response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
    return
}

// ❌ Bad - Direct binding bypasses validation
var req CreateUserRequest
if err := c.ShouldBindJSON(&req); err != nil {
    // This bypasses the validation middleware
}
```

**Type Assertion:**
```go
// ✅ Good - Proper type assertion with error handling
req, ok := validatedBody.(*struct {
    Name string `json:"name" validate:"required"`
})
if !ok {
    response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
    return
}

// ❌ Bad - Direct type assertion without error handling
req := validatedBody.(*CreateUserRequest) // Panic if wrong type
```

### 3. Validation Rules Design

**Comprehensive Validation:**
```go
// ✅ Good - Comprehensive validation
var CreateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name         string  `json:"name" validate:"required,min=2,max=100"`
        Email        string  `json:"email" validate:"required,email,max=255"`
        UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
        Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
    }{},
}

// ❌ Bad - Minimal validation
var CreateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required"`
    }{},
}
```

### 4. Error Handling

**Consistent Error Response:**
```go
// ✅ Good - Use centralized response
response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())

// ❌ Bad - Manual error response
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
```

## Testing Validation

### Unit Testing Validation Rules

```go
func TestCreateUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        body    interface{}
        wantErr bool
    }{
        {
            name: "valid user data",
            body: struct {
                Name  string `json:"name" validate:"required,min=2,max=100"`
                Email string `json:"email" validate:"required,email"`
            }{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        {
            name: "invalid email",
            body: struct {
                Name  string `json:"name" validate:"required,min=2,max=100"`
                Email string `json:"email" validate:"required,email"`
            }{
                Name:  "John Doe",
                Email: "invalid-email",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate.Struct(tt.body)
            if (err != nil) != tt.wantErr {
                t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Testing

```go
func TestCreateUserEndpoint(t *testing.T) {
    // Test valid request
    body := `{"name":"John Doe","email":"john@example.com"}`
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    
    // Test invalid request
    invalidBody := `{"name":"","email":"invalid-email"}`
    req = httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(invalidBody))
    req.Header.Set("Content-Type", "application/json")
    
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

## Kesimpulan

Sistem validasi modular memberikan:

- ✅ **Clean Architecture** - Validation rules terpisah dari routes
- ✅ **Reusability** - Rules dapat digunakan kembali
- ✅ **Maintainability** - Mudah update dan maintain
- ✅ **Consistency** - Format error yang konsisten
- ✅ **Developer Experience** - Routes yang lebih bersih dan readable
- ✅ **Type Safety** - Compile-time validation rule checking
- ✅ **Comprehensive Validation** - Parameter, query, dan body validation

Sistem ini memberikan foundation yang solid untuk API validation yang scalable dan maintainable.