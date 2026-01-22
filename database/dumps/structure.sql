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

ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_user_id_fkey;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_role_id_fkey;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_company_id_fkey;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_branch_id_fkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_parent_id_fkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_branch_id_fkey;
ALTER TABLE IF EXISTS ONLY public.unit_roles DROP CONSTRAINT IF EXISTS unit_roles_unit_id_fkey;
ALTER TABLE IF EXISTS ONLY public.unit_roles DROP CONSTRAINT IF EXISTS unit_roles_role_id_fkey;
ALTER TABLE IF EXISTS ONLY public.unit_role_modules DROP CONSTRAINT IF EXISTS unit_role_modules_unit_role_id_fkey;
ALTER TABLE IF EXISTS ONLY public.unit_role_modules DROP CONSTRAINT IF EXISTS unit_role_modules_module_id_fkey;
ALTER TABLE IF EXISTS ONLY public.subscriptions DROP CONSTRAINT IF EXISTS subscriptions_plan_id_fkey;
ALTER TABLE IF EXISTS ONLY public.subscriptions DROP CONSTRAINT IF EXISTS subscriptions_company_id_fkey;
ALTER TABLE IF EXISTS ONLY public.role_modules DROP CONSTRAINT IF EXISTS role_modules_role_id_fkey;
ALTER TABLE IF EXISTS ONLY public.role_modules DROP CONSTRAINT IF EXISTS role_modules_module_id_fkey;
ALTER TABLE IF EXISTS ONLY public.plan_modules DROP CONSTRAINT IF EXISTS plan_modules_plan_id_fkey;
ALTER TABLE IF EXISTS ONLY public.plan_modules DROP CONSTRAINT IF EXISTS plan_modules_module_id_fkey;
ALTER TABLE IF EXISTS ONLY public.modules DROP CONSTRAINT IF EXISTS modules_parent_id_fkey;
ALTER TABLE IF EXISTS ONLY public.branches DROP CONSTRAINT IF EXISTS branches_parent_id_fkey;
ALTER TABLE IF EXISTS ONLY public.branches DROP CONSTRAINT IF EXISTS branches_company_id_fkey;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_user_id_fkey;
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;
DROP TRIGGER IF EXISTS update_units_updated_at ON public.units;
DROP TRIGGER IF EXISTS update_unit_roles_updated_at ON public.unit_roles;
DROP TRIGGER IF EXISTS update_unit_role_modules_updated_at ON public.unit_role_modules;
DROP TRIGGER IF EXISTS update_unit_hierarchy_trigger ON public.units;
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON public.subscriptions;
DROP TRIGGER IF EXISTS update_subscription_plans_updated_at ON public.subscription_plans;
DROP TRIGGER IF EXISTS update_roles_updated_at ON public.roles;
DROP TRIGGER IF EXISTS update_modules_updated_at ON public.modules;
DROP TRIGGER IF EXISTS update_companies_updated_at ON public.companies;
DROP TRIGGER IF EXISTS update_branches_updated_at ON public.branches;
DROP TRIGGER IF EXISTS update_branch_hierarchy_trigger ON public.branches;
DROP INDEX IF EXISTS public.idx_users_user_identity;
DROP INDEX IF EXISTS public.idx_users_is_active;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_users_deleted_at;
DROP INDEX IF EXISTS public.idx_user_roles_user_id;
DROP INDEX IF EXISTS public.idx_user_roles_unit_id;
DROP INDEX IF EXISTS public.idx_user_roles_company_id;
DROP INDEX IF EXISTS public.idx_user_effective_permissions_user_id;
DROP INDEX IF EXISTS public.idx_user_effective_permissions_module;
DROP INDEX IF EXISTS public.idx_units_path;
DROP INDEX IF EXISTS public.idx_units_parent_id;
DROP INDEX IF EXISTS public.idx_units_is_active;
DROP INDEX IF EXISTS public.idx_units_code;
DROP INDEX IF EXISTS public.idx_units_branch_id;
DROP INDEX IF EXISTS public.idx_unit_roles_unit_id;
DROP INDEX IF EXISTS public.idx_unit_roles_role_id;
DROP INDEX IF EXISTS public.idx_unit_role_modules_unit_role_id;
DROP INDEX IF EXISTS public.idx_unit_role_modules_module_id;
DROP INDEX IF EXISTS public.idx_subscriptions_status;
DROP INDEX IF EXISTS public.idx_subscriptions_end_date;
DROP INDEX IF EXISTS public.idx_subscriptions_company_id;
DROP INDEX IF EXISTS public.idx_subscription_plans_name;
DROP INDEX IF EXISTS public.idx_subscription_plans_is_active;
DROP INDEX IF EXISTS public.idx_roles_name;
DROP INDEX IF EXISTS public.idx_roles_is_active;
DROP INDEX IF EXISTS public.idx_role_modules_role_id;
DROP INDEX IF EXISTS public.idx_role_modules_module_id;
DROP INDEX IF EXISTS public.idx_role_modules_can_approve;
DROP INDEX IF EXISTS public.idx_plan_modules_plan_id;
DROP INDEX IF EXISTS public.idx_plan_modules_module_id;
DROP INDEX IF EXISTS public.idx_modules_url;
DROP INDEX IF EXISTS public.idx_modules_subscription_tier;
DROP INDEX IF EXISTS public.idx_modules_parent_id;
DROP INDEX IF EXISTS public.idx_modules_is_active;
DROP INDEX IF EXISTS public.idx_modules_category_parent;
DROP INDEX IF EXISTS public.idx_modules_category;
DROP INDEX IF EXISTS public.idx_companies_is_active;
DROP INDEX IF EXISTS public.idx_companies_code;
DROP INDEX IF EXISTS public.idx_branches_path;
DROP INDEX IF EXISTS public.idx_branches_parent_id;
DROP INDEX IF EXISTS public.idx_branches_is_active;
DROP INDEX IF EXISTS public.idx_branches_company_id;
DROP INDEX IF EXISTS public.idx_audit_logs_user_id;
DROP INDEX IF EXISTS public.idx_audit_logs_resource;
DROP INDEX IF EXISTS public.idx_audit_logs_created_at;
DROP INDEX IF EXISTS public.idx_audit_logs_action;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_user_identity_key;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE IF EXISTS ONLY public.users DROP CONSTRAINT IF EXISTS users_email_key;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_user_id_role_id_company_id_branch_id_key;
ALTER TABLE IF EXISTS ONLY public.user_roles DROP CONSTRAINT IF EXISTS user_roles_pkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_pkey;
ALTER TABLE IF EXISTS ONLY public.units DROP CONSTRAINT IF EXISTS units_branch_id_code_key;
ALTER TABLE IF EXISTS ONLY public.unit_roles DROP CONSTRAINT IF EXISTS unit_roles_unit_id_role_id_key;
ALTER TABLE IF EXISTS ONLY public.unit_roles DROP CONSTRAINT IF EXISTS unit_roles_pkey;
ALTER TABLE IF EXISTS ONLY public.unit_role_modules DROP CONSTRAINT IF EXISTS unit_role_modules_unit_role_id_module_id_key;
ALTER TABLE IF EXISTS ONLY public.unit_role_modules DROP CONSTRAINT IF EXISTS unit_role_modules_pkey;
ALTER TABLE IF EXISTS ONLY public.subscriptions DROP CONSTRAINT IF EXISTS subscriptions_pkey;
ALTER TABLE IF EXISTS ONLY public.subscriptions DROP CONSTRAINT IF EXISTS subscriptions_company_id_key;
ALTER TABLE IF EXISTS ONLY public.subscription_plans DROP CONSTRAINT IF EXISTS subscription_plans_pkey;
ALTER TABLE IF EXISTS ONLY public.subscription_plans DROP CONSTRAINT IF EXISTS subscription_plans_name_key;
ALTER TABLE IF EXISTS ONLY public.schema_migrations DROP CONSTRAINT IF EXISTS schema_migrations_pkey;
ALTER TABLE IF EXISTS ONLY public.roles DROP CONSTRAINT IF EXISTS roles_pkey;
ALTER TABLE IF EXISTS ONLY public.roles DROP CONSTRAINT IF EXISTS roles_name_key;
ALTER TABLE IF EXISTS ONLY public.role_modules DROP CONSTRAINT IF EXISTS role_modules_role_id_module_id_key;
ALTER TABLE IF EXISTS ONLY public.role_modules DROP CONSTRAINT IF EXISTS role_modules_pkey;
ALTER TABLE IF EXISTS ONLY public.plan_modules DROP CONSTRAINT IF EXISTS plan_modules_plan_id_module_id_key;
ALTER TABLE IF EXISTS ONLY public.plan_modules DROP CONSTRAINT IF EXISTS plan_modules_pkey;
ALTER TABLE IF EXISTS ONLY public.modules DROP CONSTRAINT IF EXISTS modules_url_key;
ALTER TABLE IF EXISTS ONLY public.modules DROP CONSTRAINT IF EXISTS modules_pkey;
ALTER TABLE IF EXISTS ONLY public.companies DROP CONSTRAINT IF EXISTS companies_pkey;
ALTER TABLE IF EXISTS ONLY public.companies DROP CONSTRAINT IF EXISTS companies_code_key;
ALTER TABLE IF EXISTS ONLY public.branches DROP CONSTRAINT IF EXISTS branches_pkey;
ALTER TABLE IF EXISTS ONLY public.branches DROP CONSTRAINT IF EXISTS branches_company_id_code_key;
ALTER TABLE IF EXISTS ONLY public.audit_logs DROP CONSTRAINT IF EXISTS audit_logs_pkey;
ALTER TABLE IF EXISTS public.users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.user_roles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.units ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.unit_roles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.unit_role_modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.subscriptions ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.subscription_plans ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.roles ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.role_modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.plan_modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.modules ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.companies ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.branches ALTER COLUMN id DROP DEFAULT;
ALTER TABLE IF EXISTS public.audit_logs ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.users;
DROP SEQUENCE IF EXISTS public.user_roles_id_seq;
DROP TABLE IF EXISTS public.user_roles_backup;
DROP VIEW IF EXISTS public.user_effective_permissions;
DROP TABLE IF EXISTS public.user_roles;
DROP SEQUENCE IF EXISTS public.units_id_seq;
DROP TABLE IF EXISTS public.units;
DROP SEQUENCE IF EXISTS public.unit_roles_id_seq;
DROP TABLE IF EXISTS public.unit_roles;
DROP SEQUENCE IF EXISTS public.unit_role_modules_id_seq;
DROP TABLE IF EXISTS public.unit_role_modules;
DROP SEQUENCE IF EXISTS public.subscriptions_id_seq;
DROP TABLE IF EXISTS public.subscriptions;
DROP SEQUENCE IF EXISTS public.subscription_plans_id_seq;
DROP TABLE IF EXISTS public.subscription_plans;
DROP TABLE IF EXISTS public.schema_migrations;
DROP SEQUENCE IF EXISTS public.roles_id_seq;
DROP TABLE IF EXISTS public.roles;
DROP SEQUENCE IF EXISTS public.role_modules_id_seq;
DROP TABLE IF EXISTS public.role_modules;
DROP SEQUENCE IF EXISTS public.plan_modules_id_seq;
DROP TABLE IF EXISTS public.plan_modules;
DROP SEQUENCE IF EXISTS public.modules_id_seq;
DROP TABLE IF EXISTS public.modules;
DROP SEQUENCE IF EXISTS public.companies_id_seq;
DROP TABLE IF EXISTS public.companies;
DROP SEQUENCE IF EXISTS public.branches_id_seq;
DROP TABLE IF EXISTS public.branches;
DROP SEQUENCE IF EXISTS public.audit_logs_id_seq;
DROP TABLE IF EXISTS public.audit_logs;
DROP FUNCTION IF EXISTS public.update_updated_at_column();
DROP FUNCTION IF EXISTS public.update_unit_hierarchy();
DROP FUNCTION IF EXISTS public.update_branch_hierarchy();
--
-- Name: update_branch_hierarchy(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_branch_hierarchy() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level = 0;
        NEW.path = NEW.id::TEXT;
    ELSE
        SELECT level + 1, path || '.' || NEW.id::TEXT
        INTO NEW.level, NEW.path
        FROM branches
        WHERE id = NEW.parent_id;
    END IF;
    RETURN NEW;
END;
$$;


--
-- Name: update_unit_hierarchy(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_unit_hierarchy() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.level = 0;
        NEW.path = NEW.id::TEXT;
    ELSE
        SELECT level + 1, path || '.' || NEW.id::TEXT
        INTO NEW.level, NEW.path
        FROM units
        WHERE id = NEW.parent_id;
    END IF;
    RETURN NEW;
END;
$$;


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_logs (
    id bigint NOT NULL,
    user_id bigint,
    action character varying(100) NOT NULL,
    resource character varying(100) NOT NULL,
    resource_id bigint,
    details jsonb DEFAULT '{}'::jsonb,
    ip_address inet,
    user_agent text,
    success boolean DEFAULT true,
    error_message text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.audit_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: audit_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.audit_logs_id_seq OWNED BY public.audit_logs.id;


--
-- Name: branches; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.branches (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    parent_id bigint,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    level integer DEFAULT 0,
    path text DEFAULT ''::text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: branches_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.branches_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: branches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.branches_id_seq OWNED BY public.branches.id;


--
-- Name: companies; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.companies (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: companies_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.companies_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: companies_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.companies_id_seq OWNED BY public.companies.id;


--
-- Name: modules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.modules (
    id bigint NOT NULL,
    category character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    icon character varying(100),
    description text,
    parent_id bigint,
    subscription_tier character varying(20) DEFAULT 'basic'::character varying,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT modules_subscription_tier_check CHECK (((subscription_tier)::text = ANY ((ARRAY['basic'::character varying, 'pro'::character varying, 'enterprise'::character varying])::text[])))
);


--
-- Name: COLUMN modules.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.modules.parent_id IS 'Reference to parent module for hierarchical structure';


--
-- Name: modules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.modules_id_seq OWNED BY public.modules.id;


--
-- Name: plan_modules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.plan_modules (
    id bigint NOT NULL,
    plan_id bigint NOT NULL,
    module_id bigint NOT NULL,
    is_included boolean DEFAULT true
);


--
-- Name: plan_modules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.plan_modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: plan_modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.plan_modules_id_seq OWNED BY public.plan_modules.id;


--
-- Name: role_modules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.role_modules (
    id bigint NOT NULL,
    role_id bigint NOT NULL,
    module_id bigint NOT NULL,
    can_read boolean DEFAULT false,
    can_write boolean DEFAULT false,
    can_delete boolean DEFAULT false,
    can_approve boolean DEFAULT false
);


--
-- Name: COLUMN role_modules.can_approve; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.role_modules.can_approve IS 'Permission to approve requests/transactions in this module';


--
-- Name: role_modules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.role_modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: role_modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.role_modules_id_seq OWNED BY public.role_modules.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version integer NOT NULL,
    name character varying(255) NOT NULL,
    applied_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: subscription_plans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscription_plans (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    display_name character varying(100) NOT NULL,
    description text,
    price_monthly numeric(10,2) NOT NULL,
    price_yearly numeric(10,2) NOT NULL,
    max_users integer,
    max_branches integer,
    features jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.subscription_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.subscription_plans_id_seq OWNED BY public.subscription_plans.id;


--
-- Name: subscriptions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subscriptions (
    id bigint NOT NULL,
    company_id bigint NOT NULL,
    plan_id bigint NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying,
    billing_cycle character varying(10) DEFAULT 'monthly'::character varying,
    start_date date DEFAULT CURRENT_DATE NOT NULL,
    end_date date NOT NULL,
    price numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'IDR'::character varying,
    payment_status character varying(20) DEFAULT 'pending'::character varying,
    last_payment_date date,
    next_payment_date date,
    auto_renew boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT subscriptions_billing_cycle_check CHECK (((billing_cycle)::text = ANY ((ARRAY['monthly'::character varying, 'yearly'::character varying])::text[]))),
    CONSTRAINT subscriptions_payment_status_check CHECK (((payment_status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'failed'::character varying, 'refunded'::character varying])::text[]))),
    CONSTRAINT subscriptions_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'expired'::character varying, 'cancelled'::character varying, 'suspended'::character varying, 'trial'::character varying])::text[])))
);


