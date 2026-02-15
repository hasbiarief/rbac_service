#!/bin/bash

# migrate-postman.sh - Script for migrating Postman collections to Swagger annotations
# This script converts Postman collections to Swagger annotations, generates initial
# documentation, and validates the conversion results.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
POSTMAN_FILE="${POSTMAN_FILE:-postman_collection.json}"
OUTPUT_DIR="${OUTPUT_DIR:-.}"
SWAGGER_CMD="${SWAGGER_CMD:-./bin/swagger}"
DOCS_DIR="${DOCS_DIR:-docs}"

# Functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Check if required tools are available
check_requirements() {
    print_header "Checking Requirements"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go first."
        exit 1
    fi
    print_success "Go is installed"
    
    # Check if swag is installed
    if ! command -v swag &> /dev/null; then
        print_warning "swag is not installed. Installing..."
        go install github.com/swaggo/swag/cmd/swag@latest
        if ! command -v swag &> /dev/null; then
            print_error "Failed to install swag"
            exit 1
        fi
        print_success "swag installed successfully"
    else
        print_success "swag is installed"
    fi
    
    # Check if swagger CLI tool exists
    if [ ! -f "$SWAGGER_CMD" ]; then
        print_warning "Swagger CLI tool not found. Building..."
        make build-swagger || {
            print_error "Failed to build swagger CLI tool"
            exit 1
        }
        print_success "Swagger CLI tool built successfully"
    else
        print_success "Swagger CLI tool found"
    fi
    
    echo ""
}

# Validate Postman collection file
validate_postman_file() {
    print_header "Validating Postman Collection"
    
    if [ ! -f "$POSTMAN_FILE" ]; then
        print_error "Postman collection file not found: $POSTMAN_FILE"
        print_info "Usage: POSTMAN_FILE=path/to/collection.json $0"
        exit 1
    fi
    
    # Check if file is valid JSON
    if ! jq empty "$POSTMAN_FILE" 2>/dev/null; then
        print_error "Invalid JSON in Postman collection file"
        exit 1
    fi
    
    # Check if it's a Postman collection (has required fields)
    if ! jq -e '.info.name' "$POSTMAN_FILE" &> /dev/null; then
        print_error "File does not appear to be a valid Postman collection"
        exit 1
    fi
    
    COLLECTION_NAME=$(jq -r '.info.name' "$POSTMAN_FILE")
    ENDPOINT_COUNT=$(jq '[.. | .request? | select(. != null)] | length' "$POSTMAN_FILE")
    
    print_success "Valid Postman collection: $COLLECTION_NAME"
    print_info "Found $ENDPOINT_COUNT endpoints"
    echo ""
}

# Run the converter
run_converter() {
    print_header "Converting Postman Collection to Swagger Annotations"
    
    # Create a temporary Go program to run the converter
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    cat > "$TEMP_DIR/main.go" << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/huminor-erp/gin-scalable-api/pkg/swagger"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: converter <postman-file> <output-dir>")
		os.Exit(1)
	}
	
	postmanFile := os.Args[1]
	outputDir := os.Args[2]
	
	converter := swagger.NewConverter()
	report, err := converter.Convert(postmanFile, outputDir)
	if err != nil {
		fmt.Printf("Conversion failed: %v\n", err)
		os.Exit(1)
	}
	
	// Print report as JSON
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(reportJSON))
}
EOF
    
    # Build and run the converter
    cd "$TEMP_DIR"
    go mod init converter 2>/dev/null || true
    go mod edit -replace github.com/huminor-erp/gin-scalable-api="$OUTPUT_DIR"
    go mod tidy
    
    print_info "Running converter..."
    REPORT=$(go run main.go "$POSTMAN_FILE" "$OUTPUT_DIR")
    
    if [ $? -ne 0 ]; then
        print_error "Conversion failed"
        echo "$REPORT"
        exit 1
    fi
    
    # Parse and display report
    TOTAL=$(echo "$REPORT" | jq -r '.TotalEndpoints')
    CONVERTED=$(echo "$REPORT" | jq -r '.ConvertedEndpoints')
    FAILED=$(echo "$REPORT" | jq -r '.FailedEndpoints | length')
    FILES=$(echo "$REPORT" | jq -r '.GeneratedFiles | length')
    
    print_success "Conversion completed"
    print_info "Total endpoints: $TOTAL"
    print_info "Converted: $CONVERTED"
    
    if [ "$FAILED" -gt 0 ]; then
        print_warning "Failed: $FAILED"
        echo ""
        print_warning "Failed endpoints:"
        echo "$REPORT" | jq -r '.FailedEndpoints[] | "  - \(.Name): \(.Reason)"'
    fi
    
    echo ""
    print_success "Generated $FILES annotation files:"
    echo "$REPORT" | jq -r '.GeneratedFiles[] | "  - \(.)"'
    echo ""
}

