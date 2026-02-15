-- Drop API Documentation tables (replaced by Swagger)
-- Migration: 019_drop_api_documentation_tables
-- Description: Remove API Documentation System tables as we've migrated to Swagger/OpenAPI

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS api_endpoint_tags CASCADE;
DROP TABLE IF EXISTS api_tags CASCADE;
DROP TABLE IF EXISTS api_tests CASCADE;
DROP TABLE IF EXISTS api_environment_variables CASCADE;
DROP TABLE IF EXISTS api_environments CASCADE;
DROP TABLE IF EXISTS api_responses CASCADE;
DROP TABLE IF EXISTS api_request_bodies CASCADE;
DROP TABLE IF EXISTS api_parameters CASCADE;
DROP TABLE IF EXISTS api_headers CASCADE;
DROP TABLE IF EXISTS api_endpoints CASCADE;
DROP TABLE IF EXISTS api_folders CASCADE;
DROP TABLE IF EXISTS api_collections CASCADE;

-- Remove API Documentation modules from modules table
DELETE FROM modules WHERE id IN (139, 140, 141, 142, 143);

-- Remove API Documentation permissions
DELETE FROM role_modules WHERE module_id IN (139, 140, 141, 142, 143);
DELETE FROM unit_role_modules WHERE module_id IN (139, 140, 141, 142, 143);

-- Note: We keep the Swagger documentation system which is now the primary API documentation method
