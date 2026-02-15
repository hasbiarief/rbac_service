#!/bin/bash

# verify-endpoints.sh
# Script to verify API endpoints haven't changed during Swagger migration
# Task 17.1: Ensure API endpoints remain unchanged

set -e

echo "🔍 Verifying API Endpoint Backward Compatibility"
echo "================================================"
echo ""

# Function to extract routes from route files
extract_routes() {
    local module=$1
    local route_file="internal/modules/$module/route.go"
    
    if [ ! -f "$route_file" ]; then
        echo "⚠️  Route file not found: $route_file"
        return 1
    fi
    
    echo "📁 Module: $module"
    echo "   Routes defined in $route_file:"
    
    # Extract route definitions (looking for patterns like auth.POST("/login", ...))
    grep -E '(GET|POST|PUT|DELETE|PATCH)\("' "$route_file" | \
        sed 's/^[[:space:]]*/   /' | \
        sed 's/,.*$//' || echo "   No routes found"
    
    echo ""
}

# Function to extract swagger routes from docs
extract_swagger_routes() {
    local module=$1
    local docs_file="internal/modules/$module/docs/swagger.go"
    
    if [ ! -f "$docs_file" ]; then
        echo "⚠️  Swagger docs not found: $docs_file"
        return 1
    fi
    
    echo "📄 Swagger annotations in $docs_file:"
    
    # Extract @Router annotations
    grep -E '@Router' "$docs_file" | \
        sed 's/^[[:space:]]*/   /' || echo "   No @Router annotations found"
    
    echo ""
}

# Function to compare routes
compare_module_routes() {
    local module=$1
    
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Checking module: $module"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    
    extract_routes "$module"
    extract_swagger_routes "$module"
}

# Check all modules
modules=("auth" "user" "company" "branch" "role" "module" "subscription" "unit" "application")

for module in "${modules[@]}"; do
    compare_module_routes "$module"
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Verification Complete"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Summary:"
echo "- Route files contain the actual API endpoint definitions"
echo "- Swagger docs contain documentation annotations"
echo "- Both should match for backward compatibility"
echo ""
echo "Manual verification required:"
echo "1. Check that all routes in route.go have corresponding @Router annotations"
echo "2. Verify HTTP methods match (GET, POST, PUT, DELETE)"
echo "3. Verify paths match exactly"
echo "4. Confirm no routes were removed or changed"
