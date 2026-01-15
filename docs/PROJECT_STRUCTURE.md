# Project Structure - Module-Based Architecture

## 🎯 Filosofi Desain

Project ini menggunakan **vertical module-based structure** (Express.js style), bukan horizontal layer-based.

**Prinsip Utama: 1 fitur = 1 folder**

Setiap module berisi semua layer yang dibutuhkan (route, handler, service, repository, model, dto, validator) dalam satu folder.

## 📁 Struktur Folder

```
rbac-service/
├── cmd/
│   ├── api/
│   │   └── main.go              # Entry point aplikasi
│   ├── migrate/
│   │   └── main.go              # Migration tool
│   └── sql-migrate/
│       └── main.go              # SQL migration tool
│
├── config/
│   ├── config.go                # Configuration loader
│   └── redis.go                 # Redis configuration
│
├── internal/
│   ├── app/
│   │   ├── server.go            # Server initialization
│   │   └── routes.go            # Route registration
│   │
│   ├── modules/                 # 🔥 SEMUA FITUR DI SINI
│   │   │
│   │   ├── auth/                # Authentication module
│   │   │   ├── route.go         # Routes: /api/v1/auth/*
│   │   │   ├── handler.go       # HTTP handlers
│   │   │   ├── service.go       # Business logic
│   │   │   ├── repository.go    # Database queries (user data)
│   │   │   ├── model.go         # Local User model
│   │   │   ├── dto.go           # Request/Response DTOs
│   │   │   └── validator.go    # Validation rules
│   │   │
│   │   ├── user/                # User management
│   │   │   ├── route.go         # Routes: /api/v1/users/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local User model
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── role/                # Role management
│   │   │   ├── route.go         # Routes: /api/v1/roles/*, /api/v1/role-management/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Role, UserRole models
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── company/             # Company management
│   │   │   ├── route.go         # Routes: /api/v1/companies/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Company model
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── branch/              # Branch management (hierarchical)
│   │   │   ├── route.go         # Routes: /api/v1/branches/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Branch model
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── module/              # Module system (menu/features)
│   │   │   ├── route.go         # Routes: /api/v1/modules/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Module, UserModule models
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── unit/                # Unit management (unit-based RBAC)
│   │   │   ├── route.go         # Routes: /api/v1/units/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Unit model
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   ├── subscription/        # Subscription system
│   │   │   ├── route.go         # Routes: /api/v1/plans/*, /api/v1/subscription/*
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── model.go         # Local Plan, Subscription models
│   │   │   ├── dto.go
│   │   │   └── validator.go
│   │   │
│   │   └── audit/               # Audit logging
│   │       ├── route.go         # Routes: /api/v1/audit/*
│   │       ├── handler.go
│   │       ├── service.go
│   │       ├── repository.go
│   │       ├── model.go         # Local AuditLog model
│   │       ├── dto.go
│   │       └── validator.go
│   │
│   └── constants/               # Shared constants
│       └── constants.go         # API messages, status codes
│
├── middleware/                  # HTTP middleware
│   ├── auth.go                  # JWT authentication
│   ├── cors.go                  # CORS configuration
│   ├── rate_limit.go            # Rate limiting
│   └── validation.go            # Request validation
│
├── migrations/                  # SQL migrations
│   ├── 001_init.sql
│   ├── 002_rbac.sql
│   └── ...
│
├── pkg/                         # Reusable utilities (generic)
│   ├── model/
│   │   └── repository.go        # Base repository helper
│   ├── response/
│   │   └── response.go          # Response helpers
│   └── utils/
│       └── utils.go             # General utilities
│
├── scripts/                     # Development scripts
│   ├── dev.sh
│   ├── prod.sh
│   └── ...
│
├── docs/                        # Documentation
│   ├── ENGINEER_RULES.md
│   ├── PROJECT_STRUCTURE.md
│   ├── API_OVERVIEW.md
│   └── *.postman_collection.json
│
├── .env                         # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔑 Key Principles

### 1. No Cross-Module Imports
❌ **TIDAK BOLEH:**
```go
import "gin-scalable-api/internal/modules/user"  // dari module lain
```

✅ **BOLEH:**
- Duplicate models jika context berbeda
- Query database langsung dengan minimal fields
- Import dari `pkg/` (utilities)
- Import dari `internal/constants/`
- Import dari `middleware/`

### 2. Local Models per Module
Setiap module punya model lokalnya sendiri. Tidak ada shared models.

**Contoh:** Module `auth` dan `user` sama-sama punya model `User`, tapi dengan fields yang berbeda sesuai kebutuhan:

```go
// internal/modules/auth/model.go
type User struct {
    ID           int64
    Email        string
    PasswordHash string  // auth butuh password
    IsActive     bool
}

// internal/modules/user/model.go
type User struct {
    ID           int64
    Name         string
    Email        string
    UserIdentity *string
    IsActive     bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
    // Tidak ada PasswordHash - user management tidak butuh
}
```

### 3. Repository per Module
Setiap module punya repository lokalnya sendiri.

**Tidak ada:**
- `internal/repository/` global
- Interface + implementation pattern (over-engineering)

**Ada:**
- `internal/modules/{module}/repository.go` dengan concrete struct

### 4. No Interface Folder
Interface dibuat inline di file yang membutuhkan, bukan di folder terpisah.

### 5. No Mapper Folder
Konversi Model ↔ DTO dilakukan inline di service, tidak perlu mapper terpisah.

## 📦 Module Structure (7 Files)

Setiap module memiliki 7 file standar:

1. **route.go** - Route registration
2. **handler.go** - HTTP handlers
3. **service.go** - Business logic
4. **repository.go** - Database queries (raw SQL)
5. **model.go** - Database entities (local)
6. **dto.go** - Request/Response structures
7. **validator.go** - Custom validation rules

## 🔄 Data Flow

```
HTTP Request
    ↓
route.go (+ validation middleware)
    ↓
handler.go (parse request)
    ↓
service.go (business logic)
    ↓
repository.go (database query)
    ↓
database
    ↓
repository.go (return model)
    ↓
service.go (convert to DTO)
    ↓
handler.go (return response)
    ↓
HTTP Response
```

## 🚫 Folder yang TIDAK Ada

Folder-folder ini **TIDAK ADA** karena sudah diganti dengan module-based structure:

- ❌ `internal/interfaces/` - Interface dibuat inline
- ❌ `internal/mapper/` - Mapping dilakukan inline di service
- ❌ `internal/dto/` (global) - DTO per module
- ❌ `internal/handlers/` (global) - Handler per module
- ❌ `internal/service/` (global) - Service per module
- ❌ `internal/repository/` (global) - Repository per module
- ❌ `internal/models/` (global) - Model per module
- ❌ `internal/validation/` (global) - Validator per module
- ❌ `internal/routes/` (global) - Route per module
- ❌ `internal/shared/` - Tidak digunakan

## 🎯 Kapan Membuat Module Baru?

Buat module baru ketika:
- Fitur baru yang independent
- Punya domain logic sendiri
- Punya database table sendiri
- Punya endpoint API sendiri

**Contoh:**
- ✅ `employee` - Fitur employee management
- ✅ `attendance` - Fitur attendance tracking
- ✅ `payroll` - Fitur payroll processing
- ❌ `helpers` - Bukan fitur, taruh di `pkg/`
- ❌ `utils` - Bukan fitur, taruh di `pkg/`

## 📚 Related Documentation

- [Backend Engineer Rules](ENGINEER_RULES.md) - Panduan development
- [API Overview](API_OVERVIEW.md) - API documentation
- [README](../README.md) - Project overview
