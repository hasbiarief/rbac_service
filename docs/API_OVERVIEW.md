# API Overview - RBAC Service

## 📍 Base Information

- **Base URL**: `http://localhost:8081/api/v1`
- **Authentication**: JWT Bearer Token
- **Content-Type**: `application/json`
- **Postman Collection**: `docs/HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json`

## 🔐 Authentication

### Login Endpoints

**Login with user_identity:**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "user_identity": "100000001",
  "password": "password123"
}
```

**Login with email:**
```http
POST /api/v1/auth/login-email
Content-Type: application/json

{
  "email": "admin@system.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "name": "Admin User",
      "email": "admin@system.com",
      "user_identity": "100000001"
    }
  }
}
```

### Token Management

**Refresh Token:**
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

**Logout:**
```http
POST /api/v1/auth/logout
Content-Type: application/json

{
  "user_id": 1
}
```

**Check User Tokens:**
```http
GET /api/v1/auth/check-tokens?user_id=1
```

**Get Session Count:**
```http
GET /api/v1/auth/session-count?user_id=1
```

## 👥 User Management

### User CRUD

**Get All Users:**
```http
GET /api/v1/users?limit=10&offset=0&search=admin
Authorization: Bearer {token}
```

**Get User by ID:**
```http
GET /api/v1/users/1
Authorization: Bearer {token}
```

**Create User:**
```http
POST /api/v1/users
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "user_identity": "100000010",
  "password": "password123"
}
```

**Update User:**
```http
PUT /api/v1/users/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Updated",
  "email": "john.updated@example.com",
  "is_active": true
}
```

**Delete User:**
```http
DELETE /api/v1/users/1
Authorization: Bearer {token}
```

### User Module Access

**Get User Modules by ID:**
```http
GET /api/v1/users/3/modules?limit=20&category=Core%20HR
Authorization: Bearer {token}
```

**Get User Modules by Identity:**
```http
GET /api/v1/users/identity/100000001/modules?limit=20
Authorization: Bearer {token}
```

**Check User Access:**
```http
POST /api/v1/users/check-access
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_identity": "100000001",
  "module_url": "/core-hr/employees"
}
```

### Password Management

**Change Password:**
```http
PUT /api/v1/users/1/password
Authorization: Bearer {token}
Content-Type: application/json

{
  "current_password": "password123",
  "new_password": "newpassword456",
  "confirm_password": "newpassword456"
}
```

## 🎭 Role Management

### Basic Role Operations

**Get All Roles:**
```http
GET /api/v1/roles?limit=10&offset=0&search=admin
Authorization: Bearer {token}
```

**Get Role by ID:**
```http
GET /api/v1/roles/1
Authorization: Bearer {token}
```

**Create Role:**
```http
POST /api/v1/roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "CUSTOM_ROLE",
  "description": "Custom role for testing"
}
```

### Role Assignment System

**Assign User Role (Company/Branch Level):**
```http
POST /api/v1/role-management/assign-user-role
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": 2,
  "role_id": 4,
  "company_id": 1,
  "branch_id": 1
}
```

**Assign User Role (Unit Level):**
```http
POST /api/v1/role-management/assign-user-role
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": 2,
  "role_id": 4,
  "company_id": 1,
  "branch_id": 1,
  "unit_id": 1
}
```

**Remove User Role:**
```http
DELETE /api/v1/role-management/user/10/role/4
Authorization: Bearer {token}
```

**Bulk Assign Roles:**
```http
POST /api/v1/role-management/bulk-assign-roles
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_ids": [2, 3],
  "role_id": 4,
  "company_id": 1,
  "branch_id": 1,
  "unit_id": 1
}
```

**Update Role Modules:**
```http
PUT /api/v1/role-management/role/4/modules
Authorization: Bearer {token}
Content-Type: application/json

{
  "modules": [
    {
      "module_id": 1,
      "can_read": true,
      "can_write": true,
      "can_delete": false
    }
  ]
}
```

**Get Users by Role:**
```http
GET /api/v1/role-management/role/3/users?limit=10
Authorization: Bearer {token}
```

**Get User Roles:**
```http
GET /api/v1/role-management/user/3/roles
Authorization: Bearer {token}
```

**Get User Access Summary:**
```http
GET /api/v1/role-management/user/3/access-summary
Authorization: Bearer {token}
```

### Debug Endpoints

**Get All User Role Assignments:**
```http
GET /api/v1/role-management/debug/all-assignments
Authorization: Bearer {token}
```

**Get User Role Assignments:**
```http
GET /api/v1/role-management/debug/user/3/roles
Authorization: Bearer {token}
```

**Get Role-Users Mapping:**
```http
GET /api/v1/role-management/debug/role-users-mapping
Authorization: Bearer {token}
```

## 📦 Module System

### Module CRUD

**Get All Modules:**
```http
GET /api/v1/modules?limit=20&offset=0&category=Core%20HR&search=employee
Authorization: Bearer {token}
```

**Get Module by ID:**
```http
GET /api/v1/modules/1
Authorization: Bearer {token}
```

**Create Module:**
```http
POST /api/v1/modules
Authorization: Bearer {token}
Content-Type: application/json

