package handler

import (
	"context"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

// RefreshToken verifies the refresh token and obtains a new set of access and refresh tokens
func RefreshToken(_ *gorm.DB, rc *redis.Client, w http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	ctx := context.Background()
	currentRefreshToken := ExtractToken(r)
	if IsValidJWT(currentRefreshToken, KF(os.Getenv("JWT_SECRET"))) {
		claims := GetTokenClaims(currentRefreshToken)
		username := fmt.Sprintf("%v", claims["sub"])
		tokenInfo, err := GetTokenPair(username, accessDuration, refreshDuration)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to generate token pair")
		}
		c := GenerateCookie(model.RefreshTokenVar, tokenInfo.RefreshToken)
		if refreshToken, err := rc.Get(ctx, "refresh_"+username).Result(); refreshToken != currentRefreshToken && err != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to get refresh token from cache")
		}
		_, err = rc.Set(ctx, "access_"+username, tokenInfo.AccessToken, time.Duration(accessDuration*minuteToNanosecond)).Result()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		_, err = rc.Set(ctx, "refresh_"+username, tokenInfo.RefreshToken, time.Duration(refreshDuration*minuteToNanosecond)).Result()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		http.SetCookie(w, c)
		type login struct {
			Token string `json:"token"`
		}
		s.Code = status.SuccessCode
		s.Message = status.TokenPairGenerateSuccess
		s.Data = login{Token: tokenInfo.AccessToken}
		return http.StatusOK, nil
	}
	return http.StatusForbidden, nil
}
