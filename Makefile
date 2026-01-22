# Makefile for ERP RBAC System without GORM

.PHONY: build run migrate-up migrate-status clean test newmodule removemodule listmodules db-dump db-seed

# Build the application
build:
	go build -o bin/server cmd/api/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run the application
run: build
	./bin/server

# Generate new module with boilerplate code
newmodule:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make newmodule name=<module_name>"; \
		echo "📝 Example: make newmodule name=dashboard"; \
		exit 1; \
	fi
	@if [ -d "internal/modules/$(name)" ]; then \
		echo "❌ Module '$(name)' already exists!"; \
		exit 1; \
	fi
	@echo "🚀 Generating module '$(name)'..."
	@mkdir -p internal/modules/$(name)
	@$(MAKE) -s _generate-files name=$(name)
	@echo "✅ Module '$(name)' generated successfully!"
	@echo ""
	@echo "📝 Next steps:"
	@echo "   1. Register module in internal/app/routes.go"
	@echo "   2. Register module in internal/app/server.go"
	@echo "   3. Create database migration if needed"
	@echo "   4. Implement business logic"
	@echo "   5. Test with: go build ./cmd/api"

# Remove/rollback generated module
removemodule:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make removemodule name=<module_name>"; \
		echo "📝 Example: make removemodule name=dashboard"; \
		exit 1; \
	fi
	@if [ ! -d "internal/modules/$(name)" ]; then \
		echo "❌ Module '$(name)' does not exist!"; \
		exit 1; \
	fi
	@echo "🗑️  Removing module '$(name)'..."
	@echo "📁 Files to be deleted:"
	@ls -la internal/modules/$(name)/
	@echo ""
	@echo "⚠️  This action cannot be undone!"
	@read -p "Are you sure you want to delete module '$(name)'? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		rm -rf internal/modules/$(name); \
		echo "✅ Module '$(name)' removed successfully!"; \
		echo ""; \
		echo "📝 Don't forget to:"; \
		echo "   1. Remove module registration from internal/app/routes.go"; \
		echo "   2. Remove module registration from internal/app/server.go"; \
		echo "   3. Remove related database migrations if any"; \
		echo "   4. Test with: go build ./cmd/api"; \
	else \
		echo "❌ Operation cancelled."; \
	fi

# Force remove module without confirmation (for scripts/automation)
removemodule-force:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make removemodule-force name=<module_name>"; \
		exit 1; \
	fi
	@if [ ! -d "internal/modules/$(name)" ]; then \
		echo "❌ Module '$(name)' does not exist!"; \
		exit 1; \
	fi
	@echo "🗑️  Force removing module '$(name)'..."
	@rm -rf internal/modules/$(name)
	@echo "✅ Module '$(name)' removed successfully!"

