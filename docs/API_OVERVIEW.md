# ERP RBAC API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require JWT authentication:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

## Quick Start
1. Import Postman collection: `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json`
2. Set environment: `docs/ERP_RBAC_Environment_Module_Based.postman_environment.json`
3. Login to get access token
4. Use protected endpoints

## API Endpoints

### 🔐 Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login dengan user_identity |
| POST | `/auth/login-email` | Login dengan email |
| POST | `/auth/refresh` | Refresh JWT token |
| POST | `/auth/logout` | Logout user |

### 👥 User Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | Get users list |
| GET | `/users/{id}` | Get user by ID |
| POST | `/users` | Create new user |
| PUT | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |
| GET | `/users/{id}/modules` | Get user modules |
| POST | `/users/check-access` | Check user access |
| PUT | `/users/{id}/password` | Change user password |

### 🏢 Company Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/companies` | Get companies list |
| GET | `/companies/{id}` | Get company by ID |
| POST | `/companies` | Create new company |
| PUT | `/companies/{id}` | Update company |
| DELETE | `/companies/{id}` | Delete company |

### 🏪 Branch Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/branches` | Get branches list |
| GET | `/branches/{id}` | Get branch by ID |
| POST | `/branches` | Create new branch |
| PUT | `/branches/{id}` | Update branch |
| DELETE | `/branches/{id}` | Delete branch |
| GET | `/branches/company/{companyId}` | Get company branches |

### 🎭 Role Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/roles` | Get roles list |
| GET | `/roles/{id}` | Get role by ID |
| POST | `/roles` | Create new role |
| PUT | `/roles/{id}` | Update role |
| DELETE | `/roles/{id}` | Delete role |
| POST | `/role-management/assign-user-role` | Assign role to user |
| PUT | `/role-management/role/{roleId}/modules` | Update role permissions |
| GET | `/role-management/user/{userId}/roles` | Get user roles |

### 📦 Module Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/modules` | Get modules list |
| GET | `/modules/{id}` | Get module by ID |
| POST | `/modules` | Create new module |
| PUT | `/modules/{id}` | Update module |
| DELETE | `/modules/{id}` | Delete module |
| GET | `/modules/tree` | Get module tree |

### 💳 Subscription Management

#### Public Endpoints (No Auth)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/plans` | Get all subscription plans |
| GET | `/plans/{id}` | Get subscription plan by ID |

#### Protected Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/subscription/subscriptions` | Get subscriptions list |
| GET | `/subscription/subscriptions/{id}` | Get subscription by ID |
| POST | `/subscription/subscriptions` | Create new subscription |
| PUT | `/subscription/subscriptions/{id}` | Update subscription |
| **POST** | **`/subscription/subscriptions/{id}/renew`** | **🆕 Renew subscription** |
| POST | `/subscription/subscriptions/{id}/cancel` | Cancel subscription |
| GET | `/subscription/companies/{id}/subscription` | Get company subscription |

#### Admin Endpoints (Admin Auth Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| **POST** | **`/admin/subscription-plans`** | **🆕 Create subscription plan** |
| **PUT** | **`/admin/subscription-plans/{id}`** | **🆕 Update subscription plan** |
| **DELETE** | **`/admin/subscription-plans/{id}`** | **🆕 Delete subscription plan** |

### 📊 Audit Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/audit/logs` | Get audit logs |
| POST | `/audit/logs` | Create audit log |
| GET | `/audit/users/{userId}/logs` | Get user audit logs |

## New Features

### Subscription Renewal API
**Endpoint**: `POST /subscription/subscriptions/{id}/renew`

**Request Body**:
```json
{
    "billing_cycle": "yearly",  // "monthly" or "yearly"
    "plan_id": 2               // optional: untuk upgrade plan
}
```

**Use Cases**:
- Perpanjang langganan dari monthly ke yearly
- Upgrade plan saat renewal
- Self-service renewal

### Admin Subscription Plan Management
**Authentication**: Requires admin JWT token

**Create Plan**:
```json
{
    "name": "premium",
    "display_name": "Premium Plan",
    "description": "Premium features",
    "price_monthly": 99000,
    "price_yearly": 990000,
    "max_users": 50,
    "max_branches": 10,
    "features": {
        "advanced_reporting": true,
        "api_access": true,
        "priority_support": true
    }
}
```

## Response Format

### Success Response
```json
{
    "success": true,
    "message": "Operation successful",
    "data": { ... }
}
```

### Error Response
```json
{
    "success": false,
    "message": "Error message",
    "error": "Detailed error description"
}
```

## HTTP Status Codes
- **200 OK**: Request successful
- **201 Created**: Resource created
- **400 Bad Request**: Invalid request
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Access denied
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

## Testing with Postman
1. Import collection: `docs/ERP_RBAC_API_MODULE_BASED.postman_collection.json`
2. Import environment: `docs/ERP_RBAC_Environment_Module_Based.postman_environment.json`
3. Login to get access token (auto-saved to environment)
4. Test endpoints using the organized folders

## Additional Documentation
- [Backend Engineer SOP](./BACKEND_ENGINEER_SOP.md)
- [Clean Architecture Guide](./CLEAN_ARCHITECTURE.md)
- [Project Structure](./PROJECT_STRUCTURE.md)