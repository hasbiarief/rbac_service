# Struktur Project - ERP RBAC API

## Gambaran Umum

Project ini menggunakan arsitektur modular dengan Go dan PostgreSQL, tanpa GORM. Sistem menggunakan raw SQL untuk performa optimal dan kontrol penuh atas database operations.

## Struktur Direktori

```
.
├── cmd/                    # Entry points aplikasi
│   ├── api/               # Main API server
│   ├── migrate/           # Migration runner
│   └── sql-migrate/       # SQL migration tool
├── config/                # Konfigurasi aplikasi
│   ├── config.go         # Main config
│   └── redis.go          # Redis config
├── internal/              # Private application code
│   ├── handlers/         # HTTP request handlers
│   ├── service/          # Business logic layer
│   ├── repository/       # Data access layer (Raw SQL)
│   ├── models/           # Data models
│   ├── routes/           # Route definitions
│   ├── server/           # Server initialization
│   └── validation/       # Request validation rules
├── middleware/            # HTTP middleware
│   ├── auth.go           # Custom token authentication
│   ├── cors.go           # CORS handling
│   ├── rate_limit.go     # Rate limiting
│   └── validation.go     # Request validation
├── pkg/                   # Public packages
│   ├── database/         # Database connection
│   ├── migration/        # Migration system
│   ├── model/            # Base model helpers
│   ├── response/         # Response utilities
│   ├── token/            # Custom token management
│   ├── password/         # Password utilities
│   └── jobs/             # Background jobs
├── migrations/            # Database migrations
├── docs/                 # Documentation
├── .env                  # Environment variables
├── go.mod               # Go modules
├── Makefile             # Build commands
└── README.md            # Project overview
```

## Komponen Utama

### 1. Handlers (`internal/handlers/`)
- **Tanggung jawab**: Memproses HTTP requests dan responses
- **Pattern**: Menggunakan validated body dari middleware
- **Response**: Centralized response format

**Contoh Handler:**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    validatedBody, _ := c.Get("validated_body")
    req := validatedBody.(*CreateUserRequest)
    
    result, err := h.userService.CreateUser(req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Operation failed", err.Error())
        return
    }
    
    response.Success(c, http.StatusCreated, "User created successfully", result)
}
```

### 2. Services (`internal/service/`)
- **Tanggung jawab**: Business logic dan orchestration
- **Pattern**: Tidak tahu tentang HTTP, hanya domain logic
- **Dependencies**: Repository layer untuk data access

### 3. Repository (`internal/repository/`)
- **Tanggung jawab**: Data access dengan raw SQL
- **Pattern**: Prepared statements untuk security
- **Database**: PostgreSQL dengan lib/pq driver

**Contoh Repository:**
```go
func (r *UserRepository) Create(user *models.User) error {
    query := `
        INSERT INTO users (name, email, user_identity, password_hash, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRow(query, 
        user.Name, user.Email, user.UserIdentity, user.PasswordHash, user.IsActive,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    return err
}
```

### 4. Models (`internal/models/`)
- **Tanggung jawab**: Data structure definitions
- **Pattern**: Struct tags untuk JSON dan database mapping
- **Features**: TableName() method untuk custom table names

### 5. Validation (`internal/validation/`)
- **Tanggung jawab**: Request validation rules yang modular
- **Pattern**: Validation rules terpisah dari routes untuk reusability
- **Structure**: Grouped berdasarkan domain/module

**Contoh Validation:**
```go
// internal/validation/user_validation.go
var CreateUserValidation = middleware.ValidationRules{
    Body: &struct {
        Name         string  `json:"name" validate:"required,min=2,max=100"`
        Email        string  `json:"email" validate:"required,email,max=255"`
        UserIdentity *string `json:"user_identity" validate:"omitempty,min=3,max=50"`
        Password     string  `json:"password" validate:"omitempty,min=6,max=100"`
    }{},
}
```

### 6. Middleware
- **Auth**: Custom token validation dengan Redis lookup
- **CORS**: Cross-origin request handling  
- **Validation**: Request validation dengan struct tags
- **Rate Limiting**: IP-based rate limiting (100 req/min)

### 7. Database System
- **Auth**: JWT token validation
- **CORS**: Cross-origin request handling  
- **Validation**: Request validation dengan struct tags
- **Rate Limiting**: IP-based rate limiting (100 req/min)

### 6. Database System
- **Driver**: lib/pq (PostgreSQL)
- **Migrations**: File-based SQL migrations
- **Connection**: Connection pooling
- **Transactions**: Manual transaction management

## Fitur Utama

### 1. Authentication System
- **Login Methods**: user_identity (primary), email (fallback)
- **Token Management**: Simple 1:1 custom token approach (1 access + 1 refresh per user)
- **Storage**: Redis untuk token storage dengan TTL otomatis
- **Security**: bcrypt password hashing
- **Revocation**: Mudah revoke token karena disimpan di Redis

### 2. Authorization (RBAC)
- **Roles**: Company dan branch-specific roles
- **Modules**: Hierarchical module system
- **Permissions**: Read, write, delete per module
- **Subscription**: Module access berdasarkan subscription tier

### 3. Multi-Company System
- **Companies**: Multi-tenant support
- **Branches**: 4-level hierarchical structure (Pusat → Area → City → District)
- **Isolation**: Data isolation per company

### 4. Subscription Management
- **Plans**: Tiered subscription plans (Basic, Standard, Premium, Enterprise)
- **Billing**: Monthly/yearly billing cycles
- **Module Access**: Module availability berdasarkan subscription tier
- **Payment Tracking**: Payment status dan renewal management

## Custom Token System

Project ini menggunakan custom token system yang disimpan di Redis, bukan JWT:

### Keuntungan Custom Token:
- **Easy Revocation**: Token dapat di-revoke dengan mudah dari Redis
- **Centralized Control**: Semua token tersimpan terpusat di Redis
- **TTL Management**: Automatic expiration dengan Redis TTL
- **Session Limiting**: 1 user = 1 access token + 1 refresh token
- **Real-time Validation**: Token validation langsung dari Redis

### Token Structure:
```
Redis Keys:
- access:user:{user_id}   # Access token (TTL: 15 minutes)
- refresh:user:{user_id}  # Refresh token (TTL: 7 days)

Token Data:
{
  "token": "generated_token_string",
  "metadata": {
    "user_id": 123,
    "user_agent": "...",
    "ip": "...",
    "abilities": ["module1", "module2"],
    "expires_at": 1640995200
  }
}
```

### Token Flow:
1. **Login**: Generate access + refresh token, store di Redis
2. **API Request**: Validate access token dari Redis
3. **Refresh**: Generate new tokens, overwrite existing di Redis
4. **Logout**: Delete semua tokens user dari Redis

## API Response Format

Semua API menggunakan format response yang konsisten:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... },
  "timestamp": "2025-01-01T12:00:00Z"
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Detailed error message",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=huminor_rbac
DB_SSL_MODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Token System
TOKEN_SECRET=your-secret-key

