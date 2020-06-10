package model

type TokenInfo struct {
	AccessToken  string
	RefreshToken string
}

func NewTokenInfo() *TokenInfo{
	return &TokenInfo{}
}
