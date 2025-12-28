-- Migration 005: Create subscription system
-- Tiered pricing with module access control

CREATE TABLE IF NOT EXISTS subscription_plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    price_monthly DECIMAL(10,2) NOT NULL,
    price_yearly DECIMAL(10,2) NOT NULL,
    max_users INTEGER,
    max_branches INTEGER,
    features JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    plan_id BIGINT NOT NULL REFERENCES subscription_plans(id),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'cancelled', 'suspended', 'trial')),
    billing_cycle VARCHAR(10) DEFAULT 'monthly' CHECK (billing_cycle IN ('monthly', 'yearly')),
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'IDR',
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded')),
    last_payment_date DATE,
    next_payment_date DATE,
    auto_renew BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id)
);

CREATE TABLE IF NOT EXISTS plan_modules (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
    module_id BIGINT NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    is_included BOOLEAN DEFAULT true,
    UNIQUE(plan_id, module_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id BIGINT,
    details JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_subscription_plans_name ON subscription_plans(name);
CREATE INDEX IF NOT EXISTS idx_subscription_plans_is_active ON subscription_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_subscriptions_company_id ON subscriptions(company_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions(end_date);
CREATE INDEX IF NOT EXISTS idx_plan_modules_plan_id ON plan_modules(plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_modules_module_id ON plan_modules(module_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

-- Create triggers for updated_at
CREATE TRIGGER update_subscription_plans_updated_at 
    BEFORE UPDATE ON subscription_plans 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at 
    BEFORE UPDATE ON subscriptions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert subscription plans
INSERT INTO subscription_plans (name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features) VALUES
('basic', 'Basic Plan', 'Essential HR features for small businesses', 99000.00, 990000.00, 25, 3, '{"support": "email", "storage": "5GB", "reports": "basic"}'),
('pro', 'Professional Plan', 'Advanced HR features for growing companies', 299000.00, 2990000.00, 100, 10, '{"support": "priority", "storage": "50GB", "reports": "advanced", "api_access": true}'),
('enterprise', 'Enterprise Plan', 'Complete HR solution for large organizations', 599000.00, 5990000.00, NULL, NULL, '{"support": "dedicated", "storage": "unlimited", "reports": "custom", "api_access": true, "white_label": true}')
ON CONFLICT (name) DO NOTHING;

-- Create plan-module mappings
-- Basic Plan: Basic modules only
INSERT INTO plan_modules (plan_id, module_id, is_included)
SELECT 1, m.id, true
FROM modules m
WHERE m.subscription_tier = 'basic' AND m.is_active = true
ON CONFLICT (plan_id, module_id) DO NOTHING;

-- Pro Plan: Basic + Pro modules
INSERT INTO plan_modules (plan_id, module_id, is_included)
SELECT 2, m.id, true
FROM modules m
WHERE m.subscription_tier IN ('basic', 'pro') AND m.is_active = true
ON CONFLICT (plan_id, module_id) DO NOTHING;

-- Enterprise Plan: All modules
INSERT INTO plan_modules (plan_id, module_id, is_included)
SELECT 3, m.id, true
FROM modules m
WHERE m.is_active = true
ON CONFLICT (plan_id, module_id) DO NOTHING;