# ERP RBAC API - Modular Architecture

Sistem ERP dengan Role-Based Access Control (RBAC) yang dibangun menggunakan Go, Gin, PostgreSQL, dan Redis. Project ini menggunakan arsitektur modular dengan raw SQL (tanpa ORM) untuk performa optimal.

## 🚀 Fitur Utama

- **Authentication System** - JWT-based authentication dengan refresh token
- **Role-Based Access Control** - Sistem RBAC yang fleksibel dengan module-based permissions
- **Module System** - Hierarchical module system untuk granular access control
- **Company & Branch Management** - Multi-company dengan hierarchical branch structure
- **Subscription System** - Tiered subscription dengan module access control
- **Audit Logging** - Comprehensive audit trail untuk semua aktivitas
- **User Management** - Complete user lifecycle management
- **Password Management** - Secure password handling dengan bcrypt

## 🏗️ Arsitektur

Project menggunakan **Clean Architecture** dengan struktur modular:

```
├── cmd/                    # Application entry points
│   ├── api/               # Main API server
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── handlers/          # HTTP request handlers (8 modules)
│   ├── routes/            # Route configuration
│   ├── server/            # Server initialization
│   ├── service/           # Business logic layer (8 services)
│   ├── repository/        # Data access layer (7 repositories)
│   └── models/            # Data models (7 domains)
├── middleware/            # HTTP middlewares
├── pkg/                   # Reusable packages
│   ├── database/          # Database connection
│   ├── migration/         # Migration system
│   └── token/             # JWT token management
├── migrations/            # SQL migration files
├── config/                # Configuration management
└── docs/                  # Documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Authentication**: JWT
- **Password**: bcrypt
- **Migration**: Custom file-based system
- **Architecture**: Clean Architecture + Modular Design

## 📦 Models

Sistem memiliki 7 domain model utama:

- **User** - User management dengan authentication
- **Company** - Multi-company support
- **Branch** - Hierarchical branch structure
- **Role** - Role dan permission management
- **Module** - Hierarchical module system
- **Subscription** - Tiered subscription system
- **Audit** - Comprehensive audit logging

## 🚦 Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 12+
- Redis 6+

### Installation

1. **Clone repository**
```bash
git clone <repository-url>
cd huminor_rbac
```

2. **Setup environment**
```bash
cp .env.example .env
# Edit .env dengan konfigurasi database dan Redis
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run migrations**
```bash
make migrate-up
```

5. **Start server**

#### Development (with Live Reload)
```bash
# Install Air (if not installed)
go install github.com/cosmtrek/air@latest

# Start with live reload
air

# Server akan berjalan di http://localhost:8081
# Auto restart saat ada perubahan file .go
```

#### Development (Manual)
```bash
make run
# atau
go run cmd/api/main.go
```

#### Production

**Option 1: Binary Deployment**
```bash
# Build binary
make build
# atau
go build -o server cmd/api/main.go

# Run production server
GIN_MODE=release ./server

# Dengan environment variables
GIN_MODE=release \
DB_HOST=your-db-host \
DB_USER=your-db-user \
DB_PASSWORD=your-db-password \
REDIS_HOST=your-redis-host \
./server
```

**Option 2: Docker Compose (Recommended)**
```bash
# Setup environment
cp .env.production .env
# Edit .env dengan konfigurasi production

# Start all services (App + PostgreSQL + Redis)
./scripts/docker-prod.sh start

# Start with Nginx reverse proxy
./scripts/docker-prod.sh start-nginx

# Run migrations
./scripts/docker-prod.sh migrate

# View logs
./scripts/docker-prod.sh logs

# Stop services
./scripts/docker-prod.sh stop
```

## � Producbtion Deployment

### Option 1: Binary Deployment

#### Environment Configuration

```bash
# Required environment variables
export GIN_MODE=release
export DB_HOST=your-production-db-host
export DB_PORT=5432
export DB_USER=your-db-user
export DB_PASSWORD=your-secure-password
export DB_NAME=your-db-name
export REDIS_HOST=your-redis-host
export REDIS_PORT=6379
export JWT_SECRET=your-very-secure-jwt-secret
export SERVER_PORT=8081
```

### Build & Deploy

```bash
# 1. Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/api/main.go

# 2. Run migrations on production
./server migrate

# 3. Start production server
GIN_MODE=release ./server

# 4. With systemd (recommended)
sudo systemctl start huminor-rbac
sudo systemctl enable huminor-rbac
```

