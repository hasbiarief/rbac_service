# Arsitektur Modular - ERP RBAC API

## Gambaran Umum

Project ini telah direstrukturisasi menjadi arsitektur modular yang lebih terorganisir dan mudah dipelihara. Setiap modul memiliki tanggung jawab yang jelas dan terpisah dengan infrastruktur yang terpusat.

## Struktur Direktori Terbaru

```
internal/
├── handlers/           # HTTP request handlers berdasarkan modul
│   ├── auth_handler.go
│   ├── module_handler.go
│   ├── user_handler.go
│   ├── company_handler.go
│   ├── role_handler.go
│   ├── subscription_handler.go
│   ├── audit_handler.go
│   └── branch_handler.go
├── routes/             # Konfigurasi routing (clean & minimal)
│   └── routes.go
├── server/             # Server initialization dan dependency injection
│   └── server.go
├── validation/         # Modular validation rules (NEW!)
│   ├── common_validation.go
│   ├── auth_validation.go
│   ├── user_validation.go
│   ├── company_validation.go
│   ├── role_validation.go
│   ├── module_validation.go
│   ├── subscription_validation.go
│   ├── audit_validation.go
│   └── branch_validation.go
├── service/            # Business logic layer
├── repository/         # Data access layer (Raw SQL)
└── models/             # Data models

pkg/
├── database/           # Database connection (Raw SQL)
├── migration/          # File-based migration system
├── model/              # Base model dengan helper methods
├── response/           # Centralized response system
└── ratelimiter/        # Rate limiting system

middleware/
├── auth.go             # JWT authentication
├── cors.go             # CORS configuration
├── rate_limit.go       # Rate limiting middleware
└── validation.go       # Request validation middleware
```

## Komponen Utama

### 1. Handlers (`internal/handlers/`)

Setiap handler bertanggung jawab untuk:
- Memproses HTTP requests
- Menggunakan validated body dari middleware
- Memanggil service layer
- Mengembalikan centralized response format

**Modul Handlers:**
- `AuthHandler` - Authentication (login dual method, logout, refresh token)
- `ModuleHandler` - Module management (CRUD, tree, hierarchy)
- `UserHandler` - User management (CRUD, password, module access)
- `CompanyHandler` - Company management (CRUD)
- `RoleHandler` - Role management (basic & advanced RBAC)
- `SubscriptionHandler` - Subscription management (plans, billing)
- `AuditHandler` - Audit logging (logs, statistics)
- `BranchHandler` - Branch management (hierarchical structure)

**Handler Pattern (Updated):**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Get validated body from middleware
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }

    // Type assert to expected struct
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

    result, err := h.userService.CreateUser(createReq)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "User created successfully", result)
}
```

### 2. Modular Validation System (`internal/validation/`)

**Keuntungan Baru:**
- ✅ **Separation of Concerns** - Validation rules terpisah dari routes
- ✅ **Reusability** - Rules dapat digunakan kembali
- ✅ **Maintainability** - Mudah update dan maintain
- ✅ **Clean Routes** - Routes file lebih bersih dan readable

**Struktur Validation:**
```go
// common_validation.go - Shared rules
var IDValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
    },
}

// auth_validation.go - Authentication specific
var LoginValidation = middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
    }{},
}

// user_validation.go - User management specific
var CreateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name         string  `json:"name" validate:"required,min=2,max=100"`
        Email        string  `json:"email" validate:"required,email,max=255"`
        UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
        Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
    }{},
}
```

### 3. Clean Routes (`internal/routes/`)

Routes sekarang jauh lebih bersih dan readable:

**Before (1400+ lines dengan validation inline):**
```go
// Validation rules defined inline di routes
loginValidation := middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
    }{},
}
auth.POST("/login", middleware.ValidateRequest(loginValidation), authHandler.Login)
```

**After (Clean & Modular):**
```go
// Import validation rules
import "gin-scalable-api/internal/validation"

