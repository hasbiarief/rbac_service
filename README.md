# RBAC API - Module-Based Architecture (5-File Structure)

> **Made with ❤️ by Hasbi**  
> 📧 hasbiarief7@gmail.com

Sistem ERP dengan Role-Based Access Control (RBAC) menggunakan **Module-Based Architecture** (Express.js style), Go, PostgreSQL, dan Redis. Menggunakan raw SQL tanpa ORM untuk performa optimal.

**✅ REFACTORING COMPLETED**: Struktur telah berhasil direfactor dari 7-file menjadi 5-file per module untuk meningkatkan developer experience dan mengurangi cognitive load.

## 🏗️ Arsitektur

Project ini menggunakan **vertical module-based structure** dengan prinsip: **1 fitur = 1 folder**

```
internal/
├── app/           # Server initialization & route registration
├── modules/       # 🔥 SEMUA FITUR DI SINI (5 files per module)
│   ├── auth/      # Authentication (dto, model, repository, route, service)
│   ├── user/      # User management
│   ├── role/      # Role management
│   ├── company/   # Company management
│   ├── branch/    # Branch management (hierarchical)
│   ├── module/    # Module system (menu/features)
│   ├── unit/      # Unit management (unit-based RBAC)
│   ├── subscription/  # Subscription system
│   └── audit/     # Audit logging
└── constants/     # Shared constants

middleware/        # HTTP middleware (auth, CORS, rate limit)
pkg/              # Reusable utilities (generic only)
migrations/       # SQL migrations
```

**Refactoring Benefits:**
- ✅ File count berkurang: 63 → 45 files (28% reduction)
- ✅ Faster navigation: Less file switching untuk developer
- ✅ Cleaner structure: Logical grouping of related code
- ✅ Easier onboarding: New developers less overwhelmed
- ✅ Maintained modularity: Zero impact ke cross-module dependencies

**Key Differences:**
- ❌ Tidak ada folder `interfaces/`, `mapper/`, `dto/` global
- ❌ Tidak ada separation `handlers/`, `service/`, `repository/` terpisah
- ❌ Tidak ada `handler.go` dan `validator.go` terpisah (merged ke route.go dan dto.go)
- ✅ Setiap module punya semua layer-nya sendiri (5 files: dto, model, repository, route, service)
- ✅ Model lokal per module (tidak shared)
- ✅ No cross-module imports

## ✨ Fitur Utama

- **Module-Based Architecture** - 1 fitur = 1 folder (5-file structure)
- **Module Generator** - Otomatis generate module dengan Makefile commands
- **Authentication** - JWT token system dengan Redis storage
- **Unit-Based RBAC** - Hierarchical RBAC: Company → Branch → Unit → Role → User
- **Multi-Company** - Company dan branch hierarchy dengan unit management
- **Subscription** - Tiered subscription dengan module access control
- **Audit Logging** - Complete audit trail untuk semua actions
- **User Management** - User lifecycle dengan soft delete
- **Module System** - Dynamic menu/feature management
- **Raw SQL** - Tanpa ORM untuk performa optimal
- **Validation Middleware** - Centralized input validation dengan `ValidateRequest`
- **CORS & Rate Limiting** - Production-ready security

## 🛠️ Tech Stack

- **Go 1.21+** dengan Gin framework
- **PostgreSQL 13+** dengan raw SQL (tanpa ORM)
- **Redis 6+** untuk token storage dan caching
- **Module-Based Architecture** (Express.js style)
- **JWT Authentication** dengan refresh token
- **Validation Middleware** dengan go-playground/validator
- **bcrypt** untuk password hashing
- **Air** untuk hot reload development

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- PostgreSQL 12+
- Redis 6+

### Development Setup

```bash
# 1. Clone dan setup
git clone <repository>
cd rbac-service
make dev-setup  # Install dependencies + create .env

# 2. Edit environment variables
# Edit .env dengan konfigurasi database dan Redis

# 3. Setup database
make db-create
make db-seed    # Seed with template data (recommended)
# atau
make migrate-up # Run migrations only

# 4. Start development server
air  # dengan live reload
# atau
make run
```

Server akan berjalan di `http://localhost:8081`

### Quick Database Setup (New Projects)

```bash
# One-command setup for new projects
make db-seed-fresh  # Drop, create, and seed database

# Or step by step
make db-create      # Create database
make db-seed        # Seed with template data
```

**Template includes:**
- ✅ 11 sample users with different roles
- ✅ 3 companies with hierarchical branches
- ✅ 128+ modules across 12 categories
- ✅ Complete RBAC with unit-based permissions
- ✅ 3 subscription plans (Basic, Pro, Enterprise)
- ✅ Ready-to-use data for testing

