package export

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestCacheKeyGeneration(t *testing.T) {
	// Create a mock cache (without Redis for unit testing)
	config := DefaultCacheConfig()
	cache := &ExportCache{
		redis:  nil, // No Redis for unit test
		config: config,
	}

	// Test cache key generation
	collectionID := int64(123)
	format := FormatPostman
	options := DefaultExportOptions()
	options.Format = format

	key1 := cache.GenerateCacheKey(collectionID, format, options)
	key2 := cache.GenerateCacheKey(collectionID, format, options)

	// Same parameters should generate same key
	if key1 != key2 {
		t.Errorf("Expected same cache key for same parameters, got %s and %s", key1, key2)
	}

	// Different collection should generate different key
	key3 := cache.GenerateCacheKey(456, format, options)
	if key1 == key3 {
		t.Errorf("Expected different cache key for different collection, got same key: %s", key1)
	}

	// Different format should generate different key
	key4 := cache.GenerateCacheKey(collectionID, FormatOpenAPI, options)
	if key1 == key4 {
		t.Errorf("Expected different cache key for different format, got same key: %s", key1)
	}

	// Key should contain collection ID and format
	if !contains(key1, "123") {
		t.Errorf("Expected cache key to contain collection ID, got: %s", key1)
	}

	if !contains(key1, string(format)) {
		t.Errorf("Expected cache key to contain format, got: %s", key1)
	}
}

func TestCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	// Test default values
	if !config.Enabled {
		t.Error("Expected cache to be enabled by default")
	}

	if config.TTL != 30*time.Minute {
		t.Errorf("Expected default TTL to be 30 minutes, got %v", config.TTL)
	}

	if config.KeyPrefix != "apidoc:export:" {
		t.Errorf("Expected default key prefix to be 'apidoc:export:', got %s", config.KeyPrefix)
	}

	if config.MaxSize != 10*1024*1024 {
		t.Errorf("Expected default max size to be 10MB, got %d", config.MaxSize)
	}

	if !config.Compress {
		t.Error("Expected compression to be enabled by default")
	}
}

func TestCacheDisabled(t *testing.T) {
	config := DefaultCacheConfig()
	config.Enabled = false

	cache := NewExportCache(nil, config)

	// Test that disabled cache returns empty keys
	key := cache.GenerateCacheKey(123, FormatPostman, DefaultExportOptions())
	if key != "" {
		t.Errorf("Expected empty key when cache disabled, got: %s", key)
	}

	// Test that disabled cache operations return nil/no-op
	result, err := cache.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error from disabled cache Get, got: %v", err)
	}
	if result != nil {
		t.Error("Expected nil result from disabled cache Get")
	}

	err = cache.Set("test-key", &ExportResult{})
	if err != nil {
		t.Errorf("Expected no error from disabled cache Set, got: %v", err)
	}

	err = cache.Delete("test-key")
	if err != nil {
		t.Errorf("Expected no error from disabled cache Delete, got: %v", err)
	}
}

