-- Permission Resolution Testing Script
-- This script tests how permissions are resolved for users with unit assignments

-- 1. Test permission resolution for each user
SELECT 
    'USER PERMISSION RESOLUTION TEST' as test_section,
    u.name as user_name,
    r.name as role_name,
    un.name as unit_name,
    m.name as module_name,
    m.category as module_category,
    COALESCE(urm.can_read, rm.can_read, false) as effective_read,
    COALESCE(urm.can_write, rm.can_write, false) as effective_write,
    COALESCE(urm.can_delete, rm.can_delete, false) as effective_delete,
    COALESCE(urm.can_approve, rm.can_approve, false) as effective_approve,
    CASE 
        WHEN urm.id IS NOT NULL THEN 'UNIT_CUSTOM'
        WHEN rm.id IS NOT NULL THEN 'ROLE_DEFAULT'
        ELSE 'NO_ACCESS'
    END as permission_source
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id AND (urm.module_id IS NULL OR rm.module_id = urm.module_id)
LEFT JOIN modules m ON COALESCE(urm.module_id, rm.module_id) = m.id
WHERE m.is_active = true
ORDER BY u.name, m.category, m.name
LIMIT 20; -- Limit for readability

-- 2. Permission summary per user
SELECT 
    'USER PERMISSION SUMMARY' as test_section,
    u.name as user_name,
    r.name as role_name,
    un.name as unit_name,
    COUNT(DISTINCT m.id) as total_modules,
    COUNT(CASE WHEN COALESCE(urm.can_read, rm.can_read, false) = true THEN 1 END) as read_permissions,
    COUNT(CASE WHEN COALESCE(urm.can_write, rm.can_write, false) = true THEN 1 END) as write_permissions,
    COUNT(CASE WHEN COALESCE(urm.can_delete, rm.can_delete, false) = true THEN 1 END) as delete_permissions,
    COUNT(CASE WHEN COALESCE(urm.can_approve, rm.can_approve, false) = true THEN 1 END) as approve_permissions,
    COUNT(CASE WHEN urm.id IS NOT NULL THEN 1 END) as custom_permissions,
    COUNT(CASE WHEN rm.id IS NOT NULL AND urm.id IS NULL THEN 1 END) as default_permissions
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
LEFT JOIN modules m ON COALESCE(urm.module_id, rm.module_id) = m.id
WHERE m.is_active = true
GROUP BY u.name, r.name, un.name
ORDER BY u.name;

-- 3. Test specific scenarios
-- Scenario A: HR Admin user in HR Admin unit vs default HR Admin permissions
SELECT 
    'SCENARIO A: HR ADMIN UNIT CUSTOMIZATION' as test_section,
    'Haruno Sakura (HR_ADMIN in HR-ADM unit)' as scenario_description,
    m.category,
    COUNT(*) as modules_in_category,
    COUNT(CASE WHEN urm.can_read = true THEN 1 END) as unit_read_access,
    COUNT(CASE WHEN rm.can_read = true THEN 1 END) as default_read_access,
    COUNT(CASE WHEN urm.can_approve = true THEN 1 END) as unit_approve_access,
    COUNT(CASE WHEN rm.can_approve = true THEN 1 END) as default_approve_access
FROM user_roles ur
JOIN users u ON ur.user_id = u.id AND u.name = 'Haruno Sakura'
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id AND rm.module_id = urm.module_id
LEFT JOIN modules m ON urm.module_id = m.id
WHERE m.is_active = true
GROUP BY m.category
ORDER BY m.category;

-- Scenario B: Employee user in Admin unit
SELECT 
    'SCENARIO B: EMPLOYEE IN ADMIN UNIT' as test_section,
    'Yamanaka Ino (EMPLOYEE in ADM unit)' as scenario_description,
    m.category,
    COUNT(*) as accessible_modules,
    COUNT(CASE WHEN COALESCE(urm.can_read, rm.can_read, false) = true THEN 1 END) as read_access,
    COUNT(CASE WHEN COALESCE(urm.can_write, rm.can_write, false) = true THEN 1 END) as write_access
FROM user_roles ur
JOIN users u ON ur.user_id = u.id AND u.name = 'Yamanaka Ino'
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
LEFT JOIN modules m ON COALESCE(urm.module_id, rm.module_id) = m.id
WHERE m.is_active = true
GROUP BY m.category
HAVING COUNT(*) > 0
ORDER BY m.category;

-- 4. Compare unit-level vs branch-level permissions
SELECT 
    'PERMISSION LEVEL COMPARISON' as test_section,
    'Unit Level Users' as user_type,
    COUNT(DISTINCT ur.user_id) as user_count,
    AVG(perm_count.total_permissions) as avg_permissions_per_user
FROM user_roles ur
JOIN (
    SELECT 
        ur.user_id,
        COUNT(DISTINCT m.id) as total_permissions
    FROM user_roles ur
    LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
    LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
    LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
    LEFT JOIN modules m ON COALESCE(urm.module_id, rm.module_id) = m.id
    WHERE ur.unit_id IS NOT NULL AND m.is_active = true
    GROUP BY ur.user_id
) perm_count ON ur.user_id = perm_count.user_id
WHERE ur.unit_id IS NOT NULL

UNION ALL

SELECT 
    'PERMISSION LEVEL COMPARISON' as test_section,
    'Branch Level Users' as user_type,
    COUNT(DISTINCT ur.user_id) as user_count,
    AVG(perm_count.total_permissions) as avg_permissions_per_user
FROM user_roles ur
JOIN (
    SELECT 
        ur.user_id,
        COUNT(DISTINCT m.id) as total_permissions
    FROM user_roles ur
    LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
    LEFT JOIN modules m ON rm.module_id = m.id
    WHERE ur.unit_id IS NULL AND m.is_active = true
    GROUP BY ur.user_id
) perm_count ON ur.user_id = perm_count.user_id
WHERE ur.unit_id IS NULL;

-- 5. Test the user_effective_permissions view
SELECT 
    'EFFECTIVE PERMISSIONS VIEW TEST' as test_section,
    user_id,
    COUNT(*) as total_module_permissions,
    COUNT(CASE WHEN can_read = true THEN 1 END) as read_permissions,
    COUNT(CASE WHEN can_write = true THEN 1 END) as write_permissions,
    COUNT(CASE WHEN can_delete = true THEN 1 END) as delete_permissions,
    COUNT(CASE WHEN can_approve = true THEN 1 END) as approve_permissions,
    COUNT(CASE WHEN is_customized = true THEN 1 END) as customized_permissions
FROM user_effective_permissions
GROUP BY user_id
ORDER BY user_id;

-- 6. Performance test - measure query execution time
EXPLAIN ANALYZE
SELECT 
    u.name,
    COUNT(DISTINCT m.id) as accessible_modules
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
LEFT JOIN modules m ON COALESCE(urm.module_id, rm.module_id) = m.id
WHERE m.is_active = true
AND (urm.can_read = true OR rm.can_read = true)
GROUP BY u.name;

-- 7. Final validation summary
SELECT 
    'MIGRATION VALIDATION SUMMARY' as final_section,
    COUNT(DISTINCT ur.user_id) as total_users,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) as users_with_units,
    COUNT(CASE WHEN ur.unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(DISTINCT ur.unit_id) as unique_units_used,
    COUNT(DISTINCT urm.id) as total_custom_permissions,
    CASE 
        WHEN COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) > 0 
        THEN 'MIGRATION SUCCESSFUL'
        ELSE 'NO USERS MIGRATED'
    END as migration_status
FROM user_roles ur
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id;