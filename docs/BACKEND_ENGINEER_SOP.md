# Panduan Operasional Backend Engineer - RBAC API

## 🚀 Getting Started - Mulai dari Mana?

### Untuk Developer Baru di Project Ini

#### 1. Setup Awal (Wajib)
```bash
# 1. Clone dan setup dependencies
git clone <repository-url>
cd rbac-service
go mod download

# 2. Install tools
go install github.com/cosmtrek/air@latest

# 3. Setup database
createdb huminor_rbac

# 4. Setup environment
cp .env.example .env
# Edit .env dengan konfigurasi database dan Redis

# 5. Jalankan migrasi
make migrate-up

# 6. Start development server
air
```

#### 2. Test Setup (Pastikan Berjalan)
```bash
# Test health endpoint
curl http://localhost:8081/health

# Test login dengan user default
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@system.com", "password": "password123"}'
```

#### 3. Import Postman Collection
1. Import `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json`
2. Import `docs/ERP_RBAC_Environment_Module_Based.postman_environment.json`
3. Test beberapa endpoint untuk memahami flow

### 📋 Workflow Membuat Fitur Baru

#### Urutan Wajib (Jangan Diubah):

**1. Analisis & Planning**
- Pahami requirement fitur baru
- Tentukan endpoint yang dibutuhkan
- Rancang database schema (jika perlu tabel baru)

**2. Database First (Jika Perlu Tabel Baru)**
```bash
# Buat migration file
touch migrations/008_create_employees_table.sql
# Edit file dengan CREATE TABLE statement
make migrate-up
```

**3. Buat DTO (Data Transfer Objects)**

**4. Buat DTO (Data Transfer Objects)**
DTO adalah struktur data yang digunakan untuk mewakili data yang akan dikirim atau diterima melalui API.
Dalam contoh ini, kita membuat dua DTO: `CreateEmployeeRequest` untuk mendefinisikan data yang akan dikirimkan 
saat membuat employee baru, dan `EmployeeResponse` untuk mendefinisikan data yang akan dikirimkan saat 
mengambil data employee.

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

**4. Buat Interface Contracts**
Interface adalah sebuah kontrak yang menentukan apa yang harus dilakukan oleh sebuah objek. 
Dalam kasus ini, setiap objek yang mengimplementasikan `EmployeeServiceInterface` 
harus memiliki dua method: `CreateEmployee` dan `GetEmployeeByID`. 
Ini penting karena ini akan membantu kita untuk melakukan dependency injection,
sehingga kita dapat mengganti implementasi `EmployeeService` tanpa perlu mengubah client-nya.

```go
// internal/interfaces/service.go - tambahkan
type EmployeeServiceInterface interface {
    CreateEmployee(req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
    GetEmployeeByID(id int64) (*dto.EmployeeResponse, error)
}

// internal/interfaces/repository.go - tambahkan  
type EmployeeRepositoryInterface interface {
    Create(employee *models.Employee) error
    GetByID(id int64) (*models.Employee, error)
}
```

**5. Buat Mapper**
Mapper adalah fungsi atau objek yang digunakan untuk mengubah data dari salah satu model (data struct) 
menjadi model lainnya. Dalam kasus ini, Mapper digunakan untuk mengubah data dari model `models.Employee` 
menjadi model `dto.EmployeeResponse`. Ini berguna karena kita ingin mengirimkan data yang hanya diperlukan 
dalam respon HTTP, bukan semua data yang ada di database.

```go
// internal/mapper/employee_mapper.go
type EmployeeMapper struct{}

func (m *EmployeeMapper) ToResponse(employee *models.Employee) *dto.EmployeeResponse {
    return &dto.EmployeeResponse{
        ID:        employee.ID,
        Name:      employee.Name,
        Email:     employee.Email,
        CreatedAt: employee.CreatedAt.Format(time.RFC3339),
    }
}
```

**6. Implementasi Repository (Data Layer)**
Repository adalah tempat dimana kita berinteraksi dengan database. Dalam kasus ini, 
kita membuat `EmployeeRepository` untuk berinteraksi dengan tabel `employees` di database. 
Dengan menggunakan Repository ini, kita dapat melakukan operasi seperti membuat data baru, 
mendapatkan data berdasarkan ID, dan lain-lain.

