package handler

import (
	"context"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	accessDuration  = 5       // 5 minutes
	refreshDuration = 60 * 24 // 24 hours
)

// Login - User Login
//   See model.credentials or docs for username and email constraints
func Login(db *gorm.DB, rc *redis.Client, w http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	ctx := context.Background()
	creds := model.NewCredentials()
	if extractBody(r, creds) != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	creds.Username = strings.ToLower(creds.Username)
	hash, passFound := getPasswordHash(db, creds.Username)
	if !passFound {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	err := bcrypt.CompareHashAndPassword(hash, []byte(creds.Password))
	if err != nil {
		s.Message = status.LoginFailure
		return http.StatusUnauthorized, nil
	}
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, creds.Username, model.NewUser())
	if tp, err := GetTokenPair(creds.Username, accessDuration, refreshDuration); err == nil && userExists {
		c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
		http.SetCookie(w, c)
		type login struct {
			Token string
		}
		s.Code = status.SuccessCode
		s.Message = status.LoginSuccess
		s.Data = login{Token: tp.AccessToken}
		fmt.Println(creds.Username)
		fmt.Println(tp.AccessToken)
		val, err := rc.Set(ctx, creds.Username, tp.AccessToken, 300000000000).Result()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		fmt.Println(val)
		return http.StatusOK, nil
	}
	s.Message = status.UserNotApproved
	return http.StatusUnauthorized, nil
}

/*
Gets password hash for a user given the username.
*/
func getPasswordHash(db *gorm.DB, userName string) ([]byte, bool) {
	type p struct {
		Password string
	}
	pass := &p{}
	result := db.Table(model.UserTable).Where("Username = ?", userName).Find(pass)
	if result.Error != nil {
		return []byte(""), false
	}
	return []byte(pass.Password), true

}

// GenerateCookie - Generating http cookie where the refresh token will be embedded.
func GenerateCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
	}
}

// GenerateJWT - Generating a JWT based on the user's username
func GenerateJWT(subject string, duration time.Duration, jwtSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sysTime := time.Now()
	claims["iat"] = sysTime
	claims["exp"] = sysTime.Add(time.Minute * duration).Unix()
	claims["sub"] = subject // Usernames are unique to each user
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetTokenPair - Obtaining an access and refresh token pair
// Note: Access tokens have a lifespan of 5 minutes
//	  Refresh tokens have a lifespan of 24 hours (1 day)
func GetTokenPair(subject string, accessDuration time.Duration, refreshDuration time.Duration) (*model.TokenInfo, error) {
	if accessToken, atErr := GenerateJWT(subject, accessDuration, os.Getenv("JWT_SECRET")); atErr == nil {
		if refreshToken, rtErr := GenerateJWT(subject, refreshDuration, os.Getenv("JWT_SECRET")); rtErr == nil {
			token := model.NewTokenInfo()
			token.AccessToken = accessToken
			token.RefreshToken = refreshToken
			return token, nil
		}
	}
	return nil, fmt.Errorf(status.ErrTokenGen)
}
