# Authentication Integration Quick Start Guide

## Overview

Panduan ini akan membantu Anda mengintegrasikan authentication service dari Huminor RBAC System ke aplikasi project management dalam waktu 5 menit. RBAC service menyediakan authentication, authorization, dan user management yang lengkap.

## Prerequisites

- RBAC Service berjalan di `http://localhost:8082`
- Node.js atau bahasa pemrograman pilihan Anda
- Basic understanding tentang JWT tokens
- HTTP client library (axios, fetch, dll)

## Quick Setup (5 Menit)

### Step 1: Test RBAC Service Connection

```bash
# Test apakah RBAC service dapat diakses
curl http://localhost:8082/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}'
```

**Expected Response:**
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
      "role_assignments": [...]
    }
  }
}
```

### Step 2: Basic Frontend Integration (JavaScript)

```javascript
// auth-service.js
class AuthService {
  constructor() {
    this.baseURL = 'http://localhost:8082/api/v1';
    this.token = localStorage.getItem('access_token');
  }

  async login(userIdentity, password) {
    try {
      const response = await fetch(`${this.baseURL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_identity: userIdentity,
          password: password
        })
      });

      const data = await response.json();
      
      if (data.success) {
        this.token = data.data.access_token;
        localStorage.setItem('access_token', this.token);
        localStorage.setItem('refresh_token', data.data.refresh_token);
        localStorage.setItem('user_data', JSON.stringify(data.data.user));
        return data.data;
      } else {
        throw new Error(data.message);
      }
    } catch (error) {
      console.error('Login failed:', error);
      throw error;
    }
  }

  async logout() {
    try {
      const userData = JSON.parse(localStorage.getItem('user_data'));
      await fetch(`${this.baseURL}/auth/logout`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`
        },
        body: JSON.stringify({
          user_id: userData.id
        })
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      this.token = null;
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user_data');
    }
  }

  getAuthHeaders() {
    return {
      'Authorization': `Bearer ${this.token}`,
      'Content-Type': 'application/json'
    };
  }

  isAuthenticated() {
    return !!this.token;
  }

  getCurrentUser() {
    const userData = localStorage.getItem('user_data');
    return userData ? JSON.parse(userData) : null;
  }
}

// Usage example
const authService = new AuthService();

// Login
authService.login('800000001', 'password123')
  .then(userData => {
    console.log('Login successful:', userData);
    // Redirect to dashboard
  })
  .catch(error => {
    console.error('Login failed:', error);
  });
```

### Step 3: Basic Backend Integration (Node.js/Express)

```javascript
// middleware/auth.js
const jwt = require('jsonwebtoken');
const axios = require('axios');

const RBAC_SERVICE_URL = 'http://localhost:8082/api/v1';

// JWT Validation Middleware
const authenticateToken = async (req, res, next) => {
  const authHeader = req.headers['authorization'];
  const token = authHeader && authHeader.split(' ')[1];

  if (!token) {
    return res.status(401).json({
      success: false,
      message: 'Access token required'
    });
  }

  try {
    // Validate token with RBAC service
    const response = await axios.get(`${RBAC_SERVICE_URL}/auth/check-tokens`, {
      params: { user_id: getUserIdFromToken(token) },
      headers: { 'Authorization': `Bearer ${token}` }
    });

    if (response.data.success) {
      req.user = response.data.user;
      req.token = token;
      next();
    } else {
      return res.status(401).json({
        success: false,
        message: 'Invalid token'
      });
    }
  } catch (error) {
    return res.status(401).json({
      success: false,
      message: 'Token validation failed'
    });
  }
};

// Permission Check Middleware
const requirePermission = (moduleId, permission) => {
  return async (req, res, next) => {
    try {
      const response = await axios.post(`${RBAC_SERVICE_URL}/rbac/check-permission`, {
        user_id: req.user.id,
        module_id: moduleId,
        permission: permission
      }, {
        headers: { 'Authorization': `Bearer ${req.token}` }
      });

      if (response.data.hasPermission) {
        next();
      } else {
        return res.status(403).json({
          success: false,
          message: 'Insufficient permissions'
        });
      }
    } catch (error) {
      return res.status(500).json({
        success: false,
        message: 'Permission check failed'
      });
    }
  };
};

