-- Migration 018: Add application_id to roles table
-- This allows roles to be associated with specific applications

-- Add application_id column to roles table
ALTER TABLE roles ADD COLUMN application_id BIGINT REFERENCES applications(id) ON DELETE SET NULL;

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_roles_application_id ON roles(application_id);

-- Update existing roles to assign them to appropriate applications
-- RBAC roles (system management roles)
UPDATE roles SET application_id = (SELECT id FROM applications WHERE code = 'RBAC') 
WHERE id IN (13); -- CONSOLE ADMIN

-- HRIS roles (HR management roles)  
UPDATE roles SET application_id = (SELECT id FROM applications WHERE code = 'HRIS')
WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15); -- All other roles

-- Add comment for documentation
COMMENT ON COLUMN roles.application_id IS 'Reference to application that this role belongs to. NULL means role can be used across applications.';

-- Verify the updates
-- SELECT r.id, r.name, a.name as application_name, a.code as application_code 
-- FROM roles r 
-- LEFT JOIN applications a ON r.application_id = a.id 
-- ORDER BY r.id;