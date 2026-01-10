# Clean Architecture - RBAC API

## Gambaran Umum

Project ini mengimplementasikan **Clean Architecture** dengan prinsip **Dependency Inversion** dan **Separation of Concerns** yang jelas. 
Setiap layer memiliki tanggung jawab yang spesifik dan tidak bergantung pada layer luar.

## Prinsip Clean Architecture

### 1. **Dependency Rule**
- **Inner layers** tidak boleh tahu tentang **outer layers**
- **Outer layers** bergantung pada **inner layers** melalui **interfaces**
- **Business logic** tidak tahu tentang HTTP, database, atau framework

### 2. **Separation of Concerns**
- Setiap layer memiliki **tanggung jawab tunggal**
- **Handler** hanya untuk HTTP concerns
- **Service** hanya untuk business logic
- **Repository** hanya untuk data access

### 3. **Interface Segregation**
- Interface yang **spesifik** dan **kecil**
- Satu interface untuk satu tanggung jawab
- Dependency injection melalui interfaces

## Struktur Layer

```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │   Routes    │  │  Handlers   │  │   Validation    │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                 Interface Adapters Layer                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │     DTO     │  │   Mapper    │  │   Interfaces    │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                    Use Case Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │   Service   │  │  Constants  │  │  Custom Errors  │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                     Domain Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │   Models    │  │ Business    │  │   Entities      │  │
│  │             │  │   Rules     │  │                 │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                    │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │ Repository  │  │  Database   │  │   External      │  │
│  │             │  │             │  │   Services      │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Implementasi per Layer

### 1. Domain Layer (`internal/models/`)

**Tanggung jawab**: Business entities dan domain logic
- Database entities dengan struct tags
- Domain-specific business rules
- TableName() methods

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
    // Validasi business rules
    if req.Email == "" {
        return nil, errors.NewValidationError("Email wajib diisi")
    }
    
    // Konversi DTO ke Model
    user := s.mapper.ToModel(req)
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.NewInternalServerError("Gagal hash password")
    }
    user.PasswordHash = string(hashedPassword)
    
    // Simpan ke database
    if err := s.userRepo.Create(user); err != nil {
        return nil, errors.NewInternalServerError(err.Error())
    }
    
    // Konversi Model ke DTO Response
    return s.mapper.ToResponse(user), nil
}
```

### 3. Interface Adapters Layer

#### A. DTOs (`internal/dto/`)

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
```

#### B. Mappers (`internal/mapper/`)

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
        IsActive:     true, // Default aktif
    }
}
```

#### C. Interfaces (`internal/interfaces/`)

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

### 4. Presentation Layer (`internal/handlers/`)

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
    // Extract validated request
    var req *dto.CreateUserRequest
    if err := h.validationHelper.GetValidatedBody(c, &req); err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    // Call business logic
    result, err := h.userService.CreateUser(req)
    if err != nil {
        h.responseHelper.HandleError(c, err)
        return
    }
    
    // Send response
    h.responseHelper.Created(c, constants.MsgUserCreated, result)
}
```

### 5. Infrastructure Layer (`internal/repository/`)

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

func NewValidationError(details string) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: "Validasi gagal",
        Details: details,
    }
}
```

### 2. Response Helper (`pkg/utils/`)

```go
type ResponseHelper struct{}

func (h *ResponseHelper) Success(c *gin.Context, message string, data interface{}) {
    response.Success(c, http.StatusOK, message, data)
}

func (h *ResponseHelper) Created(c *gin.Context, message string, data interface{}) {
    response.Success(c, http.StatusCreated, message, data)
}

func (h *ResponseHelper) HandleError(c *gin.Context, err error) {
    if appErr, ok := err.(*errors.AppError); ok {
        response.Error(c, appErr.Code, appErr.Message, appErr.Details)
        return
    }
    response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
}
```

### 3. Query Builder (`pkg/query/`)

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

func (qb *QueryBuilder) Build() (string, []any) {
    query := qb.baseQuery
    if len(qb.conditions) > 0 {
        query += " WHERE " + strings.Join(qb.conditions, " AND ")
    }
    return query, qb.args
}
```

## Data Flow

### 1. Request Flow

```
HTTP Request → Middleware → Handler → Service → Repository → Database
     ↓             ↓          ↓         ↓          ↓
 Validation → DTO → Business → Model → SQL Query
```

### 2. Response Flow

```
Database → Repository → Service → Handler → HTTP Response
    ↓          ↓         ↓        ↓
  Model → Model → DTO → JSON
```

### 3. Error Flow

```
Repository → Custom Error → Service → Handler → Response Helper → Client
     ↓            ↓           ↓        ↓            ↓
Database Error → App Error → App Error → HTTP Error → JSON Error
```

## Keuntungan Clean Architecture

### 1. **Testability**
- Setiap layer dapat ditest secara terpisah
- Mock interfaces untuk unit testing
- Business logic terpisah dari framework

```go
// Test service tanpa database
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    mapper := mapper.NewUserMapper()
    service := service.NewUserService(mockRepo, mapper)
    
    req := &dto.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    result, err := service.CreateUser(req)
    
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", result.Name)
}
```

### 2. **Maintainability**
- Perubahan di satu layer tidak mempengaruhi layer lain
- Code yang modular dan terorganisir
- Easy to extend dan modify

### 3. **Scalability**
- Mudah menambah fitur baru
- Consistent patterns
- Reusable components

### 4. **Independence**
- Framework independence
- Database independence
- External service independence

## Best Practices

### 1. **Dependency Injection**
```go
// Gunakan interfaces untuk dependency injection
type UserService struct {
    userRepo interfaces.UserRepositoryInterface  // Interface, bukan concrete type
    mapper   *mapper.UserMapper
}
```

### 2. **Error Handling**
```go
// Gunakan custom error types
if user == nil {
    return nil, errors.NewNotFoundError("User")
}

// Handle errors di handler layer
if err != nil {
    h.responseHelper.HandleError(c, err)
    return
}
```

### 3. **Validation**
```go
// Validation di DTO level
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}

// Business validation di service level
if req.Email == "" {
    return nil, errors.NewValidationError("Email wajib diisi")
}
```

### 4. **Consistent Patterns**
- Semua handler menggunakan pattern yang sama
- Semua service menggunakan interfaces
- Semua repository menggunakan raw SQL
- Semua response menggunakan helper

## Testing Strategy

### 1. **Unit Testing**
```go
// Test service dengan mock repository
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, mapper.NewUserMapper())
    
    // Test business logic
}

// Test mapper
func TestUserMapper_ToResponse(t *testing.T) {
    mapper := NewUserMapper()
    user := &models.User{Name: "John"}
    
    result := mapper.ToResponse(user)
    
    assert.Equal(t, "John", result.Name)
}
```

### 2. **Integration Testing**
```go
// Test handler dengan real service dan mock repository
func TestUserHandler_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, mapper.NewUserMapper())
    handler := NewUserHandler(service, responseHelper, validationHelper)
    
    // Test HTTP flow
}
```

### 3. **End-to-End Testing**
```go
// Test dengan real database
func TestUserAPI_CreateUser(t *testing.T) {
    // Setup test database
    // Test complete flow
}
```

---

**Kesimpulan**: Clean Architecture memberikan struktur yang solid, maintainable, dan scalable untuk pengembangan jangka panjang dengan separation of concerns yang jelas.