function getUserIdFromToken(token) {
  // Simple JWT decode (in production, use proper JWT library)
  const payload = JSON.parse(Buffer.from(token.split('.')[1], 'base64'));
  return payload.user_id;
}

module.exports = { authenticateToken, requirePermission };
```

### Step 4: Usage in Your Project Management API

```javascript
// app.js
const express = require('express');
const { authenticateToken, requirePermission } = require('./middleware/auth');

const app = express();
app.use(express.json());

// Public routes
app.post('/api/login', async (req, res) => {
  // Proxy to RBAC service
  try {
    const response = await axios.post(`${RBAC_SERVICE_URL}/auth/login`, req.body);
    res.json(response.data);
  } catch (error) {
    res.status(error.response?.status || 500).json({
      success: false,
      message: 'Login failed'
    });
  }
});

// Protected routes
app.use('/api/projects', authenticateToken);

// Get all projects (requires read permission)
app.get('/api/projects', requirePermission(200, 'read'), (req, res) => {
  // Your project logic here
  res.json({
    success: true,
    data: {
      projects: [],
      user: req.user
    }
  });
});

// Create project (requires write permission)
app.post('/api/projects', requirePermission(200, 'write'), (req, res) => {
  // Your project creation logic here
  res.json({
    success: true,
    message: 'Project created successfully'
  });
});

app.listen(3000, () => {
  console.log('Project Management API running on port 3000');
});
```

## Testing Your Integration

### 1. Test Login
```bash
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"user_identity": "800000001", "password": "password123"}'
```

### 2. Test Protected Endpoint
```bash
# Use the token from login response
curl -X GET http://localhost:3000/api/projects \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Common Issues & Solutions

### Issue 1: CORS Error
**Problem**: Browser blocks requests due to CORS policy

**Solution**: Configure CORS in your backend
```javascript
const cors = require('cors');
app.use(cors({
  origin: 'http://localhost:3000',
  credentials: true
}));
```

### Issue 2: Token Expired
**Problem**: 401 error after some time

**Solution**: Implement token refresh
```javascript
async function refreshToken() {
  const refreshToken = localStorage.getItem('refresh_token');
  const response = await fetch(`${baseURL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken })
  });
  
  if (response.ok) {
    const data = await response.json();
    localStorage.setItem('access_token', data.data.access_token);
    return data.data.access_token;
  }
  throw new Error('Token refresh failed');
}
```

### Issue 3: Permission Denied
**Problem**: 403 error when accessing resources

**Solution**: Check user roles and permissions
```javascript
// Get user permissions
const user = authService.getCurrentUser();
console.log('User roles:', user.role_assignments);

// Check specific permission before action
if (user.role_assignments.some(role => role.role_name === 'CONSOLE ADMIN')) {
  // User has admin access
} else {
  // Check specific module permissions
}
```

## Next Steps

1. **Read the Complete Guides**:
   - [Authentication API Reference](./AUTH_API_REFERENCE.md)
   - [Client Integration Guide](./CLIENT_INTEGRATION_GUIDE.md)
   - [Backend Integration Guide](./BACKEND_INTEGRATION_GUIDE.md)

2. **Implement Security Best Practices**:
   - [Security Best Practices Guide](./SECURITY_BEST_PRACTICES.md)

3. **Add RBAC Integration**:
   - [RBAC Integration Guide](./RBAC_INTEGRATION_GUIDE.md)

4. **Explore Examples**:
   - Check `/examples` folder for complete implementations
   - Import Postman collections for API testing

## Support

- **Documentation**: Check other guides in `/docs` folder
- **Examples**: See `/examples` folder for complete implementations
- **Issues**: Common problems and solutions in each guide
- **API Testing**: Use provided Postman collections

---

**Congratulations!** 🎉 You now have basic authentication integration working. Your project management system can authenticate users through the RBAC service and protect routes based on user permissions.