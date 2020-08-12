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
		if isAdmin(db, r) {
			db.Model(user).Update(model.IsApprovedColumn, !user.IsApproved)
			status.Message = model.UserUpdated
			status.Code = model.SuccessCode
		} else {
			status.Message = model.AdminRequired
		}
	} else {
		status.Message = model.UserNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func ToggleClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	clubID := getVar(r, model.ClubIDVar)
	status := model.NewStatus()
	club := model.NewClub()
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if clubExists {
		if isAdmin(db, r) {
			db.Model(club).Update(model.IsActiveColumn, !club.Active)
			status.Message = model.ClubUpdateSuccess
			status.Code = model.SuccessCode
		} else {
			status.Message = model.AdminRequired
		}
	} else {
		status.Message = model.ClubNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func GetToggleUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	users := []*model.User{}
	result := db.Where(model.IsApprovedColumn+"= ?", false).Find(&users)
	status := model.NewStatus()
	if isAdmin(db, r) {
		if result.Error != nil {
			status.Message = model.GetNonToggledUsersFailure
		} else {
			status.Code = model.SuccessCode
			status.Message = model.GetNonToggledUsersSuccess
			status.Data = users
			for _, v := range users {
				fmt.Println(v)
			}
		}
	} else {
		status.Message = model.AdminRequired
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
