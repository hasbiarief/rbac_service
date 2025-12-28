# Database Migrations Guide - ERP RBAC API

## Overview

Sistem migration menggunakan file-based approach dengan PostgreSQL untuk mengelola database schema evolution. Setiap migration file berisi SQL commands untuk forward dan rollback operations.

## Migration System Architecture

### File Structure
```
migrations/
├── 001_create_users_table.sql
├── 002_create_companies_and_branches.sql
├── 003_create_roles_and_modules.sql
├── 004_seed_modules_data.sql
├── 005_create_subscription_system.sql
└── 006_seed_initial_data.sql
```

### Migration Tools
- **Migration Runner**: `cmd/migrate/main.go`
- **SQL Migrate Tool**: `cmd/sql-migrate/main.go`
- **Package**: `pkg/migration/migration.go`

## Migration Files Detail

### 1. 001_create_users_table.sql
**Purpose**: Core user management system

**Creates:**
- `users` table dengan authentication fields
- Indexes untuk performance
- Constraints untuk data integrity

**Key Features:**
- User identity system (9-digit format)
- Password hashing support
- Soft delete capability
- Audit timestamps

### 2. 002_create_companies_and_branches.sql
**Purpose**: Multi-company dan hierarchical branch structure

**Creates:**
- `companies` table untuk multi-tenancy
- `branches` table dengan hierarchical support
- Foreign key relationships
- Path-based hierarchy tracking

**Key Features:**
- Unlimited branch hierarchy levels
- Company isolation
- Efficient tree traversal

### 3. 003_create_roles_and_modules.sql
**Purpose**: Core RBAC system

**Creates:**
- `roles` table untuk role definitions
- `modules` table dengan hierarchical structure
- `user_roles` table untuk user-role assignments
- `role_modules` table untuk permission mapping

**Key Features:**
- Granular permissions (read, write, delete)
- Company and branch-specific assignments
- Module hierarchy support

### 4. 004_seed_modules_data.sql
**Purpose**: Default module data

**Seeds:**
- 59 modules across 8 categories
- Hierarchical module relationships
- Subscription tier assignments

**Module Categories:**
- Core HR / Master Data
- Employee Self Service
- Payroll & Compensation
- Time & Attendance
- Performance Management
- Learning & Development
- Recruitment
- Reports & Analytics

### 5. 005_create_subscription_system.sql
**Purpose**: Subscription management

**Creates:**
- `subscription_plans` table
- `subscriptions` table
- Plan-module relationships

**Key Features:**
- Tiered subscription system
- Automatic renewal support
- Module access control

### 6. 006_seed_initial_data.sql
**Purpose**: Sample data untuk testing

**Seeds:**
- Sample companies dan branches
- Test users dengan different roles
- Subscription assignments
- Role-module permissions

## Running Migrations

### Using Make Commands
```bash
# Run all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Reset database (drop all tables)
make migrate-reset

# Check migration status
make migrate-status
```

### Using Migration Tool
```bash
# Build migration tool
go build -o migrate cmd/migrate/main.go

# Run specific migration
./migrate -file migrations/001_create_users_table.sql

# Run all migrations
./migrate -all
```

### Manual Execution
```bash
# Run specific migration file
psql -d database_name -f migrations/001_create_users_table.sql

# Run all migrations in order
for file in migrations/*.sql; do
    psql -d database_name -f "$file"
done
```

## Migration Best Practices

### 1. **Naming Convention**
- Use sequential numbering: `001`, `002`, etc.
- Descriptive names: `create_users_table`, `add_audit_columns`
- Consistent format: `{number}_{description}.sql`

### 2. **File Structure**
```sql
-- Migration: 001_create_users_table.sql
-- Description: Create users table with authentication support
-- Author: Developer Name
-- Date: 2024-01-01

-- Forward Migration
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    -- table definition
);

-- Rollback Migration (commented)
-- DROP TABLE IF EXISTS users;
```

### 3. **Safety Guidelines**
- **Always backup** before running migrations
- **Test migrations** on development environment first
- **Include rollback** instructions in comments
- **Avoid destructive operations** in production
- **Use transactions** for complex migrations

