package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *Config) *redis.Client {
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  10 * time.Second,
	}

	// Enable TLS if configured
	if cfg.Redis.UseTLS {
		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true, // Skip certificate verification for cloud Redis
		}
		log.Println("Redis TLS enabled (InsecureSkipVerify: true)")
	}

	rdb := redis.NewClient(options)

	// Test connection with retry
	ctx := context.Background()
	maxRetries := 5
	var pingErr error

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to Redis (attempt %d/%d)...", i+1, maxRetries)
		_, pingErr = rdb.Ping(ctx).Result()
		if pingErr == nil {
			log.Println("✅ Redis connected successfully")
			return rdb
		}
		log.Printf("❌ Redis ping attempt %d failed: %v", i+1, pingErr)
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}

	log.Fatalf("Failed to connect to Redis after %d attempts: %v", maxRetries, pingErr)
	return nil
}
