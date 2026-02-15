-- Migration 017: Create Application Layer
-- Description: Add application layer above modules for better organization
-- Date: 2026-02-01

-- Create applications table
CREATE TABLE applications (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    icon VARCHAR(100),
    url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for applications
CREATE INDEX idx_applications_code ON applications(code);
CREATE INDEX idx_applications_is_active ON applications(is_active);
CREATE INDEX idx_applications_sort_order ON applications(sort_order);

-- Create plan_applications table (many-to-many relationship)
CREATE TABLE plan_applications (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
    application_id BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    is_included BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(plan_id, application_id)
);

-- Create indexes for plan_applications
CREATE INDEX idx_plan_applications_plan_id ON plan_applications(plan_id);
CREATE INDEX idx_plan_applications_application_id ON plan_applications(application_id);
CREATE INDEX idx_plan_applications_is_included ON plan_applications(is_included);

-- Add application_id column to modules table
ALTER TABLE modules ADD COLUMN application_id BIGINT REFERENCES applications(id);
CREATE INDEX idx_modules_application_id ON modules(application_id);

-- Create trigger for updating updated_at timestamp
CREATE OR REPLACE FUNCTION update_applications_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_applications_updated_at 
    BEFORE UPDATE ON applications 
    FOR EACH ROW 
    EXECUTE FUNCTION update_applications_updated_at_column();

-- Insert initial applications
INSERT INTO applications (name, code, description, icon, url, sort_order, is_active) VALUES
('HRIS', 'HRIS', 'Human Resources Information System - Complete HR management solution', 'fas fa-users', '/hris', 1, true),
('RBAC', 'RBAC', 'Role-Based Access Control - User management and API documentation system', 'fas fa-shield-alt', '/rbac', 2, true);

-- Map existing modules to applications
-- HRIS Application (all modules except 137-146)
UPDATE modules 
SET application_id = (SELECT id FROM applications WHERE code = 'HRIS')
WHERE id NOT BETWEEN 137 AND 146;

-- RBAC Application (modules 137-146)
UPDATE modules 
SET application_id = (SELECT id FROM applications WHERE code = 'RBAC')
WHERE id BETWEEN 137 AND 146;

-- Add applications to all existing subscription plans
INSERT INTO plan_applications (plan_id, application_id, is_included)
SELECT sp.id, a.id, true
FROM subscription_plans sp
CROSS JOIN applications a
WHERE sp.is_active = true AND a.is_active = true;

-- Add comments for documentation
COMMENT ON TABLE applications IS 'Applications that group related modules together';
COMMENT ON TABLE plan_applications IS 'Many-to-many relationship between subscription plans and applications';
COMMENT ON COLUMN modules.application_id IS 'Reference to the application this module belongs to';