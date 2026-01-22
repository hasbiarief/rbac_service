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
-- Data for Name: companies; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (2, 'PT. Teknologi Maju', 'TEKNO', true, '2025-12-28 17:31:43.740456', '2025-12-28 17:31:43.740456');
INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (3, 'CV. Dagang Sukses', 'DAGANG', true, '2025-12-28 17:31:43.740456', '2025-12-28 17:31:43.740456');
INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (4, 'PT. Test API', 'TESTAPI', true, '2025-12-28 18:28:35.489431', '2025-12-28 18:28:35.489431');
INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (5, 'PT. Test Company Tech', 'TEST', true, '2025-12-30 15:06:38.441337', '2025-12-30 15:06:38.441337');
INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (1, 'PT. Huminor Tech', 'HMN', true, '2025-12-28 17:31:43.740456', '2026-01-02 07:52:15.14697');
INSERT INTO public.companies (id, name, code, is_active, created_at, updated_at) VALUES (6, 'PT. Test Review', 'REVIEW', true, '2026-01-03 16:21:35.740167', '2026-01-03 16:21:35.740167');


--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.companies_id_seq', 6, true);


--
-- PostgreSQL database dump complete
--

