# 📋 SOP Backend Engineer - ERP RBAC API

## Gambaran Umum

Dokumen ini adalah Standard Operating Procedure (SOP) untuk backend engineer yang bekerja dengan project ERP RBAC API. Project ini menggunakan arsitektur modular dengan raw SQL, PostgreSQL, dan Clean Architecture principles.

---

## 🚀 **Setup Development Environment**

### Prerequisites
- Go 1.25+
- PostgreSQL 12+
- Redis 6+
- Git

### Initial Setup
```bash
# 1. Clone repository
git clone <repository-url>
cd huminor_rbac

# 2. Setup environment
cp .env.example .env
# Edit .env dengan konfigurasi database dan Redis

# 3. Install dependencies
go mod download

# 4. Install Air for live reload (recommended)
go install github.com/cosmtrek/air@latest

# 5. Run migrations
make migrate-up

# 6. Start development server
air  # With live reload (recommended)
# atau
make run  # Manual mode
```

---

## 🌪️ **Development dengan Air (Live Reload)**

### Setup Air
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Verify installation
air -v
```

### Menggunakan Air
```bash
# Start development server dengan live reload
air

# Features yang didapat:
# ✅ Auto rebuild saat file .go berubah
# ✅ Auto restart server setelah build
# ✅ Build error langsung terlihat
# ✅ Fast development cycle
```

### Konfigurasi Air
File `.air.toml` sudah dikonfigurasi dengan optimal:
- **Build Command**: `go build -o ./tmp/main ./cmd/api`
- **Watch Extensions**: `.go`, `.tpl`, `.tmpl`, `.html`
- **Excluded**: `*_test.go`, `tmp/`, `vendor/`
- **Auto Restart**: Ya

### Troubleshooting Air
```bash
# Jika Air tidak ditemukan
export PATH=$PATH:$(go env GOPATH)/bin

# Clean build artifacts
rm -rf tmp/

# Check build errors
cat build-errors.log
```

---

## 🗄️ **Database Management**

### Migration System

Project menggunakan file-based migration system yang terletak di folder `migrations/`.

#### Struktur Migration Files
```
migrations/
├── 001_create_users_table.sql
├── 002_create_companies_and_branches.sql
├── 003_create_roles_and_modules.sql
├── 004_seed_modules_data.sql
├── 005_create_subscription_system.sql
└── 006_seed_initial_data.sql
```

#### Naming Convention
- **Format**: `{number}_{description}.sql`
- **Number**: 3 digit sequential (001, 002, 003, ...)
- **Description**: snake_case, descriptive action

**Contoh:**
- `007_create_employee_table.sql`
- `008_add_department_id_to_users.sql`
- `009_seed_default_departments.sql`

### Membuat Migration Baru

#### 1. Buat File Migration
```bash
# Buat file migration baru dengan nomor urut berikutnya
touch migrations/007_create_employee_table.sql
```

#### 2. Template Migration File
```sql
-- migrations/007_create_employee_table.sql
-- Description: Create employee table with department relationship

-- Create employee table
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    department_id BIGINT REFERENCES departments(id) ON DELETE SET NULL,
    position VARCHAR(100),
    hire_date DATE NOT NULL,
    salary DECIMAL(15,2),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'terminated')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_employees_employee_id ON employees(employee_id);
CREATE INDEX idx_employees_user_id ON employees(user_id);
CREATE INDEX idx_employees_department_id ON employees(department_id);
CREATE INDEX idx_employees_status ON employees(status);

-- Add comments
COMMENT ON TABLE employees IS 'Employee information and employment details';
COMMENT ON COLUMN employees.employee_id IS 'Unique employee identifier';
COMMENT ON COLUMN employees.salary IS 'Employee salary in IDR';
```

#### 3. Run Migration
```bash
# Run single migration
make migrate-up

# Check migration status
make migrate-status

# Rollback if needed (be careful!)
make migrate-down
```

### Modifikasi Tabel Existing

#### 1. Buat Migration untuk Perubahan
```sql
-- migrations/008_add_department_id_to_users.sql
-- Description: Add department_id column to users table

-- Add new column
ALTER TABLE users ADD COLUMN department_id BIGINT REFERENCES departments(id) ON DELETE SET NULL;

