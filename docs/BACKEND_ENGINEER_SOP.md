# Backend Engineer SOP - Huminor RBAC API

## Panduan Operasional untuk Backend Engineer

### 1. Setup Development Environment

#### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Redis 6+
- Air (untuk live reload)

#### Installation
```bash
# Clone repository
git clone <repository-url>
cd huminor_rbac

# Install dependencies
go mod download

# Install Air untuk live reload
go install github.com/cosmtrek/air@latest

# Setup database
createdb huminor_rbac
```

#### Environment Setup
```bash
# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi database dan Redis
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=huminor_rbac
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key
```

### 2. Database Management

#### Menjalankan Migrasi
```bash
# Jalankan semua migrasi
make migrate-up

# Rollback migrasi terakhir
make migrate-down

# Reset database (drop semua tabel dan jalankan ulang migrasi)
make migrate-reset

# Check status migrasi
make migrate-status
```

#### Membuat Tabel Baru

1. **Buat file migrasi baru:**
```bash
# Format: migrations/XXX_description.sql
# Contoh: migrations/007_create_employees_table.sql
```

2. **Struktur file migrasi:**
```sql
-- migrations/007_create_employees_table.sql
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    position VARCHAR(100),
    department VARCHAR(100),
    hire_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_employee_id ON employees(employee_id);
CREATE INDEX idx_employees_is_active ON employees(is_active);

-- Comments
COMMENT ON TABLE employees IS 'Employee data table';
COMMENT ON COLUMN employees.employee_id IS 'Unique employee identifier';
```

3. **Jalankan migrasi:**
```bash
make migrate-up
```

#### Menyesuaikan Tabel yang Sudah Ada

1. **Buat file migrasi untuk perubahan:**
```sql
-- migrations/008_alter_employees_add_salary.sql
ALTER TABLE employees ADD COLUMN salary DECIMAL(15,2);
ALTER TABLE employees ADD COLUMN currency VARCHAR(3) DEFAULT 'IDR';

-- Update existing records if needed
UPDATE employees SET currency = 'IDR' WHERE currency IS NULL;

-- Add constraints
ALTER TABLE employees ADD CONSTRAINT chk_salary_positive CHECK (salary >= 0);
```

2. **Untuk perubahan yang kompleks, gunakan transaction:**
```sql
-- migrations/009_restructure_user_roles.sql
BEGIN;

-- Create new table
CREATE TABLE user_role_assignments (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    company_id BIGINT NOT NULL REFERENCES companies(id),
    assigned_by BIGINT REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(user_id, role_id, company_id)
);

-- Migrate data from old table
INSERT INTO user_role_assignments (user_id, role_id, company_id, assigned_at)
SELECT user_id, role_id, company_id, created_at 
FROM user_roles 
WHERE is_active = true;

-- Drop old table (optional, bisa di-comment dulu untuk safety)
-- DROP TABLE user_roles;

COMMIT;
```

### 3. Model Development

#### Membuat Model Baru

1. **Buat file model di `internal/models/`:**
```go
// internal/models/employee.go
package models

import (
    "time"
    "gin-scalable-api/pkg/model"
)

type Employee struct {
    model.BaseModel
    CompanyID    int64     `json:"company_id" db:"company_id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    EmployeeID   string    `json:"employee_id" db:"employee_id"`
    Position     *string   `json:"position" db:"position"`
    Department   *string   `json:"department" db:"department"`
    HireDate     *time.Time `json:"hire_date" db:"hire_date"`
    Salary       *float64  `json:"salary" db:"salary"`
    Currency     string    `json:"currency" db:"currency"`
    IsActive     bool      `json:"is_active" db:"is_active"`
}

func (e *Employee) TableName() string {
    return "employees"
}
```

2. **Implementasi Repository:**
```go
// internal/repository/employee_repository.go
package repository

import (
    "database/sql"
    "fmt"
    "gin-scalable-api/internal/models"
    "gin-scalable-api/pkg/model"
)

type EmployeeRepository struct {
    *model.Repository
    db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
    return &EmployeeRepository{
        Repository: model.NewRepository(db),
        db:         db,
    }
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
    query := `
        INSERT INTO employees (company_id, name, email, employee_id, position, department, hire_date, salary, currency, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(query, 
        employee.CompanyID, employee.Name, employee.Email, employee.EmployeeID,
        employee.Position, employee.Department, employee.HireDate, employee.Salary,
        employee.Currency, employee.IsActive,
    ).Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create employee: %w", err)
    }
    
    return nil
}

func (r *EmployeeRepository) GetByID(id int64) (*models.Employee, error) {
    employee := &models.Employee{}
    query := `
        SELECT id, company_id, name, email, employee_id, position, department, 
               hire_date, salary, currency, is_active, created_at, updated_at
        FROM employees 
        WHERE id = $1 AND is_active = true
    `
    
    err := r.db.QueryRow(query, id).Scan(
        &employee.ID, &employee.CompanyID, &employee.Name, &employee.Email,
        &employee.EmployeeID, &employee.Position, &employee.Department,
        &employee.HireDate, &employee.Salary, &employee.Currency,
        &employee.IsActive, &employee.CreatedAt, &employee.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("employee not found")
        }
        return nil, fmt.Errorf("failed to get employee: %w", err)
    }
    
    return employee, nil
}

func (r *EmployeeRepository) Update(employee *models.Employee) error {
    query, values := r.BuildUpdateQuery(employee, employee.ID)
    
    err := r.db.QueryRow(query, values...).Scan(&employee.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to update employee: %w", err)
    }
    
    return nil
}

func (r *EmployeeRepository) Delete(id int64) error {
    query := `UPDATE employees SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete employee: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("employee not found")
    }
    
    return nil
}
```

3. **Implementasi Service:**
```go
// internal/service/employee_service.go
package service

