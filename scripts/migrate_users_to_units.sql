-- Migration Script: Assign Users to Units
-- This script migrates existing users to appropriate units based on their roles

-- Create backup table first
CREATE TABLE IF NOT EXISTS user_roles_backup AS 
SELECT * FROM user_roles WHERE 1=0;

-- Backup current user_roles before migration
INSERT INTO user_roles_backup 
SELECT * FROM user_roles WHERE unit_id IS NULL;

-- Display backup confirmation
SELECT 
    'BACKUP CREATED' as status,
    COUNT(*) as backed_up_records
FROM user_roles_backup;

-- Migration Plan Display
SELECT 
    'MIGRATION PLAN' as section,
    ur.id as user_role_id,
    u.name as user_name,
    r.name as role_name,
    b.name as branch_name,
    CASE 
        WHEN r.name IN ('SUPER_ADMIN', 'ADMIN') THEN 'KEEP AT BRANCH LEVEL'
        WHEN r.name = 'HR_ADMIN' THEN 'MIGRATE TO HR ADMIN UNIT'
        WHEN r.name = 'HR_MANAGER' THEN 'MIGRATE TO HR DEPARTMENT'
        WHEN r.name = 'RECRUITER' THEN 'MIGRATE TO RECRUITMENT UNIT'
        WHEN r.name = 'PAYROLL_OFFICER' THEN 'MIGRATE TO PAYROLL UNIT'
        WHEN r.name = 'IT_ADMIN' THEN 'MIGRATE TO IT DEPARTMENT'
        WHEN r.name = 'LINE_MANAGER' AND ur.branch_id = 1 THEN 'MIGRATE TO OPERATIONS UNIT'
        WHEN r.name = 'LINE_MANAGER' AND ur.branch_id != 1 THEN 'MIGRATE TO SALES UNIT'
        WHEN r.name = 'EMPLOYEE' AND ur.branch_id != 1 THEN 'MIGRATE TO ADMIN UNIT'
        WHEN r.name = 'EMPLOYEE' AND ur.branch_id = 1 THEN 'MIGRATE TO HR ADMIN UNIT'
        ELSE 'KEEP AT BRANCH LEVEL'
    END as migration_action
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN branches b ON ur.branch_id = b.id
WHERE ur.unit_id IS NULL
ORDER BY ur.id;

-- 1. Migrate EMPLOYEE users to appropriate Admin units
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'ADM' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'EMPLOYEE')
AND branch_id != 1; -- Non-headquarters branches

-- For employees at headquarters, assign to HR Admin unit if available
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'HR-ADM' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'EMPLOYEE')
AND branch_id = 1; -- Headquarters

-- 2. Migrate HR_ADMIN users to HR Admin units
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'HR-ADM' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'HR_ADMIN');

-- 3. Migrate HR_MANAGER users to HR Department
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'HR' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'HR_MANAGER');

-- 4. Migrate RECRUITER users to Recruitment units
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'HR-REC' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'RECRUITER');

-- 5. Migrate PAYROLL_OFFICER users to Payroll units
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'HR-PAY' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'PAYROLL_OFFICER');

-- 6. Migrate IT_ADMIN users to IT Department
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'IT' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'IT_ADMIN');

-- 7. Migrate LINE_MANAGER users to appropriate units
-- For headquarters: Operations unit
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'OPS' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'LINE_MANAGER')
AND branch_id = 1; -- Headquarters

-- For branches: Sales unit
UPDATE user_roles 
SET unit_id = (
    SELECT un.id 
    FROM units un 
    WHERE un.code = 'SALES' 
    AND un.branch_id = user_roles.branch_id 
    LIMIT 1
)
WHERE unit_id IS NULL 
AND role_id = (SELECT id FROM roles WHERE name = 'LINE_MANAGER')
AND branch_id != 1; -- Branch offices

-- 8. Keep SUPER_ADMIN and other high-level roles at branch level (no unit assignment)
-- These roles intentionally remain with unit_id = NULL for branch-wide access

-- Migration Results Summary
SELECT 
    'MIGRATION RESULTS' as section,
    r.name as role_name,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) as migrated_to_units,
    COUNT(CASE WHEN ur.unit_id IS NULL THEN 1 END) as kept_at_branch_level,
    COUNT(*) as total_users
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
GROUP BY r.name
ORDER BY r.name;

-- Detailed migration results
SELECT 
    'DETAILED RESULTS' as section,
    u.name as user_name,
    r.name as role_name,
    b.name as branch_name,
    un.name as unit_name,
    un.code as unit_code,
    CASE 
        WHEN ur.unit_id IS NOT NULL THEN 'MIGRATED TO UNIT'
        ELSE 'KEPT AT BRANCH LEVEL'
    END as status
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN branches b ON ur.branch_id = b.id
LEFT JOIN units un ON ur.unit_id = un.id
ORDER BY u.name, r.name;

-- Validation: Check for any orphaned assignments
SELECT 
    'VALIDATION CHECK' as section,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.id IS NULL THEN 1 END) as orphaned_unit_assignments,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.branch_id != ur.branch_id THEN 1 END) as mismatched_branch_units
FROM user_roles ur
LEFT JOIN units un ON ur.unit_id = un.id;

-- Final summary
SELECT 
    'FINAL SUMMARY' as section,
    COUNT(CASE WHEN unit_id IS NOT NULL THEN 1 END) as users_with_unit_assignment,
    COUNT(CASE WHEN unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(*) as total_user_roles
FROM user_roles;