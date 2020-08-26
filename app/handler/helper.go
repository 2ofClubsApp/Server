package handler

import (
	"context"
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
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

// GetTokenClaims - Extract the Token Claims from the HTTP Request Header
func GetTokenClaims(token string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, &claims, KF(os.Getenv("JWT_SECRET"))); err != nil {
		return jwt.MapClaims{}
	}
	return claims

}

// VerifyJWT - Return true if the JWT is valid, false otherwise
func VerifyJWT(r *http.Request) bool {
	if token := ExtractToken(r); token != "" {
		return IsValidJWT(token, KF(os.Getenv("JWT_SECRET")))
	}
	return false
}

// ExtractToken returns the JWT token if it's provided otherwise, an empty string will be returned
func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return strings.TrimSpace(splitToken[1])
}

// IsValidJWT - Returning true whether the JWT is valid
func IsValidJWT(token string, kf jwt.Keyfunc) bool {
	if t, err := jwt.Parse(token, kf); err == nil {
		if _, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			return true
		}
	}
	return false
}

// IsActiveToken returns whether the token is active or not (in the redis cache)
func IsActiveToken(rc *redis.Client, r *http.Request) bool {
	ctx := context.Background()
	token := ExtractToken(r)
	claims := GetTokenClaims(token)
	uname := fmt.Sprintf("%v", claims["sub"])
	return rc.Get(ctx, "access_"+uname).Val() == token
}

// KF - Key Function to verify the token signing method (Used in conjunction with IsValidJWT)
func KF(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Verifying that the signing method is the same before continuing any further
		if _, accepted := token.Method.(*jwt.SigningMethodHMAC); !accepted {
			return nil, fmt.Errorf(status.ErrGeneric)
		}
		return []byte(secret), nil
	}
}

// IsSingleRecordActive - Returns true if the record is active
// This is used for verifying users and clubs as they need to be activated upon creation
func IsSingleRecordActive(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	exists := SingleRecordExists(db, tableName, column, val, t)
	if exists {
		switch m := t.(type) {
		case *model.Club:
			return m.Active
		case *model.User:
			return m.IsApproved
		}
	}
	return false
}

// SingleRecordExists - Returning true if the record already exists in the table, false otherwise.
func SingleRecordExists(db *gorm.DB, tableName string, column string, val string, t interface{}) bool {
	result := db.Table(tableName).Where(column+"= ?", val).First(t)
	return result.Error == nil
}

// Returning true whether the user is an admin or not
func isAdmin(db *gorm.DB, r *http.Request) bool {
	claims := GetTokenClaims(ExtractToken(r))
	subject := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	// If the user is an admin, it would already be active by default (No need to check for it's active state)
	if SingleRecordExists(db, model.UserTable, model.UsernameColumn, subject, user) {
		return user.IsAdmin
	}
	return false
}