# Generate initial Swagger documentation
generate_swagger() {
    print_header "Generating Swagger Documentation"
    
    print_info "Running swag init..."
    swag init \
        --generalInfo cmd/api/main.go \
        --output "$DOCS_DIR" \
        --dir ./ \
        --parseInternal \
        --exclude vendor,tmp,bin
    
    if [ $? -ne 0 ]; then
        print_error "Swagger generation failed"
        exit 1
    fi
    
    print_success "Swagger documentation generated"
    
    # Check generated files
    if [ -f "$DOCS_DIR/swagger.json" ]; then
        print_success "Generated: $DOCS_DIR/swagger.json"
        ENDPOINT_COUNT=$(jq '.paths | length' "$DOCS_DIR/swagger.json")
        print_info "Documented endpoints: $ENDPOINT_COUNT"
    else
        print_error "swagger.json not generated"
        exit 1
    fi
    
    if [ -f "$DOCS_DIR/swagger.yaml" ]; then
        print_success "Generated: $DOCS_DIR/swagger.yaml"
    fi
    
    echo ""
}

# Validate conversion results
validate_conversion() {
    print_header "Validating Conversion Results"
    
    # Validate annotation files exist
    print_info "Checking annotation files..."
    ANNOTATION_FILES=$(find internal/modules -name "swagger.go" -path "*/docs/swagger.go" 2>/dev/null | wc -l)
    
    if [ "$ANNOTATION_FILES" -eq 0 ]; then
        print_warning "No annotation files found in internal/modules/*/docs/"
    else
        print_success "Found $ANNOTATION_FILES annotation files"
    fi
    
    # Validate generated documentation
    print_info "Validating generated documentation..."
    
    if [ ! -f "$DOCS_DIR/swagger.json" ]; then
        print_error "swagger.json not found"
        exit 1
    fi
    
    # Check if swagger.json is valid
    if ! jq empty "$DOCS_DIR/swagger.json" 2>/dev/null; then
        print_error "Invalid JSON in swagger.json"
        exit 1
    fi
    
    # Validate OpenAPI structure
    if ! jq -e '.openapi' "$DOCS_DIR/swagger.json" &> /dev/null && \
       ! jq -e '.swagger' "$DOCS_DIR/swagger.json" &> /dev/null; then
        print_error "Generated file is not a valid OpenAPI/Swagger specification"
        exit 1
    fi
    
    print_success "swagger.json is valid"
    
    # Check for paths
    PATH_COUNT=$(jq '.paths | length' "$DOCS_DIR/swagger.json")
    if [ "$PATH_COUNT" -eq 0 ]; then
        print_warning "No paths found in generated documentation"
    else
        print_success "Found $PATH_COUNT paths in documentation"
    fi
    
    # Run annotation validator if available
    if [ -f "$SWAGGER_CMD" ]; then
        print_info "Running annotation validator..."
        if $SWAGGER_CMD validate; then
            print_success "All annotations are valid"
        else
            print_warning "Some annotations have validation issues"
        fi
    fi
    
    echo ""
}

# Generate migration report
generate_report() {
    print_header "Migration Report"
    
    REPORT_FILE="migration-report.txt"
    
    {
        echo "Postman to Swagger Migration Report"
        echo "===================================="
        echo ""
        echo "Date: $(date)"
        echo "Postman Collection: $POSTMAN_FILE"
        echo "Collection Name: $COLLECTION_NAME"
        echo ""
        echo "Results:"
        echo "--------"
        echo "Annotation files generated: $ANNOTATION_FILES"
        echo "Endpoints documented: $PATH_COUNT"
        echo ""
        echo "Generated Files:"
        echo "----------------"
        find internal/modules -name "swagger.go" -path "*/docs/swagger.go" 2>/dev/null || echo "None"
        echo ""
        echo "$DOCS_DIR/swagger.json"
        echo "$DOCS_DIR/swagger.yaml"
        echo "$DOCS_DIR/docs.go"
        echo ""
        echo "Next Steps:"
        echo "-----------"
        echo "1. Review generated annotation files in internal/modules/*/docs/swagger.go"
        echo "2. Update annotations with accurate type information and descriptions"
        echo "3. Run 'make swagger-gen' to regenerate documentation after changes"
        echo "4. Test Swagger UI at http://localhost:8081/swagger/index.html"
        echo "5. Validate all endpoints work correctly"
        echo ""
    } | tee "$REPORT_FILE"
    
    print_success "Report saved to: $REPORT_FILE"
    echo ""
}

# Main execution
main() {
    print_header "Postman to Swagger Migration Tool"
    echo ""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -f|--file)
                POSTMAN_FILE="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -d|--docs)
                DOCS_DIR="$2"
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  -f, --file FILE     Postman collection file (default: postman_collection.json)"
                echo "  -o, --output DIR    Output directory (default: .)"
                echo "  -d, --docs DIR      Documentation directory (default: docs)"
                echo "  -h, --help          Show this help message"
                echo ""
                echo "Environment Variables:"
                echo "  POSTMAN_FILE        Postman collection file path"
                echo "  OUTPUT_DIR          Output directory for generated files"
                echo "  DOCS_DIR            Documentation output directory"
                echo "  SWAGGER_CMD         Path to swagger CLI tool"
                echo ""
                echo "Examples:"
                echo "  $0 -f collection.json"
                echo "  POSTMAN_FILE=my-api.json $0"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use -h or --help for usage information"
                exit 1
                ;;
        esac
    done
    
    # Execute migration steps
    check_requirements
    validate_postman_file
    run_converter
    generate_swagger
    validate_conversion
    generate_report
    
    print_header "Migration Completed Successfully!"
    print_success "Your Postman collection has been migrated to Swagger annotations"
    print_info "Review the migration report for next steps"
    echo ""
}

# Run main function
main "$@"
