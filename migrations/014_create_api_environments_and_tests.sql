-- Migration 014: Create API Documentation Environment Tables
-- Environment and variable management

CREATE TABLE IF NOT EXISTS api_environments (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT NOT NULL REFERENCES api_collections(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_id, name)
);

CREATE TABLE IF NOT EXISTS api_environment_variables (
    id BIGSERIAL PRIMARY KEY,
    environment_id BIGINT NOT NULL REFERENCES api_environments(id) ON DELETE CASCADE,
    key_name VARCHAR(255) NOT NULL,
    value TEXT,
    description TEXT,
    is_secret BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(environment_id, key_name)
);

CREATE TABLE IF NOT EXISTS api_tests (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id BIGINT NOT NULL REFERENCES api_endpoints(id) ON DELETE CASCADE,
    test_type VARCHAR(20) NOT NULL CHECK (test_type IN ('pre_request', 'test')),
    script_content TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_api_environments_collection_id ON api_environments(collection_id);
CREATE INDEX IF NOT EXISTS idx_api_environments_is_default ON api_environments(is_default);
CREATE INDEX IF NOT EXISTS idx_api_environments_name ON api_environments(name);

CREATE INDEX IF NOT EXISTS idx_api_env_vars_environment_id ON api_environment_variables(environment_id);
CREATE INDEX IF NOT EXISTS idx_api_env_vars_key_name ON api_environment_variables(key_name);
CREATE INDEX IF NOT EXISTS idx_api_env_vars_is_secret ON api_environment_variables(is_secret);

CREATE INDEX IF NOT EXISTS idx_api_tests_endpoint_id ON api_tests(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_api_tests_type ON api_tests(test_type);

-- Create triggers for updated_at
CREATE TRIGGER update_api_environments_updated_at 
    BEFORE UPDATE ON api_environments 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_api_tests_updated_at 
    BEFORE UPDATE ON api_tests 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE api_environments IS 'Environment configurations for API collections (dev, staging, prod)';
COMMENT ON TABLE api_environment_variables IS 'Variables for each environment (base_url, tokens, etc.)';
COMMENT ON TABLE api_tests IS 'Test scripts for API endpoints (pre-request and test scripts)';

COMMENT ON COLUMN api_environments.is_default IS 'Whether this is the default environment for the collection';
COMMENT ON COLUMN api_environment_variables.is_secret IS 'Whether this variable contains sensitive data (passwords, tokens)';
COMMENT ON COLUMN api_tests.test_type IS 'Type of test script: pre_request or test';
COMMENT ON COLUMN api_tests.script_content IS 'JavaScript code for test execution';