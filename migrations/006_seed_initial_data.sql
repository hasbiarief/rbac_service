-- Migration 006: Seed initial data
-- Create initial users, roles, and subscriptions

-- Update existing users with user_identity
UPDATE users SET user_identity = '100000001' WHERE id = 1;

-- Insert additional test users
INSERT INTO users (name, email, user_identity, password_hash, is_active) VALUES
('HR Manager', 'hr@company.com', '100000002', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VcSAg/9qm', true),
('Super Admin', 'superadmin@company.com', '100000003', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VcSAg/9qm', true),
('HR Staff', 'hrstaff@company.com', '100000004', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VcSAg/9qm', true),
('Manager', 'manager@company.com', '100000005', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VcSAg/9qm', true),
('Employee', 'employee@company.com', '100000006', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/VcSAg/9qm', true)
ON CONFLICT (email) DO NOTHING;

-- Assign roles to users
INSERT INTO user_roles (user_id, role_id, company_id, branch_id) VALUES
-- System Admin - full access
(1, 1, 1, NULL),
-- HR Manager - company level access
(2, 2, 1, NULL),
-- Super Admin - full access
(3, 1, 1, NULL),
-- HR Staff - branch level access
(4, 3, 1, 1),
-- Manager - branch level access
(5, 4, 1, 2),
-- Employee - branch level access
(6, 5, 1, 2)
ON CONFLICT (user_id, role_id, company_id, branch_id) DO NOTHING;

-- Assign modules to roles
-- Super Admin gets all modules with full access
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
SELECT 1, id, true, true, true FROM modules WHERE is_active = true
ON CONFLICT (role_id, module_id) DO NOTHING;

-- HR Manager gets HR-related modules with full access
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
SELECT 2, id, true, true, true FROM modules 
WHERE category IN ('Core HR / Master Data', 'Recruitment', 'Attendance & Time', 'Leave Management', 
                   'Payroll & Compensation', 'Performance Management', 'Training & Development', 
                   'Employee Self Service', 'Disciplinary & Relations', 'Offboarding & Exit', 'Reporting & Analytics')
AND is_active = true
ON CONFLICT (role_id, module_id) DO NOTHING;

-- HR Staff gets operational HR modules with limited access
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
SELECT 3, id, true, true, false FROM modules 
WHERE category IN ('Core HR / Master Data', 'Attendance & Time', 'Leave Management', 'Employee Self Service')
AND is_active = true
ON CONFLICT (role_id, module_id) DO NOTHING;

-- Manager gets management modules
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
SELECT 4, id, true, true, false FROM modules 
WHERE category IN ('Core HR / Master Data', 'Attendance & Time', 'Leave Management', 'Performance Management', 'Employee Self Service')
AND is_active = true
ON CONFLICT (role_id, module_id) DO NOTHING;

-- Employee gets self-service modules only
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete)
SELECT 5, id, true, true, false FROM modules 
WHERE category IN ('Employee Self Service')
AND is_active = true
ON CONFLICT (role_id, module_id) DO NOTHING;

-- Create sample subscriptions for existing companies
INSERT INTO subscriptions (company_id, plan_id, status, billing_cycle, start_date, end_date, price, payment_status, last_payment_date, next_payment_date)
SELECT 
    c.id,
    CASE 
        WHEN c.id = 1 THEN 2 -- PT. Cinta Sejati -> Pro Plan
        WHEN c.id = 2 THEN 3 -- PT. Teknologi Maju -> Enterprise
        ELSE 1 -- Others -> Basic
    END as plan_id,
    'active' as status,
    'yearly' as billing_cycle,
    CURRENT_DATE as start_date,
    CURRENT_DATE + INTERVAL '1 year' as end_date,
    CASE 
        WHEN c.id = 1 THEN 2990000.00 -- Pro yearly
        WHEN c.id = 2 THEN 5990000.00 -- Enterprise yearly
        ELSE 990000.00 -- Basic yearly
    END as price,
    'paid' as payment_status,
    CURRENT_DATE as last_payment_date,
    CURRENT_DATE + INTERVAL '1 year' as next_payment_date
FROM companies c
WHERE c.is_active = true
ON CONFLICT (company_id) DO NOTHING;