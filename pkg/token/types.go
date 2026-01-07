package token

type TokenMetadata struct {
	UserID    int64    `json:"user_id"`
	UserAgent string   `json:"user_agent"`
	IP        string   `json:"ip"`
	Abilities []string `json:"abilities"`
	ExpiresAt int64    `json:"expires_at"`
}

type RefreshTokenMetadata struct {
	UserID   int64  `json:"user_id"`
	FamilyID string `json:"family_id"`
}

type TokenInfo struct {
	Type      string `json:"type"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	TTL       int64  `json:"ttl,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`
	FamilyID  string `json:"family_id,omitempty"`
}

type UserTokensResponse struct {
	UserID        int64       `json:"user_id"`
	AccessTokens  []TokenInfo `json:"access_tokens"`
	RefreshTokens []TokenInfo `json:"refresh_tokens"`
	HasValidToken bool        `json:"has_valid_token"`
}
