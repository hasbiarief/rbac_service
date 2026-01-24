# API Documentation Export Formats

## Overview

Dokumen ini menjelaskan format export yang didukung oleh sistem API Documentation dan struktur data yang dihasilkan untuk setiap format.

## Supported Export Formats

### 1. Postman Collection v2.1

#### Format Structure
```json
{
  "info": {
    "_postman_id": "uuid",
    "name": "Collection Name",
    "description": "Collection Description",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
    "_exporter_id": "system_id"
  },
  "item": [
    {
      "name": "Folder Name",
      "item": [
        {
          "name": "Endpoint Name",
          "event": [
            {
              "listen": "prerequest",
              "script": {
                "exec": ["script content"],
                "type": "text/javascript"
              }
            },
            {
              "listen": "test",
              "script": {
                "exec": ["test script"],
                "type": "text/javascript"
              }
            }
          ],
          "request": {
            "method": "GET|POST|PUT|DELETE|PATCH",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{accessToken}}",
                "type": "text"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"key\": \"value\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/api/v1/endpoint",
              "host": ["{{base_url}}"],
              "path": ["api", "v1", "endpoint"],
              "query": [
                {
                  "key": "param",
                  "value": "value"
                }
              ]
            }
          },
          "response": [
            {
              "name": "Success Response",
              "originalRequest": {},
              "status": "OK",
              "code": 200,
              "header": [],
              "body": "{\n  \"success\": true\n}"
            }
          ]
        }
      ]
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080",
      "type": "string"
    }
  ]
}
```

#### Export Endpoint
```
GET /api/v1/api-docs/collections/{id}/export/postman
```

#### Response Headers
```
Content-Type: application/json
Content-Disposition: attachment; filename="collection_name.postman_collection.json"
```

### 2. OpenAPI 3.0 Specification

#### Format Structure
```yaml
openapi: 3.0.0
info:
  title: API Title
  description: API Description
  version: 1.0.0
  contact:
    name: API Support
    email: support@example.com
servers:
  - url: http://localhost:8080
    description: Development server
  - url: https://api.example.com
    description: Production server
paths:
  /api/v1/endpoint:
    get:
      summary: Endpoint Summary
      description: Endpoint Description
      tags:
        - Tag Name
      parameters:
        - name: param
          in: query
          required: false
          schema:
            type: string
          description: Parameter description
      responses:
        '200':
          description: Success response
          content:
            application/json:
              schema:
                type: object
                properties:
                  success:
                    type: boolean
                  data:
                    type: object
              example:
                success: true
                data: {}
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    post:
      summary: Create endpoint
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                  description: Name field
              required:
                - name
            example:
              name: "Example"
      responses:
        '201':
          description: Created successfully
components:
  schemas:
    Error:
      type: object
      properties:
        success:
          type: boolean
          example: false
        message:
          type: string
          example: "Error message"
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
security:
  - BearerAuth: []
```

#### Export Endpoints
```
GET /api/v1/api-docs/collections/{id}/export/openapi?format=json
GET /api/v1/api-docs/collections/{id}/export/openapi?format=yaml
```

### 3. Insomnia Collection

#### Format Structure
```json
{
  "_type": "export",
  "__export_format": 4,
  "__export_date": "2024-01-01T00:00:00.000Z",
  "__export_source": "insomnia.desktop.app:v2023.5.8",
  "resources": [
    {
      "_id": "req_123",
      "_type": "request",
      "parentId": "fld_123",
      "modified": 1640995200000,
      "created": 1640995200000,
      "url": "{{ _.base_url }}/api/v1/endpoint",
      "name": "Endpoint Name",
      "description": "Endpoint Description",
      "method": "GET",
      "body": {
        "mimeType": "application/json",
        "text": "{\n  \"key\": \"value\"\n}"
      },
      "parameters": [
        {
          "name": "param",
          "value": "value",
          "description": "Parameter description",
          "disabled": false
        }
      ],
      "headers": [
        {
          "name": "Authorization",
          "value": "Bearer {{ _.accessToken }}",
          "description": "JWT Token"
        }
      ],
      "authentication": {},
      "metaSortKey": -1640995200000,
      "isPrivate": false,
      "settingStoreCookies": true,
      "settingSendCookies": true,
      "settingDisableRenderRequestBody": false,
      "settingEncodeUrl": true,
      "settingRebuildPath": true,
      "settingFollowRedirects": "global"
    },
    {
      "_id": "fld_123",
      "_type": "request_group",
      "parentId": "wrk_123",
      "modified": 1640995200000,
      "created": 1640995200000,
      "name": "Folder Name",
      "description": "Folder Description",
      "environment": {},
      "environmentPropertyOrder": null,
      "metaSortKey": -1640995200000
    },
    {
      "_id": "wrk_123",
      "_type": "workspace",
      "parentId": null,
      "modified": 1640995200000,
      "created": 1640995200000,
      "name": "Workspace Name",
      "description": "Workspace Description",
      "scope": "collection"
    },
    {
      "_id": "env_123",
      "_type": "environment",
      "parentId": "wrk_123",
      "modified": 1640995200000,
      "created": 1640995200000,
      "name": "Base Environment",
      "data": {
        "base_url": "http://localhost:8080",
        "accessToken": "",
        "refreshToken": ""
      },
      "dataPropertyOrder": {
        "&": ["base_url", "accessToken", "refreshToken"]
      },
      "color": null,
      "isPrivate": false,
      "metaSortKey": 1640995200000
    }
  ]
}
```

#### Export Endpoint
```
GET /api/v1/api-docs/collections/{id}/export/insomnia
```

### 4. Swagger JSON

