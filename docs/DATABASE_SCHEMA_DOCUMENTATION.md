# Database Schema Documentation - ERP RBAC API

## Overview

Database schema untuk ERP RBAC API menggunakan PostgreSQL dengan 6 migration files yang mengimplementasikan complete RBAC system dengan subscription management.

## Migration Files

### 1. 001_create_users_table.sql
**Purpose**: User management dan authentication system

**Tables Created:**
- `users` - Core user information dengan authentication

**Key Features:**
- User identity system untuk login
- Password hashing dengan bcrypt
- Soft delete dengan `is_active` flag
- Timestamp tracking

### 2. 002_create_companies_and_branches.sql
**Purpose**: Multi-company support dengan hierarchical branch structure

**Tables Created:**
- `companies` - Company master data
- `branches` - Hierarchical branch structure

**Key Features:**
- Multi-company isolation
- Unlimited level branch hierarchy
- Parent-child relationships
- Path-based hierarchy tracking

### 3. 003_create_roles_and_modules.sql
**Purpose**: Core RBAC system implementation

**Tables Created:**
- `roles` - Role definitions
- `modules` - Hierarchical module system
- `user_roles` - User-role assignments
- `role_modules` - Role-module permissions

**Key Features:**
- Granular permission system (read, write, delete)
- Company and branch-specific role assignments
- Hierarchical module structure
- Subscription tier integration

### 4. 004_seed_modules_data.sql
**Purpose**: Default module data seeding

**Data Seeded:**
- 59 modules across 8 categories
- Hierarchical module relationships
- Subscription tier assignments

**Module Categories:**
- Core HR / Master Data (12 modules)
- Employee Self Service (8 modules)
- Payroll & Compensation (10 modules)
- Time & Attendance (8 modules)
- Performance Management (6 modules)
- Learning & Development (5 modules)
- Recruitment (6 modules)
- Reports & Analytics (4 modules)

### 5. 005_create_subscription_system.sql
**Purpose**: Subscription management system

**Tables Created:**
- `subscription_plans` - Available subscription plans
- `subscriptions` - Company subscriptions

**Key Features:**
- Tiered subscription system
- Automatic renewal support
- Billing cycle management
- Module access control

### 6. 006_seed_initial_data.sql
**Purpose**: Sample data untuk testing

**Data Seeded:**
- Sample companies dan branches
- Test users dengan different roles
- Subscription assignments
- Role-module permissions

## Database Schema Relationships

```
users
├── user_roles (many-to-many via user_roles)
│   ├── roles
│   ├── companies
│   └── branches
└── audit_logs (one-to-many)

companies
├── branches (one-to-many, hierarchical)
├── subscriptions (one-to-many)
└── user_roles (one-to-many)

modules (hierarchical)
├── role_modules (many-to-many via role_modules)
└── subscription_plans (many-to-many)

roles
├── user_roles (one-to-many)
└── role_modules (one-to-many)

subscription_plans
└── subscriptions (one-to-many)
```

## Key Design Patterns

### 1. Hierarchical Data
- **Modules**: Parent-child relationships dengan unlimited levels
- **Branches**: Company organizational structure
- **Path-based queries**: Efficient hierarchy traversal

### 2. Multi-tenancy
- **Company isolation**: Data segregation per company
- **Branch-level access**: Granular access control
- **Subscription-based features**: Module access control

### 3. RBAC Implementation
- **Role-based access**: Users assigned to roles
- **Permission granularity**: Read, write, delete permissions
- **Context-aware**: Company and branch-specific assignments

### 4. Audit Trail
- **Comprehensive logging**: All user activities tracked
- **Metadata storage**: JSON-based flexible data storage
- **Performance optimized**: Indexed for fast queries

## Performance Optimizations

### Indexes
- Primary keys pada semua tables
- Foreign key indexes untuk joins
- Composite indexes untuk common queries
- Path indexes untuk hierarchical queries

### Query Patterns
- **Efficient joins**: Optimized relationship queries
- **Hierarchical queries**: CTE-based tree traversal
- **Pagination support**: Limit/offset patterns
- **Search optimization**: Text search indexes

## Security Features

### Data Protection
- **Password hashing**: bcrypt dengan proper salt
- **Soft deletes**: Data preservation dengan is_active flags
- **Audit logging**: Complete activity tracking
- **Input validation**: Database-level constraints

### Access Control
- **Row-level security**: Company data isolation
- **Permission checking**: Role-based access validation
- **Session management**: JWT token validation
- **Rate limiting**: API endpoint protection

## Migration Strategy

### Forward Migrations
```bash
# Run all migrations
make migrate-up

# Run specific migration
psql -d database_name -f migrations/001_create_users_table.sql
```

### Rollback Strategy
- Each migration includes rollback instructions
- Data preservation during schema changes
- Backup recommendations before major changes

### Version Control
- Sequential numbering (001, 002, etc.)
- Descriptive naming convention
- Change documentation in each file

## Data Types and Constraints

### Common Patterns
- **IDs**: `BIGSERIAL PRIMARY KEY`
- **Timestamps**: `TIMESTAMP DEFAULT CURRENT_TIMESTAMP`
- **Soft Delete**: `is_active BOOLEAN DEFAULT true`
- **JSON Data**: `JSONB` untuk flexible metadata

### Validation Rules
- **Email format**: Valid email constraints
- **Password strength**: Minimum length requirements
- **Unique constraints**: Prevent duplicate data
- **Foreign key constraints**: Data integrity

## Backup and Recovery

### Backup Strategy
```bash
# Full database backup
pg_dump -h localhost -U username -d database_name > backup.sql

# Schema only backup
pg_dump -h localhost -U username -d database_name --schema-only > schema.sql
```

### Recovery Procedures
- Point-in-time recovery support
- Transaction log backup
- Regular backup scheduling
- Disaster recovery planning

## Monitoring and Maintenance

### Performance Monitoring
- Query performance analysis
- Index usage statistics
- Connection pool monitoring
- Slow query identification

### Maintenance Tasks
- Regular VACUUM and ANALYZE
- Index maintenance
- Statistics updates
- Log rotation

## Future Considerations

### Scalability
- **Horizontal scaling**: Read replicas support
- **Partitioning**: Large table partitioning strategy
- **Caching**: Redis integration for performance
- **Connection pooling**: Optimized connection management

### Feature Extensions
- **Multi-language support**: Internationalization ready
- **Custom fields**: Flexible schema extensions
- **Integration APIs**: External system connectivity
- **Advanced reporting**: Analytics and BI support

## Conclusion

Database schema dirancang untuk:
- **Scalability**: Handle enterprise-level data
- **Security**: Comprehensive access control
- **Performance**: Optimized for common operations
- **Flexibility**: Easy to extend and modify
- **Reliability**: ACID compliance dan data integrity

Schema ini mendukung complete ERP RBAC system dengan subscription management yang siap untuk production deployment.