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
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_id_seq', 23, true);


--
-- PostgreSQL database dump complete
--

