# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.5] - 2025-12-25

### Added
- **Documentation Reorganization Complete**
  - Consolidated all documentation in `docs/` folder
  - Created comprehensive documentation index in `docs/README.md`
  - Added dedicated documentation for each system component
  - Updated main `README.md` with complete system overview and monetization focus

### Documentation Structure
- **System Architecture**
  - `docs/MODULE_SYSTEM_DOCUMENTATION.md` - RBAC module system
  - `docs/DATABASE_SCHEMA_DOCUMENTATION.md` - Complete database schema
  - `docs/BRANCH_HIERARCHY_API.md` - Branch hierarchy system

- **Security & Authentication**
  - `docs/PASSWORD_SYSTEM.md` - Bcrypt password management
  - `docs/CORS_INTEGRATION.md` - Cross-origin resource sharing

- **Subscription & Monetization**
  - `docs/SUBSCRIPTION_SYSTEM.md` - Tiered pricing system
  - `docs/MODULE_LOCKING_SYSTEM.md` - Subscription-based access control

- **API & Integration**
  - `docs/API_DOCUMENTATION_MODULE_BASED.md` - Complete API reference
  - `docs/FRONTEND_INTEGRATION_GUIDE.md` - Next.js integration guide
  - `docs/POSTMAN_DOCUMENTATION_SUMMARY.md` - Testing guide

### Removed
- **Obsolete Documentation Files**
  - `IMPLEMENTATION_SUMMARY.md` - Consolidated into docs folder
  - `MODULE_LOCKING_IMPLEMENTATION_SUMMARY.md` - Moved to `docs/MODULE_LOCKING_SYSTEM.md`
  - `SUBSCRIPTION_SYSTEM_IMPLEMENTATION_SUMMARY.md` - Moved to `docs/SUBSCRIPTION_SYSTEM.md`
  - `PASSWORD_SYSTEM_IMPLEMENTATION_SUMMARY.md` - Moved to `docs/PASSWORD_SYSTEM.md`
  - `CORS_IMPLEMENTATION_SUMMARY.md` - Moved to `docs/CORS_INTEGRATION.md`

### Improved
- **Documentation Quality**
  - Standardized documentation format across all files
  - Added comprehensive system overview in main docs README
  - Improved navigation with clear categorization
  - Added production deployment guides
  - Enhanced testing documentation
  - Updated main README.md with complete monetization focus

### Technical Improvements
- **Module Locking System**: Production-ready subscription-based access control working correctly
- **Documentation Structure**: Clean, organized, and maintainable documentation
- **System Status**: All components documented and production-ready
- **Business Focus**: Clear monetization strategy and revenue model documented
- **Test Data**: Updated company 1 to Basic Plan for proper module locking demonstration

### Database Updates
- **Migration 013**: Updated test company (PT. Cinta Sejati) to Basic Plan for module locking testing
- **Module Access Control**: User 100000003 now properly demonstrates Basic Plan module restrictions

### Business Value
- **Revenue Model**: Clear subscription tiers with pricing (Basic: Rp 990K/year, Pro: Rp 2.99M/year, Enterprise: Rp 5.99M/year)
- **Module Access Control**: Automatic module locking based on subscription plans
- **Production Ready**: Complete system ready for deployment and monetization
- **Market Positioning**: Clear value proposition for SaaS providers and enterprises

## [2.0.4] - 2025-12-25

### Added
- **CORS Middleware for Frontend Integration**
  - Complete CORS configuration for Next.js frontend support
  - Preflight OPTIONS request handling
  - Configurable allowed origins via environment variables
  - Development and production CORS modes

### Features
- **CORS Configuration**
  - Allowed origins: localhost:3000, localhost:3001, 127.0.0.1:3000, 127.0.0.1:3001
  - Allowed methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
  - Credentials support: `Access-Control-Allow-Credentials: true`
  - 24-hour preflight cache: `Access-Control-Max-Age: 86400`

- **Environment Configuration**
  - `CORS_ORIGINS` environment variable for flexible origin configuration
  - `ENVIRONMENT` variable for development/production mode switching
  - Automatic origin validation and security

- **Frontend Integration Ready**
  - Complete Next.js integration guide with examples
  - React components for authentication and module access
  - API client with automatic token refresh
  - Protected route components with module-based access control

### Documentation
- **Frontend Integration Guide**: Complete guide for Next.js integration
- **API Documentation**: Updated with CORS configuration details
- **React Examples**: Login forms, API clients, protected routes
- **Environment Setup**: Development and production configuration

### Technical Implementation
- **Middleware**: `middleware/cors.go` with flexible configuration
- **Config Updates**: Added CORS_ORIGINS and ENVIRONMENT to config
- **Router Integration**: CORS middleware integrated in router setup
- **Testing**: Comprehensive CORS testing with curl and browser simulation

