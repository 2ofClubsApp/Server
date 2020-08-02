package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

const (
	HashErr     = "hashing Error"
	ErrTokenGen = "token generation error"
)

// TODO: Prevent login many times (if user tries to brute force this)
func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	creds := model.NewCredentials()
	extractBody(r, creds)
	creds.Username = strings.ToLower(creds.Username)
	hash, isFound := getPasswordHash(db, creds.Username)
	var err error
	if isFound {
		err = bcrypt.CompareHashAndPassword(hash, []byte(creds.Password))
	}
	status := model.NewStatus()
	if err != nil {
		status.Message = model.LoginFailure
	} else {
		if tp, err := GetTokenPair(creds.Username, 5, 60*24); err == nil && isApproved(db, creds.Username) {
			c := GenerateCookie(model.RefreshToken, tp.RefreshToken)
			http.SetCookie(w, c)
			type login struct {
				Token string
			}
			status.Code = model.SuccessCode
			status.Message = model.LoginSuccess
			status.Data = login{Token: tp.AccessToken}
		} else {
			status.Message = model.NotApproved
		}
	}
	WriteData(GetJSON(status), http.StatusOK, w)

}
func isApproved(db *gorm.DB, username string) bool {
	u := model.NewUser()
	SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, u)
	return u.IsApproved
}

//func getCredentials(r *http.Request) *model.Credentials {
//	decoder := json.NewDecoder(r.Body)
//	cred := model.NewCredentials()
//	decoder.Decode(cred)
//	cred.Username = strings.ToLower(cred.Username)
//	return cred
//}


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
