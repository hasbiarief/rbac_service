-- Add API Documentation module to the system
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) 
VALUES (139, 'System & Security', 'API Documentation', '/api-docs', 'FileText', 'Manage API documentation and collections', 131, 'pro', true, NOW(), NOW());

-- Add sub-modules for API Documentation
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) 
VALUES 
(140, 'System & Security', 'Collections Management', '/api-docs/collections', 'Folder', 'Manage API documentation collections', 139, 'pro', true, NOW(), NOW()),
(141, 'System & Security', 'Endpoints Management', '/api-docs/endpoints', 'Globe', 'Manage API endpoints documentation', 139, 'pro', true, NOW(), NOW()),
(142, 'System & Security', 'Environments Management', '/api-docs/environments', 'Settings', 'Manage API environments and variables', 139, 'pro', true, NOW(), NOW()),
(143, 'System & Security', 'Export Documentation', '/api-docs/export', 'Download', 'Export API documentation to various formats', 139, 'pro', true, NOW(), NOW());

-- Update sequence
SELECT pg_catalog.setval('public.modules_id_seq', 143, true);