### Fixed
- **Frontend Access**: Resolved CORS policy blocking frontend API calls
- **Browser Security**: Proper handling of browser preflight requests
- **Credentials**: Fixed credential passing for authenticated requests

## [2.0.3] - 2025-12-25

### Added
- **Subscription Management System**
  - Tabel `subscription_plans` dengan 3 paket (Basic, Pro, Enterprise)
  - Tabel `subscriptions` untuk manajemen langganan perusahaan
  - Tabel `plan_modules` untuk mapping paket ke modul
  - Module subscription tiers (basic, pro, enterprise)
  - Monthly dan yearly billing cycles dengan auto-renewal
  - Subscription status tracking (active, expired, cancelled, suspended, trial)

### Features
- **Subscription Plans**
  - Basic Plan: Rp 99.000/bulan, Rp 990.000/tahun (25 users, 3 branches)
  - Professional Plan: Rp 299.000/bulan, Rp 2.990.000/tahun (100 users, 10 branches)
  - Enterprise Plan: Rp 599.000/bulan, Rp 5.990.000/tahun (unlimited)
  - Automatic discount calculation (16.7% yearly vs monthly)

- **Module Access Control**
  - Modules locked/unlocked based on subscription tier
  - Real-time access checking per company per module
  - Graceful degradation for expired subscriptions

- **Subscription Management**
  - Create, update, renew, cancel subscriptions
  - Company subscription status with module breakdown
  - Subscription statistics and analytics
  - Expiring subscription alerts

### API Endpoints
- `GET /api/v1/subscription/plans` - Get subscription plans (public)
- `POST /api/v1/subscription/subscriptions` - Create subscription
- `GET /api/v1/subscription/companies/{id}/status` - Company subscription status
- `GET /api/v1/subscription/companies/{id}/modules/{id}/access` - Check module access
- `GET /api/v1/subscription/stats` - Subscription statistics
- `POST /api/v1/subscription/subscriptions/{id}/renew` - Renew subscription
- `POST /api/v1/subscription/subscriptions/{id}/cancel` - Cancel subscription

### Database Schema
- **Migration 012**: Complete subscription system tables
- **Module tiers**: Added subscription_tier column to modules table
- **Plan features**: JSONB storage for flexible plan features
- **Billing cycles**: Monthly/yearly with proper date calculations
- **Payment tracking**: Payment status and date management

### Documentation
- **API Documentation**: Complete subscription endpoints documentation
- **Database Documentation**: Subscription tables with query examples
- **Subscription tiers**: Module access control explanation
- **Postman Collection**: Subscription Management section (14 endpoints)
  - Public subscription plans endpoints
  - Complete subscription CRUD operations
  - Company subscription status and module access checking
  - Subscription analytics and expiring alerts
  - Automated test scripts for subscription workflows

## [2.0.2] - 2025-12-25

### Added
- **Password Management System**
  - Bcrypt password hashing dengan cost 12 untuk keamanan optimal
  - Route untuk ubah password: `PUT /api/v1/users/{id}/password` (admin)
  - Route untuk ubah password sendiri: `PUT /api/v1/users/me/password` (user)
  - Password validation dengan minimum 6 karakter, maksimum 100 karakter
  - Default password "password123" untuk user baru (ter-hash dengan bcrypt)

### Changed
- **Database Schema**
  - Kolom `password_hash` sekarang menyimpan bcrypt hash (cost 12) bukan plain text
  - Migration 011: Update semua password existing ke bcrypt hash
  - Comment database diperbaharui untuk menjelaskan bcrypt hashing

### Security
- **Password Security**
  - Implementasi bcrypt hashing untuk semua password
  - Password verification menggunakan bcrypt.CompareHashAndPassword
  - Validasi password strength di level aplikasi
  - Audit logging untuk password changes

### Documentation
- **API Documentation**
  - Menambahkan dokumentasi lengkap Password Management endpoints
  - Menambahkan dokumentasi User Management CRUD operations
  - Contoh request/response untuk semua password operations
  - Error handling documentation untuk password operations

- **Database Documentation**
  - Update dokumentasi password_hash column dengan bcrypt information
  - Menambahkan query patterns untuk password management
  - Security best practices untuk password handling

- **Postman Collection**
  - Menambahkan User Management section (6 endpoints)
  - Menambahkan Password Management section (6 endpoints)
  - Test scenarios untuk password validation dan error handling
  - Environment variables untuk user management testing
  - Automated token management dan test scripts

## [2.0.1] - 2025-12-25

### Added
- **Database Schema Documentation** (`docs/DATABASE_SCHEMA_DOCUMENTATION.md`)
  - Dokumentasi lengkap struktur tabel dan atribut untuk semua 9 tabel utama
  - Penjelasan detail setiap kolom dengan constraint, index, dan relationship
  - Query patterns yang sering digunakan untuk authentication, authorization, dan audit
  - Views (`user_module_access`, `module_hierarchy`) untuk mempermudah query kompleks
  - Best practices untuk maintenance, performance, dan security
  - Panduan lengkap untuk developer dan DBA

