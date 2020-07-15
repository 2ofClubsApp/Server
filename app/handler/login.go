package handler

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

const (
	LoginFailure = "Username or Password is Incorrect"
	HashErr      = "hashing Error"
	ErrTokenGen  = "token generation error"
)

// TODO: Prevent login many times (if user tries to brute force this)
func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	creds := getCredentials(r)
	hash, isFound := getPasswordHash(db, creds.Username)
	var err error
	if isFound {
		err = bcrypt.CompareHashAndPassword(hash, []byte(creds.Password))
	}
	if err != nil {
		s := model.NewStatus()
		s.Message = LoginFailure
		s.Code = model.FailureCode
		WriteData(GetJSON(s), http.StatusOK, w)
	} else {
		if tp, err := GetTokenPair(creds.Username, 5, 60*24); err == nil {
			c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
			http.SetCookie(w, c)
			WriteData(tp.AccessToken, http.StatusOK, w)
		}
	}
}

func getCredentials(r *http.Request) *model.Credentials {
	decoder := json.NewDecoder(r.Body)
	cred := model.NewCredentials()
	decoder.Decode(cred)
	cred.Username = strings.ToLower(cred.Username)
	return cred
}

/*
	Gets password hash for both clubs and users provided the username.
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

func GenerateJWT(subject string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	sysTime := time.Now()
	claims["iat"] = sysTime
	claims["exp"] = sysTime.Add(time.Minute * duration).Unix()
	claims["sub"] = subject // Subject usually as a number (unique value?)
	// Note: This must be changed to an env variable later
	tokenString, err := token.SignedString([]byte("2ofClubs"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetTokenPair(subject string, accessDuration time.Duration, refreshDuration time.Duration) (*model.TokenInfo, error) {
	if accessToken, atErr := GenerateJWT(subject, accessDuration); atErr == nil {
		if refreshToken, rtErr := GenerateJWT(subject, refreshDuration); rtErr == nil {
			token := model.NewTokenInfo()
			token.AccessToken = accessToken
			token.RefreshToken = refreshToken
			return token, nil
		}
	}
	return nil, fmt.Errorf(ErrTokenGen)
}
