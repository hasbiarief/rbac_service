# API Documentation System - Migration Plan

## Database Migration Strategy

### Migration Files Sequence

#### 1. Migration 012: Core Tables
File: `migrations/012_create_api_documentation_tables.sql`

```sql
-- Create core API documentation tables
CREATE TABLE api_collections (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) DEFAULT '1.0.0',
    base_url VARCHAR(500),
    schema_version VARCHAR(50) DEFAULT 'v2.1.0',
    created_by BIGINT,
    company_id BIGINT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_company_id (company_id),
    INDEX idx_created_by (created_by),
    INDEX idx_is_active (is_active),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE TABLE api_folders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    collection_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (collection_id) REFERENCES api_collections(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES api_folders(id) ON DELETE CASCADE,
    INDEX idx_collection_id (collection_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort_order (sort_order)
);

CREATE TABLE api_endpoints (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    collection_id BIGINT NOT NULL,
    folder_id BIGINT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    method ENUM('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS') NOT NULL,
    url VARCHAR(1000) NOT NULL,
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (collection_id) REFERENCES api_collections(id) ON DELETE CASCADE,
    FOREIGN KEY (folder_id) REFERENCES api_folders(id) ON DELETE SET NULL,
    INDEX idx_collection_id (collection_id),
    INDEX idx_folder_id (folder_id),
    INDEX idx_method (method),
    INDEX idx_sort_order (sort_order)
);
```

#### 2. Migration 013: Supporting Tables
File: `migrations/013_create_api_documentation_details.sql`

```sql
-- Create supporting tables for API documentation details
CREATE TABLE api_headers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    header_type ENUM('request', 'response') DEFAULT 'request',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    INDEX idx_endpoint_id (endpoint_id),
    INDEX idx_header_type (header_type)
);

CREATE TABLE api_parameters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type ENUM('query', 'path', 'form', 'file') NOT NULL,
    data_type VARCHAR(50) DEFAULT 'string',
    description TEXT,
    default_value TEXT,
    example_value TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    INDEX idx_endpoint_id (endpoint_id),
    INDEX idx_type (type)
);

CREATE TABLE api_request_bodies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    content_type VARCHAR(100) DEFAULT 'application/json',
    body_content TEXT,
    description TEXT,
    schema_definition JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    INDEX idx_endpoint_id (endpoint_id),
    INDEX idx_content_type (content_type)
);

CREATE TABLE api_responses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    status_code INT NOT NULL,
    status_text VARCHAR(100),
    content_type VARCHAR(100) DEFAULT 'application/json',
    response_body TEXT,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    INDEX idx_endpoint_id (endpoint_id),
    INDEX idx_status_code (status_code),
    INDEX idx_is_default (is_default)
);
```

#### 3. Migration 014: Environment & Testing
File: `migrations/014_create_api_environments_and_tests.sql`

```sql
-- Create environment and testing tables
CREATE TABLE api_environments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    collection_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (collection_id) REFERENCES api_collections(id) ON DELETE CASCADE,
    INDEX idx_collection_id (collection_id),
    INDEX idx_is_default (is_default),
    UNIQUE KEY unique_collection_name (collection_id, name)
);

CREATE TABLE api_environment_variables (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    environment_id BIGINT NOT NULL,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_secret BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (environment_id) REFERENCES api_environments(id) ON DELETE CASCADE,
    INDEX idx_environment_id (environment_id),
    INDEX idx_is_secret (is_secret),
    UNIQUE KEY unique_env_key (environment_id, key_name)
);

CREATE TABLE api_tests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    test_type ENUM('pre_request', 'test') NOT NULL,
    script_content TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    INDEX idx_endpoint_id (endpoint_id),
    INDEX idx_test_type (test_type)
);
```

#### 4. Migration 015: Tagging System
File: `migrations/015_create_api_tagging_system.sql`

```sql
-- Create tagging system for API documentation
CREATE TABLE api_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#007bff',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_name (name)
);

CREATE TABLE api_endpoint_tags (
    endpoint_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (endpoint_id, tag_id),
    FOREIGN KEY (endpoint_id) REFERENCES api_endpoints(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES api_tags(id) ON DELETE CASCADE
);

-- Insert default tags
INSERT INTO api_tags (name, color, description) VALUES
('Authentication', '#dc3545', 'Authentication related endpoints'),
('CRUD', '#28a745', 'Create, Read, Update, Delete operations'),
('Admin', '#ffc107', 'Administrative endpoints'),
('Public', '#17a2b8', 'Public accessible endpoints'),
('Deprecated', '#6c757d', 'Deprecated endpoints');
```

#### 5. Migration 016: Seed Initial Data
File: `migrations/016_seed_api_documentation_data.sql`

