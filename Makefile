# Makefile for ERP RBAC System without GORM

.PHONY: build run migrate-up migrate-status clean test

# Build the application
build:
	go build -o bin/server cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run the application
run: build
	./bin/server

# Run migrations
migrate-up:
	./bin/migrate -action=up -dir=migrations

# Check migration status
migrate-status:
	./bin/migrate -action=status -dir=migrations

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f server migrate

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Please configure your .env file"

# Database operations
db-create:
	createdb huminor_rbac

db-drop:
	dropdb huminor_rbac

db-reset: db-drop db-create migrate-up

# Docker operations
docker-build:
	docker build -t huminor-rbac .

docker-run:
	docker-compose -f docker-compose.prod.yml up -d

docker-stop:
	docker-compose -f docker-compose.prod.yml down

docker-logs:
	docker-compose -f docker-compose.prod.yml logs -f

docker-clean:
	docker-compose -f docker-compose.prod.yml down -v --remove-orphans
	docker system prune -f

# Production Docker operations
prod-start:
	./scripts/docker-prod.sh start

prod-stop:
	./scripts/docker-prod.sh stop

prod-logs:
	./scripts/docker-prod.sh logs

prod-migrate:
	./scripts/docker-prod.sh migrate

prod-backup:
	./scripts/docker-prod.sh backup

prod-status:
	./scripts/docker-prod.sh status

# Linting and formatting
fmt:
	go fmt ./...

lint:
	golangci-lint run

# Generate documentation
docs:
	@echo "Generating API documentation..."
	@echo "Documentation available in docs/ folder"

# Production build
prod-build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/api/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/migrate cmd/migrate/main.go

# Help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  deps         - Install dependencies"
	@echo "  dev-setup    - Setup development environment"
	@echo "  test         - Run tests"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-status - Check migration status"
	@echo "  db-create    - Create database"
	@echo "  db-drop      - Drop database"
	@echo "  db-reset     - Reset database and run migrations"
	@echo ""
	@echo "Production:"
	@echo "  prod-build   - Build for production"
	@echo "  prod-start   - Start production with Docker Compose"
	@echo "  prod-stop    - Stop production services"
	@echo "  prod-logs    - View production logs"
	@echo "  prod-migrate - Run migrations in production"
	@echo "  prod-backup  - Create database backup"
	@echo "  prod-status  - Check production service status"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker services"
	@echo "  docker-logs  - View Docker logs"
	@echo "  docker-clean - Clean up Docker resources"
	@echo ""
	@echo "Utilities:"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help"