import (
    "gin-scalable-api/internal/models"
    "gin-scalable-api/internal/repository"
)

type EmployeeService struct {
    employeeRepo *repository.EmployeeRepository
}

func NewEmployeeService(employeeRepo *repository.EmployeeRepository) *EmployeeService {
    return &EmployeeService{
        employeeRepo: employeeRepo,
    }
}

type EmployeeResponse struct {
    ID         int64   `json:"id"`
    CompanyID  int64   `json:"company_id"`
    Name       string  `json:"name"`
    Email      string  `json:"email"`
    EmployeeID string  `json:"employee_id"`
    Position   *string `json:"position"`
    Department *string `json:"department"`
    HireDate   *string `json:"hire_date"`
    Salary     *float64 `json:"salary"`
    Currency   string  `json:"currency"`
    IsActive   bool    `json:"is_active"`
    CreatedAt  string  `json:"created_at"`
    UpdatedAt  string  `json:"updated_at"`
}

type CreateEmployeeRequest struct {
    CompanyID  int64   `json:"company_id" binding:"required"`
    Name       string  `json:"name" binding:"required"`
    Email      string  `json:"email" binding:"required,email"`
    EmployeeID string  `json:"employee_id" binding:"required"`
    Position   *string `json:"position"`
    Department *string `json:"department"`
    HireDate   *string `json:"hire_date"`
    Salary     *float64 `json:"salary"`
    Currency   string  `json:"currency"`
}

func (s *EmployeeService) CreateEmployee(req *CreateEmployeeRequest) (*EmployeeResponse, error) {
    employee := &models.Employee{
        CompanyID:  req.CompanyID,
        Name:       req.Name,
        Email:      req.Email,
        EmployeeID: req.EmployeeID,
        Position:   req.Position,
        Department: req.Department,
        Salary:     req.Salary,
        Currency:   req.Currency,
        IsActive:   true,
    }
    
    if req.Currency == "" {
        employee.Currency = "IDR"
    }
    
    if err := s.employeeRepo.Create(employee); err != nil {
        return nil, err
    }
    
    return &EmployeeResponse{
        ID:         employee.ID,
        CompanyID:  employee.CompanyID,
        Name:       employee.Name,
        Email:      employee.Email,
        EmployeeID: employee.EmployeeID,
        Position:   employee.Position,
        Department: employee.Department,
        Salary:     employee.Salary,
        Currency:   employee.Currency,
        IsActive:   employee.IsActive,
        CreatedAt:  employee.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:  employee.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }, nil
}
```

4. **Implementasi Handler:**
```go
// internal/handlers/employee_handler.go
package handlers

import (
    "gin-scalable-api/internal/service"
    "gin-scalable-api/pkg/response"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
    employeeService *service.EmployeeService
}

func NewEmployeeHandler(employeeService *service.EmployeeService) *EmployeeHandler {
    return &EmployeeHandler{
        employeeService: employeeService,
    }
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
    var req service.CreateEmployeeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
        return
    }
    
    result, err := h.employeeService.CreateEmployee(&req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, "Employee created successfully", result)
}
```

5. **Tambahkan ke Routes:**
```go
// internal/routes/routes.go - tambahkan di SetupRoutes function
func SetupRoutes(r *gin.Engine, h *Handlers, jwtSecret string, redis *redis.Client) {
    // ... existing code ...
    
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware(jwtSecret, redis))
    {
        // ... existing routes ...
        
        // Employee routes
        setupEmployeeRoutes(protected, h.Employee)
    }
}

