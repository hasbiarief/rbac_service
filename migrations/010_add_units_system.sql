-- Migration 010: Add Units System
-- Add unit layer between branch and role for granular access control

-- Create units table
CREATE TABLE IF NOT EXISTS units (
    id BIGSERIAL PRIMARY KEY,
    branch_id BIGINT NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES units(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    level INTEGER DEFAULT 0,
    path TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(branch_id, code)
);

-- Create unit_roles table to map roles to units
CREATE TABLE IF NOT EXISTS unit_roles (
    id BIGSERIAL PRIMARY KEY,
    unit_id BIGINT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(unit_id, role_id)
);

-- Create unit_role_modules table for unit-specific role permissions
CREATE TABLE IF NOT EXISTS unit_role_modules (
    id BIGSERIAL PRIMARY KEY,
    unit_role_id BIGINT NOT NULL REFERENCES unit_roles(id) ON DELETE CASCADE,
    module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    can_read BOOLEAN DEFAULT false,
    can_write BOOLEAN DEFAULT false,
    can_delete BOOLEAN DEFAULT false,
    can_approve BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(unit_role_id, module_id)
);

-- Update user_roles table to include unit_id
ALTER TABLE user_roles ADD COLUMN unit_id BIGINT REFERENCES units(id) ON DELETE SET NULL;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_units_branch_id ON units(branch_id);
CREATE INDEX IF NOT EXISTS idx_units_parent_id ON units(parent_id);
CREATE INDEX IF NOT EXISTS idx_units_path ON units(path);
CREATE INDEX IF NOT EXISTS idx_units_is_active ON units(is_active);
CREATE INDEX IF NOT EXISTS idx_units_code ON units(code);

CREATE INDEX IF NOT EXISTS idx_unit_roles_unit_id ON unit_roles(unit_id);
CREATE INDEX IF NOT EXISTS idx_unit_roles_role_id ON unit_roles(role_id);

CREATE INDEX IF NOT EXISTS idx_unit_role_modules_unit_role_id ON unit_role_modules(unit_role_id);
CREATE INDEX IF NOT EXISTS idx_unit_role_modules_module_id ON unit_role_modules(module_id);

CREATE INDEX IF NOT EXISTS idx_user_roles_unit_id ON user_roles(unit_id);

-- Create triggers for updated_at
CREATE TRIGGER update_units_updated_at 
    BEFORE UPDATE ON units 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_unit_roles_updated_at 
    BEFORE UPDATE ON unit_roles 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_unit_role_modules_updated_at 
    BEFORE UPDATE ON unit_role_modules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update unit hierarchy (similar to branch hierarchy)
CREATE OR REPLACE FUNCTION update_unit_hierarchy()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level = 0;
        NEW.path = NEW.id::TEXT;
    ELSE
        SELECT level + 1, path || '.' || NEW.id::TEXT
        INTO NEW.level, NEW.path
        FROM units
        WHERE id = NEW.parent_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for unit hierarchy
CREATE TRIGGER update_unit_hierarchy_trigger
    BEFORE INSERT OR UPDATE ON units
    FOR EACH ROW
    EXECUTE FUNCTION update_unit_hierarchy();

-- Add comments for documentation
COMMENT ON TABLE units IS 'Units within branches - represents departments, divisions, or teams';
COMMENT ON TABLE unit_roles IS 'Maps roles to specific units for granular access control';
COMMENT ON TABLE unit_role_modules IS 'Unit-specific role permissions that can override default role permissions';
COMMENT ON COLUMN user_roles.unit_id IS 'Optional unit assignment for user role - if NULL, role applies to entire branch';
COMMENT ON COLUMN units.parent_id IS 'Reference to parent unit for hierarchical structure within branch';
COMMENT ON COLUMN units.path IS 'Hierarchical path for efficient tree queries';
COMMENT ON COLUMN units.level IS 'Depth level in unit hierarchy (0 = root unit)';

-- Insert sample units for existing branches
INSERT INTO units (branch_id, parent_id, name, code, description, is_active) VALUES
-- Units for Kantor Pusat (branch_id = 1)
(1, NULL, 'Human Resources', 'HR', 'Departemen Sumber Daya Manusia', true),
(1, NULL, 'Finance & Accounting', 'FIN', 'Departemen Keuangan dan Akuntansi', true),
(1, NULL, 'Information Technology', 'IT', 'Departemen Teknologi Informasi', true),
(1, NULL, 'Operations', 'OPS', 'Departemen Operasional', true),

-- Sub-units for HR
(1, (SELECT id FROM units WHERE code = 'HR' AND branch_id = 1), 'HR Admin', 'HR-ADM', 'Tim Administrasi HR', true),
(1, (SELECT id FROM units WHERE code = 'HR' AND branch_id = 1), 'Recruitment', 'HR-REC', 'Tim Rekrutmen', true),
(1, (SELECT id FROM units WHERE code = 'HR' AND branch_id = 1), 'Payroll', 'HR-PAY', 'Tim Penggajian', true),

-- Sub-units for Finance
(1, (SELECT id FROM units WHERE code = 'FIN' AND branch_id = 1), 'Accounting', 'FIN-ACC', 'Tim Akuntansi', true),
(1, (SELECT id FROM units WHERE code = 'FIN' AND branch_id = 1), 'Treasury', 'FIN-TRS', 'Tim Treasury', true),

-- Units for Cabang Jakarta (branch_id = 2)
(2, NULL, 'Sales', 'SALES', 'Departemen Penjualan', true),
(2, NULL, 'Customer Service', 'CS', 'Layanan Pelanggan', true),
(2, NULL, 'Admin', 'ADM', 'Administrasi Cabang', true),

-- Units for Cabang Surabaya (branch_id = 3)
(3, NULL, 'Sales', 'SALES', 'Departemen Penjualan', true),
(3, NULL, 'Operations', 'OPS', 'Operasional Cabang', true),
(3, NULL, 'Admin', 'ADM', 'Administrasi Cabang', true)
ON CONFLICT (branch_id, code) DO NOTHING;

-- Create default unit-role mappings
-- Map HR roles to HR units
INSERT INTO unit_roles (unit_id, role_id) VALUES
-- HR Admin role to HR units
((SELECT id FROM units WHERE code = 'HR-ADM' AND branch_id = 1), 2), -- HR_ADMIN
((SELECT id FROM units WHERE code = 'HR-REC' AND branch_id = 1), 9), -- RECRUITER
((SELECT id FROM units WHERE code = 'HR-PAY' AND branch_id = 1), 4), -- PAYROLL_OFFICER

-- Manager roles to department units
((SELECT id FROM units WHERE code = 'HR' AND branch_id = 1), 3), -- HR_MANAGER
((SELECT id FROM units WHERE code = 'FIN' AND branch_id = 1), 5), -- LINE_MANAGER
((SELECT id FROM units WHERE code = 'IT' AND branch_id = 1), 12), -- IT_ADMIN

-- Sales units get LINE_MANAGER role
((SELECT id FROM units WHERE code = 'SALES' AND branch_id = 2), 5), -- LINE_MANAGER
((SELECT id FROM units WHERE code = 'SALES' AND branch_id = 3), 5)  -- LINE_MANAGER
ON CONFLICT (unit_id, role_id) DO NOTHING;