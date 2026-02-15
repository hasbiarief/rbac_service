#!/bin/bash

# Script to generate seed data from local database
# Usage: ./scripts/generate-seed-data.sh

DB_HOST="localhost"
DB_USER="hasbi"
DB_NAME="huminor_rbac"
OUTPUT_FILE="migrations/006_seed_initial_data.sql"

echo "Generating seed data from local database..."
echo "-- Seed Initial Data" > $OUTPUT_FILE
echo "-- Generated from local database: $(date)" >> $OUTPUT_FILE
echo "" >> $OUTPUT_FILE

# Users
echo "-- Insert Users" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO users (id, name, email, user_identity, password_hash, is_active, created_at, updated_at) VALUES (' || 
  id || ', ' || 
  quote_literal(name) || ', ' || 
  quote_literal(email) || ', ' || 
  quote_literal(user_identity) || ', ' || 
  quote_literal(password_hash) || ', ' || 
  is_active || ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);'
FROM users ORDER BY id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE

# Companies
echo "-- Insert Companies" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO companies (id, name, code, is_active, created_at, updated_at) VALUES (' || 
  id || ', ' || 
  quote_literal(name) || ', ' || 
  quote_literal(code) || ', ' || 
  is_active || ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);'
FROM companies ORDER BY id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE

# Branches
echo "-- Insert Branches" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO branches (id, company_id, parent_id, name, code, level, path, is_active, created_at, updated_at) VALUES (' || 
  id || ', ' || 
  company_id || ', ' || 
  COALESCE(parent_id::text, 'NULL') || ', ' || 
  quote_literal(name) || ', ' || 
  quote_literal(code) || ', ' || 
  COALESCE(level::text, '0') || ', ' || 
  quote_literal(COALESCE(path, '')) || ', ' || 
  is_active || ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);'
FROM branches ORDER BY id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE

# Roles
echo "-- Insert Roles" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO roles (id, name, description, is_active, created_at, updated_at) VALUES (' || 
  id || ', ' || 
  quote_literal(name) || ', ' || 
  COALESCE(quote_literal(description), 'NULL') || ', ' || 
  is_active || ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);'
FROM roles ORDER BY id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE

# Units
echo "-- Insert Units" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO units (id, branch_id, name, code, description, is_active, created_at, updated_at) VALUES (' || 
  id || ', ' || 
  branch_id || ', ' || 
  quote_literal(name) || ', ' || 
  quote_literal(code) || ', ' || 
  COALESCE(quote_literal(description), 'NULL') || ', ' || 
  is_active || ', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);'
FROM units ORDER BY id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE

# User Roles
echo "-- Insert User Roles" >> $OUTPUT_FILE
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -t -c "
SELECT 'INSERT INTO user_roles (user_id, role_id, company_id, branch_id, unit_id, created_at) VALUES (' || 
  user_id || ', ' || 
  role_id || ', ' || 
  company_id || ', ' || 
  COALESCE(branch_id::text, 'NULL') || ', ' || 
  COALESCE(unit_id::text, 'NULL') || ', CURRENT_TIMESTAMP);'
FROM user_roles ORDER BY user_id, role_id;" | grep INSERT >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE
echo "-- Update sequences" >> $OUTPUT_FILE
echo "SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));" >> $OUTPUT_FILE
echo "SELECT setval('companies_id_seq', (SELECT MAX(id) FROM companies));" >> $OUTPUT_FILE
echo "SELECT setval('branches_id_seq', (SELECT MAX(id) FROM branches));" >> $OUTPUT_FILE
echo "SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));" >> $OUTPUT_FILE
echo "SELECT setval('units_id_seq', (SELECT MAX(id) FROM units));" >> $OUTPUT_FILE

echo "✅ Seed data generated: $OUTPUT_FILE"
