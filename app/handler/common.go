package handler

import (
	"encoding/json"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
)

/*
	Common methods shared amongst the different models
*/

func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteData(http.StatusText(http.StatusNotFound), http.StatusNotFound, w)
	})
}

/*
Extracting variables from request URL
*/
func getVar(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}

/*
Extract the Token Claims from the HTTP Request Header
*/
func GetTokenClaims(r *http.Request) jwt.MapClaims {
	t := r.Header.Get("Authorization")
	splitToken := strings.Split(t, "Bearer")
	token := strings.TrimSpace(splitToken[1])
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, &claims, KF(os.Getenv("JWT_SECRET")))
	return claims
}

// Note: Need to add more authentication checks later (This is temporary)

/*
Return true if the JWT is valid, false otherwise
*/
func VerifyJWT(r *http.Request) bool {
	if bearerToken := r.Header.Get("Authorization"); bearerToken != "" {
		splitToken := strings.Split(bearerToken, "Bearer ")
		token := strings.TrimSpace(splitToken[1])
		return IsValidJWT(token, KF(os.Getenv("JWT_SECRET")))
	}
	return false
}

func IsValidJWT(token string, kf jwt.Keyfunc) bool {
	if t, err := jwt.Parse(token, kf); err == nil {
		if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			return true
		}
	}
	return false
}

func KF(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Verifying that the signing method is the same before continuing any further
		if _, accepted := token.Method.(*jwt.SigningMethodHMAC); !accepted {
			return nil, fmt.Errorf(model.ErrGeneric)
		}
		return []byte(secret), nil
	}
}

/*
Returning true if the record already exists in the table, false otherwise.
*/
// When soft deleted, the record won't exist anymore
func SingleRecordExists(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	result := db.Table(tableName).Where(column+"= ?", val).First(t)
	return result.Error == nil
}

/*
Returning the JSON representation of a struct.
*/
func GetJSON(response interface{}) string {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

/*
Return response message and an HTTP Status Code upon receiving a request.
*/
func WriteData(data string, code int, w http.ResponseWriter) int {
	w.WriteHeader(code)
	n, err := fmt.Fprint(w, data)
	if err != nil {
		return -1
	}
	return n
}
