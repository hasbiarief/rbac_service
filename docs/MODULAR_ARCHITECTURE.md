# Arsitektur Modular - ERP RBAC API

## Gambaran Umum

Project ini telah direstrukturisasi menjadi arsitektur modular yang lebih terorganisir dan mudah dipelihara. Setiap modul memiliki tanggung jawab yang jelas dan terpisah.

## Struktur Direktori

```
internal/
├── handlers/           # HTTP request handlers berdasarkan modul
│   ├── auth_handler.go
│   ├── module_handler.go
│   ├── user_handler.go
│   ├── company_handler.go
│   ├── role_handler.go
│   ├── subscription_handler.go
│   ├── audit_handler.go
│   └── branch_handler.go
├── routes/             # Konfigurasi routing
│   └── routes.go
├── server/             # Server initialization dan dependency injection
│   └── server.go
├── service/            # Business logic layer
├── repository/         # Data access layer
└── models/             # Data models
```

## Komponen Utama

### 1. Handlers (`internal/handlers/`)

Setiap handler bertanggung jawab untuk:
- Memproses HTTP requests
- Validasi input
- Memanggil service layer
- Mengembalikan HTTP responses

**Modul Handlers:**
- `AuthHandler` - Authentication (login, logout, refresh token)
- `ModuleHandler` - Module management (CRUD, tree, hierarchy)
- `UserHandler` - User management (CRUD, password, module access)
- `CompanyHandler` - Company management (CRUD)
- `RoleHandler` - Role management (basic & advanced RBAC)
- `SubscriptionHandler` - Subscription management (plans, billing)
- `AuditHandler` - Audit logging (logs, statistics)
- `BranchHandler` - Branch management (hierarchical structure)

### 2. Routes (`internal/routes/`)

File `routes.go` mengorganisir semua routing berdasarkan modul:

```go
func SetupRoutes(r *gin.Engine, h *Handlers, jwtSecret string, redis *redis.Client)
```

**Pengelompokan Routes:**
- Public routes (auth, public subscription plans)
- Protected routes (semua endpoint yang memerlukan authentication)
- Modular route setup functions untuk setiap modul

### 3. Server (`internal/server/`)

File `server.go` mengelola:
- Database connection initialization
- Repository initialization
- Service initialization
- Handler initialization
- Dependency injection
- Server startup

**Struktur Utama:**
```go
type Server struct {
    router *gin.Engine
    config *config.Config
}

type Repositories struct {
    User, Module, Company, Role, Subscription, Audit, Branch
}

type Services struct {
    Auth, Module, Company, Role, User, Subscription, Audit, Branch
}
```

## Keuntungan Arsitektur Modular

### 1. **Separation of Concerns**
- Setiap handler hanya menangani satu domain
- Business logic terpisah dari HTTP handling
- Database access terpisah dari business logic

### 2. **Maintainability**
- Kode lebih mudah dibaca dan dipahami
- Perubahan pada satu modul tidak mempengaruhi modul lain
- Testing lebih mudah dilakukan per modul

### 3. **Scalability**
- Mudah menambah modul baru
- Mudah menambah endpoint baru pada modul yang ada
- Dependency injection memudahkan testing dan mocking

### 4. **Code Reusability**
- Handler dapat digunakan kembali
- Service layer dapat digunakan oleh multiple handlers
- Repository pattern memudahkan database abstraction

## Cara Menambah Modul Baru

### 1. Buat Handler Baru
```go
// internal/handlers/new_module_handler.go
type NewModuleHandler struct {
    newModuleService *service.NewModuleService
}

func NewNewModuleHandler(service *service.NewModuleService) *NewModuleHandler {
    return &NewModuleHandler{newModuleService: service}
}

func (h *NewModuleHandler) GetItems(c *gin.Context) {
    // Implementation
}
```

### 2. Tambahkan ke Routes
```go
// internal/routes/routes.go
func setupNewModuleRoutes(protected *gin.RouterGroup, handler *handlers.NewModuleHandler) {
    newModule := protected.Group("/new-module")
    {
        newModule.GET("", handler.GetItems)
        newModule.POST("", handler.CreateItem)
        // ... other routes
    }
}
```

### 3. Update Server Initialization
```go
// internal/server/server.go
type Repositories struct {
    // ... existing repos
    NewModule *repository.NewModuleRepository
}

type Services struct {
    // ... existing services
    NewModule *service.NewModuleService
}

func (s *Server) initializeHandlers(services *Services, repos *Repositories) *routes.Handlers {
    return &routes.Handlers{
        // ... existing handlers
        NewModule: handlers.NewNewModuleHandler(services.NewModule),
    }
}
```

## Best Practices

### 1. **Handler Guidelines**
- Minimal business logic di handler
- Selalu validasi input
- Gunakan consistent response format
- Handle errors dengan proper HTTP status codes

### 2. **Service Guidelines**
- Semua business logic ada di service layer
- Service tidak boleh tahu tentang HTTP
- Gunakan repository untuk database access
- Return domain errors, bukan HTTP errors

### 3. **Repository Guidelines**
- Hanya database operations
- Return domain models
- Handle database-specific errors
- Use transactions untuk complex operations

### 4. **Error Handling**
- Consistent error response format
- Proper HTTP status codes
- Log errors untuk debugging
- Don't expose internal errors ke client

## Testing Strategy

### 1. **Unit Testing**
- Test setiap handler secara terpisah
- Mock dependencies (services, repositories)
- Test happy path dan error cases

### 2. **Integration Testing**
- Test end-to-end flow
- Test dengan real database
- Test authentication dan authorization

### 3. **API Testing**
- Gunakan Postman collection yang sudah ada
- Test semua endpoints
- Test dengan different user roles

## Migration dari Struktur Lama

Struktur lama (monolithic `main.go` dengan 1400+ lines) telah berhasil dipecah menjadi:

- **8 Handler files** (~100-200 lines each)
- **1 Routes file** (~200 lines)
- **1 Server file** (~150 lines)
- **1 Main file** (~20 lines)

**Total reduction:** Dari 1400+ lines menjadi ~1000 lines yang terdistribusi dengan baik.

## Kesimpulan

Arsitektur modular ini memberikan:
- ✅ **Better organization** - Kode terstruktur berdasarkan domain
- ✅ **Easier maintenance** - Perubahan terisolasi per modul
- ✅ **Better testing** - Unit testing per modul
- ✅ **Scalability** - Mudah menambah fitur baru
- ✅ **Team collaboration** - Developer dapat bekerja pada modul berbeda
- ✅ **Code reusability** - Handler dan service dapat digunakan kembali

Project sekarang siap untuk pengembangan lebih lanjut dengan struktur yang solid dan maintainable.