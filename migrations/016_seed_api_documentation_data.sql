-- Migration 016: Seed Initial API Documentation Data
-- Initial data and additional indexes for API documentation system

-- Create additional composite indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_api_collections_company_active ON api_collections(company_id, is_active);
CREATE INDEX IF NOT EXISTS idx_api_folders_collection_parent ON api_folders(collection_id, parent_id);
CREATE INDEX IF NOT EXISTS idx_api_endpoints_collection_folder ON api_endpoints(collection_id, folder_id);
CREATE INDEX IF NOT EXISTS idx_api_endpoints_method_active ON api_endpoints(method, is_active);
CREATE INDEX IF NOT EXISTS idx_api_headers_endpoint_type ON api_headers(endpoint_id, header_type);
CREATE INDEX IF NOT EXISTS idx_api_parameters_endpoint_type ON api_parameters(endpoint_id, type);
CREATE INDEX IF NOT EXISTS idx_api_responses_endpoint_default ON api_responses(endpoint_id, is_default);

-- Create a sample API collection for demonstration
INSERT INTO api_collections (name, description, version, base_url, created_by, company_id, is_active) 
SELECT 
    'Huminor RBAC API Documentation',
    'Complete API documentation for ERP with Unit-Based Role Access Control (RBAC)',
    '1.0.0',
    '{{base_url}}',
    u.id,
    c.id,
    true
FROM users u, companies c 
WHERE u.email = 'hasbi@company.com' 
AND c.code = 'CINTA'
ON CONFLICT DO NOTHING;

-- Get the collection ID for further seeding
DO $$
DECLARE
    collection_id BIGINT;
    env_id BIGINT;
    auth_folder_id BIGINT;
    user_folder_id BIGINT;
    endpoint_id BIGINT;
BEGIN
    -- Get collection ID
    SELECT id INTO collection_id 
    FROM api_collections 
    WHERE name = 'Huminor RBAC API Documentation' 
    LIMIT 1;
    
    IF collection_id IS NOT NULL THEN
        -- Create default environment
        INSERT INTO api_environments (collection_id, name, description, is_default) 
        VALUES (collection_id, 'Development', 'Development environment configuration', TRUE)
        ON CONFLICT (collection_id, name) DO NOTHING
        RETURNING id INTO env_id;
        
        -- If environment already exists, get its ID
        IF env_id IS NULL THEN
            SELECT id INTO env_id 
            FROM api_environments 
            WHERE collection_id = collection_id AND name = 'Development';
        END IF;
        
        -- Add environment variables
        INSERT INTO api_environment_variables (environment_id, key_name, value, description, is_secret) VALUES
        (env_id, 'base_url', 'http://localhost:8081', 'Base URL for API requests', FALSE),
        (env_id, 'accessToken', '', 'JWT Access Token for authentication', TRUE),
        (env_id, 'refreshToken', '', 'JWT Refresh Token', TRUE),
        (env_id, 'userIdentity', 'admin', 'User identity for login', FALSE),
        (env_id, 'userPassword', 'password123', 'User password for login', TRUE),
        (env_id, 'userId', '', 'Current authenticated user ID', FALSE),
        (env_id, 'companyId', '1', 'Company ID for multi-tenant operations', FALSE)
        ON CONFLICT (environment_id, key_name) DO NOTHING;
        
        -- Create folder structure
        INSERT INTO api_folders (collection_id, parent_id, name, description, sort_order) VALUES
        (collection_id, NULL, 'Authentication', 'User authentication and authorization endpoints', 1),
        (collection_id, NULL, 'User Management', 'User CRUD and profile management', 2),
        (collection_id, NULL, 'Role Management', 'Role-based access control management', 3),
        (collection_id, NULL, 'Company Management', 'Company and organization management', 4),
        (collection_id, NULL, 'Module System', 'System modules and permissions', 5),
        (collection_id, NULL, 'API Documentation', 'API documentation management endpoints', 6)
        ON CONFLICT DO NOTHING;
        
        -- Get folder IDs
        SELECT id INTO auth_folder_id FROM api_folders WHERE collection_id = collection_id AND name = 'Authentication';
        SELECT id INTO user_folder_id FROM api_folders WHERE collection_id = collection_id AND name = 'User Management';
        
        -- Create sample endpoints
        IF auth_folder_id IS NOT NULL THEN
            INSERT INTO api_endpoints (collection_id, folder_id, name, description, method, url, sort_order) VALUES
            (collection_id, auth_folder_id, 'Login with User Identity', 'Authenticate user with user_identity and password', 'POST', '/api/v1/auth/login', 1),
            (collection_id, auth_folder_id, 'Login with Email', 'Authenticate user with email and password', 'POST', '/api/v1/auth/login-email', 2),
            (collection_id, auth_folder_id, 'Refresh Token', 'Refresh JWT access token using refresh token', 'POST', '/api/v1/auth/refresh', 3),
            (collection_id, auth_folder_id, 'Logout', 'Logout user and invalidate tokens', 'POST', '/api/v1/auth/logout', 4)
            ON CONFLICT DO NOTHING;
        END IF;
        
        IF user_folder_id IS NOT NULL THEN
            INSERT INTO api_endpoints (collection_id, folder_id, name, description, method, url, sort_order) VALUES
            (collection_id, user_folder_id, 'Get All Users', 'Retrieve paginated list of users', 'GET', '/api/v1/users', 1),
            (collection_id, user_folder_id, 'Get User by ID', 'Retrieve user details by ID', 'GET', '/api/v1/users/{id}', 2),
            (collection_id, user_folder_id, 'Create User', 'Create new user account', 'POST', '/api/v1/users', 3),
            (collection_id, user_folder_id, 'Update User', 'Update user information', 'PUT', '/api/v1/users/{id}', 4)
            ON CONFLICT DO NOTHING;
        END IF;
        
        -- Add sample headers for login endpoint
        SELECT id INTO endpoint_id FROM api_endpoints WHERE collection_id = collection_id AND name = 'Login with User Identity';
        IF endpoint_id IS NOT NULL THEN
            INSERT INTO api_headers (endpoint_id, key_name, value, description, is_required, header_type) VALUES
            (endpoint_id, 'Content-Type', 'application/json', 'Request content type', TRUE, 'request'),
            (endpoint_id, 'Accept', 'application/json', 'Response content type', FALSE, 'request')
            ON CONFLICT DO NOTHING;
            
            -- Add sample request body
            INSERT INTO api_request_bodies (endpoint_id, content_type, body_content, description) VALUES
            (endpoint_id, 'application/json', '{"user_identity": "{{userIdentity}}", "password": "{{userPassword}}"}', 'Login credentials with user identity')
            ON CONFLICT DO NOTHING;
            
            -- Add sample responses
            INSERT INTO api_responses (endpoint_id, status_code, status_text, content_type, response_body, description, is_default) VALUES
            (endpoint_id, 200, 'OK', 'application/json', '{"success": true, "data": {"access_token": "jwt_token", "refresh_token": "refresh_token", "user": {"id": 1, "name": "User Name"}}}', 'Successful login response', TRUE),
            (endpoint_id, 401, 'Unauthorized', 'application/json', '{"success": false, "message": "Invalid credentials"}', 'Invalid login credentials', FALSE),
            (endpoint_id, 422, 'Unprocessable Entity', 'application/json', '{"success": false, "message": "Validation failed", "errors": {"user_identity": ["required"]}}', 'Validation error response', FALSE)
            ON CONFLICT DO NOTHING;
        END IF;
        
    END IF;