-- Create index
CREATE INDEX idx_users_department_id ON users(department_id);

-- Add comment
COMMENT ON COLUMN users.department_id IS 'Reference to user department';
```

#### 2. Migration untuk Data Seeding
```sql
-- migrations/009_seed_default_departments.sql
-- Description: Seed default departments

INSERT INTO departments (name, code, description, is_active) VALUES
('Human Resources', 'HR', 'Human Resources Department', true),
('Information Technology', 'IT', 'Information Technology Department', true),
('Finance', 'FIN', 'Finance Department', true),
('Operations', 'OPS', 'Operations Department', true);
```

---

## 🏗️ **Menambah Fitur Baru (Complete Flow)**

### Langkah 1: Database Schema

#### 1.1 Buat Migration
```bash
touch migrations/010_create_departments_table.sql
```

#### 1.2 Definisi Schema
```sql
-- migrations/010_create_departments_table.sql
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL,
    description TEXT,
    manager_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    parent_id BIGINT REFERENCES departments(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(company_id, code)
);

CREATE INDEX idx_departments_company_id ON departments(company_id);
CREATE INDEX idx_departments_manager_id ON departments(manager_id);
CREATE INDEX idx_departments_parent_id ON departments(parent_id);
```

#### 1.3 Run Migration
```bash
make migrate-up
```

### Langkah 2: Model Definition

#### 2.1 Buat Model File
```bash
touch internal/models/department.go
```

#### 2.2 Definisi Model
```go
// internal/models/department.go
package models

import "time"

