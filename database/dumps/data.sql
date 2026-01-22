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
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (12, 'Aburame Shino', 'shino@company.com', '100000009', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2026-01-05 14:42:22', '2026-01-05 07:57:32.172626', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (11, 'Inuzaka Kiba', 'kiba@company.com', '100000008', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2026-01-05 14:40:23', '2026-01-05 07:57:37.1506', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (10, 'Akimichi Chouji', 'chouji@comapny.com', '100000007', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2026-01-05 14:32:25', '2026-01-05 07:57:43.419076', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (6, 'Yamanaka Ino', 'ino@company.com', '100000006', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.824744', '2026-01-05 07:57:52.604384', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (3, 'Sasuke Uchiha', 'sasuke@company.com', '100000003', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.824744', '2026-01-05 07:58:28.553687', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (4, 'Shikamaru Nara', 'shikamaru@company.com', '100000004', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.824744', '2026-01-05 07:58:28.553687', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (2, 'Haruno Sakura', 'sakura@company.com', '100000002', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.824744', '2026-01-05 07:58:28.553687', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (7, 'Hyuga Neji', 'neji@comapny.com', '100000999', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 18:25:57.864074', '2026-01-05 07:59:11.073668', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (5, 'Hyuga Hinata', 'hinata@company.com', '100000005', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.824744', '2026-01-05 07:59:11.073668', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (1, 'Uzumaki Naruto', 'naruto@company.com', '100000001', '$2a$10$CyG6BmkdxPkYTKnJ6geffOXh6QwQq9V8R9oZn5eBCX6G/zX/5V6vy', true, '2025-12-28 17:31:43.688095', '2026-01-06 03:46:11.617987', NULL);
INSERT INTO public.users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at, deleted_at) VALUES (16, 'Hasbi Due', 'hasbi@company.com', '800000001', '$2a$12$t1IUha3ZCj0E4GgcFLcajeKYf.AuRQXryjTSaYfz/whndUdNdcXZy', true, '2026-01-17 09:40:46.254582', '2026-01-17 09:40:46.254582', NULL);


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (1, 1, 'CREATE', 'test_resource', 123, '{"test_field": "test_value", "additional_info": "This is a manual audit log"}', '192.168.1.100', 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', true, NULL, '2025-12-29 15:27:07.183173');
INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (2, 1, 'UPDATE', 'user', NULL, '{}', NULL, NULL, true, NULL, '2025-12-29 15:27:28.887374');
INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (3, 1, 'DELETE', 'user', 456, '{"reason": "Account cleanup", "deleted_by": "admin"}', '10.0.0.1', 'PostmanRuntime/7.32.3', true, NULL, '2025-12-29 16:18:42.105064');
INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (4, 3, 'manual_test', 'test_resource', 123, '{}', NULL, NULL, true, NULL, '2025-12-29 16:22:12.547147');
INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (5, 16, 'manual_test', 'test_resource', 123, '{"url": "/api/v1/test", "method": "POST", "status": "success", "message": "Manual audit log for testing", "status_code": 200, "user_identity": "800000001"}', NULL, NULL, true, NULL, '2026-01-21 13:39:53.508874');
INSERT INTO public.audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, success, error_message, created_at) VALUES (6, 16, 'manual_test', 'test_resource', 123, '{"url": "/api/v1/test", "method": "POST", "status": "success", "message": "Manual audit log for testing", "status_code": 200, "user_identity": "800000001"}', NULL, NULL, true, NULL, '2026-01-21 14:03:24.365024');


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
-- Data for Name: modules; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (67, 'Asset & Facility', 'Facility Booking', '/asset/booking', 'Calendar', 'Booking fasilitas perusahaan', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (69, 'Asset & Facility', 'Inventory Management', '/asset/inventory', 'Package', 'Kelola inventori perusahaan', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (68, 'Asset & Facility', 'Vehicle Management', '/asset/vehicle', 'Truck', 'Kelola kendaraan perusahaan', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (65, 'Asset & Facility', 'Asset Assignment', '/asset/assignment', 'UserCheck', 'Kelola penugasan aset ke karyawan', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (40, 'Leave Management', 'Cuti CF & Exp', '/leave/carry-forward', 'ArrowRight', 'Kelola carry forward dan expired cuti', 128, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-07 04:08:31.16454');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (41, 'Leave Management', 'Kalender Cuti Tim', '/leave/calendar', 'Calendar', 'Kalender cuti tim dan departemen', 128, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:43:08.274204');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (45, 'Performance Management', 'Goal Setting', '/performance/goals', 'Flag', 'Kelola goal setting karyawan', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (47, 'Performance Management', 'Succession Planning', '/performance/succession', 'Users', 'Perencanaan suksesi jabatan', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (42, 'Performance Management', 'KPI & Target', '/performance/kpi', 'Target', 'Kelola KPI dan target karyawan', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (43, 'Performance Management', 'Performance Review', '/performance/review', 'Star', 'Kelola review performa karyawan', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (44, 'Performance Management', '360 Degree Feedback', '/performance/360-feedback', 'Users', 'Sistem feedback 360 derajat', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (46, 'Performance Management', 'Career Development', '/performance/career', 'TrendingUp', 'Kelola pengembangan karir', 97, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:40.167957');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (49, 'Training & Development', 'Training Schedule', '/training/schedule', 'Calendar', 'Kelola jadwal pelatihan', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (48, 'Training & Development', 'Training Program', '/training/programs', 'BookOpen', 'Kelola program pelatihan', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (50, 'Training & Development', 'Training Evaluation', '/training/evaluation', 'CheckSquare', 'Evaluasi hasil pelatihan', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (37, 'Leave Management', 'Pengajuan Cuti & Izin', '/leave/request', 'FileText', 'Kelola pengajuan cuti dan izin karyawan', 128, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:22:44.324031');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (36, 'Leave Management', 'Jenis Cuti', '/leave/types', 'Calendar', 'Kelola jenis-jenis cuti', 128, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:22:44.324031');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (38, 'Leave Management', 'Approval Workflow', '/leave/approval', 'CheckCircle', 'Kelola alur persetujuan cuti', 128, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:22:44.324031');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (39, 'Leave Management', 'Saldo & Kuota Cuti', '/leave/balance', 'BarChart3', 'Kelola saldo dan kuota cuti karyawan', 128, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:22:44.324031');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (132, 'Training & Development', 'Training & Development', '/training', 'CheckSquare', 'Manage Training & Development System', NULL, 'pro', true, '2026-01-02 07:49:14.067716', '2026-01-07 04:05:58.788413');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (133, 'Recruitment', 'Workforce & Headcount', '/recruitment/workforce-headcount', 'Users', 'Perencanaan biaya dan kebutuhan sumber daya', 94, 'pro', true, '2026-01-08 10:56:20', '2026-01-08 08:07:09.214931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (104, 'Master Data', 'Branch Management', '/core-hr/master_data/branch', 'Building2', 'Manage departments', 92, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (6, 'Master Data', 'Grade & Level', '/core-hr/master_data/grades', 'TrendingUp', 'Manage employment grade and level', 92, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 09:02:45.593605');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (136, 'Test Category', 'Test Module', '/test-module', 'TestIcon', 'Test module description', NULL, 'enterprise', true, '2026-01-17 12:56:58.841821', '2026-01-17 12:56:58.841821');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (129, 'Offboarding & Exit', 'Offboarding & Exit', '/offboarding', 'Award', 'Manage Offboarding & Exit', NULL, 'enterprise', true, '2026-01-02 07:43:50.262836', '2026-01-07 04:06:13.366355');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (99, 'Disciplinary & Relations', 'Disciplinary & Relations', '/disciplinary', 'BarChart3', 'Manage Disciplinary & Relations', NULL, 'enterprise', true, '2025-12-29 13:52:18.807576', '2026-01-07 04:08:40.665255');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (98, 'Asset & Facility', 'Asset & Facility', '/asset', 'BookOpen', 'Manage Asset & Facility', NULL, 'enterprise', true, '2025-12-29 13:52:18.807576', '2026-01-07 04:08:48.088001');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (51, 'Training & Development', 'Competency Management', '/training/competency', 'Award', 'Kelola kompetensi karyawan', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (52, 'Training & Development', 'Learning Management', '/training/learning', 'GraduationCap', 'Sistem manajemen pembelajaran', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (53, 'Training & Development', 'Training Budget', '/training/budget', 'DollarSign', 'Kelola budget pelatihan', 132, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:49:32.681647');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (14, 'Master Data', 'Family & Biodata', '/core-hr/employees_management/family', 'Users', 'Kelola data keluarga karyawan', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:33:47.233049');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (88, 'System & Security', 'Audit Log', '/system/audit', 'FileText', 'Log aktivitas sistem', 131, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:05:13.642783');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (91, 'System & Security', 'Backup & Restore', '/system/backup', 'Database', 'Backup dan restore data', 131, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:05:25.054402');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (95, 'Attendance & Time', 'Attendance System', '/attendance', 'Clock', 'Time and attendance management', NULL, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-02 12:13:24.287804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (112, 'Dashboard & Analytic', 'Dashboard Overview', '/dashboard', 'Test', 'Overview dashboard', NULL, 'pro', true, '2025-12-30 13:44:23.518858', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (128, 'Leave Management', 'Manage Leave Management', '/leave', 'BookOpen', 'Manage Leave Management', NULL, 'enterprise', true, '2026-01-02 07:42:45.460885', '2026-01-02 12:22:44.324031');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (3, 'Master Data', 'Jabatan & Posisi', '/core-hr/master_data/positions', 'Briefcase', 'Kelola jabatan dan posisi kerja', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (4, 'Master Data', 'Department', '/core-hr/master_data/departments', 'Building2', 'Kelola departemen perusahaan', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (7, 'Master Data', 'Work Location', '/core-hr/master_data/locations', 'MapPin', 'Kelola lokasi kerja', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (17, 'Master Data', 'Bank Reference', '/core-hr/master_data/bank-accounts', 'CreditCard', 'Kelola rekening bank karyawan', 92, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:36:13.955167');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (66, 'Asset & Facility', 'Asset Maintenance', '/asset/maintenance', 'Tool', 'Kelola maintenance aset', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (64, 'Asset & Facility', 'Asset Management', '/asset/management', 'Package', 'Kelola aset perusahaan', 98, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:38:30.462879');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (70, 'Disciplinary & Relations', 'Disciplinary Action', '/disciplinary/action', 'AlertTriangle', 'Kelola tindakan disipliner', 99, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:40:03.850428');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (72, 'Disciplinary & Relations', 'Grievance Management', '/disciplinary/grievance', 'MessageSquare', 'Kelola pengaduan karyawan', 99, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:40:03.850428');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (71, 'Disciplinary & Relations', 'Employee Relations', '/disciplinary/relations', 'Users', 'Kelola hubungan karyawan', 99, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:40:03.850428');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (93, 'Employee Self Service', 'Employee Portal', '/ess', 'User', 'Employee self-service portal', NULL, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-02 04:24:24.622239');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (73, 'Disciplinary & Relations', 'Investigation', '/disciplinary/investigation', 'Search', 'Kelola investigasi kasus', 99, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:40:03.850428');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (74, 'Offboarding & Exit', 'Resignation Process', '/offboarding/resignation', 'UserMinus', 'Proses pengunduran diri', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (79, 'Offboarding & Exit', 'Knowledge Transfer', '/offboarding/knowledge-transfer', 'BookOpen', 'Transfer pengetahuan', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (78, 'Offboarding & Exit', 'Alumni Network', '/offboarding/alumni', 'Users', 'Jaringan alumni karyawan', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (96, 'Payroll & Compensation', 'Payroll System', '/payroll', 'DollarSign', 'Payroll and compensation management', NULL, 'enterprise', true, '2025-12-29 13:52:18.807576', '2026-01-02 07:44:42.782534');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (56, 'Payroll & Compensation', 'Proses Payroll', '/payroll/process', 'Calculator', 'Proses penggajian bulanan', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (58, 'Payroll & Compensation', 'Pajak PPh 21', '/payroll/tax', 'FileText', 'Kelola pajak PPh 21', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (60, 'Payroll & Compensation', 'BPJS Kesehatan', '/payroll/bpjs-kes', 'Heart', 'Kelola BPJS Kesehatan', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (62, 'Payroll & Compensation', 'Potongan', '/payroll/deductions', 'Minus', 'Kelola potongan gaji', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (55, 'Payroll & Compensation', 'Komponen Gaji', '/payroll/components', 'List', 'Kelola komponen gaji', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (54, 'Payroll & Compensation', 'Struktur Gaji', '/payroll/salary-structure', 'DollarSign', 'Kelola struktur gaji', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (97, 'Performance Management', 'Performance System', '/performance', 'TrendingUp', 'Performance management system', NULL, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-02 07:45:54.800191');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (82, 'Reporting & Analytics', 'Attendance Reports', '/reporting/attendance', 'Clock', 'Laporan absensi', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (80, 'Reporting & Analytics', 'HR Dashboard', '/reporting/dashboard', 'BarChart3', 'Dashboard HR analytics', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (81, 'Reporting & Analytics', 'Employee Reports', '/reporting/employee', 'Users', 'Laporan data karyawan', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (84, 'Reporting & Analytics', 'Performance Reports', '/reporting/performance', 'TrendingUp', 'Laporan performa', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (85, 'Reporting & Analytics', 'Custom Reports', '/reporting/custom', 'FileText', 'Laporan kustom', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (83, 'Reporting & Analytics', 'Payroll Reports', '/reporting/payroll', 'DollarSign', 'Laporan payroll', 130, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:47:35.35637');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (89, 'System & Security', 'System Settings', '/system/settings', 'Settings', 'Pengaturan sistem', 131, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:48:31.329303');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (86, 'System & Security', 'User Management', '/system/users', 'Users', 'Kelola pengguna sistem', 131, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:48:31.329303');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (87, 'System & Security', 'Role Management', '/system/roles', 'Shield', 'Kelola role dan permission', 131, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:48:31.329303');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (33, 'Attendance & Time', 'Shift Management', '/attendance/shift', 'RotateCcw', 'Kelola jadwal shift karyawan', 95, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:13:24.287804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (30, 'Attendance & Time', 'Absensi Harian', '/attendance/presence', 'Clock', 'Kelola presensi harian karyawan', 95, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:13:24.287804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (24, 'Recruitment', 'Job Requisition', '/recruitment/job_equisition', 'Target', 'Perencanaan kebutuhan tenaga kerja', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:07:09.214931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (26, 'Recruitment', 'Candidate Management', '/recruitment/candidate-management', 'Users', 'Sistem pelacakan pelamar kerja', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:07:43.344185');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (28, 'Recruitment', 'Offer & Hiring', '/recruitment/offer-hiring', 'UserCheck', 'Kelola penawaran kerja dan proses hiring', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:08:45.23193');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (29, 'Recruitment', 'Pre-Employment Onboarding', '/recruitment/onboarding', 'UserPlus', 'Proses orientasi karyawan baru', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:09:01.228993');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (111, 'Employee Self Service', 'Permission', '/ess/requests/permission', 'Clock', 'Request permission', 93, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (16, 'Master Data', 'Certifications', '/core-hr/master_data/certifications', 'Award', 'Kelola sertifikasi karyawan', 92, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (110, 'Employee Self Service', 'Sick Leave', '/ess/requests/sick', 'Heart', 'Request sick leave', 93, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (137, 'Module Management', 'Module Management', '/module', 'Test', 'Manage modules on huminor system', NULL, 'enterprise', true, '2026-01-17 13:14:55.929865', '2026-01-17 13:14:55.929865');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (23, 'Employee Self Service', 'Feedback & Saran', '/ess/feedback', 'MessageSquare', 'Berikan feedback dan saran', 93, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (20, 'Employee Self Service', 'Slip Gaji', '/ess/payslip', 'Receipt', 'Lihat slip gaji', 93, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (18, 'Employee Self Service', 'Update Profil', '/ess/profile', 'User', 'Update profil pribadi', 93, 'basic', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (21, 'Employee Self Service', 'Dokumen Pribadi', '/ess/documents', 'Folder', 'Kelola dokumen pribadi', 93, 'basic', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (22, 'Employee Self Service', 'Pengumuman', '/ess/announcements', 'Megaphone', 'Lihat pengumuman perusahaan', 93, 'basic', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (25, 'Recruitment', 'Job Vacancy Management', '/recruitment/vacancy', 'Briefcase', 'Kelola lowongan pekerjaan', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (27, 'Recruitment', 'Interview & Assessment', '/recruitment/interview', 'MessageSquare', 'Kelola proses wawancara dan penilaian', 94, 'pro', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (138, 'Module Management', 'Create Modules', '/module/create-module', 'TestChild', 'Create new module', 137, 'enterprise', true, '2026-01-17 13:16:35.46013', '2026-01-17 06:19:21.385828');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (31, 'Attendance & Time', 'Lembur', '/attendance/overtime', 'Clock', 'Kelola lembur karyawan', 95, 'pro', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (34, 'Attendance & Time', 'Rekap Absensi', '/attendance/recap', 'BarChart3', 'Rekap dan laporan absensi', 95, 'pro', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (35, 'Attendance & Time', 'Kalender Libur Nasional', '/attendance/holiday', 'Calendar', 'Kelola kalender libur nasional', 95, 'pro', true, '2025-12-28 17:31:43.793239', '2025-12-29 13:52:18.807576');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (106, 'Employee Self Service', 'Basic Profile', '/ess/profile/basic', 'User', 'Update basic profile information', 93, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-02 07:40:48.327533');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (107, 'Employee Self Service', 'Emergency Contacts', '/ess/profile/emergency', 'Phone', 'Manage emergency contacts', 93, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-02 07:40:48.327533');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (108, 'Employee Self Service', 'Bank Information', '/ess/profile/bank', 'CreditCard', 'Update bank information', 93, 'enterprise', true, '2025-12-29 13:52:18.807576', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (19, 'Employee Self Service', 'Pengajuan Cuti & Izin', '/ess/requests', 'FileText', 'Ajukan cuti dan izin', 93, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (109, 'Employee Self Service', 'Annual Leave', '/ess/requests/annual', 'Calendar', 'Request annual leave', 93, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-02 12:20:23.677931');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (77, 'Offboarding & Exit', 'Final Settlement', '/offboarding/settlement', 'Calculator', 'Perhitungan settlement akhir', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (75, 'Offboarding & Exit', 'Exit Interview', '/offboarding/exit-interview', 'MessageSquare', 'Wawancara keluar karyawan', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (76, 'Offboarding & Exit', 'Asset Return', '/offboarding/asset-return', 'RotateCcw', 'Pengembalian aset perusahaan', 129, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:44:15.30284');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (32, 'Attendance & Time', 'Keterlambatan', '/attendance/late', 'AlertCircle', 'Laporan keterlambatan karyawan', 95, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-02 06:15:29.428341');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (57, 'Payroll & Compensation', 'Slip Gaji', '/payroll/payslip', 'Receipt', 'Generate slip gaji', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (63, 'Payroll & Compensation', 'Bonus & Insentif', '/payroll/bonus', 'Gift', 'Kelola bonus dan insentif', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (61, 'Payroll & Compensation', 'Tunjangan', '/payroll/allowances', 'Plus', 'Kelola tunjangan karyawan', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (134, 'Master Data', 'Bank Account', '/core-hr/employees_management/bank-account', 'Users', 'Employment bank account management', 92, 'pro', true, '2026-01-08 15:38:06', '2026-01-08 15:38:08');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (59, 'Payroll & Compensation', 'BPJS Ketenagakerjaan', '/payroll/bpjs-tk', 'Shield', 'Kelola BPJS Ketenagakerjaan', 96, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-02 07:45:12.909804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (9, 'Master Data', 'Holiday & Calendar', '/core-hr/company_settings/holidays', 'Calendar', 'Kelola hari libur dan cuti bersama', 92, 'enterprise', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:38:59.447395');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (94, 'Recruitment', 'Recruitment System', '/recruitment', 'UserPlus', 'Complete recruitment management', NULL, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-02 04:24:24.622239');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (124, 'Attendance & Time', 'Jadwal Kerja', '/attendance/jadwal-kerja', 'Calendar', 'Penugasan jadwal kerja karyawan', 95, 'basic', true, '2026-01-02 04:31:54.848414', '2026-01-02 04:31:54.848414');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (125, 'Attendance & Time', 'Attendance Log', '/attendance/attendance-log', 'Calendar', 'Laporan log kehadiran', 95, 'basic', true, '2026-01-02 06:13:33.274328', '2026-01-02 06:13:33.274328');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (92, 'Master Data', 'Core HR Management', '/core-hr', 'Building', 'Core HR management system', NULL, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-02 11:48:33.651133');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (5, 'Master Data', 'Employment Status', '/core-hr/master_data/employment-status', 'UserCheck', 'Kelola status kepegawaian', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (15, 'Master Data', 'Skills', '/core-hr/master_data/skills', 'Award', 'Kelola keterampilan karyawan', 92, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (126, 'Attendance & Time', 'Koreksi Absensi', '/attendance/correction', 'Calendar', 'Koreksi data kehadiran', 95, 'pro', true, '2026-01-02 06:14:42.954963', '2026-01-02 12:13:24.287804');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (114, 'Dashboard & Analytic', 'Statistik Kehadiran', '/dashboard/statistik-kehadiran', 'Clock', 'Statistik Kehadiran (harian / bulanan)', 112, 'pro', true, '2026-01-02 03:28:20.245722', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (113, 'Dashboard & Analytic', 'Ringkasan Headcount', '/dashboard/headcount', 'TestChild', 'Ringkasan Headcount', 112, 'pro', true, '2025-12-30 13:47:14.120922', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (119, 'Dashboard & Analytic', 'Status Cuti & Izin', '/dashboard/cuti-izin', 'Clock', 'Status Cuti & Izin', 112, 'pro', true, '2026-01-02 03:31:07.803302', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (120, 'Dashboard & Analytic', 'Reminder', '/dashboard/reminder', 'Clock', 'Reminder (kontrak, cuti, appraisal)', 112, 'pro', true, '2026-01-02 03:31:07.803302', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (121, 'Dashboard & Analytic', 'Grafik HR', '/dashboard/grafik-hr', 'Clock', 'Grafik HR (turnover, absenteeism)', 112, 'pro', true, '2026-01-02 03:31:07.803302', '2026-01-02 12:14:32.541582');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (130, 'Reporting & Analytics', 'Reporting & Analytics', '/reporting', 'MessageSquare', 'Reporting & Analytics', NULL, 'enterprise', true, '2026-01-02 07:47:21.735744', '2026-01-07 04:05:44.976561');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (131, 'System & Security', 'System & Security', '/system', 'Code', 'Manage System & Security', NULL, 'basic', true, '2026-01-02 07:48:14.08102', '2026-01-07 04:05:52.264626');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (105, 'Master Data', 'Position Hierarchy', '/core-hr/master_data/positions-hierarchy', 'TrendingUp', 'Manage position hierarchy', 92, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (8, 'Master Data', 'Shift Kerja', '/core-hr/master_data/shifts', 'Clock', 'Kelola shift kerja', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (103, 'Master Data', 'Organizational Structure
', '/core-hr/master_data/company', 'Building', 'Manage company structure', 92, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:24:30.464264');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (1, 'Master Data', 'Employee Data', '/core-hr/employees_management/data', 'Users', 'Kelola data karyawan', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:32:46.667172');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (101, 'Master Data', 'Employment Details', '/core-hr/employees_management/employment', 'Briefcase', 'Manage employment details', 92, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:33:47.233049');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (100, 'Master Data', 'Personal Information', '/core-hr/employees_management/personal', 'User', 'Manage personal information', 92, 'basic', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:33:47.233049');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (102, 'Master Data', 'Career History', '/core-hr/employees_management/career', 'Phone', 'Manage employee career', 92, 'pro', true, '2025-12-29 13:52:18.807576', '2026-01-08 08:36:13.955167');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (12, 'Master Data', 'Education History', '/core-hr/employees_management/education', 'GraduationCap', 'Kelola riwayat pendidikan', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:36:13.955167');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (127, 'Master Data', 'Employment Contract', '/core-hr/employees_management/contract', 'User', 'Kelola Kontrak Kerja', 92, 'pro', true, '2026-01-02 07:28:51.690792', '2026-01-08 08:36:13.955167');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (13, 'Master Data', 'Work History', '/core-hr/employees_management/work-history', 'Briefcase', 'Kelola riwayat pekerjaan sebelumnya', 92, 'pro', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:36:13.955167');
INSERT INTO public.modules (id, category, name, url, icon, description, parent_id, subscription_tier, is_active, created_at, updated_at) VALUES (10, 'Master Data', 'Employment Document', '/core-hr/employees_management/documents', 'FileText', 'Kelola dokumen karyawan', 92, 'basic', true, '2025-12-28 17:31:43.793239', '2026-01-08 08:36:13.955167');


--
-- Data for Name: subscription_plans; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (1, 'basic', 'Basic Plan', 'Essential HR features for small businesses', 99000.00, 990000.00, 25, 3, '{"reports": "basic", "storage": "5GB", "support": "email"}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');
INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (2, 'pro', 'Professional Plan', 'Advanced HR features for growing companies', 299000.00, 2990000.00, 100, 10, '{"reports": "advanced", "storage": "50GB", "support": "priority", "api_access": true}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');
INSERT INTO public.subscription_plans (id, name, display_name, description, price_monthly, price_yearly, max_users, max_branches, features, is_active, created_at, updated_at) VALUES (3, 'enterprise', 'Enterprise Plan', 'Complete HR solution for large organizations', 599000.00, 5990000.00, NULL, NULL, '{"reports": "custom", "storage": "unlimited", "support": "dedicated", "api_access": true, "white_label": true}', true, '2025-12-28 17:31:43.800256', '2025-12-28 17:31:43.800256');


--
-- Data for Name: plan_modules; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (1, 1, 1, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (3, 1, 3, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (4, 1, 4, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (5, 1, 5, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (7, 1, 7, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (8, 1, 8, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (29, 1, 38, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (10, 1, 10, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (343, 3, 87, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (12, 1, 12, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (182, 1, 39, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (14, 1, 14, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (183, 1, 86, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (184, 1, 87, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (185, 1, 89, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (18, 1, 18, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (186, 1, 92, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (187, 1, 93, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (21, 1, 21, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (22, 1, 22, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (188, 1, 95, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (189, 1, 100, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (190, 1, 101, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (191, 1, 103, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (192, 1, 106, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (193, 1, 107, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (344, 3, 29, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (194, 1, 124, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (195, 1, 125, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (196, 1, 131, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (197, 2, 1, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (199, 2, 3, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (200, 2, 4, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (201, 2, 5, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (202, 2, 6, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (203, 2, 7, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (204, 2, 8, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (205, 2, 38, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (206, 2, 10, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (207, 2, 12, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (208, 2, 39, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (209, 2, 14, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (210, 2, 86, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (211, 2, 87, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (212, 2, 89, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (213, 2, 18, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (214, 2, 92, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (215, 2, 93, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (216, 2, 21, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (217, 2, 22, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (218, 2, 95, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (219, 2, 100, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (220, 2, 101, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (221, 2, 103, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (222, 2, 106, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (223, 2, 107, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (224, 2, 124, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (225, 2, 125, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (226, 2, 131, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (227, 2, 30, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (228, 2, 33, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (229, 2, 36, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (230, 2, 37, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (345, 3, 4, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (232, 2, 13, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (233, 2, 15, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (234, 2, 16, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (235, 2, 17, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (236, 2, 19, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (237, 2, 23, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (238, 2, 24, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (239, 2, 25, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (240, 2, 26, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (241, 2, 27, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (242, 2, 28, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (243, 2, 29, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (244, 2, 31, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (245, 2, 32, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (246, 2, 34, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (247, 2, 35, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (248, 2, 40, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (249, 2, 41, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (250, 2, 42, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (251, 2, 43, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (252, 2, 44, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (253, 2, 45, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (254, 2, 46, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (255, 2, 47, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (256, 2, 48, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (257, 2, 49, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (258, 2, 50, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (259, 2, 51, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (260, 2, 52, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (261, 2, 53, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (262, 2, 88, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (263, 2, 91, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (264, 2, 94, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (265, 2, 97, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (266, 2, 102, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (267, 2, 104, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (268, 2, 105, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (269, 2, 109, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (270, 2, 110, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (271, 2, 111, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (272, 2, 112, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (273, 2, 113, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (274, 2, 114, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (275, 2, 119, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (276, 2, 120, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (277, 2, 121, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (278, 2, 126, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (279, 2, 127, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (280, 2, 132, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (346, 3, 34, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (347, 3, 51, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (348, 3, 52, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (349, 3, 10, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (350, 3, 105, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (351, 3, 35, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (352, 3, 132, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (353, 3, 45, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (354, 3, 107, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (355, 3, 6, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (356, 3, 86, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (357, 3, 39, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (358, 3, 92, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (359, 3, 101, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (360, 3, 93, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (24, 1, 30, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (25, 1, 33, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (26, 1, 36, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (27, 1, 37, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (361, 3, 89, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (362, 3, 36, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (363, 3, 31, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (364, 3, 114, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (365, 3, 50, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (366, 3, 102, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (367, 3, 97, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (368, 3, 14, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (369, 3, 112, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (370, 3, 109, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (371, 3, 22, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (372, 3, 13, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (373, 3, 111, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (374, 3, 127, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (376, 3, 16, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (377, 3, 124, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (378, 3, 126, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (379, 3, 44, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (380, 3, 103, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (381, 3, 42, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (382, 3, 121, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (383, 3, 88, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (384, 3, 41, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (385, 3, 119, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (386, 3, 113, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (387, 3, 46, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (388, 3, 40, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (389, 3, 125, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (390, 3, 43, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (391, 3, 53, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (392, 3, 32, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (393, 3, 7, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (394, 3, 120, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (395, 3, 100, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (396, 3, 38, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (397, 3, 15, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (398, 3, 48, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (399, 3, 26, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (400, 3, 12, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (401, 3, 95, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (402, 3, 24, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (403, 3, 19, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (404, 3, 25, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (405, 3, 94, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (406, 3, 30, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (407, 3, 21, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (408, 3, 49, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (409, 3, 47, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (410, 3, 131, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (411, 3, 3, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (412, 3, 17, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (413, 3, 28, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (414, 3, 37, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (415, 3, 33, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (416, 3, 1, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (417, 3, 104, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (418, 3, 106, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (419, 3, 5, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (420, 3, 18, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (421, 3, 110, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (422, 3, 27, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (423, 3, 23, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (424, 3, 91, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (425, 3, 8, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (426, 3, 9, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (427, 3, 20, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (428, 3, 67, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (429, 3, 69, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (430, 3, 68, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (431, 3, 65, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (432, 3, 99, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (433, 3, 128, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (434, 3, 66, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (435, 3, 64, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (436, 3, 98, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (437, 3, 70, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (438, 3, 72, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (439, 3, 71, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (440, 3, 73, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (441, 3, 129, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (442, 3, 74, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (443, 3, 79, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (444, 3, 78, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (445, 3, 96, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (446, 3, 56, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (447, 3, 58, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (448, 3, 60, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (449, 3, 62, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (450, 3, 55, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (451, 3, 54, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (452, 3, 82, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (453, 3, 80, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (454, 3, 81, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (455, 3, 84, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (456, 3, 85, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (457, 3, 83, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (459, 3, 108, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (460, 3, 77, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (461, 3, 75, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (462, 3, 76, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (463, 3, 57, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (464, 3, 63, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (465, 3, 61, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (466, 3, 59, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (467, 3, 130, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (468, 3, 134, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (469, 3, 137, true);
INSERT INTO public.plan_modules (id, plan_id, module_id, is_included) VALUES (471, 3, 138, true);


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
-- Data for Name: role_modules; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (416, 2, 67, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (417, 2, 69, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (418, 2, 68, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (419, 2, 65, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (420, 2, 66, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (421, 2, 64, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (422, 2, 98, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (423, 2, 95, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (424, 2, 33, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (425, 2, 30, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (426, 2, 31, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (427, 2, 34, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (428, 2, 35, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (429, 2, 32, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (430, 2, 124, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (431, 2, 125, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (432, 2, 126, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (433, 2, 112, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (434, 2, 114, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (435, 2, 113, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (436, 2, 119, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (437, 2, 120, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (438, 2, 121, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (439, 2, 99, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (440, 2, 70, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (441, 2, 72, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (442, 2, 71, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (443, 2, 73, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (444, 2, 93, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (445, 2, 111, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (446, 2, 110, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (447, 2, 23, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (448, 2, 20, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (449, 2, 18, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (450, 2, 21, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (451, 2, 22, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (452, 2, 106, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (453, 2, 107, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (454, 2, 108, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (455, 2, 19, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (456, 2, 109, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (457, 2, 41, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (458, 2, 40, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (459, 2, 37, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (460, 2, 36, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (461, 2, 38, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (462, 2, 39, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (463, 2, 128, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (464, 2, 129, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (465, 2, 74, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (466, 2, 79, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (467, 2, 78, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (468, 2, 77, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (469, 2, 75, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (470, 2, 76, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (471, 2, 96, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (472, 2, 56, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (473, 2, 58, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (474, 2, 60, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (475, 2, 62, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (476, 2, 55, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (477, 2, 54, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (478, 2, 57, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (479, 2, 63, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (480, 2, 61, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (481, 2, 59, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (482, 2, 45, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (483, 2, 47, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (484, 2, 42, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (485, 2, 43, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (486, 2, 44, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (487, 2, 46, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (488, 2, 97, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (489, 2, 24, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (490, 2, 25, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (491, 2, 26, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (492, 2, 27, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (493, 2, 28, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (494, 2, 29, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (495, 2, 94, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (496, 2, 82, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (497, 2, 80, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (498, 2, 81, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (499, 2, 84, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (500, 2, 85, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (501, 2, 83, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (502, 2, 130, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (503, 2, 132, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (504, 2, 49, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (505, 2, 48, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (506, 2, 50, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (507, 2, 51, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (508, 2, 52, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (509, 2, 53, true, true, true, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (510, 3, 67, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (511, 3, 69, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (512, 3, 68, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (513, 3, 65, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (514, 3, 66, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (515, 3, 64, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (516, 3, 98, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (517, 3, 95, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (518, 3, 33, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (519, 3, 30, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (520, 3, 31, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (521, 3, 34, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (522, 3, 35, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (523, 3, 32, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (524, 3, 124, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (525, 3, 125, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (526, 3, 126, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (527, 3, 112, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (528, 3, 114, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (529, 3, 113, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (530, 3, 119, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (531, 3, 120, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (532, 3, 121, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (533, 3, 99, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (534, 3, 70, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (535, 3, 72, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (536, 3, 71, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (537, 3, 73, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (538, 3, 41, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (539, 3, 40, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (540, 3, 37, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (541, 3, 36, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (542, 3, 38, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (543, 3, 39, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (544, 3, 128, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (545, 3, 129, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (546, 3, 74, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (547, 3, 79, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (548, 3, 78, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (549, 3, 77, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (550, 3, 75, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (551, 3, 76, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (552, 3, 96, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (553, 3, 56, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (554, 3, 58, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (555, 3, 60, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (556, 3, 62, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (557, 3, 55, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (558, 3, 54, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (559, 3, 57, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (560, 3, 63, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (561, 3, 61, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (562, 3, 59, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (563, 3, 45, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (564, 3, 47, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (565, 3, 42, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (566, 3, 43, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (567, 3, 44, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (568, 3, 46, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (569, 3, 97, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (570, 3, 24, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (571, 3, 25, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (572, 3, 26, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (573, 3, 27, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (574, 3, 28, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (575, 3, 29, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (576, 3, 94, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (577, 3, 82, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (578, 3, 80, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (579, 3, 81, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (580, 3, 84, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (581, 3, 85, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (582, 3, 83, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (583, 3, 130, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (584, 3, 132, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (585, 3, 49, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (586, 3, 48, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (587, 3, 50, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (588, 3, 51, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (589, 3, 52, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (590, 3, 53, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (591, 4, 95, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (592, 4, 33, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (593, 4, 30, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (594, 4, 31, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (595, 4, 34, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (596, 4, 35, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (597, 4, 32, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (598, 4, 124, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (599, 4, 125, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (600, 4, 126, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (601, 4, 41, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (602, 4, 40, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (603, 4, 37, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (604, 4, 36, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (605, 4, 38, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (606, 4, 39, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (607, 4, 128, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (608, 4, 96, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (609, 4, 56, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (610, 4, 58, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (611, 4, 60, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (612, 4, 62, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (613, 4, 55, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (614, 4, 54, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (615, 4, 57, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (616, 4, 63, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (617, 4, 61, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (618, 4, 59, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (619, 4, 82, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (620, 4, 80, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (621, 4, 81, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (622, 4, 84, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (623, 4, 85, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (624, 4, 83, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (625, 4, 130, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (626, 5, 95, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (627, 5, 33, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (628, 5, 30, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (629, 5, 31, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (630, 5, 34, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (260, 1, 112, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (257, 1, 132, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (258, 1, 104, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (259, 1, 95, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (261, 1, 99, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (262, 1, 128, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (263, 1, 98, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (264, 1, 93, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (265, 1, 129, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (266, 1, 96, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (267, 1, 97, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (268, 1, 111, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (269, 1, 110, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (270, 1, 106, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (271, 1, 107, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (272, 1, 108, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (285, 1, 113, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (631, 5, 35, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (632, 5, 32, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (633, 5, 124, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (634, 5, 125, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (635, 5, 126, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (636, 5, 41, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (637, 5, 40, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (638, 5, 37, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (639, 5, 36, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (640, 5, 38, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (641, 5, 39, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (642, 5, 128, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (643, 5, 45, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (644, 5, 47, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (645, 5, 42, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (646, 5, 43, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (647, 5, 44, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (648, 5, 46, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (649, 5, 97, true, true, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (650, 5, 132, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (651, 5, 49, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (652, 5, 48, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (653, 5, 50, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (654, 5, 51, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (655, 5, 52, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (656, 5, 53, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (657, 5, 129, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (658, 5, 74, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (659, 5, 79, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (660, 5, 78, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (661, 5, 77, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (662, 5, 75, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (663, 5, 76, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (664, 5, 112, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (665, 5, 114, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (666, 5, 113, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (667, 5, 119, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (668, 5, 120, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (669, 5, 121, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (670, 8, 93, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (671, 8, 111, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (672, 8, 110, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (673, 8, 23, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (674, 8, 20, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (675, 8, 18, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (676, 8, 21, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (677, 8, 22, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (678, 8, 106, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (679, 8, 107, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (680, 8, 108, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (1, 1, 1, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (3, 1, 3, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (4, 1, 4, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (5, 1, 5, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (6, 1, 6, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (7, 1, 7, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (8, 1, 8, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (9, 1, 9, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (10, 1, 10, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (12, 1, 12, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (13, 1, 13, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (14, 1, 14, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (15, 1, 15, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (16, 1, 16, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (17, 1, 17, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (18, 1, 18, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (19, 1, 19, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (20, 1, 20, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (21, 1, 21, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (22, 1, 22, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (23, 1, 23, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (24, 1, 24, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (25, 1, 25, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (26, 1, 26, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (27, 1, 27, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (28, 1, 28, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (29, 1, 29, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (30, 1, 30, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (31, 1, 31, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (32, 1, 32, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (33, 1, 33, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (34, 1, 34, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (35, 1, 35, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (36, 1, 36, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (37, 1, 37, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (38, 1, 38, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (39, 1, 39, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (40, 1, 40, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (41, 1, 41, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (42, 1, 42, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (43, 1, 43, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (44, 1, 44, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (45, 1, 45, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (46, 1, 46, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (47, 1, 47, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (48, 1, 48, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (49, 1, 49, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (50, 1, 50, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (51, 1, 51, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (52, 1, 52, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (53, 1, 53, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (54, 1, 54, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (55, 1, 55, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (56, 1, 56, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (57, 1, 57, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (58, 1, 58, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (59, 1, 59, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (60, 1, 60, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (61, 1, 61, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (62, 1, 62, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (63, 1, 63, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (64, 1, 64, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (65, 1, 65, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (66, 1, 66, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (67, 1, 67, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (68, 1, 68, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (69, 1, 69, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (70, 1, 70, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (71, 1, 71, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (72, 1, 72, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (73, 1, 73, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (74, 1, 74, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (75, 1, 75, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (76, 1, 76, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (77, 1, 77, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (78, 1, 78, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (79, 1, 79, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (80, 1, 80, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (81, 1, 81, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (82, 1, 82, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (83, 1, 83, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (84, 1, 84, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (85, 1, 85, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (86, 1, 86, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (87, 1, 87, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (88, 1, 88, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (89, 1, 89, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (91, 1, 91, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (273, 1, 109, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (274, 1, 130, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (275, 1, 131, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (276, 1, 94, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (277, 1, 124, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (278, 1, 125, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (279, 1, 92, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (280, 1, 100, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (281, 1, 101, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (282, 1, 103, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (283, 1, 126, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (289, 1, 102, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (290, 1, 127, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (291, 1, 105, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (286, 1, 119, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (284, 1, 114, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (287, 1, 120, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (288, 1, 121, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (681, 8, 19, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (682, 8, 109, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (683, 8, 95, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (684, 8, 33, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (685, 8, 30, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (686, 8, 31, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (687, 8, 34, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (688, 8, 35, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (689, 8, 32, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (690, 8, 124, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (691, 8, 125, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (692, 8, 126, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (693, 8, 41, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (694, 8, 40, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (695, 8, 37, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (696, 8, 36, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (697, 8, 38, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (698, 8, 39, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (699, 8, 128, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (700, 8, 45, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (701, 8, 47, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (702, 8, 42, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (703, 8, 43, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (704, 8, 44, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (705, 8, 46, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (706, 8, 97, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (707, 8, 132, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (708, 8, 49, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (709, 8, 48, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (710, 8, 50, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (711, 8, 51, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (712, 8, 52, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (713, 8, 53, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (714, 9, 24, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (715, 9, 25, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (716, 9, 26, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (717, 9, 27, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (718, 9, 28, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (719, 9, 29, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (720, 9, 94, true, true, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (722, 9, 82, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (723, 9, 80, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (724, 9, 81, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (725, 9, 84, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (726, 9, 85, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (727, 9, 83, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (728, 9, 130, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (729, 10, 67, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (730, 10, 69, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (731, 10, 68, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (732, 10, 65, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (733, 10, 66, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (734, 10, 64, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (735, 10, 98, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (736, 10, 76, true, false, false, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (737, 10, 82, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (738, 10, 80, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (739, 10, 81, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (740, 10, 84, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (741, 10, 85, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (742, 10, 83, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (743, 10, 130, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (744, 11, 82, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (745, 11, 80, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (746, 11, 81, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (747, 11, 84, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (748, 11, 85, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (749, 11, 83, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (750, 11, 130, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (751, 11, 96, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (752, 11, 56, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (753, 11, 58, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (754, 11, 60, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (755, 11, 62, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (756, 11, 55, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (757, 11, 54, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (758, 11, 57, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (759, 11, 63, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (760, 11, 61, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (761, 11, 59, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (762, 11, 95, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (763, 11, 33, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (764, 11, 30, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (765, 11, 31, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (766, 11, 34, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (767, 11, 35, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (768, 11, 32, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (769, 11, 124, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (770, 11, 125, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (771, 11, 126, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (772, 11, 99, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (773, 11, 70, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (774, 11, 72, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (775, 11, 71, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (776, 11, 73, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (777, 11, 88, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (778, 12, 88, true, false, false, false);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (779, 1, 134, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (781, 13, 138, true, true, true, true);
INSERT INTO public.role_modules (id, role_id, module_id, can_read, can_write, can_delete, can_approve) VALUES (780, 13, 137, true, true, true, true);


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (1, 'create_users_table', '2025-12-28 17:31:43.688095');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (2, 'create_companies_and_branches', '2025-12-28 17:31:43.740456');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (3, 'create_roles_and_modules', '2025-12-28 17:31:43.767333');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (4, 'seed_modules_data', '2025-12-28 17:31:43.793239');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (5, 'create_subscription_system', '2025-12-28 17:31:43.800256');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (6, 'seed_initial_data', '2025-12-28 17:31:43.824744');
INSERT INTO public.schema_migrations (version, name, applied_at) VALUES (7, 'add_module_hierarchy', '2025-12-29 13:52:18.807576');


--
-- Data for Name: subscriptions; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (2, 2, 3, 'active', 'yearly', '2025-12-28', '2026-12-28', 5990000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2025-12-28 17:31:43.824744');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (3, 3, 2, 'active', 'monthly', '2025-12-28', '2026-01-31', 299000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2025-12-31 15:18:49.681511');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (8, 4, 3, 'active', 'monthly', '2025-12-31', '2026-01-31', 599000.00, 'IDR', 'paid', '2025-12-31', '2026-01-31', false, '2025-12-31 15:03:49.811229', '2025-12-31 15:26:34.213118');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (6, 5, 2, 'active', 'yearly', '2025-12-31', '2026-12-31', 2990000.00, 'IDR', 'paid', '2025-12-31', '2026-12-31', true, '2025-12-31 15:01:41.777448', '2025-12-31 15:32:44.393663');
INSERT INTO public.subscriptions (id, company_id, plan_id, status, billing_cycle, start_date, end_date, price, currency, payment_status, last_payment_date, next_payment_date, auto_renew, created_at, updated_at) VALUES (1, 1, 3, 'active', 'yearly', '2025-12-28', '2026-01-31', 5990000.00, 'IDR', 'paid', '2025-12-28', '2026-12-28', true, '2025-12-28 17:31:43.824744', '2026-01-04 11:50:44.296087');


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
-- Data for Name: unit_roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (1, 5, 2, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (2, 6, 9, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (3, 7, 4, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (4, 1, 3, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (5, 2, 5, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (6, 3, 12, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (7, 10, 5, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (8, 13, 5, '2026-01-10 21:49:31.139975', '2026-01-10 21:49:31.139975');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (10, 16, 13, '2026-01-17 23:29:00.555546', '2026-01-17 23:29:00.555546');
INSERT INTO public.unit_roles (id, unit_id, role_id, created_at, updated_at) VALUES (11, 17, 14, '2026-01-18 04:16:55.626198', '2026-01-18 04:16:55.626198');


--
-- Data for Name: unit_role_modules; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (1, 1, 40, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (2, 1, 41, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (3, 1, 45, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (4, 1, 47, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (5, 1, 42, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (6, 1, 43, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (7, 1, 44, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (8, 1, 46, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (9, 1, 49, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (10, 1, 48, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (11, 1, 50, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (12, 1, 37, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (13, 1, 36, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (14, 1, 38, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (15, 1, 39, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (16, 1, 132, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (17, 1, 99, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (18, 1, 51, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (19, 1, 52, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (20, 1, 53, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (21, 1, 95, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (22, 1, 128, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (23, 1, 70, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (24, 1, 72, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (25, 1, 71, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (26, 1, 93, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (27, 1, 73, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (28, 1, 97, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (29, 1, 33, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (30, 1, 30, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (31, 1, 111, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (32, 1, 110, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (33, 1, 23, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (34, 1, 20, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (35, 1, 18, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (36, 1, 21, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (37, 1, 22, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (38, 1, 31, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (39, 1, 34, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (40, 1, 35, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (41, 1, 106, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (42, 1, 107, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (43, 1, 108, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (44, 1, 19, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (45, 1, 109, true, true, true, false, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (46, 1, 32, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (47, 1, 124, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (48, 1, 125, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (49, 1, 126, true, true, true, true, '2026-01-10 21:49:46.339099', '2026-01-10 21:49:46.339099');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (50, 2, 133, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (51, 2, 24, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (52, 2, 26, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (53, 2, 28, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (54, 2, 29, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (55, 2, 25, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (56, 2, 27, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (57, 2, 94, true, true, false, false, '2026-01-10 21:49:46.355749', '2026-01-10 21:49:46.355749');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (58, 3, 40, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (59, 3, 41, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (60, 3, 37, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (61, 3, 36, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (62, 3, 38, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (63, 3, 39, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (64, 3, 95, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (65, 3, 128, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (66, 3, 96, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (67, 3, 56, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (68, 3, 58, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (69, 3, 60, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (70, 3, 62, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (71, 3, 55, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (72, 3, 54, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (73, 3, 82, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (74, 3, 80, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (75, 3, 81, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (76, 3, 84, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (77, 3, 85, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (78, 3, 83, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (79, 3, 33, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (80, 3, 30, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (81, 3, 31, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (82, 3, 34, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (83, 3, 35, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (84, 3, 32, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (85, 3, 57, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (86, 3, 63, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (87, 3, 61, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (88, 3, 59, true, true, false, true, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (89, 3, 124, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (90, 3, 125, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (91, 3, 126, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (92, 3, 130, true, false, false, false, '2026-01-10 21:49:46.357774', '2026-01-10 21:49:46.357774');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (93, 7, 40, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (94, 7, 41, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (95, 7, 45, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (96, 7, 47, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (97, 7, 42, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (98, 7, 43, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (99, 7, 44, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (100, 7, 46, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (101, 7, 37, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (102, 7, 36, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (103, 7, 38, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (104, 7, 39, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (105, 7, 95, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (106, 7, 112, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (107, 7, 128, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (108, 7, 93, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (109, 7, 97, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (110, 7, 33, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (111, 7, 30, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (112, 7, 111, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (113, 7, 110, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (114, 7, 23, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (115, 7, 20, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (116, 7, 18, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (117, 7, 21, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (118, 7, 22, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (119, 7, 31, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (120, 7, 34, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (121, 7, 35, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (122, 7, 106, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (123, 7, 107, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (124, 7, 108, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (125, 7, 19, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (126, 7, 109, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (127, 7, 32, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (128, 7, 124, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (129, 7, 125, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (130, 7, 126, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (131, 7, 114, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (132, 7, 113, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (133, 7, 119, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (134, 7, 120, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (135, 7, 121, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (136, 8, 40, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (137, 8, 41, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (138, 8, 45, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (139, 8, 47, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (140, 8, 42, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (141, 8, 43, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (142, 8, 44, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (143, 8, 46, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (144, 8, 37, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (145, 8, 36, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (146, 8, 38, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (147, 8, 39, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (148, 8, 95, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (149, 8, 112, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (150, 8, 128, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (151, 8, 93, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (152, 8, 97, true, true, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (153, 8, 33, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (154, 8, 30, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (155, 8, 111, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (156, 8, 110, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (157, 8, 23, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (158, 8, 20, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (159, 8, 18, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (160, 8, 21, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (161, 8, 22, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (162, 8, 31, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (163, 8, 34, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (164, 8, 35, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (165, 8, 106, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (166, 8, 107, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (167, 8, 108, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (168, 8, 19, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (169, 8, 109, true, true, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (170, 8, 32, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (171, 8, 124, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (172, 8, 125, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (173, 8, 126, true, false, false, true, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (174, 8, 114, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (175, 8, 113, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (176, 8, 119, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (177, 8, 120, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (178, 8, 121, true, false, false, false, '2026-01-10 21:49:46.361861', '2026-01-10 21:49:46.361861');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (179, 6, 88, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (180, 6, 91, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (181, 6, 82, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (182, 6, 80, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (183, 6, 81, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (184, 6, 84, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (185, 6, 85, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (186, 6, 83, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (187, 6, 89, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (188, 6, 86, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (189, 6, 87, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (190, 6, 130, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (191, 6, 131, true, true, true, false, '2026-01-10 21:49:46.366444', '2026-01-10 21:49:46.366444');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (192, 10, 137, true, true, true, true, '2026-01-17 23:30:26.29535', '2026-01-17 23:30:26.29535');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (193, 10, 138, true, true, true, true, '2026-01-17 23:30:26.29535', '2026-01-17 23:30:26.29535');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (194, 11, 137, true, true, true, true, '2026-01-18 11:46:56.643139', '2026-01-18 11:46:56.643139');
INSERT INTO public.unit_role_modules (id, unit_role_id, module_id, can_read, can_write, can_delete, can_approve, created_at, updated_at) VALUES (195, 11, 138, true, true, true, true, '2026-01-18 11:46:56.643139', '2026-01-18 11:46:56.643139');


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
-- Data for Name: user_roles_backup; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.user_roles_backup (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (1, 1, 1, 1, 1, '2025-12-28 17:31:43.824744', NULL);
INSERT INTO public.user_roles_backup (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (21, 6, 8, 1, 2, '2026-01-06 11:17:37.447977', NULL);
INSERT INTO public.user_roles_backup (id, user_id, role_id, company_id, branch_id, created_at, unit_id) VALUES (22, 3, 1, 1, 1, '2026-01-06 14:49:55.185813', NULL);


--
-- Name: audit_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.audit_logs_id_seq', 6, true);


--
-- Name: branches_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.branches_id_seq', 18, true);


--
-- Name: companies_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.companies_id_seq', 6, true);


--
-- Name: modules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.modules_id_seq', 138, true);


--
-- Name: plan_modules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.plan_modules_id_seq', 471, true);


--
-- Name: role_modules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.role_modules_id_seq', 781, true);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.roles_id_seq', 14, true);


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.subscription_plans_id_seq', 3, true);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.subscriptions_id_seq', 9, true);


--
-- Name: unit_role_modules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.unit_role_modules_id_seq', 195, true);


--
-- Name: unit_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.unit_roles_id_seq', 11, true);


--
-- Name: units_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.units_id_seq', 17, true);


--
-- Name: user_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.user_roles_id_seq', 35, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_id_seq', 23, true);


--
-- PostgreSQL database dump complete
--

