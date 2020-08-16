package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)


// TODO: Prevent login many times (if user tries to brute force this)
func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	creds := model.NewCredentials()
	extractBody(r, creds)
	creds.Username = strings.ToLower(creds.Username)
	hash, isFound := getPasswordHash(db, creds.Username)
	var err error
	s := status.New()
	if isFound {
		err = bcrypt.CompareHashAndPassword(hash, []byte(creds.Password))
	}
	if err != nil {
		s.Message = status.LoginFailure
		httpStatus = http.StatusUnauthorized
	} else {
		userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, creds.Username, model.NewUser())
		if tp, err := GetTokenPair(creds.Username, 5, 60*24); err == nil && userExists {
			c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
			http.SetCookie(w, c)
			type login struct {
				Token string
			}
			s.Code = status.SuccessCode
			s.Message = status.LoginSuccess
			s.Data = login{Token: tp.AccessToken}
		} else {
			s.Message = status.UserNotApproved
		}
		httpStatus = http.StatusOK
	}
	WriteData(GetJSON(s), httpStatus, w)
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

/*
Generating http cookie where the refresh token will be embedded.
*/
func GenerateCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
	}
}

func GenerateJWT(subject string, duration time.Duration, jwtSecret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sysTime := time.Now()
	claims["iat"] = sysTime
	claims["exp"] = sysTime.Add(time.Minute * duration).Unix()
	claims["sub"] = subject // Subject usually as a number (unique value?)
	// Note: This must be changed to an env variable later
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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
