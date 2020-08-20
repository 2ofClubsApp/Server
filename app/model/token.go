package model

// TokenInfo struct containing the Access and Refresh token
type TokenInfo struct {
	// Life span of 5 minutes
	AccessToken string

	// Life span of 1 hour
	RefreshToken string
}

// NewTokenInfo Create new default TokenInfo
func NewTokenInfo() *TokenInfo {
	return &TokenInfo{}
}

// Token variable constants
const (
	RefreshToken = "RefreshToken"
	TokenVar     = "token"
)
