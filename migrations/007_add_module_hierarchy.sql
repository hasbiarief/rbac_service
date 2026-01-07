-- Migration 007: Add module hierarchy structure
-- This migration adds parent-child relationships to existing modules

-- First, let's add some parent modules (main categories)
INSERT INTO modules (category, name, url, icon, description, parent_id, subscription_tier, is_active) VALUES
-- Main category modules
('System', 'Core HR Management', '/core-hr', 'Building', 'Core HR management system', NULL, 'basic', true),
('System', 'Employee Portal', '/ess', 'User', 'Employee self-service portal', NULL, 'basic', true),
('System', 'Recruitment System', '/recruitment', 'UserPlus', 'Complete recruitment management', NULL, 'pro', true),
('System', 'Attendance System', '/attendance', 'Clock', 'Time and attendance management', NULL, 'pro', true),
('System', 'Payroll System', '/payroll', 'DollarSign', 'Payroll and compensation management', NULL, 'pro', true),
('System', 'Performance System', '/performance', 'TrendingUp', 'Performance management system', NULL, 'enterprise', true),
('System', 'Learning System', '/learning', 'BookOpen', 'Learning and development system', NULL, 'enterprise', true),
('System', 'Reports & Analytics', '/reports', 'BarChart3', 'Comprehensive reporting system', NULL, 'pro', true);

-- Now let's update existing modules to have proper parent relationships
-- Get the IDs of parent modules we just created
-- Core HR Management modules (assuming parent ID = last_insert_id)

-- Update Core HR modules to have parent relationship
UPDATE modules SET parent_id = (
    SELECT id FROM modules WHERE name = 'Core HR Management' AND category = 'System'
) WHERE category = 'Core HR / Master Data';

-- Update ESS modules to have parent relationship  
UPDATE modules SET parent_id = (
    SELECT id FROM modules WHERE name = 'Employee Portal' AND category = 'System'
) WHERE category = 'Employee Self Service';

-- Update Recruitment modules to have parent relationship
UPDATE modules SET parent_id = (
    SELECT id FROM modules WHERE name = 'Recruitment System' AND category = 'System'
) WHERE category = 'Recruitment';

-- Update Attendance modules to have parent relationship
UPDATE modules SET parent_id = (
    SELECT id FROM modules WHERE name = 'Attendance System' AND category = 'System'
) WHERE category = 'Attendance & Time';

-- Add some sub-sub modules (level 3) for demonstration
-- Employee Data sub-modules
INSERT INTO modules (category, name, url, icon, description, parent_id, subscription_tier, is_active) VALUES
('Core HR / Master Data', 'Personal Information', '/core-hr/employees/personal', 'User', 'Manage personal information', 
    (SELECT id FROM modules WHERE name = 'Data Karyawan' AND category = 'Core HR / Master Data'), 'basic', true),
('Core HR / Master Data', 'Employment Details', '/core-hr/employees/employment', 'Briefcase', 'Manage employment details', 
    (SELECT id FROM modules WHERE name = 'Data Karyawan' AND category = 'Core HR / Master Data'), 'basic', true),
('Core HR / Master Data', 'Contact Information', '/core-hr/employees/contact', 'Phone', 'Manage contact information', 
    (SELECT id FROM modules WHERE name = 'Data Karyawan' AND category = 'Core HR / Master Data'), 'basic', true);

-- Organization Structure sub-modules
INSERT INTO modules (category, name, url, icon, description, parent_id, subscription_tier, is_active) VALUES
('Core HR / Master Data', 'Company Structure', '/core-hr/organization/company', 'Building', 'Manage company structure', 
    (SELECT id FROM modules WHERE name = 'Struktur Organisasi' AND category = 'Core HR / Master Data'), 'basic', true),
('Core HR / Master Data', 'Department Management', '/core-hr/organization/departments', 'Building2', 'Manage departments', 
    (SELECT id FROM modules WHERE name = 'Struktur Organisasi' AND category = 'Core HR / Master Data'), 'basic', true),
('Core HR / Master Data', 'Position Hierarchy', '/core-hr/organization/positions', 'TrendingUp', 'Manage position hierarchy', 
    (SELECT id FROM modules WHERE name = 'Struktur Organisasi' AND category = 'Core HR / Master Data'), 'basic', true);

-- ESS Profile sub-modules
INSERT INTO modules (category, name, url, icon, description, parent_id, subscription_tier, is_active) VALUES
('Employee Self Service', 'Basic Profile', '/ess/profile/basic', 'User', 'Update basic profile information', 
    (SELECT id FROM modules WHERE name = 'Update Profil' AND category = 'Employee Self Service'), 'basic', true),
('Employee Self Service', 'Emergency Contacts', '/ess/profile/emergency', 'Phone', 'Manage emergency contacts', 
    (SELECT id FROM modules WHERE name = 'Update Profil' AND category = 'Employee Self Service'), 'basic', true),
('Employee Self Service', 'Bank Information', '/ess/profile/bank', 'CreditCard', 'Update bank information', 
    (SELECT id FROM modules WHERE name = 'Update Profil' AND category = 'Employee Self Service'), 'basic', true);

-- Leave Request sub-modules
INSERT INTO modules (category, name, url, icon, description, parent_id, subscription_tier, is_active) VALUES
('Employee Self Service', 'Annual Leave', '/ess/requests/annual', 'Calendar', 'Request annual leave', 
    (SELECT id FROM modules WHERE name = 'Pengajuan Cuti & Izin' AND category = 'Employee Self Service'), 'basic', true),
('Employee Self Service', 'Sick Leave', '/ess/requests/sick', 'Heart', 'Request sick leave', 
    (SELECT id FROM modules WHERE name = 'Pengajuan Cuti & Izin' AND category = 'Employee Self Service'), 'basic', true),
('Employee Self Service', 'Permission', '/ess/requests/permission', 'Clock', 'Request permission', 
    (SELECT id FROM modules WHERE name = 'Pengajuan Cuti & Izin' AND category = 'Employee Self Service'), 'basic', true);

-- Add indexes for better performance on hierarchical queries
CREATE INDEX IF NOT EXISTS idx_modules_parent_id ON modules(parent_id);
CREATE INDEX IF NOT EXISTS idx_modules_category_parent ON modules(category, parent_id);

-- Add comments
COMMENT ON COLUMN modules.parent_id IS 'Reference to parent module for hierarchical structure';