#### Format Structure
```json
{
  "swagger": "2.0",
  "info": {
    "title": "API Title",
    "description": "API Description",
    "version": "1.0.0",
    "contact": {
      "name": "API Support",
      "email": "support@example.com"
    }
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "schemes": ["http", "https"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "JWT Authorization header using the Bearer scheme"
    }
  },
  "security": [
    {
      "BearerAuth": []
    }
  ],
  "paths": {
    "/endpoint": {
      "get": {
        "summary": "Endpoint Summary",
        "description": "Endpoint Description",
        "tags": ["Tag Name"],
        "parameters": [
          {
            "name": "param",
            "in": "query",
            "type": "string",
            "required": false,
            "description": "Parameter description"
          }
        ],
        "responses": {
          "200": {
            "description": "Success response",
            "schema": {
              "type": "object",
              "properties": {
                "success": {
                  "type": "boolean"
                },
                "data": {
                  "type": "object"
                }
              }
            },
            "examples": {
              "application/json": {
                "success": true,
                "data": {}
              }
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      },
      "post": {
        "summary": "Create endpoint",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "properties": {
                "name": {
                  "type": "string",
                  "description": "Name field"
                }
              },
              "required": ["name"]
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Created successfully"
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "success": {
          "type": "boolean",
          "example": false
        },
        "message": {
          "type": "string",
          "example": "Error message"
        }
      }
    }
  }
}
```

#### Export Endpoint
```
GET /api/v1/api-docs/collections/{id}/export/swagger
```

### 5. Apidog Collection

#### Format Structure
```json
{
  "apidogVersion": "2.0.0",
  "info": {
    "name": "Collection Name",
    "description": "Collection Description",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "{{base_url}}",
      "description": "Development Server"
    }
  ],
  "folders": [
    {
      "id": "folder_1",
      "name": "Folder Name",
      "description": "Folder Description",
      "parentId": null,
      "sort": 1
    }
  ],
  "apis": [
    {
      "id": "api_1",
      "name": "Endpoint Name",
      "description": "Endpoint Description",
      "folderId": "folder_1",
      "method": "GET",
      "path": "/api/v1/endpoint",
      "parameters": {
        "query": [
          {
            "name": "param",
            "type": "string",
            "required": false,
            "description": "Parameter description",
            "example": "value"
          }
        ],
        "header": [
          {
            "name": "Authorization",
            "type": "string",
            "required": true,
            "description": "JWT Token",
            "example": "Bearer {{accessToken}}"
          }
        ]
      },
      "requestBody": {
        "type": "application/json",
        "jsonSchema": {
          "type": "object",
          "properties": {
            "key": {
              "type": "string"
            }
          }
        },
        "example": {
          "key": "value"
        }
      },
      "responses": [
        {
          "statusCode": 200,
          "description": "Success response",
          "contentType": "application/json",
          "jsonSchema": {
            "type": "object",
            "properties": {
              "success": {
                "type": "boolean"
              },
              "data": {
                "type": "object"
              }
            }
          },
          "example": {
            "success": true,
            "data": {}
          }
        }
      ],
      "sort": 1
    }
  ],
  "environments": [
    {
      "id": "env_1",
      "name": "Development",
      "variables": [
        {
          "key": "base_url",
          "value": "http://localhost:8080",
          "description": "Base URL"
        },
        {
          "key": "accessToken",
          "value": "",
          "description": "JWT Access Token"
        }
      ]
    }
  ]
}
```

#### Export Endpoint
```
GET /api/v1/api-docs/collections/{id}/export/apidog
```

## Export Configuration Options

### Common Parameters
- `environment_id`: ID environment yang akan digunakan untuk variabel
- `include_tests`: Include test scripts dalam export (default: true)
- `include_examples`: Include response examples (default: true)
- `format`: Format output (json/yaml untuk OpenAPI)

### Postman Specific
- `include_prerequest`: Include pre-request scripts (default: true)
- `collection_id`: UUID untuk collection (auto-generated jika kosong)

### OpenAPI Specific
- `spec_version`: OpenAPI version (3.0.0, 3.0.1, 3.1.0)
- `include_servers`: Include server definitions (default: true)
- `include_security`: Include security schemes (default: true)

### Example Usage
```
GET /api/v1/api-docs/collections/1/export/postman?environment_id=1&include_tests=true
GET /api/v1/api-docs/collections/1/export/openapi?format=yaml&spec_version=3.0.1
GET /api/v1/api-docs/collections/1/export/insomnia?environment_id=2
```

## File Naming Convention

### Postman
`{collection_name}.postman_collection.json`

### OpenAPI
`{collection_name}.openapi.{json|yaml}`

### Insomnia
`{collection_name}.insomnia_collection.json`

### Swagger
`{collection_name}.swagger.json`

### Apidog
`{collection_name}.apidog_collection.json`

## Error Handling

### Common Errors
- `404`: Collection not found
- `403`: Insufficient permissions
- `422`: Invalid export parameters
- `500`: Export generation failed

### Error Response Format
```json
{
  "success": false,
  "message": "Export generation failed",
  "error": {
    "code": "EXPORT_ERROR",
    "details": "Specific error details"
  }
}
```

## Performance Considerations

### Caching Strategy
- Cache exported files untuk collection yang tidak berubah
- Cache key: `export:{collection_id}:{format}:{environment_id}:{hash}`
- TTL: 1 hour atau sampai collection diupdate

### Large Collections
- Implement streaming untuk collection besar (>1000 endpoints)
- Background job untuk export yang memakan waktu lama
- Progress tracking untuk long-running exports

### Rate Limiting
- Limit export requests: 10 per minute per user
- Implement queue untuk concurrent exports
- Priority queue untuk premium users