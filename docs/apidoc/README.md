# API Documentation Module

Dokumentasi lengkap untuk API Documentation System dan semua API endpoints dalam Huminor RBAC Service.

## 📋 Daftar Dokumentasi

### 🏗️ Core API Documentation
- **[API Overview](API_OVERVIEW.md)** - Dokumentasi lengkap semua API endpoints dengan request/response examples
- **[API Documentation System](API_DOCUMENTATION_SYSTEM.md)** - Sistem dokumentasi API terintegrasi
- **[API Workflow Guide](API_WORKFLOW_GUIDE.md)** - Complete setup process dari user creation sampai module visibility

### 🔧 API Documentation Features
- **[API Doc Authentication](API_DOC_AUTHENTICATION.md)** - Authentication untuk API Documentation
- **[API Doc Auto Discovery](API_DOC_AUTO_DISCOVERY.md)** - Auto discovery sistem untuk API endpoints
- **[API Doc Export Formats](API_DOC_EXPORT_FORMATS.md)** - Format export (Postman, OpenAPI, Insomnia, Swagger, Apidog)
- **[API Doc Migration Plan](API_DOC_MIGRATION_PLAN.md)** - Migration plan untuk API Documentation System

## 🚀 Quick Navigation

### Untuk Developer
1. **Mulai dengan** → [API Overview](API_OVERVIEW.md)
2. **Pahami sistem** → [API Documentation System](API_DOCUMENTATION_SYSTEM.md)
3. **Setup authentication** → [API Doc Authentication](API_DOC_AUTHENTICATION.md)
4. **Complete workflow** → [API Workflow Guide](API_WORKFLOW_GUIDE.md)
5. **Export collections** → [API Doc Export Formats](API_DOC_EXPORT_FORMATS.md)

### Untuk DevOps/Admin
1. **Migration setup** → [API Doc Migration Plan](API_DOC_MIGRATION_PLAN.md)
2. **Auto discovery** → [API Doc Auto Discovery](API_DOC_AUTO_DISCOVERY.md)
3. **Authentication config** → [API Doc Authentication](API_DOC_AUTHENTICATION.md)

## 📡 API Endpoints Overview

### Authentication Module
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh
- `POST /api/v1/auth/logout` - User logout

### API Documentation System
- `GET /api/v1/api-docs/collections` - Get all collections
- `POST /api/v1/api-docs/collections` - Create collection
- `GET /api/v1/api-docs/collections/:id/export/:format` - Export collection
- `GET /api/v1/api-docs/endpoints` - Get all endpoints
- `POST /api/v1/api-docs/endpoints` - Create endpoint

### RBAC Module
- `GET /api/v1/users` - Get users
- `GET /api/v1/roles` - Get roles
- `POST /api/v1/rbac/check-permission` - Check user permission

### Management Modules
- `GET /api/v1/companies` - Company management
- `GET /api/v1/branches` - Branch management
- `GET /api/v1/units` - Unit management
- `GET /api/v1/modules` - Module management

## 🔗 Related Documentation

- **[Integration Guide](../integration/)** - Authentication integration untuk external services
- **[Project Structure](../PROJECT_STRUCTURE.md)** - Arsitektur project
- **[Engineer Rules](../ENGINEER_RULES.md)** - Development guidelines
- **[Quick Start](../QUICK_START.md)** - Setup dan development guide

## 📦 Postman Collections

Import collections untuk testing:
- `docs/HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json`
- `docs/HUMINOR_RBAC_Environment_Module_Based.postman_environment.json`

## 🛠️ Tools & Testing

### API Testing
1. Import Postman collection
2. Setup environment variables
3. Test authentication endpoints
4. Test API Documentation System endpoints
5. Test RBAC permissions

### Export Formats
- **Postman Collection** - Native Postman format
- **OpenAPI 3.0** - Industry standard API specification
- **Insomnia** - REST client format
- **Swagger** - API documentation format
- **Apidog** - API development platform format

---

**Navigation**: [← Back to Docs](../) | [API Overview →](API_OVERVIEW.md)