# Backend Engineer SOP - RBAC Service

Panduan lengkap untuk bekerja dengan RBAC Service yang menggunakan Clean Architecture dengan Unit-Based RBAC system.

## 📋 Table of Contents
1. [Quick Start](#quick-start)
2. [Project Structure](#project-structure)
3. [Development Workflow](#development-workflow)
4. [Component Details](#component-details)
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

## 🏗️ Project Structure

```
rbac-service/
├── cmd/api/                 # Entry point aplikasi
├── internal/
│   ├── constants/          # Konstanta aplikasi
│   ├── dto/               # Data Transfer Objects
│   ├── handlers/          # HTTP handlers (controllers)
│   ├── interfaces/        # Interface definitions
│   ├── mapper/           # Data mapping logic
│   ├── models/           # Database models
│   ├── repository/       # Data access layer
│   ├── routes/           # Route definitions
│   ├── server/           # Server setup
│   ├── service/          # Business logic layer
│   └── validation/       # Request validation
├── middleware/           # HTTP middleware
├── migrations/          # Database migrations
├── pkg/                # Shared packages
└── docs/               # Documentation
```

## 🔄 Development Workflow

Ketika mengembangkan fitur baru, ikuti urutan ini:

### 1. Analisis & Desain
- Tentukan endpoint yang dibutuhkan
- Desain struktur data request/response
- Identifikasi business logic yang diperlukan

### 2. Database (jika diperlukan)
- Buat migration file baru
- Definisikan model database

### 3. DTO (Data Transfer Objects)
- Buat request/response DTOs
- Tambahkan validation tags

### 4. Interface
- Definisikan interface untuk repository dan service
- Tambahkan method signatures

### 5. Repository
- Implementasi data access logic
- Handle database operations

### 6. Service
- Implementasi business logic
- Orchestrate repository calls

### 7. Mapper
- Convert antara models dan DTOs
- Handle data transformation

### 8. Handler
- Handle HTTP requests
- Call service methods
- Return responses

### 9. Routes
- Register endpoints
- Apply middleware

### 10. Testing
- Test dengan Postman/curl
- Verify business logic

## 🧩 Component Details

### 📝 DTO (Data Transfer Objects)
**Location**: `internal/dto/`

DTOs mendefinisikan struktur data untuk request dan response API.

**Contoh membuat DTO baru:**

```go
// internal/dto/example_dto.go
package dto

// Request DTO
type CreateExampleRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=100"`
    Description string `json:"description" validate:"required"`
    IsActive    bool   `json:"is_active"`
}

// Response DTO
type ExampleResponse struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    IsActive    bool   `json:"is_active"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

// List Response DTO
type ExampleListResponse struct {
    Data    []*ExampleResponse `json:"data"`
    Total   int64              `json:"total"`
    Limit   int                `json:"limit"`
    Offset  int                `json:"offset"`
    HasMore bool               `json:"has_more"`
}
```

**Validation Tags:**
- `required`: Field wajib diisi
- `email`: Validasi format email
- `min=n,max=n`: Panjang minimum/maksimum
- `omitempty`: Field opsional dalam JSON

### 🔌 Interfaces
**Location**: `internal/interfaces/`

Interfaces mendefinisikan kontrak untuk repository dan service layers.

**Contoh interface:**

```go
// internal/interfaces/repository.go
type ExampleRepositoryInterface interface {
    Create(example *models.Example) error
    GetByID(id int64) (*models.Example, error)
    GetAll(limit, offset int, search string) ([]*models.Example, error)
    Update(example *models.Example) error
    Delete(id int64) error
    Count(search string) (int64, error)
}

// internal/interfaces/service.go
type ExampleServiceInterface interface {
    CreateExample(req *dto.CreateExampleRequest) (*dto.ExampleResponse, error)
    GetExampleByID(id int64) (*dto.ExampleResponse, error)
    GetExamples(req *dto.ExampleListRequest) (*dto.ExampleListResponse, error)
    UpdateExample(id int64, req *dto.UpdateExampleRequest) (*dto.ExampleResponse, error)
    DeleteExample(id int64) error
}
```

### 🗃️ Repository Layer
**Location**: `internal/repository/`

Repository menangani semua operasi database.

**Contoh repository:**

```go
// internal/repository/example_repository.go
package repository

import (
    "database/sql"
    "gin-scalable-api/internal/models"
)

type ExampleRepository struct {
    db *sql.DB
}

func NewExampleRepository(db *sql.DB) *ExampleRepository {
    return &ExampleRepository{db: db}
}

func (r *ExampleRepository) Create(example *models.Example) error {
    query := `
        INSERT INTO examples (name, description, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(query, example.Name, example.Description, example.IsActive).Scan(
        &example.ID, &example.CreatedAt, &example.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("failed to create example: %w", err)
    }
    
    return nil
}

func (r *ExampleRepository) GetByID(id int64) (*models.Example, error) {
    example := &models.Example{}
    query := `
        SELECT id, name, description, is_active, created_at, updated_at
        FROM examples 
        WHERE id = $1 AND deleted_at IS NULL
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &example.ID, &example.Name, &example.Description,
        &example.IsActive, &example.CreatedAt, &example.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("example not found")
        }
        return nil, fmt.Errorf("failed to get example: %w", err)
    }
    
    return example, nil
}
```

**Best Practices Repository:**
- Gunakan prepared statements untuk keamanan
- Handle `sql.ErrNoRows` untuk data tidak ditemukan
- Gunakan soft delete dengan `deleted_at`
- Return error yang descriptive

### 🏢 Service Layer
**Location**: `internal/service/`

Service layer berisi business logic dan orchestrate repository calls.

**Contoh service:**

```go
// internal/service/example_service.go
package service

import (
    "gin-scalable-api/internal/dto"
    "gin-scalable-api/internal/interfaces"
    "gin-scalable-api/internal/mapper"
)

type ExampleService struct {
    exampleRepo   interfaces.ExampleRepositoryInterface
    exampleMapper *mapper.ExampleMapper
}

func NewExampleService(exampleRepo interfaces.ExampleRepositoryInterface) *ExampleService {
    return &ExampleService{
        exampleRepo:   exampleRepo,
        exampleMapper: mapper.NewExampleMapper(),
    }
}

func (s *ExampleService) CreateExample(req *dto.CreateExampleRequest) (*dto.ExampleResponse, error) {
    // Business logic validation
    if req.Name == "" {
        return nil, errors.New("name cannot be empty")
    }
    
    // Convert DTO to model
    example := s.exampleMapper.ToModel(req)
    
    // Save to database
    err := s.exampleRepo.Create(example)
    if err != nil {
        return nil, fmt.Errorf("failed to create example: %w", err)
    }
    
    // Convert model to response DTO
    return s.exampleMapper.ToResponse(example), nil
}

func (s *ExampleService) GetExampleByID(id int64) (*dto.ExampleResponse, error) {
    example, err := s.exampleRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    return s.exampleMapper.ToResponse(example), nil
}
```

**Best Practices Service:**
- Jangan ada struct definitions di service layer
- Semua struct harus di DTO layer
- Focus pada business logic, bukan data access
- Use mapper untuk konversi data

### 🔄 Mapper Layer
**Location**: `internal/mapper/`

Mapper menangani konversi antara models dan DTOs.

**Contoh mapper:**

```go
// internal/mapper/example_mapper.go
package mapper

import (
    "gin-scalable-api/internal/dto"
    "gin-scalable-api/internal/models"
    "time"
)

type ExampleMapper struct{}

func NewExampleMapper() *ExampleMapper {
    return &ExampleMapper{}
}

// Convert DTO to Model
func (m *ExampleMapper) ToModel(req *dto.CreateExampleRequest) *models.Example {
    return &models.Example{
        Name:        req.Name,
        Description: req.Description,
        IsActive:    req.IsActive,
    }
}

// Convert Model to Response DTO
func (m *ExampleMapper) ToResponse(example *models.Example) *dto.ExampleResponse {
    return &dto.ExampleResponse{
        ID:          example.ID,
        Name:        example.Name,
        Description: example.Description,
        IsActive:    example.IsActive,
        CreatedAt:   example.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   example.UpdatedAt.Format(time.RFC3339),
    }
}

// Convert Models to List Response DTO
func (m *ExampleMapper) ToListResponse(examples []*models.Example, total int64, limit, offset int) *dto.ExampleListResponse {
    var responses []*dto.ExampleResponse
    for _, example := range examples {
        responses = append(responses, m.ToResponse(example))
    }
    
    return &dto.ExampleListResponse{
        Data:    responses,
        Total:   total,
        Limit:   limit,
        Offset:  offset,
        HasMore: int64(offset+limit) < total,
    }
}
```

### 🎮 Handler Layer
**Location**: `internal/handlers/`

Handlers menangani HTTP requests dan responses.

**Contoh handler:**

```go
// internal/handlers/example_handler.go
package handlers

import (
    "gin-scalable-api/internal/constants"
    "gin-scalable-api/internal/dto"
    "gin-scalable-api/internal/interfaces"
    "gin-scalable-api/pkg/response"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type ExampleHandler struct {
    exampleService interfaces.ExampleServiceInterface
}

func NewExampleHandler(exampleService interfaces.ExampleServiceInterface) *ExampleHandler {
    return &ExampleHandler{
        exampleService: exampleService,
    }
}

func (h *ExampleHandler) CreateExample(c *gin.Context) {
    // Get validated body from middleware
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }
    
    req, ok := validatedBody.(*dto.CreateExampleRequest)
    if !ok {
        response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
        return
    }
    
    result, err := h.exampleService.CreateExample(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create example", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, constants.MsgExampleCreated, result)
}

func (h *ExampleHandler) GetExampleByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Bad request", "Invalid ID")
        return
    }
    
    result, err := h.exampleService.GetExampleByID(id)
    if err != nil {
        response.Error(c, http.StatusNotFound, constants.MsgExampleNotFound, err.Error())
        return
    }
    
    response.Success(c, http.StatusOK, constants.MsgExampleRetrieved, result)
}
```

**Best Practices Handler:**
- Gunakan validation middleware untuk request validation
- Gunakan `response` package untuk consistent responses
- Handle errors dengan appropriate HTTP status codes
- Gunakan constants untuk messages

### 🛣️ Routes
**Location**: `internal/routes/`

Routes mendefinisikan endpoint dan middleware.

**Contoh routes:**

```go
// internal/routes/example_routes.go
func setupExampleRoutes(router *gin.RouterGroup, exampleHandler *handlers.ExampleHandler) {
    examples := router.Group("/examples")
    {
        examples.GET("", exampleHandler.GetExamples)
        examples.POST("", 
            middleware.ValidateJSON(&dto.CreateExampleRequest{}),
            exampleHandler.CreateExample,
        )
        examples.GET("/:id", exampleHandler.GetExampleByID)
        examples.PUT("/:id",
            middleware.ValidateJSON(&dto.UpdateExampleRequest{}),
            exampleHandler.UpdateExample,
        )
        examples.DELETE("/:id", exampleHandler.DeleteExample)
    }
}
```

### ✅ Validation
**Location**: `internal/validation/`

Validation rules untuk request DTOs.

**Contoh validation:**

```go
// internal/validation/example_validation.go
package validation

import (
    "gin-scalable-api/internal/dto"
    "github.com/go-playground/validator/v10"
)

func ValidateCreateExampleRequest(req *dto.CreateExampleRequest) error {
    validate := validator.New()
    
    // Custom validation rules
    validate.RegisterValidation("example_name", validateExampleName)
    
    return validate.Struct(req)
}

func validateExampleName(fl validator.FieldLevel) bool {
    name := fl.Field().String()
    // Custom validation logic
    return len(name) >= 2 && len(name) <= 100
}
```

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