// Integration test with Redis (requires Redis to be running)
func TestCacheWithRedis(t *testing.T) {
	// Skip if Redis is not available
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use test database
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	// Clean up test data
	defer func() {
		rdb.FlushDB(ctx)
		rdb.Close()
	}()

	config := DefaultCacheConfig()
	config.TTL = 1 * time.Second // Short TTL for testing
	cache := NewExportCache(rdb, config)

	// Test cache operations
	key := "test:export:123"
	testResult := &ExportResult{
		Content:     "test content",
		ContentType: "application/json",
		Filename:    "test.json",
		Size:        12,
		GeneratedAt: time.Now(),
	}

	// Test Set
	err := cache.Set(key, testResult)
	if err != nil {
		t.Errorf("Failed to set cache: %v", err)
	}

	// Test Get
	result, err := cache.Get(key)
	if err != nil {
		t.Errorf("Failed to get cache: %v", err)
	}

	if result == nil {
		t.Error("Expected cached result, got nil")
	} else {
		if result.Content != testResult.Content {
			t.Errorf("Expected content %s, got %s", testResult.Content, result.Content)
		}
		if result.Filename != testResult.Filename {
			t.Errorf("Expected filename %s, got %s", testResult.Filename, result.Filename)
		}
	}

	// Test TTL expiration
	time.Sleep(2 * time.Second)
	result, err = cache.Get(key)
	if err != nil {
		t.Errorf("Failed to get expired cache: %v", err)
	}
	if result != nil {
		t.Error("Expected nil result after TTL expiration, got result")
	}

	// Test Delete
	cache.Set(key, testResult)
	err = cache.Delete(key)
	if err != nil {
		t.Errorf("Failed to delete cache: %v", err)
	}

	result, err = cache.Get(key)
	if err != nil {
		t.Errorf("Failed to get deleted cache: %v", err)
	}
	if result != nil {
		t.Error("Expected nil result after deletion, got result")
	}
}

func TestCacheInvalidation(t *testing.T) {
	// Skip if Redis is not available
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use test database
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	// Clean up test data
	defer func() {
		rdb.FlushDB(ctx)
		rdb.Close()
	}()

	cache := NewExportCache(rdb, DefaultCacheConfig())

	// Set up test data
	collectionID := int64(123)
	testResult := &ExportResult{
		Content:     "test content",
		ContentType: "application/json",
		Filename:    "test.json",
		Size:        12,
		GeneratedAt: time.Now(),
	}

	// Create cache entries for different formats
	formats := []ExportFormat{FormatPostman, FormatOpenAPI, FormatInsomnia}
	keys := make([]string, len(formats))

	for i, format := range formats {
		key := cache.GenerateCacheKey(collectionID, format, DefaultExportOptions())
		keys[i] = key
		cache.Set(key, testResult)
	}

	// Verify all entries exist
	for _, key := range keys {
		result, err := cache.Get(key)
		if err != nil || result == nil {
			t.Errorf("Expected cached result for key %s", key)
		}
	}

	// Test collection invalidation
	err := cache.InvalidateCollection(collectionID)
	if err != nil {
		t.Errorf("Failed to invalidate collection cache: %v", err)
	}

	// Verify all entries are gone
	for _, key := range keys {
		result, err := cache.Get(key)
		if err != nil {
			t.Errorf("Error getting invalidated cache: %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result after invalidation for key %s", key)
		}
	}
}

func TestCacheStats(t *testing.T) {
	// Skip if Redis is not available
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use test database
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping integration test")
		return
	}

	// Clean up test data
	defer func() {
		rdb.FlushDB(ctx)
		rdb.Close()
	}()

	cache := NewExportCache(rdb, DefaultCacheConfig())

	// Get initial stats
	stats, err := cache.GetCacheStats()
	if err != nil {
		t.Errorf("Failed to get cache stats: %v", err)
	}

	if !stats.Enabled {
		t.Error("Expected cache to be enabled")
	}

	if stats.KeyPrefix != "apidoc:export:" {
		t.Errorf("Expected key prefix 'apidoc:export:', got %s", stats.KeyPrefix)
	}

	// Add some cache entries
	testResult := &ExportResult{
		Content:     "test content",
		ContentType: "application/json",
		Filename:    "test.json",
		Size:        12,
		GeneratedAt: time.Now(),
	}

	for i := 0; i < 3; i++ {
		key := cache.GenerateCacheKey(int64(i), FormatPostman, DefaultExportOptions())
		cache.Set(key, testResult)
	}

	// Get updated stats
	stats, err = cache.GetCacheStats()
	if err != nil {
		t.Errorf("Failed to get updated cache stats: %v", err)
	}

	if stats.KeyCount < 3 {
		t.Errorf("Expected at least 3 keys, got %d", stats.KeyCount)
	}
}
