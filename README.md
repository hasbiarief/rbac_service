# RBAC API - Clean Architecture

Sistem ERP dengan Role-Based Access Control (RBAC) menggunakan **Clean Architecture**, Go, PostgreSQL, dan Redis. Menggunakan raw SQL tanpa ORM untuk performa optimal.

## 🏗️ Arsitektur

Project ini menggunakan **Clean Architecture** dengan struktur:

```
internal/
├── dto/           # Data Transfer Objects (Request/Response)
├── interfaces/    # Contracts untuk Service & Repository  
├── mapper/        # Konversi antara Model dan DTO
├── handlers/      # HTTP Request/Response handlers
├── service/       # Business logic layer
├── repository/    # Data access layer
├── models/        # Database entities
├── validation/    # Request validation rules
├── routes/        # Route definitions
└── constants/     # Konstanta aplikasi

pkg/
├── errors/        # Custom error types
├── query/         # Query builder utilities
├── pagination/    # Pagination helpers
├── utils/         # Response & validation helpers
└── response/      # Standardized API responses
```

## ✨ Fitur Utama

- **Clean Architecture** - Separation of concerns yang jelas
- **Authentication** - JWT token system dengan Redis storage
- **RBAC System** - Role-based access control dengan module permissions
- **Multi-Company** - Company dan branch hierarchy (4 levels)
- **Subscription** - Tiered subscription dengan module access control
- **Audit Logging** - Complete audit trail
- **User Management** - User lifecycle dengan soft delete
- **DTO Pattern** - Consistent request/response structure
- **Custom Error Handling** - Structured error responses
- **Query Builder** - Dynamic SQL query building
- **Validation Middleware** - Centralized input validation

## 🛠️ Tech Stack

- **Go 1.25+** dengan Gin framework
- **PostgreSQL** dengan raw SQL (tanpa ORM)
- **Redis** untuk token storage dan caching
- **Clean Architecture** dengan DTO, Mapper, Interface patterns
- **Custom Error Types** untuk error handling
- **Validation Middleware** untuk input validation
- **bcrypt** untuk password hashing

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

### Membuat Fitur Baru

1. **Buat DTO** di `internal/dto/`
2. **Buat Interface** di `internal/interfaces/`
3. **Buat Mapper** di `internal/mapper/`
4. **Implementasi Repository** di `internal/repository/`
5. **Implementasi Service** di `internal/service/`
6. **Implementasi Handler** di `internal/handlers/`
7. **Tambah Validation** di `internal/validation/`
8. **Tambah Routes** di `internal/routes/`

### Contoh Implementasi

**DTO:**
```go
// internal/dto/employee_dto.go
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

**Service Interface:**
```go
// internal/interfaces/service.go
type EmployeeServiceInterface interface {
    CreateEmployee(req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
    GetEmployeeByID(id int64) (*dto.EmployeeResponse, error)
}
```

**Handler:**
```go
// internal/handlers/employee_handler.go
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
    var req *dto.CreateEmployeeRequest
    if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    result, err := h.employeeService.CreateEmployee(req)
    if err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    h.responseHelper.Created(c, constants.MsgEmployeeCreated, result)
}
```

## 📖 Documentation

- **[Documentation Index](docs/INDEX.md)** - Daftar dan navigasi dokumentasi
- **[Backend SOP](docs/BACKEND_ENGINEER_SOP.md)** - Panduan lengkap pengembangan
- **[Migration Guide](docs/MIGRATIONS.md)** - Prosedur database migration
- **[Project Structure](docs/PROJECT_STRUCTURE.md)** - Overview arsitektur
- **[Role Permissions](docs/ROLE_PERMISSIONS_MAPPING.md)** - Dokumentasi RBAC
- **[Postman Collection](docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json)** - API testing

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