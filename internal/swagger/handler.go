package swagger

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler handles Swagger UI and documentation serving
type Handler struct {
	config *Config
	cache  *specCache
}

// specCache caches the generated specification
type specCache struct {
	mu         sync.RWMutex
	spec       []byte
	lastUpdate time.Time
	ttl        time.Duration
}

// NewHandler creates a new Swagger handler
func NewHandler(config *Config) *Handler {
	return &Handler{
		config: config,
		cache: &specCache{
			ttl: config.CacheTTL,
		},
	}
}

// RegisterRoutes registers Swagger UI and spec routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	if !h.config.EnableUI {
		return
	}

	// Swagger UI with custom configuration
	router.GET(h.config.UIPath+"/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL(h.config.SpecPath),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DocExpansion("list"),
		ginSwagger.PersistAuthorization(true),
	))

	// Swagger JSON spec
	router.GET(h.config.SpecPath, h.serveSpec)
}

// serveSpec serves the swagger.json file with caching
func (h *Handler) serveSpec(c *gin.Context) {
	// Check cache if enabled
	if h.config.EnableCache {
		h.cache.mu.RLock()
		if h.cache.spec != nil && time.Since(h.cache.lastUpdate) < h.cache.ttl {
			spec := h.cache.spec
			h.cache.mu.RUnlock()
			c.Data(http.StatusOK, "application/json", spec)
			return
		}
		h.cache.mu.RUnlock()
	}

	// Read spec file
	specPath := filepath.Join(h.config.OutputDir, "swagger.json")
	spec, err := os.ReadFile(specPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Swagger specification not found. Please generate it first.",
		})
		return
	}

	// Update cache
	if h.config.EnableCache {
		h.cache.mu.Lock()
		h.cache.spec = spec
		h.cache.lastUpdate = time.Now()
		h.cache.mu.Unlock()
	}

	c.Data(http.StatusOK, "application/json", spec)
}

// InvalidateCache clears the cached specification
func (h *Handler) InvalidateCache() {
	h.cache.mu.Lock()
	defer h.cache.mu.Unlock()
	h.cache.spec = nil
	h.cache.lastUpdate = time.Time{}
}