```go
// internal/repository/employee_repository.go
type EmployeeRepository struct {
    *model.Repository
    db *sql.DB
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
    query := `INSERT INTO employees (name, email) VALUES ($1, $2) RETURNING id, created_at`
    return r.db.QueryRow(query, employee.Name, employee.Email).Scan(&employee.ID, &employee.CreatedAt)
}
```

**7. Implementasi Service (Business Logic)**
Service adalah tempat dimana kita menangani logika bisnis. Dalam kasus ini, 
kita membuat `EmployeeService` untuk menangani logika bisnis seperti membuat data baru, 
mendapatkan data berdasarkan ID, dan lain-lain.

```go
// internal/service/employee_service.go
type EmployeeService struct {
    employeeRepo interfaces.EmployeeRepositoryInterface
    mapper       *mapper.EmployeeMapper
}

func (s *EmployeeService) CreateEmployee(req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
    employee := &models.Employee{
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := s.employeeRepo.Create(employee); err != nil {
        return nil, err
    }
    
    return s.mapper.ToResponse(employee), nil
}
```

**8. Implementasi Handler (HTTP Layer)**
Handler adalah bagian yang berhubungan dengan HTTP. Di dalamnya, 
kita membuat fungsi-fungsi yang akan dipanggil ketika ada permintaan HTTP diterima. 
Disini, kita membuat fungsi `CreateEmployee` yang akan dipanggil ketika ada 
permintaan POST untuk membuat data baru employee.

```go
// internal/handlers/employee_handler.go
type EmployeeHandler struct {
    employeeService interfaces.EmployeeServiceInterface
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
    validatedBody, exists := c.Get("validated_body")
    if !exists {
        response.Error(c, http.StatusBadRequest, "Bad request", "validation failed")
        return
    }

    req, ok := validatedBody.(*dto.CreateEmployeeRequest)
    if !ok {
        response.Error(c, http.StatusBadRequest, "Bad request", "invalid body structure")
        return
    }

    result, err := h.employeeService.CreateEmployee(req)
    if err != nil {
        response.ErrorWithAutoStatus(c, "Failed to create employee", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, constants.MsgEmployeeCreated, result)
}
```

**9. Buat Validation Rules**
Validation adalah bagian yang berhubungan dengan validasi. Di dalamnya, 
kita membuat fungsi-fungsi yang akan dipanggil ketika ada permintaan HTTP diterima. 
Disini, kita membuat fungsi `CreateEmployee` yang akan dipanggil ketika ada 
permintaan POST untuk membuat data baru employee.

```go
// internal/validation/employee_validation.go
var CreateEmployeeValidation = middleware.ValidationRules{
    Body: &dto.CreateEmployeeRequest{},
}
```

**10. Tambah Constants**
Constants adalah bagian yang berhubungan dengan konstanta. Di dalamnya, 
kita membuat fungsi-fungsi yang akan dipanggil ketika ada permintaan HTTP diterima. 
Disini, kita membuat fungsi `CreateEmployee` yang akan dipanggil ketika ada 
permintaan POST untuk membuat data baru employee.

```go
// internal/constants/constants.go - tambahkan
const (
    MsgEmployeeCreated = "Karyawan berhasil dibuat"
    MsgEmployeeNotFound = "Karyawan tidak ditemukan"
)
```

**11. Setup Routes**
Routes adalah bagian yang berhubungan dengan routing. Di dalamnya, 
kita membuat fungsi-fungsi yang akan dipanggil ketika ada permintaan HTTP diterima. 
Disini, kita membuat fungsi `CreateEmployee` yang akan dipanggil ketika ada 
permintaan POST untuk membuat data baru employee.

```go
// internal/routes/routes.go - tambahkan function
func setupEmployeeRoutes(protected *gin.RouterGroup, employeeHandler *handlers.EmployeeHandler) {
    employees := protected.Group("/employees")
    {
        employees.POST("", middleware.ValidateRequest(validation.CreateEmployeeValidation), employeeHandler.CreateEmployee)
        employees.GET("/:id", middleware.ValidateRequest(validation.IDValidation), employeeHandler.GetEmployeeByID)
    }
}

// Di SetupRoutes function, tambahkan:
setupEmployeeRoutes(protected, h.Employee)
```