### Module Development

```bash
# Generate module baru dengan boilerplate code
make newmodule name=employee

# List semua modules
make listmodules

# Remove module jika diperlukan
make removemodule name=employee

# Build dan test
make build
go build ./cmd/api
```

### Production Deployment

```bash
# Build binary
make prod-build

# Run migrations
./bin/migrate -action=up -dir=migrations

# Start server
GIN_MODE=release ./bin/server
```

**Docker Compose (Recommended):**
```bash
cp .env.production .env
make prod-start
make prod-migrate
make prod-status
```

## 📚 API Documentation

### Test Users
Default users untuk testing (password: `password123`):
- `admin@system.com` - System Admin (akses penuh)
- `hr@company.com` - HR Manager (modul HR)
- `superadmin@company.com` - Super Admin (akses penuh)

### Key Endpoints

**Authentication:**
```bash
# Login
POST /api/v1/auth/login
{
  "email": "admin@system.com",
  "password": "password123"
}

# Refresh token
POST /api/v1/auth/refresh
{
  "refresh_token": "your_refresh_token"
}

# Logout
POST /api/v1/auth/logout
```

**User Management:**
```bash
# List users dengan pagination
GET /api/v1/users?limit=10&offset=0&search=admin

# Create user
POST /api/v1/users
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}

# Update user
PUT /api/v1/users/1
{
  "name": "John Updated",
  "is_active": true
}
```

**Company & Branch:**
```bash
# List companies
GET /api/v1/companies

# Company branches (hierarchy)
GET /api/v1/branches/company/1?nested=true

# Create branch
POST /api/v1/branches
{
  "company_id": 1,
  "name": "Jakarta Branch",
  "code": "JKT",
  "parent_id": null
}
```

**Subscription:**
```bash
# List plans (public)
GET /api/v1/plans

# Create subscription
POST /api/v1/subscription/subscriptions
{
  "company_id": 1,
  "plan_id": 1,
  "billing_cycle": "monthly"
}

# Check module access
GET /api/v1/subscription/module-access/1/1
```

### Response Format

**Success Response:**
```json
{
  "success": true,
  "message": "Pengguna berhasil dibuat",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Validasi gagal",
  "error": "Email sudah ada"
}
```

**Paginated Response:**
```json
{
  "success": true,
  "message": "Pengguna berhasil diambil",
  "data": {
    "data": [...],
    "total": 100,
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}
```

## 🏗️ Development Guide

### Membuat Module Baru

#### ⚡ Otomatis dengan Makefile (Recommended)
```bash
# 1. Generate module dengan boilerplate code
make newmodule name=employee

# 2. Register module di internal/app/routes.go dan internal/app/server.go

# 3. Test build
go build ./cmd/api

# 4. Implement business logic sesuai kebutuhan
```

#### 🔧 Manual (jika diperlukan)
1. **Buat folder module** di `internal/modules/feature_name/`
2. **Buat 5 file standar:**
   - `dto.go` - Request/Response structures + validation logic
   - `model.go` - Database entities (local)
   - `repository.go` - Database queries (raw SQL)
   - `route.go` - Route registration + HTTP handlers
   - `service.go` - Business logic
3. **Register module** di `internal/app/routes.go`

### Module Structure (5 Files)

Setiap module memiliki 5 file standar setelah refactoring:

**Model:**
```go
// internal/modules/employee/model.go
package employee

import "time"

type employee struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (employee) TableName() string {
    return "employees"
}
```

**DTO:**
```go
// internal/modules/employee/dto.go
package employee

type createemployeerequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

type employeeresponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

// Validation logic (merged from validator.go)
func ValidateEmployeeEmail(email string) bool {
    // Custom validation logic
    return len(email) > 0
}
```

