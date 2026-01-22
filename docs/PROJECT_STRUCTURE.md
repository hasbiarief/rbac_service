# Project Structure - Module-Based Architecture (5-File Structure)

## 🎯 Filosofi Desain

Project ini menggunakan **vertical module-based structure** (Express.js style), bukan horizontal layer-based.

**Prinsip Utama: 1 fitur = 1 folder**

**✅ REFACTORING COMPLETED**: Struktur telah berhasil direfactor dari 7-file menjadi 5-file per module untuk meningkatkan developer experience dan mengurangi cognitive load.

Setiap module berisi semua layer yang dibutuhkan (route + handler, service, repository, model, dto + validation) dalam satu folder.

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
│   ├── modules/                 # 🔥 SEMUA FITUR DI SINI (5 files per module)
│   │   │
│   │   ├── auth/                # Authentication module
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local User model
│   │   │   ├── repository.go    # Database queries (user data)
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/auth/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── user/                # User management
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local User model
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/users/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── role/                # Role management
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Role, UserRole models
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/roles/*, /api/v1/role-management/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── company/             # Company management
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Company model
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/companies/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── branch/              # Branch management (hierarchical)
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Branch model
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/branches/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── module/              # Module system (menu/features)
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Module, UserModule models
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/modules/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── unit/                # Unit management (unit-based RBAC)
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Unit model
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/units/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   ├── subscription/        # Subscription system
│   │   │   ├── dto.go           # Request/Response DTOs + validation logic
│   │   │   ├── model.go         # Local Plan, Subscription models
│   │   │   ├── repository.go    # Database queries
│   │   │   ├── route.go         # Routes + HTTP handlers: /api/v1/plans/*, /api/v1/subscription/*
│   │   │   └── service.go       # Business logic
│   │   │
│   │   └── audit/               # Audit logging
│   │       ├── dto.go           # Request/Response DTOs + validation logic
│   │       ├── model.go         # Local AuditLog model
│   │       ├── repository.go    # Database queries
│   │       ├── route.go         # Routes + HTTP handlers: /api/v1/audit/*
│   │       └── service.go       # Business logic
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

## 📦 Module Structure (5 Files)

Setiap module memiliki 5 file standar setelah refactoring:

1. **dto.go** - Request/Response structures + validation logic (merged dari validator.go)
2. **model.go** - Database entities (local)
3. **repository.go** - Database queries (raw SQL)
4. **route.go** - Route registration + HTTP handlers (merged dari handler.go)
5. **service.go** - Business logic

**Refactoring Benefits:**
- ✅ File count berkurang: 63 → 45 files (28% reduction)
- ✅ Faster navigation: Less file switching untuk developer
- ✅ Cleaner structure: Logical grouping of related code
- ✅ Easier onboarding: New developers less overwhelmed
- ✅ Maintained modularity: Zero impact ke cross-module dependencies

## 🔄 Data Flow

```
HTTP Request
    ↓
route.go (+ validation middleware + handler logic)
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
route.go (return response)
    ↓
HTTP Response
```

## 🚫 Folder yang TIDAK Ada

Folder-folder ini **TIDAK ADA** karena sudah diganti dengan module-based structure dan refactoring:

- ❌ `internal/interfaces/` - Interface dibuat inline
- ❌ `internal/mapper/` - Mapping dilakukan inline di service
- ❌ `internal/dto/` (global) - DTO per module
- ❌ `internal/handlers/` (global) - Handler merged ke route.go per module
- ❌ `internal/service/` (global) - Service per module
- ❌ `internal/repository/` (global) - Repository per module
- ❌ `internal/models/` (global) - Model per module
- ❌ `internal/validation/` (global) - Validator merged ke dto.go per module
- ❌ `internal/routes/` (global) - Route per module
- ❌ `internal/shared/` - Tidak digunakan
- ❌ `internal/modules/{module}/handler.go` - Merged ke route.go
- ❌ `internal/modules/{module}/validator.go` - Merged ke dto.go

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

- [Module Structure Refactoring](MODULE_STRUCTURE_REFACTORING.md) - Completed refactoring details
- [Backend Engineer Rules](ENGINEER_RULES.md) - Panduan development
- [API Overview](API_OVERVIEW.md) - API documentation
- [Quick Start Guide](QUICK_START.md) - Setup and development guide
- [README](../README.md) - Project overview