func setupEmployeeRoutes(protected *gin.RouterGroup, employeeHandler *handlers.EmployeeHandler) {
    employees := protected.Group("/employees")
    {
        employees.GET("", employeeHandler.GetEmployees)
        employees.GET("/:id", employeeHandler.GetEmployeeByID)
        employees.POST("", employeeHandler.CreateEmployee)
        employees.PUT("/:id", employeeHandler.UpdateEmployee)
        employees.DELETE("/:id", employeeHandler.DeleteEmployee)
    }
}
```

6. **Update Dependency Injection di Server:**
```go
// internal/server/server.go - tambahkan di Initialize method
func (s *Server) Initialize() {
    // ... existing code ...
    
    // Repositories
    employeeRepo := repository.NewEmployeeRepository(s.DB)
    
    // Services
    employeeService := service.NewEmployeeService(employeeRepo)
    
    // Handlers
    employeeHandler := handlers.NewEmployeeHandler(employeeService)
    
    // Update handlers struct
    handlers := &routes.Handlers{
        // ... existing handlers ...
        Employee: employeeHandler,
    }
}
```

### 4. Development Workflow

#### Menjalankan Development Server

**Opsi 1: Manual Run**
```bash
# Build dan jalankan
go build -o server cmd/api/main.go
./server
```

**Opsi 2: Live Reload dengan Air**
```bash
# Jalankan dengan Air untuk auto-reload
air

# Atau dengan verbose mode
air -v
```

#### Testing

```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run specific test
go test ./internal/service -v

# Run tests dengan race detection
go test -race ./...
```

#### API Testing

```bash
# Test health endpoint
curl -X GET "http://localhost:8081/health"

# Test login dengan user_identity (primary method)
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "user_identity": "100000003",
    "password": "password123"
  }'

# Test login dengan email (alternative method)
curl -X POST "http://localhost:8081/api/v1/auth/login-email" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@company.com",
    "password": "password123"
  }'

# Test protected endpoint dengan token
curl -X GET "http://localhost:8081/api/v1/users?limit=5" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Test create user
curl -X POST "http://localhost:8081/api/v1/users" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New User",
    "email": "newuser@example.com",
    "user_identity": "100000999"
  }'

# Test create company
curl -X POST "http://localhost:8081/api/v1/companies" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PT. New Company",
    "code": "NEWCO"
  }'
```

**Available Test Users:**
- `100000001` - System Admin (full access) | Email: `admin@system.com`
- `100000002` - HR Manager (HR modules) | Email: `hr@company.com`
- `100000003` - Super Admin (full access) | Email: `superadmin@company.com`
- `100000004` - HR Staff (limited HR access) | Email: `hrstaff@company.com`
- `100000005` - Manager (management modules) | Email: `manager@company.com`
- `100000006` - Employee (self-service only) | Email: `employee@company.com`

**Default Password:** `password123`

**Login Methods:**
1. **Primary**: `/api/v1/auth/login` dengan `user_identity`
2. **Alternative**: `/api/v1/auth/login-email` dengan `email`
3. **Fallback**: `/api/v1/auth/login` dengan `email` sebagai `user_identity` (backward compatibility)

### 5. Production Deployment

#### Opsi 1: Binary Deployment
```bash
# Build untuk production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/api/main.go

# Copy binary dan jalankan di server
./server
```

#### Opsi 2: Docker Compose
```bash
# Setup environment production
cp .env.production .env

# Build dan jalankan dengan Docker Compose
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop services
docker-compose -f docker-compose.prod.yml down
```

### 6. Best Practices

#### Code Structure
- Gunakan modular architecture (handler -> service -> repository)
- Pisahkan model berdasarkan domain
- Implementasi proper error handling
- Gunakan centralized response format

#### Database
- Selalu gunakan transaction untuk operasi kompleks
- Buat index untuk kolom yang sering di-query
- Gunakan soft delete (is_active = false) daripada hard delete
- Tambahkan comment pada tabel dan kolom penting

#### Security
- Validasi semua input menggunakan middleware validation
- Implementasi rate limiting
- Gunakan prepared statements untuk mencegah SQL injection
- Hash password dengan bcrypt

#### Performance
- Implementasi pagination untuk list endpoints
- Gunakan connection pooling untuk database
- Cache data yang sering diakses di Redis
- Monitor query performance

### 7. Troubleshooting

#### Common Issues

**Database Connection Error:**
```bash
# Check database status
pg_isready -h localhost -p 5432

# Check connection string di .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=huminor_rbac
```

**Migration Error:**
```bash
# Check migration status
make migrate-status

# Force version (hati-hati!)
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" force VERSION
```

**Route Conflicts:**
- Pastikan tidak ada route yang sama
- Gunakan path yang spesifik untuk menghindari konflik parameter
- Check route registration order

**Build Errors:**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### 8. Monitoring & Logging

#### Health Check
```bash
# Check server health
curl http://localhost:8081/health
```

#### Logs
- Server logs ditampilkan di console
- Gunakan structured logging untuk production
- Monitor error rates dan response times

### 9. Documentation

#### API Documentation
- Update Postman collection setelah menambah endpoint baru
- Dokumentasikan request/response format
- Sertakan contoh payload

#### Code Documentation
- Tambahkan comment untuk function yang kompleks
- Update README.md jika ada perubahan setup
- Dokumentasikan breaking changes di CHANGELOG.md

---

**Catatan Penting:**
- Selalu test di development environment sebelum deploy ke production
- Backup database sebelum menjalankan migration di production
- Monitor server performance setelah deployment
- Koordinasi dengan team untuk perubahan yang mempengaruhi API contract