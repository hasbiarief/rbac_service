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
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (21, 6, 8, 1, 2, '2026-01-06 11:17:37.447977', 12);
INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (1, 1, 1, 1, 1, '2025-12-28 17:31:43.824744', 5);
INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (20, 2, 2, 1, 1, '2026-01-05 15:13:53', 7);
INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (22, 3, 9, 1, 1, '2026-01-06 14:49:55.185813', 6);
INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (34, 16, 13, 1, 1, '2026-01-17 23:03:29.063362', 16);
INSERT INTO public.user_roles (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (35, 12, 14, 1, 1, '2026-01-18 04:19:38.239794', 17);


--
-- Name: user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.user_roles_id_seq', 35, true);


--
-- PostgreSQL database dump complete
--

