package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenService struct {
	redis *redis.Client
}

func NewTokenService(redis *redis.Client) *TokenService {
	return &TokenService{
		redis: redis,
	}
}

// GenerateToken generates a random token
func (ts *TokenService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates SHA-256 hash of token
func (ts *TokenService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// StoreAccessToken stores access token in Redis
func (ts *TokenService) StoreAccessToken(token string, metadata TokenMetadata, ttl time.Duration) error {
	ctx := context.Background()
	hash := ts.HashToken(token)
	key := fmt.Sprintf("access:v1:%s", hash)

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return ts.redis.Set(ctx, key, data, ttl).Err()
}

// StoreRefreshToken stores refresh token in Redis
func (ts *TokenService) StoreRefreshToken(token string, metadata RefreshTokenMetadata, ttl time.Duration) error {
	ctx := context.Background()
	hash := ts.HashToken(token)
	key := fmt.Sprintf("refresh:v1:%s", hash)

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return ts.redis.Set(ctx, key, data, ttl).Err()
}

// GetAccessToken retrieves access token metadata from Redis
func (ts *TokenService) GetAccessToken(token string) (*TokenMetadata, error) {
	ctx := context.Background()
	hash := ts.HashToken(token)
	key := fmt.Sprintf("access:v1:%s", hash)

	data, err := ts.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var metadata TokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// GetRefreshToken retrieves refresh token metadata from Redis
func (ts *TokenService) GetRefreshToken(token string) (*RefreshTokenMetadata, error) {
	ctx := context.Background()
	hash := ts.HashToken(token)
	key := fmt.Sprintf("refresh:v1:%s", hash)

	data, err := ts.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var metadata RefreshTokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// RevokeToken removes token from Redis
func (ts *TokenService) RevokeToken(token string, tokenType string) error {
	ctx := context.Background()
	hash := ts.HashToken(token)
	key := fmt.Sprintf("%s:v1:%s", tokenType, hash)

	return ts.redis.Del(ctx, key).Err()
}

// RevokeAllUserTokens removes all tokens for a user
func (ts *TokenService) RevokeAllUserTokens(userID int64) error {
	ctx := context.Background()

	// Find all access tokens for user
	accessPattern := "access:v1:*"
	accessKeys, err := ts.redis.Keys(ctx, accessPattern).Result()
	if err != nil {
		return err
	}

	for _, key := range accessKeys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var metadata TokenMetadata
		if err := json.Unmarshal([]byte(data), &metadata); err != nil {
			continue
		}

		if metadata.UserID == userID {
			ts.redis.Del(ctx, key)
		}
	}

	// Find all refresh tokens for user
	refreshPattern := "refresh:v1:*"
	refreshKeys, err := ts.redis.Keys(ctx, refreshPattern).Result()
	if err != nil {
		return err
	}

	for _, key := range refreshKeys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var metadata RefreshTokenMetadata
		if err := json.Unmarshal([]byte(data), &metadata); err != nil {
			continue
		}

		if metadata.UserID == userID {
			ts.redis.Del(ctx, key)
		}
	}

	return nil
}
