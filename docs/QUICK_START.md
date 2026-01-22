# Quick Start Guide - RBAC Service (5-File Module Structure)

## 🚀 Setup Development (5 menit)

```bash
# 1. Clone repository
git clone <repo-url>
cd rbac-service

# 2. Setup development environment (install dependencies + create .env)
make dev-setup
# Edit .env: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, REDIS_HOST, JWT_SECRET

# 3. Create database & run migrations
make db-create
make migrate-up

# 4. Build & run server
make run
# atau dengan hot reload (jika sudah install air)
air
```

Server: `http://localhost:8081`

## 🏗️ Architecture Overview

**✅ REFACTORING COMPLETED**: Project telah berhasil direfactor dari 7-file menjadi 5-file per module.

**Benefits:**
- ✅ File count berkurang: 63 → 45 files (28% reduction)
- ✅ Faster navigation: Less file switching untuk developer
- ✅ Cleaner structure: Logical grouping of related code
- ✅ Easier onboarding: New developers less overwhelmed

## 🐳 Setup dengan Docker (Production)

```bash
# 1. Setup environment
cp .env.production .env

# 2. Start services (PostgreSQL, Redis, API)
make prod-start

# 3. Run migrations
make prod-migrate

# 4. Check status
make prod-status

# 5. View logs
make prod-logs
```

Server: `http://localhost:8081`

## 📁 Module Structure (1 fitur = 1 folder, 5 files)

```
internal/modules/feature_name/
├── dto.go          # Request/Response structures + validation logic
├── model.go        # Database entities (local)
├── repository.go   # Database queries (raw SQL)
├── route.go        # Routes + HTTP handlers
└── service.go      # Business logic
```

**Refactoring Changes:**
- ❌ `handler.go` → merged ke `route.go`
- ❌ `validator.go` → merged ke `dto.go`
- ✅ Validation menggunakan `middleware.ValidateRequest` dengan `ValidationRules`
- ✅ Semua routes memiliki dokumentasi komentar yang lengkap

## 🔨 Membuat Module Baru

```bash
# 1. Buat folder
mkdir -p internal/modules/employee

# 2. Buat 5 file
cd internal/modules/employee
touch dto.go model.go repository.go route.go service.go

# 3. Implementasi (copy template dari module lain)

# 4. Register di internal/app/routes.go
```

## 📝 Template Minimal

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

// Validation logic (merged from validator.go)
func ValidateEmployeeName(name string) bool {
    return len(name) >= 2 && len(name) <= 100
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

// Handler struct (merged from handler.go)
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// Handler methods (merged from handler.go)
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

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    employees := router.Group("/employees")
    {
        // POST /api/v1/employees - Create new employee
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

Edit `internal/app/routes.go`:

```go
// Import
employeeModule "gin-scalable-api/internal/modules/employee"

// Di SetupNewModuleRoutes:
employeeModule.RegisterRoutes(protected, h.Employee)
```

Edit `internal/app/server.go`:

```go
// Di NewModuleHandlers:
Employee: employeeModule.NewHandler(
    employeeModule.NewService(
        employeeModule.NewRepository(s.db),
    ),
),
```

## ✅ Testing

```bash
# Build
go build ./cmd/api

# Test endpoint
curl -X POST http://localhost:8081/api/v1/employees \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'
```

## 🚫 Aturan Penting

1. ❌ **TIDAK BOLEH** import module lain: `import "gin-scalable-api/internal/modules/user"`
2. ✅ **BOLEH** duplicate model jika context berbeda
3. ✅ **BOLEH** query database langsung dengan minimal fields
4. ❌ **TIDAK PERLU** interface + implementation pattern
5. ❌ **TIDAK PERLU** mapper terpisah (konversi inline di service)
6. ✅ **GUNAKAN** `middleware.ValidateRequest` dengan `ValidationRules` (bukan `ValidateJSON`)
7. ✅ **TAMBAHKAN** komentar deskriptif untuk setiap route
8. ✅ **MERGE** handler logic ke route.go dan validation logic ke dto.go

## 📚 Dokumentasi Lengkap

- [Module Structure Refactoring](MODULE_STRUCTURE_REFACTORING.md) - Completed refactoring details
- [Backend Engineer Rules](ENGINEER_RULES.md) - Panduan lengkap
- [Project Structure](PROJECT_STRUCTURE.md) - Arsitektur detail
- [API Overview](API_OVERVIEW.md) - API documentation
- [README](../README.md) - Project overview

## 🔧 Makefile Commands

### Development
```bash
make dev-setup      # Setup development environment (deps + .env)
make deps           # Install dependencies
make build          # Build application (bin/server, bin/migrate)
make run            # Build and run server
make test           # Run tests
make fmt            # Format code
make lint           # Run linter
make clean          # Clean build artifacts
```

### Database
```bash
make db-create      # Create database (huminor_rbac)
make db-drop        # Drop database
make db-reset       # Drop, create, and migrate
make migrate-up     # Run migrations
make migrate-status # Check migration status
```

### Production (Docker)
```bash
make prod-build     # Build for production (Linux binary)
make prod-start     # Start production services
make prod-stop      # Stop production services
make prod-logs      # View production logs
make prod-migrate   # Run migrations in production
make prod-backup    # Create database backup
make prod-status    # Check service status
```

### Docker (Development)
```bash
make docker-build   # Build Docker image
make docker-run     # Run with docker-compose
make docker-stop    # Stop Docker services
make docker-logs    # View Docker logs
make docker-clean   # Clean up Docker resources
```

### Help
```bash
make help           # Show all available commands
```

## 🐛 Troubleshooting

**Database connection error:**
```bash
# Check PostgreSQL
pg_isready -h localhost -p 5432

# Check .env
cat .env | grep DB_
```

**Redis connection error:**
```bash
# Check Redis
redis-cli ping

# Check .env
cat .env | grep REDIS_
```

**Build error:**
```bash
go clean -modcache
go mod download
go mod tidy
```
