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
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (2, 2, 3, 'active', 'yearly', '2025-12-28', '2026-12-28', 5990000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2025-12-28 17:31:43.824744');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (3, 3, 2, 'active', 'monthly', '2025-12-28', '2026-01-31', 299000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2025-12-31 15:18:49.681511');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (8, 4, 3, 'active', 'monthly', '2025-12-31', '2026-01-31', 599000.00, 'IDR', 'paid', '2025-12-31', '2026-01-31', false, '2025-12-31 15:03:49.811229', '2025-12-31 15:26:34.213118');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (6, 5, 2, 'active', 'yearly', '2025-12-31', '2026-12-31', 2990000.00, 'IDR', 'paid', '2025-12-31', '2026-12-31', true, '2025-12-31 15:01:41.777448', '2025-12-31 15:32:44.393663');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (1, 1, 3, 'active', 'yearly', '2025-12-28', '2026-01-31', 5990000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2026-01-04 11:50:44.296087');


--
-- Name: subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.subscriptions_id_seq', 9, true);


--
-- PostgreSQL database dump complete
--

