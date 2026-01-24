-- Migration 015: Create API Documentation Tagging System
-- Tag system for endpoint categorization

CREATE TABLE IF NOT EXISTS api_tags (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7) DEFAULT '#007bff',
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_endpoint_tags (
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES api_tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (endpoint_id, tag_id)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_tags_name ON api_tags(name);
CREATE INDEX IF NOT EXISTS idx_api_endpoint_tags_endpoint_id ON api_endpoint_tags(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_endpoint_tags_tag_id ON api_endpoint_tags(tag_id);

-- Insert default tags
INSERT INTO api_tags (name, color, description) VALUES
('Authentication', '#dc3545', 'Authentication related endpoints'),
('CRUD', '#28a745', 'Create, Read, Update, Delete operations'),
('Admin', '#ffc107', 'Administrative endpoints'),
('Public', '#17a2b8', 'Public accessible endpoints'),
('Deprecated', '#6c757d', 'Deprecated endpoints'),
('Internal', '#6f42c1', 'Internal system endpoints'),
('External', '#fd7e14', 'External API integrations'),
('Reporting', '#20c997', 'Reporting and analytics endpoints')
ON CONFLICT (name) DO NOTHING;

-- Add comments for documentation
COMMENT ON TABLE api_tags IS 'Tags for categorizing and organizing API endpoints';
COMMENT ON TABLE api_endpoint_tags IS 'Many-to-many relationship between endpoints and tags';

COMMENT ON COLUMN api_tags.color IS 'Hex color code for tag display (e.g., #007bff)';
COMMENT ON COLUMN api_tags.name IS 'Unique tag name for endpoint categorization';