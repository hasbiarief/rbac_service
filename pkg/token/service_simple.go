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

// StoreAccessToken stores access token - 1 token per user (overwrites existing)
func (ts *SimpleTokenService) StoreAccessToken(token string, metadata TokenMetadata, ttl time.Duration) error {
	ctx := context.Background()

	// Single key per user - this ensures only 1 access token per user
	key := fmt.Sprintf("access:user:%d", metadata.UserID)

	// Store both the token and metadata
	sessionData := map[string]interface{}{
		"token":    token,
		"metadata": metadata,
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return ts.redis.Set(ctx, key, data, ttl).Err()
}

// StoreRefreshToken stores refresh token - 1 token per user (overwrites existing)
func (ts *SimpleTokenService) StoreRefreshToken(token string, metadata RefreshTokenMetadata, ttl time.Duration) error {
	ctx := context.Background()

	// Single key per user - this ensures only 1 refresh token per user
	key := fmt.Sprintf("refresh:user:%d", metadata.UserID)

	// Store both the token and metadata
	sessionData := map[string]interface{}{
		"token":    token,
		"metadata": metadata,
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return ts.redis.Set(ctx, key, data, ttl).Err()
}

// GetAccessToken retrieves access token metadata from Redis
func (ts *SimpleTokenService) GetAccessToken(token string) (*TokenMetadata, error) {
	ctx := context.Background()

	// Scan all access keys to find the matching token
	pattern := "access:user:*"
	keys, err := ts.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
			continue
		}

		// Check if this is the token we're looking for
		if storedToken, ok := sessionData["token"].(string); ok && storedToken == token {
			// Extract metadata
			metadataBytes, err := json.Marshal(sessionData["metadata"])
			if err != nil {
				continue
			}

			var metadata TokenMetadata
			if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
				continue
			}

			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("token not found")
}

// GetRefreshToken retrieves refresh token metadata from Redis
func (ts *SimpleTokenService) GetRefreshToken(token string) (*RefreshTokenMetadata, error) {
	ctx := context.Background()

	// Scan all refresh keys to find the matching token
	pattern := "refresh:user:*"
	keys, err := ts.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
			continue
		}

		// Check if this is the token we're looking for
		if storedToken, ok := sessionData["token"].(string); ok && storedToken == token {
			// Extract metadata
			metadataBytes, err := json.Marshal(sessionData["metadata"])
			if err != nil {
				continue
			}

			var metadata RefreshTokenMetadata
			if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
				continue
			}

			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("token not found")
}

// RevokeToken removes token from Redis
func (ts *SimpleTokenService) RevokeToken(token string, tokenType string) error {
	ctx := context.Background()

	// Find the token key by scanning
	pattern := fmt.Sprintf("%s:user:*", tokenType)
	keys, err := ts.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
			continue
		}

		// Check if this is the token we're looking for
		if storedToken, ok := sessionData["token"].(string); ok && storedToken == token {
			return ts.redis.Del(ctx, key).Err()
		}
	}

	return nil
}

// RevokeAllUserTokens removes all tokens for a user
func (ts *SimpleTokenService) RevokeAllUserTokens(userID int64) error {
	ctx := context.Background()

	// Delete access token
	accessKey := fmt.Sprintf("access:user:%d", userID)
	ts.redis.Del(ctx, accessKey)

	// Delete refresh token
	refreshKey := fmt.Sprintf("refresh:user:%d", userID)
	ts.redis.Del(ctx, refreshKey)

	return nil
}

// GetUserTokens retrieves tokens for a user (for frontend token check)
func (ts *SimpleTokenService) GetUserTokens(userID int64) (*UserTokensResponse, error) {
	ctx := context.Background()

	var accessTokens []TokenInfo
	var refreshTokens []TokenInfo

	// Check access token
	accessKey := fmt.Sprintf("access:user:%d", userID)
	accessData, err := ts.redis.Get(ctx, accessKey).Result()
	if err == nil {
		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(accessData), &sessionData); err == nil {
			if metadataRaw, ok := sessionData["metadata"]; ok {
				metadataBytes, _ := json.Marshal(metadataRaw)
				var metadata TokenMetadata
				if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
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
	}

	// Check refresh token
	refreshKey := fmt.Sprintf("refresh:user:%d", userID)
	refreshData, err := ts.redis.Get(ctx, refreshKey).Result()
	if err == nil {
		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(refreshData), &sessionData); err == nil {
			if metadataRaw, ok := sessionData["metadata"]; ok {
				metadataBytes, _ := json.Marshal(metadataRaw)
				var metadata RefreshTokenMetadata
				if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
					// Get TTL for refresh token
					ttl, err := ts.redis.TTL(ctx, refreshKey).Result()
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
func (ts *SimpleTokenService) CleanupExpiredTokens() error {
	ctx := context.Background()

	// Clean up expired access tokens
	accessPattern := "access:user:*"
	accessKeys, err := ts.redis.Keys(ctx, accessPattern).Result()
	if err != nil {
		return err
	}

	for _, key := range accessKeys {
		data, err := ts.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var sessionData map[string]interface{}
		if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
			continue
		}

		if metadataRaw, ok := sessionData["metadata"]; ok {
			metadataBytes, _ := json.Marshal(metadataRaw)
			var metadata TokenMetadata
			if err := json.Unmarshal(metadataBytes, &metadata); err == nil {
				// Check if token is expired
				if time.Now().Unix() > metadata.ExpiresAt {
					ts.redis.Del(ctx, key)
				}
			}
		}
	}

	return nil
}
