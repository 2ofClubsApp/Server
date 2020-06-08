package model

type TokenInfo struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}
