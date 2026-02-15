# Developer Guide - RBAC Service

Panduan lengkap untuk development RBAC Service dengan module-based architecture.

## 🚀 Quick Setup

```bash
# 1. Clone & install
git clone <repo-url> && cd rbac-service
make dev-setup

# 2. Configure .env
cp .env.example .env
# Edit: DB_HOST, DB_USER, DB_PASSWORD, REDIS_HOST, JWT_SECRET

# 3. Setup database
make db-seed-fresh

# 4. Run server
make dev
```

Server: `http://localhost:8081`
Swagger: `http://localhost:8081/swagger/index.html`

**Default Login:** `800000001` / `password123`

## 🏗️ Architecture

### Module Structure (5 files per module)

```
internal/modules/<module>/
├── dto.go          # Request/Response + validation
├── model.go        # Database entities
├── repository.go   # Database queries (raw SQL)
├── route.go        # Routes + HTTP handlers
└── service.go      # Business logic
```

### Existing Modules
auth, user, company, branch, role, module, subscription, unit, application, audit

## 📝 Creating New Module

### Automatic
```bash
make newmodule name=employee
make listmodules
make removemodule name=employee
```

### Manual
```bash
mkdir -p internal/modules/employee
cd internal/modules/employee
touch dto.go model.go repository.go route.go service.go
```

## 📋 Code Templates

### model.go
```go
package employee

import "time"

type Employee struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### dto.go
```go
package employee

type CreateEmployeeRequest struct {
    Name string `json:"name" validate:"required,min=2,max=100"`
}

type EmployeeResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
}
```

### repository.go
```go
package employee

import "database/sql"

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(emp *Employee) error {
    query := `INSERT INTO employees (name, is_active, created_at, updated_at)
              VALUES ($1, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              RETURNING id, created_at, updated_at`
    return r.db.QueryRow(query, emp.Name).Scan(&emp.ID, &emp.CreatedAt, &emp.UpdatedAt)
}
```

### service.go
```go
package employee

import "time"

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateEmployee(req *CreateEmployeeRequest) (*EmployeeResponse, error) {
    emp := &Employee{Name: req.Name}
    if err := s.repo.Create(emp); err != nil {
        return nil, err
    }
    return &EmployeeResponse{
        ID:        emp.ID,
        Name:      emp.Name,
        CreatedAt: emp.CreatedAt.Format(time.RFC3339),
    }, nil
}
```

### route.go
```go
package employee

