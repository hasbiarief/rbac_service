# Authentication API Reference

## Overview

Huminor RBAC Authentication API menyediakan endpoint untuk authentication, authorization, dan session management. API ini menggunakan JWT tokens untuk authentication dan mendukung role-based access control (RBAC).

**Base URL**: `http://localhost:8082/api/v1`

## Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant RBAC
    participant Redis
    participant DB

    Client->>API: POST /auth/login
    API->>DB: Validate credentials
    API->>DB: Get user roles & permissions
    API->>Redis: Store access token
    API->>Redis: Store refresh token
    API->>Client: Return tokens + user data

    Client->>API: GET /protected-resource
    Note over Client,API: Authorization: Bearer <access_token>
    API->>Redis: Validate access token
    API->>RBAC: Check permissions
    API->>Client: Return protected data

    Client->>API: POST /auth/refresh
    API->>Redis: Validate refresh token
    API->>Redis: Store new access token
    API->>Client: Return new tokens
```

## Endpoints

### 1. Login with User Identity

**Endpoint**: `POST /auth/login`

**Description**: Authenticate user menggunakan user_identity dan password.

**Request Headers**:
```http
Content-Type: application/json
User-Agent: <client-user-agent>
```

**Request Body**:
```json
{
  "user_identity": "800000001",
  "password": "password123",
  "user_agent": "Mozilla/5.0...",  // Optional
  "ip": "192.168.1.1"              // Optional
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `user_identity` | string | Yes | User identity number |
| `password` | string | Yes | User password (min 6 characters) |
| `user_agent` | string | No | Client user agent |
| `ip` | string | No | Client IP address |

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a3f485997fd448775128c5b9f5011ee3...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": 16,
      "name": "Hasbi Due",
      "email": "hasbi@company.com",
      "user_identity": "800000001",
      "is_active": true,
      "created_at": "2026-01-17T09:40:46Z",
      "updated_at": "2026-01-17T09:40:46Z",
      "modules": {
        "Module Management": [
          ["Module Management", "/module", "Test", "Manage modules on huminor system"],
          ["Create Modules", "/module/create-module", "TestChild", "Create new module"]
        ]
      },
      "role_assignments": [
        {
          "assignment_id": 34,
          "assignment_level": "unit",
          "branch_id": 1,
          "branch_name": "Pusat",
          "company_id": 1,
          "company_name": "PT. Huminor Tech",
          "role_description": "Console admin role for manage system",
          "role_id": 13,
          "role_name": "CONSOLE ADMIN",
          "unit_id": 16,
          "unit_name": "Dewa"
        }
      ],
      "total_roles": 1
    }
  }
}
```

**Error Responses**:

*401 Unauthorized*:
```json
{
  "success": false,
  "message": "Login failed",
  "error": "kredensial tidak valid"
}
```

*400 Bad Request*:
```json
{
  "success": false,
  "message": "Invalid request",
  "error": "user_identity must be provided"
}
```

### 2. Login with Email

**Endpoint**: `POST /auth/login-email`

**Description**: Authenticate user menggunakan email dan password.

**Request Body**:
```json
{
  "email": "hasbi@company.com",
  "password": "password123",
  "user_agent": "Mozilla/5.0...",  // Optional
  "ip": "192.168.1.1"              // Optional
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | Valid email address |
| `password` | string | Yes | User password (min 6 characters) |
| `user_agent` | string | No | Client user agent |
| `ip` | string | No | Client IP address |

**Response**: Same as login with user_identity

### 3. Refresh Token

**Endpoint**: `POST /auth/refresh`

**Description**: Generate new access token menggunakan refresh token.

**Request Body**:
```json
{
  "refresh_token": "a3f485997fd448775128c5b9f5011ee3..."
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Token successfully refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "b4g596008ge559886239c6a0f5022ff4...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Error Response** (401 Unauthorized):
```json
{
  "success": false,
  "message": "Token refresh failed",
  "error": "refresh token tidak valid atau kedaluwarsa"
}
```

### 4. Logout

**Endpoint**: `POST /auth/logout`

**Description**: Logout user dan revoke semua tokens.

**Request Body** (Option 1 - by token):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Request Body** (Option 2 - by user_id):
```json
{
  "user_id": 16
}
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

### 5. Check User Tokens

**Endpoint**: `GET /auth/check-tokens`

**Description**: Validate dan check user tokens.

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `user_id` | integer | Yes | User ID to check tokens |

**Example Request**:
```http
GET /auth/check-tokens?user_id=16
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Token check completed",
  "data": {
    "access_tokens": [
      {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": "2026-01-25T10:30:00Z",
        "user_agent": "Mozilla/5.0...",
        "ip": "192.168.1.1"
      }
    ],
    "refresh_tokens": [
      {
        "token": "a3f485997fd448775128c5b9f5011ee3...",
        "family_id": "abc123def456",
        "created_at": "2026-01-25T09:30:00Z"
      }
    ]
  }
}
```

### 6. Get User Session Count

**Endpoint**: `GET /auth/session-count`

**Description**: Get jumlah active sessions untuk user.

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `user_id` | integer | Yes | User ID to get session count |

**Example Request**:
```http
GET /auth/session-count?user_id=16
```

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Session count retrieved",
  "data": {
    "user_id": 16,
    "access_token_count": 2,
    "refresh_token_count": 1,
    "total_sessions": 2
  }
}
```

### 7. Cleanup Expired Tokens

**Endpoint**: `POST /auth/cleanup-expired`

**Description**: Cleanup expired tokens dari Redis (admin only).

**Success Response** (200 OK):
```json
{
  "success": true,
  "message": "Expired tokens cleaned up successfully",
  "data": null
}
```

## Authentication Headers

Untuk protected endpoints, gunakan Bearer token di header:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

## Token Information

### Access Token
- **Type**: JWT (JSON Web Token)
- **Expiration**: 15 minutes (900 seconds)
- **Storage**: Redis
- **Usage**: Include in Authorization header untuk API calls

### Refresh Token
- **Type**: Random string
- **Expiration**: 7 days
- **Storage**: Redis
- **Usage**: Generate new access token when expired

## Error Codes

| HTTP Status | Error Code | Description |
|-------------|------------|-------------|
| 400 | `BAD_REQUEST` | Invalid request format atau missing required fields |
| 401 | `UNAUTHORIZED` | Invalid credentials atau expired token |
| 403 | `FORBIDDEN` | Insufficient permissions |
| 429 | `TOO_MANY_REQUESTS` | Rate limit exceeded |
| 500 | `INTERNAL_SERVER_ERROR` | Server error |

## Rate Limiting

- **Authentication endpoints**: 10 requests per minute per IP
- **Token refresh**: 5 requests per minute per user
- **Session check**: 20 requests per minute per user

## Security Considerations

### Token Security
- Access tokens expire dalam 15 menit
- Refresh tokens expire dalam 7 hari
- Tokens disimpan di Redis dengan TTL
- Logout akan revoke semua user tokens

### Password Security
- Minimum 6 characters
- Hashed menggunakan bcrypt
- No password in response

### Session Management
- Multiple sessions per user supported
- Family-based refresh token rotation
- Automatic cleanup expired tokens

## Integration Examples

### JavaScript/Fetch
```javascript
// Login
const loginResponse = await fetch('http://localhost:8082/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    user_identity: '800000001',
    password: 'password123'
  })
});

