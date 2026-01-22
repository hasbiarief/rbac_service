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
-- Name: unit_roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.unit_roles_id_seq', 11, true);


--
-- PostgreSQL database dump complete
--