END $$;

-- Create view for easy collection statistics
CREATE OR REPLACE VIEW api_collection_stats AS
SELECT 
    c.id,
    c.name,
    c.company_id,
    COUNT(DISTINCT f.id) as total_folders,
    COUNT(DISTINCT e.id) as total_endpoints,
    COUNT(DISTINCT env.id) as total_environments,
    COUNT(DISTINCT CASE WHEN e.method = 'GET' THEN e.id END) as get_endpoints,
    COUNT(DISTINCT CASE WHEN e.method = 'POST' THEN e.id END) as post_endpoints,
    COUNT(DISTINCT CASE WHEN e.method = 'PUT' THEN e.id END) as put_endpoints,
    COUNT(DISTINCT CASE WHEN e.method = 'DELETE' THEN e.id END) as delete_endpoints,
    c.created_at,
    c.updated_at
FROM api_collections c
LEFT JOIN api_folders f ON c.id = f.collection_id
LEFT JOIN api_endpoints e ON c.id = e.collection_id AND e.is_active = TRUE
LEFT JOIN api_environments env ON c.id = env.collection_id
WHERE c.is_active = TRUE
GROUP BY c.id, c.name, c.company_id, c.created_at, c.updated_at;

COMMENT ON VIEW api_collection_stats IS 'Statistics view for API collections showing counts of folders, endpoints, and environments';

-- Create function to get endpoint details with all related data
CREATE OR REPLACE FUNCTION get_endpoint_details(endpoint_id_param BIGINT)
RETURNS JSON AS $$
DECLARE
    result JSON;
BEGIN
    SELECT json_build_object(
        'endpoint', row_to_json(e),
        'headers', COALESCE(headers.headers_array, '[]'::json),
        'parameters', COALESCE(params.params_array, '[]'::json),
        'request_body', COALESCE(row_to_json(rb), 'null'::json),
        'responses', COALESCE(responses.responses_array, '[]'::json),
        'tags', COALESCE(tags.tags_array, '[]'::json)
    ) INTO result
    FROM api_endpoints e
    LEFT JOIN (
        SELECT endpoint_id, json_agg(row_to_json(h)) as headers_array
        FROM api_headers h
        WHERE h.endpoint_id = endpoint_id_param
        GROUP BY endpoint_id
    ) headers ON e.id = headers.endpoint_id
    LEFT JOIN (
        SELECT endpoint_id, json_agg(row_to_json(p)) as params_array
        FROM api_parameters p
        WHERE p.endpoint_id = endpoint_id_param
        GROUP BY endpoint_id
    ) params ON e.id = params.endpoint_id
    LEFT JOIN api_request_bodies rb ON e.id = rb.endpoint_id
    LEFT JOIN (
        SELECT endpoint_id, json_agg(row_to_json(r)) as responses_array
        FROM api_responses r
        WHERE r.endpoint_id = endpoint_id_param
        GROUP BY endpoint_id
    ) responses ON e.id = responses.endpoint_id
    LEFT JOIN (
        SELECT et.endpoint_id, json_agg(row_to_json(t)) as tags_array
        FROM api_endpoint_tags et
        JOIN api_tags t ON et.tag_id = t.id
        WHERE et.endpoint_id = endpoint_id_param
        GROUP BY et.endpoint_id
    ) tags ON e.id = tags.endpoint_id
    WHERE e.id = endpoint_id_param;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_endpoint_details(BIGINT) IS 'Get complete endpoint details with all related data as JSON';