type Department struct {
    ID          int64     `json:"id" db:"id"`
    CompanyID   int64     `json:"company_id" db:"company_id"`
    Name        string    `json:"name" db:"name"`
    Code        string    `json:"code" db:"code"`
    Description *string   `json:"description" db:"description"`
    ManagerID   *int64    `json:"manager_id" db:"manager_id"`
    ParentID    *int64    `json:"parent_id" db:"parent_id"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (Department) TableName() string {
    return "departments"
}

// Helper types
type DepartmentWithManager struct {
    Department
    ManagerName *string `json:"manager_name" db:"manager_name"`
}

type DepartmentHierarchy struct {
    Department
    Children []*DepartmentHierarchy `json:"children,omitempty"`
}
```

### Langkah 3: Repository Layer

#### 3.1 Buat Repository
```bash
touch internal/repository/department_repository.go
```

#### 3.2 Implementasi Repository
```go
// internal/repository/department_repository.go
package repository

import (
    "database/sql"
    "gin-scalable-api/internal/models"
)

type DepartmentRepository struct {
    db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
    return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) GetAll(companyID int64, limit, offset int) ([]*models.Department, error) {
    query := `
        SELECT id, company_id, name, code, description, manager_id, parent_id, 
               is_active, created_at, updated_at
        FROM departments 
        WHERE company_id = $1 AND is_active = true
        ORDER BY name
        LIMIT $2 OFFSET $3`
    
    rows, err := r.db.Query(query, companyID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var departments []*models.Department
    for rows.Next() {
        dept := &models.Department{}
        err := rows.Scan(
            &dept.ID, &dept.CompanyID, &dept.Name, &dept.Code,
            &dept.Description, &dept.ManagerID, &dept.ParentID,
            &dept.IsActive, &dept.CreatedAt, &dept.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        departments = append(departments, dept)
    }

    return departments, nil
}

func (r *DepartmentRepository) GetByID(id int64) (*models.Department, error) {
    query := `
        SELECT id, company_id, name, code, description, manager_id, parent_id,
               is_active, created_at, updated_at
        FROM departments WHERE id = $1`
    
    dept := &models.Department{}
    err := r.db.QueryRow(query, id).Scan(
        &dept.ID, &dept.CompanyID, &dept.Name, &dept.Code,
        &dept.Description, &dept.ManagerID, &dept.ParentID,
        &dept.IsActive, &dept.CreatedAt, &dept.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return dept, nil
}

func (r *DepartmentRepository) Create(dept *models.Department) (*models.Department, error) {
    query := `
        INSERT INTO departments (company_id, name, code, description, manager_id, parent_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRow(
        query, dept.CompanyID, dept.Name, dept.Code,
        dept.Description, dept.ManagerID, dept.ParentID,
    ).Scan(&dept.ID, &dept.CreatedAt, &dept.UpdatedAt)
    
    if err != nil {
        return nil, err
    }
    
    return dept, nil
}

// Add Update, Delete, and other methods...
```

### Langkah 4: Service Layer

#### 4.1 Buat Service
```bash
touch internal/service/department_service.go
```

#### 4.2 Implementasi Service
```go
// internal/service/department_service.go
package service

import (
    "gin-scalable-api/internal/models"
    "gin-scalable-api/internal/repository"
)

type DepartmentService struct {
    departmentRepo *repository.DepartmentRepository
}

func NewDepartmentService(departmentRepo *repository.DepartmentRepository) *DepartmentService {
    return &DepartmentService{
        departmentRepo: departmentRepo,
    }
}

// Request/Response types
type DepartmentListRequest struct {
    CompanyID int64 `form:"company_id" binding:"required"`
    Limit     int   `form:"limit"`
    Offset    int   `form:"offset"`
}

type CreateDepartmentRequest struct {
    CompanyID   int64   `json:"company_id" binding:"required"`
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    Code        string  `json:"code" binding:"required,min=2,max=20"`
    Description *string `json:"description"`
    ManagerID   *int64  `json:"manager_id"`
    ParentID    *int64  `json:"parent_id"`
}

type DepartmentResponse struct {
    ID          int64   `json:"id"`
    CompanyID   int64   `json:"company_id"`
    Name        string  `json:"name"`
    Code        string  `json:"code"`
    Description *string `json:"description"`
    ManagerID   *int64  `json:"manager_id"`
    ParentID    *int64  `json:"parent_id"`
    IsActive    bool    `json:"is_active"`
    CreatedAt   string  `json:"created_at"`
    UpdatedAt   string  `json:"updated_at"`
}

func (s *DepartmentService) GetDepartments(req *DepartmentListRequest) ([]*DepartmentResponse, error) {
    if req.Limit == 0 {
        req.Limit = 10
    }

    departments, err := s.departmentRepo.GetAll(req.CompanyID, req.Limit, req.Offset)
    if err != nil {
        return nil, err
    }

    var response []*DepartmentResponse
    for _, dept := range departments {
        response = append(response, &DepartmentResponse{
            ID:          dept.ID,
            CompanyID:   dept.CompanyID,
            Name:        dept.Name,
            Code:        dept.Code,
            Description: dept.Description,
            ManagerID:   dept.ManagerID,
            ParentID:    dept.ParentID,
            IsActive:    dept.IsActive,
            CreatedAt:   dept.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt:   dept.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
        })
    }

    return response, nil
}

func (s *DepartmentService) CreateDepartment(req *CreateDepartmentRequest) (*DepartmentResponse, error) {
    dept := &models.Department{
        CompanyID:   req.CompanyID,
        Name:        req.Name,
        Code:        req.Code,
        Description: req.Description,
        ManagerID:   req.ManagerID,
        ParentID:    req.ParentID,
        IsActive:    true,
    }

    createdDept, err := s.departmentRepo.Create(dept)
    if err != nil {
        return nil, err
    }

    return &DepartmentResponse{
        ID:          createdDept.ID,
        CompanyID:   createdDept.CompanyID,
        Name:        createdDept.Name,
        Code:        createdDept.Code,
        Description: createdDept.Description,
        ManagerID:   createdDept.ManagerID,
        ParentID:    createdDept.ParentID,
        IsActive:    createdDept.IsActive,
        CreatedAt:   createdDept.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:   createdDept.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }, nil
}

// Add other service methods...
```

### Langkah 5: Handler Layer

#### 5.1 Buat Handler
```bash
touch internal/handlers/department_handler.go
```

#### 5.2 Implementasi Handler
```go
// internal/handlers/department_handler.go
package handlers

import (
    "net/http"
    "strconv"

    "gin-scalable-api/internal/service"
    "gin-scalable-api/pkg/response"

    "github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
    departmentService *service.DepartmentService
}

func NewDepartmentHandler(departmentService *service.DepartmentService) *DepartmentHandler {
    return &DepartmentHandler{
        departmentService: departmentService,
    }
}

func (h *DepartmentHandler) GetDepartments(c *gin.Context) {
    var req service.DepartmentListRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request parameters", err.Error())
        return
    }

    result, err := h.departmentService.GetDepartments(&req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to get departments", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Departments retrieved successfully", result)
}

