# Struktur Project - RBAC API Clean Architecture

## Gambaran Umum

Project ini menggunakan **Clean Architecture** dengan Go dan PostgreSQL, tanpa ORM. Sistem menggunakan raw SQL untuk performa optimal dan kontrol penuh atas database operations.

## Struktur Direktori Clean Architecture

```
.
├── cmd/                    # Entry points aplikasi
│   ├── api/               # Main API server
│   ├── migrate/           # Migration runner
│   └── sql-migrate/       # SQL migration tool
├── config/                # Konfigurasi aplikasi
│   ├── config.go         # Main config
│   └── redis.go          # Redis config
├── internal/              # Private application code (Clean Architecture)
│   ├── dto/              # Data Transfer Objects (Request/Response)
│   ├── interfaces/       # Contracts untuk Service & Repository
│   ├── mapper/           # Konversi antara Model dan DTO
│   ├── handlers/         # HTTP request handlers (Presentation Layer)
│   ├── service/          # Business logic layer (Use Case Layer)
│   ├── repository/       # Data access layer (Infrastructure Layer)
│   ├── models/           # Database entities (Domain Layer)
│   ├── validation/       # Request validation rules
│   ├── routes/           # Route definitions
│   ├── server/           # Server initialization
│   └── constants/        # Konstanta aplikasi
├── pkg/                   # Shared utilities (Clean Architecture)
│   ├── errors/           # Custom error types
│   ├── query/            # Query builder utilities
│   ├── pagination/       # Pagination helpers
│   ├── utils/            # Response & validation helpers
│   ├── database/         # Database connection
│   ├── response/         # Response utilities
│   ├── token/            # Token management
│   └── jobs/             # Background jobs
├── middleware/            # HTTP middleware
│   ├── auth.go           # JWT authentication
│   ├── cors.go           # CORS handling
│   ├── rate_limit.go     # Rate limiting
│   └── validation.go     # Request validation
├── migrations/            # Database migrations
├── docs/                 # Documentation
├── .env                  # Environment variables
├── go.mod               # Go modules
├── Makefile             # Build commands
└── README.md            # Project overview
```

## Clean Architecture Layers

### 1. Domain Layer (`internal/models/`)
**Tanggung jawab**: Business entities dan domain logic
- Database entities dengan struct tags
- TableName() methods untuk custom table names
- Domain-specific business rules

```go
// internal/models/user.go
type User struct {
    ID           int64     `json:"id" db:"id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    UserIdentity *string   `json:"user_identity" db:"user_identity"`
    PasswordHash string    `json:"-" db:"password_hash"`
    IsActive     bool      `json:"is_active" db:"is_active"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (User) TableName() string {
    return "users"
}
```

### 2. Use Case Layer (`internal/service/`)
**Tanggung jawab**: Business logic dan orchestration
- Menggunakan repository interfaces
- Menggunakan mapper untuk konversi
- Business rules enforcement

```go
// internal/service/user_service.go
type UserService struct {
    userRepo interfaces.UserRepositoryInterface
    mapper   *mapper.UserMapper
}

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Business logic
    user := s.mapper.ToModel(req)
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, errors.NewInternalServerError(err.Error())
    }
    
    return s.mapper.ToResponse(user), nil
}
```

### 3. Interface Adapters Layer

#### A. Controllers (`internal/handlers/`)
**Tanggung jawab**: HTTP request/response handling
- Menggunakan validation helpers
- Menggunakan response helpers
- Error handling

```go
// internal/handlers/user_handler.go
type UserHandler struct {
    userService      interfaces.UserServiceInterface
    responseHelper   *utils.ResponseHelper
    validationHelper *utils.ValidationHelper
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req *dto.CreateUserRequest
    if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    result, err := h.userService.CreateUser(req)
    if err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    h.responseHelper.Created(c, constants.MsgUserCreated, result)
}
```

