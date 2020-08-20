package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"gorm.io/gorm"
	"net/http"
)

// Toggling users as active or inactive
func ToggleUser(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	user := model.NewUser()
	username := getVar(r, model.UsernameVar)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
	if userExists {
		if isAdmin(db, r) && !user.IsAdmin {
			res := db.Model(user).Update(model.IsApprovedColumn, !user.IsApproved)
			if res.Error != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Code = status.SuccessCode
			s.Message = status.ToggleUserSuccess
			return http.StatusOK, nil
		}
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
	s.Message = status.UserNotFound
	return http.StatusNotFound, nil
}

// Toggling clubs as active or inactive
func ToggleClub(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	club := model.NewClub()
	clubID := getVar(r, model.ClubIDVar)
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if clubExists {
		if isAdmin(db, r) {
			res := db.Model(club).Update(model.ActiveColumn, !club.Active)
			if res.Error != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Code = status.SuccessCode
			s.Message = status.ClubToggleSuccess
			return http.StatusOK, nil
		}
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
	s.Message = status.ClubNotFound
	return http.StatusNotFound, nil
}

// Obtaining all users that need to be activated (toggled)
func GetToggleUser(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	users := []*model.User{}
	result := db.Where(model.IsApprovedColumn+"= ?", false).Find(&users)
	if isAdmin(db, r) {
		if result.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
		}
		s.Code = status.SuccessCode
		s.Message = status.GetNonToggledUsersSuccess
		s.Data = users
		for _, v := range users {
			fmt.Println(v)
		}
		return http.StatusOK, nil
	}
	s.Message = status.AdminRequired
	return http.StatusForbidden, nil
}
