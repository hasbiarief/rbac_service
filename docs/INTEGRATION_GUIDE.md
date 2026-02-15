# Integration Guide - RBAC Authentication

Panduan integrasi RBAC Authentication Service untuk external applications.

## Quick Start

### Base URL
```
Development: http://localhost:8081/api/v1
Production: https://your-domain.com/api/v1
```

### Swagger Documentation
- **Swagger UI**: `http://localhost:8081/swagger/index.html`
- **OpenAPI JSON**: `http://localhost:8081/swagger/doc.json`
- **Files**: `docs/swagger.json` atau `docs/swagger.yaml`

### Import to Postman/Insomnia
1. Download `docs/swagger.json` atau `docs/swagger.yaml`
2. Import ke Postman atau Insomnia
3. Set environment variables (`base_url`, `access_token`)
4. Ready to test!

### Authentication Flow

1. **Login** → Get access token & refresh token
2. **Use Token** → Include in Authorization header
3. **Refresh** → Get new access token when expired
4. **Logout** → Revoke tokens

## API Endpoints

### 1. Login

```http
POST /auth/login
Content-Type: application/json

{
  "user_identity": "800000001",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "a3f485...",
    "token_type": "Bearer",
    "expires_in": 900,
    "user": {
      "id": 16,
      "name": "User Name",
      "email": "user@company.com",
      "modules": {...},
      "role_assignments": [...]
    }
  }
}
```

### 2. Refresh Token

```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "a3f485..."
}
```

### 3. Logout

```http
POST /auth/logout
Content-Type: application/json

{
  "token": "eyJhbGc..."
}
```

### 4. Get Profile

```http
GET /auth/profile?user_identity=800000001&application_code=APP001
Authorization: Bearer eyJhbGc...
```

## Token Information

- **Access Token**: JWT, expires in 15 minutes
- **Refresh Token**: Random string, expires in 7 days
- **Storage**: Redis with TTL

## Integration Examples

### JavaScript/Fetch

```javascript
// Login
const login = async (userIdentity, password) => {
  const response = await fetch('http://localhost:8081/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_identity: userIdentity, password })
  });
  
  const data = await response.json();
  return data.data.access_token;
};

// Use token
const fetchData = async (token) => {
  const response = await fetch('http://localhost:8081/api/v1/users', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  });
  
  return response.json();
};
```

### cURL

```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity":"800000001","password":"password123"}'

# Use token
curl -X GET http://localhost:8081/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Security Best Practices

### Token Storage

✅ **Recommended:**
- httpOnly cookies (most secure)
- sessionStorage for short-lived tokens
- Memory storage (lost on refresh)

❌ **Avoid:**
- localStorage (vulnerable to XSS)
- URL parameters
- Global variables

### Token Transmission

✅ **Always:**
- Use HTTPS in production
- Include token in Authorization header
- Validate SSL certificates

❌ **Never:**
- Send tokens in URL
- Use HTTP in production
- Ignore SSL errors

### CORS Configuration

**Production:**
```javascript
const corsOptions = {
  origin: ['https://yourdomain.com'],
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  allowedHeaders: ['Content-Type', 'Authorization']
};
```

**Development:**
```javascript
const corsOptions = {
  origin: ['http://localhost:3000'],
  credentials: true
};
```

## Error Handling

| Status | Error | Description |
|--------|-------|-------------|
| 400 | BAD_REQUEST | Invalid request format |
| 401 | UNAUTHORIZED | Invalid credentials or expired token |
| 403 | FORBIDDEN | Insufficient permissions |
| 429 | TOO_MANY_REQUESTS | Rate limit exceeded |
| 500 | INTERNAL_SERVER_ERROR | Server error |

## Rate Limiting

- Authentication endpoints: 10 requests/minute per IP
- Token refresh: 5 requests/minute per user
- Session check: 20 requests/minute per user

## Testing Credentials

| User Identity | Email | Password | Role |
|---------------|-------|----------|------|
| 800000001 | hasbi@company.com | password123 | CONSOLE ADMIN |
| 100000001 | naruto@company.com | password123 | User |

## Environment Variables

```bash
# .env
RBAC_SERVICE_URL=http://localhost:8081/api/v1
JWT_SECRET=your-secret-key
CORS_ORIGINS=http://localhost:3000
```

## Troubleshooting

**Token Expired:**
```json
{"success": false, "error": "token expired"}
```
→ Use refresh token to get new access token

**CORS Error:**
```
Access blocked by CORS policy
```
→ Add your domain to CORS_ORIGINS in .env

**Invalid Credentials:**
```json
{"success": false, "error": "kredensial tidak valid"}
```
→ Check user_identity and password

## Support

- Swagger UI: `http://localhost:8081/swagger/index.html`
- Documentation: `docs/`
- Issues: [GitHub Issues]
