package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
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

// Returning true whether the JWT is valid
func IsValidJWT(token string, kf jwt.Keyfunc) bool {
	if t, err := jwt.Parse(token, kf); err == nil {
		if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			return true
		}
	}
	return false
}

// Key Function to verify the token signing method (Used in conjunction with IsValidJWT)
func KF(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Verifying that the signing method is the same before continuing any further
		if _, accepted := token.Method.(*jwt.SigningMethodHMAC); !accepted {
			return nil, fmt.Errorf(status.ErrGeneric)
		}
		return []byte(secret), nil
	}
}

// Returns true if the record is active
// This is used for verifying users and clubs as they need to be activated upon creation
func IsSingleRecordActive(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	exists := SingleRecordExists(db, tableName, column, val, t)
	if exists {
		switch model := t.(type) {
		case *model.Club:
			return model.Active
		case *model.User:
			return model.IsApproved
		}
	}
	return false
}

//Returning true if the record already exists in the table, false otherwise.
func SingleRecordExists(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	result := db.Table(tableName).Where(column+"= ?", val).First(t)
	return result.Error == nil
}

// Returning true whether the user is an admin or not
func isAdmin(db *gorm.DB, r *http.Request) bool {
	claims := GetTokenClaims(r)
	subject := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	// If the user is an admin, it would already be active by default (No need to check for it's active state)
	if SingleRecordExists(db, model.UserTable, model.UsernameColumn, subject, user) {
		return user.IsAdmin
	}
	return false
}
