-- Migration 002: Create companies and branches tables
-- Support for multi-company and hierarchical branch structure

CREATE TABLE IF NOT EXISTS companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS branches (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES branches(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    level INTEGER DEFAULT 0,
    path TEXT DEFAULT '',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, code)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_companies_code ON companies(code);
CREATE INDEX IF NOT EXISTS idx_companies_is_active ON companies(is_active);
CREATE INDEX IF NOT EXISTS idx_branches_company_id ON branches(company_id);
CREATE INDEX IF NOT EXISTS idx_branches_parent_id ON branches(parent_id);
CREATE INDEX IF NOT EXISTS idx_branches_path ON branches(path);
CREATE INDEX IF NOT EXISTS idx_branches_is_active ON branches(is_active);

-- Create triggers for updated_at
CREATE TRIGGER update_companies_updated_at 
    BEFORE UPDATE ON companies 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_branches_updated_at 
    BEFORE UPDATE ON branches 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to update branch path and level
CREATE OR REPLACE FUNCTION update_branch_hierarchy()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level = 0;
        NEW.path = NEW.id::TEXT;
    ELSE
        SELECT level + 1, path || '.' || NEW.id::TEXT
        INTO NEW.level, NEW.path
        FROM branches
        WHERE id = NEW.parent_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for branch hierarchy
CREATE TRIGGER update_branch_hierarchy_trigger
    BEFORE INSERT OR UPDATE ON branches
    FOR EACH ROW
    EXECUTE FUNCTION update_branch_hierarchy();

-- Insert sample companies and branches
INSERT INTO companies (name, code, is_active) VALUES
('PT. Cinta Sejati', 'CINTA', true),
('PT. Teknologi Maju', 'TEKNO', true),
('CV. Dagang Sukses', 'DAGANG', true)
ON CONFLICT (code) DO NOTHING;

-- Insert sample branches
INSERT INTO branches (company_id, parent_id, name, code, is_active) VALUES
(1, NULL, 'Kantor Pusat', 'PUSAT', true),
(1, 1, 'Cabang Jakarta', 'JKT', true),
(1, 1, 'Cabang Surabaya', 'SBY', true),
(1, 2, 'Sub Cabang Jakarta Selatan', 'JKTSEL', true)
ON CONFLICT (company_id, code) DO NOTHING;