**12. Wire Dependencies di Server**
Wire adalah bagian yang berhubungan dengan wire. Di dalamnya, 
kita membuat fungsi-fungsi yang akan dipanggil ketika ada permintaan HTTP diterima. 
Disini, kita membuat fungsi `CreateEmployee` yang akan dipanggil ketika ada 
permintaan POST untuk membuat data baru employee.

```go
// internal/server/server.go
// Di initializeRepositories:
employeeRepo := repository.NewEmployeeRepository(db.DB)

// Di initializeServices:
Employee: service.NewEmployeeService(repos.Employee),

// Di initializeHandlers:
Employee: handlers.NewEmployeeHandler(services.Employee),
```

**13. Test dengan Postman**
- Tambahkan request baru di Postman collection
- Test create, get, update, delete
- Pastikan validation bekerja
- Test error scenarios

### ⚠️ Aturan Penting

#### Yang WAJIB Dilakukan:
1. **Selalu ikuti urutan di atas** - jangan loncat-loncat
2. **Test setiap layer** sebelum lanjut ke layer berikutnya
3. **Gunakan DTO** untuk semua request/response
4. **Gunakan Interface** untuk dependency injection
5. **Gunakan Mapper** untuk konversi Model ↔ DTO
6. **Validation middleware** untuk semua endpoint
7. **Constants** untuk semua pesan

#### Yang DILARANG:
1. ❌ Inline struct di handler (gunakan DTO)
2. ❌ Direct model return dari service (gunakan DTO)
3. ❌ Hardcode string messages (gunakan constants)
4. ❌ Skip validation middleware
5. ❌ Concrete types di constructor (gunakan interface)

### 🔍 Debugging Tips

#### Jika Error "invalid body structure":
- Pastikan validation menggunakan DTO yang sama dengan handler
- Check type assertion di handler

#### Jika Error "method not found":
- Pastikan interface method signature sama dengan implementasi
- Check dependency injection di server.go

#### Jika Database Error:
- Check migration sudah dijalankan
- Pastikan connection string benar di .env

## Panduan Lengkap untuk Backend Engineer

### 1. Arsitektur Clean Code

Project ini menggunakan **Clean Architecture** dengan struktur sebagai berikut:

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

### 2. Setup Development Environment

#### Prerequisites
- Go 1.25+
- PostgreSQL 13+
- Redis 6+
- Air (untuk live reload)

#### Installation
```bash
# Clone repository
git clone <repository-url>
cd rbac-service

# Install dependencies
go mod download

# Install Air untuk live reload
go install github.com/cosmtrek/air@latest

# Setup database
createdb huminor_rbac
```

#### Environment Setup
```bash
# Copy environment file
cp .env.example .env

# Edit .env dengan konfigurasi database dan Redis
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=huminor_rbac
REDIS_ADDR=localhost:6379
JWT_SECRET=your-secret-key
```

### 3. Aturan Pengembangan

#### A. Struktur Kode Wajib

**1. DTO (Data Transfer Objects)**
- Semua request/response harus menggunakan DTO
- Lokasi: `internal/dto/`
- Naming: `{Entity}Request`, `{Entity}Response`, `{Entity}ListResponse`

```go
// internal/dto/user_dto.go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

**2. Interfaces (Contracts)**
- Semua service dan repository harus memiliki interface
- Lokasi: `internal/interfaces/`

```go
// internal/interfaces/service.go
type UserServiceInterface interface {
    GetUsers(req *dto.UserListRequest) (*dto.UserListResponse, error)
    CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
    // ... methods lainnya
}
```

**3. Mappers**
- Konversi antara Model dan DTO harus menggunakan mapper
- Lokasi: `internal/mapper/`

```go
// internal/mapper/user_mapper.go
type UserMapper struct{}