```sql
-- Seed initial API documentation data
INSERT INTO api_collections (name, description, version, base_url, created_by, company_id) VALUES
('Huminor RBAC API', 'Complete API Collection for ERP with Unit-Based Role Access Control (RBAC)', '1.0.0', '{{base_url}}', 1, 1);

SET @collection_id = LAST_INSERT_ID();

-- Create default environment
INSERT INTO api_environments (collection_id, name, description, is_default) VALUES
(@collection_id, 'Development', 'Development environment', TRUE);

SET @env_id = LAST_INSERT_ID();

-- Add environment variables
INSERT INTO api_environment_variables (environment_id, key_name, value, description) VALUES
(@env_id, 'base_url', 'http://localhost:8080', 'Base URL for API'),
(@env_id, 'accessToken', '', 'JWT Access Token'),
(@env_id, 'refreshToken', '', 'JWT Refresh Token'),
(@env_id, 'userIdentity', 'admin', 'User Identity for login'),
(@env_id, 'userPassword', 'password123', 'User Password for login'),
(@env_id, 'userId', '', 'Current User ID');

-- Create folders
INSERT INTO api_folders (collection_id, name, description, sort_order) VALUES
(@collection_id, 'Authentication', 'Authentication related endpoints', 1),
(@collection_id, 'Module System', 'Module management endpoints', 2),
(@collection_id, 'User Management', 'User management endpoints', 3),
(@collection_id, 'Role Management', 'Role management endpoints', 4),
(@collection_id, 'Company Management', 'Company management endpoints', 5),
(@collection_id, 'Branch Management', 'Branch management endpoints', 6),
(@collection_id, 'Unit Management', 'Unit management endpoints', 7),
(@collection_id, 'Audit Logging', 'Audit logging endpoints', 8);
```

## Implementation Timeline

### Week 1: Database Setup
- [ ] Create migration files (012-016)
- [ ] Run migrations pada development environment
- [ ] Verify database schema
- [ ] Create seed data

### Week 2: Core Models & Repositories
- [ ] Create model structs
- [ ] Implement repository interfaces
- [ ] Create base CRUD operations
- [ ] Write unit tests untuk repositories

### Week 3: Service Layer
- [ ] Implement business logic services
- [ ] Add validation rules
- [ ] Implement authorization checks
- [ ] Write service unit tests

### Week 4: API Endpoints - Basic CRUD
- [ ] Collection management endpoints
- [ ] Folder management endpoints
- [ ] Endpoint management endpoints
- [ ] Environment management endpoints

### Week 5: API Endpoints - Advanced Features
- [ ] Parameter management
- [ ] Header management
- [ ] Request/Response body management
- [ ] Test script management

### Week 6: Export Functionality - Phase 1
- [ ] Postman collection export
- [ ] Basic OpenAPI export
- [ ] File download endpoints
- [ ] Export validation

### Week 7: Export Functionality - Phase 2
- [ ] Insomnia collection export
- [ ] Swagger JSON export
- [ ] Apidog collection export
- [ ] Export optimization

### Week 8: Auto-Discovery Features
- [ ] Route scanning functionality
- [ ] Parameter extraction
- [ ] Auto-populate endpoints
- [ ] Integration testing

### Week 9: Testing & Documentation
- [ ] Integration tests
- [ ] API documentation
- [ ] Performance testing
- [ ] Security testing

### Week 10: Deployment & Monitoring
- [ ] Production deployment
- [ ] Monitoring setup
- [ ] Performance optimization
- [ ] Bug fixes

## Rollback Strategy

### Database Rollback
```sql
-- Rollback migration 016
DELETE FROM api_environment_variables;
DELETE FROM api_environments;
DELETE FROM api_folders;
DELETE FROM api_collections;
DELETE FROM api_tags WHERE name IN ('Authentication', 'CRUD', 'Admin', 'Public', 'Deprecated');

-- Rollback migration 015
DROP TABLE api_endpoint_tags;
DROP TABLE api_tags;

-- Rollback migration 014
DROP TABLE api_tests;
DROP TABLE api_environment_variables;
DROP TABLE api_environments;

-- Rollback migration 013
DROP TABLE api_responses;
DROP TABLE api_request_bodies;
DROP TABLE api_parameters;
DROP TABLE api_headers;

-- Rollback migration 012
DROP TABLE api_endpoints;
DROP TABLE api_folders;
DROP TABLE api_collections;
```

### Code Rollback
1. Remove API documentation module dari routes
2. Remove module folder
3. Revert database migrations
4. Update documentation

## Risk Assessment

### High Risk
- **Data Loss**: Backup database sebelum migration
- **Performance Impact**: Monitor query performance
- **Breaking Changes**: Ensure backward compatibility

### Medium Risk
- **Export Format Compatibility**: Test dengan berbagai client versions
- **Large File Handling**: Implement streaming untuk file besar
- **Concurrent Access**: Handle concurrent editing

### Low Risk
- **UI Changes**: Frontend changes tidak mempengaruhi API
- **Minor Bug Fixes**: Easy to fix dan deploy
- **Documentation Updates**: Non-breaking changes

## Success Metrics

### Technical Metrics
- [ ] All migrations run successfully
- [ ] All tests pass (>95% coverage)
- [ ] API response time < 200ms
- [ ] Export file generation < 5 seconds

### Business Metrics
- [ ] User adoption rate > 80%
- [ ] Export success rate > 99%
- [ ] User satisfaction score > 4.5/5
- [ ] Reduction in documentation maintenance time > 50%

## Post-Implementation Tasks

### Week 11-12: Monitoring & Optimization
- [ ] Monitor system performance
- [ ] Optimize slow queries
- [ ] Fix reported bugs
- [ ] Gather user feedback

### Week 13-14: Enhancement Planning
- [ ] Plan advanced route scanning
- [ ] Analyze usage patterns
- [ ] Plan auto-documentation features
- [ ] Roadmap untuk integration features