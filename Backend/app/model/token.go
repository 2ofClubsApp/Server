package model

type TokenInfo struct {
	AccessToken  string
	RefreshToken string
}

func NewTokenInfo() *TokenInfo{
	return &TokenInfo{}
}

const RefreshToken = "RefreshToken"
const AccessToken = "AccessToken"