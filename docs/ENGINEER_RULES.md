# Backend Engineer Rules - RBAC Service

Panduan lengkap untuk bekerja dengan RBAC Service yang menggunakan **Module-Based Architecture** (Express.js style) dengan **5-file structure per module** setelah refactoring dari 7-file structure.

## 📋 Table of Contents
1. [Quick Start](#quick-start)
2. [Module-Based Architecture](#module-based-architecture)
3. [Development Workflow](#development-workflow)
4. [Module Structure (5 Files)](#module-structure-5-files)
5. [Database & Migrations](#database--migrations)
6. [Testing](#testing)
7. [Deployment](#deployment)

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- Air (untuk hot reload): `go install github.com/cosmtrek/air@latest`

### Setup
```bash
# 1. Clone repository
git clone <repo-url>
cd rbac-service

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env dengan konfigurasi database dan redis

# 4. Setup database
createdb huminor_rbac
make migrate-up

# 5. Run server
air  # atau go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8081`

## 🏗️ Module-Based Architecture

Project ini menggunakan **vertical module-based structure** (Express.js style), bukan horizontal layer-based.

**Prinsip: 1 fitur = 1 folder**

**✅ REFACTORING COMPLETED**: Struktur telah berhasil direfactor dari 7-file menjadi 5-file per module untuk meningkatkan developer experience.

```
rbac-service/
├── cmd/
│   ├── api/main.go          # Entry point aplikasi
│   └── migrate/main.go      # Migration tool
│
├── internal/
│   ├── app/
│   │   ├── server.go        # Server initialization
│   │   └── routes.go        # Route registration
│   │
│   ├── modules/             # 🔥 SEMUA FITUR DI SINI (5 files per module)
│   │   ├── auth/            # Authentication module
│   │   │   ├── dto.go       # Request/Response DTOs + validation logic
│   │   │   ├── model.go     # Local models
│   │   │   ├── repository.go # Database queries
│   │   │   ├── route.go     # Routes + HTTP handlers
│   │   │   └── service.go   # Business logic
│   │   │
│   │   ├── user/            # User management
│   │   ├── role/            # Role management
│   │   ├── company/         # Company management
│   │   ├── branch/          # Branch management
│   │   ├── module/          # Module system
│   │   ├── unit/            # Unit management
│   │   ├── subscription/    # Subscription system
│   │   └── audit/           # Audit logging
│   │
│   └── constants/           # Shared constants
│
├── middleware/              # HTTP middleware (auth, CORS, rate limit)
├── migrations/              # SQL migrations
├── pkg/                     # Reusable utilities
│   ├── model/              # Base repository helper
│   ├── response/           # Response helpers
│   └── utils/              # General utilities
│
├── config/                  # Configuration
└── docs/                    # Documentation
```

**Key Changes dari Refactoring:**
- ❌ Tidak ada lagi `handler.go` terpisah → digabung ke `route.go`
- ❌ Tidak ada lagi `validator.go` terpisah → digabung ke `dto.go`
- ✅ File count berkurang: 63 → 45 files (28% reduction)
- ✅ Validation menggunakan `middleware.ValidateRequest` dengan `ValidationRules`
- ✅ Semua routes memiliki dokumentasi komentar yang lengkap

## 🔄 Development Workflow

Ketika mengembangkan fitur baru dalam **module-based architecture**:

### 1. Buat Module Folder
```bash
mkdir -p internal/modules/feature_name
```

### 2. Buat File Structure (5 files dalam 1 folder)
```bash
cd internal/modules/feature_name
touch dto.go model.go repository.go route.go service.go
```

### 3. Implementasi (urutan yang disarankan)

**a. Model (`model.go`)** - Database entities
```go
package feature_name

type Feature struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**b. DTO (`dto.go`)** - Request/Response structures + validation logic
```go
package feature_name

type CreateFeatureRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
}

type FeatureResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
}

// Validation rules (merged from validator.go)
func (r *CreateFeatureRequest) ValidationRules() map[string]string {
    return map[string]string{
        "name": "required,min=2,max=100",
    }
}
```

**c. Repository (`repository.go`)** - Database queries (raw SQL)
```go
package feature_name

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(feature *Feature) error {
    // Raw SQL query
}
```

**d. Service (`service.go`)** - Business logic
```go
package feature_name

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateFeature(req *CreateFeatureRequest) (*FeatureResponse, error) {
    // Business logic
}
```

**e. Route (`route.go`)** - Route registration + HTTP handlers
```go
package feature_name

import (
    "gin-scalable-api/middleware"
    "gin-scalable-api/pkg/response"
    "net/http"
    "github.com/gin-gonic/gin"
)

// Handler struct (merged from handler.go)
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// Handler methods (merged from handler.go)
func (h *Handler) CreateFeature(c *gin.Context) {
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }

    req, ok := validatedBody.(*CreateFeatureRequest)
    if !ok {
        response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
        return
    }

    result, err := h.service.CreateFeature(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create feature", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "Feature created", result)
}

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    features := router.Group("/features")
    {
        // POST /api/v1/features - Create new feature
        features.POST("", 
            middleware.ValidateRequest(middleware.ValidationRules{
                Body: &CreateFeatureRequest{},
            }),
            handler.CreateFeature,
        )
        
        // GET /api/v1/features/:id - Get feature by ID
        features.GET("/:id", handler.GetFeature)
    }
}
```

### 4. Register Module di `internal/app/routes.go`
```go
featureModule "gin-scalable-api/internal/modules/feature_name"

// Di SetupNewModuleRoutes:
featureModule.RegisterRoutes(protected, h.Feature)
```

### 5. Testing
- Test dengan Postman collection
- Verify dengan `go build ./cmd/api`

## 🧩 Module Structure (5 Files)

Setiap module memiliki 5 file standar dalam 1 folder setelah refactoring:

### 📝 1. model.go - Database Entities
**Lokasi**: `internal/modules/{module}/model.go`

Model mendefinisikan struktur data database. Setiap module punya model lokalnya sendiri.

```go
package user

import "time"

type User struct {
    ID           int64     `json:"id" db:"id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    UserIdentity *string   `json:"user_identity" db:"user_identity"`
    PasswordHash string    `json:"-" db:"password_hash"`
    IsActive     bool      `json:"is_active" db:"is_active"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (User) TableName() string {
    return "users"
}
```

**Best Practices:**
- Model lokal per module (tidak shared)
- Gunakan pointer untuk nullable fields (`*int64`, `*string`)
- Tag `json:"-"` untuk field sensitif (password)
- Method `TableName()` untuk explicit table name

### 📦 2. dto.go - Request/Response DTOs + Validation
**Lokasi**: `internal/modules/{module}/dto.go`

DTOs mendefinisikan struktur data untuk API request dan response, serta validation logic yang sebelumnya ada di `validator.go`.

```go
package user

// Request DTO
type CreateUserRequest struct {
    Name         string  `json:"name" validate:"required,min=2,max=100"`
    Email        string  `json:"email" validate:"required,email"`
    UserIdentity *string `json:"user_identity" validate:"omitempty"`
    Password     string  `json:"password" validate:"required,min=8"`
}

// Response DTO
type UserResponse struct {
    ID           int64   `json:"id"`
    Name         string  `json:"name"`
    Email        string  `json:"email"`
    UserIdentity *string `json:"user_identity,omitempty"`
    IsActive     bool    `json:"is_active"`
    CreatedAt    string  `json:"created_at"`
    UpdatedAt    string  `json:"updated_at"`
}

// List Response DTO
type UserListResponse struct {
    Data    []*UserResponse `json:"data"`
    Total   int64           `json:"total"`
    Limit   int             `json:"limit"`
    Offset  int             `json:"offset"`
    HasMore bool            `json:"has_more"`
}

// Validation logic (merged from validator.go)
func ValidateUserIdentity(fl validator.FieldLevel) bool {
    identity := fl.Field().String()
    return len(identity) == 9 // Example: must be 9 digits
}
```

**Validation Tags:**
- `required` - Field wajib diisi
- `email` - Validasi format email
- `min=n,max=n` - Panjang minimum/maksimum
- `omitempty` - Field opsional

### 🗃️ 3. repository.go - Database Queries
**Lokasi**: `internal/modules/{module}/repository.go`

Repository menangani semua operasi database dengan raw SQL.

```go
package user

import (
    "database/sql"
    "fmt"
    "gin-scalable-api/pkg/model"
)

type Repository struct {
    *model.Repository
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{
        Repository: model.NewRepository(db),
        db:         db,
    }
}

func (r *Repository) Create(user *User) error {
    query := `
        INSERT INTO users (name, email, user_identity, password_hash, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(query, user.Name, user.Email, user.UserIdentity, user.PasswordHash, user.IsActive).Scan(
        &user.ID, &user.CreatedAt, &user.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func (r *Repository) GetByID(id int64) (*User, error) {
    user := &User{}
    query := `
        SELECT id, name, email, user_identity, password_hash, is_active, created_at, updated_at
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL
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

**Best Practices:**
- Gunakan raw SQL (tidak pakai ORM)
- Handle `sql.ErrNoRows` untuk data tidak ditemukan
- Gunakan soft delete dengan `deleted_at`
- Return error yang descriptive

### 🏢 5. service.go - Business Logic
**Lokasi**: `internal/modules/{module}/service.go`

Service layer berisi business logic dan orchestrate repository calls.

```go
package user

import (
    "fmt"
    "time"
    "golang.org/x/crypto/bcrypt"
)

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user model
    user := &User{
        Name:         req.Name,
        Email:        req.Email,
        UserIdentity: req.UserIdentity,
        PasswordHash: string(hashedPassword),
        IsActive:     true,
    }
    
    // Save to database
    if err := s.repo.Create(user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Convert to response
    return &UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        UserIdentity: user.UserIdentity,
        IsActive:     user.IsActive,
        CreatedAt:    user.CreatedAt.Format(time.RFC3339),
        UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
    }, nil
}
```

**Best Practices:**
- Tidak ada struct definitions di service (harus di DTO/Model)
- Focus pada business logic
- Konversi Model ↔ DTO dilakukan inline (tidak perlu mapper terpisah)

### 🛣️ 4. route.go - Route Registration + HTTP Handlers
**Lokasi**: `internal/modules/{module}/route.go`

Routes mendefinisikan endpoint, middleware, dan HTTP handlers (merged dari `handler.go`).

```go
package user

import (
    "gin-scalable-api/internal/constants"
    "gin-scalable-api/middleware"
    "gin-scalable-api/pkg/response"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

// Handler struct (merged from handler.go)
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// Handler methods (merged from handler.go)
func (h *Handler) CreateUser(c *gin.Context) {
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }
    
    req, ok := validatedBody.(*CreateUserRequest)
    if !ok {
        response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
        return
    }
    
    result, err := h.service.CreateUser(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create user", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, constants.MsgUserCreated, result)
}

func (h *Handler) GetUserByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Bad request", "Invalid ID")
        return
    }
    
    result, err := h.service.GetUserByID(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, constants.MsgUserNotFound, err.Error())
        return
    }
    
    response.Success(c, http.StatusOK, constants.MsgUserRetrieved, result)
}

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    users := router.Group("/users")
    {
        // GET /api/v1/users - Get all users with optional filters and pagination
        users.GET("", handler.GetUsers)
        
        // POST /api/v1/users - Create new user with role assignments
        users.POST("", 
            middleware.ValidateRequest(middleware.ValidationRules{
                Body: &CreateUserRequest{},
            }),
            handler.CreateUser,
        )
        
        // GET /api/v1/users/:id - Get user by ID with roles and permissions
        users.GET("/:id", handler.GetUserByID)
        
        // PUT /api/v1/users/:id - Update user information and role assignments
        users.PUT("/:id",
            middleware.ValidateRequest(middleware.ValidationRules{
                Body: &UpdateUserRequest{},
            }),
            handler.UpdateUser,
        )
        
        // DELETE /api/v1/users/:id - Delete user and remove all associations
        users.DELETE("/:id", handler.DeleteUser)
    }
}
```

**Best Practices:**
- Gunakan `middleware.ValidateRequest` dengan `ValidationRules` (bukan `ValidateJSON`)
- Tambahkan komentar deskriptif untuk setiap route
- Handle errors dengan appropriate HTTP status codes
- Gunakan constants untuk messages

## 🗄️ Database & Migrations

### Membuat Migration Baru
```bash
# Buat migration file
migrate create -ext sql -dir migrations -seq add_examples_table

# Edit file migration
# migrations/XXX_add_examples_table.up.sql
# migrations/XXX_add_examples_table.down.sql

# Jalankan migration
make migrate-up

# Rollback migration
make migrate-down
```

### Unit-Based RBAC Schema
Project ini menggunakan hierarchical RBAC:
- **Company** → **Branch** → **Unit** → **Role** → **User**

**Key Tables:**
- `companies` - Perusahaan
- `branches` - Cabang dalam perusahaan
- `units` - Unit dalam cabang (HR, Finance, etc.)
- `roles` - Role/jabatan
- `users` - User/pengguna
- `user_roles` - Assignment user ke role dengan scope (company/branch/unit)

## 🧪 Testing

### Manual Testing dengan curl
```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "100000003", "password": "password123"}'

# Get user with role assignments
curl -X GET http://localhost:8081/api/v1/users/3 \
  -H "Authorization: Bearer <token>"
```

### Postman Collection
Gunakan file `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json` untuk testing lengkap.

## 🚀 Deployment

### Build Production
```bash
# Build binary
go build -o server cmd/api/main.go

# Run production
GIN_MODE=release ./server
```

### Docker
```bash
# Build image
docker build -t rbac-service .

# Run container
docker run -p 8081:8081 rbac-service
```

## 📚 Additional Resources

- [Clean Architecture Guide](docs/CLEAN_ARCHITECTURE.md)
- [Unit-Based RBAC Documentation](docs/UNIT_BASED_RBAC.md)
- [API Documentation](docs/API_OVERVIEW.md)
- [Migration Guide](docs/MIGRATIONS.md)

## 🔧 Common Issues

### Database Connection Issues
- Pastikan PostgreSQL running
- Check connection string di `.env`
- Verify database exists

### Redis Connection Issues
- Pastikan Redis running
- Check Redis configuration di `.env`

### Migration Issues
- Check migration files syntax
- Verify database permissions
- Use `make migrate-force <version>` untuk fix dirty state

## 💡 Tips & Best Practices

1. **Selalu gunakan interfaces** untuk dependency injection
2. **Jangan skip validation** - gunakan validation middleware
3. **Handle errors dengan baik** - return descriptive error messages
4. **Gunakan constants** untuk messages dan status codes
5. **Follow naming conventions** - konsisten dengan existing code
6. **Test setiap endpoint** sebelum commit
7. **Dokumentasikan API changes** di Postman collection