{
  "category": "Test Category",
  "name": "Test Module",
  "url": "/test/module",
  "icon": "Test",
  "description": "Test module",
  "subscription_tier": "basic",
  "is_active": true
}
```

**Update Module:**
```http
PUT /api/v1/modules/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Updated Module",
  "description": "Updated description",
  "is_active": true
}
```

**Delete Module:**
```http
DELETE /api/v1/modules/1
Authorization: Bearer {token}
```

### Module Hierarchy

**Get Module Tree by Parent:**
```http
GET /api/v1/modules/tree?parent=Core%20HR%20Management
Authorization: Bearer {token}
```

**Get Module Children:**
```http
GET /api/v1/modules/1/children
Authorization: Bearer {token}
```

**Get Module Ancestors:**
```http
GET /api/v1/modules/5/ancestors
Authorization: Bearer {token}
```

## 🏢 Company Management

**Get All Companies:**
```http
GET /api/v1/companies?limit=10&offset=0&search=PT
Authorization: Bearer {token}
```

**Get Company by ID:**
```http
GET /api/v1/companies/1
Authorization: Bearer {token}
```

**Create Company:**
```http
POST /api/v1/companies
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "PT. Test Company",
  "code": "TEST"
}
```

**Update Company:**
```http
PUT /api/v1/companies/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "PT. Updated Company",
  "code": "UPDATED",
  "is_active": true
}
```

**Delete Company:**
```http
DELETE /api/v1/companies/1
Authorization: Bearer {token}
```

## 🌳 Branch Management (Hierarchical)

**Get All Branches:**
```http
GET /api/v1/branches?limit=10&offset=0&company_id=1
Authorization: Bearer {token}
```

**Get Branch Hierarchy (Nested):**
```http
GET /api/v1/branches?nested=true&company_id=1
Authorization: Bearer {token}
```

**Get Company Branches:**
```http
GET /api/v1/branches/company/1?nested=true
Authorization: Bearer {token}
```

**Get Branch Hierarchy by ID:**
```http
GET /api/v1/branches/3/hierarchy?nested=true
Authorization: Bearer {token}
```

**Create Branch:**
```http
POST /api/v1/branches
Authorization: Bearer {token}
Content-Type: application/json

{
  "company_id": 1,
  "name": "Jakarta Branch",
  "code": "JKT",
  "parent_id": null
}
```

**Create Sub-Branch:**
```http
POST /api/v1/branches
Authorization: Bearer {token}
Content-Type: application/json

{
  "company_id": 1,
  "name": "Jakarta Pusat",
  "code": "JKT_PUSAT",
  "parent_id": 1
}
```

## 🏭 Unit Management (Unit-Based RBAC)

**Get All Units:**
```http
GET /api/v1/units?limit=20&offset=0&branch_id=1
Authorization: Bearer {token}
```

**Get Unit by ID:**
```http
GET /api/v1/units/1
Authorization: Bearer {token}
```

**Create Unit:**
```http
POST /api/v1/units
Authorization: Bearer {token}
Content-Type: application/json

{
  "branch_id": 1,
  "name": "HR Department",
  "code": "HR",
  "description": "Human Resources",
  "parent_id": null
}
```

**Update Unit:**
```http
PUT /api/v1/units/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Updated HR Department",
  "description": "Updated description",
  "is_active": true
}
```

**Delete Unit:**
```http
DELETE /api/v1/units/1
Authorization: Bearer {token}
```

**Get Unit with Statistics:**
```http
GET /api/v1/units/1/stats
Authorization: Bearer {token}
```

**Get Unit Hierarchy:**
```http
GET /api/v1/branches/1/units/hierarchy
Authorization: Bearer {token}
```

### Unit Role Management

**Get Unit Roles:**
```http
GET /api/v1/units/1/roles
Authorization: Bearer {token}
```

**Assign Role to Unit:**
```http
POST /api/v1/units/1/roles/3
Authorization: Bearer {token}
```

**Remove Role from Unit:**
```http
DELETE /api/v1/units/1/roles/3
Authorization: Bearer {token}
```

**Get Unit Permissions:**
```http
GET /api/v1/units/1/roles/3/permissions
Authorization: Bearer {token}
```

**Get User Effective Permissions:**
```http
GET /api/v1/users/1/effective-permissions
Authorization: Bearer {token}
```

## 💳 Subscription Management

### Subscription Plans (Public)

**Get All Plans:**
```http
GET /api/v1/plans
```

**Get Plan by ID:**
```http
GET /api/v1/plans/1
```

### Subscription Operations

**Get All Subscriptions:**
```http
GET /api/v1/subscription/subscriptions?limit=10&offset=0
Authorization: Bearer {token}
```

**Get Subscription by ID:**
```http
GET /api/v1/subscription/subscriptions/1
Authorization: Bearer {token}
```

**Create Subscription:**
```http
POST /api/v1/subscription/subscriptions
Authorization: Bearer {token}
Content-Type: application/json