### Documentation
- Menambahkan dokumentasi komprehensif untuk semua tabel database sistem
- Penjelasan cara membaca dan menggunakan setiap tabel dengan contoh query
- Referensi lengkap struktur database module-based RBAC system
- Panduan maintenance dan monitoring database performance

## [2.0.0] - 2025-12-24

### 🚀 MAJOR RELEASE: Module-Based RBAC System

**BREAKING CHANGES**: Complete system redesign from permission-based to module-based access control with user_identity authentication.

### Added
- **Module-Based Access Control System**
  - 89+ modules across 12 categories (Core HR, Recruitment, Attendance, etc.)
  - Hierarchical module structure with parent-child relationships
  - Granular permissions per module (read, write, delete)
  - Dynamic role-module assignments

- **user_identity Authentication**
  - Login using 9-digit user_identity instead of email
  - JWT-based authentication with 15-minute access tokens
  - 7-day refresh tokens with family ID tracking
  - User roles and modules loaded during authentication

- **Comprehensive Module Categories**
  1. Core HR / Master Data (17 modules)
  2. Recruitment (6 modules)
  3. Attendance & Time (6 modules)
  4. Leave Management (6 modules)
  5. Payroll & Compensation (10 modules)
  6. Performance Management (6 modules)
  7. Training & Development (6 modules)
  8. Employee Self Service (6 modules)
  9. Asset & Facility (6 modules)
  10. Disciplinary & Relations (4 modules)
  11. Offboarding & Exit (6 modules)
  12. Reporting & Analytics (6 modules)

- **Dynamic Role Management System**
  - Assign/remove user roles with bulk operations
  - Update role-module mappings dynamically
  - User access summaries with complete module breakdown
  - Role-based access control with company/branch scoping

- **Comprehensive Audit Logging**
  - Track all system activities with user identity
  - Request/response logging with performance metrics
  - Audit log retrieval with filtering capabilities
  - User-specific audit trails

- **Hierarchical Branch Structure**
  - Unlimited level branch hierarchy (pusat → cabang → sub-cabang)
  - Automatic path generation with database triggers
  - Parent-child relationships with proper constraints

### Changed
- **Authentication Method**: Email → user_identity (9 digits)
- **Access Control**: Permission-based → Module-based hierarchical system
- **Database Schema**: New tables (modules, role_modules, audit_logs)
- **API Endpoints**: New module-based endpoints replacing permission endpoints
- **Test Users**: Updated to use user_identity format

### Removed
- **Old Permission System**: Removed permissions and role_permissions tables
- **Email Authentication**: No longer supported
- **Old API Endpoints**: Permission-based endpoints removed
- **Legacy Documentation**: Old documentation files cleaned up

### Database Migrations
- `007_create_module_system.sql` - Module-based RBAC schema
- `008_seed_modules_data.sql` - 89+ modules across 12 categories
- `009_seed_dummy_data.sql` - Test users with user_identity
- `010_create_audit_logs.sql` - Comprehensive audit logging system

### API Endpoints (New)
```
# Authentication (user_identity)
POST   /api/v1/auth/login          # Login with user_identity
POST   /api/v1/auth/refresh        # Refresh JWT token
POST   /api/v1/auth/logout         # Logout

# Module System
GET    /api/v1/modules             # Get modules with filtering
POST   /api/v1/modules             # Create new module
GET    /api/v1/modules/tree        # Get hierarchical module tree
GET    /api/v1/modules/:id         # Get specific module

# User Access
GET    /api/v1/users/identity/:identity/modules  # Get user's modules
POST   /api/v1/users/check-access  # Check specific access

# Role Management
POST   /api/v1/role-management/assign-user-role     # Assign role to user
POST   /api/v1/role-management/bulk-assign-roles    # Bulk role assignments
GET    /api/v1/role-management/user/:id/access-summary  # User access summary

# Audit Logging
GET    /api/v1/audit/logs          # Get audit logs with filtering
GET    /api/v1/audit/stats         # Get audit statistics
```

### Documentation (Updated)
- **Complete Rewrite**: All documentation updated for module-based system
- **New API Documentation**: `API_DOCUMENTATION_MODULE_BASED.md`
- **New Quick Start**: `QUICK_START_MODULE_BASED.md`
- **Module System Guide**: `MODULE_SYSTEM_DOCUMENTATION.md`
- **Updated Postman Collection**: Module-based API collection
- **Cleanup Summary**: Documentation cleanup and migration guide

