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
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (1, 'SUPER_ADMIN', 'Mengelola seluruh sistem, konfigurasi, keamanan, dan akses tanpa batas', true, '2025-12-28 17:31:43.767333', '2026-01-05 08:03:20.819173');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (8, 'EMPLOYEE', 'Mengakses layanan HR mandiri dan mengelola data pribadi', true, '2026-01-05 15:01:32', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (5, 'LINE_MANAGER', 'Mengelola dan menyetujui aktivitas HR untuk tim yang dipimpin', true, '2025-12-28 17:31:43.767333', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (3, 'HR_MANAGER', 'Mengawasi, menyetujui, dan menganalisis proses HR secara strategis', true, '2025-12-28 17:31:43.767333', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (9, 'RECRUITER', 'Mengelola proses rekrutmen dari perencanaan hingga onboarding', true, '2026-01-05 08:03:20.819173', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (4, 'PAYROLL_OFFICER', 'Mengelola penggajian, kompensasi, pajak, dan kepatuhan payroll', true, '2025-12-28 17:31:43.767333', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (2, 'HR_ADMIN', 'Menjalankan operasional HR harian dan mengelola data serta proses HR', true, '2025-12-28 17:31:43.767333', '2026-01-05 08:04:40.161024');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (10, 'ASSET_OFFICER', 'Mengelola aset, inventaris, dan fasilitas perusahaan', true, '2026-01-05 15:05:21', '2026-01-05 08:05:47.899374');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (11, 'AUDITOR', 'Mengakses laporan dan data HR secara read-only untuk audit', true, '2026-01-05 15:06:11', '2026-01-05 15:06:12');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (12, 'IT_ADMIN', 'Mengelola keamanan, user, dan konfigurasi teknis sistem', true, '2026-01-05 15:06:46', '2026-01-05 15:06:47');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (13, 'CONSOLE ADMIN', 'Console admin role for manage system', true, '2026-01-17 13:23:44.581398', '2026-01-17 13:31:03.793899');
INSERT INTO public.roles (id, name, description, is_active, created_at, updated_at) VALUES (14, 'TEST COPY', 'test copy role', true, '2026-01-17 13:23:44.581398', '2026-01-17 13:31:03.793899');


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.roles_id_seq', 14, true);


--
-- PostgreSQL database dump complete
--

