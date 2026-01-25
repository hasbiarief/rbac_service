package export

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheConfig contains cache configuration
type CacheConfig struct {
	Enabled   bool          `json:"enabled"`
	TTL       time.Duration `json:"ttl"`
	KeyPrefix string        `json:"key_prefix"`
	MaxSize   int64         `json:"max_size"` // Maximum cache entry size in bytes
	Compress  bool          `json:"compress"`
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:   true,
		TTL:       30 * time.Minute, // 30 minutes default TTL
		KeyPrefix: "apidoc:export:",
		MaxSize:   10 * 1024 * 1024, // 10MB max size per entry
		Compress:  true,
	}
}

// ExportCache handles caching of export results
type ExportCache struct {
	redis  *redis.Client
	config *CacheConfig
	ctx    context.Context
}

// NewExportCache creates a new export cache instance
func NewExportCache(redisClient *redis.Client, config *CacheConfig) *ExportCache {
	if config == nil {
		config = DefaultCacheConfig()
	}

	return &ExportCache{
		redis:  redisClient,
		config: config,
		ctx:    context.Background(),
	}
}

// CacheKey represents a cache key with metadata
type CacheKey struct {
	CollectionID  int64          `json:"collection_id"`
	Format        ExportFormat   `json:"format"`
	Options       *ExportOptions `json:"options"`
	EnvironmentID *int64         `json:"environment_id,omitempty"`
	Hash          string         `json:"hash"`
	GeneratedAt   time.Time      `json:"generated_at"`
}

// GenerateCacheKey generates a cache key for export results
func (c *ExportCache) GenerateCacheKey(collectionID int64, format ExportFormat, options *ExportOptions) string {
	if !c.config.Enabled {
		return ""
	}

	// Create a deterministic key based on collection, format, and options
	keyData := map[string]interface{}{
		"collection_id": collectionID,
		"format":        string(format),
		"options":       options,
	}

	// Convert to JSON for consistent hashing
	jsonData, err := json.Marshal(keyData)
	if err != nil {
		// Fallback to simple key if JSON marshaling fails
		return fmt.Sprintf("%s%d:%s", c.config.KeyPrefix, collectionID, format)
	}

	// Generate MD5 hash for compact key
	hash := fmt.Sprintf("%x", md5.Sum(jsonData))

	return fmt.Sprintf("%s%d:%s:%s", c.config.KeyPrefix, collectionID, format, hash[:8])
}

// Get retrieves cached export result
func (c *ExportCache) Get(key string) (*ExportResult, error) {
	if !c.config.Enabled || key == "" {
		return nil, nil
	}

	// Get cached data
	data, err := c.redis.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	// Deserialize cached result
	var result ExportResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		// If deserialization fails, delete the corrupted cache entry
		c.redis.Del(c.ctx, key)
		return nil, fmt.Errorf("failed to deserialize cached result: %w", err)
	}

	return &result, nil
}

// Set stores export result in cache
func (c *ExportCache) Set(key string, result *ExportResult) error {
	if !c.config.Enabled || key == "" || result == nil {
		return nil
	}

	// Serialize result
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to serialize export result: %w", err)
	}

	// Check size limit
	if c.config.MaxSize > 0 && int64(len(data)) > c.config.MaxSize {
		return fmt.Errorf("export result too large for cache: %d bytes (max: %d)", len(data), c.config.MaxSize)
	}

	// Store in cache with TTL
	err = c.redis.Set(c.ctx, key, data, c.config.TTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes cached export result
func (c *ExportCache) Delete(key string) error {
	if !c.config.Enabled || key == "" {
		return nil
	}

	err := c.redis.Del(c.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	return nil
}

// InvalidateCollection removes all cached exports for a collection
func (c *ExportCache) InvalidateCollection(collectionID int64) error {
	if !c.config.Enabled {
		return nil
	}

	// Find all keys for this collection
	pattern := fmt.Sprintf("%s%d:*", c.config.KeyPrefix, collectionID)
	keys, err := c.redis.Keys(c.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find cache keys: %w", err)
	}

	if len(keys) == 0 {
		return nil // No keys to delete
	}

	// Delete all matching keys
	err = c.redis.Del(c.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}

// InvalidateCollectionFormat removes cached exports for specific collection and format
func (c *ExportCache) InvalidateCollectionFormat(collectionID int64, format ExportFormat) error {
	if !c.config.Enabled {
		return nil
	}

	// Find all keys for this collection and format
	pattern := fmt.Sprintf("%s%d:%s:*", c.config.KeyPrefix, collectionID, format)
	keys, err := c.redis.Keys(c.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find cache keys: %w", err)
	}

	if len(keys) == 0 {
		return nil // No keys to delete
	}

	// Delete all matching keys
	err = c.redis.Del(c.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}

// GetCacheStats returns cache statistics
func (c *ExportCache) GetCacheStats() (*CacheStats, error) {
	if !c.config.Enabled {
		return &CacheStats{Enabled: false}, nil
	}

	// Get Redis info
	info, err := c.redis.Info(c.ctx, "memory").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	// Count keys with our prefix
	pattern := fmt.Sprintf("%s*", c.config.KeyPrefix)
	keys, err := c.redis.Keys(c.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to count cache keys: %w", err)
	}

	// Parse memory usage from Redis info
	var memoryUsed int64
	lines := strings.Split(info, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "used_memory:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				if mem, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					memoryUsed = mem
				}
			}
			break
		}
	}

	return &CacheStats{
		Enabled:    true,
		KeyCount:   int64(len(keys)),
		MemoryUsed: memoryUsed,
		TTL:        c.config.TTL,
		MaxSize:    c.config.MaxSize,
		KeyPrefix:  c.config.KeyPrefix,
	}, nil
}

// CacheStats contains cache statistics
type CacheStats struct {
	Enabled    bool          `json:"enabled"`
	KeyCount   int64         `json:"key_count"`
	MemoryUsed int64         `json:"memory_used_bytes"`
	TTL        time.Duration `json:"ttl"`
	MaxSize    int64         `json:"max_size_bytes"`
	KeyPrefix  string        `json:"key_prefix"`
}

// ClearAll removes all cached export results
func (c *ExportCache) ClearAll() error {
	if !c.config.Enabled {
		return nil
	}

	// Find all keys with our prefix
	pattern := fmt.Sprintf("%s*", c.config.KeyPrefix)
	keys, err := c.redis.Keys(c.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to find cache keys: %w", err)
	}

	if len(keys) == 0 {
		return nil // No keys to delete
	}

	// Delete all matching keys
	err = c.redis.Del(c.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}

// IsEnabled returns whether caching is enabled
func (c *ExportCache) IsEnabled() bool {
	return c.config.Enabled
}

// GetTTL returns the cache TTL
func (c *ExportCache) GetTTL() time.Duration {
	return c.config.TTL
}

// SetTTL updates the cache TTL
func (c *ExportCache) SetTTL(ttl time.Duration) {
	c.config.TTL = ttl
}
