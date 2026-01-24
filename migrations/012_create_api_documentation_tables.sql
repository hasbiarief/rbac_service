-- Migration 012: Create API Documentation Core Tables
-- Core tables for API documentation system

CREATE TABLE IF NOT EXISTS api_collections (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) DEFAULT '1.0.0',
    base_url VARCHAR(500),
    schema_version VARCHAR(50) DEFAULT 'v2.1.0',
    created_by BIGINT NOT NULL REFERENCES users(id),
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_folders (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES api_folders(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_endpoints (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    folder_id BIGINT REFERENCES api_folders(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    method VARCHAR(10) NOT NULL CHECK (method IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS')),
    url VARCHAR(1000) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_collections_company_id ON api_collections(company_id);
CREATE INDEX IF NOT EXISTS idx_api_collections_created_by ON api_collections(created_by);
CREATE INDEX IF NOT EXISTS idx_api_collections_is_active ON api_collections(is_active);

CREATE INDEX IF NOT EXISTS idx_api_folders_collection_id ON api_folders(collection_id);
CREATE INDEX IF NOT EXISTS idx_api_folders_parent_id ON api_folders(parent_id);
CREATE INDEX IF NOT EXISTS idx_api_folders_sort_order ON api_folders(sort_order);

CREATE INDEX IF NOT EXISTS idx_api_endpoints_collection_id ON api_endpoints(collection_id);
CREATE INDEX IF NOT EXISTS idx_api_endpoints_folder_id ON api_endpoints(folder_id);
CREATE INDEX IF NOT EXISTS idx_api_endpoints_method ON api_endpoints(method);
CREATE INDEX IF NOT EXISTS idx_api_endpoints_sort_order ON api_endpoints(sort_order);

-- Create triggers for updated_at
CREATE TRIGGER update_api_collections_updated_at 
    BEFORE UPDATE ON api_collections 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_folders_updated_at 
    BEFORE UPDATE ON api_folders 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_endpoints_updated_at 
    BEFORE UPDATE ON api_endpoints 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE api_collections IS 'API documentation collections - main container for API documentation';
COMMENT ON TABLE api_folders IS 'Hierarchical folder structure for organizing API endpoints';
COMMENT ON TABLE api_endpoints IS 'Individual API endpoints with method, URL and basic information';

COMMENT ON COLUMN api_collections.schema_version IS 'Version of the collection schema (e.g., v2.1.0 for Postman)';
COMMENT ON COLUMN api_collections.base_url IS 'Default base URL for the API collection';
COMMENT ON COLUMN api_folders.sort_order IS 'Order of folders within parent folder or collection';
COMMENT ON COLUMN api_endpoints.sort_order IS 'Order of endpoints within folder';
COMMENT ON COLUMN api_endpoints.method IS 'HTTP method: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS';