# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2025-12-28

### Added
- **Modular Validation System** - Separated validation rules into dedicated files
- **Centralized Response System** - Consistent API response format across all endpoints
- **Rate Limiting System** - Global rate limiting (100 requests/minute per IP)
- **Dual Authentication Methods** - Support for both user_identity and email login
- **Validation Middleware Integration** - All handlers now use validated body from middleware

### Changed
- **Routes Architecture** - Moved from inline validation (1400+ lines) to modular validation (~200 lines)
- **Handler Pattern** - All handlers now use validated body from context instead of direct binding
- **Validation Middleware** - Updated to use `c.ShouldBind()` instead of `c.ShouldBindJSON()`
- **Authentication Flow** - Refresh token endpoint no longer requires Bearer token

### Fixed
- **Struct Field Mismatches** - Fixed BulkAssignRolesRequest and UpdateModuleRequest field mappings
- **Validation Consistency** - All endpoints now use centralized validation system
- **Response Format** - Consistent error and success response format

### Technical Improvements
- **Code Organization** - Validation rules separated into 8 modular files
- **Maintainability** - Routes file reduced from 1400+ lines to ~200 lines
- **Reusability** - Validation rules can be reused across endpoints
- **Developer Experience** - Cleaner and more readable route definitions

### File Structure Changes
```
internal/validation/          # NEW - Modular validation rules
├── common_validation.go      # Common rules (ID, pagination, etc.)
├── auth_validation.go        # Authentication validation
├── user_validation.go        # User management validation
├── company_validation.go     # Company validation
├── role_validation.go        # Role & RBAC validation
├── module_validation.go      # Module validation
├── subscription_validation.go # Subscription validation
├── audit_validation.go       # Audit logging validation
└── branch_validation.go      # Branch validation
```

### API Changes
- **Consistent Response Format** - All endpoints now return standardized JSON responses
- **Enhanced Validation** - Comprehensive validation similar to Joi/Zod in Node.js
- **Rate Limiting** - All endpoints protected with rate limiting
- **Dual Login Methods** - `/api/v1/auth/login` (user_identity) and `/api/v1/auth/login-email` (email)

### Documentation Updates
- **API Infrastructure Documentation** - New comprehensive API infrastructure guide
- **Validation System Documentation** - Complete validation system documentation
- **Modular Architecture Documentation** - Updated with latest changes
- **Backend Engineer SOP** - Updated development guidelines

## [2.0.0] - 2025-12-27

### Added
- **Raw SQL Implementation** - Complete replacement of GORM with raw SQL using lib/pq
- **File-based Migration System** - Custom migration system without GORM
- **Base Model System** - Helper methods for raw SQL operations
- **Centralized Response System** - Consistent API response format
- **Rate Limiting** - IP-based rate limiting with Redis storage

### Removed
- **GORM Dependency** - Completely removed GORM and all related dependencies
- **GORM Migrations** - Replaced with file-based SQL migrations
- **ORM Abstractions** - Direct SQL queries for better performance and control

### Changed
- **Database Layer** - All repositories now use raw SQL queries
- **Migration System** - File-based migrations in `migrations/` directory
- **Model Definitions** - Updated models to work with raw SQL scanning
- **Repository Pattern** - Implemented with prepared statements and proper error handling

### Technical Improvements
- **Performance** - Significant performance improvement with raw SQL
- **Transparency** - All SQL queries are visible and optimizable
- **Control** - Full control over database operations and query optimization
- **Debugging** - Easier debugging with visible SQL queries

### Migration Guide
1. **Database Reset** - All existing data needs to be migrated
2. **Model Updates** - Models now use `db` tags for SQL column mapping
3. **Repository Changes** - All repository methods rewritten with raw SQL
4. **Query Optimization** - Queries can now be optimized for specific use cases

## [1.2.0] - 2025-12-26

### Added
- **Modular Architecture** - Restructured monolithic main.go into modular handlers
- **Handler Separation** - 8 dedicated handler files for different domains
- **Server Initialization** - Centralized server setup with dependency injection
- **Route Organization** - Modular route setup functions

### Changed
- **Code Structure** - From monolithic (1400+ lines) to modular architecture
- **Dependency Injection** - Proper DI pattern for handlers, services, and repositories
- **Route Management** - Organized routes by domain/module

### Technical Improvements
- **Maintainability** - Easier to maintain and extend
- **Team Collaboration** - Multiple developers can work on different modules
- **Code Reusability** - Handler and service patterns can be reused
- **Testing** - Easier unit testing per module

### File Structure
```
internal/
├── handlers/           # 8 handler files (~100-200 lines each)
├── routes/            # Centralized route configuration
├── server/            # Server initialization
├── service/           # Business logic layer
├── repository/        # Data access layer
└── models/            # Data models
```

## [1.1.0] - 2025-12-25

### Added
- **Complete RBAC System** - Role-based access control with 59 modules
- **Module Hierarchy** - Hierarchical module system with parent-child relationships
- **Subscription System** - Tiered subscription with module access control
- **Audit Logging** - Comprehensive audit trail for all activities
- **Branch Management** - Hierarchical branch structure for companies

### Features
- **Authentication** - JWT-based authentication with refresh tokens
- **User Management** - Complete user lifecycle management
- **Company Management** - Multi-company support with branch hierarchy
- **Role Management** - Advanced role and permission management
- **Module System** - 59 predefined modules across 8 categories
- **Subscription Tiers** - Basic, Professional, Enterprise, and Custom plans

## [1.0.0] - 2025-12-24

### Added
- **Initial Release** - Basic ERP RBAC API structure
- **Database Schema** - PostgreSQL database with GORM
- **Basic Authentication** - User login and registration
- **Basic RBAC** - Simple role and permission system
- **API Endpoints** - Basic CRUD operations for users and roles

### Technical Stack
- **Go 1.21+** - Programming language
- **Gin Framework** - HTTP web framework
- **GORM** - ORM for database operations (later replaced with raw SQL)
- **PostgreSQL** - Primary database
- **Redis** - Caching and session storage
- **JWT** - Authentication tokens

---

## Migration Notes

### From v2.0.0 to v2.1.0
- **No Breaking Changes** - All existing endpoints remain functional
- **Enhanced Validation** - Better error messages and validation coverage
- **Improved Developer Experience** - Cleaner code structure and better documentation

### From v1.x.x to v2.0.0
- **Breaking Change** - Database schema changes due to GORM removal
- **Migration Required** - Run new migration files to update database
- **Performance Improvement** - Significant performance gains with raw SQL
- **API Compatibility** - All API endpoints remain the same

### From v1.1.0 to v1.2.0
- **No Breaking Changes** - Internal restructuring only
- **Code Organization** - Better organized codebase
- **Development Experience** - Improved development workflow

---

## Upcoming Features

### v2.2.0 (Planned)
- **Advanced Caching** - Redis-based query result caching
- **Bulk Operations** - Batch processing for large datasets
- **Advanced Reporting** - Enhanced analytics and reporting features
- **API Versioning** - Support for multiple API versions

### v3.0.0 (Future)
- **Microservices Architecture** - Split into domain-specific services
- **Event Sourcing** - Event-driven architecture implementation
- **GraphQL Support** - GraphQL API alongside REST
- **Real-time Features** - WebSocket support for real-time updates

---

**Note**: This changelog follows semantic versioning. Major version changes indicate breaking changes, minor versions add new features, and patch versions include bug fixes and improvements.