-- Migration 013: Create API Documentation Detail Tables
-- Supporting tables for endpoint details

CREATE TABLE IF NOT EXISTS api_headers (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    header_type VARCHAR(10) DEFAULT 'request' CHECK (header_type IN ('request', 'response')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_parameters (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('query', 'path', 'form', 'file')),
    data_type VARCHAR(50) DEFAULT 'string',
    description TEXT,
    default_value TEXT,
    example_value TEXT,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_request_bodies (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    content_type VARCHAR(100) DEFAULT 'application/json',
    body_content TEXT,
    description TEXT,
    schema_definition JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_responses (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    status_code INTEGER NOT NULL,
    status_text VARCHAR(100),
    content_type VARCHAR(100) DEFAULT 'application/json',
    response_body TEXT,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_headers_endpoint_id ON api_headers(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_headers_type ON api_headers(header_type);
CREATE INDEX IF NOT EXISTS idx_api_headers_key_name ON api_headers(key_name);

CREATE INDEX IF NOT EXISTS idx_api_parameters_endpoint_id ON api_parameters(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_parameters_type ON api_parameters(type);
CREATE INDEX IF NOT EXISTS idx_api_parameters_name ON api_parameters(name);

CREATE INDEX IF NOT EXISTS idx_api_request_bodies_endpoint_id ON api_request_bodies(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_request_bodies_content_type ON api_request_bodies(content_type);

CREATE INDEX IF NOT EXISTS idx_api_responses_endpoint_id ON api_responses(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_responses_status_code ON api_responses(status_code);
CREATE INDEX IF NOT EXISTS idx_api_responses_is_default ON api_responses(is_default);

-- Create triggers for updated_at
CREATE TRIGGER update_api_request_bodies_updated_at 
    BEFORE UPDATE ON api_request_bodies 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE api_headers IS 'HTTP headers for API endpoints (both request and response headers)';
COMMENT ON TABLE api_parameters IS 'Parameters for API endpoints (query, path, form, file parameters)';
COMMENT ON TABLE api_request_bodies IS 'Request body documentation with content type and schema definition';
COMMENT ON TABLE api_responses IS 'Response examples with status codes and response bodies';

COMMENT ON COLUMN api_headers.header_type IS 'Type of header: request or response';
COMMENT ON COLUMN api_headers.key_name IS 'Header name (e.g., Authorization, Content-Type)';
COMMENT ON COLUMN api_parameters.type IS 'Parameter type: query, path, form, or file';
COMMENT ON COLUMN api_parameters.data_type IS 'Data type of parameter (string, integer, boolean, etc.)';
COMMENT ON COLUMN api_request_bodies.schema_definition IS 'JSON schema definition for request body validation';
COMMENT ON COLUMN api_responses.status_code IS 'HTTP status code (200, 404, 500, etc.)';
COMMENT ON COLUMN api_responses.is_default IS 'Whether this is the default response example for the endpoint';