# Server
PORT=8081
GIN_MODE=debug
```

## Development Workflow

### 1. Setup Development
```bash
# Clone dan setup
git clone <repository>
cd project
go mod download

# Setup database
createdb huminor_rbac
cp .env.example .env

# Run migrations
make migrate-up

# Start development server
air  # atau go run cmd/api/main.go
```

### 2. Menambah Fitur Baru
1. Buat migration file di `migrations/`
2. Buat model di `internal/models/`
3. Buat repository di `internal/repository/`
4. Buat service di `internal/service/`
5. Buat handler di `internal/handlers/`
6. Tambahkan routes di `internal/routes/`
7. Update dependency injection di `internal/server/`

### 3. Testing
```bash
# Unit tests
go test ./...

# API testing dengan Postman
# Import: docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json
# Environment: docs/ERP_RBAC_Environment_Module_Based.postman_environment.json
```

## Production Deployment

### Docker Compose
```bash
# Setup production environment
cp .env.production .env

# Deploy dengan Docker
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Manual Deployment
```bash
# Build binary
CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

# Run migrations
./migrate -all

# Start server
./server
```

## Monitoring & Maintenance

### Health Check
```bash
curl http://localhost:8081/health
```

### Database Maintenance
```bash
# Check migration status
make migrate-status

# Backup database
pg_dump -h localhost -U postgres huminor_rbac > backup.sql

# Monitor performance
psql -d huminor_rbac -c "SELECT * FROM pg_stat_activity;"
```

### Redis Monitoring
```bash
# Check Redis connection
redis-cli ping

# Monitor token usage
redis-cli keys "access:user:*" | wc -l
redis-cli keys "refresh:user:*" | wc -l
```

## Best Practices

### 1. Code Organization
- Gunakan modular architecture
- Pisahkan concerns dengan jelas
- Implementasi proper error handling
- Gunakan centralized response format

### 2. Database
- Selalu gunakan prepared statements
- Implementasi proper indexing
- Gunakan soft delete untuk user data
- Backup database secara regular

### 3. Security
- Validasi semua input
- Gunakan rate limiting
- Hash password dengan bcrypt
- Implementasi proper authentication

### 4. Performance
- Gunakan connection pooling
- Cache data di Redis
- Implementasi pagination
- Monitor query performance

---

**Catatan**: Project ini dirancang untuk scalability dan maintainability dengan arsitektur yang clean dan modular.