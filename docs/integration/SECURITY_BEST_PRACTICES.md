# Security Best Practices Guide

## Overview

Panduan ini menyediakan security best practices untuk mengintegrasikan Huminor RBAC Authentication System dengan aplikasi project management Anda. Mengikuti panduan ini akan memastikan implementasi yang aman dan sesuai dengan standar keamanan industri.

## JWT Token Security Guidelines

### 1. Token Storage

#### ✅ Recommended Approaches

**Frontend (Browser)**:
```javascript
// ✅ GOOD: Use httpOnly cookies (most secure)
// Set by server with httpOnly flag
document.cookie = "access_token=...; httpOnly; secure; sameSite=strict";

// ✅ ACCEPTABLE: sessionStorage for short-lived tokens
sessionStorage.setItem('access_token', token);

// ✅ ACCEPTABLE: Memory storage (lost on refresh)
class TokenManager {
  constructor() {
    this.accessToken = null;
    this.refreshToken = null;
  }
  
  setTokens(access, refresh) {
    this.accessToken = access;
    // Store refresh token in httpOnly cookie
    this.storeRefreshToken(refresh);
  }
}
```

#### ❌ Security Risks to Avoid

```javascript
// ❌ BAD: localStorage vulnerable to XSS
localStorage.setItem('access_token', token);

// ❌ BAD: Storing tokens in regular cookies without httpOnly
document.cookie = "access_token=" + token;

// ❌ BAD: Storing tokens in URL parameters
window.location.href = "/dashboard?token=" + token;

// ❌ BAD: Storing tokens in global variables
window.accessToken = token;
```

### 2. Token Transmission

#### ✅ Secure Token Transmission

```javascript
// ✅ GOOD: Always use HTTPS in production
const API_BASE_URL = process.env.NODE_ENV === 'production' 
  ? 'https://api.yourdomain.com' 
  : 'http://localhost:8082';

// ✅ GOOD: Proper Authorization header
const headers = {
  'Authorization': `Bearer ${accessToken}`,
  'Content-Type': 'application/json'
};

// ✅ GOOD: Validate SSL certificates
const httpsAgent = new https.Agent({
  rejectUnauthorized: true // Don't ignore SSL errors
});
```

#### ❌ Insecure Transmission

```javascript
// ❌ BAD: HTTP in production
const API_BASE_URL = 'http://api.yourdomain.com';

// ❌ BAD: Token in URL
fetch(`/api/data?token=${accessToken}`);

// ❌ BAD: Token in request body for GET requests
fetch('/api/data', {
  method: 'GET',
  body: JSON.stringify({ token: accessToken })
});
```

### 3. Token Validation

#### Backend Token Validation

```javascript
// ✅ GOOD: Comprehensive token validation
const validateToken = async (token) => {
  try {
    // 1. Check token format
    if (!token || !token.startsWith('Bearer ')) {
      throw new Error('Invalid token format');
    }
    
    const actualToken = token.substring(7);
    
    // 2. Validate with RBAC service
    const response = await axios.get(`${RBAC_SERVICE_URL}/auth/check-tokens`, {
      params: { user_id: getUserIdFromToken(actualToken) },
      headers: { 'Authorization': token },
      timeout: 5000 // Prevent hanging requests
    });
    
    // 3. Check response validity
    if (!response.data.success) {
      throw new Error('Token validation failed');
    }
    
    // 4. Check token expiration
    const tokenData = parseJWT(actualToken);
    if (tokenData.exp < Date.now() / 1000) {
      throw new Error('Token expired');
    }
    
    return response.data;
  } catch (error) {
    console.error('Token validation error:', error);
    throw new Error('Invalid token');
  }
};
```

## CORS Configuration Examples

### 1. Production CORS Setup

