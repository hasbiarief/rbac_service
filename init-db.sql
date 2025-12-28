-- Database initialization script for production
-- This script will be executed when PostgreSQL container starts for the first time

-- Create database if not exists (handled by POSTGRES_DB env var)
-- CREATE DATABASE huminor_rbac;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
SET timezone = 'UTC';

-- Create application user (optional, for better security)
-- CREATE USER huminor_app WITH PASSWORD 'secure_password';
-- GRANT ALL PRIVILEGES ON DATABASE huminor_rbac TO huminor_app;