func (m *UserMapper) ToResponse(user *models.User) *dto.UserResponse {
    return &dto.UserResponse{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        IsActive:  user.IsActive,
        CreatedAt: user.CreatedAt.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
    }
}
```

**4. Error Handling**
- Gunakan custom error types dari `pkg/errors/`
- Gunakan response helper dari `pkg/utils/`

```go
// Di service layer
if user == nil {
    return nil, errors.NewNotFoundError("User")
}

// Di handler layer
if err != nil {
    responseHelper.HandleError(c, err)
    return
}
```

#### B. Validation Rules

**1. Request Validation**
- Semua endpoint harus menggunakan validation middleware
- Lokasi: `internal/validation/`

```go
// internal/validation/user_validation.go
var CreateUserValidation = middleware.ValidationRules{
    Body: &dto.CreateUserRequest{},
}
```

**2. Penggunaan di Routes**
```go
// internal/routes/routes.go
users.POST("", middleware.ValidateRequest(validation.CreateUserValidation), userHandler.CreateUser)
```

#### C. Constants

**1. Pesan dan Konstanta**
- Semua pesan harus didefinisikan di `internal/constants/`
- Gunakan bahasa Indonesia untuk pesan

```go
// internal/constants/constants.go
const (
    MsgUserCreated = "Pengguna berhasil dibuat"
    MsgUserNotFound = "Pengguna tidak ditemukan"
)
```

### 4. Workflow Pengembangan Fitur Baru

#### Langkah 1: Buat DTO
```go
// internal/dto/employee_dto.go
type CreateEmployeeRequest struct {
    CompanyID  int64  `json:"company_id" validate:"required,min=1"`
    Name       string `json:"name" validate:"required,min=2,max=100"`
    Email      string `json:"email" validate:"required,email"`
    EmployeeID string `json:"employee_id" validate:"required,min=2,max=50"`
}

