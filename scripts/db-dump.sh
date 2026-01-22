#!/bin/bash

# Database Dump Script for RBAC Service
# Creates a complete database dump with structure and data

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASS=${DB_PASS:-password}
DB_NAME=${DB_NAME:-huminor_rbac}
DUMP_DIR="database/dumps"
SEEDER_DIR="database/seeders"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🗄️  RBAC Database Dump Generator${NC}"
echo -e "${BLUE}=================================${NC}"
echo ""

# Create dump directories
mkdir -p $DUMP_DIR
mkdir -p $SEEDER_DIR

echo -e "${YELLOW}📊 Database Info:${NC}"
echo "  Host: $DB_HOST:$DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Dump Directory: $DUMP_DIR"
echo "  Seeder Directory: $SEEDER_DIR"
echo ""

# Check if database exists
echo -e "${YELLOW}🔍 Checking database connection...${NC}"
if ! PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo -e "${RED}❌ Cannot connect to database. Please check your connection settings.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Database connection successful${NC}"
echo ""

# Create structure-only dump
echo -e "${YELLOW}🏗️  Creating structure dump...${NC}"
PGPASSWORD=$DB_PASS pg_dump \
    -h $DB_HOST -p $DB_PORT -U $DB_USER \
    --schema-only \
    --no-owner \
    --no-privileges \
    --clean \
    --if-exists \
    $DB_NAME > $DUMP_DIR/structure.sql

echo -e "${GREEN}✅ Structure dump created: $DUMP_DIR/structure.sql${NC}"

# Create data-only dump
echo -e "${YELLOW}📦 Creating data dump...${NC}"
PGPASSWORD=$DB_PASS pg_dump \
    -h $DB_HOST -p $DB_PORT -U $DB_USER \
    --data-only \
    --no-owner \
    --no-privileges \
    --column-inserts \
    $DB_NAME > $DUMP_DIR/data.sql

echo -e "${GREEN}✅ Data dump created: $DUMP_DIR/data.sql${NC}"

# Create complete dump
echo -e "${YELLOW}🎯 Creating complete dump...${NC}"
PGPASSWORD=$DB_PASS pg_dump \
    -h $DB_HOST -p $DB_PORT -U $DB_USER \
    --no-owner \
    --no-privileges \
    --clean \
    --if-exists \
    --column-inserts \
    $DB_NAME > $DUMP_DIR/complete_${TIMESTAMP}.sql

echo -e "${GREEN}✅ Complete dump created: $DUMP_DIR/complete_${TIMESTAMP}.sql${NC}"

# Create template dump (structure + essential data only)
echo -e "${YELLOW}📋 Creating template dump...${NC}"
PGPASSWORD=$DB_PASS pg_dump \
    -h $DB_HOST -p $DB_PORT -U $DB_USER \
    --no-owner \
    --no-privileges \
    --clean \
    --if-exists \
    --column-inserts \
    --exclude-table-data=audit_logs \
    $DB_NAME > $SEEDER_DIR/template.sql

echo -e "${GREEN}✅ Template dump created: $SEEDER_DIR/template.sql${NC}"

# Create seeder files for each table
echo -e "${YELLOW}🌱 Creating individual seeder files...${NC}"

# Get list of tables
TABLES=$(PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "
SELECT tablename FROM pg_tables 
WHERE schemaname = 'public' 
AND tablename NOT IN ('schema_migrations', 'audit_logs', 'user_roles_backup')
ORDER BY tablename;
")

for table in $TABLES; do
    table=$(echo $table | xargs) # trim whitespace
    if [ ! -z "$table" ]; then
        echo "  � Creating seeder for: $table"
        PGPASSWORD=$DB_PASS pg_dump \
            -h $DB_HOST -p $DB_PORT -U $DB_USER \
            --data-only \
            --no-owner \
            --no-privileges \
            --column-inserts \
            --table=$table \
            $DB_NAME > $SEEDER_DIR/${table}_seeder.sql
    fi
done

echo -e "${GREEN}✅ Individual seeder files created in: $SEEDER_DIR${NC}"

# Get database statistics
echo ""
echo -e "${YELLOW}📊 Database Statistics:${NC}"
PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
SELECT 
    schemaname,
    relname as tablename,
    n_live_tup as rows
FROM pg_stat_user_tables 
WHERE schemaname = 'public'
ORDER BY n_live_tup DESC;
" -t

echo ""
echo -e "${GREEN}🎉 Database dump completed successfully!${NC}"
echo ""
echo -e "${BLUE}📁 Generated files:${NC}"
echo "  📄 $DUMP_DIR/structure.sql - Database structure only"
echo "  📦 $DUMP_DIR/data.sql - All data"
echo "  🎯 $DUMP_DIR/complete_${TIMESTAMP}.sql - Complete dump with timestamp"
echo "  📋 $SEEDER_DIR/template.sql - Template for new projects (no audit logs)"
echo "  🌱 $SEEDER_DIR/*_seeder.sql - Individual table seeders"
echo ""
echo -e "${YELLOW}💡 Usage:${NC}"
echo "  # Restore structure only:"
echo "  psql -h localhost -U $DB_USER -d new_db < $DUMP_DIR/structure.sql"
echo ""
echo "  # Restore complete database:"
echo "  psql -h localhost -U $DB_USER -d new_db < $SEEDER_DIR/template.sql"
echo ""
echo "  # Or use the seeder commands:"
echo "  make db-seed"
echo "  make db-seed-fresh"