# List all existing modules
listmodules:
	@echo "📦 Existing modules:"
	@echo ""
	@if [ -d "internal/modules" ]; then \
		for dir in internal/modules/*/; do \
			if [ -d "$$dir" ]; then \
				module_name=$$(basename "$$dir"); \
				file_count=$$(find "$$dir" -name "*.go" | wc -l); \
				echo "  📁 $$module_name ($$file_count files)"; \
				find "$$dir" -name "*.go" -exec basename {} \; | sed 's/^/    - /'; \
				echo ""; \
			fi; \
		done; \
	else \
		echo "  No modules directory found."; \
	fi
	@echo "💡 Usage:"
	@echo "  make newmodule name=<name>     - Create new module"
	@echo "  make removemodule name=<name>  - Remove existing module"

# Generate all 5 files with templates
_generate-files:
	@$(MAKE) -s _gen-model name=$(name)
	@$(MAKE) -s _gen-dto name=$(name)
	@$(MAKE) -s _gen-repository name=$(name)
	@$(MAKE) -s _gen-service name=$(name)
	@$(MAKE) -s _gen-route name=$(name)

# Generate model.go
_gen-model:
	@echo "package $(name)" > internal/modules/$(name)/model.go
	@echo "" >> internal/modules/$(name)/model.go
	@echo "import \"time\"" >> internal/modules/$(name)/model.go
	@echo "" >> internal/modules/$(name)/model.go
	@echo "type $(name) struct {" >> internal/modules/$(name)/model.go
	@echo "	ID        int64     \`json:\"id\" db:\"id\"\`" >> internal/modules/$(name)/model.go
	@echo "	Name      string    \`json:\"name\" db:\"name\"\`" >> internal/modules/$(name)/model.go
	@echo "	IsActive  bool      \`json:\"is_active\" db:\"is_active\"\`" >> internal/modules/$(name)/model.go
	@echo "	CreatedAt time.Time \`json:\"created_at\" db:\"created_at\"\`" >> internal/modules/$(name)/model.go
	@echo "	UpdatedAt time.Time \`json:\"updated_at\" db:\"updated_at\"\`" >> internal/modules/$(name)/model.go
	@echo "}" >> internal/modules/$(name)/model.go
	@echo "" >> internal/modules/$(name)/model.go
	@echo "func ($(name)) TableName() string {" >> internal/modules/$(name)/model.go
	@echo "	return \"$(name)s\"" >> internal/modules/$(name)/model.go
	@echo "}" >> internal/modules/$(name)/model.go

# Generate dto.go
_gen-dto:
	@echo "package $(name)" > internal/modules/$(name)/dto.go
	@echo "" >> internal/modules/$(name)/dto.go
	@echo "// Request DTOs" >> internal/modules/$(name)/dto.go
	@echo "type create$(name)request struct {" >> internal/modules/$(name)/dto.go
	@echo "	Name string \`json:\"name\" validate:\"required,min=2,max=100\"\`" >> internal/modules/$(name)/dto.go
	@echo "}" >> internal/modules/$(name)/dto.go
	@echo "" >> internal/modules/$(name)/dto.go
	@echo "type update$(name)request struct {" >> internal/modules/$(name)/dto.go
	@echo "	Name     *string \`json:\"name,omitempty\" validate:\"omitempty,min=2,max=100\"\`" >> internal/modules/$(name)/dto.go
	@echo "	IsActive *bool   \`json:\"is_active,omitempty\"\`" >> internal/modules/$(name)/dto.go
	@echo "}" >> internal/modules/$(name)/dto.go
	@echo "" >> internal/modules/$(name)/dto.go
	@echo "// Response DTOs" >> internal/modules/$(name)/dto.go
	@echo "type $(name)response struct {" >> internal/modules/$(name)/dto.go
	@echo "	ID        int64  \`json:\"id\"\`" >> internal/modules/$(name)/dto.go
	@echo "	Name      string \`json:\"name\"\`" >> internal/modules/$(name)/dto.go
	@echo "	IsActive  bool   \`json:\"is_active\"\`" >> internal/modules/$(name)/dto.go
	@echo "	CreatedAt string \`json:\"created_at\"\`" >> internal/modules/$(name)/dto.go
	@echo "	UpdatedAt string \`json:\"updated_at\"\`" >> internal/modules/$(name)/dto.go
	@echo "}" >> internal/modules/$(name)/dto.go
	@echo "" >> internal/modules/$(name)/dto.go
	@echo "type $(name)listresponse struct {" >> internal/modules/$(name)/dto.go
	@echo "	Data    []*$(name)response \`json:\"data\"\`" >> internal/modules/$(name)/dto.go
	@echo "	Total   int64              \`json:\"total\"\`" >> internal/modules/$(name)/dto.go
	@echo "	Limit   int                \`json:\"limit\"\`" >> internal/modules/$(name)/dto.go
	@echo "	Offset  int                \`json:\"offset\"\`" >> internal/modules/$(name)/dto.go
	@echo "	HasMore bool               \`json:\"has_more\"\`" >> internal/modules/$(name)/dto.go
	@echo "}" >> internal/modules/$(name)/dto.go

# Generate repository.go
_gen-repository:
	@echo "package $(name)" > internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "import (" >> internal/modules/$(name)/repository.go
	@echo "	\"database/sql\"" >> internal/modules/$(name)/repository.go
	@echo "	\"fmt\"" >> internal/modules/$(name)/repository.go
	@echo "	\"gin-scalable-api/pkg/model\"" >> internal/modules/$(name)/repository.go
	@echo ")" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "type repository struct {" >> internal/modules/$(name)/repository.go
	@echo "	*model.Repository" >> internal/modules/$(name)/repository.go
	@echo "	db *sql.DB" >> internal/modules/$(name)/repository.go
	@echo "}" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "func newrepository(db *sql.DB) *repository {" >> internal/modules/$(name)/repository.go
	@echo "	return &repository{" >> internal/modules/$(name)/repository.go
	@echo "		Repository: model.NewRepository(db)," >> internal/modules/$(name)/repository.go
	@echo "		db:         db," >> internal/modules/$(name)/repository.go
	@echo "	}" >> internal/modules/$(name)/repository.go
	@echo "}" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "func (r *repository) create(item *$(name)) error {" >> internal/modules/$(name)/repository.go
	@echo "	query := \`" >> internal/modules/$(name)/repository.go
	@echo "		INSERT INTO $(name)s (name, is_active, created_at, updated_at)" >> internal/modules/$(name)/repository.go
	@echo "		VALUES ($$1, $$2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)" >> internal/modules/$(name)/repository.go
	@echo "		RETURNING id, created_at, updated_at" >> internal/modules/$(name)/repository.go
	@echo "	\`" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "	err := r.db.QueryRow(query, item.Name, item.IsActive).Scan(" >> internal/modules/$(name)/repository.go
	@echo "		&item.ID, &item.CreatedAt, &item.UpdatedAt," >> internal/modules/$(name)/repository.go
	@echo "	)" >> internal/modules/$(name)/repository.go
	@echo "	if err != nil {" >> internal/modules/$(name)/repository.go
	@echo "		return fmt.Errorf(\"failed to create $(name): %w\", err)" >> internal/modules/$(name)/repository.go
	@echo "	}" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "	return nil" >> internal/modules/$(name)/repository.go
	@echo "}" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "func (r *repository) getbyid(id int64) (*$(name), error) {" >> internal/modules/$(name)/repository.go
	@echo "	item := &$(name){}" >> internal/modules/$(name)/repository.go
	@echo "	query := \`" >> internal/modules/$(name)/repository.go
	@echo "		SELECT id, name, is_active, created_at, updated_at" >> internal/modules/$(name)/repository.go
	@echo "		FROM $(name)s" >> internal/modules/$(name)/repository.go
	@echo "		WHERE id = $$1 AND deleted_at IS NULL" >> internal/modules/$(name)/repository.go
	@echo "	\`" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "	err := r.db.QueryRow(query, id).Scan(" >> internal/modules/$(name)/repository.go
	@echo "		&item.ID, &item.Name, &item.IsActive," >> internal/modules/$(name)/repository.go
	@echo "		&item.CreatedAt, &item.UpdatedAt," >> internal/modules/$(name)/repository.go
	@echo "	)" >> internal/modules/$(name)/repository.go
	@echo "	if err != nil {" >> internal/modules/$(name)/repository.go
	@echo "		if err == sql.ErrNoRows {" >> internal/modules/$(name)/repository.go
	@echo "			return nil, fmt.Errorf(\"$(name) not found\")" >> internal/modules/$(name)/repository.go
	@echo "		}" >> internal/modules/$(name)/repository.go
	@echo "		return nil, fmt.Errorf(\"failed to get $(name): %w\", err)" >> internal/modules/$(name)/repository.go
	@echo "	}" >> internal/modules/$(name)/repository.go
	@echo "" >> internal/modules/$(name)/repository.go
	@echo "	return item, nil" >> internal/modules/$(name)/repository.go
	@echo "}" >> internal/modules/$(name)/repository.go

# Generate service.go
_gen-service:
	@echo "package $(name)" > internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "import (" >> internal/modules/$(name)/service.go
	@echo "	\"fmt\"" >> internal/modules/$(name)/service.go
	@echo "	\"time\"" >> internal/modules/$(name)/service.go
	@echo ")" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "type service struct {" >> internal/modules/$(name)/service.go
	@echo "	repo *repository" >> internal/modules/$(name)/service.go
	@echo "}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "func newservice(repo *repository) *service {" >> internal/modules/$(name)/service.go
	@echo "	return &service{repo: repo}" >> internal/modules/$(name)/service.go
	@echo "}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "func (s *service) create$(name)(req *create$(name)request) (*$(name)response, error) {" >> internal/modules/$(name)/service.go
	@echo "	// Create $(name) model" >> internal/modules/$(name)/service.go
	@echo "	item := &$(name){" >> internal/modules/$(name)/service.go
	@echo "		Name:     req.Name," >> internal/modules/$(name)/service.go
	@echo "		IsActive: true," >> internal/modules/$(name)/service.go
	@echo "	}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "	// Save to database" >> internal/modules/$(name)/service.go
	@echo "	if err := s.repo.create(item); err != nil {" >> internal/modules/$(name)/service.go
	@echo "		return nil, fmt.Errorf(\"failed to create $(name): %w\", err)" >> internal/modules/$(name)/service.go
	@echo "	}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "	// Convert to response" >> internal/modules/$(name)/service.go
	@echo "	return &$(name)response{" >> internal/modules/$(name)/service.go
	@echo "		ID:        item.ID," >> internal/modules/$(name)/service.go
	@echo "		Name:      item.Name," >> internal/modules/$(name)/service.go
	@echo "		IsActive:  item.IsActive," >> internal/modules/$(name)/service.go
	@echo "		CreatedAt: item.CreatedAt.Format(time.RFC3339)," >> internal/modules/$(name)/service.go
	@echo "		UpdatedAt: item.UpdatedAt.Format(time.RFC3339)," >> internal/modules/$(name)/service.go
	@echo "	}, nil" >> internal/modules/$(name)/service.go
	@echo "}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "func (s *service) get$(name)byid(id int64) (*$(name)response, error) {" >> internal/modules/$(name)/service.go
	@echo "	item, err := s.repo.getbyid(id)" >> internal/modules/$(name)/service.go
	@echo "	if err != nil {" >> internal/modules/$(name)/service.go
	@echo "		return nil, fmt.Errorf(\"failed to get $(name): %w\", err)" >> internal/modules/$(name)/service.go
	@echo "	}" >> internal/modules/$(name)/service.go
	@echo "" >> internal/modules/$(name)/service.go
	@echo "	return &$(name)response{" >> internal/modules/$(name)/service.go
	@echo "		ID:        item.ID," >> internal/modules/$(name)/service.go
	@echo "		Name:      item.Name," >> internal/modules/$(name)/service.go
	@echo "		IsActive:  item.IsActive," >> internal/modules/$(name)/service.go
	@echo "		CreatedAt: item.CreatedAt.Format(time.RFC3339)," >> internal/modules/$(name)/service.go
	@echo "		UpdatedAt: item.UpdatedAt.Format(time.RFC3339)," >> internal/modules/$(name)/service.go
	@echo "	}, nil" >> internal/modules/$(name)/service.go
	@echo "}" >> internal/modules/$(name)/service.go

# Generate route.go
_gen-route:
	@echo "package $(name)" > internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "import (" >> internal/modules/$(name)/route.go
	@echo "	\"gin-scalable-api/internal/constants\"" >> internal/modules/$(name)/route.go
	@echo "	\"gin-scalable-api/middleware\"" >> internal/modules/$(name)/route.go
	@echo "	\"gin-scalable-api/pkg/response\"" >> internal/modules/$(name)/route.go
	@echo "	\"net/http\"" >> internal/modules/$(name)/route.go
	@echo "	\"strconv\"" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	\"github.com/gin-gonic/gin\"" >> internal/modules/$(name)/route.go
	@echo ")" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "// Handler struct" >> internal/modules/$(name)/route.go
	@echo "type handler struct {" >> internal/modules/$(name)/route.go
	@echo "	service *service" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "func newhandler(service *service) *handler {" >> internal/modules/$(name)/route.go
	@echo "	return &handler{service: service}" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "// Handler methods" >> internal/modules/$(name)/route.go
	@echo "func (h *handler) get$(name)s(c *gin.Context) {" >> internal/modules/$(name)/route.go
	@echo "	// TODO: Implement pagination and filtering" >> internal/modules/$(name)/route.go
	@echo "	response.Success(c, http.StatusOK, \"$(name)s retrieved\", nil)" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "func (h *handler) get$(name)byid(c *gin.Context) {" >> internal/modules/$(name)/route.go
	@echo "	id, err := strconv.ParseInt(c.Param(\"id\"), 10, 64)" >> internal/modules/$(name)/route.go
	@echo "	if err != nil {" >> internal/modules/$(name)/route.go
	@echo "		response.Error(c, http.StatusBadRequest, \"Bad request\", \"Invalid ID\")" >> internal/modules/$(name)/route.go
	@echo "		return" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	result, err := h.service.get$(name)byid(id)" >> internal/modules/$(name)/route.go
	@echo "	if err != nil {" >> internal/modules/$(name)/route.go
	@echo "		response.Error(c, http.StatusNotFound, \"$(name) not found\", err.Error())" >> internal/modules/$(name)/route.go
	@echo "		return" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	response.Success(c, http.StatusOK, \"$(name) retrieved\", result)" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "func (h *handler) create$(name)(c *gin.Context) {" >> internal/modules/$(name)/route.go
	@echo "	validatedBody, exists := c.Get(\"validated_body\")" >> internal/modules/$(name)/route.go
	@echo "	if !exists {" >> internal/modules/$(name)/route.go
	@echo "		response.Error(c, http.StatusBadRequest, \"Bad request\", \"validation failed\")" >> internal/modules/$(name)/route.go
	@echo "		return" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	req, ok := validatedBody.(*create$(name)request)" >> internal/modules/$(name)/route.go
	@echo "	if !ok {" >> internal/modules/$(name)/route.go
	@echo "		response.Error(c, http.StatusBadRequest, \"Bad request\", \"invalid body structure\")" >> internal/modules/$(name)/route.go
	@echo "		return" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	result, err := h.service.create$(name)(req)" >> internal/modules/$(name)/route.go
	@echo "	if err != nil {" >> internal/modules/$(name)/route.go
	@echo "		response.ErrorWithAutoStatus(c, \"Failed to create $(name)\", err.Error())" >> internal/modules/$(name)/route.go
	@echo "		return" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "	response.Success(c, http.StatusCreated, constants.MsgCreated, result)" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "// Route registration" >> internal/modules/$(name)/route.go
	@echo "func RegisterRoutes(router *gin.RouterGroup, handler *handler) {" >> internal/modules/$(name)/route.go
	@echo "	$(name)s := router.Group(\"/$(name)s\")" >> internal/modules/$(name)/route.go
	@echo "	{" >> internal/modules/$(name)/route.go
	@echo "		// GET /api/v1/$(name)s - Get all $(name)s with pagination" >> internal/modules/$(name)/route.go
	@echo "		$(name)s.GET(\"\", handler.get$(name)s)" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "		// GET /api/v1/$(name)s/:id - Get $(name) by ID" >> internal/modules/$(name)/route.go
	@echo "		$(name)s.GET(\"/:id\", handler.get$(name)byid)" >> internal/modules/$(name)/route.go
	@echo "" >> internal/modules/$(name)/route.go
	@echo "		// POST /api/v1/$(name)s - Create new $(name)" >> internal/modules/$(name)/route.go
	@echo "		$(name)s.POST(\"\"," >> internal/modules/$(name)/route.go
	@echo "			middleware.ValidateRequest(middleware.ValidationRules{" >> internal/modules/$(name)/route.go
	@echo "				Body: &create$(name)request{}," >> internal/modules/$(name)/route.go
	@echo "			})," >> internal/modules/$(name)/route.go
	@echo "			handler.create$(name)," >> internal/modules/$(name)/route.go
	@echo "		)" >> internal/modules/$(name)/route.go
	@echo "	}" >> internal/modules/$(name)/route.go
	@echo "}" >> internal/modules/$(name)/route.go

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
	createdb $(DB_NAME)

db-drop:
	dropdb $(DB_NAME)

db-reset: db-drop db-create migrate-up

# Database dump and seeding
db-dump:
	@chmod +x scripts/db-dump.sh
	@./scripts/db-dump.sh

db-seed:
	@chmod +x scripts/db-seed.sh
	@./scripts/db-seed.sh

db-seed-fresh: db-drop db-create db-seed

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
	@echo "  newmodule    - Generate new module (Usage: make newmodule name=<module_name>)"
	@echo "  removemodule - Remove existing module (Usage: make removemodule name=<module_name>)"
	@echo "  listmodules  - List all existing modules"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-status - Check migration status"
	@echo "  db-create    - Create database"
	@echo "  db-drop      - Drop database"
	@echo "  db-reset     - Reset database and run migrations"
	@echo "  db-dump      - Create database dump and seeder files"
	@echo "  db-seed      - Seed database with template data"
	@echo "  db-seed-fresh - Drop, create, and seed database"
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