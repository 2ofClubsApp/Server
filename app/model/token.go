package model

// TokenInfo struct containing the Access and Refresh token
type TokenInfo struct {
	// Life span of 5 minutes
	AccessToken string

	// Life span of 1 hour
	RefreshToken string
}

func NewTokenInfo() *TokenInfo {
	return &TokenInfo{}
}

const RefreshToken = "RefreshToken"
const AccessToken = "AccessToken"
const TokenVar = "token"
