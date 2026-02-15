#!/bin/bash

# generate-swagger.sh
# Script untuk generate dan validate Swagger documentation
# Requirements: 6.1, 6.4

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
OUTPUT_DIR="docs"
SEARCH_DIR="./"
VALIDATE_ONLY=false
VERBOSE=false
CLEAN=false

# Function to print colored messages
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Function to show usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Generate and validate Swagger/OpenAPI documentation from code annotations.

OPTIONS:
    -h, --help              Show this help message
    -o, --output DIR        Output directory (default: docs)
    -d, --dir DIR           Search directory (default: ./)
    -v, --validate          Validate annotations only (no generation)
    -c, --clean             Clean generated files before generation
    --verbose               Enable verbose output
    
EXAMPLES:
    # Generate documentation
    $0
    
    # Generate with custom output directory
    $0 -o ./api-docs
    
    # Validate annotations only
    $0 -v
    
    # Clean and regenerate
    $0 -c
    
    # Verbose mode
    $0 --verbose

EOF
    exit 0
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        -o|--output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -d|--dir)
            SEARCH_DIR="$2"
            shift 2
            ;;
        -v|--validate)
            VALIDATE_ONLY=true
            shift
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Function to check if swag is installed
check_swag() {
    if ! command -v swag &> /dev/null; then
        print_error "swag command not found"
        print_info "Installing swag..."
        
        if ! go install github.com/swaggo/swag/cmd/swag@latest; then
            print_error "Failed to install swag"
            print_info "Please install manually: go install github.com/swaggo/swag/cmd/swag@latest"
            exit 1
        fi
        
        print_success "swag installed successfully"
    fi
}

# Function to clean generated files
clean_files() {
    if [ "$CLEAN" = true ]; then
        print_info "Cleaning generated files..."
        
        if [ -d "$OUTPUT_DIR" ]; then
            rm -f "$OUTPUT_DIR/swagger.json" "$OUTPUT_DIR/swagger.yaml" "$OUTPUT_DIR/docs.go"
            print_success "Cleaned generated files"
        else
            print_warning "Output directory does not exist: $OUTPUT_DIR"
        fi
    fi
}

# Function to validate annotations
validate_annotations() {
    print_info "Validating Swagger annotations..."
    
    local temp_dir=$(mktemp -d)
    local exit_code=0
    
    # Run swag init to validate (output to temp directory)
    if [ "$VERBOSE" = true ]; then
        swag init -g cmd/api/main.go \
            --output "$temp_dir" \
            --dir "$SEARCH_DIR" \
            --parseDependency \
            --parseInternal 2>&1 || exit_code=$?
    else
        swag init -g cmd/api/main.go \
            --output "$temp_dir" \
            --dir "$SEARCH_DIR" \
            --parseDependency \
            --parseInternal > /dev/null 2>&1 || exit_code=$?
    fi
    
    # Clean up temp directory
    rm -rf "$temp_dir"
    
    if [ $exit_code -eq 0 ]; then
        print_success "All annotations are valid"
        return 0
    else
        print_error "Validation failed - annotations contain errors"
        print_info "Run with --verbose flag to see detailed error messages"
        return 1
    fi
}

# Function to generate documentation
generate_docs() {
    print_info "Generating Swagger documentation..."
    
    # Create output directory if it doesn't exist
    mkdir -p "$OUTPUT_DIR"
    
    # Run swag init
    local exit_code=0
    if [ "$VERBOSE" = true ]; then
        swag init -g cmd/api/main.go \
            --output "$OUTPUT_DIR" \
            --dir "$SEARCH_DIR" \
            --parseDependency \
            --parseInternal || exit_code=$?
    else
        swag init -g cmd/api/main.go \
            --output "$OUTPUT_DIR" \
            --dir "$SEARCH_DIR" \
            --parseDependency \
            --parseInternal > /dev/null 2>&1 || exit_code=$?
    fi
    
    if [ $exit_code -ne 0 ]; then
        print_error "Documentation generation failed"
        print_info "Run with --verbose flag to see detailed error messages"
        return 1
    fi
    
    # Verify generated files
    if [ ! -f "$OUTPUT_DIR/swagger.json" ] || [ ! -f "$OUTPUT_DIR/swagger.yaml" ]; then
        print_error "Generated files not found"
        return 1
    fi
    
    print_success "Documentation generated successfully"
    print_info "Output files:"
    echo "  - $OUTPUT_DIR/swagger.json"
    echo "  - $OUTPUT_DIR/swagger.yaml"
    echo "  - $OUTPUT_DIR/docs.go"
    
    # Show file sizes
    local json_size=$(du -h "$OUTPUT_DIR/swagger.json" | cut -f1)
    local yaml_size=$(du -h "$OUTPUT_DIR/swagger.yaml" | cut -f1)
    print_info "File sizes: JSON=$json_size, YAML=$yaml_size"
    
    return 0
}

# Function to count endpoints
count_endpoints() {
    if [ -f "$OUTPUT_DIR/swagger.json" ]; then
        local count=$(grep -o '"paths"' "$OUTPUT_DIR/swagger.json" | wc -l)
        if [ $count -gt 0 ]; then
            local endpoint_count=$(grep -o '"\/' "$OUTPUT_DIR/swagger.json" | wc -l)
            print_info "Total endpoints documented: $endpoint_count"
        fi
    fi
}

# Main execution
main() {
    echo ""
    print_info "Swagger Documentation Generator"
    echo ""
    
    # Check prerequisites
    check_swag
    
    # Clean if requested
    clean_files
    
    # Validate or generate
    if [ "$VALIDATE_ONLY" = true ]; then
        if validate_annotations; then
            exit 0
        else
            exit 1
        fi
    else
        # Validate first
        if ! validate_annotations; then
            print_error "Cannot generate documentation - validation failed"
            exit 1
        fi
        
        # Generate documentation
        if generate_docs; then
            count_endpoints
            echo ""
            print_success "Done!"
            print_info "Access Swagger UI at: http://localhost:8081/swagger/index.html"
            exit 0
        else
            exit 1
        fi
    fi
}

# Run main function
main