### Option 2: Docker Compose (Recommended)

#### Quick Start

```bash
# 1. Setup environment
cp .env.production .env
# Edit .env dengan konfigurasi production Anda

# 2. Start all services
./scripts/docker-prod.sh start

# 3. Run migrations
./scripts/docker-prod.sh migrate

# 4. Check status
./scripts/docker-prod.sh status
```

#### Available Services

- **Application**: Go API server (port 8081)
- **PostgreSQL**: Database (port 5432)
- **Redis**: Cache & sessions (port 6379)
- **Nginx**: Reverse proxy (port 80/443) - optional

#### Docker Commands

```bash
# Start services
./scripts/docker-prod.sh start

# Start with Nginx reverse proxy
./scripts/docker-prod.sh start-nginx

# View logs
./scripts/docker-prod.sh logs [service]

# Run migrations
./scripts/docker-prod.sh migrate

# Create database backup
./scripts/docker-prod.sh backup

# Stop services
./scripts/docker-prod.sh stop

# Clean up everything
./scripts/docker-prod.sh clean
```

### Docker Deployment (Optional)

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

### Performance Tuning

- **Database Connection Pool**: Configure max connections based on load
- **Redis Connection**: Use connection pooling for high traffic
- **Rate Limiting**: Adjust limits based on usage patterns
- **Logging**: Set appropriate log levels for production
- **Monitoring**: Implement health checks and metrics

## 🔧 Available Commands

```bash
# Development
make run              # Start development server (manual)
air                   # Start with live reload (recommended)
make build            # Build production binary
make test             # Run tests

# Database
make migrate-up       # Run all migrations
make migrate-down     # Rollback last migration
make migrate-reset    # Reset database

# Production
make build            # Build optimized binary
GIN_MODE=release ./server  # Run production server

# Utilities
make clean            # Clean build artifacts
make lint             # Run linter
```

## 📚 API Documentation

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - User logout

### Module System
- `GET /api/v1/modules` - List modules
- `GET /api/v1/modules/tree` - Module hierarchy
- `POST /api/v1/modules` - Create module
- `PUT /api/v1/modules/:id` - Update module
- `DELETE /api/v1/modules/:id` - Delete module

### User Management
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/password` - Change password

### Role Management
- `GET /api/v1/roles` - List roles
- `POST /api/v1/roles` - Create role
- `POST /api/v1/role-management/assign-user-role` - Assign role to user
- `PUT /api/v1/role-management/role/:id/modules` - Update role permissions

### Company & Branch
- `GET /api/v1/companies` - List companies
- `GET /api/v1/branches` - List branches
- `GET /api/v1/branches/company/:id` - Company branches with hierarchy

### Subscription
- `GET /api/v1/subscription/plans` - List subscription plans (public)
- `GET /api/v1/subscription/subscriptions` - List subscriptions
- `POST /api/v1/subscription/subscriptions` - Create subscription
- `GET /api/v1/subscription/module-access/:companyId/:moduleId` - Check module access

### Audit
- `GET /api/v1/audit/logs` - List audit logs
- `GET /api/v1/audit/stats` - Audit statistics

## 📚 Documentation

### For Developers
- **[Backend Engineer SOP](docs/BACKEND_ENGINEER_SOP.md)** - Complete guide for backend development
- **[API Infrastructure](docs/API_INFRASTRUCTURE.md)** - Response handling, rate limiting, validation
- **[Modular Architecture](docs/MODULAR_ARCHITECTURE.md)** - Project structure and patterns
- **[Models Documentation](docs/MODELS_DOCUMENTATION.md)** - Data models and relationships

### For Database
- **[Database Schema](docs/DATABASE_SCHEMA_DOCUMENTATION.md)** - Complete database documentation
- **[Migrations Guide](docs/MIGRATIONS.md)** - Database migration procedures

### For Features
- **[Password System](docs/PASSWORD_SYSTEM.md)** - Authentication and password management
- **[Subscription System](docs/SUBSCRIPTION_SYSTEM.md)** - Tiered subscription management

**Complete API Documentation**: Import `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json` ke Postman

### 🌪️ Development with Air (Live Reload)

Air provides automatic rebuilding and restarting during development:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Start development server with live reload
air

# Features:
# ✅ Auto rebuild on file changes
# ✅ Auto restart server
# ✅ Fast development cycle
# ✅ Build error display
```

