package handler

import (
	"encoding/json"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const (
	ErrGeneric  = "an error occurred"
	ErrTokenGen = "token generation error"
)

/*
	Common methods shared amongst the different models
*/

func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteData(http.StatusText(http.StatusNotFound), http.StatusNotFound, w)
	})
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

}

// Validate the user request to ensure that they can only access/modify their own respective data
func ValidateUserReq(username string, r *http.Request) bool {
	t := r.Header["Token"][0]
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(t, &claims, kf)
	sub := fmt.Sprintf("%v", claims["sub"])
	fmt.Println(username)
	fmt.Println(sub)
	return sub == username
}

/*
Note: Need to add more authentication checks later (This is temporary)
*/
func IsValidJWT(w http.ResponseWriter, r *http.Request) bool {
	if token := r.Header["Token"]; token != nil {
		if t, err := jwt.Parse(token[0], kf); err == nil {
			if t.Valid {
				return true
			}
		}
	}
	WriteData(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, w)
	return false
}

func kf(token *jwt.Token) (interface{}, error) {
	// Verifying that the signing method is the same before continuing any further
	if _, accepted := token.Method.(*jwt.SigningMethodHMAC); !accepted {
		return nil, fmt.Errorf(ErrGeneric)
	}
	// Note: This must be changed to an env variable later
	return []byte("2ofClubs"), nil
}

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

func Hash(info string) (string, bool) {
	// Change cost to 10+ (try to find a way to scale it with hardware?)
	saltedHashPass, err := bcrypt.GenerateFromPassword([]byte(info), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}
	return string(saltedHashPass), true
}
func RecordExists(db *gorm.DB, column string, val string, t interface{}) bool {
	return !db.Where(column+"= ?", val).First(t).RecordNotFound()
}

func ExtractPersonInfo(r *http.Request) model.Person {
	decoder := json.NewDecoder(r.Body)
	p := model.NewPerson()
	decoder.Decode(&p)
	return p
}

func ParseJSON(response interface{}) string {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func WriteData(data string, code int, w http.ResponseWriter) int {
	w.WriteHeader(code)
	n, err := fmt.Fprint(w, data)
	if err != nil {
		return -1
	}
	return n
}
