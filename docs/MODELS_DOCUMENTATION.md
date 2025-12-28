# Models Documentation - ERP RBAC API

## Gambaran Umum

Sekarang semua model telah didefinisikan dengan lengkap dan terorganisir berdasarkan domain. Setiap file model berisi struct yang berkaitan dengan domain tersebut beserta helper types dan methods.

## Struktur Models

```
internal/models/
├── audit.go        # Audit logging models
├── branch.go       # Branch management models  
├── company.go      # Company management models
├── module.go       # Module system models
├── role.go         # Role and permission models
├── subscription.go # Subscription system models
└── user.go         # User management models
```

## Detail Models

### 1. User Models (`user.go`)
```go
type User struct {
    ID           int64     `json:"id" db:"id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    UserIdentity *string   `json:"user_identity" db:"user_identity"`
    Password     string    `json:"-" db:"password"`
    IsActive     bool      `json:"is_active" db:"is_active"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `UserWithRoles` - User dengan informasi roles
- `UserWithModules` - User dengan informasi modules yang dapat diakses

### 2. Company Models (`company.go`)
```go
type Company struct {
    ID        int64     `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Code      string    `json:"code" db:"code"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `CompanyWithStats` - Company dengan statistik users dan branches

### 3. Branch Models (`branch.go`)
```go
type Branch struct {
    ID        int64     `json:"id" db:"id"`
    CompanyID int64     `json:"company_id" db:"company_id"`
    ParentID  *int64    `json:"parent_id" db:"parent_id"`
    Name      string    `json:"name" db:"name"`
    Code      string    `json:"code" db:"code"`
    Level     int       `json:"level" db:"level"`
    Path      string    `json:"path" db:"path"`
    IsActive  bool      `json:"is_active" db:"is_active"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `BranchWithCompany` - Branch dengan informasi company
- `BranchHierarchy` - Branch dengan struktur hierarki children
- `BranchWithStats` - Branch dengan statistik users dan sub-branches

### 4. Role Models (`role.go`)
```go
type Role struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole struct {
    ID        int64     `json:"id" db:"id"`
    UserID    int64     `json:"user_id" db:"user_id"`
    RoleID    int64     `json:"role_id" db:"role_id"`
    CompanyID int64     `json:"company_id" db:"company_id"`
    BranchID  *int64    `json:"branch_id" db:"branch_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RoleModule struct {
    ID        int64     `json:"id" db:"id"`
    RoleID    int64     `json:"role_id" db:"role_id"`
    ModuleID  int64     `json:"module_id" db:"module_id"`
    CanRead   bool      `json:"can_read" db:"can_read"`
    CanWrite  bool      `json:"can_write" db:"can_write"`
    CanDelete bool      `json:"can_delete" db:"can_delete"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `RoleWithPermissions` - Role dengan informasi modules dan permissions
- `RoleModulePermission` - Permission detail untuk role-module mapping

### 5. Module Models (`module.go`)
```go
type Module struct {
    ID               int64     `json:"id" db:"id"`
    Category         string    `json:"category" db:"category"`
    Name             string    `json:"name" db:"name"`
    URL              string    `json:"url" db:"url"`
    Icon             string    `json:"icon" db:"icon"`
    Description      string    `json:"description" db:"description"`
    ParentID         *int64    `json:"parent_id" db:"parent_id"`
    SubscriptionTier string    `json:"subscription_tier" db:"subscription_tier"`
    IsActive         bool      `json:"is_active" db:"is_active"`
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `ModuleWithChildren` - Module dengan struktur hierarki children
- `UserModule` - Module dengan permission information untuk user

### 6. Subscription Models (`subscription.go`)
```go
type SubscriptionPlan struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    DisplayName string    `json:"display_name" db:"display_name"`
    Description string    `json:"description" db:"description"`
    Price       float64   `json:"price" db:"price"`
    Currency    string    `json:"currency" db:"currency"`
    Duration    int       `json:"duration" db:"duration"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Subscription struct {
    ID           int64     `json:"id" db:"id"`
    CompanyID    int64     `json:"company_id" db:"company_id"`
    PlanID       int64     `json:"plan_id" db:"plan_id"`
    Status       string    `json:"status" db:"status"`
    StartDate    time.Time `json:"start_date" db:"start_date"`
    EndDate      time.Time `json:"end_date" db:"end_date"`
    BillingCycle string    `json:"billing_cycle" db:"billing_cycle"`
    AutoRenew    bool      `json:"auto_renew" db:"auto_renew"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
```

**Helper Types:**
- `SubscriptionWithPlan` - Subscription dengan informasi plan
- `SubscriptionWithCompany` - Subscription dengan informasi company

### 7. Audit Models (`audit.go`)
```go
type AuditLog struct {
    ID           int64                  `json:"id" db:"id"`
    UserID       *int64                 `json:"user_id" db:"user_id"`
    UserIdentity *string                `json:"user_identity" db:"user_identity"`
    Action       string                 `json:"action" db:"action"`
    Resource     string                 `json:"resource" db:"resource"`
    ResourceID   *string                `json:"resource_id" db:"resource_id"`
    Method       string                 `json:"method" db:"method"`
    URL          string                 `json:"url" db:"url"`
    UserAgent    *string                `json:"user_agent" db:"user_agent"`
    IP           *string                `json:"ip" db:"ip"`
    Status       string                 `json:"status" db:"status"`
    StatusCode   int                    `json:"status_code" db:"status_code"`
    Message      string                 `json:"message" db:"message"`
    Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
    CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}
```

**Helper Types:**
- `AuditLogWithUser` - Audit log dengan informasi user
- `AuditStats` - Statistik audit logs
- `ActionCount`, `UserActivityCount`, `HourlyActivity`, `StatusCount` - Helper types untuk statistik
- `JSONB` - Custom type untuk PostgreSQL JSON fields

## Fitur Model

### 1. **TableName Methods**
Setiap model utama memiliki method `TableName()` untuk mapping ke database table:
```go
func (User) TableName() string {
    return "users"
}
```

### 2. **JSON Tags**
Semua field memiliki JSON tags untuk API response:
```go
Name string `json:"name" db:"name"`
```

### 3. **Database Tags**
Field memiliki database tags untuk SQL mapping:
```go
Name string `json:"name" db:"name"`
```

### 4. **Pointer Fields**
Field yang nullable menggunakan pointer:
```go
UserIdentity *string `json:"user_identity" db:"user_identity"`
ParentID     *int64  `json:"parent_id" db:"parent_id"`
```

### 5. **Time Handling**
Semua timestamp menggunakan `time.Time`:
```go
CreatedAt time.Time `json:"created_at" db:"created_at"`
UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
```

### 6. **Password Security**
Password field menggunakan tag `json:"-"` untuk tidak di-serialize:
```go
Password string `json:"-" db:"password"`
```

### 7. **Custom Types**
- `JSONB` type untuk PostgreSQL JSON fields dengan custom `Value()` dan `Scan()` methods
- Helper structs untuk complex queries dan responses

## Konsistensi Model

### 1. **Naming Convention**
- Struct names menggunakan PascalCase
- Field names menggunakan PascalCase
- JSON tags menggunakan snake_case
- Database tags menggunakan snake_case

### 2. **Standard Fields**
Semua model utama memiliki:
- `ID int64` - Primary key
- `CreatedAt time.Time` - Creation timestamp
- `UpdatedAt time.Time` - Update timestamp (jika applicable)
- `IsActive bool` - Soft delete flag (jika applicable)

### 3. **Foreign Key Convention**
- `CompanyID int64` untuk reference ke companies
- `UserID int64` untuk reference ke users
- `ParentID *int64` untuk self-referencing hierarchy

## Keuntungan Struktur Model Modular

### 1. **Separation of Concerns**
- Setiap domain memiliki model terpisah
- Mudah maintenance dan development
- Clear ownership per domain

### 2. **Type Safety**
- Strong typing untuk semua fields
- Compile-time error detection
- Clear data contracts

### 3. **Database Mapping**
- Consistent database mapping
- Support untuk complex queries
- Efficient data serialization

### 4. **API Response**
- Clean JSON serialization
- Consistent response format
- Support untuk nested data

### 5. **Extensibility**
- Mudah menambah model baru
- Helper types untuk complex scenarios
- Support untuk future requirements

## Kesimpulan

Struktur model sekarang sudah lengkap dan terorganisir dengan baik:

- ✅ **7 Domain Models** - User, Company, Branch, Role, Module, Subscription, Audit
- ✅ **Helper Types** - 15+ helper structs untuk complex scenarios
- ✅ **Type Safety** - Strong typing dengan proper nullable fields
- ✅ **Database Integration** - Proper mapping dengan TableName methods
- ✅ **API Ready** - Clean JSON serialization
- ✅ **PostgreSQL Support** - Custom JSONB type untuk JSON fields
- ✅ **Security** - Password field protection
- ✅ **Consistency** - Standard naming dan field conventions

Model sekarang siap untuk mendukung semua fitur ERP RBAC system dengan struktur yang maintainable dan scalable.