**Configuration**: `.air.toml` - Pre-configured for optimal development experience

### 📋 API Changes & Updates

**Recent Changes:**
- **Module Access Endpoint**: Changed from `/companies/:companyId/modules/:moduleId/access` to `/module-access/:companyId/:moduleId` to avoid route conflicts
- **Centralized Response**: All endpoints now use consistent response format
- **Rate Limiting**: Applied globally to all endpoints
- **Parameter Validation**: Comprehensive validation on key endpoints

## 📋 Changelog

### v2.0.0 - API Infrastructure Update
- ✅ **Centralized Response System** - Consistent response format across all endpoints
- ✅ **Rate Limiting** - Global rate limiting (100 requests/minute per IP)
- ✅ **Parameter Validation** - Comprehensive validation like Joi/Zod
- ✅ **Route Fix** - Module access endpoint moved to `/module-access/:companyId/:moduleId`
- ✅ **Live Reload** - Air integration for faster development
- ✅ **Production Ready** - Optimized build and deployment process

### v1.0.0 - Initial Release
- ✅ **Raw SQL Implementation** - Removed GORM, implemented raw SQL with PostgreSQL
- ✅ **Modular Architecture** - Clean Architecture with 8 handler modules
- ✅ **Complete RBAC System** - Role-based access control with 59 modules
- ✅ **Subscription System** - Tiered subscription with module access control
- ✅ **Authentication** - JWT-based auth with refresh tokens
- ✅ **Database Migrations** - File-based migration system

## 🔐 Security Features

- **JWT Authentication** dengan access & refresh tokens
- **Password Hashing** menggunakan bcrypt
- **Role-Based Access Control** dengan granular permissions
- **Module-Based Authorization** untuk fine-grained access
- **Audit Logging** untuk semua aktivitas
- **CORS Protection** dengan configurable origins
- **Rate Limiting** untuk API endpoints

## 🗄️ Database Schema

Database menggunakan PostgreSQL dengan 6 migration files:

1. **001_create_users_table.sql** - Users dan authentication
2. **002_create_companies_and_branches.sql** - Company structure
3. **003_create_roles_and_modules.sql** - RBAC system
4. **004_seed_modules_data.sql** - Default modules (59 modules)
5. **005_create_subscription_system.sql** - Subscription management
6. **006_seed_initial_data.sql** - Sample data

## 📊 Module System

Sistem memiliki 59 modules yang terbagi dalam 8 kategori:

- **Core HR / Master Data** (12 modules)
- **Employee Self Service** (8 modules)
- **Payroll & Compensation** (10 modules)
- **Time & Attendance** (8 modules)
- **Performance Management** (6 modules)
- **Learning & Development** (5 modules)
- **Recruitment** (6 modules)
- **Reports & Analytics** (4 modules)

## 🏢 Multi-Company Support

- **Hierarchical Company Structure**
- **Branch Management** dengan parent-child relationships
- **Company-Specific Subscriptions**
- **Branch-Level User Assignments**
- **Company-Isolated Data**

## 📈 Subscription Tiers

- **Basic** - Core HR modules
- **Professional** - Extended HR features
- **Enterprise** - Full feature access
- **Custom** - Tailored solutions

## 🔍 Audit System

Comprehensive audit logging meliputi:

- **User Activities** - Login, logout, data changes
- **System Events** - Module access, permission changes
- **API Calls** - All HTTP requests dengan metadata
- **Statistics** - Activity trends dan usage analytics

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific test
go test ./internal/service/...

# Run with coverage
go test -cover ./...
```

## 📝 Development Guidelines

### Adding New Module

1. **Create Handler** di `internal/handlers/`
2. **Create Service** di `internal/service/`
3. **Create Repository** di `internal/repository/`
4. **Add Routes** di `internal/routes/routes.go`
5. **Update Server** di `internal/server/server.go`

### Code Standards

- **Clean Architecture** principles
- **Dependency Injection** pattern
- **Error Handling** yang konsisten
- **Logging** untuk debugging
- **Documentation** untuk public APIs

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Untuk pertanyaan atau dukungan:

- **Documentation**: Lihat folder `docs/`
- **Issues**: Buat issue di repository
- **API Testing**: Gunakan Postman collection yang disediakan

---

**Built with ❤️ using Go and Clean Architecture principles**