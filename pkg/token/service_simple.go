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

type SimpleTokenService struct {
	redis *redis.Client
}

func NewSimpleTokenService(redis *redis.Client) *SimpleTokenService {
	return &SimpleTokenService{
		redis: redis,
	}
}

// GenerateToken generates a random token
func (ts *SimpleTokenService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashToken creates SHA-256 hash of token
func (ts *SimpleTokenService) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// StoreAccessToken stores access token with token as key
func (ts *SimpleTokenService) StoreAccessToken(token string, metadata TokenMetadata, ttl time.Duration) error {
	ctx := context.Background()

	// Use token as key directly
	tokenKey := fmt.Sprintf("access:token:%s", token)
	userKey := fmt.Sprintf("access:user:%d", metadata.UserID)

	// Store token metadata
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Store token with metadata
	if err := ts.redis.Set(ctx, tokenKey, data, ttl).Err(); err != nil {
		return err
	}

	// Store user -> token mapping (for logout all)
	return ts.redis.Set(ctx, userKey, token, ttl).Err()
}

// StoreRefreshToken stores refresh token with token as key
func (ts *SimpleTokenService) StoreRefreshToken(token string, metadata RefreshTokenMetadata, ttl time.Duration) error {
	ctx := context.Background()

	// Use token as key directly
	tokenKey := fmt.Sprintf("refresh:token:%s", token)
	userKey := fmt.Sprintf("refresh:user:%d", metadata.UserID)

	// Store token metadata
	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Store token with metadata
	if err := ts.redis.Set(ctx, tokenKey, data, ttl).Err(); err != nil {
		return err
	}

	// Store user -> token mapping (for logout all)
	return ts.redis.Set(ctx, userKey, token, ttl).Err()
}

// GetAccessToken retrieves access token metadata from Redis
func (ts *SimpleTokenService) GetAccessToken(token string) (*TokenMetadata, error) {
	ctx := context.Background()

	// Direct key lookup - no KEYS command needed
	tokenKey := fmt.Sprintf("access:token:%s", token)
	data, err := ts.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("token not found or expired")
	}

	var metadata TokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// GetRefreshToken retrieves refresh token metadata from Redis
func (ts *SimpleTokenService) GetRefreshToken(token string) (*RefreshTokenMetadata, error) {
	ctx := context.Background()

	// Direct key lookup - no KEYS command needed
	tokenKey := fmt.Sprintf("refresh:token:%s", token)
	data, err := ts.redis.Get(ctx, tokenKey).Result()
	if err != nil {
		return nil, fmt.Errorf("token not found or expired")
	}

	var metadata RefreshTokenMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// RevokeToken removes token from Redis
func (ts *SimpleTokenService) RevokeToken(token string, tokenType string) error {
	ctx := context.Background()

	// Direct key deletion - no KEYS command needed
	tokenKey := fmt.Sprintf("%s:token:%s", tokenType, token)
	return ts.redis.Del(ctx, tokenKey).Err()
}

// RevokeAllUserTokens removes all tokens for a user
func (ts *SimpleTokenService) RevokeAllUserTokens(userID int64) error {
	ctx := context.Background()

	// Get user's tokens from user mapping
	accessUserKey := fmt.Sprintf("access:user:%d", userID)
	refreshUserKey := fmt.Sprintf("refresh:user:%d", userID)

	// Get access token
	if accessToken, err := ts.redis.Get(ctx, accessUserKey).Result(); err == nil {
		accessTokenKey := fmt.Sprintf("access:token:%s", accessToken)
		ts.redis.Del(ctx, accessTokenKey)
	}

	// Get refresh token
	if refreshToken, err := ts.redis.Get(ctx, refreshUserKey).Result(); err == nil {
		refreshTokenKey := fmt.Sprintf("refresh:token:%s", refreshToken)
		ts.redis.Del(ctx, refreshTokenKey)
	}

	// Delete user mappings
	ts.redis.Del(ctx, accessUserKey)
	ts.redis.Del(ctx, refreshUserKey)

	return nil
}

// GetUserTokens retrieves tokens for a user (for frontend token check)
func (ts *SimpleTokenService) GetUserTokens(userID int64) (*UserTokensResponse, error) {
	ctx := context.Background()

	var accessTokens []TokenInfo
	var refreshTokens []TokenInfo

	// Check access token via user mapping
	accessUserKey := fmt.Sprintf("access:user:%d", userID)
	if accessToken, err := ts.redis.Get(ctx, accessUserKey).Result(); err == nil {
		accessTokenKey := fmt.Sprintf("access:token:%s", accessToken)
		accessData, err := ts.redis.Get(ctx, accessTokenKey).Result()
		if err == nil {
			var metadata TokenMetadata
			if err := json.Unmarshal([]byte(accessData), &metadata); err == nil {
				// Check if token is still valid
				if time.Now().Unix() < metadata.ExpiresAt {
					accessTokens = append(accessTokens, TokenInfo{
						Type:      "access",
						ExpiresAt: metadata.ExpiresAt,
						UserAgent: metadata.UserAgent,
						IP:        metadata.IP,
					})
				}
			}
		}
	}

	// Check refresh token via user mapping
	refreshUserKey := fmt.Sprintf("refresh:user:%d", userID)
	if refreshToken, err := ts.redis.Get(ctx, refreshUserKey).Result(); err == nil {
		refreshTokenKey := fmt.Sprintf("refresh:token:%s", refreshToken)
		refreshData, err := ts.redis.Get(ctx, refreshTokenKey).Result()
		if err == nil {
			var metadata RefreshTokenMetadata
			if err := json.Unmarshal([]byte(refreshData), &metadata); err == nil {
				// Get TTL for refresh token
				ttl, err := ts.redis.TTL(ctx, refreshTokenKey).Result()
				if err == nil && ttl > 0 {
					refreshTokens = append(refreshTokens, TokenInfo{
						Type:     "refresh",
						TTL:      int64(ttl.Seconds()),
						FamilyID: metadata.FamilyID,
					})
				}
			}
		}
	}

	return &UserTokensResponse{
		UserID:        userID,
		AccessTokens:  accessTokens,
		RefreshTokens: refreshTokens,
		HasValidToken: len(accessTokens) > 0 || len(refreshTokens) > 0,
	}, nil
}

// GetUserSessionCount returns the number of active sessions for a user (always 0 or 1)
func (ts *SimpleTokenService) GetUserSessionCount(userID int64) (int, error) {
	tokensResponse, err := ts.GetUserTokens(userID)
	if err != nil {
		return 0, err
	}

	if len(tokensResponse.AccessTokens) > 0 {
		return 1, nil
	}

	return 0, nil
}

// CleanupExpiredTokens removes expired tokens from Redis
// Note: This is a no-op since Redis TTL handles expiration automatically
func (ts *SimpleTokenService) CleanupExpiredTokens() error {
	// Redis automatically removes expired keys with TTL
	// No manual cleanup needed
	return nil
}