### 4. **Data Integrity**
- Add constraints untuk data validation
- Create indexes untuk performance
- Maintain foreign key relationships
- Handle existing data carefully

## Migration Commands Reference

### Environment Setup
```bash
# Set database connection
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=huminor_rbac
```

### Development Workflow
```bash
# 1. Create new migration
touch migrations/007_new_feature.sql

# 2. Write migration SQL
vim migrations/007_new_feature.sql

# 3. Test migration
make migrate-up

# 4. Verify changes
psql -d huminor_rbac -c "\dt"

# 5. Test rollback (if needed)
make migrate-down
```

### Production Deployment
```bash
# 1. Backup database
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Run migrations
make migrate-up

# 3. Verify deployment
make migrate-status

# 4. Test application
curl http://localhost:8081/health
```

## Troubleshooting

### Common Issues

#### 1. **Migration Failed**
```bash
# Check current state
psql -d database_name -c "SELECT * FROM schema_migrations;"

# Manual rollback
psql -d database_name -f rollback_commands.sql

# Fix and retry
make migrate-up
```

#### 2. **Duplicate Key Errors**
```sql
-- Add IF NOT EXISTS clauses
CREATE TABLE IF NOT EXISTS users (...);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

#### 3. **Foreign Key Violations**
```sql
-- Disable foreign key checks temporarily
SET foreign_key_checks = 0;
-- Run migration
SET foreign_key_checks = 1;
```

#### 4. **Large Data Migrations**
```sql
-- Use batched updates
UPDATE users SET status = 'active' 
WHERE id BETWEEN 1 AND 1000;

-- Use LIMIT for large datasets
DELETE FROM audit_logs 
WHERE created_at < '2023-01-01' 
LIMIT 1000;
```

### Recovery Procedures

#### 1. **Rollback Strategy**
```bash
# Identify last successful migration
make migrate-status

# Rollback to specific version
psql -d database_name -f rollback_to_version_005.sql

# Verify rollback
make migrate-status
```

#### 2. **Data Recovery**
```bash
# Restore from backup
pg_restore -h $DB_HOST -U $DB_USER -d $DB_NAME backup_file.sql

# Partial restore
pg_restore -h $DB_HOST -U $DB_USER -d $DB_NAME -t users backup_file.sql
```

## Advanced Migration Patterns

### 1. **Schema Changes**
```sql
-- Add column with default value
ALTER TABLE users ADD COLUMN phone VARCHAR(20) DEFAULT '';

-- Modify column type
ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(50);

-- Drop column (careful!)
ALTER TABLE users DROP COLUMN IF EXISTS old_column;
```

### 2. **Data Transformations**
```sql
-- Migrate data format
UPDATE users 
SET user_identity = LPAD(user_identity::text, 9, '0')
WHERE LENGTH(user_identity::text) < 9;

-- Split data into new table
INSERT INTO user_profiles (user_id, profile_data)
SELECT id, jsonb_build_object('name', name, 'email', email)
FROM users;
```

### 3. **Index Management**
```sql
-- Create index concurrently (PostgreSQL)
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- Drop index if exists
DROP INDEX IF EXISTS old_index_name;

-- Rebuild index
REINDEX INDEX idx_users_email;
```

## Monitoring and Maintenance

### 1. **Migration Tracking**
```sql
-- Check migration history
SELECT * FROM schema_migrations ORDER BY version;

-- Verify table structure
\d+ users

-- Check constraints
\d+ users
```

### 2. **Performance Monitoring**
```sql
-- Check query performance
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';

-- Monitor index usage
SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
FROM pg_stat_user_indexes;
```

### 3. **Maintenance Tasks**
```sql
-- Update table statistics
ANALYZE users;

-- Vacuum tables
VACUUM ANALYZE users;

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables WHERE schemaname = 'public';
```

## Conclusion

Migration system dirancang untuk:
- **Safety**: Backup dan rollback support
- **Reliability**: Transaction-based execution
- **Maintainability**: Clear naming dan documentation
- **Performance**: Optimized schema changes
- **Scalability**: Support untuk large datasets

Ikuti best practices untuk memastikan smooth database evolution dan minimal downtime pada production environment.