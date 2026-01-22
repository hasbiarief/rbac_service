# Database Dumps & Seeders

This directory contains database dumps and seeder files for the RBAC Service project.

## 📁 Directory Structure

```
database/
├── dumps/           # Complete database dumps
│   ├── structure.sql           # Database structure only
│   ├── data.sql               # Data only
│   └── complete_YYYYMMDD_HHMMSS.sql  # Complete dump with timestamp
│
├── seeders/         # Seeder files for new projects
│   ├── template.sql           # Complete template (structure + data, no audit logs)
│   ├── users_seeder.sql       # Users table seeder
│   ├── companies_seeder.sql   # Companies table seeder
│   ├── modules_seeder.sql     # Modules table seeder (89+ modules)
│   └── *_seeder.sql          # Individual table seeders
│
└── README.md        # This file
```

## 🚀 Quick Start

### Generate Database Dumps
```bash
# Generate all dump and seeder files
make db-dump
```

### Seed New Database
```bash
# Create and seed new database
make db-create
make db-seed

# Or do it in one command
make db-seed-fresh
```

## 📊 Database Statistics

Current database contains:
- **488 role-module permissions** - Granular access control
- **240 plan-module mappings** - Subscription-based access
- **195 unit-role-module permissions** - Unit-specific permissions
- **128 modules** - Complete HR module system
- **18 branches** - Multi-branch hierarchy
- **17 units** - Department/team structure
- **12 roles** - RBAC roles
- **11 users** - Sample users with different roles
- **6 companies** - Multi-company support
- **3 subscription plans** - Basic, Pro, Enterprise

## 🔑 Default Login Credentials

After seeding, you can login with these accounts:

| Name | Email | User Identity | Role | Password |
|------|-------|---------------|------|----------|
| Uzumaki Naruto | naruto@company.com | 100000001 | Super Admin | password123 |
| Haruno Sakura | sakura@company.com | 100000002 | HR Admin | password123 |
| Sasuke Uchiha | sasuke@company.com | 100000003 | Recruiter | password123 |
| Yamanaka Ino | ino@company.com | 100000006 | Employee | password123 |
| Aburame Shino | shino@company.com | 100000009 | Test Copy | password123 |
| Hasbi Due | hasbi@company.com | 800000001 | Console Admin | password123 |

**Note:** All users can login with either email or user_identity

## 📦 What's Included in Template

The template seeder includes:

### 🏢 **Company Structure**
- 3 sample companies with branches
- Hierarchical branch structure
- Unit-based organization (HR, Finance, IT, Sales, etc.)

### 👥 **User Management**
- 11 sample users with different roles (Naruto characters theme)
- Proper role assignments at company/branch/unit levels
- Password hashing with bcrypt
- Both email and user_identity login supported

### 🎭 **RBAC System**
- 12 predefined roles (Super Admin, HR Manager, etc.)
- Granular permissions per module
- Unit-specific role permissions
- Role inheritance and overrides

### 📦 **Module System**
- 128+ modules across 12 categories:
  - Core HR / Master Data (17 modules)
  - Employee Self Service (6 modules)
  - Recruitment (6 modules)
  - Attendance & Time (6 modules)
  - Leave Management (6 modules)
  - Performance Management (6 modules)
  - Training & Development (6 modules)
  - Payroll & Compensation (10 modules)
  - Asset & Facility (6 modules)
  - Disciplinary & Relations (4 modules)
  - Offboarding & Exit (6 modules)
  - Reporting & Analytics (6 modules)
  - System & Security (6 modules)

### 💳 **Subscription System**
- 3 subscription plans (Basic, Pro, Enterprise)
- Module access control based on subscription tier
- Plan-module mappings for access control
- Sample active subscriptions for companies

### 🏭 **Unit-Based RBAC**
- Hierarchical unit structure within branches
- Unit-specific role assignments
- Granular permissions at unit level
- Department-based access control

## 🛠️ Usage Examples

### Restore Complete Database
```bash
# Method 1: Using seeder (recommended)
createdb my_new_rbac_db
psql -h localhost -U postgres -d my_new_rbac_db < database/seeders/template.sql

# Method 2: Using complete dump
createdb my_new_rbac_db
psql -h localhost -U postgres -d my_new_rbac_db < database/dumps/complete_YYYYMMDD_HHMMSS.sql
```

### Restore Structure Only
```bash
createdb my_new_rbac_db
psql -h localhost -U postgres -d my_new_rbac_db < database/dumps/structure.sql
```

### Restore Specific Tables
```bash
# Restore only users
psql -h localhost -U postgres -d my_db < database/seeders/users_seeder.sql

# Restore only modules
psql -h localhost -U postgres -d my_db < database/seeders/modules_seeder.sql
```

## ⚠️ Important Notes

### Circular Foreign Key Constraints
Some tables (branches, modules, units) have circular foreign key constraints due to hierarchical relationships. When restoring data-only dumps, you might need to:

```bash
# Disable triggers temporarily
psql -d my_db -c "SET session_replication_role = replica;"
psql -d my_db < database/dumps/data.sql
psql -d my_db -c "SET session_replication_role = DEFAULT;"
```

### Password Security
- All default passwords are `password123`
- **Change default passwords in production!**
- Passwords are hashed with bcrypt (cost 10)

### Environment Variables
Make sure your `.env` file has correct database settings:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASS=your_password
DB_NAME=your_database
```

## 🔄 Updating Dumps

When you make changes to the database structure or data:

```bash
# Regenerate all dumps and seeders
make db-dump
```

This will create new timestamped dumps and update the template seeder.

## 📚 Related Documentation

- [Quick Start Guide](../docs/QUICK_START.md) - Setup and development
- [API Overview](../docs/API_OVERVIEW.md) - API documentation
- [Engineer Rules](../docs/ENGINEER_RULES.md) - Development guidelines
- [Project Structure](../docs/PROJECT_STRUCTURE.md) - Architecture overview

---

**Generated by RBAC Service Database Dump Tool**  
**Made with ❤️ by Hasbi**