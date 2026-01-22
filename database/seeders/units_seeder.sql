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
-- Data for Name: units; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (1, 1, NULL, 'Human Resources', 'HR', 'Departemen Sumber Daya Manusia', 0, '1', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (2, 1, NULL, 'Finance & Accounting', 'FIN', 'Departemen Keuangan dan Akuntansi', 0, '2', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (3, 1, NULL, 'Information Technology', 'IT', 'Departemen Teknologi Informasi', 0, '3', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (4, 1, NULL, 'Operations', 'OPS', 'Departemen Operasional', 0, '4', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (5, 1, NULL, 'HR Admin', 'HR-ADM', 'Tim Administrasi HR', 0, '5', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (6, 1, NULL, 'Recruitment', 'HR-REC', 'Tim Rekrutmen', 0, '6', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (7, 1, NULL, 'Payroll', 'HR-PAY', 'Tim Penggajian', 0, '7', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (8, 1, NULL, 'Accounting', 'FIN-ACC', 'Tim Akuntansi', 0, '8', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (9, 1, NULL, 'Treasury', 'FIN-TRS', 'Tim Treasury', 0, '9', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (10, 2, NULL, 'Sales', 'SALES', 'Departemen Penjualan', 0, '10', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (11, 2, NULL, 'Customer Service', 'CS', 'Layanan Pelanggan', 0, '11', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (12, 2, NULL, 'Admin', 'ADM', 'Administrasi Cabang', 0, '12', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (13, 3, NULL, 'Sales', 'SALES', 'Departemen Penjualan', 0, '13', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (14, 3, NULL, 'Operations', 'OPS', 'Operasional Cabang', 0, '14', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (15, 3, NULL, 'Admin', 'ADM', 'Administrasi Cabang', 0, '15', true, '2026-01-10 21:49:31.129922', '2026-01-10 21:49:31.129922');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (16, 1, NULL, 'Dewa', 'DWA', 'Unit Dewa', 0, '16', true, '2026-01-17 06:47:35.939791', '2026-01-17 06:47:35.939791');
INSERT INTO public.units (id, branch_id, parent_id, name, code, description, level, path, is_active, created_at, updated_at) VALUES (17, 1, NULL, 'Test Unit', 'TEU', 'Test Unit only', 0, '17', true, '2026-01-18 04:14:42.285897', '2026-01-18 04:14:42.285897');


--
-- Name: units_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.units_id_seq', 17, true);


--
-- PostgreSQL database dump complete
--