#### B. DTOs (`internal/dto/`)
**Tanggung jawab**: Data transfer objects untuk API
- Request structures dengan validation tags
- Response structures untuk API output
- List responses dengan pagination

```go
// internal/dto/user_dto.go
type CreateUserRequest struct {
    Name         string  `json:"name" validate:"required,min=2,max=100"`
    Email        string  `json:"email" validate:"required,email"`
    UserIdentity *string `json:"user_identity"`
    Password     string  `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
    ID           int64   `json:"id"`
    Name         string  `json:"name"`
    Email        string  `json:"email"`
    UserIdentity *string `json:"user_identity"`
    IsActive     bool    `json:"is_active"`
    CreatedAt    string  `json:"created_at"`
    UpdatedAt    string  `json:"updated_at"`
}

type UserListResponse struct {
    Data    []*UserResponse `json:"data"`
    Total   int64           `json:"total"`
    Limit   int             `json:"limit"`
    Offset  int             `json:"offset"`
    HasMore bool            `json:"has_more"`
}
```

#### C. Mappers (`internal/mapper/`)
**Tanggung jawab**: Konversi antara Domain dan DTO
- Model to DTO conversion
- DTO to Model conversion
- List conversions dengan pagination

```go
// internal/mapper/user_mapper.go
type UserMapper struct{}

func (m *UserMapper) ToResponse(user *models.User) *dto.UserResponse {
    if user == nil {
        return nil
    }
    
    return &dto.UserResponse{
        ID:           user.ID,
        Name:         user.Name,
        Email:        user.Email,
        UserIdentity: user.UserIdentity,
        IsActive:     user.IsActive,
        CreatedAt:    user.CreatedAt.Format(time.RFC3339),
        UpdatedAt:    user.UpdatedAt.Format(time.RFC3339),
    }
}

