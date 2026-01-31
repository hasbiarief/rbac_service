-- Fix: Add API Documentation modules to all subscription plans
-- This ensures API Documentation modules are available to all subscription tiers

-- Insert API Documentation modules into plan_modules for all plans
INSERT INTO plan_modules (plan_id, module_id, is_included) VALUES
-- Basic Plan (plan_id = 1)
(1, 139, true),  -- API Documentation
(1, 140, true),  -- Collections Management
(1, 141, true),  -- Endpoints Management
(1, 142, true),  -- Environments Management
(1, 143, true),  -- Export Documentation

-- Professional Plan (plan_id = 2)
(2, 139, true),  -- API Documentation
(2, 140, true),  -- Collections Management
(2, 141, true),  -- Endpoints Management
(2, 142, true),  -- Environments Management
(2, 143, true),  -- Export Documentation

-- Enterprise Plan (plan_id = 3)
(3, 139, true),  -- API Documentation
(3, 140, true),  -- Collections Management
(3, 141, true),  -- Endpoints Management
(3, 142, true),  -- Environments Management
(3, 143, true)   -- Export Documentation

ON CONFLICT (plan_id, module_id) DO UPDATE SET
    is_included = EXCLUDED.is_included;

-- Verify the insertion
SELECT 
    sp.name as plan_name,
    m.name as module_name,
    m.category,
    pm.is_included
FROM plan_modules pm
JOIN subscription_plans sp ON pm.plan_id = sp.id
JOIN modules m ON pm.module_id = m.id
WHERE pm.module_id IN (139, 140, 141, 142, 143)
ORDER BY sp.id, m.id;