func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
    var req service.CreateDepartmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request format", err.Error())
        return
    }

    result, err := h.departmentService.CreateDepartment(&req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create department", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "Department created successfully", result)
}

// Add other handler methods...
```

### Langkah 6: Routes Configuration

#### 6.1 Update Routes
```go
// internal/routes/routes.go - Add to SetupRoutes function

// Department routes
setupDepartmentRoutes(protected, h.Department)

// Add new function
func setupDepartmentRoutes(protected *gin.RouterGroup, departmentHandler *handlers.DepartmentHandler) {
    departments := protected.Group("/departments")
    {
        // List departments with validation
        listValidation := middleware.ValidationRules{
            Query: []middleware.QueryValidation{
                {Name: "company_id", Type: "int", Required: true, Min: intPtr(1)},
                {Name: "limit", Type: "int", Default: 10, Min: intPtr(1), Max: intPtr(100)},
                {Name: "offset", Type: "int", Default: 0, Min: intPtr(0)},
            },
        }
        departments.GET("", middleware.ValidateRequest(listValidation), departmentHandler.GetDepartments)

        // Create department with validation
        createValidation := middleware.ValidationRules{
            Body: &service.CreateDepartmentRequest{},
        }
        departments.POST("", middleware.ValidateRequest(createValidation), departmentHandler.CreateDepartment)

        // Other routes...
    }
}
```

### Langkah 7: Server Integration

#### 7.1 Update Server
```go
// internal/server/server.go - Update structs

type Repositories struct {
    // ... existing repos
    Department *repository.DepartmentRepository
}

type Services struct {
    // ... existing services
    Department *service.DepartmentService
}

// Update initializeRepositories
func (s *Server) initializeRepositories(db *sql.DB) *Repositories {
    return &Repositories{
        // ... existing repos
        Department: repository.NewDepartmentRepository(db),
    }
}

// Update initializeServices
func (s *Server) initializeServices(repos *Repositories) *Services {
    return &Services{
        // ... existing services
        Department: service.NewDepartmentService(repos.Department),
    }
}

// Update Handlers struct in routes package
type Handlers struct {
    // ... existing handlers
    Department *handlers.DepartmentHandler
}

// Update initializeHandlers
func (s *Server) initializeHandlers(services *Services, repos *Repositories) *routes.Handlers {
    return &routes.Handlers{
        // ... existing handlers
        Department: handlers.NewDepartmentHandler(services.Department),
    }
}
```

### Langkah 8: Testing

#### 8.1 Test dengan Postman
```json
// POST /api/v1/departments
{
    "company_id": 1,
    "name": "Information Technology",
    "code": "IT",
    "description": "IT Department",
    "manager_id": 2,
    "parent_id": null
}
```

#### 8.2 Verify Database
```sql
SELECT * FROM departments WHERE company_id = 1;
```

---

## 🚀 **Production Deployment**

### Option 1: Binary Deployment

#### Build untuk Production
```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/api/main.go

# Set production mode
export GIN_MODE=release
```

#### Environment Variables Production
```bash
# Required environment variables
export GIN_MODE=release
export DB_HOST=production-db-host
export DB_PORT=5432
export DB_USER=db-user
export DB_PASSWORD=secure-password
export DB_NAME=production-db
export REDIS_HOST=redis-host
export REDIS_PORT=6379
export JWT_SECRET=very-secure-jwt-secret
export SERVER_PORT=8081
```

#### Deployment Steps
```bash
# 1. Upload binary ke server
scp server user@server:/opt/huminor-rbac/

# 2. Run migrations
./server migrate

# 3. Start production server
GIN_MODE=release ./server

# 4. Setup systemd service (recommended)
sudo systemctl start huminor-rbac
sudo systemctl enable huminor-rbac
```

### Option 2: Docker Compose (Recommended)

#### Setup Environment
```bash
# 1. Copy production environment template
cp .env.production .env

# 2. Edit environment variables
nano .env
# Update DB_PASSWORD, JWT_SECRET, REDIS_PASSWORD, dll
```

#### Deploy dengan Docker Compose
```bash
# 1. Start all services
./scripts/docker-prod.sh start