func (m *UserMapper) ToModel(req *dto.CreateUserRequest) *models.User {
    return &models.User{
        Name:         req.Name,
        Email:        req.Email,
        UserIdentity: req.UserIdentity,
        IsActive:     true,
    }
}
```

#### D. Interfaces (`internal/interfaces/`)
**Tanggung jawab**: Contracts untuk dependency inversion
- Service interfaces
- Repository interfaces
- Dependency injection contracts

```go
// internal/interfaces/service.go
type UserServiceInterface interface {
    GetUsers(req *dto.UserListRequest) (*dto.UserListResponse, error)
    GetUserByID(id int64) (*dto.UserResponse, error)
    CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
    UpdateUser(id int64, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
    DeleteUser(id int64) error
}

// internal/interfaces/repository.go
type UserRepositoryInterface interface {
    GetAll(limit, offset int, search string, isActive *bool) ([]*models.User, error)
    GetByID(id int64) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Create(user *models.User) error
    Update(user *models.User) error
    Delete(id int64) error
}
```

### 4. Infrastructure Layer (`internal/repository/`)
**Tanggung jawab**: Data access dengan raw SQL
- Database operations
- Query building
- Transaction management

```go
// internal/repository/user_repository.go
type UserRepository struct {
    *model.Repository
    db *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (name, email, user_identity, password_hash, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRow(query, 
        user.Name, user.Email, user.UserIdentity, 
        user.PasswordHash, user.IsActive,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
```

## Shared Utilities (`pkg/`)

### 1. Custom Error Types (`pkg/errors/`)
```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func NewNotFoundError(resource string) *AppError {
    return &AppError{
        Code:    http.StatusNotFound,
        Message: fmt.Sprintf("%s tidak ditemukan", resource),
    }
}
```

### 2. Query Builder (`pkg/query/`)
```go
type QueryBuilder struct {
    baseQuery  string
    conditions []string
    args       []any
    argIndex   int
}

func (qb *QueryBuilder) AddCondition(condition string, value any) *QueryBuilder {
    qb.conditions = append(qb.conditions, fmt.Sprintf(condition, qb.argIndex))
    qb.args = append(qb.args, value)
    qb.argIndex++
    return qb
}
```

### 3. Response Helper (`pkg/utils/`)
```go
type ResponseHelper struct{}

func (h *ResponseHelper) Success(c *gin.Context, message string, data interface{}) {
    response.Success(c, http.StatusOK, message, data)
}

func (h *ResponseHelper) HandleError(c *gin.Context, err error) {
    if appErr, ok := err.(*errors.AppError); ok {
        response.Error(c, appErr.Code, appErr.Message, appErr.Details)
        return
    }
    response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
}
```

## Validation System (`internal/validation/`)

**Tanggung jawab**: Centralized validation rules
- Menggunakan DTO types untuk validation
- Modular validation rules
- Reusable validation patterns

```go
// internal/validation/user_validation.go
var CreateUserValidation = middleware.ValidationRules{
    Body: &dto.CreateUserRequest{},
}

var UpdateUserValidation = middleware.ValidationRules{
    Params: IDValidation.Params,
    Body:   &dto.UpdateUserRequest{},
}
```

## Constants (`internal/constants/`)

**Tanggung jawab**: Centralized constants dan messages
- HTTP status messages dalam bahasa Indonesia
- Validation constants
- Business constants

```go
// internal/constants/constants.go
const (
    MsgUserCreated   = "Pengguna berhasil dibuat"
    MsgUserUpdated   = "Pengguna berhasil diperbarui"
    MsgUserNotFound  = "Pengguna tidak ditemukan"
)
```

## API Response Format

Semua API menggunakan format response yang konsisten:

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

## Development Workflow

### 1. Membuat Fitur Baru (Clean Architecture)

1. **Buat DTO** di `internal/dto/`
2. **Buat Interface** di `internal/interfaces/`
3. **Buat Mapper** di `internal/mapper/`
4. **Implementasi Repository** di `internal/repository/`
5. **Implementasi Service** di `internal/service/`
6. **Implementasi Handler** di `internal/handlers/`
7. **Tambah Validation** di `internal/validation/`
8. **Tambah Routes** di `internal/routes/`
9. **Update Constants** di `internal/constants/`

### 2. Dependency Flow

```
Handler -> Service -> Repository -> Database
   ↓         ↓           ↓
  DTO    Interface    Model
   ↓         ↓           ↓
Mapper -> Mapper -> Mapper
```

### 3. Error Flow

```
Repository -> Custom Error -> Service -> Handler -> Response Helper -> Client
```

## Best Practices Clean Architecture

### 1. Dependency Rule
- Inner layers tidak boleh tahu tentang outer layers
- Outer layers bergantung pada inner layers melalui interfaces
- Business logic tidak tahu tentang HTTP atau database

### 2. Interface Segregation
- Buat interface yang spesifik dan kecil
- Satu interface untuk satu tanggung jawab
- Gunakan dependency injection

### 3. Single Responsibility
- Setiap layer memiliki tanggung jawab yang jelas
- Handler hanya untuk HTTP concerns
- Service hanya untuk business logic
- Repository hanya untuk data access

### 4. Consistent Patterns
- Gunakan DTO untuk semua API communication
- Gunakan mapper untuk semua conversions
- Gunakan custom errors untuk error handling
- Gunakan constants untuk messages

## Monitoring & Maintenance

### Health Check
```bash
curl http://localhost:8081/health
```

### Database Performance
```bash
# Check query performance
psql -d huminor_rbac -c "SELECT * FROM pg_stat_activity;"

# Check index usage
psql -d huminor_rbac -c "SELECT * FROM pg_stat_user_indexes;"
```

### Code Quality
```bash
# Run tests
go test ./...

# Check coverage
go test -cover ./...

# Static analysis
golangci-lint run
```

---

**Catatan**: Project ini menggunakan Clean Architecture untuk maintainability, testability, dan scalability yang optimal.