const loginData = await loginResponse.json();
const accessToken = loginData.data.access_token;

// Use token for protected requests
const protectedResponse = await fetch('http://localhost:8082/api/v1/protected-endpoint', {
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json'
  }
});
```

### cURL Examples

**Login**:
```bash
curl -X POST http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_identity": "800000001",
    "password": "password123"
  }'
```

**Refresh Token**:
```bash
curl -X POST http://localhost:8082/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
  }'
```

**Protected Request**:
```bash
curl -X GET http://localhost:8082/api/v1/protected-endpoint \
  -H "Authorization: Bearer your-access-token-here" \
  -H "Content-Type: application/json"
```

## Postman Collection

Import collection untuk testing: [Download Postman Collection](../postman/auth-endpoints.json)

Environment variables yang diperlukan:
- `base_url`: `http://localhost:8082`
- `access_token`: (akan di-set otomatis setelah login)
- `refresh_token`: (akan di-set otomatis setelah login)

## Testing Credentials

Untuk testing, gunakan credentials berikut:

| User Identity | Email | Password | Role |
|---------------|-------|----------|------|
| `800000001` | `hasbi@company.com` | `password123` | CONSOLE ADMIN |
| `100000001` | `naruto@company.com` | `password123` | User |
| `100000002` | `sakura@company.com` | `password123` | User |

## Troubleshooting

### Common Issues

**1. CORS Error**
```
Access to fetch at 'http://localhost:8082' from origin 'http://localhost:3000' has been blocked by CORS policy
```
**Solution**: Pastikan CORS_ORIGINS di .env includes your frontend URL

**2. Token Expired**
```json
{
  "success": false,
  "message": "Token validation failed",
  "error": "token expired"
}
```
**Solution**: Use refresh token untuk generate new access token

**3. Invalid Credentials**
```json
{
  "success": false,
  "message": "Login failed",
  "error": "kredensial tidak valid"
}
```
**Solution**: Check user_identity/email dan password

**4. User Inactive**
```json
{
  "success": false,
  "message": "Login failed",
  "error": "akun pengguna tidak aktif"
}
```
**Solution**: Contact admin untuk activate user account

### Debug Tips

1. **Check server logs** untuk detailed error messages
2. **Validate request format** sesuai dengan schema
3. **Check Redis connection** untuk token storage issues
4. **Verify database connection** untuk user validation
5. **Test dengan Postman** untuk isolate client-side issues

---

**Next Steps**: 
- [Client Integration Guide](./CLIENT_INTEGRATION_GUIDE.md) - Frontend integration
- [Backend Integration Guide](./BACKEND_INTEGRATION_GUIDE.md) - Server-side integration
- [RBAC Integration Guide](./RBAC_INTEGRATION_GUIDE.md) - Permission system