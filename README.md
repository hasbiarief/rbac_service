# ERP RBAC API

Sistem ERP dengan Role-Based Access Control (RBAC) menggunakan Go, PostgreSQL, dan Redis. Menggunakan raw SQL tanpa ORM untuk performa optimal.

## Fitur Utama

- **Authentication** - Custom token system dengan Redis storage (1 access + 1 refresh per user)
- **RBAC System** - Role-based access control dengan module permissions
- **Multi-Company** - Company dan branch hierarchy (4 levels)
- **Subscription** - Tiered subscription dengan module access control
- **Audit Logging** - Complete audit trail
- **User Management** - User lifecycle dengan soft delete

## Tech Stack

- **Go 1.25+** dengan Gin framework
- **PostgreSQL** dengan raw SQL (tanpa ORM)
- **Redis** untuk custom token storage dan caching
- **Custom Token System** untuk authentication (mudah revoke)
- **bcrypt** untuk password hashing

## Quick Start

### Prerequisites
- Go 1.25+
- PostgreSQL 12+
- Redis 6+

### Development Setup

```bash
# 1. Clone dan setup
git clone <repository>
cd project
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

## API Documentation

### Test Users
Default users untuk testing (password: `password123`):
- `100000001` - System Admin | `admin@system.com`
- `100000002` - HR Manager | `hr@company.com`  
- `100000003` - Super Admin | `superadmin@company.com`

### Key Endpoints

**Authentication:**
- `POST /api/v1/auth/login` - Login (user_identity atau email)
- `POST /api/v1/auth/refresh` - Refresh custom token
- `POST /api/v1/auth/logout` - Logout (revoke tokens dari Redis)
- `GET /api/v1/auth/check-tokens?user_id=1` - Check user tokens di Redis

**User Management:**
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Soft delete user

**Company & Branch:**
- `GET /api/v1/companies` - List companies
- `GET /api/v1/branches/company/:id` - Company branches (hierarchy)

**Subscription:**
- `GET /api/v1/plans` - List plans (public)
- `POST /api/v1/subscription/subscriptions` - Create subscription
- `GET /api/v1/subscription/module-access/:companyId/:moduleId` - Check access

**Complete API Collection:** Import `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json` ke Postman

## Documentation

- **[Backend SOP](docs/BACKEND_ENGINEER_SOP.md)** - Complete development guide
- **[Migration Guide](docs/MIGRATIONS.md)** - Database migration procedures  
- **[Project Structure](docs/PROJECT_STRUCTURE.md)** - Architecture overview
- **[Role Permissions Mapping](docs/ROLE_PERMISSIONS_MAPPING.md)** - Complete role-based access control documentation
- **[Postman Collection](docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json)** - API testing

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=huminor_rbac

# Redis
REDIS_ADDR=localhost:6379

# Token System
TOKEN_SECRET=your-secret-key

# Server
PORT=8081
GIN_MODE=debug
```

---

**Built with Go + PostgreSQL + Redis**