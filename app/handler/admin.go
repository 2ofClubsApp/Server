package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"gorm.io/gorm"
	"net/http"
)

func isAdmin(db *gorm.DB, r *http.Request) bool {
	claims := GetTokenClaims(r)
	subject := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	if SingleRecordExists(db, model.UserTable, model.UsernameColumn, subject, user) {
		return user.IsAdmin
	}
	return false
}

func ToggleUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	username := getVar(r, model.UsernameVar)
	status := model.NewStatus()
	user := model.NewUser()
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
	if userExists {
		db.Model(user).Update(model.IsApprovedColumn, !user.IsApproved)
		status.Message = model.UserUpdated
		status.Code = model.SuccessCode
	} else {
		status.Message = model.UserNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
