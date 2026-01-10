package server

import (
	"fmt"
	"gin-scalable-api/config"
	"gin-scalable-api/internal/handlers"
	"gin-scalable-api/internal/interfaces"
	"gin-scalable-api/internal/repository"
	"gin-scalable-api/internal/routes"
	"gin-scalable-api/internal/service"
	"gin-scalable-api/middleware"
	"gin-scalable-api/pkg/database"
	"gin-scalable-api/pkg/jobs"
	"gin-scalable-api/pkg/rbac"
	"gin-scalable-api/pkg/token"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	cleanupJob *jobs.CleanupJob
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

	// Initialize repositories
	repositories := s.initializeRepositories(db)

	// Initialize services
	services := s.initializeServices(repositories, redis)

	// Initialize handlers
	handlers := s.initializeHandlers(services, repositories)

	// Initialize Gin router
	s.router = gin.Default()

	// Add CORS middleware
	s.router.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(s.router, handlers, s.config.JWT.Secret, redis)

	// Initialize and start cleanup job (runs every 30 minutes)
	// Type assert to concrete type for cleanup job
	authService, ok := services.Auth.(*service.AuthService)
	if !ok {
		return fmt.Errorf("failed to type assert auth service")
	}
	s.cleanupJob = jobs.NewCleanupJob(authService, 30*time.Minute)
	go s.cleanupJob.Start()

	return nil
}

func (s *Server) initializeRepositories(db *database.DB) *Repositories {
	// Create concrete repositories
	userRepo := repository.NewUserRepository(db.DB)
	moduleRepo := repository.NewModuleRepository(db.DB)
	companyRepo := repository.NewCompanyRepository(db.DB)
	roleRepo := repository.NewRoleRepository(db.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(db.DB)
	auditRepo := repository.NewAuditRepository(db.DB)
	branchRepo := repository.NewBranchRepository(db.DB)

	return &Repositories{
		User:         userRepo,
		Module:       moduleRepo,
		Company:      companyRepo,
		Role:         roleRepo,
		Subscription: subscriptionRepo,
		Audit:        auditRepo,
		Branch:       branchRepo,
	}
}

func (s *Server) initializeServices(repos *Repositories, redis *redis.Client) *Services {
	tokenService := token.NewSimpleTokenService(redis)
	rbacService := rbac.NewRBACService(repos.User.GetDB())

	return &Services{
		Auth:         service.NewAuthService(repos.User, tokenService, s.config.JWT.Secret),
		Module:       service.NewModuleService(repos.Module, repos.User, rbacService),
		Company:      service.NewCompanyService(repos.Company),
		Role:         service.NewRoleService(repos.Role, repos.User),
		User:         service.NewUserService(repos.User, rbacService),
		Subscription: service.NewSubscriptionService(repos.Subscription),
		Audit:        service.NewAuditService(repos.Audit),
		Branch:       service.NewBranchService(repos.Branch),
	}
}

func (s *Server) initializeHandlers(services *Services, repos *Repositories) *routes.Handlers {
	return &routes.Handlers{
		Auth:         handlers.NewAuthHandler(services.Auth),
		Module:       handlers.NewModuleHandler(services.Module),
		User:         handlers.NewUserHandler(services.User, services.Module, repos.User),
		Company:      handlers.NewCompanyHandler(services.Company),
		Role:         handlers.NewRoleHandler(services.Role),
		Subscription: handlers.NewSubscriptionHandler(services.Subscription),
		Audit:        handlers.NewAuditHandler(services.Audit),
		Branch:       handlers.NewBranchHandler(services.Branch),
	}
}

func (s *Server) Run() error {
	log.Printf("Server starting on port %s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}

// Repositories struct to group all repositories
type Repositories struct {
	User         *repository.UserRepository
	Module       *repository.ModuleRepository
	Company      *repository.CompanyRepository
	Role         *repository.RoleRepository
	Subscription *repository.SubscriptionRepository
	Audit        *repository.AuditRepository
	Branch       *repository.BranchRepository
}

// Services struct to group all services
type Services struct {
	Auth         interfaces.AuthServiceInterface
	Module       interfaces.ModuleServiceInterface
	Company      interfaces.CompanyServiceInterface
	Role         interfaces.RoleServiceInterface
	User         interfaces.UserServiceInterface
	Subscription interfaces.SubscriptionServiceInterface
	Audit        interfaces.AuditServiceInterface
	Branch       interfaces.BranchServiceInterface
}