**Repository:**
```go
// internal/modules/employee/repository.go
package employee

import (
    "database/sql"
    "gin-scalable-api/pkg/model"
)

type repository struct {
    *model.Repository
    db *sql.DB
}

func newrepository(db *sql.DB) *repository {
    return &repository{
        Repository: model.NewRepository(db),
        db:         db,
    }
}

func (r *repository) create(emp *employee) error {
    query := `INSERT INTO employees (name, email, is_active, created_at, updated_at)
              VALUES ($1, $2, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              RETURNING id, created_at, updated_at`
    return r.db.QueryRow(query, emp.Name, emp.Email).Scan(&emp.ID, &emp.CreatedAt, &emp.UpdatedAt)
}
```

**Service:**
```go
// internal/modules/employee/service.go
package employee

import "time"

type service struct {
    repo *repository
}

func newservice(repo *repository) *service {
    return &service{repo: repo}
}

func (s *service) createemployee(req *createemployeerequest) (*employeeresponse, error) {
    emp := &employee{Name: req.Name, Email: req.Email}
    if err := s.repo.create(emp); err != nil {
        return nil, err
    }
    return &employeeresponse{
        ID:        emp.ID,
        Name:      emp.Name,
        Email:     emp.Email,
        IsActive:  emp.IsActive,
        CreatedAt: emp.CreatedAt.Format(time.RFC3339),
        UpdatedAt: emp.UpdatedAt.Format(time.RFC3339),
    }, nil
}
```

**Route (with Handler):**
```go
// internal/modules/employee/route.go
package employee

import (
    "gin-scalable-api/middleware"
    "gin-scalable-api/pkg/response"
    "net/http"
    "github.com/gin-gonic/gin"
)

// Handler struct (merged from handler.go)
type handler struct {
    service *service
}

func newhandler(service *service) *handler {
    return &handler{service: service}
}

// Handler methods (merged from handler.go)
func (h *handler) createemployee(c *gin.Context) {
    validatedBody, _ := c.Get("validated_body")
    req := validatedBody.(*createemployeerequest)
    
    result, err := h.service.createemployee(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create employee", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, "Employee created", result)
}

// Route registration
func RegisterRoutes(router *gin.RouterGroup, handler *handler) {
    employees := router.Group("/employees")
    {
        // POST /api/v1/employees - Create new employee
        employees.POST("", 
            middleware.ValidateRequest(middleware.ValidationRules{
                Body: &createemployeerequest{},
            }),
            handler.createemployee,
        )
        
        // GET /api/v1/employees/:id - Get employee by ID
        employees.GET("/:id", handler.getemployeebyid)
    }
}
```

## 📖 Documentation

- **[📚 Documentation Index](docs/INDEX.md)** - Navigasi dokumentasi lengkap
- **[🚀 Quick Start Guide](docs/QUICK_START.md)** - Setup cepat & Makefile commands
- **[🏗️ Project Structure](docs/PROJECT_STRUCTURE.md)** - Module-based architecture (5-file structure)
- **[👨‍💻 Backend Engineer Rules](docs/ENGINEER_RULES.md)** - Development guide
- **[📡 API Overview](docs/API_OVERVIEW.md)** - API documentation
- **[🔄 Module Structure Refactoring](docs/MODULE_STRUCTURE_REFACTORING.md)** - Refactoring details
- **[🗄️ Database Documentation](database/README.md)** - Database dumps and seeders
- **[🧪 Postman Collection](docs/HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json)** - API testing

## ⚙️ Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=huminor_rbac

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key

# Server
PORT=8081
GIN_MODE=debug

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
ENVIRONMENT=development
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run specific test
go test ./internal/modules/user -v

# Test dengan race detection
go test -race ./...

# Build validation
go build ./cmd/api
```

### Makefile Commands

```bash
# Development
make dev-setup      # Setup development environment
make build          # Build application
make run            # Build and run server
make test           # Run tests
make clean          # Clean build artifacts

# Module Management
make newmodule name=<name>     # Generate new module
make removemodule name=<name>  # Remove existing module
make listmodules               # List all modules

# Database
make db-create      # Create database
make db-drop        # Drop database
make db-reset       # Reset database and run migrations
make db-dump        # Create database dump and seeder files
make db-seed        # Seed database with template data
make db-seed-fresh  # Drop, create, and seed database
make migrate-up     # Run migrations
make migrate-status # Check migration status

# Production
make prod-build     # Build for production
make prod-start     # Start production services
make prod-stop      # Stop production services
make prod-logs      # View production logs

# Help
make help           # Show all available commands
```

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8081/health
```

### Metrics
- Response time monitoring
- Error rate tracking
- Database connection pooling
- Redis cache hit rates

## 🔧 Troubleshooting

### Common Issues

**Database Connection:**
```bash
# Check database
pg_isready -h localhost -p 5432

# Check environment
cat .env | grep DB_
```

**Migration Error:**
```bash
# Check status
make migrate-status

# Force version (hati-hati!)
migrate -path migrations -database "postgres://..." force VERSION
```

**Build Error:**
```bash
# Clean cache
go clean -modcache
go mod download
go mod tidy
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Follow clean architecture patterns
4. Add tests untuk fitur baru
5. Update documentation
6. Commit changes (`git commit -m 'Add amazing feature'`)
7. Push branch (`git push origin feature/amazing-feature`)
8. Open Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with Module-Based Architecture (5-File Structure) + Go + PostgreSQL + Redis**

**Key Features:**
- ✅ 28% file reduction (63 → 45 files)
- ✅ Automated module generation with Makefile
- ✅ Consistent 5-file structure per module
- ✅ Production-ready with Docker support