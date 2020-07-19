package handler

import (
	"../model"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

func isAdmin(db *gorm.DB, r *http.Request) bool{
	claims := GetTokenClaims(r)
	subject := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	if SingleRecordExists(db, model.UserTable, model.UsernameColumn, subject, user){
		return user.IsAdmin
	}
	return false
}
