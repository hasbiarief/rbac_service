# API Infrastructure Documentation

## Gambaran Umum

Project ini menggunakan infrastruktur API yang terpusat dan terstandarisasi dengan fokus pada:
- **Raw SQL** tanpa ORM (menggantikan GORM)
- **Centralized Response System** untuk konsistensi
- **Modular Validation System** yang terpisah dari routes
- **Rate Limiting** berbasis IP
- **Dual Authentication Methods** (user_identity & email)

## Arsitektur Infrastruktur

### 1. Database Layer (Raw SQL)

**Teknologi:** PostgreSQL dengan lib/pq driver (tanpa GORM)

**Struktur:**
```
pkg/
├── database/
│   └── connection.go      # Database connection management
├── migration/
│   └── migration.go       # File-based migration system
└── model/
    └── base.go           # Base model dengan helper methods
```

**Keuntungan Raw SQL:**
- ✅ **Performance** - Query yang lebih optimal
- ✅ **Control** - Full control atas SQL queries
- ✅ **Transparency** - SQL queries terlihat jelas
- ✅ **Debugging** - Mudah debug query issues
- ✅ **Migration** - File-based migration yang simple

**Contoh Implementation:**
```go
// Repository dengan Raw SQL
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
    user := &models.User{}
    query := `
        SELECT id, name, email, user_identity, password_hash, is_active, 
               created_at, updated_at
        FROM users 
        WHERE id = $1 AND is_active = true
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Name, &user.Email, &user.UserIdentity,
        &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}
```

### 2. Centralized Response System

**Lokasi:** `pkg/response/response.go`

**Format Standar:**
```json
{
  "success": true,
  "message": "Success message",
  "data": { ... },
  "timestamp": "2025-12-28T22:39:39Z"
}
```

**Error Format:**
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

**Usage:**
```go
// Success response
response.Success(c, http.StatusOK, "User created successfully", user)

// Error response
response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
```

### 3. Modular Validation System

**Lokasi:** `internal/validation/`

**Struktur:**
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

**Keuntungan:**
- ✅ **Modular** - Validation rules terpisah per domain
- ✅ **Reusable** - Rules dapat digunakan kembali
- ✅ **Maintainable** - Mudah update dan maintain
- ✅ **Clean Routes** - Routes file lebih bersih dan readable

**Contoh Usage:**
```go
// Di routes.go
auth.POST("/login", middleware.ValidateRequest(validation.LoginValidation), authHandler.Login)

// Di validation/auth_validation.go
var LoginValidation = middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
    }{},
}
```

### 4. Rate Limiting System

**Lokasi:** `pkg/ratelimiter/limiter.go` & `middleware/rate_limit.go`

**Konfigurasi:**
- **100 requests per minute** per IP address
- **Redis-based** storage untuk distributed systems
- **Automatic cleanup** expired entries

**Implementation:**
```go
// Global rate limiting
r.Use(middleware.RateLimit())

