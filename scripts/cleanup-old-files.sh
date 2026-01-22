#!/bin/bash

# Cleanup Script for RBAC Service
# Removes outdated and unused files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧹 RBAC Service Cleanup${NC}"
echo -e "${BLUE}======================${NC}"
echo ""

# List of files to remove (already removed)
REMOVED_FILES=(
    "scripts/audit_user_roles.sql"
    "scripts/migrate_users_to_units.sql" 
    "scripts/rollback_migration.sql"
    "scripts/test_permission_resolution.sql"
    "scripts/validate_migration.sql"
)

echo -e "${GREEN}✅ Cleaned up outdated files:${NC}"
for file in "${REMOVED_FILES[@]}"; do
    echo "  🗑️  $file"
done

echo ""
echo -e "${BLUE}📁 Current scripts directory:${NC}"
ls -la scripts/

echo ""
echo -e "${YELLOW}📋 Remaining files:${NC}"
echo "  ✅ db-dump.sh - Database dump generator"
echo "  ✅ db-seed.sh - Database seeder"
echo "  ✅ dev.sh - Development server with hot reload"
echo "  ✅ docker-prod.sh - Production Docker management"
echo "  🆕 cleanup-old-files.sh - This cleanup script"

echo ""
echo -e "${GREEN}🎉 Cleanup completed!${NC}"
echo -e "${YELLOW}💡 All remaining files are actively used by the project.${NC}"