--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Homebrew)
-- Dumped by pg_dump version 16.9 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: subscription_plans; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (1, 'basic', 'Basic Plan', 'Essential HR features for small businesses', 99000.00, 990000.00, 25, 3, '{"reports": "basic", "storage": "5GB", "support": "email"}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');
INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (2, 'pro', 'Professional Plan', 'Advanced HR features for growing companies', 299000.00, 2990000.00, 100, 10, '{"reports": "advanced", "storage": "50GB", "support": "priority", "api_access": true}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');
INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (3, 'enterprise', 'Enterprise Plan', 'Complete HR solution for large organizations', 599000.00, 5990000.00, NULL, NULL, '{"reports": "custom", "storage": "unlimited", "support": "dedicated", "api_access": true, "white_label": true}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.subscription_plans_id_seq', 3, true);


--
-- PostgreSQL database dump complete
--

