# RBAC Service - Huminor

Role-Based Access Control (RBAC) service dengan module-based architecture untuk sistem manajemen akses dan permission.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Make

### Setup Development

```bash
# Clone repository
git clone <repository-url>
cd rbac-service

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Setup database
make db-setup

# Run migrations
make migrate-up

# Run seeders
make seed

# Start development server
make dev
```

Server akan berjalan di `http://localhost:8081`

## 📖 Dokumentasi API

### Swagger UI (Recommended)
Akses dokumentasi interaktif dan test endpoints:
```
http://localhost:8081/swagger/index.html
```

### Import ke Postman/Insomnia
1. Download file: `docs/swagger.json` atau `docs/swagger.yaml`
2. Import ke Postman atau Insomnia
3. Set environment variables
4. Ready to test!

**Via URL (jika server running):**
```
http://localhost:8081/swagger/doc.json
```

### Generate Swagger Documentation
```bash
make swagger-gen
```

## 🏗️ Arsitektur

### Module-Based Architecture
Setiap fitur adalah satu module dengan struktur:
```
internal/modules/<module-name>/
├── route.go       # HTTP routes & handlers
├── service.go     # Business logic
├── repository.go  # Database operations
├── model.go       # Database models
├── dto.go         # Request/Response DTOs
```

### Modules
- **auth** - Authentication & authorization
- **user** - User management
- **company** - Company management
- **branch** - Branch management (hierarchical)
- **role** - Role management
- **module** - Module management
- **subscription** - Subscription & plan management
- **unit** - Unit management
- **application** - Application management
- **audit** - Audit logging

## 🛠️ Development

### Makefile Commands

```bash
# Development
make dev              # Run with hot reload
make build            # Build binary
make run              # Run binary

# Database
make db-setup         # Create database
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make seed             # Run seeders

# Swagger
make swagger-gen      # Generate Swagger docs
make swagger-validate # Validate annotations

# Testing
make test             # Run tests
make test-coverage    # Run with coverage

# Docker
make docker-build     # Build Docker image
make docker-up        # Start containers
make docker-down      # Stop containers
```

### Adding New Module

1. Create module directory:
```bash
mkdir -p internal/modules/<module-name>
```

2. Create required files (route.go, service.go, repository.go, model.go, dto.go)

3. Add Swagger annotations directly above handler methods

4. Register routes in `internal/app/routes.go`

5. Generate Swagger docs:
```bash
make swagger-gen
```

## 📝 Swagger Annotations

Tambahkan anotasi langsung di atas handler method:

```go
// @Summary      Get user by ID
// @Description  Mendapatkan detail user berdasarkan ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=user.UserResponse}
// @Failure      404  {object}  response.Response
// @Router       /api/v1/users/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
    // handler implementation
}
```

Template lengkap: `docs/SWAGGER_ANNOTATION_TEMPLATE.go`

## 🔑 Authentication

Service menggunakan JWT dengan refresh token:

1. **Login**: `POST /api/v1/auth/login`
2. **Refresh**: `POST /api/v1/auth/refresh`
3. **Logout**: `POST /api/v1/auth/logout`

Gunakan Bearer token di header:
```
Authorization: Bearer <access_token>
```

## 🗄️ Database

### Migrations
Migrations menggunakan `golang-migrate`:
```bash
# Create new migration
migrate create -ext sql -dir database/migrations -seq <migration_name>

# Run migrations
make migrate-up

# Rollback
make migrate-down
```

### Seeders
Seeder files di `database/seeders/`:
```bash
make seed
```

## 🐳 Docker

### Development
```bash
docker-compose up -d
```

### Production
```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 📊 API Documentation System

Service memiliki built-in API documentation system yang mendukung:
- Collection management
- Folder organization
- Endpoint documentation
- Environment variables
- Export ke Postman, OpenAPI, Insomnia, Swagger, Apidog

Dokumentasi lengkap: `docs/apidoc/`

## 🔒 Security

- JWT authentication dengan refresh token
- Password hashing menggunakan bcrypt
- Role-based access control (RBAC)
- Module-based permissions
- Subscription-based feature access
- Audit logging untuk semua operasi

## 📚 Dokumentasi

- **[Developer Guide](docs/DEVELOPER_GUIDE.md)** - Setup, architecture, dan development workflow
- **[Swagger Guide](docs/SWAGGER_GUIDE.md)** - API documentation dengan Swagger
- **[Integration Guide](docs/INTEGRATION_GUIDE.md)** - Integrasi dengan external apps

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test ./internal/modules/auth/...
```

## 📦 Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL (raw SQL)
- **Cache**: Redis
- **Auth**: JWT
- **Documentation**: Swagger/OpenAPI
- **Migration**: golang-migrate
- **Hot Reload**: Air

## 🤝 Contributing

1. Create feature branch
2. Follow module-based architecture
3. Add Swagger annotations
4. Write tests
5. Update documentation
6. Submit pull request

## 📄 License

[Your License]

## 👥 Team

[Your Team Info]
