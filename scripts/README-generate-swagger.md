# Swagger Documentation Generation Scripts

This directory contains scripts and tools for generating and managing Swagger/OpenAPI documentation.

## Available Tools

### 1. Bash Script: `generate-swagger.sh`

A comprehensive bash script for generating and validating Swagger documentation.

**Features:**
- Generate Swagger documentation from code annotations
- Validate annotations before generation
- Clean generated files
- Verbose output mode
- Error handling with clear messages

**Usage:**

```bash
# Generate documentation
./scripts/generate-swagger.sh

# Validate annotations only
./scripts/generate-swagger.sh -v

# Clean and regenerate
./scripts/generate-swagger.sh -c

# Verbose mode
./scripts/generate-swagger.sh --verbose

# Custom output directory
./scripts/generate-swagger.sh -o ./api-docs

# Show help
./scripts/generate-swagger.sh --help
```

**Options:**
- `-h, --help` - Show help message
- `-o, --output DIR` - Output directory (default: docs)
- `-d, --dir DIR` - Search directory (default: ./)
- `-v, --validate` - Validate annotations only (no generation)
- `-c, --clean` - Clean generated files before generation
- `--verbose` - Enable verbose output

### 2. Go CLI Tool: `cmd/swagger/main.go`

A Go-based CLI tool for Swagger documentation management.

**Features:**
- Generate documentation
- Validate annotations
- Watch mode for auto-regeneration
- Configurable output and search directories

**Build:**

```bash
# Build the CLI tool
make build-swagger

# Or manually
go build -o bin/swagger cmd/swagger/main.go
```

**Usage:**

```bash
# Generate documentation
./bin/swagger -generate

# Validate annotations
./bin/swagger -validate

# Watch for changes and auto-regenerate
./bin/swagger -watch

# Custom output directory
./bin/swagger -generate -output ./api-docs

# Custom search directory
./bin/swagger -generate -dir ./internal

# Show help
./bin/swagger -help
```

**Options:**
- `-generate` - Generate Swagger documentation
- `-validate` - Validate Swagger annotations
- `-watch` - Watch for changes and regenerate
- `-output string` - Output directory (default: "docs")
- `-dir string` - Directory to search for annotations (default: "./")
- `-help` - Show help message

## Makefile Targets

The project Makefile includes convenient targets for Swagger operations:

```bash
# Generate documentation
make swagger-gen

# Validate annotations
make swagger-validate

# Watch mode (checks every 5 seconds)
make swagger-watch

# Clean generated files
make swagger-clean

# Generate all documentation (alias)
make docs
```

## Generated Files

Both tools generate the following files in the output directory (default: `docs/`):

- `swagger.json` - OpenAPI specification in JSON format
- `swagger.yaml` - OpenAPI specification in YAML format
- `docs.go` - Go code for embedding the specification

## Requirements

- Go 1.16 or higher
- [swaggo/swag](https://github.com/swaggo/swag) tool installed

**Install swag:**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

The bash script will automatically install swag if it's not found.

## Accessing Documentation

After generating the documentation, you can access the Swagger UI at:

```
http://localhost:8081/swagger/index.html
```

The API specification is also available at:

```
http://localhost:8081/swagger/doc.json
```

## Error Handling

Both tools include comprehensive error handling:

- **Validation errors**: Shows file, line number, and error message
- **Generation errors**: Clear error messages with suggestions
- **Missing dependencies**: Automatic installation (bash script) or clear instructions

**Example error output:**

```
Found 2 validation errors:
  internal/modules/auth/docs/swagger.go:15 - missing required field: @Summary
  internal/modules/user/docs/swagger.go:42 - invalid parameter type
```

## Watch Mode

The watch mode continuously monitors for changes and regenerates documentation:

**Using CLI tool:**
```bash
./bin/swagger -watch
```

**Using Makefile:**
```bash
make swagger-watch
```

Watch mode checks for changes every 2 seconds (CLI) or 5 seconds (Makefile) and automatically regenerates documentation when annotations are modified.

## Best Practices

1. **Validate before committing**: Always run validation before committing changes
   ```bash
   ./scripts/generate-swagger.sh -v
   ```

2. **Use watch mode during development**: Enable watch mode while working on annotations
   ```bash
   ./bin/swagger -watch
   ```

3. **Clean regeneration**: Use clean flag when troubleshooting
   ```bash
   ./scripts/generate-swagger.sh -c
   ```

4. **Verbose output for debugging**: Enable verbose mode to see detailed error messages
   ```bash
   ./scripts/generate-swagger.sh --verbose
   ```

## Integration with CI/CD

Add validation to your CI/CD pipeline:

```yaml
# Example GitHub Actions
- name: Validate Swagger Annotations
  run: ./scripts/generate-swagger.sh -v

- name: Generate Documentation
  run: make swagger-gen
```

## Troubleshooting

**Problem**: `swag: command not found`

**Solution**: Install swag tool
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Problem**: Validation fails with syntax errors

**Solution**: Run with verbose mode to see detailed errors
```bash
./scripts/generate-swagger.sh --verbose
```

**Problem**: Generated files are empty or incomplete

**Solution**: 
1. Check that annotations are in the correct format
2. Ensure all modules have `docs/swagger.go` files
3. Verify the main API file has general info annotations

## Related Documentation

- [Swagger Annotation Guide](../internal/swagger/README.md)
- [Migration from Postman](./README-migrate-postman.md)
- [API Documentation](../docs/INDEX.md)

## Requirements Validation

This implementation satisfies the following requirements:

- **Requirement 6.1**: Provides make commands and scripts for generating documentation
- **Requirement 6.3**: Implements watch mode for auto-regeneration during development
- **Requirement 6.4**: Includes comprehensive error handling and validation
