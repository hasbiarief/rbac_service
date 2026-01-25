-- Add API Documentation permissions for CONSOLE ADMIN role (role_id=13)
-- This allows users with CONSOLE ADMIN role to access API Documentation System

-- API Documentation Main Module (139)
INSERT INTO public.role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve) 
VALUES (13, 139, true, true, true, true);

-- API Documentation Collections Module (140)
INSERT INTO public.role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve) 
VALUES (13, 140, true, true, true, true);

-- API Documentation Endpoints Module (141)
INSERT INTO public.role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve) 
VALUES (13, 141, true, true, true, true);

-- API Documentation Environments Module (142)
INSERT INTO public.role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve) 
VALUES (13, 142, true, true, true, true);

-- API Documentation Export Module (143)
INSERT INTO public.role_modules (role_id, module_id, can_read, can_write, can_delete, can_approve) 
VALUES (13, 143, true, true, true, true);