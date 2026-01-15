package app

import (
	"database/sql"
	"gin-scalable-api/config"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/database"
	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/token"
	"log"

	// Module imports
	auditModule "gin-scalable-api/internal/modules/audit"
	authModule "gin-scalable-api/internal/modules/auth"
	branchModule "gin-scalable-api/internal/modules/branch"
	companyModule "gin-scalable-api/internal/modules/company"
	moduleModule "gin-scalable-api/internal/modules/module"
	roleModule "gin-scalable-api/internal/modules/role"
	subscriptionModule "gin-scalable-api/internal/modules/subscription"
	unitModule "gin-scalable-api/internal/modules/unit"
	userModule "gin-scalable-api/internal/modules/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Initialize() error {
	// Connect to database
	dbConfig := database.Config{
		Host:     s.config.Database.Host,
		Port:     s.config.Database.Port,
		User:     s.config.Database.User,
		Password: s.config.Database.Password,
		DBName:   s.config.Database.Name,
		SSLMode:  s.config.Database.SSLMode,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		return err
	}

	// Initialize Redis
	redis := config.InitRedis(s.config)

	// Initialize NEW module handlers
	newModuleHandlers := s.initializeNewModuleHandlers(redis, db.DB)

	// Initialize Gin router
	s.router = gin.Default()

	// Add CORS middleware
	s.router.Use(middleware.CORSMiddleware())

	// Setup NEW module routes
	SetupNewModuleRoutes(s.router, newModuleHandlers, s.config.JWT.Secret, redis)

	return nil
}

func (s *Server) initializeNewModuleHandlers(redis *redis.Client, db *sql.DB) *NewModuleHandlers {
	// Initialize token service
	tokenService := token.NewSimpleTokenService(redis)

	// Initialize RBAC service
	rbacService := rbac.NewRBACService(db)

	// Initialize module repositories (using module implementations)
	userRepo := userModule.NewUserRepository(db)
	roleRepo := roleModule.NewRoleRepository(db)
	companyRepo := companyModule.NewCompanyRepository(db)
	branchRepo := branchModule.NewBranchRepository(db)
	moduleRepo := moduleModule.NewModuleRepository(db)
	subscriptionRepo := subscriptionModule.NewRepository(db)
	auditRepo := auditModule.NewRepository(db)
	unitRepo := unitModule.NewRepository(db)

	// Initialize module services
	authRepo := authModule.NewRepository(db)
	authService := authModule.NewService(authRepo, tokenService, s.config.JWT.Secret)
	userService := userModule.NewService(userRepo, rbacService)
	roleService := roleModule.NewService(roleRepo)
	companyService := companyModule.NewService(companyRepo)
	branchService := branchModule.NewService(branchRepo)
	moduleService := moduleModule.NewService(moduleRepo)
	unitService := unitModule.NewService(unitRepo)
	subscriptionService := subscriptionModule.NewService(subscriptionRepo)
	auditService := auditModule.NewService(auditRepo)

	// Initialize module handlers
	return &NewModuleHandlers{
		Auth:         authModule.NewHandler(authService),
		User:         userModule.NewHandler(userService, userRepo),
		Role:         roleModule.NewHandler(roleService),
		Company:      companyModule.NewHandler(companyService),
		Branch:       branchModule.NewHandler(branchService),
		Module:       moduleModule.NewHandler(moduleService),
		Unit:         unitModule.NewHandler(unitService),
		Subscription: subscriptionModule.NewHandler(subscriptionService),
		Audit:        auditModule.NewHandler(auditService),
	}
}

func (s *Server) Run() error {
	log.Printf("Server starting on port %s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}

// NewModuleHandlers struct for new module-based handlers
type NewModuleHandlers struct {
	Auth         *authModule.Handler
	User         *userModule.Handler
	Role         *roleModule.Handler
	Company      *companyModule.Handler
	Branch       *branchModule.Handler
	Module       *moduleModule.Handler
	Unit         *unitModule.Handler
	Subscription *subscriptionModule.Handler
	Audit        *auditModule.Handler
}