--
-- Name: subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.subscriptions_id_seq OWNED BY public.subscriptions.id;


--
-- Name: unit_role_modules; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.unit_role_modules (
    id bigint NOT NULL,
    unit_role_id bigint NOT NULL,
    module_id bigint NOT NULL,
    can_read boolean DEFAULT false,
    can_write boolean DEFAULT false,
    can_delete boolean DEFAULT false,
    can_approve boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: TABLE unit_role_modules; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.unit_role_modules IS 'Unit-specific role permissions that can override default role permissions for granular access control';


--
-- Name: unit_role_modules_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.unit_role_modules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: unit_role_modules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.unit_role_modules_id_seq OWNED BY public.unit_role_modules.id;


--
-- Name: unit_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.unit_roles (
    id bigint NOT NULL,
    unit_id bigint NOT NULL,
    role_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: TABLE unit_roles; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.unit_roles IS 'Maps roles to specific units for granular access control';


--
-- Name: unit_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.unit_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: unit_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.unit_roles_id_seq OWNED BY public.unit_roles.id;


--
-- Name: units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.units (
    id bigint NOT NULL,
    branch_id bigint NOT NULL,
    parent_id bigint,
    name character varying(255) NOT NULL,
    code character varying(50) NOT NULL,
    description text,
    level integer DEFAULT 0,
    path text DEFAULT ''::text,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: TABLE units; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.units IS 'Units within branches - represents departments, divisions, or teams';


--
-- Name: COLUMN units.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.units.parent_id IS 'Reference to parent unit for hierarchical structure within branch';


--
-- Name: COLUMN units.level; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.units.level IS 'Depth level in unit hierarchy (0 = root unit)';


--
-- Name: COLUMN units.path; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.units.path IS 'Hierarchical path for efficient tree queries';


--
-- Name: units_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.units_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.units_id_seq OWNED BY public.units.id;


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    role_id bigint NOT NULL,
    company_id bigint NOT NULL,
    branch_id bigint,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    unit_id bigint
);


--
-- Name: COLUMN user_roles.unit_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_roles.unit_id IS 'Optional unit assignment for user role - if NULL, role applies to entire branch';


--
-- Name: user_effective_permissions; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.user_effective_permissions AS
 SELECT DISTINCT ur.user_id,
    ur.company_id,
    ur.branch_id,
    ur.unit_id,
    m.id AS module_id,
    m.name AS module_name,
    m.category AS module_category,
    m.url AS module_url,
    COALESCE(urm.can_read, rm.can_read, false) AS can_read,
    COALESCE(urm.can_write, rm.can_write, false) AS can_write,
    COALESCE(urm.can_delete, rm.can_delete, false) AS can_delete,
    COALESCE(urm.can_approve, rm.can_approve, false) AS can_approve,
        CASE
            WHEN (urm.id IS NOT NULL) THEN true
            ELSE false
        END AS is_customized
   FROM (((((public.user_roles ur
     JOIN public.roles r ON ((ur.role_id = r.id)))
     JOIN public.role_modules rm ON ((r.id = rm.role_id)))
     JOIN public.modules m ON ((rm.module_id = m.id)))
     LEFT JOIN public.unit_roles unit_r ON (((ur.unit_id = unit_r.unit_id) AND (ur.role_id = unit_r.role_id))))
     LEFT JOIN public.unit_role_modules urm ON (((unit_r.id = urm.unit_role_id) AND (m.id = urm.module_id))))
  WHERE ((ur.user_id IS NOT NULL) AND (r.is_active = true) AND (m.is_active = true));


--
-- Name: VIEW user_effective_permissions; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON VIEW public.user_effective_permissions IS 'Consolidated view of user permissions combining default role permissions with unit-specific overrides';


--
-- Name: user_roles_backup; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_roles_backup (
    id bigint,
    user_id bigint,
    role_id bigint,
    company_id bigint,
    branch_id bigint,
    created_at timestamp without time zone,
    unit_id bigint
);


--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_roles_id_seq OWNED BY public.user_roles.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    user_identity character varying(9),
    password_hash character varying(255) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: audit_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs ALTER COLUMN id SET DEFAULT nextval('public.audit_logs_id_seq'::regclass);


--
-- Name: branches id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.branches ALTER COLUMN id SET DEFAULT nextval('public.branches_id_seq'::regclass);


--
-- Name: companies id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies ALTER COLUMN id SET DEFAULT nextval('public.companies_id_seq'::regclass);


--
-- Name: modules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.modules ALTER COLUMN id SET DEFAULT nextval('public.modules_id_seq'::regclass);


--
-- Name: plan_modules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plan_modules ALTER COLUMN id SET DEFAULT nextval('public.plan_modules_id_seq'::regclass);


--
-- Name: role_modules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_modules ALTER COLUMN id SET DEFAULT nextval('public.role_modules_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: subscription_plans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans ALTER COLUMN id SET DEFAULT nextval('public.subscription_plans_id_seq'::regclass);


--
-- Name: subscriptions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions ALTER COLUMN id SET DEFAULT nextval('public.subscriptions_id_seq'::regclass);


--
-- Name: unit_role_modules id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_role_modules ALTER COLUMN id SET DEFAULT nextval('public.unit_role_modules_id_seq'::regclass);


--
-- Name: unit_roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_roles ALTER COLUMN id SET DEFAULT nextval('public.unit_roles_id_seq'::regclass);


--
-- Name: units id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units ALTER COLUMN id SET DEFAULT nextval('public.units_id_seq'::regclass);


--
-- Name: user_roles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles ALTER COLUMN id SET DEFAULT nextval('public.user_roles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: branches branches_company_id_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_company_id_code_key UNIQUE (company_id, code);


--
-- Name: branches branches_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_pkey PRIMARY KEY (id);


--
-- Name: companies companies_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_code_key UNIQUE (code);


--
-- Name: companies companies_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.companies
    ADD CONSTRAINT companies_pkey PRIMARY KEY (id);


--
-- Name: modules modules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT modules_pkey PRIMARY KEY (id);


--
-- Name: modules modules_url_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT modules_url_key UNIQUE (url);


--
-- Name: plan_modules plan_modules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plan_modules
    ADD CONSTRAINT plan_modules_pkey PRIMARY KEY (id);


--
-- Name: plan_modules plan_modules_plan_id_module_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plan_modules
    ADD CONSTRAINT plan_modules_plan_id_module_id_key UNIQUE (plan_id, module_id);


--
-- Name: role_modules role_modules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_modules
    ADD CONSTRAINT role_modules_pkey PRIMARY KEY (id);


--
-- Name: role_modules role_modules_role_id_module_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_modules
    ADD CONSTRAINT role_modules_role_id_module_id_key UNIQUE (role_id, module_id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: subscription_plans subscription_plans_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_name_key UNIQUE (name);


--
-- Name: subscription_plans subscription_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);


--
-- Name: subscriptions subscriptions_company_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_company_id_key UNIQUE (company_id);


--
-- Name: subscriptions subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_pkey PRIMARY KEY (id);


--
-- Name: unit_role_modules unit_role_modules_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_role_modules
    ADD CONSTRAINT unit_role_modules_pkey PRIMARY KEY (id);


--
-- Name: unit_role_modules unit_role_modules_unit_role_id_module_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_role_modules
    ADD CONSTRAINT unit_role_modules_unit_role_id_module_id_key UNIQUE (unit_role_id, module_id);


--
-- Name: unit_roles unit_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_roles
    ADD CONSTRAINT unit_roles_pkey PRIMARY KEY (id);


--
-- Name: unit_roles unit_roles_unit_id_role_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_roles
    ADD CONSTRAINT unit_roles_unit_id_role_id_key UNIQUE (unit_id, role_id);


--
-- Name: units units_branch_id_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_branch_id_code_key UNIQUE (branch_id, code);


--
-- Name: units units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);


--
-- Name: user_roles user_roles_user_id_role_id_company_id_branch_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_role_id_company_id_branch_id_key UNIQUE (user_id, role_id, company_id, branch_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_user_identity_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_user_identity_key UNIQUE (user_identity);


--
-- Name: idx_audit_logs_action; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_action ON public.audit_logs USING btree (action);


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_created_at ON public.audit_logs USING btree (created_at);


--
-- Name: idx_audit_logs_resource; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_resource ON public.audit_logs USING btree (resource);


--
-- Name: idx_audit_logs_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_user_id ON public.audit_logs USING btree (user_id);


--
-- Name: idx_branches_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_branches_company_id ON public.branches USING btree (company_id);


--
-- Name: idx_branches_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_branches_is_active ON public.branches USING btree (is_active);


--
-- Name: idx_branches_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_branches_parent_id ON public.branches USING btree (parent_id);


--
-- Name: idx_branches_path; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_branches_path ON public.branches USING btree (path);


--
-- Name: idx_companies_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_code ON public.companies USING btree (code);


--
-- Name: idx_companies_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_companies_is_active ON public.companies USING btree (is_active);


--
-- Name: idx_modules_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_category ON public.modules USING btree (category);


--
-- Name: idx_modules_category_parent; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_category_parent ON public.modules USING btree (category, parent_id);


--
-- Name: idx_modules_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_is_active ON public.modules USING btree (is_active);


--
-- Name: idx_modules_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_parent_id ON public.modules USING btree (parent_id);


--
-- Name: idx_modules_subscription_tier; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_subscription_tier ON public.modules USING btree (subscription_tier);


--
-- Name: idx_modules_url; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_modules_url ON public.modules USING btree (url);


--
-- Name: idx_plan_modules_module_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_plan_modules_module_id ON public.plan_modules USING btree (module_id);


--
-- Name: idx_plan_modules_plan_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_plan_modules_plan_id ON public.plan_modules USING btree (plan_id);


--
-- Name: idx_role_modules_can_approve; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_modules_can_approve ON public.role_modules USING btree (can_approve);


--
-- Name: idx_role_modules_module_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_modules_module_id ON public.role_modules USING btree (module_id);


--
-- Name: idx_role_modules_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_role_modules_role_id ON public.role_modules USING btree (role_id);


--
-- Name: idx_roles_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_is_active ON public.roles USING btree (is_active);


--
-- Name: idx_roles_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_name ON public.roles USING btree (name);


--
-- Name: idx_subscription_plans_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_plans_is_active ON public.subscription_plans USING btree (is_active);


--
-- Name: idx_subscription_plans_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscription_plans_name ON public.subscription_plans USING btree (name);


--
-- Name: idx_subscriptions_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscriptions_company_id ON public.subscriptions USING btree (company_id);


--
-- Name: idx_subscriptions_end_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscriptions_end_date ON public.subscriptions USING btree (end_date);


--
-- Name: idx_subscriptions_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_subscriptions_status ON public.subscriptions USING btree (status);


--
-- Name: idx_unit_role_modules_module_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_unit_role_modules_module_id ON public.unit_role_modules USING btree (module_id);


--
-- Name: idx_unit_role_modules_unit_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_unit_role_modules_unit_role_id ON public.unit_role_modules USING btree (unit_role_id);


--
-- Name: idx_unit_roles_role_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_unit_roles_role_id ON public.unit_roles USING btree (role_id);


--
-- Name: idx_unit_roles_unit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_unit_roles_unit_id ON public.unit_roles USING btree (unit_id);


--
-- Name: idx_units_branch_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_branch_id ON public.units USING btree (branch_id);


--
-- Name: idx_units_code; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_code ON public.units USING btree (code);


--
-- Name: idx_units_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_is_active ON public.units USING btree (is_active);


--
-- Name: idx_units_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_parent_id ON public.units USING btree (parent_id);


--
-- Name: idx_units_path; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_units_path ON public.units USING btree (path);


--
-- Name: idx_user_effective_permissions_module; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_effective_permissions_module ON public.user_roles USING btree (user_id) INCLUDE (company_id, branch_id, unit_id);


--
-- Name: idx_user_effective_permissions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_effective_permissions_user_id ON public.user_roles USING btree (user_id);


--
-- Name: idx_user_roles_company_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_company_id ON public.user_roles USING btree (company_id);


--
-- Name: idx_user_roles_unit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_unit_id ON public.user_roles USING btree (unit_id);


--
-- Name: idx_user_roles_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_roles_user_id ON public.user_roles USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_is_active ON public.users USING btree (is_active);


--
-- Name: idx_users_user_identity; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_user_identity ON public.users USING btree (user_identity);


--
-- Name: branches update_branch_hierarchy_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_branch_hierarchy_trigger BEFORE INSERT OR UPDATE ON public.branches FOR EACH ROW EXECUTE FUNCTION public.update_branch_hierarchy();


--
-- Name: branches update_branches_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_branches_updated_at BEFORE UPDATE ON public.branches FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: companies update_companies_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON public.companies FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: modules update_modules_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_modules_updated_at BEFORE UPDATE ON public.modules FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: roles update_roles_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON public.roles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: subscription_plans update_subscription_plans_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_subscription_plans_updated_at BEFORE UPDATE ON public.subscription_plans FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: subscriptions update_subscriptions_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON public.subscriptions FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: units update_unit_hierarchy_trigger; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_unit_hierarchy_trigger BEFORE INSERT OR UPDATE ON public.units FOR EACH ROW EXECUTE FUNCTION public.update_unit_hierarchy();


--
-- Name: unit_role_modules update_unit_role_modules_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_unit_role_modules_updated_at BEFORE UPDATE ON public.unit_role_modules FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: unit_roles update_unit_roles_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_unit_roles_updated_at BEFORE UPDATE ON public.unit_roles FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: units update_units_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_units_updated_at BEFORE UPDATE ON public.units FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: branches branches_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--
-- Name: branches branches_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.branches
    ADD CONSTRAINT branches_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.branches(id) ON DELETE SET NULL;


--
-- Name: modules modules_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.modules
    ADD CONSTRAINT modules_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.modules(id) ON DELETE SET NULL;


--
-- Name: plan_modules plan_modules_module_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plan_modules
    ADD CONSTRAINT plan_modules_module_id_fkey FOREIGN KEY (module_id) REFERENCES public.modules(id) ON DELETE CASCADE;


--
-- Name: plan_modules plan_modules_plan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.plan_modules
    ADD CONSTRAINT plan_modules_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.subscription_plans(id) ON DELETE CASCADE;


--
-- Name: role_modules role_modules_module_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_modules
    ADD CONSTRAINT role_modules_module_id_fkey FOREIGN KEY (module_id) REFERENCES public.modules(id) ON DELETE CASCADE;


--
-- Name: role_modules role_modules_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.role_modules
    ADD CONSTRAINT role_modules_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--
-- Name: subscriptions subscriptions_plan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subscriptions
    ADD CONSTRAINT subscriptions_plan_id_fkey FOREIGN KEY (plan_id) REFERENCES public.subscription_plans(id);


--
-- Name: unit_role_modules unit_role_modules_module_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_role_modules
    ADD CONSTRAINT unit_role_modules_module_id_fkey FOREIGN KEY (module_id) REFERENCES public.modules(id) ON DELETE CASCADE;


--
-- Name: unit_role_modules unit_role_modules_unit_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_role_modules
    ADD CONSTRAINT unit_role_modules_unit_role_id_fkey FOREIGN KEY (unit_role_id) REFERENCES public.unit_roles(id) ON DELETE CASCADE;


--
-- Name: unit_roles unit_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_roles
    ADD CONSTRAINT unit_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: unit_roles unit_roles_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.unit_roles
    ADD CONSTRAINT unit_roles_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE CASCADE;


--
-- Name: units units_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE CASCADE;


--
-- Name: units units_parent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.units
    ADD CONSTRAINT units_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.units(id) ON DELETE SET NULL;


--
-- Name: user_roles user_roles_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE SET NULL;


--
-- Name: user_roles user_roles_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.companies(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: user_roles user_roles_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES public.units(id) ON DELETE SET NULL;


--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