#### Node.js/Express
```javascript
// ✅ GOOD: Restrictive CORS for production
const cors = require('cors');

const corsOptions = {
  origin: function (origin, callback) {
    // Allow requests from your domains only
    const allowedOrigins = [
      'https://yourdomain.com',
      'https://app.yourdomain.com',
      'https://admin.yourdomain.com'
    ];
    
    // Allow requests with no origin (mobile apps, etc.)
    if (!origin) return callback(null, true);
    
    if (allowedOrigins.indexOf(origin) !== -1) {
      callback(null, true);
    } else {
      callback(new Error('Not allowed by CORS'));
    }
  },
  credentials: true, // Allow cookies
  methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowedHeaders: ['Content-Type', 'Authorization'],
  maxAge: 86400 // Cache preflight for 24 hours
};

app.use(cors(corsOptions));
```

#### Go/Gin
```go
// ✅ GOOD: Secure CORS configuration
func CORSMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        allowedOrigins := []string{
            "https://yourdomain.com",
            "https://app.yourdomain.com",
        }
        
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin {
                c.Header("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Max-Age", "86400")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
}
```

### 2. Development CORS Setup

```javascript
// ✅ ACCEPTABLE: Relaxed CORS for development only
if (process.env.NODE_ENV === 'development') {
  app.use(cors({
    origin: ['http://localhost:3000', 'http://localhost:3001'],
    credentials: true
  }));
}
```

## Environment-Specific Security Settings

### 1. Environment Configuration

#### .env.production
```bash
# ✅ GOOD: Production environment variables
NODE_ENV=production
RBAC_SERVICE_URL=https://api.yourdomain.com/api/v1
JWT_SECRET=your-super-secure-secret-key-here
CORS_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
SESSION_SECRET=another-super-secure-secret
SECURE_COOKIES=true
HTTPS_ONLY=true
RATE_LIMIT_WINDOW=900000
RATE_LIMIT_MAX=100
```

#### .env.development
```bash
# ✅ ACCEPTABLE: Development environment variables
NODE_ENV=development
RBAC_SERVICE_URL=http://localhost:8082/api/v1
JWT_SECRET=dev-secret-key
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
SESSION_SECRET=dev-session-secret
SECURE_COOKIES=false
HTTPS_ONLY=false
RATE_LIMIT_WINDOW=60000
RATE_LIMIT_MAX=1000
```

### 2. Security Headers

```javascript
// ✅ GOOD: Security headers middleware
const helmet = require('helmet');

app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"],
      scriptSrc: ["'self'"],
      imgSrc: ["'self'", "data:", "https:"],
      connectSrc: ["'self'", process.env.RBAC_SERVICE_URL],
      fontSrc: ["'self'"],
      objectSrc: ["'none'"],
      mediaSrc: ["'self'"],
      frameSrc: ["'none'"],
    },
  },
  hsts: {
    maxAge: 31536000,
    includeSubDomains: true,
    preload: true
  }
}));

// Additional security headers
app.use((req, res, next) => {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('X-XSS-Protection', '1; mode=block');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  next();
});
```

## Common Security Vulnerabilities

### 1. Cross-Site Scripting (XSS) Prevention

#### ✅ XSS Prevention Techniques

```javascript
// ✅ GOOD: Input sanitization
const DOMPurify = require('dompurify');
const { JSDOM } = require('jsdom');

const window = new JSDOM('').window;
const purify = DOMPurify(window);

const sanitizeInput = (input) => {
  return purify.sanitize(input);
};

// ✅ GOOD: Output encoding
const escapeHtml = (unsafe) => {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
};

// ✅ GOOD: React JSX automatically escapes
const UserProfile = ({ user }) => (
  <div>
    <h1>{user.name}</h1> {/* Automatically escaped */}
    <div dangerouslySetInnerHTML={{
      __html: purify.sanitize(user.bio) // Sanitize HTML content
    }} />
  </div>
);
```

### 2. Cross-Site Request Forgery (CSRF) Prevention

```javascript
// ✅ GOOD: CSRF protection
const csrf = require('csurf');

// Configure CSRF protection
const csrfProtection = csrf({
  cookie: {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'strict'
  }
});

app.use(csrfProtection);

// Provide CSRF token to frontend
app.get('/api/csrf-token', (req, res) => {
  res.json({ csrfToken: req.csrfToken() });
});

// Frontend usage
const getCsrfToken = async () => {
  const response = await fetch('/api/csrf-token');
  const data = await response.json();
  return data.csrfToken;
};

const makeSecureRequest = async (url, data) => {
  const csrfToken = await getCsrfToken();
  
  return fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken,
      'Authorization': `Bearer ${accessToken}`
    },
    body: JSON.stringify(data)
  });
};
```

