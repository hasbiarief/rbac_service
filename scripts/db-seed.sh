#!/bin/bash

# Database Seeder Script for RBAC Service
# Seeds database with template data for new projects

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
SEEDER_DIR="database/seeders"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🌱 RBAC Database Seeder${NC}"
echo -e "${BLUE}======================${NC}"
echo ""

echo -e "${YELLOW}📊 Database Info:${NC}"
echo "  Host: $DB_HOST:$DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Seeder Directory: $SEEDER_DIR"
echo ""

# Check if seeder directory exists
if [ ! -d "$SEEDER_DIR" ]; then
    echo -e "${RED}❌ Seeder directory not found: $SEEDER_DIR${NC}"
    echo -e "${YELLOW}💡 Run 'make db-dump' first to generate seeder files${NC}"
    exit 1
fi

# Check if template file exists
if [ ! -f "$SEEDER_DIR/template.sql" ]; then
    echo -e "${RED}❌ Template seeder not found: $SEEDER_DIR/template.sql${NC}"
    echo -e "${YELLOW}💡 Run 'make db-dump' first to generate seeder files${NC}"
    exit 1
fi

# Check if database exists
echo -e "${YELLOW}🔍 Checking database connection...${NC}"
if ! PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
    echo -e "${RED}❌ Cannot connect to database. Please check your connection settings.${NC}"
    echo -e "${YELLOW}💡 Create database first with: make db-create${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Database connection successful${NC}"
echo ""

# Seed database with template
echo -e "${YELLOW}🌱 Seeding database with template data...${NC}"
PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $SEEDER_DIR/template.sql > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Database seeded successfully!${NC}"
else
    echo -e "${RED}❌ Failed to seed database${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Database seeding completed successfully!${NC}"
echo ""
echo -e "${YELLOW}🔑 Default Login Credentials:${NC}"
echo "  📧 naruto@company.com (Super Admin) - ID: 100000001"
echo "  📧 sakura@company.com (HR Admin) - ID: 100000002"
echo "  📧 sasuke@company.com (Recruiter) - ID: 100000003"
echo "  � ino@company.com (Employee) - ID: 100000006"
echo "  📧 hasbi@company.com (Console Admin) - ID: 800000001"
echo "  🔐 Password: password123 (for all users)"
echo ""
echo -e "${YELLOW}� Login Options:${NC}"
echo "  • Use either email or user_identity (ID) to login"
echo "  • Example: user_identity: '800000001', password: 'password123'"
echo ""
echo -e "${YELLOW}� Next Steps:${NC}"
echo "  1. Start the server: make run"
echo "  2. Test API: http://localhost:8081/api/v1/auth/login"
echo "  3. View all users: See database/README.md for complete list"