# API Documentation System

## Overview

Sistem dokumentasi API yang memungkinkan pengelolaan dokumentasi API secara terpusat dengan kemampuan export ke berbagai format client seperti Postman, Insomnia, dan Apidog.

## Tujuan

1. **Centralized Documentation**: Menyimpan semua dokumentasi API dalam database
2. **Multi-format Export**: Generate file yang bisa diimport ke berbagai API client
3. **Auto-Documentation**: Otomatis mendokumentasikan API endpoints yang ada
4. **Environment Management**: Mendukung multiple environment (dev, staging, prod)
5. **Export Ready**: Siap export ke format yang dibutuhkan tim

## Fitur Utama

### 1. Collection Management
- Membuat dan mengelola koleksi API
- Versioning collection
- Multi-tenant support (per company)

### 2. Folder Organization
- Hierarchical folder structure
- Nested folders untuk organisasi yang lebih baik
- Drag & drop reordering

### 3. Endpoint Documentation
- Complete HTTP method support (GET, POST, PUT, PATCH, DELETE, etc.)
- URL parameters dan query parameters
- Request headers dan response headers
- Request body dengan schema validation
- Multiple response examples dengan status codes

### 4. Environment Variables
- Multiple environments per collection
- Secure variable storage
- Variable interpolation dalam URLs dan bodies

### 5. Testing & Scripts
- Pre-request scripts
- Test scripts untuk validation
- Environment variable manipulation

### 6. Export Formats
- Postman Collection v2.1
- Insomnia Collection
- OpenAPI 3.0 Specification
- Swagger JSON
- Apidog Collection

### 7. Auto-Discovery
- Scan existing routes untuk auto-populate endpoints
- Extract parameter information dari route handlers
- Generate basic documentation dari code comments

## Database Schema

### Core Tables

#### 1. api_collections
Menyimpan informasi koleksi API utama.

```sql
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 2. api_folders
Untuk organisasi endpoint dalam folder hierarkis.

```sql
CREATE TABLE api_folders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    collection_id BIGINT NOT NULL,
    parent_id BIGINT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 3. api_endpoints
Menyimpan informasi endpoint API.

```sql
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Supporting Tables

#### 4. api_headers
```sql
CREATE TABLE api_headers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    header_type ENUM('request', 'response') DEFAULT 'request',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 5. api_parameters
```sql
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 6. api_request_bodies
```sql
CREATE TABLE api_request_bodies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    content_type VARCHAR(100) DEFAULT 'application/json',
    body_content TEXT,
    description TEXT,
    schema_definition JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 7. api_responses
