-- Audit Script: Current User Role Assignments
-- This script analyzes existing user roles and suggests unit assignments

-- 1. Overview of current user roles without unit assignment
SELECT 
    'USERS WITHOUT UNIT ASSIGNMENT' as report_section,
    COUNT(*) as total_count
FROM user_roles 
WHERE unit_id IS NULL;

-- 2. Detailed breakdown of users without unit assignment
SELECT 
    u.id as user_id,
    u.name as user_name,
    u.email,
    r.name as role_name,
    c.name as company_name,
    b.name as branch_name,
    ur.created_at as role_assigned_date
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
JOIN companies c ON ur.company_id = c.id
LEFT JOIN branches b ON ur.branch_id = b.id
WHERE ur.unit_id IS NULL
ORDER BY c.name, b.name, r.name, u.name;

-- 3. Available units per branch for assignment
SELECT 
    b.name as branch_name,
    COUNT(un.id) as total_units,
    STRING_AGG(un.name || ' (' || un.code || ')', ', ' ORDER BY un.name) as available_units
FROM branches b
LEFT JOIN units un ON b.id = un.branch_id AND un.is_active = true
GROUP BY b.id, b.name
ORDER BY b.name;

-- 4. Current unit-role mappings
SELECT 
    un.name as unit_name,
    un.code as unit_code,
    r.name as role_name,
    COUNT(urm.id) as custom_permissions
FROM unit_roles ur
JOIN units un ON ur.unit_id = un.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN unit_role_modules urm ON ur.id = urm.unit_role_id
GROUP BY un.name, un.code, r.name
ORDER BY un.name, r.name;

-- 5. Suggested unit assignments based on role and branch
WITH role_suggestions AS (
    SELECT 
        ur.id as user_role_id,
        ur.user_id,
        u.name as user_name,
        ur.role_id,
        r.name as role_name,
        ur.branch_id,
        b.name as branch_name,
        CASE 
            -- HR roles should go to HR units
            WHEN r.name = 'HR_ADMIN' THEN 
                (SELECT id FROM units WHERE code = 'HR-ADM' AND branch_id = ur.branch_id LIMIT 1)
            WHEN r.name = 'HR_MANAGER' THEN 
                (SELECT id FROM units WHERE code = 'HR' AND branch_id = ur.branch_id LIMIT 1)
            WHEN r.name = 'RECRUITER' THEN 
                (SELECT id FROM units WHERE code = 'HR-REC' AND branch_id = ur.branch_id LIMIT 1)
            WHEN r.name = 'PAYROLL_OFFICER' THEN 
                (SELECT id FROM units WHERE code = 'HR-PAY' AND branch_id = ur.branch_id LIMIT 1)
            
            -- IT roles should go to IT units
            WHEN r.name = 'IT_ADMIN' THEN 
                (SELECT id FROM units WHERE code = 'IT' AND branch_id = ur.branch_id LIMIT 1)
            
            -- Line managers go to appropriate department units
            WHEN r.name = 'LINE_MANAGER' AND ur.branch_id = 1 THEN -- Pusat
                (SELECT id FROM units WHERE code = 'OPS' AND branch_id = ur.branch_id LIMIT 1)
            WHEN r.name = 'LINE_MANAGER' AND ur.branch_id != 1 THEN -- Branches
                (SELECT id FROM units WHERE code = 'SALES' AND branch_id = ur.branch_id LIMIT 1)
            
            -- Employees can stay at branch level or go to admin units
            WHEN r.name = 'EMPLOYEE' THEN 
                (SELECT id FROM units WHERE code = 'ADM' AND branch_id = ur.branch_id LIMIT 1)
            
            -- Super admin and other roles stay at branch level
            ELSE NULL
        END as suggested_unit_id,
        CASE 
            WHEN r.name = 'HR_ADMIN' THEN 'HR Admin Unit'
            WHEN r.name = 'HR_MANAGER' THEN 'HR Department'
            WHEN r.name = 'RECRUITER' THEN 'Recruitment Unit'
            WHEN r.name = 'PAYROLL_OFFICER' THEN 'Payroll Unit'
            WHEN r.name = 'IT_ADMIN' THEN 'IT Department'
            WHEN r.name = 'LINE_MANAGER' AND ur.branch_id = 1 THEN 'Operations Unit'
            WHEN r.name = 'LINE_MANAGER' AND ur.branch_id != 1 THEN 'Sales Unit'
            WHEN r.name = 'EMPLOYEE' THEN 'Admin Unit'
            ELSE 'Keep at Branch Level'
        END as suggested_assignment
    FROM user_roles ur
    JOIN users u ON ur.user_id = u.id
    JOIN roles r ON ur.role_id = r.id
    LEFT JOIN branches b ON ur.branch_id = b.id
    WHERE ur.unit_id IS NULL
)
SELECT 
    rs.user_name,
    rs.role_name,
    rs.branch_name,
    rs.suggested_assignment,
    un.name as suggested_unit_name,
    un.code as suggested_unit_code,
    CASE 
        WHEN rs.suggested_unit_id IS NOT NULL THEN 'MIGRATE TO UNIT'
        ELSE 'KEEP AT BRANCH LEVEL'
    END as migration_action
FROM role_suggestions rs
LEFT JOIN units un ON rs.suggested_unit_id = un.id
ORDER BY rs.branch_name, rs.role_name, rs.user_name;

-- 6. Migration readiness check
SELECT 
    'MIGRATION READINESS CHECK' as report_section,
    COUNT(CASE WHEN ur.unit_id IS NULL THEN 1 END) as users_to_migrate,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) as users_already_migrated,
    COUNT(*) as total_user_roles
FROM user_roles ur;

-- 7. Potential conflicts or issues
SELECT 
    'POTENTIAL ISSUES' as report_section,
    u.name as user_name,
    COUNT(ur.id) as role_count,
    STRING_AGG(r.name, ', ') as roles
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE ur.unit_id IS NULL
GROUP BY u.id, u.name
HAVING COUNT(ur.id) > 1
ORDER BY role_count DESC;