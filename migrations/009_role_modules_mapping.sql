-- Migration 009: Role Modules Mapping
-- Map all roles to modules based on role_akses.txt

-- Clear existing role_modules except SUPER_ADMIN (role_id = 1)
DELETE FROM role_modules WHERE role_id != 1;

-- HR_ADMIN (role_id = 2) Mapping
-- Asset & Facility: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, false
FROM modules WHERE category = 'Asset & Facility' AND is_active = true;

-- Attendance & Time: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, true
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Dashboard & Analytic: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, false, false, false
FROM modules WHERE category = 'Dashboard & Analytic' AND is_active = true;

-- Disciplinary & Relations: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, true
FROM modules WHERE category = 'Disciplinary & Relations' AND is_active = true;

-- Employee Self Service: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, false
FROM modules WHERE category = 'Employee Self Service' AND is_active = true;

-- Leave Management: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, true
FROM modules WHERE category = 'Leave Management' AND is_active = true;

-- Core HR / Master Data: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, false
FROM modules WHERE category = 'Core HR / Master Data' AND is_active = true;

-- Offboarding & Exit: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, true
FROM modules WHERE category = 'Offboarding & Exit' AND is_active = true;

-- Payroll & Compensation: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, false, false, false
FROM modules WHERE category = 'Payroll & Compensation' AND is_active = true;

-- Performance Management: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, true
FROM modules WHERE category = 'Performance Management' AND is_active = true;

-- Recruitment: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, false
FROM modules WHERE category = 'Recruitment' AND is_active = true;

-- Reporting & Analytics: R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- Training & Development: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 2, id, true, true, true, false
FROM modules WHERE category = 'Training & Development' AND is_active = true;

-- HR_MANAGER (role_id = 3) Mapping
-- Asset & Facility: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Asset & Facility' AND is_active = true;

-- Attendance & Time: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, true
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Dashboard & Analytic: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Dashboard & Analytic' AND is_active = true;

-- Disciplinary & Relations: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, true
FROM modules WHERE category = 'Disciplinary & Relations' AND is_active = true;

-- Leave Management: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, true
FROM modules WHERE category = 'Leave Management' AND is_active = true;

-- Core HR / Master Data: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Core HR / Master Data' AND is_active = true;

-- Offboarding & Exit: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, true
FROM modules WHERE category = 'Offboarding & Exit' AND is_active = true;

-- Payroll & Compensation: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Payroll & Compensation' AND is_active = true;

-- Performance Management: R=✓, W=✓, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, true, false, true
FROM modules WHERE category = 'Performance Management' AND is_active = true;

-- Recruitment: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, true
FROM modules WHERE category = 'Recruitment' AND is_active = true;

-- Reporting & Analytics: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- Training & Development: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 3, id, true, false, false, false
FROM modules WHERE category = 'Training & Development' AND is_active = true;

-- PAYROLL_OFFICER (role_id = 4) Mapping
-- Attendance & Time: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 4, id, true, false, false, false
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Leave Management: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 4, id, true, false, false, false
FROM modules WHERE category = 'Leave Management' AND is_active = true;

-- Payroll & Compensation: R=✓, W=✓, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 4, id, true, true, false, true
FROM modules WHERE category = 'Payroll & Compensation' AND is_active = true;

-- Reporting & Analytics: R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 4, id, true, true, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- LINE_MANAGER (role_id = 5) Mapping
-- Attendance & Time: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, false, false, true
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Leave Management: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, false, false, true
FROM modules WHERE category = 'Leave Management' AND is_active = true;

-- Performance Management: R=✓, W=✓, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, true, false, true
FROM modules WHERE category = 'Performance Management' AND is_active = true;

-- Training & Development: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, false, false, false
FROM modules WHERE category = 'Training & Development' AND is_active = true;

-- Offboarding & Exit: R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, false, false, true
FROM modules WHERE category = 'Offboarding & Exit' AND is_active = true;

-- Dashboard & Analytic: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 5, id, true, false, false, false
FROM modules WHERE category = 'Dashboard & Analytic' AND is_active = true;

-- EMPLOYEE (role_id = 8) Mapping
-- Employee Self Service: R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, true, false, false
FROM modules WHERE category = 'Employee Self Service' AND is_active = true;

-- Attendance & Time: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, false, false, false
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Leave Management: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, false, false, false
FROM modules WHERE category = 'Leave Management' AND is_active = true;

-- Performance Management: R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, true, false, false
FROM modules WHERE category = 'Performance Management' AND is_active = true;

-- Training & Development: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, false, false, false
FROM modules WHERE category = 'Training & Development' AND is_active = true;

-- Payroll & Compensation (Slip Gaji): R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 8, id, true, false, false, false
FROM modules WHERE category = 'Payroll & Compensation' AND name LIKE '%Payslip%' AND is_active = true;

-- RECRUITER (role_id = 9) Mapping
-- Recruitment: R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 9, id, true, true, false, false
FROM modules WHERE category = 'Recruitment' AND is_active = true;

-- Onboarding (part of Recruitment): R=✓, W=✓, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 9, id, true, true, false, false
FROM modules WHERE name LIKE '%Onboarding%' AND is_active = true;

-- Reporting & Analytics: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 9, id, true, false, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- ASSET_OFFICER (role_id = 10) Mapping
-- Asset & Facility: R=✓, W=✓, D=✓, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 10, id, true, true, true, true
FROM modules WHERE category = 'Asset & Facility' AND is_active = true;

-- Offboarding (Asset Return): R=✓, W=✗, D=✗, A=✓
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 10, id, true, false, false, true
FROM modules WHERE category = 'Offboarding & Exit' AND name LIKE '%Asset%' AND is_active = true;

-- Reporting & Analytics: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 10, id, true, false, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- AUDITOR (role_id = 11) Mapping
-- Reporting & Analytics: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 11, id, true, false, false, false
FROM modules WHERE category = 'Reporting & Analytics' AND is_active = true;

-- Payroll & Compensation: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 11, id, true, false, false, false
FROM modules WHERE category = 'Payroll & Compensation' AND is_active = true;

-- Attendance & Time: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 11, id, true, false, false, false
FROM modules WHERE category = 'Attendance & Time' AND is_active = true;

-- Disciplinary: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 11, id, true, false, false, false
FROM modules WHERE category = 'Disciplinary & Relations' AND is_active = true;

-- Audit Log: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 11, id, true, false, false, false
FROM modules WHERE name LIKE '%Audit%' AND is_active = true;

-- IT_ADMIN (role_id = 12) Mapping
-- System & Security: R=✓, W=✓, D=✓, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 12, id, true, true, true, false
FROM modules WHERE category = 'System Management' AND is_active = true;

-- Master Data (technical): R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 12, id, true, false, false, false
FROM modules WHERE category = 'Core HR / Master Data' AND name LIKE '%System%' AND is_active = true;

-- Audit Log: R=✓, W=✗, D=✗, A=✗
INSERT INTO role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve)
SELECT 12, id, true, false, false, false
FROM modules WHERE name LIKE '%Audit%' AND is_active = true;