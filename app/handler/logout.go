package handler

import (
	"context"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
)

// Logout logs out a user
// Any access tokens will be revoked and can no longer be used (even if they are still valid)
func Logout(_ *gorm.DB, rc *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	ctx := context.Background()
	requestUsername := getVar(r, model.UsernameVar)
	claims := GetTokenClaims(ExtractToken(r))
	tokenUsername := fmt.Sprintf("%v", claims["sub"])
	if tokenUsername == requestUsername {
		s.Code = status.SuccessCode
		s.Message = status.LogoutSuccess
		rc.Del(ctx, "access_"+requestUsername)
		return http.StatusOK, nil
	}
	s.Message = status.LogoutFailure
	return http.StatusForbidden, nil
}