### 3. SQL Injection Prevention

```javascript
// ✅ GOOD: Parameterized queries
const getUserProjects = async (userId) => {
  // Using parameterized query
  const query = 'SELECT * FROM projects WHERE user_id = ? AND deleted_at IS NULL';
  const results = await db.query(query, [userId]);
  return results;
};

// ✅ GOOD: ORM usage
const getUserProjects = async (userId) => {
  return await Project.findAll({
    where: {
      user_id: userId,
      deleted_at: null
    }
  });
};

// ❌ BAD: String concatenation
const getUserProjects = async (userId) => {
  const query = `SELECT * FROM projects WHERE user_id = ${userId}`;
  return await db.query(query);
};
```

## Security Checklist for Production

### Pre-Deployment Security Checklist

#### ✅ Authentication & Authorization
- [ ] JWT tokens expire within 15 minutes
- [ ] Refresh tokens expire within 7 days
- [ ] Token rotation implemented
- [ ] Proper logout functionality (token revocation)
- [ ] Rate limiting on authentication endpoints
- [ ] Account lockout after failed attempts
- [ ] Strong password requirements enforced

#### ✅ Data Protection
- [ ] All API calls use HTTPS
- [ ] Sensitive data encrypted at rest
- [ ] Database connections encrypted
- [ ] No sensitive data in logs
- [ ] Proper input validation and sanitization
- [ ] Output encoding implemented

#### ✅ Network Security
- [ ] CORS properly configured
- [ ] Security headers implemented
- [ ] CSP (Content Security Policy) configured
- [ ] HSTS enabled
- [ ] No unnecessary ports exposed

#### ✅ Application Security
- [ ] Dependencies updated and scanned
- [ ] No hardcoded secrets
- [ ] Environment variables secured
- [ ] Error messages don't leak information
- [ ] File upload restrictions implemented
- [ ] Session management secure

#### ✅ Infrastructure Security
- [ ] Firewall rules configured
- [ ] Database access restricted
- [ ] Redis access secured
- [ ] Monitoring and alerting setup
- [ ] Backup and recovery tested
- [ ] SSL certificates valid

### Security Monitoring

```javascript
// ✅ GOOD: Security event logging
const securityLogger = require('./security-logger');

const logSecurityEvent = (event, details) => {
  securityLogger.warn({
    event,
    timestamp: new Date().toISOString(),
    ip: details.ip,
    userAgent: details.userAgent,
    userId: details.userId,
    details
  });
};

// Usage examples
app.use('/api', (req, res, next) => {
  // Log suspicious activity
  if (req.headers['user-agent'].includes('bot')) {
    logSecurityEvent('SUSPICIOUS_USER_AGENT', {
      ip: req.ip,
      userAgent: req.headers['user-agent'],
      path: req.path
    });
  }
  
  next();
});

// Failed login attempts
const handleFailedLogin = (userIdentity, ip, userAgent) => {
  logSecurityEvent('FAILED_LOGIN', {
    userIdentity,
    ip,
    userAgent,
    timestamp: new Date().toISOString()
  });
};
```

## Incident Response Plan

### 1. Security Incident Detection

```javascript
// ✅ GOOD: Automated threat detection
const detectThreats = (req) => {
  const threats = [];
  
  // Multiple failed login attempts
  if (getFailedLoginCount(req.ip) > 5) {
    threats.push('BRUTE_FORCE_ATTACK');
  }
  
  // Suspicious user agent
  if (req.headers['user-agent'].match(/sqlmap|nikto|nmap/i)) {
    threats.push('SCANNING_ATTEMPT');
  }
  
  // Rate limit exceeded
  if (getRateLimit(req.ip) > 100) {
    threats.push('RATE_LIMIT_EXCEEDED');
  }
  
  return threats;
};
```

### 2. Incident Response Actions

