-- Rollback Script: Revert User Migration to Units
-- This script can be used to rollback the migration if needed

-- 1. Verify backup exists
SELECT 
    'BACKUP VERIFICATION' as status,
    COUNT(*) as backup_records,
    MIN(created_at) as oldest_record,
    MAX(created_at) as newest_record
FROM user_roles_backup;

-- 2. Show current state before rollback
SELECT 
    'CURRENT STATE BEFORE ROLLBACK' as status,
    COUNT(CASE WHEN unit_id IS NOT NULL THEN 1 END) as users_with_units,
    COUNT(CASE WHEN unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(*) as total_user_roles
FROM user_roles;

-- 3. Rollback confirmation prompt (comment out the rollback commands below to execute)
SELECT 
    'ROLLBACK READY' as status,
    'Uncomment the rollback commands below to proceed' as instruction;

/*
-- UNCOMMENT THE SECTION BELOW TO EXECUTE ROLLBACK

-- 4. Reset unit assignments to NULL for backed up records
UPDATE user_roles 
SET unit_id = NULL 
WHERE id IN (SELECT id FROM user_roles_backup);

-- 5. Verify rollback
SELECT 
    'ROLLBACK COMPLETED' as status,
    COUNT(CASE WHEN unit_id IS NOT NULL THEN 1 END) as users_with_units,
    COUNT(CASE WHEN unit_id IS NULL THEN 1 END) as users_at_branch_level,
    COUNT(*) as total_user_roles
FROM user_roles;

-- 6. Show rollback details
SELECT 
    'ROLLBACK DETAILS' as status,
    u.name as user_name,
    r.name as role_name,
    b.name as branch_name,
    ur.unit_id as current_unit_id,
    'ROLLED BACK TO BRANCH LEVEL' as new_status
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
LEFT JOIN branches b ON ur.branch_id = b.id
WHERE ur.id IN (SELECT id FROM user_roles_backup)
ORDER BY u.name;

-- END OF ROLLBACK SECTION
*/

-- 7. Alternative: Selective rollback for specific users
-- Example: Rollback only EMPLOYEE users
/*
UPDATE user_roles 
SET unit_id = NULL 
WHERE role_id = (SELECT id FROM roles WHERE name = 'EMPLOYEE')
AND id IN (SELECT id FROM user_roles_backup);
*/

-- 8. Cleanup backup table (only run after confirming rollback is successful)
/*
DROP TABLE IF EXISTS user_roles_backup;
*/