// Clean route definition
auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)
```

### 4. Raw SQL Infrastructure (`pkg/database/`, `pkg/model/`)

**Menggantikan GORM dengan Raw SQL:**
- ✅ **Better Performance** - Query yang lebih optimal
- ✅ **Full Control** - Complete control atas SQL queries
- ✅ **Transparency** - SQL queries terlihat jelas
- ✅ **Easier Debugging** - Mudah debug query issues
- ✅ **Simple Migration** - File-based migration system

**Repository Pattern:**
```go
func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (name, email, user_identity, password_hash, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(query, 
        user.Name, user.Email, user.UserIdentity, user.PasswordHash, user.IsActive,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}
```

### 5. Centralized Response System (`pkg/response/`)

**Consistent API Response Format:**
```go
// Success response
response.Success(c, http.StatusOK, "Success message", data)

// Error response  
response.Error(c, http.StatusBadRequest, "Error message", errorDetails)
```

**Response Format:**
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... },
  "timestamp": "2025-12-28T22:39:39Z"
}
```

### 6. Rate Limiting System (`pkg/ratelimiter/`, `middleware/rate_limit.go`)

**Features:**
- **100 requests per minute** per IP address
- **Redis-based** storage untuk distributed systems
- **Automatic cleanup** expired entries
- **Global application** pada semua endpoints

## Keuntungan Arsitektur Modular Terbaru

### 1. **Enhanced Separation of Concerns**
- Validation rules terpisah dari routes
- Handler hanya fokus pada HTTP processing
- Service layer untuk business logic
- Repository layer untuk data access

### 2. **Improved Maintainability**
- Routes file lebih bersih (dari 1400+ lines menjadi ~200 lines)
- Validation rules modular dan reusable
- Handler pattern yang consistent
- Centralized response format

### 3. **Better Performance**
- Raw SQL queries yang optimal
- Connection pooling
- Redis-based caching dan rate limiting
- Efficient validation middleware

### 4. **Enhanced Security**
- Centralized validation system
- Rate limiting untuk semua endpoints
- Dual authentication methods
- SQL injection prevention dengan prepared statements

### 5. **Developer Experience**
- Clean dan readable code structure
- Easy to add new modules
- Consistent patterns across modules
- Better error handling dan debugging

## Migration dari Struktur Lama

**Perubahan Utama:**

1. **GORM → Raw SQL**
   - Menghapus semua dependency GORM
   - Implementasi raw SQL dengan lib/pq
   - File-based migration system

2. **Inline Validation → Modular Validation**
   - Memindahkan validation rules ke `internal/validation/`
   - Routes menjadi lebih clean dan readable
   - Validation rules dapat digunakan kembali

3. **Manual Response → Centralized Response**
   - Semua handler menggunakan `pkg/response`
   - Consistent response format
   - Better error handling

4. **No Rate Limiting → Global Rate Limiting**
   - Implementasi rate limiting berbasis IP
   - Redis-based storage
   - 100 requests/minute limit

**Statistik Improvement:**
- **Routes file:** 1400+ lines → ~200 lines (85% reduction)
- **Validation:** Inline → Modular (8 validation files)
- **Response:** Manual → Centralized (100% consistent)
- **Database:** GORM → Raw SQL (Better performance)
- **Security:** Basic → Enhanced (Rate limiting + validation)

## Cara Menambah Modul Baru

### 1. Buat Validation Rules
```go
// internal/validation/new_module_validation.go
package validation

var CreateNewItemValidation = middleware.ValidationRules{
    Body: &struct {
        Name        string `json:"name" validate:"required,min=2,max=100"`
        Description string `json:"description" validate:"max=500"`
    }{},
}
```

### 2. Buat Handler
```go
// internal/handlers/new_module_handler.go
func (h *NewModuleHandler) CreateItem(c *gin.Context) {
    // Get validated body from middleware
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }

    // Type assert dan convert ke service request
    // ... implementation

    response.Success(c, http.StatusCreated, "Item created successfully", result)
}
```

### 3. Tambahkan Routes
```go
// internal/routes/routes.go
func setupNewModuleRoutes(protected *gin.RouterGroup, handler *handlers.NewModuleHandler) {
    items := protected.Group("/items")
    {
        items.POST("", middleware.ValidateRequest(validation.CreateNewItemValidation), handler.CreateItem)
        items.GET("", middleware.ValidateRequest(validation.ListValidation), handler.GetItems)
        items.GET("/:id", middleware.ValidateRequest(validation.IDValidation), handler.GetItemByID)
    }
}
```

## Best Practices Terbaru

### 1. **Handler Guidelines**
- Selalu gunakan validated body dari middleware
- Gunakan centralized response format
- Convert request struct ke service struct
- Handle errors dengan proper HTTP status codes

### 2. **Validation Guidelines**
- Buat validation rules di `internal/validation/`
- Gunakan common validation untuk rules yang sama
- Group validation berdasarkan domain/module
- Reuse validation rules sebisa mungkin

### 3. **Service Guidelines**
- Semua business logic ada di service layer
- Service tidak boleh tahu tentang HTTP
- Gunakan repository untuk database access
- Return domain errors, bukan HTTP errors

### 4. **Repository Guidelines**
- Gunakan raw SQL dengan prepared statements
- Return domain models
- Handle database-specific errors
- Use transactions untuk complex operations

## Testing Strategy

### 1. **Unit Testing**
- Test validation rules secara terpisah
- Test handler dengan mock services
- Test service dengan mock repositories
- Test repository dengan test database

### 2. **Integration Testing**
- Test end-to-end API flow
- Test dengan real database dan Redis
- Test authentication dan authorization
- Test rate limiting behavior

### 3. **API Testing**
- Gunakan Postman collection yang updated
- Test semua validation scenarios
- Test error handling
- Test dengan different user roles

## Kesimpulan

Arsitektur modular terbaru memberikan:

- ✅ **Clean Architecture** - Separation of concerns yang jelas
- ✅ **Better Performance** - Raw SQL dan optimized queries
- ✅ **Enhanced Security** - Rate limiting dan centralized validation
- ✅ **Improved DX** - Developer experience yang lebih baik
- ✅ **Maintainable Code** - Modular dan easy to extend
- ✅ **Consistent API** - Centralized response format
- ✅ **Scalable Infrastructure** - Ready untuk growth

Project sekarang memiliki foundation yang solid untuk pengembangan jangka panjang dengan infrastruktur yang modern dan maintainable.