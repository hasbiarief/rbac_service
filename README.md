# RBAC API - Module-Based Architecture

Sistem ERP dengan Role-Based Access Control (RBAC) menggunakan **Module-Based Architecture** (Express.js style), Go, PostgreSQL, dan Redis. Menggunakan raw SQL tanpa ORM untuk performa optimal.

## 🏗️ Arsitektur

Project ini menggunakan **vertical module-based structure** dengan prinsip: **1 fitur = 1 folder**

```
internal/
├── app/           # Server initialization & route registration
├── modules/       # 🔥 SEMUA FITUR DI SINI
│   ├── auth/      # Authentication (7 files: route, handler, service, repository, model, dto, validator)
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

**Key Differences:**
- ❌ Tidak ada folder `interfaces/`, `mapper/`, `dto/` global
- ❌ Tidak ada separation `handlers/`, `service/`, `repository/` terpisah
- ✅ Setiap module punya semua layer-nya sendiri (route, handler, service, repository, model, dto, validator)
- ✅ Model lokal per module (tidak shared)
- ✅ No cross-module imports

## ✨ Fitur Utama

- **Module-Based Architecture** - 1 fitur = 1 folder (Express.js style)
- **Authentication** - JWT token system dengan Redis storage
- **Unit-Based RBAC** - Hierarchical RBAC: Company → Branch → Unit → Role → User
- **Multi-Company** - Company dan branch hierarchy dengan unit management
- **Subscription** - Tiered subscription dengan module access control
- **Audit Logging** - Complete audit trail untuk semua actions
- **User Management** - User lifecycle dengan soft delete
- **Module System** - Dynamic menu/feature management
- **Raw SQL** - Tanpa ORM untuk performa optimal
- **Validation Middleware** - Centralized input validation
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
go mod download

# 2. Setup environment
cp .env.example .env
# Edit .env dengan konfigurasi database dan Redis

# 3. Setup database
createdb huminor_rbac
make migrate-up

# 4. Start development server
air  # dengan live reload
# atau
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8081`

### Production Deployment

```bash
# Build binary
CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

# Run migrations
./server migrate

# Start server
GIN_MODE=release ./server
```

**Docker Compose (Recommended):**
```bash
cp .env.production .env
docker-compose -f docker-compose.prod.yml up -d
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

1. **Buat folder module** di `internal/modules/feature_name/`
2. **Buat 7 file standar:**
   - `route.go` - Route registration
   - `handler.go` - HTTP handlers
   - `service.go` - Business logic
   - `repository.go` - Database queries (raw SQL)
   - `model.go` - Database entities (local)
   - `dto.go` - Request/Response structures
   - `validator.go` - Custom validation rules
3. **Register module** di `internal/app/routes.go`

### Contoh Implementasi

**Model:**
```go
// internal/modules/employee/model.go
package employee

type Employee struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

**DTO:**
```go
// internal/modules/employee/dto.go
package employee

type CreateEmployeeRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

type EmployeeResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}
```

**Repository:**
```go
// internal/modules/employee/repository.go
package employee

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(emp *Employee) error {
    query := `INSERT INTO employees (name, email, is_active, created_at, updated_at)
              VALUES ($1, $2, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
              RETURNING id, created_at`
    return r.db.QueryRow(query, emp.Name, emp.Email).Scan(&emp.ID, &emp.CreatedAt)
}
```

**Service:**
```go
// internal/modules/employee/service.go
package employee

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateEmployee(req *CreateEmployeeRequest) (*EmployeeResponse, error) {
    emp := &Employee{Name: req.Name, Email: req.Email}
    if err := s.repo.Create(emp); err != nil {
        return nil, err
    }
    return &EmployeeResponse{
        ID:        emp.ID,
        Name:      emp.Name,
        Email:     emp.Email,
        CreatedAt: emp.CreatedAt.Format(time.RFC3339),
    }, nil
}
```

**Handler:**
```go
// internal/modules/employee/handler.go
package employee

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

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
```

**Route:**
```go
// internal/modules/employee/route.go
package employee

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
    employees := router.Group("/employees")
    {
        employees.POST("", middleware.ValidateJSON(&CreateEmployeeRequest{}), handler.CreateEmployee)
        employees.GET("/:id", handler.GetEmployee)
    }
}
```

## 📖 Documentation

- **[📚 Documentation Index](docs/INDEX.md)** - Navigasi dokumentasi lengkap
- **[🚀 Quick Start Guide](docs/QUICK_START.md)** - Setup cepat & Makefile commands
- **[🏗️ Project Structure](docs/PROJECT_STRUCTURE.md)** - Module-based architecture
- **[👨‍💻 Backend Engineer SOP](docs/BACKEND_ENGINEER_SOP.md)** - Development guide
- **[📡 API Overview](docs/API_OVERVIEW.md)** - API documentation
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
go test ./internal/service -v

# Test dengan race detection
go test -race ./...
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

**Built with Clean Architecture + Go + PostgreSQL + Redis**