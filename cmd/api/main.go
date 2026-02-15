package main

import (
	"gin-scalable-api/config"
	"gin-scalable-api/internal/app"
	"log"

	_ "gin-scalable-api/docs" // Import generated docs
)

// @title           Huminor Console API
// @version         1.0
// @description     Complete API for ERP with RBAC and API Documentation System. This API provides comprehensive endpoints for managing users, companies, branches, roles, modules, subscriptions, units, applications, and authentication with JWT-based security.
// @description
// @description     ## Environments
// @description     - **Development**: http://localhost:8081
// @description     - **Staging**: https://staging-api.huminor.com
// @description     - **Production**: https://api.huminor.com
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name 🔐 Authentication
// @tag.description Authentication endpoints - Login, logout, token refresh
// @tag.name Users
// @tag.description User management endpoints
// @tag.name Companies
// @tag.description Company management endpoints
// @tag.name Branches
// @tag.description Branch management endpoints with hierarchy support
// @tag.name Roles
// @tag.description Role management endpoints
// @tag.name Role Management
// @tag.description Advanced role management and assignments
// @tag.name Modules
// @tag.description Module/menu management endpoints
// @tag.name Subscriptions
// @tag.description Subscription management endpoints
// @tag.name Subscription Plans
// @tag.description Subscription plan management endpoints
// @tag.name Units
// @tag.description Unit management endpoints
// @tag.name Applications
// @tag.description Application management endpoints
// @tag.name Audit
// @tag.description Audit log endpoints
// @tag.name System
// @tag.description System health and status endpoints

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @Summary Health check
// @Description Check if the API is running
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize server with module-based structure
	srv := app.NewServer(cfg)

	// Initialize all components (database, repositories, services, handlers, routes)
	if err := srv.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
