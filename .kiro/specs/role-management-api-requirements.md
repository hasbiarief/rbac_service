# Role Management API - Requirements Specification

## Overview
This specification defines the requirements for a comprehensive Role-Based Access Control (RBAC) module-based API system that manages users, roles, and permissions within a multi-tenant environment supporting companies, branches, and units.

## User Stories

### US-001: Role Creation and Management
**As a** system administrator  
**I want to** create and manage roles with specific permissions  
**So that** I can control access to different parts of the system

**Acceptance Criteria:**
- ✅ Create new roles with name, description, and active status
- ✅ Update existing roles (name, description, is_active)
- ✅ Delete roles when no longer needed
- ✅ All role operations include proper validation
- ✅ Consistent error handling and response format

### US-002: User Role Assignment
**As a** system administrator  
**I want to** assign roles to users within specific organizational contexts  
**So that** users have appropriate permissions for their responsibilities

**Acceptance Criteria:**
- ✅ Assign single role to user with company, branch, and unit context
- ✅ Bulk assign roles to multiple users simultaneously
- ✅ Remove user roles with proper context validation
- ✅ Support unit_id in role assignments
- ✅ Validate company_id requirements for role operations

### US-003: Role Module Permissions
**As a** system administrator  
**I want to** configure module-level permissions for roles  
**So that** I can define granular access control

**Acceptance Criteria:**
- ✅ Update role modules with read, write, delete permissions
- ✅ Support multiple modules per role
- ✅ Proper field mapping (modules vs permissions)
- ✅ Validate module permissions structure

## Technical Requirements

### API Endpoints
1. **POST** `/api/v1/role-management/roles` - Create role
2. **PUT** `/api/v1/role-management/roles/{id}` - Update role
3. **DELETE** `/api/v1/role-management/roles/{id}` - Delete role
4. **POST** `/api/v1/role-management/assign-user-role` - Assign user role
5. **POST** `/api/v1/role-management/bulk-assign-roles` - Bulk assign roles
6. **DELETE** `/api/v1/role-management/user/{user_id}/role/{role_id}` - Remove user role
7. **PUT** `/api/v1/role-management/role/{id}/modules` - Update role modules

### Data Models

#### Role
```json
{
  "id": "integer",
  "name": "string (required)",
  "description": "string (optional)",
  "is_active": "boolean (default: true)"
}
```

#### User Role Assignment
```json
{
  "user_id": "integer (required)",
  "role_id": "integer (required)",
  "company_id": "integer (required)",
  "branch_id": "integer (optional)",
  "unit_id": "integer (optional)"
}
```

#### Role Module Permissions
```json
{
  "modules": [
    {
      "module_id": "integer (required)",
      "can_read": "boolean (required)",
      "can_write": "boolean (required)",
      "can_delete": "boolean (required)"
    }
  ]
}
```

### Validation Requirements
- All endpoints must include validation middleware
- Required fields must be validated
- Data types must be enforced
- Business logic validation (e.g., company_id requirements)

### Response Format
All APIs must follow consistent response format:
```json
{
  "success": "boolean",
  "message": "string",
  "data": "object (optional)",
  "error": "string (optional)"
}
```

## Implementation Status

### Completed Features ✅
- [x] Role CRUD operations with validation
- [x] User role assignment with unit support
- [x] Bulk role assignment functionality
- [x] Role removal with company validation
- [x] Role module permissions management
- [x] Consistent error handling across all endpoints
- [x] Updated API documentation
- [x] Postman collection updates
- [x] Debug endpoints removed for production readiness

### Architecture Patterns
- Validation middleware integration
- Consistent error response format
- Module-based RBAC architecture
- Multi-tenant support (company/branch/unit)

## Testing Criteria
- All endpoints return proper HTTP status codes
- Validation errors provide clear messages
- Database operations maintain data integrity
- Unit_id properly saved in role assignments
- Company_id validation enforced where required
- Module permissions correctly mapped and stored

## Documentation Requirements
- API documentation updated with all endpoints
- Postman collection includes all test cases
- Response examples for success and error scenarios
- Parameter validation rules documented

## Success Metrics
- All 7 role management endpoints functional
- 100% validation coverage on required fields
- Consistent response format across all APIs
- Zero data integrity issues in role assignments
- Complete API documentation and testing collection