-- Migration 008: Add can_approve column to role_modules table
-- Add approval permission to role-module relationship

-- Add can_approve column to role_modules table
ALTER TABLE role_modules ADD COLUMN can_approve BOOLEAN DEFAULT false;

-- Add comment for the new column
COMMENT ON COLUMN role_modules.can_approve IS 'Permission to approve requests/transactions in this module';

-- Update existing records to set can_approve based on role hierarchy
-- Super Admin and Admin roles get approve permission for all their modules
UPDATE role_modules 
SET can_approve = true 
WHERE role_id IN (
    SELECT id FROM roles 
    WHERE name IN ('SUPER_ADMIN', 'ADMIN', 'HR_MANAGER', 'MANAGER')
);

-- HR Staff gets approve permission only for basic HR modules
UPDATE role_modules 
SET can_approve = true 
WHERE role_id IN (
    SELECT id FROM roles WHERE name = 'HR_STAFF'
) AND module_id IN (
    SELECT id FROM modules 
    WHERE category IN ('Employee Self Service', 'Leave Management') 
    AND subscription_tier = 'basic'
);

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_role_modules_can_approve ON role_modules(can_approve);