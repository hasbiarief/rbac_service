# Unit-Based RBAC System

## Overview

Sistem RBAC telah diperluas dengan menambahkan layer "Unit" di antara Branch dan Role untuk memberikan kontrol akses yang lebih granular. Unit merepresentasikan bagian-bagian kecil dalam sebuah branch seperti departemen, divisi, atau tim kerja.

## Hierarki Baru

```
Company → Branch → Unit → Role → User
```

### Sebelumnya:
```
Company → Branch → Role → User
```

### Sekarang:
```
Company → Branch → Unit → Role → User
```

## Struktur Database Baru

### 1. Tabel `units`
Menyimpan unit-unit dalam sebuah branch dengan struktur hierarkis.

```sql
CREATE TABLE units (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL REFERENCES branches(id),
    parent_id BIGINT REFERENCES units(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    level INTEGER DEFAULT 0,
    path TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);
```

**Fitur:**
- **Hierarkis**: Unit dapat memiliki parent-child relationship
- **Path Tracking**: Menggunakan materialized path untuk query hierarkis yang efisien
- **Level Tracking**: Melacak kedalaman hierarki
- **Branch Scoped**: Setiap unit terikat pada satu branch

### 2. Tabel `unit_roles`
Memetakan role yang tersedia di setiap unit.

```sql
CREATE TABLE unit_roles (
    id BIGSERIAL PRIMARY KEY,
    unit_id BIGINT NOT NULL REFERENCES units(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(unit_id, role_id)
);
```

**Fungsi:**
- Menentukan role apa saja yang tersedia di setiap unit
- Satu unit dapat memiliki multiple roles
- Satu role dapat digunakan di multiple units

### 3. Tabel `unit_role_modules`
Menyimpan permission khusus untuk kombinasi unit-role-module.

```sql
CREATE TABLE unit_role_modules (
    id BIGSERIAL PRIMARY KEY,
    unit_role_id BIGINT NOT NULL REFERENCES unit_roles(id),
    module_id BIGINT NOT NULL REFERENCES modules(id),
    can_read BOOLEAN DEFAULT false,
    can_write BOOLEAN DEFAULT false,
    can_delete BOOLEAN DEFAULT false,
    can_approve BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(unit_role_id, module_id)
);
```

**Keunggulan:**
- **Customizable Permissions**: Setiap unit dapat memiliki permission yang berbeda untuk role yang sama
- **Override Default**: Permission di unit dapat override permission default dari role
- **Granular Control**: 4 level permission (Read, Write, Delete, Approve) per module

### 4. Update Tabel `user_roles`
Menambahkan kolom `unit_id` untuk mendukung assignment user ke unit.

```sql
ALTER TABLE user_roles ADD COLUMN unit_id BIGINT REFERENCES units(id);
```

## Contoh Implementasi

### Struktur Hierarki Unit

```
PT. Cinta Sejati (Company)
├── Kantor Pusat (Branch)
│   ├── Human Resources (Unit)
│   │   ├── HR Admin (Sub-Unit)
│   │   ├── Recruitment (Sub-Unit)
│   │   └── Payroll (Sub-Unit)
│   ├── Finance & Accounting (Unit)
│   │   ├── Accounting (Sub-Unit)
│   │   └── Treasury (Sub-Unit)
│   ├── Information Technology (Unit)
│   └── Operations (Unit)
├── Cabang Jakarta (Branch)
│   ├── Sales (Unit)
│   ├── Customer Service (Unit)
│   └── Admin (Unit)
└── Cabang Surabaya (Branch)
    ├── Sales (Unit)
    ├── Operations (Unit)
    └── Admin (Unit)
```

### Skenario Permission

#### Skenario 1: Role yang Sama, Permission Berbeda
- **HR_ADMIN** di unit "HR Admin" → Full access ke semua HR modules
- **HR_ADMIN** di unit "Recruitment" → Limited access, hanya recruitment modules
- **HR_ADMIN** di unit "Payroll" → Limited access, hanya payroll modules

#### Skenario 2: Unit Inheritance
- **MANAGER** di unit "Human Resources" → Access ke semua sub-units (HR Admin, Recruitment, Payroll)
- **EMPLOYEE** di unit "HR Admin" → Hanya access ke modules yang relevan dengan HR Admin

#### Skenario 3: Cross-Unit Access
- **AUDITOR** → Read-only access ke semua units dalam branch
- **SUPER_ADMIN** → Full access ke semua units di semua branches

## Permission Resolution Logic