# 2. Run migrations
./scripts/docker-prod.sh migrate

# 3. Check service status
./scripts/docker-prod.sh status

# 4. View logs
./scripts/docker-prod.sh logs
```

#### Available Services
- **Application**: Go API server (port 8081)
- **PostgreSQL**: Database dengan persistent storage
- **Redis**: Cache dan session storage
- **Nginx**: Reverse proxy dengan SSL support (optional)

#### Docker Management Commands
```bash
# Start services
./scripts/docker-prod.sh start

# Start with Nginx reverse proxy
./scripts/docker-prod.sh start-nginx

# View logs (all services)
./scripts/docker-prod.sh logs

# View logs (specific service)
./scripts/docker-prod.sh logs app
./scripts/docker-prod.sh logs postgres

# Run database migrations
./scripts/docker-prod.sh migrate

# Create database backup
./scripts/docker-prod.sh backup

# Stop all services
./scripts/docker-prod.sh stop

# Restart services
./scripts/docker-prod.sh restart

# Clean up (remove all containers and volumes)
./scripts/docker-prod.sh clean
```

#### Production Checklist
- [ ] Environment variables configured dengan secure values
- [ ] Database password changed dari default
- [ ] JWT secret menggunakan strong random string
- [ ] Redis password configured (jika diperlukan)
- [ ] SSL certificates configured untuk Nginx (jika menggunakan HTTPS)
- [ ] Firewall rules configured
- [ ] Backup strategy implemented
- [ ] Monitoring dan logging setup

### Docker Deployment Manual
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
CMD ["./server"]
```

---

## 🔧 **Best Practices**

### Database Design
1. **Selalu gunakan foreign key constraints**
2. **Buat index untuk kolom yang sering di-query**
3. **Gunakan appropriate data types**
4. **Tambahkan comments untuk dokumentasi**
5. **Gunakan check constraints untuk validasi**

### Code Structure
1. **Follow Clean Architecture principles**
2. **Gunakan consistent naming conventions**
3. **Implement proper error handling**
4. **Add validation di semua layers**
5. **Use centralized response format**

### Migration Guidelines
1. **Never modify existing migrations**
2. **Always create new migration for changes**
3. **Test migrations on development first**
4. **Backup database before running migrations**
5. **Document breaking changes**

---

## 🚨 **Common Issues & Solutions**

### Migration Issues
```bash
# Issue: Migration failed
# Solution: Check database connection and SQL syntax
make migrate-status
make migrate-down  # If safe to rollback
```

### Build Issues
```bash
# Issue: Import cycle or missing dependencies
go mod tidy
go clean -cache
go build ./...
```

### Database Connection Issues
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Check Redis status
sudo systemctl status redis
```

---

## 📚 **References**

- [Database Schema Documentation](./DATABASE_SCHEMA_DOCUMENTATION.md)
- [Models Documentation](./MODELS_DOCUMENTATION.md)
- [Modular Architecture](./MODULAR_ARCHITECTURE.md)
- [API Infrastructure](./API_INFRASTRUCTURE.md)
- [Migration Guide](./MIGRATIONS.md)

---

## ✅ **Checklist untuk Fitur Baru**

### Database
- [ ] Migration file dibuat dengan naming convention yang benar
- [ ] Schema includes proper constraints dan indexes
- [ ] Migration tested di development environment
- [ ] Database comments added untuk dokumentasi

### Code
- [ ] Model struct dibuat dengan proper tags
- [ ] Repository implements all CRUD operations
- [ ] Service layer includes business logic dan validation
- [ ] Handler uses centralized response format
- [ ] Routes configured dengan proper validation
- [ ] Server integration completed

### Testing
- [ ] API endpoints tested dengan Postman
- [ ] Database queries verified
- [ ] Error handling tested
- [ ] Validation rules tested

### Documentation
- [ ] API endpoints documented
- [ ] Database schema documented
- [ ] Code comments added
- [ ] README updated if needed

---

**🎯 Dengan mengikuti SOP ini, backend engineer dapat bekerja secara konsisten dan efisien dalam mengembangkan fitur baru atau memodifikasi existing features.**