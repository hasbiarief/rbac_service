package swagger

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHandler_ServeSpec_CacheHit(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory for test
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "swagger.json")

	// Create test swagger.json
	testSpec := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0"}}`
	err := os.WriteFile(specPath, []byte(testSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Create handler with caching enabled
	config := &Config{
		EnableCache: true,
		CacheTTL:    5 * time.Minute,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// First request - should read from file and cache
	req1 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}

	if w1.Body.String() != testSpec {
		t.Errorf("Expected spec content, got: %s", w1.Body.String())
	}

	// Verify cache was populated
	handler.cache.mu.RLock()
	cachePopulated := handler.cache.spec != nil
	handler.cache.mu.RUnlock()

	if !cachePopulated {
		t.Error("Expected cache to be populated after first request")
	}

	// Modify file (simulate change)
	newSpec := `{"openapi":"3.0.0","info":{"title":"Modified API","version":"2.0"}}`
	err = os.WriteFile(specPath, []byte(newSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test spec: %v", err)
	}

	// Second request - should serve from cache (not see the change)
	req2 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	// Should still return original spec (from cache)
	if w2.Body.String() != testSpec {
		t.Errorf("Expected cached spec, got: %s", w2.Body.String())
	}
}

func TestHandler_ServeSpec_CacheMiss(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory for test
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "swagger.json")

	// Create test swagger.json
	testSpec := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0"}}`
	err := os.WriteFile(specPath, []byte(testSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Create handler with very short TTL
	config := &Config{
		EnableCache: true,
		CacheTTL:    100 * time.Millisecond,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// First request - populate cache
	req1 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Modify file
	newSpec := `{"openapi":"3.0.0","info":{"title":"Modified API","version":"2.0"}}`
	err = os.WriteFile(specPath, []byte(newSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test spec: %v", err)
	}

	// Second request - cache expired, should read new file
	req2 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	// Should return modified spec (cache expired)
	if w2.Body.String() != newSpec {
		t.Errorf("Expected new spec after cache expiry, got: %s", w2.Body.String())
	}
}

func TestHandler_ServeSpec_CacheDisabled(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory for test
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "swagger.json")

	// Create test swagger.json
	testSpec := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0"}}`
	err := os.WriteFile(specPath, []byte(testSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Create handler with caching disabled
	config := &Config{
		EnableCache: false,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// First request
	req1 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}

	// Verify cache was NOT populated
	handler.cache.mu.RLock()
	cachePopulated := handler.cache.spec != nil
	handler.cache.mu.RUnlock()

	if cachePopulated {
		t.Error("Expected cache to remain empty when caching disabled")
	}

	// Modify file
	newSpec := `{"openapi":"3.0.0","info":{"title":"Modified API","version":"2.0"}}`
	err = os.WriteFile(specPath, []byte(newSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test spec: %v", err)
	}

	// Second request - should read new file immediately
	req2 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	// Should return modified spec (no caching)
	if w2.Body.String() != newSpec {
		t.Errorf("Expected new spec immediately when cache disabled, got: %s", w2.Body.String())
	}
}

func TestHandler_ServeSpec_FileNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory (but no swagger.json)
	tmpDir := t.TempDir()

	// Create handler
	config := &Config{
		EnableCache: true,
		CacheTTL:    5 * time.Minute,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// Request when file doesn't exist
	req := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// Verify error message
	if !contains(w.Body.String(), "not found") {
		t.Errorf("Expected error message about file not found, got: %s", w.Body.String())
	}
}

func TestHandler_InvalidateCache(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory for test
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "swagger.json")

	// Create test swagger.json
	testSpec := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0"}}`
	err := os.WriteFile(specPath, []byte(testSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Create handler with caching enabled
	config := &Config{
		EnableCache: true,
		CacheTTL:    5 * time.Minute,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// First request - populate cache
	req1 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w1.Code)
	}

	// Verify cache was populated
	handler.cache.mu.RLock()
	cachePopulated := handler.cache.spec != nil
	handler.cache.mu.RUnlock()

	if !cachePopulated {
		t.Error("Expected cache to be populated")
	}

	// Invalidate cache
	handler.InvalidateCache()

	// Verify cache was cleared
	handler.cache.mu.RLock()
	cacheCleared := handler.cache.spec == nil
	handler.cache.mu.RUnlock()

	if !cacheCleared {
		t.Error("Expected cache to be cleared after invalidation")
	}

	// Modify file
	newSpec := `{"openapi":"3.0.0","info":{"title":"Modified API","version":"2.0"}}`
	err = os.WriteFile(specPath, []byte(newSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test spec: %v", err)
	}

	// Second request - should read new file (cache was invalidated)
	req2 := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	// Should return modified spec (cache was invalidated)
	if w2.Body.String() != newSpec {
		t.Errorf("Expected new spec after cache invalidation, got: %s", w2.Body.String())
	}
}

func TestHandler_ConcurrentAccess(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Create temporary directory for test
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "swagger.json")

	// Create test swagger.json
	testSpec := `{"openapi":"3.0.0","info":{"title":"Test API","version":"1.0"}}`
	err := os.WriteFile(specPath, []byte(testSpec), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Create handler with caching enabled
	config := &Config{
		EnableCache: true,
		CacheTTL:    5 * time.Minute,
		OutputDir:   tmpDir,
		SpecPath:    "/api/swagger.json",
	}
	handler := NewHandler(config)

	// Setup router
	router := gin.New()
	router.GET(config.SpecPath, handler.serveSpec)

	// Concurrent requests
	const numRequests = 100
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, config.SpecPath, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			if w.Body.String() != testSpec {
				t.Errorf("Expected spec content, got: %s", w.Body.String())
			}

			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}

	// Verify cache is still valid
	handler.cache.mu.RLock()
	cacheValid := handler.cache.spec != nil
	handler.cache.mu.RUnlock()

	if !cacheValid {
		t.Error("Expected cache to remain valid after concurrent access")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