### 1. Hierarchy Resolution
```
1. Check unit_role_modules (most specific)
2. If not found, check role_modules (default role permission)
3. If not found, deny access
```

### 2. Unit Hierarchy Permission
```
1. Check permission at current unit level
2. If user has role at parent unit, inherit permissions
3. Apply most restrictive permission if conflict
```

### 3. Multiple Role Resolution
```
1. User can have multiple roles in different units
2. Union of all permissions (most permissive wins)
3. Approval permissions require explicit grant
```

## API Endpoints Baru

### Unit Management
```
GET    /api/v1/units                    # List units
POST   /api/v1/units                    # Create unit
GET    /api/v1/units/:id                # Get unit details
PUT    /api/v1/units/:id                # Update unit
DELETE /api/v1/units/:id                # Delete unit
GET    /api/v1/units/:id/hierarchy      # Get unit hierarchy
GET    /api/v1/branches/:id/units       # Get units by branch
```

### Unit Role Management
```
GET    /api/v1/unit-roles               # List unit roles
POST   /api/v1/unit-roles               # Create unit role
DELETE /api/v1/unit-roles/:id           # Delete unit role
GET    /api/v1/units/:id/roles          # Get roles by unit
GET    /api/v1/roles/:id/units          # Get units by role
```

### Unit Permission Management
```
GET    /api/v1/unit-role-modules                    # List unit role modules
POST   /api/v1/unit-role-modules                    # Create unit role module
PUT    /api/v1/unit-role-modules/:id                # Update unit role module
DELETE /api/v1/unit-role-modules/:id                # Delete unit role module
POST   /api/v1/unit-role-modules/bulk-update        # Bulk update permissions
GET    /api/v1/units/:id/permissions                # Get unit permission summary
POST   /api/v1/units/copy-permissions               # Copy permissions between units
```

### Enhanced User Role Management
```
GET    /api/v1/user-roles               # List user roles (with unit info)
POST   /api/v1/user-roles               # Create user role (with unit)
PUT    /api/v1/user-roles/:id           # Update user role (with unit)
GET    /api/v1/users/:id/permissions    # Get user permissions (unit-aware)
```

## Migration Strategy

### 1. Data Migration
```sql
-- Existing user_roles without unit_id remain valid
-- They represent branch-level role assignments
-- New assignments can specify unit_id for granular control
```

### 2. Backward Compatibility
- Existing role assignments tetap berfungsi (branch-level)
- API responses include unit information jika tersedia
- Permission checking supports both branch-level dan unit-level

### 3. Gradual Migration
1. **Phase 1**: Create units untuk existing branches
2. **Phase 2**: Migrate critical user roles ke unit-specific
3. **Phase 3**: Customize permissions per unit as needed

## Benefits

### 1. Granular Access Control
- Setiap unit dapat memiliki permission yang berbeda untuk role yang sama
- Flexibility dalam mengatur akses sesuai kebutuhan bisnis
- Reduced over-privileging

### 2. Organizational Alignment
- Struktur permission mengikuti struktur organisasi
- Easier management untuk HR dan IT Admin
- Clear separation of concerns

### 3. Scalability
- Support untuk organizational growth
- Easy to add new units dan roles
- Hierarchical permission inheritance

### 4. Audit & Compliance
- Clear audit trail per unit
- Easier compliance reporting
- Granular permission tracking

## Best Practices

### 1. Unit Design
- Align units dengan struktur organisasi aktual
- Gunakan naming convention yang konsisten
- Maintain reasonable hierarchy depth (max 3-4 levels)

### 2. Permission Management
- Start dengan default role permissions
- Customize hanya jika diperlukan
- Document permission customizations

### 3. User Assignment
- Assign users ke most specific unit yang sesuai
- Avoid multiple overlapping role assignments
- Regular review dan cleanup

### 4. Performance
- Use materialized path untuk hierarchical queries
- Index pada frequently queried columns
- Cache permission results untuk active users

## Security Considerations

### 1. Principle of Least Privilege
- Default ke minimum required permissions
- Explicit grant untuk sensitive operations
- Regular permission audits

### 2. Separation of Duties
- Different units untuk conflicting responsibilities
- Approval workflows across units
- No single point of failure

### 3. Data Isolation
- Unit-based data filtering
- Cross-unit access controls
- Audit logging per unit

---

**Implementation Status**: Ready for development  
**Migration Required**: Yes (add units, update user_roles)  
**Backward Compatible**: Yes  
**Performance Impact**: Minimal (with proper indexing)