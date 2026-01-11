-- Validation Script: Verify User Migration Results
-- This script validates the migration and checks for any issues

-- 1. Overall migration status
SELECT 
    'MIGRATION STATUS OVERVIEW' as report_section,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) as users_assigned_to_units,
    COUNT(CASE WHEN ur.unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(*) as total_user_roles,
    ROUND(
        COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) * 100.0 / COUNT(*), 
        2
    ) as migration_percentage
FROM user_roles ur;

-- 2. Detailed user assignments with full hierarchy
SELECT 
    'USER ASSIGNMENTS DETAIL' as report_section,
    u.name as user_name,
    u.email,
    r.name as role_name,
    c.name as company_name,
    b.name as branch_name,
    un.name as unit_name,
    un.code as unit_code,
    CASE 
        WHEN ur.unit_id IS NOT NULL THEN 'UNIT LEVEL'
        ELSE 'BRANCH LEVEL'
    END as assignment_level
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
JOIN companies c ON ur.company_id = c.id
LEFT JOIN branches b ON ur.branch_id = b.id
LEFT JOIN units un ON ur.unit_id = un.id
ORDER BY c.name, b.name, un.name, u.name;

-- 3. Permission validation - check if users have access to modules
SELECT 
    'PERMISSION VALIDATION' as report_section,
    u.name as user_name,
    r.name as role_name,
    un.name as unit_name,
    COUNT(DISTINCT m.id) as total_accessible_modules,
    COUNT(DISTINCT CASE WHEN urm.can_read = true OR rm.can_read = true THEN m.id END) as read_access,
    COUNT(DISTINCT CASE WHEN urm.can_write = true OR rm.can_write = true THEN m.id END) as write_access,
    COUNT(DISTINCT CASE WHEN urm.can_delete = true OR rm.can_delete = true THEN m.id END) as delete_access,
    COUNT(DISTINCT CASE WHEN urm.can_approve = true OR rm.can_approve = true THEN m.id END) as approve_access
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN unit_roles unit_r ON ur.unit_id = unit_r.unit_id AND ur.role_id = unit_r.role_id
LEFT JOIN unit_role_modules urm ON unit_r.id = urm.unit_role_id
LEFT JOIN role_modules rm ON ur.role_id = rm.role_id
LEFT JOIN modules m ON (urm.module_id = m.id OR rm.module_id = m.id) AND m.is_active = true
GROUP BY u.name, r.name, un.name
ORDER BY u.name;

-- 4. Unit utilization report
SELECT 
    'UNIT UTILIZATION' as report_section,
    un.name as unit_name,
    un.code as unit_code,
    b.name as branch_name,
    COUNT(ur.id) as assigned_users,
    COUNT(DISTINCT ur.role_id) as different_roles,
    STRING_AGG(DISTINCT r.name, ', ' ORDER BY r.name) as roles_in_unit
FROM units un
LEFT JOIN user_roles ur ON un.id = ur.unit_id
LEFT JOIN roles r ON ur.role_id = r.id
LEFT JOIN branches b ON un.branch_id = b.id
GROUP BY un.id, un.name, un.code, b.name
ORDER BY b.name, assigned_users DESC, un.name;

-- 5. Role distribution across units
SELECT 
    'ROLE DISTRIBUTION' as report_section,
    r.name as role_name,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) as users_in_units,
    COUNT(CASE WHEN ur.unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(*) as total_users_with_role
FROM roles r
LEFT JOIN user_roles ur ON r.id = ur.role_id
GROUP BY r.name
ORDER BY total_users_with_role DESC, r.name;

-- 6. Check for potential permission conflicts
SELECT 
    'POTENTIAL CONFLICTS' as report_section,
    u.name as user_name,
    COUNT(ur.id) as total_role_assignments,
    STRING_AGG(r.name, ', ' ORDER BY r.name) as assigned_roles,
    STRING_AGG(DISTINCT un.name, ', ' ORDER BY un.name) as assigned_units
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN units un ON ur.unit_id = un.id
GROUP BY u.id, u.name
HAVING COUNT(ur.id) > 1
ORDER BY total_role_assignments DESC;

-- 7. Backup verification
SELECT 
    'BACKUP VERIFICATION' as report_section,
    COUNT(*) as backed_up_records,
    MIN(created_at) as oldest_backup,
    MAX(created_at) as newest_backup
FROM user_roles_backup;

-- 8. Data integrity checks
SELECT 
    'DATA INTEGRITY CHECKS' as report_section,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.id IS NULL THEN 1 END) as orphaned_unit_references,
    COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.branch_id != ur.branch_id THEN 1 END) as unit_branch_mismatches,
    COUNT(CASE WHEN ur.branch_id IS NOT NULL AND b.id IS NULL THEN 1 END) as orphaned_branch_references,
    COUNT(CASE WHEN ur.company_id IS NOT NULL AND c.id IS NULL THEN 1 END) as orphaned_company_references
FROM user_roles ur
LEFT JOIN units un ON ur.unit_id = un.id
LEFT JOIN branches b ON ur.branch_id = b.id
LEFT JOIN companies c ON ur.company_id = c.id;

-- 9. Performance impact assessment
SELECT 
    'PERFORMANCE METRICS' as report_section,
    COUNT(*) as total_user_role_records,
    COUNT(CASE WHEN unit_id IS NOT NULL THEN 1 END) as records_with_unit_joins,
    AVG(CASE WHEN unit_id IS NOT NULL THEN 1.0 ELSE 0.0 END) as unit_join_ratio,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT role_id) as unique_roles,
    COUNT(DISTINCT unit_id) as unique_units_used
FROM user_roles;

-- 10. Migration success summary
SELECT 
    'MIGRATION SUCCESS SUMMARY' as report_section,
    CASE 
        WHEN COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.id IS NULL THEN 1 END) = 0 
        AND COUNT(CASE WHEN ur.unit_id IS NOT NULL AND un.branch_id != ur.branch_id THEN 1 END) = 0
        THEN 'SUCCESS - No data integrity issues found'
        ELSE 'WARNING - Data integrity issues detected'
    END as integrity_status,
    CASE 
        WHEN COUNT(CASE WHEN ur.unit_id IS NOT NULL THEN 1 END) > 0
        THEN 'SUCCESS - Users migrated to units'
        ELSE 'INFO - No users migrated (may be intentional)'
    END as migration_status
FROM user_roles ur
LEFT JOIN units un ON ur.unit_id = un.id;