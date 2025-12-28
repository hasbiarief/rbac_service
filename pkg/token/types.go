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
