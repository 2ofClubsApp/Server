package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"gorm.io/gorm"
	"net/http"
)

// ToggleUser - Toggling users as active or inactive
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

// ToggleClub - Toggling clubs as active or inactive
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

// GetToggleClub obtains all clubs that need to be activated (toggled)
func GetToggleClub(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getToggleModel(db, r, s, model.ClubTable)
}

// GetToggleUser obtains all users that need to be activated (toggled)
func GetToggleUser(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getToggleModel(db, r, s, model.UserTable)
}

// Helper function for obtaining all models (users/clubs) that need to be toggled
// This can be extended for future models that require approval
func getToggleModel(db *gorm.DB, r *http.Request, s *status.Status, modelType string) (int, error) {
	if isAdmin(db, r) {
		switch modelType {
		case model.ClubTable:
			var clubs []model.Club
			result := db.Where(model.ActiveColumn+"= ?", false).Find(&clubs)
			if result.Error != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Message = status.GetNonApprovedClubsSuccess
			toggleClubs := []model.ClubBaseInfo{}
			for _, c := range clubs {
				toggleClubs = append(toggleClubs, c.DisplayBaseClubInfo())
			}
			s.Data = toggleClubs
		case model.UserTable:
			var users []model.User
			result := db.Where(model.IsApprovedColumn+"= ?", false).Find(&users)
			if result.Error != nil {
				return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
			}
			s.Message = status.GetNonApprovedUsersSuccess
			toggleUsers := []model.UserBaseInfo{}
			for _, u := range users {
				toggleUsers = append(toggleUsers, u.DisplayBaseUserInfo())
			}
			s.Data = toggleUsers
		}
		s.Code = status.SuccessCode
		return http.StatusOK, nil
	}
	s.Message = status.AdminRequired
	return http.StatusForbidden, nil
}