```javascript
// ✅ GOOD: Automated response system
const respondToThreat = async (threat, details) => {
  switch (threat) {
    case 'BRUTE_FORCE_ATTACK':
      await blockIP(details.ip, '1 hour');
      await notifySecurityTeam('Brute force attack detected', details);
      break;
      
    case 'TOKEN_COMPROMISE':
      await revokeAllUserTokens(details.userId);
      await forcePasswordReset(details.userId);
      await notifyUser(details.userId, 'Security alert: Please change your password');
      break;
      
    case 'DATA_BREACH':
      await enableMaintenanceMode();
      await notifySecurityTeam('URGENT: Data breach detected', details);
      await backupCurrentState();
      break;
  }
};
```

## Security Testing Guidelines

### 1. Automated Security Testing

```javascript
// ✅ GOOD: Security test suite
describe('Security Tests', () => {
  test('should reject requests without authentication', async () => {
    const response = await request(app)
      .get('/api/projects')
      .expect(401);
    
    expect(response.body.success).toBe(false);
  });
  
  test('should reject expired tokens', async () => {
    const expiredToken = generateExpiredToken();
    
    const response = await request(app)
      .get('/api/projects')
      .set('Authorization', `Bearer ${expiredToken}`)
      .expect(401);
  });
  
  test('should prevent XSS attacks', async () => {
    const maliciousInput = '<script>alert("xss")</script>';
    
    const response = await request(app)
      .post('/api/projects')
      .set('Authorization', `Bearer ${validToken}`)
      .send({ name: maliciousInput })
      .expect(400);
  });
  
  test('should enforce rate limits', async () => {
    // Make multiple requests rapidly
    const promises = Array(101).fill().map(() =>
      request(app).post('/api/auth/login').send(validCredentials)
    );
    
    const responses = await Promise.all(promises);
    const rateLimitedResponses = responses.filter(r => r.status === 429);
    
    expect(rateLimitedResponses.length).toBeGreaterThan(0);
  });
});
```

### 2. Manual Security Testing

#### Penetration Testing Checklist
- [ ] Authentication bypass attempts
- [ ] Authorization escalation tests
- [ ] Input validation testing
- [ ] Session management testing
- [ ] CSRF protection testing
- [ ] XSS vulnerability scanning
- [ ] SQL injection testing
- [ ] File upload security testing

## Security Resources

### Tools and Libraries

#### Frontend Security
- **DOMPurify**: HTML sanitization
- **helmet.js**: Security headers
- **csurf**: CSRF protection
- **express-rate-limit**: Rate limiting

#### Backend Security
- **bcrypt**: Password hashing
- **jsonwebtoken**: JWT handling
- **express-validator**: Input validation
- **node-rate-limiter-flexible**: Advanced rate limiting

#### Security Scanning
- **npm audit**: Dependency vulnerability scanning
- **Snyk**: Continuous security monitoring
- **OWASP ZAP**: Web application security testing
- **Burp Suite**: Professional security testing

### Security Standards

- **OWASP Top 10**: Web application security risks
- **NIST Cybersecurity Framework**: Comprehensive security guidelines
- **ISO 27001**: Information security management
- **PCI DSS**: Payment card industry standards (if applicable)

### Regular Security Maintenance

#### Monthly Tasks
- [ ] Update all dependencies
- [ ] Review security logs
- [ ] Test backup and recovery
- [ ] Update security documentation

#### Quarterly Tasks
- [ ] Penetration testing
- [ ] Security training for team
- [ ] Review and update security policies
- [ ] Audit user access and permissions

#### Annual Tasks
- [ ] Comprehensive security audit
- [ ] Disaster recovery testing
- [ ] Security certification renewal
- [ ] Third-party security assessment

---

**Remember**: Security is an ongoing process, not a one-time setup. Regularly review and update your security measures as threats evolve and your application grows.

**Next Steps**:
- [Client Integration Guide](./CLIENT_INTEGRATION_GUIDE.md) - Frontend security implementation
- [Backend Integration Guide](./BACKEND_INTEGRATION_GUIDE.md) - Server-side security
- [RBAC Integration Guide](./RBAC_INTEGRATION_GUIDE.md) - Permission-based security