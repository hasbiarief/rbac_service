-- Add API Documentation modules to the modules table
-- These modules are required for the API Documentation System

-- Main API Documentation module (139)
INSERT INTO public.modules (id, category, name, description, url, icon, parent_id, is_active, created_at, updated_at) 
VALUES (139, 'System & Security', 'API Documentation', 'Main API Documentation System module', '/api-docs', 'fas fa-book', NULL, true, NOW(), NOW());

-- Collections Management sub-module (140)
INSERT INTO public.modules (id, category, name, description, url, icon, parent_id, is_active, created_at, updated_at) 
VALUES (140, 'System & Security', 'Collections Management', 'Manage API documentation collections', '/api-docs/collections', 'fas fa-folder', 139, true, NOW(), NOW());

-- Endpoints Management sub-module (141)
INSERT INTO public.modules (id, category, name, description, url, icon, parent_id, is_active, created_at, updated_at) 
VALUES (141, 'System & Security', 'Endpoints Management', 'Manage API endpoints documentation', '/api-docs/endpoints', 'fas fa-plug', 139, true, NOW(), NOW());

-- Environments Management sub-module (142)
INSERT INTO public.modules (id, category, name, description, url, icon, parent_id, is_active, created_at, updated_at) 
VALUES (142, 'System & Security', 'Environments Management', 'Manage API testing environments', '/api-docs/environments', 'fas fa-server', 139, true, NOW(), NOW());

-- Export Documentation sub-module (143)
INSERT INTO public.modules (id, category, name, description, url, icon, parent_id, is_active, created_at, updated_at) 
VALUES (143, 'System & Security', 'Export Documentation', 'Export API documentation in various formats', '/api-docs/export', 'fas fa-download', 139, true, NOW(), NOW());