{
  "company_id": 1,
  "plan_id": 2,
  "billing_cycle": "yearly",
  "auto_renew": true
}
```

**Update Subscription:**
```http
PUT /api/v1/subscription/subscriptions/1
Authorization: Bearer {token}
Content-Type: application/json

{
  "plan_id": 3,
  "auto_renew": false
}
```

**Renew Subscription:**
```http
POST /api/v1/subscription/subscriptions/1/renew
Authorization: Bearer {token}
Content-Type: application/json

{
  "billing_cycle": "yearly",
  "plan_id": 3
}
```

**Cancel Subscription:**
```http
POST /api/v1/subscription/subscriptions/1/cancel
Authorization: Bearer {token}
Content-Type: application/json

{
  "reason": "Switching to different solution",
  "cancel_immediately": false
}
```

**Mark Payment as Paid:**
```http
POST /api/v1/subscription/subscriptions/1/mark-paid
Authorization: Bearer {token}
```

**Get Company Subscription:**
```http
GET /api/v1/subscription/companies/1/subscription
Authorization: Bearer {token}
```

**Get Company Subscription Status:**
```http
GET /api/v1/subscription/companies/1/status
Authorization: Bearer {token}
```

**Check Module Access:**
```http
GET /api/v1/subscription/module-access/1/1
Authorization: Bearer {token}
```

**Get Subscription Statistics:**
```http
GET /api/v1/subscription/stats
Authorization: Bearer {token}
```

**Get Expiring Subscriptions:**
```http
GET /api/v1/subscription/expiring?days=30
Authorization: Bearer {token}
```

## 📊 Audit Logging

**Get Audit Logs:**
```http
GET /api/v1/audit/logs?limit=20&offset=0&user_id=3&action=login
Authorization: Bearer {token}
```

**Get User Audit Logs by ID:**
```http
GET /api/v1/audit/users/3/logs?limit=10
Authorization: Bearer {token}
```

**Get User Audit Logs by Identity:**
```http
GET /api/v1/audit/users/identity/100000001/logs?limit=10
Authorization: Bearer {token}
```

**Get Audit Statistics:**
```http
GET /api/v1/audit/stats
Authorization: Bearer {token}
```

**Create Audit Log (Manual):**
```http
POST /api/v1/audit/logs
Authorization: Bearer {token}
Content-Type: application/json

{
  "user_id": 1,
  "user_identity": "100000001",
  "action": "manual_test",
  "resource": "test_resource",
  "resource_id": "123",
  "method": "POST",
  "url": "/api/v1/test",
  "status": "success",
  "status_code": 200,
  "message": "Manual audit log"
}
```

## 📋 Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": 1,
    "name": "Example"
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Operation failed",
  "error": "Detailed error message"
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved",
  "data": {
    "data": [...],
    "total": 100,
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}
```

## 🔒 Middleware

### Authentication
Semua endpoint (kecuali `/auth/login`, `/auth/login-email`, `/auth/refresh`, `/plans`) memerlukan JWT token:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### CORS
CORS enabled untuk:
- `http://localhost:3000`
- `http://localhost:3001`
- `http://127.0.0.1:3000`
- `http://127.0.0.1:3001`

### Rate Limiting
- Default endpoints: 10 req/sec, burst 50
- Check-access endpoint: 30 req/sec, burst 100
- Rate limit per IP address
- Returns 429 status when exceeded

## 📚 Testing

Gunakan Postman collection untuk testing lengkap:
- `docs/HUMINOR_RBAC_API_MODULE_BASED.postman_collection.json`
- `docs/HUMINOR_RBAC_Environment_Module_Based.postman_environment.json`

### Test Users (password: `password123`)
- `admin@system.com` - System Admin
- `hr@company.com` - HR Manager
- `superadmin@company.com` - Super Admin

## 🔗 Related Documentation

- [Backend Engineer Rules](ENGINEER_RULES.md) - Development guide
- [Project Structure](PROJECT_STRUCTURE.md) - Architecture overview
- [README](../README.md) - Project overview