type EmployeeResponse struct {
    ID         int64  `json:"id"`
    CompanyID  int64  `json:"company_id"`
    Name       string `json:"name"`
    Email      string `json:"email"`
    EmployeeID string `json:"employee_id"`
    CreatedAt  string `json:"created_at"`
    UpdatedAt  string `json:"updated_at"`
}
```

#### Langkah 2: Buat Interface
```go
// internal/interfaces/service.go
type EmployeeServiceInterface interface {
    CreateEmployee(req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
    GetEmployeeByID(id int64) (*dto.EmployeeResponse, error)
}

// internal/interfaces/repository.go
type EmployeeRepositoryInterface interface {
    Create(employee *models.Employee) error
    GetByID(id int64) (*models.Employee, error)
}
```

#### Langkah 3: Buat Mapper
```go
// internal/mapper/employee_mapper.go
type EmployeeMapper struct{}

func (m *EmployeeMapper) ToResponse(employee *models.Employee) *dto.EmployeeResponse {
    return &dto.EmployeeResponse{
        ID:         employee.ID,
        CompanyID:  employee.CompanyID,
        Name:       employee.Name,
        Email:      employee.Email,
        EmployeeID: employee.EmployeeID,
        CreatedAt:  employee.CreatedAt.Format(time.RFC3339),
        UpdatedAt:  employee.UpdatedAt.Format(time.RFC3339),
    }
}

func (m *EmployeeMapper) ToModel(req *dto.CreateEmployeeRequest) *models.Employee {
    return &models.Employee{
        CompanyID:  req.CompanyID,
        Name:       req.Name,
        Email:      req.Email,
        EmployeeID: req.EmployeeID,
        IsActive:   true,
    }
}
```

#### Langkah 4: Implementasi Repository
```go
// internal/repository/employee_repository.go
type EmployeeRepository struct {
    *model.Repository
    db *sql.DB
}

func (r *EmployeeRepository) Create(employee *models.Employee) error {
    query := `
        INSERT INTO employees (company_id, name, email, employee_id, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRow(query, 
        employee.CompanyID, employee.Name, employee.Email, 
        employee.EmployeeID, employee.IsActive,
    ).Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
}
```

#### Langkah 5: Implementasi Service
```go
// internal/service/employee_service.go
type EmployeeService struct {
    employeeRepo interfaces.EmployeeRepositoryInterface
    mapper       *mapper.EmployeeMapper
}

func (s *EmployeeService) CreateEmployee(req *dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
    employee := s.mapper.ToModel(req)
    
    if err := s.employeeRepo.Create(employee); err != nil {
        return nil, errors.NewInternalServerError(err.Error())
    }
    
    return s.mapper.ToResponse(employee), nil
}
```

#### Langkah 6: Implementasi Handler
```go
// internal/handlers/employee_handler.go
type EmployeeHandler struct {
    employeeService interfaces.EmployeeServiceInterface
    responseHelper  *utils.ResponseHelper
    validationHelper *utils.ValidationHelper
}

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

#### Langkah 7: Tambah Validation
```go
// internal/validation/employee_validation.go
var CreateEmployeeValidation = middleware.ValidationRules{
    Body: &dto.CreateEmployeeRequest{},
}
```

#### Langkah 8: Tambah Routes
```go
// internal/routes/routes.go
func setupEmployeeRoutes(protected *gin.RouterGroup, employeeHandler *handlers.EmployeeHandler) {
    employees := protected.Group("/employees")
    {
        employees.POST("", middleware.ValidateRequest(validation.CreateEmployeeValidation), employeeHandler.CreateEmployee)
        employees.GET("/:id", middleware.ValidateRequest(validation.IDValidation), employeeHandler.GetEmployeeByID)
    }
}
```

### 5. Database Management

#### Membuat Migrasi Baru
```bash
# Format: migrations/XXX_description.sql
# Contoh: migrations/007_create_employees_table.sql
```

```sql
-- migrations/007_create_employees_table.sql
CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_employee_id ON employees(employee_id);

-- Comments
COMMENT ON TABLE employees IS 'Tabel data karyawan';
```

#### Menjalankan Migrasi
```bash
# Jalankan semua migrasi
make migrate-up

# Rollback migrasi terakhir
make migrate-down

# Check status migrasi
make migrate-status
```

### 6. Testing & Development

#### Menjalankan Development Server
```bash
# Live reload dengan Air
air

# Manual run
go run cmd/api/main.go
```

#### Testing API
```bash
# Test health endpoint
curl -X GET "http://localhost:8081/health"

# Test login
curl -X POST "http://localhost:8081/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@system.com",
    "password": "password123"
  }'
```

**Test Users:**
- `admin@system.com` - System Admin
- `hr@company.com` - HR Manager  
- `superadmin@company.com` - Super Admin

**Password:** `password123`

### 7. Best Practices

#### A. Kode
- **Bahasa Indonesia** untuk semua komentar dan pesan
- **Consistent naming** menggunakan camelCase untuk Go
- **Error handling** yang proper di setiap layer
- **Validation** untuk semua input

#### B. Database
- Gunakan **soft delete** (`is_active = false`)
- **Index** untuk kolom yang sering di-query
- **Transaction** untuk operasi kompleks
- **Comment** pada tabel dan kolom penting

#### C. Security
- **Validation middleware** untuk semua endpoint
- **Rate limiting** global
- **JWT token** dengan Redis storage
- **bcrypt** untuk password hashing

#### D. Performance
- **Pagination** untuk list endpoints
- **Connection pooling** untuk database
- **Redis caching** untuk data yang sering diakses
- **Query optimization** dengan proper indexing

### 8. Troubleshooting

#### Database Connection Error
```bash
# Check database status
pg_isready -h localhost -p 5432

# Check .env configuration
cat .env | grep DB_
```

#### Migration Error
```bash
# Check migration status
make migrate-status

# Force version (hati-hati!)
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" force VERSION
```

#### Build Error
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### 9. Deployment

#### Development
```bash
# Build
go build -o server cmd/api/main.go

# Run
./server
```

#### Production dengan Docker
```bash
# Setup environment
cp .env.production .env

# Build dan jalankan
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 10. Monitoring

#### Health Check
```bash
curl http://localhost:8081/health
```

#### Logs
- Server logs di console
- Error monitoring dengan structured logging
- Performance monitoring untuk response times

---

**Catatan Penting:**
- Selalu ikuti clean architecture pattern
- Test di development sebelum deploy ke production
- Backup database sebelum migration di production
- Koordinasi dengan team untuk API contract changes
- Dokumentasikan setiap perubahan di README.md