// Rate limiter configuration
const (
    RequestsPerMinute = 100
    WindowSize       = time.Minute
)
```

**Response saat limit exceeded:**
```json
{
  "success": false,
  "message": "Rate limit exceeded",
  "error": "Too many requests. Please try again later.",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

### 5. Dual Authentication System

**Primary Method:** User Identity
```bash
POST /api/v1/auth/login
{
  "user_identity": "100000001",
  "password": "password123"
}
```

**Alternative Method:** Email
```bash
POST /api/v1/auth/login-email
{
  "email": "admin@system.com", 
  "password": "password123"
}
```

**Refresh Token (No Bearer Required):**
```bash
POST /api/v1/auth/refresh
{
  "refresh_token": "your_refresh_token"
}
```

### 6. Middleware Stack

**Order & Purpose:**
1. **Rate Limiting** - Mencegah abuse
2. **CORS** - Cross-origin resource sharing
3. **Authentication** - JWT token validation (untuk protected routes)
4. **Validation** - Request validation per endpoint

**Middleware Configuration:**
```go
// Global middlewares
r.Use(middleware.RateLimit())
r.Use(middleware.CORS())

// Protected routes
protected := api.Group("")
protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
{
    // All protected endpoints here
}
```

## API Endpoints Structure

### Public Endpoints (No Auth Required)

```
GET    /health                    # Health check
POST   /api/v1/auth/login         # Login with user_identity
POST   /api/v1/auth/login-email   # Login with email
POST   /api/v1/auth/refresh       # Refresh token
POST   /api/v1/auth/logout        # Logout
GET    /api/v1/plans              # Public subscription plans
GET    /api/v1/plans/:id          # Get plan by ID
```

### Protected Endpoints (Auth Required)

**User Management:**
```
GET    /api/v1/users              # List users (paginated)
GET    /api/v1/users/:id          # Get user by ID
POST   /api/v1/users              # Create user
PUT    /api/v1/users/:id          # Update user
DELETE /api/v1/users/:id          # Delete user (soft delete)
GET    /api/v1/users/:id/modules  # Get user modules
POST   /api/v1/users/check-access # Check module access
PUT    /api/v1/users/:id/password # Change user password
PUT    /api/v1/users/me/password  # Change own password
```

**Company Management:**
```
GET    /api/v1/companies          # List companies
GET    /api/v1/companies/:id      # Get company by ID
POST   /api/v1/companies          # Create company
PUT    /api/v1/companies/:id      # Update company
DELETE /api/v1/companies/:id      # Delete company
```

**Role Management:**
```
GET    /api/v1/roles              # List roles
GET    /api/v1/roles/:id          # Get role by ID
POST   /api/v1/roles              # Create role
PUT    /api/v1/roles/:id          # Update role
DELETE /api/v1/roles/:id          # Delete role

# Advanced Role Management
POST   /api/v1/role-management/assign-user-role
POST   /api/v1/role-management/bulk-assign-roles
PUT    /api/v1/role-management/role/:roleId/modules
DELETE /api/v1/role-management/user/:userId/role/:roleId
GET    /api/v1/role-management/role/:roleId/users
GET    /api/v1/role-management/user/:userId/roles
GET    /api/v1/role-management/user/:userId/access-summary
```

**Module Management:**
```
GET    /api/v1/modules            # List modules (paginated)
GET    /api/v1/modules/:id        # Get module by ID
POST   /api/v1/modules            # Create module
PUT    /api/v1/modules/:id        # Update module
DELETE /api/v1/modules/:id        # Delete module
GET    /api/v1/modules/tree       # Get module tree
GET    /api/v1/modules/:id/children    # Get module children
GET    /api/v1/modules/:id/ancestors   # Get module ancestors
```

**Subscription Management:**
```
GET    /api/v1/subscription/subscriptions
POST   /api/v1/subscription/subscriptions
GET    /api/v1/subscription/subscriptions/:id
PUT    /api/v1/subscription/subscriptions/:id
POST   /api/v1/subscription/subscriptions/:id/renew
POST   /api/v1/subscription/subscriptions/:id/cancel
GET    /api/v1/subscription/companies/:id/subscription
GET    /api/v1/subscription/companies/:id/status
GET    /api/v1/subscription/module-access/:companyId/:moduleId
GET    /api/v1/subscription/stats
GET    /api/v1/subscription/expiring
POST   /api/v1/subscription/update-expired
```

**Audit Logging:**
```
GET    /api/v1/audit/logs         # List audit logs
POST   /api/v1/audit/logs         # Create audit log
GET    /api/v1/audit/users/:userId/logs
GET    /api/v1/audit/users/identity/:identity/logs
GET    /api/v1/audit/stats        # Audit statistics
```

**Branch Management:**
```
GET    /api/v1/branches           # List branches
GET    /api/v1/branches/:id       # Get branch by ID
POST   /api/v1/branches           # Create branch
PUT    /api/v1/branches/:id       # Update branch
DELETE /api/v1/branches/:id       # Delete branch
GET    /api/v1/branches/company/:companyId
GET    /api/v1/branches/:id/children
```

## Validation Rules

### Common Validations

**ID Validation:**
```go
var IDValidation = middleware.ValidationRules{
    Params: []middleware.ParamValidation{
        {Name: "id", Type: "int", Required: true, Min: IntPtr(1)},
    },
}
```

**Pagination:**
```go
var ListValidation = middleware.ValidationRules{
    Query: []middleware.QueryValidation{
        {Name: "page", Type: "int", Default: 1, Min: IntPtr(1)},
        {Name: "limit", Type: "int", Default: 10, Min: IntPtr(1), Max: IntPtr(100)},
        {Name: "search", Type: "string"},
    },
}
```

### Authentication Validations

**Login Validation:**
```go
var LoginValidation = middleware.ValidationRules{
    Body: &struct {
        UserIdentity string `json:"user_identity" validate:"required"`
        Password     string `json:"password" validate:"required,min=6"`
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

## Error Handling

### Standard Error Codes

- **400 Bad Request** - Validation errors, malformed requests
- **401 Unauthorized** - Authentication required or failed
- **403 Forbidden** - Insufficient permissions
- **404 Not Found** - Resource not found
- **409 Conflict** - Duplicate resource (email, user_identity)
- **429 Too Many Requests** - Rate limit exceeded
- **500 Internal Server Error** - Server errors

### Error Response Format

```json
{
  "success": false,
  "message": "Human readable error message",
  "error": "Technical error details",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

### Validation Error Example

```json
{
  "success": false,
  "message": "Invalid request body",
  "error": "field 'user_identity' is required; field 'password' must be at least 6 characters",
  "timestamp": "2025-12-28T22:39:39Z"
}
```

## Performance Optimizations

### Database Optimizations

1. **Connection Pooling** - Configured di database connection
2. **Prepared Statements** - Untuk query yang sering digunakan
3. **Indexes** - Pada kolom yang sering di-query
4. **Pagination** - Semua list endpoints menggunakan pagination

### Caching Strategy

1. **Redis Caching** - Untuk session management dan rate limiting
2. **Query Result Caching** - Untuk data yang jarang berubah
3. **Connection Pooling** - Database dan Redis connections

### Rate Limiting

- **100 requests/minute** per IP address
- **Redis-based** untuk distributed systems
- **Automatic cleanup** untuk expired entries

## Security Features

### Authentication & Authorization

1. **JWT Tokens** - Access token (15 menit) + Refresh token (7 hari)
2. **Role-Based Access Control** - Granular permissions per module
3. **Password Hashing** - bcrypt dengan cost 12
4. **Session Management** - Redis-based session storage

### Input Validation

1. **Centralized Validation** - Middleware-based validation
2. **SQL Injection Prevention** - Prepared statements
3. **XSS Prevention** - Input sanitization
4. **CORS Configuration** - Proper cross-origin settings

### Rate Limiting & Monitoring

1. **IP-based Rate Limiting** - 100 requests/minute
2. **Audit Logging** - Semua API calls logged
3. **Error Monitoring** - Structured error logging
4. **Health Checks** - `/health` endpoint untuk monitoring

## Development Guidelines

### Adding New Endpoints

1. **Create Validation Rules** di `internal/validation/`
2. **Implement Handler** di `internal/handlers/`
3. **Add Routes** di `internal/routes/`
4. **Update Documentation** dan Postman collection

### Testing Strategy

1. **Unit Tests** - Handler, service, repository layers
2. **Integration Tests** - End-to-end API testing
3. **Postman Collection** - Manual API testing
4. **Load Testing** - Performance validation

### Monitoring & Debugging

1. **Structured Logging** - JSON format untuk production
2. **Health Checks** - Database, Redis, application health
3. **Metrics Collection** - Response times, error rates
4. **Debug Mode** - Detailed logging untuk development

---

**Catatan Penting:**
- Semua endpoints menggunakan centralized response format
- Validation rules terpisah dan modular
- Rate limiting aktif untuk semua endpoints
- Authentication menggunakan dual method (user_identity/email)
- Database menggunakan raw SQL tanpa ORM