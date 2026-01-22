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
-- Data for Name: branches; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (1, 1, NULL, 'Pusat', 'PST', 0, '1', true, '2025-12-30 16:29:01.228679', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (2, 1, 1, 'Area 1 Jakarta', 'JKT', 1, '1.2', true, '2025-12-30 16:29:01.242648', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (3, 1, 1, 'Area 2 Bandung', 'BDG', 1, '1.3', true, '2025-12-30 16:29:01.242648', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (4, 1, 2, 'Jakarta Pusat', 'JKTPST', 2, '1.2.4', true, '2025-12-30 16:29:01.245411', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (5, 1, 2, 'Jakarta Selatan', 'JKTSEL', 2, '1.2.5', true, '2025-12-30 16:29:01.245411', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (6, 1, 3, 'Bandung', 'BDG01', 2, '1.3.6', true, '2025-12-30 16:29:01.245411', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (7, 1, 3, 'Cimahi', 'CMH', 2, '1.3.7', true, '2025-12-30 16:29:01.245411', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (8, 1, 4, 'Menteng', 'MTG', 3, '1.2.4.8', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (9, 1, 4, 'Banteng', 'BTG', 3, '1.2.4.9', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (10, 1, 6, 'Cilaki', 'CLK', 3, '1.3.6.10', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (11, 1, 7, 'Sangkuriang', 'SKR', 3, '1.3.7.11', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (12, 1, 7, 'Cisangkan', 'CSK', 3, '1.3.7.12', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (13, 1, 7, 'Ciageran', 'CGR', 3, '1.3.7.13', true, '2025-12-30 16:29:01.247201', '2025-12-30 16:31:18.749469');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (14, 1, 7, 'Padalarang', 'PDL', 3, '1.3.7.14', true, '2025-12-31 10:50:46.830701', '2025-12-31 10:50:46.830701');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (15, 1, 7, 'Batujajar', 'BTJ', 3, '1.3.7.15', true, '2025-12-31 10:51:05.224477', '2025-12-31 10:51:05.224477');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (16, 1, 1, 'Area 3 Surabaya', 'SBY', 1, '1.16', true, '2025-12-31 10:53:18.307166', '2025-12-31 10:53:18.307166');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (17, 1, 16, 'Surabaya Utara', 'SBYUTR', 2, '1.16.17', true, '2025-12-31 10:53:38.462547', '2025-12-31 10:53:38.462547');
INSERT INTO public.branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (18, 1, 1, 'Cabang Review', 'REVIEW', 1, '1.18', true, '2026-01-03 16:21:57.957933', '2026-01-03 16:21:57.957933');


--
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.branches_id_seq', 18, true);


--
-- PostgreSQL database dump complete
--

