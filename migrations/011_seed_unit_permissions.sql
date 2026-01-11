-- Migration 011: Seed Unit Permissions
-- Create sample unit-specific permissions to demonstrate the new system

-- Create unit-specific role permissions that differ from default role permissions
-- This demonstrates how units can have customized access

-- HR Admin unit gets full HR access but limited to HR modules only
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    true,
    true,
    CASE 
        WHEN m.category IN ('Leave Management', 'Attendance & Time', 'Performance Management') THEN true
        ELSE false
    END
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN (
    'Core HR / Master Data',
    'Employee Self Service', 
    'Leave Management',
    'Attendance & Time',
    'Performance Management',
    'Training & Development',
    'Disciplinary & Relations'
)
WHERE u.code = 'HR-ADM' AND ur.role_id = 2 -- HR_ADMIN role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- Recruitment unit gets recruitment-focused permissions
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    true,
    false, -- No delete access for recruitment
    false  -- No approval access (requires HR Manager approval)
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN ('Recruitment', 'Core HR / Master Data')
WHERE u.code = 'HR-REC' AND ur.role_id = 9 -- RECRUITER role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- Payroll unit gets payroll-specific permissions
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    CASE 
        WHEN m.category = 'Payroll & Compensation' THEN true
        ELSE false -- Read-only for other modules
    END,
    false, -- No delete access
    CASE 
        WHEN m.category = 'Payroll & Compensation' THEN true
        ELSE false
    END
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN (
    'Payroll & Compensation',
    'Attendance & Time',
    'Leave Management',
    'Reporting & Analytics'
)
WHERE u.code = 'HR-PAY' AND ur.role_id = 4 -- PAYROLL_OFFICER role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- Finance units get finance-specific access
-- Accounting unit
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    true,
    false,
    true
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN (
    'Payroll & Compensation',
    'Reporting & Analytics',
    'Asset & Facility'
)
WHERE u.code = 'FIN-ACC' AND ur.role_id = 5 -- LINE_MANAGER role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- Sales units get limited access focused on employee management
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    CASE 
        WHEN m.category IN ('Employee Self Service', 'Performance Management') THEN true
        ELSE false
    END,
    false,
    CASE 
        WHEN m.category IN ('Leave Management', 'Attendance & Time', 'Performance Management') THEN true
        ELSE false
    END
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN (
    'Employee Self Service',
    'Attendance & Time',
    'Leave Management',
    'Performance Management',
    'Dashboard & Analytic'
)
WHERE u.code = 'SALES' AND ur.role_id = 5 -- LINE_MANAGER role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- IT Admin unit gets system access
INSERT INTO unit_role_modules (unit_role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 
    ur.id,
    m.id,
    true,
    true,
    true,
    false -- IT Admin can't approve business processes
FROM unit_roles ur
JOIN units u ON ur.unit_id = u.id
JOIN modules m ON m.category IN (
    'System & Security',
    'Reporting & Analytics'
)
WHERE u.code = 'IT' AND ur.role_id = 12 -- IT_ADMIN role
AND m.is_active = true
ON CONFLICT (unit_role_id, module_id) DO NOTHING;

-- Update some existing user_roles to include unit assignments
-- Assign HR users to specific HR units
UPDATE user_roles 
SET unit_id = (SELECT id FROM units WHERE code = 'HR-ADM' AND branch_id = 1)
WHERE role_id = 2 AND company_id = 1; -- HR_ADMIN

UPDATE user_roles 
SET unit_id = (SELECT id FROM units WHERE code = 'HR' AND branch_id = 1)
WHERE role_id = 3 AND company_id = 1; -- HR_MANAGER

-- Assign branch managers to sales units
UPDATE user_roles 
SET unit_id = (SELECT id FROM units WHERE code = 'SALES' AND branch_id = 2)
WHERE role_id = 5 AND branch_id = 2; -- LINE_MANAGER in Jakarta

UPDATE user_roles 
SET unit_id = (SELECT id FROM units WHERE code = 'SALES' AND branch_id = 3)
WHERE role_id = 5 AND branch_id = 3; -- LINE_MANAGER in Surabaya

-- Create some additional sample users with unit assignments
INSERT INTO users (name, email, user_identity, password_hash, is_active) VALUES
('Recruitment Specialist', 'recruiter@company.com', '100000007', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true),
('Payroll Officer', 'payroll@company.com', '100000008', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true),
('Sales Manager Jakarta', 'sales.jkt@company.com', '100000009', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true),
('IT Administrator', 'it.admin@company.com', '100000010', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true)
ON CONFLICT (email) DO NOTHING;

-- Assign these users to specific units
INSERT INTO user_roles (user_id, role_id, company_id, branch_id, unit_id) VALUES
-- Recruitment Specialist to Recruitment unit
((SELECT id FROM users WHERE email = 'recruiter@company.com'), 9, 1, 1, 
 (SELECT id FROM units WHERE code = 'HR-REC' AND branch_id = 1)),

-- Payroll Officer to Payroll unit
((SELECT id FROM users WHERE email = 'payroll@company.com'), 4, 1, 1, 
 (SELECT id FROM units WHERE code = 'HR-PAY' AND branch_id = 1)),

-- Sales Manager to Sales unit in Jakarta
((SELECT id FROM users WHERE email = 'sales.jkt@company.com'), 5, 1, 2, 
 (SELECT id FROM units WHERE code = 'SALES' AND branch_id = 2)),

-- IT Administrator to IT unit
((SELECT id FROM users WHERE email = 'it.admin@company.com'), 12, 1, 1, 
 (SELECT id FROM units WHERE code = 'IT' AND branch_id = 1))
ON CONFLICT (user_id, role_id, company_id, branch_id) DO NOTHING;

-- Add some comments for documentation
COMMENT ON TABLE unit_role_modules IS 'Unit-specific role permissions that can override default role permissions for granular access control';

-- Create a view for easy permission checking
CREATE OR REPLACE VIEW user_effective_permissions AS
SELECT DISTINCT
    ur.user_id,
    ur.company_id,
    ur.branch_id,
    ur.unit_id,
    m.id as module_id,
    m.name as module_name,
    m.category as module_category,
    m.url as module_url,
    COALESCE(urm.can_read, rm.can_read, false) as can_read,
    COALESCE(urm.can_write, rm.can_write, false) as can_write,
    COALESCE(urm.can_delete, rm.can_delete, false) as can_delete,
    COALESCE(urm.can_approve, rm.can_approve, false) as can_approve,
    CASE WHEN urm.id IS NOT NULL THEN true ELSE false END as is_customized
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
JOIN role_modules rm ON r.id = rm.role_id
JOIN modules m ON rm.module_id = m.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id AND m.id = urm.module_id
WHERE ur.user_id IS NOT NULL
AND r.is_active = true
AND m.is_active = true;

COMMENT ON VIEW user_effective_permissions IS 'Consolidated view of user permissions combining default role permissions with unit-specific overrides';

-- Create indexes for the new view
CREATE INDEX IF NOT EXISTS idx_user_effective_permissions_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_effective_permissions_module ON user_roles(user_id) INCLUDE (company_id, branch_id, unit_id);