import (
    "gin-scalable-api/middleware"
    "gin-scalable-api/pkg/response"
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// @Summary      Create employee
// @Tags         Employees
// @Accept       json
// @Produce      json
// @Param        employee  body      CreateEmployeeRequest  true  "Employee data"
// @Success      201       {object}  response.Response{data=EmployeeResponse}
// @Router       /api/v1/employees [post]
// @Security     BearerAuth
func (h *Handler) CreateEmployee(c *gin.Context) {
    validatedBody, _ := c.Get("validated_body")
    req := validatedBody.(*CreateEmployeeRequest)
    
    result, err := h.service.CreateEmployee(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create employee", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, "Employee created", result)
}

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    employees := router.Group("/employees")
    {
        employees.POST("", 
            middleware.ValidateRequest(middleware.ValidationRules{
                Body: &CreateEmployeeRequest{},
            }),
            handler.CreateEmployee,
        )
    }
}
```

## 🔗 Register Module

**internal/app/routes.go:**
```go
employeeModule "gin-scalable-api/internal/modules/employee"

// In SetupNewModuleRoutes:
employeeModule.RegisterRoutes(protected, h.Employee)
```

**internal/app/server.go:**
```go
// In NewModuleHandlers:
Employee: employeeModule.NewHandler(
    employeeModule.NewService(
        employeeModule.NewRepository(s.db),
    ),
),
```

## 🗄️ Database

```bash
# Create migration
migrate create -ext sql -dir migrations -seq add_employees

# Run migrations
make migrate-up

# Seed data
make db-seed
```

## 📖 Swagger

### Generate Documentation
```bash
make swagger-gen
```

### Access Swagger UI
```
http://localhost:8081/swagger/index.html
```

### Export to Postman

**Method 1: Direct Import (Recommended)**
1. Buka Postman
2. Click **Import** button (top left)
3. Pilih **File** tab
4. Upload file `docs/swagger.json` atau `docs/swagger.yaml`
5. Postman akan otomatis convert ke Postman Collection
6. Done! Collection siap digunakan

**Method 2: Via URL (jika server running)**
1. Buka Postman
2. Click **Import** → **Link** tab
3. Paste URL: `http://localhost:8081/swagger/doc.json`
4. Click **Continue** → **Import**

**Setup Environment:**
```
base_url: http://localhost:8081
access_token: (akan di-set otomatis setelah login)
```

### Export to Insomnia

**Method 1: Direct Import**
1. Buka Insomnia
2. Click **Create** → **Import From** → **File**
3. Pilih file `docs/swagger.yaml` (Insomnia prefer YAML)
4. Insomnia akan import sebagai Request Collection
5. Done! Semua endpoints siap digunakan

**Method 2: Via URL**
1. Buka Insomnia
2. Click **Create** → **Import From** → **URL**
3. Paste: `http://localhost:8081/swagger/doc.json`
4. Click **Fetch and Import**

**Setup Environment:**
- Create environment dengan `base_url` dan `access_token`

### Online Converter (Alternative)

**Swagger Editor:**
- URL: https://editor.swagger.io/
- Paste isi `docs/swagger.yaml`
- Click **Generate Client** → **Postman Collection**
- Download dan import ke Postman

**API Transformer:**
- URL: https://www.apimatic.io/transformer
- Upload `docs/swagger.json`
- Convert ke format: Postman v2.1, Insomnia, dll
- Download hasil konversi

### Tips
- File `swagger.json` dan `swagger.yaml` berisi konten yang sama
- Postman support JSON dan YAML
- Insomnia lebih prefer YAML format
- Setelah import, jangan lupa set environment variables
- Authentication akan otomatis ter-handle dengan Bearer token

## 🔑 Development Rules

### ✅ BOLEH (DO)

**Architecture:**
- ✅ Gunakan 5-file module structure (dto, model, repository, route, service)
- ✅ Keep models local per module (duplicate jika perlu)
- ✅ Query database langsung dengan minimal fields yang dibutuhkan
- ✅ Import dari `pkg/` (utilities), `internal/constants/`, `middleware/`
- ✅ Buat helper functions di dalam module sendiri

**Database:**
- ✅ Gunakan raw SQL di repositories
- ✅ Gunakan transactions untuk operasi multiple
- ✅ Implement soft delete dengan `deleted_at`
- ✅ Gunakan prepared statements untuk security
- ✅ Handle `sql.ErrNoRows` dengan proper error message

**Validation:**
- ✅ Gunakan `middleware.ValidateRequest` dengan `ValidationRules`
- ✅ Tambahkan validation tags di DTO structs
- ✅ Validate business logic di service layer
- ✅ Return descriptive error messages

**API Documentation:**
- ✅ Tambahkan Swagger annotations di atas handler methods
- ✅ Dokumentasikan semua parameters, responses, dan errors
- ✅ Gunakan proper HTTP status codes
- ✅ Tambahkan komentar deskriptif untuk setiap route

**Code Quality:**
- ✅ Follow Go naming conventions
- ✅ Write clear, self-documenting code
- ✅ Handle errors properly (don't ignore errors)
- ✅ Use constants untuk messages dan status codes
- ✅ Keep functions small and focused (single responsibility)

**Testing:**
- ✅ Test dengan Swagger UI sebelum commit
- ✅ Test happy path dan error cases
- ✅ Verify dengan `go build` sebelum push

### ❌ TIDAK BOLEH (DON'T)

**Architecture:**
- ❌ Import module lain: `import "gin-scalable-api/internal/modules/user"`
- ❌ Buat shared models di folder global
- ❌ Buat interface + implementation pattern (over-engineering)
- ❌ Buat mapper folder terpisah (konversi inline di service)
- ❌ Cross-module dependencies

**Database:**
- ❌ Gunakan ORM (Gorm, XORM, dll)
- ❌ Hardcode SQL strings di handler atau service
- ❌ Ignore database errors
- ❌ Lakukan query N+1 (gunakan JOIN)
- ❌ Expose password hash di response

**File Structure:**
- ❌ Buat file `handler.go` terpisah (merge ke `route.go`)
- ❌ Buat file `validator.go` terpisah (merge ke `dto.go`)
- ❌ Buat folder `internal/repository/` global
- ❌ Buat folder `internal/models/` global
- ❌ Buat folder `internal/dto/` global

**Validation:**
- ❌ Skip validation middleware
- ❌ Gunakan `ValidateJSON` (deprecated, gunakan `ValidateRequest`)
- ❌ Validate di handler (harus di middleware atau service)
- ❌ Return generic error messages

**Security:**
- ❌ Store passwords in plain text
- ❌ Log sensitive data (passwords, tokens)
- ❌ Skip authentication middleware
- ❌ Expose internal error details ke client
- ❌ Hardcode secrets di code (gunakan .env)

**Code Quality:**
- ❌ Ignore linter warnings
- ❌ Leave commented-out code
- ❌ Use panic() di production code
- ❌ Ignore error returns
- ❌ Write functions > 50 lines (refactor!)

**Git:**
- ❌ Commit `.env` file
- ❌ Commit binary files (`bin/`, `*.exe`)
- ❌ Push broken code (test dulu!)
- ❌ Commit directly ke `main` branch

### ⚠️ PERHATIAN KHUSUS

**Module Independence:**
```go
// ❌ SALAH - Import module lain
import "gin-scalable-api/internal/modules/user"

func (s *Service) GetUserData(userID int64) (*user.User, error) {
    return s.userRepo.GetByID(userID) // Cross-module dependency!
}

// ✅ BENAR - Query langsung atau duplicate model
func (s *Service) GetUserData(userID int64) (*UserBasicInfo, error) {
    query := `SELECT id, name, email FROM users WHERE id = $1`
    var user UserBasicInfo
    err := s.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}
```

**Validation:**
```go
// ❌ SALAH - Validation di handler
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    c.BindJSON(&req)
    
    if req.Name == "" {
        response.Error(c, 400, "Name required", "")
        return
    }
}

// ✅ BENAR - Gunakan middleware
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    users.POST("", 
        middleware.ValidateRequest(middleware.ValidationRules{
            Body: &CreateUserRequest{},
        }),
        handler.CreateUser,
    )
}
```

**Error Handling:**
```go
// ❌ SALAH - Ignore errors
result, _ := s.repo.GetByID(id)

// ❌ SALAH - Expose internal errors
return fmt.Errorf("database error: %v", err)

// ✅ BENAR - Handle dan wrap errors
result, err := s.repo.GetByID(id)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

## 🛠️ Commands

```bash
# Development
make dev-setup      # Setup environment
make dev            # Run with hot reload
make build          # Build binary
make newmodule      # Generate module

# Database
make db-seed-fresh  # Drop, create, seed
make migrate-up     # Run migrations

# Swagger
make swagger-gen    # Generate docs

# Docker
make docker-up      # Start containers
```

## 🧪 Testing

```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity":"800000001","password":"password123"}'

# Use Swagger UI
open http://localhost:8081/swagger/index.html
```

## 🐛 Troubleshooting

```bash
# Check database
pg_isready -h localhost -p 5432

# Check Redis
redis-cli ping

# Fix dependencies
go clean -modcache && go mod download
```

## 📚 Resources

- Swagger Guide: `docs/SWAGGER_GUIDE.md`
- Export Guide: `docs/EXPORT_GUIDE.md` - Postman/Insomnia export
- Integration Guide: `docs/INTEGRATION_GUIDE.md`
- Template: `docs/SWAGGER_ANNOTATION_TEMPLATE.go`