### Test Users (Updated)
| user_identity | Role | Password | Access Level |
|---------------|------|----------|--------------|
| 100000003 | BRANCH_ADMIN | password123 | 41 modules across multiple categories |
| 100000004 | HR_MANAGER | password123 | 73 modules (HR focused) |
| 100000006 | EMPLOYEE | password123 | 41 modules (ESS focused) |

### Security Enhancements
- **user_identity Validation**: 9-digit format validation
- **Module Access Control**: Granular read/write/delete permissions
- **Comprehensive Audit**: All activities tracked with user identity
- **Rate Limiting**: 100 requests per minute per IP
- **Input Validation**: All inputs validated and sanitized

### Performance Optimizations
- **Hierarchical Queries**: Optimized tree traversal for modules
- **Path-based Indexing**: Fast module lookups using path field
- **Cached Permissions**: Redis-based permission caching
- **Efficient Audit**: Optimized audit log storage and retrieval

## [1.0.0] - 2025-12-22

### Added
- **RBAC System Implementation**
  - Multi-company and multi-branch support
  - Role-based access control with flexible scoping
  - User role assignments with company/branch level permissions

- **Token-Based Authentication**
  - Access token (15 minutes expiration)
  - Refresh token (7 days expiration) 
  - Token binding to device, IP, and User-Agent
  - Secure token storage in Redis with SHA-256 hashing
  - Token rotation and revocation

- **Database Schema**
  - Users table with active status
  - Companies and branches tables
  - Roles and permissions tables
  - Role-permission mapping
  - User-role assignments with scope

- **API Endpoints**
  - Authentication endpoints (login, refresh, logout)
  - User management CRUD operations
  - Health check endpoint
  - Protected routes with Bearer token authentication

### Technical Details
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis for token storage
- **Authentication**: Custom token-based system (Sanctum-like)
- **Architecture**: Clean architecture with domain separation

## Migration Guide

### From v1.0.0 to v2.0.0 (Module-Based System)

**⚠️ BREAKING CHANGES - Complete System Redesign**

#### Database Migration
1. **Backup existing data** - Create full database backup
2. **Run new migrations** - Execute migrations 007-010
3. **Update user data** - Add user_identity to existing users
4. **Verify module data** - Ensure all 89+ modules are loaded

#### Authentication Changes
```bash
# Old (v1.0.0)
curl -X POST /api/v1/auth/login \
  -d '{"email":"admin@system.com","password":"password123"}'

# New (v2.0.0)
curl -X POST /api/v1/auth/login \
  -d '{"user_identity":"100000003","password":"password123"}'
```

#### API Endpoint Changes
| Old Endpoint | New Endpoint | Notes |
|--------------|--------------|-------|
| `/api/v1/permissions` | `/api/v1/modules` | Permission → Module |
| `/api/v1/role-permissions` | `/api/v1/role-management` | New role management system |
| N/A | `/api/v1/users/identity/:id/modules` | New user access endpoint |
| N/A | `/api/v1/audit/logs` | New audit logging |

#### Client Application Updates
1. **Update Authentication**: Change from email to user_identity
2. **Update API Calls**: Use new module-based endpoints
3. **Update Test Data**: Use new user_identity test users
4. **Update Validation**: Handle new response formats

#### Documentation Updates
1. **Import New Postman Collection**: `ERP_RBAC_API_MODULE_BASED.postman_collection.json`
2. **Update API Documentation**: Use `API_DOCUMENTATION_MODULE_BASED.md`
3. **Follow New Quick Start**: Use `QUICK_START_MODULE_BASED.md`

### Production Deployment Checklist

#### Pre-Deployment
- [ ] **Database Backup**: Full backup of existing data
- [ ] **Migration Testing**: Test migrations on staging environment
- [ ] **API Testing**: Verify all endpoints with new Postman collection
- [ ] **User Data Migration**: Ensure user_identity is populated
- [ ] **Client App Updates**: Update frontend/mobile apps

#### Deployment
- [ ] **Run Migrations**: Execute migrations 007-010 in order
- [ ] **Verify Data**: Check module and user data integrity
- [ ] **Test Authentication**: Verify user_identity login works
- [ ] **Test Module Access**: Verify user permissions work correctly
- [ ] **Monitor Logs**: Check audit logs are being created

#### Post-Deployment
- [ ] **Performance Monitoring**: Monitor API response times
- [ ] **Error Monitoring**: Check for authentication/authorization errors
- [ ] **User Training**: Update user documentation and training
- [ ] **Client Rollout**: Deploy updated client applications

## Contributors

- **v2.0.0**: Complete module-based RBAC system redesign
- **v1.0.0**: Initial implementation and RBAC system design

## Support

- **Documentation**: Check `docs/` directory for complete guides
- **API Testing**: Use module-based Postman collection
- **Migration Help**: Refer to migration guide above
- **Issues**: Create GitHub issue for support