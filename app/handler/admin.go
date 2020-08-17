package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"gorm.io/gorm"
	"net/http"
)

func ToggleUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	user := model.NewUser()
	username := getVar(r, model.UsernameVar)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
	httpStatus := http.StatusForbidden
	if userExists {
		if isAdmin(db, r) && !user.IsAdmin {
			db.Model(user).Update(model.IsApprovedColumn, !user.IsApproved)
			httpStatus = http.StatusOK
			s.Code = status.SuccessCode
			s.Message = status.UserUpdated
		} else {
			s.Message = status.AdminRequired
		}
	} else {
		s.Message = status.UserNotFound
	}
	WriteData(GetJSON(s), httpStatus, w)
}

func ToggleClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	club := model.NewClub()
	clubID := getVar(r, model.ClubIDVar)
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	httpStatus := http.StatusForbidden
	if clubExists {
		if isAdmin(db, r) {
			db.Model(club).Update(model.ActiveColumn, !club.Active)
			httpStatus = http.StatusOK
			s.Code = status.SuccessCode
			s.Message = status.ClubUpdateSuccess
		} else {
			s.Message = status.AdminRequired
		}
	} else {
		s.Message = status.ClubNotFound
	}
	WriteData(GetJSON(s), httpStatus, w)
}

// In-Progress
func GetToggleUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	users := []*model.User{}
	result := db.Where(model.IsApprovedColumn+"= ?", false).Find(&users)
	if isAdmin(db, r) {
		if result.Error != nil {
			s.Message = status.GetNonToggledUsersFailure
		} else {
			s.Code = status.SuccessCode
			s.Message = status.GetNonToggledUsersSuccess
			s.Data = users
			for _, v := range users {
				fmt.Println(v)
			}
		}
	} else {
		s.Message = status.AdminRequired
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}