```sql
CREATE TABLE api_responses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    status_code INT NOT NULL,
    status_text VARCHAR(100),
    content_type VARCHAR(100) DEFAULT 'application/json',
    response_body TEXT,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 8. api_environments & api_environment_variables
```sql
CREATE TABLE api_environments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    collection_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE api_environment_variables (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    environment_id BIGINT NOT NULL,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_secret BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 9. api_tests
```sql
CREATE TABLE api_tests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    endpoint_id BIGINT NOT NULL,
    test_type ENUM('pre_request', 'test') NOT NULL,
    script_content TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 10. api_tags & api_endpoint_tags
```sql
CREATE TABLE api_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#007bff',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE api_endpoint_tags (
    endpoint_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (endpoint_id, tag_id)
);
```

## API Endpoints

### Collection Management
- `GET /api/v1/api-docs/collections` - List all collections
- `POST /api/v1/api-docs/collections` - Create new collection
- `GET /api/v1/api-docs/collections/{id}` - Get collection details
- `PUT /api/v1/api-docs/collections/{id}` - Update collection
- `DELETE /api/v1/api-docs/collections/{id}` - Delete collection

### Folder Management
- `GET /api/v1/api-docs/collections/{id}/folders` - List folders in collection
- `POST /api/v1/api-docs/collections/{id}/folders` - Create folder
- `PUT /api/v1/api-docs/folders/{id}` - Update folder
- `DELETE /api/v1/api-docs/folders/{id}` - Delete folder
- `POST /api/v1/api-docs/folders/{id}/reorder` - Reorder folders

### Endpoint Management
- `GET /api/v1/api-docs/collections/{id}/endpoints` - List endpoints
- `POST /api/v1/api-docs/collections/{id}/endpoints` - Create endpoint
- `GET /api/v1/api-docs/endpoints/{id}` - Get endpoint details
- `PUT /api/v1/api-docs/endpoints/{id}` - Update endpoint
- `DELETE /api/v1/api-docs/endpoints/{id}` - Delete endpoint

### Environment Management
- `GET /api/v1/api-docs/collections/{id}/environments` - List environments
- `POST /api/v1/api-docs/collections/{id}/environments` - Create environment
- `PUT /api/v1/api-docs/environments/{id}` - Update environment
- `DELETE /api/v1/api-docs/environments/{id}` - Delete environment

### Export Endpoints
- `GET /api/v1/api-docs/collections/{id}/export/postman` - Export to Postman
- `GET /api/v1/api-docs/collections/{id}/export/insomnia` - Export to Insomnia
- `GET /api/v1/api-docs/collections/{id}/export/openapi` - Export to OpenAPI
- `GET /api/v1/api-docs/collections/{id}/export/swagger` - Export to Swagger
- `GET /api/v1/api-docs/collections/{id}/export/apidog` - Export to Apidog

### Auto-Discovery Endpoints
- `POST /api/v1/api-docs/collections/{id}/scan-routes` - Scan dan populate endpoints dari routes
- `GET /api/v1/api-docs/routes/available` - List available routes untuk scanning

## Implementation Plan

### Phase 1: Core Infrastructure
1. Database migration files
2. Base models dan repositories
3. Basic CRUD services
4. Authentication dan authorization

### Phase 2: API Endpoints
1. Collection management endpoints
2. Folder management endpoints
3. Endpoint management endpoints
4. Environment management endpoints

### Phase 3: Export Functionality
1. Postman collection export
2. OpenAPI specification export
3. Insomnia collection export
4. Swagger JSON export
5. Apidog collection export

### Phase 4: Auto-Discovery
1. Route scanning functionality
2. Parameter extraction dari Gin routes
3. Auto-populate basic endpoint information
4. Integration dengan existing route structure

### Phase 5: Advanced Features
1. Search dan filtering
2. Tagging system
3. Bulk operations
4. Enhanced documentation
5. Performance optimization

## File Structure

```
internal/modules/apidoc/
├── dto.go              # Data Transfer Objects
├── model.go            # Database models
├── repository.go       # Database operations
├── service.go          # Business logic
├── route.go            # HTTP routes
├── scanner.go          # Route scanning logic
├── export/
│   ├── postman.go      # Postman export logic
│   ├── openapi.go      # OpenAPI export logic
│   ├── insomnia.go     # Insomnia export logic
│   ├── swagger.go      # Swagger export logic
│   └── apidog.go       # Apidog export logic
└── utils/
    ├── parser.go       # Route parsing utilities
    └── validator.go    # Validation utilities
```

## Security Considerations

1. **Authentication**: Semua endpoint memerlukan authentication
2. **Authorization**: Role-based access control
3. **Data Validation**: Input validation untuk semua endpoints
4. **Rate Limiting**: Prevent abuse pada export endpoints
5. **Audit Logging**: Log semua operasi CRUD
6. **Sensitive Data**: Encrypt environment variables yang sensitive

## Performance Considerations

1. **Caching**: Cache hasil export untuk collection yang tidak berubah
2. **Pagination**: Implement pagination untuk list endpoints
3. **Indexing**: Database indexes untuk query performance
4. **Async Processing**: Background jobs untuk export file besar
5. **File Storage**: Store exported files di cloud storage

## Testing Strategy

1. **Unit Tests**: Test semua service methods
2. **Integration Tests**: Test API endpoints
3. **Export Tests**: Validate exported file formats
4. **Import Tests**: Test import functionality
5. **Performance Tests**: Load testing untuk export endpoints

## Monitoring & Logging

1. **Metrics**: Track export/import operations
2. **Error Logging**: Comprehensive error logging
3. **Performance Monitoring**: Monitor response times
4. **Usage Analytics**: Track popular collections dan endpoints

## Future Enhancements

1. **Advanced Route Scanning**: Deep analysis dari Gin middleware dan handlers
2. **Code Comment Parsing**: Extract documentation dari Go comments
3. **Real-time Updates**: Auto-update documentation saat code berubah
4. **API Testing**: Built-in API testing capabilities
5. **Performance Monitoring**: Integration dengan monitoring tools