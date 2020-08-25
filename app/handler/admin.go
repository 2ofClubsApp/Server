package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// ToggleUser - Toggling users as active or inactive
func ToggleUser(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	user := model.NewUser()
	username := strings.ToLower(getVar(r, model.UsernameVar))
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	// Must be an admin but the admin cannot toggle themselves off
	if !(isAdmin(db, r) && !user.IsAdmin) {
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
	res := db.Model(user).Update(model.IsApprovedColumn, !user.IsApproved)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	s.Code = status.SuccessCode
	s.Message = status.ToggleUserSuccess
	return http.StatusOK, nil
}

// ToggleClub - Toggling clubs as active or inactive
func ToggleClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	club := model.NewClub()
	clubID := getVar(r, model.ClubIDVar)
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !isAdmin(db, r) {
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
	res := db.Model(club).Update(model.ActiveColumn, !club.Active)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	s.Code = status.SuccessCode
	s.Message = status.ClubToggleSuccess
	return http.StatusOK, nil
}

// GetClubPreview obtains a preview a non active club
func GetClubPreview(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if !isAdmin(db, r) {
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if club.Active {
		s.Message = status.ClubAlreadyActive
		return http.StatusNotFound, nil
	}
	s.Code = status.SuccessCode
	s.Message = status.ClubFound
	s.Data = club // Obtaining extra club info won't be needed as you can't add any tags/events without activating the club
	return http.StatusOK, nil
}

// GetToggleClub obtains all clubs that need to be activated (toggled)
func GetToggleClub(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getToggleModel(db, r, s, model.ClubTable)
}

// GetToggleUser obtains all users that need to be activated (toggled)
func GetToggleUser(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return getToggleModel(db, r, s, model.UserTable)
}

// Helper function for obtaining all models (users/clubs) that need to be toggled
// This can be extended for future models that require approval
func getToggleModel(db *gorm.DB, r *http.Request, s *status.Status, modelType string) (int, error) {
	if !isAdmin(db, r) {
		s.Message = status.AdminRequired
		return http.StatusForbidden, nil
	}
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

// GetClubManagers returns all club managers (not including the club owner)
func GetClubManagers(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	club := model.NewClub()
	clubID := getVar(r, model.ClubIDVar)
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	clubManagers := []model.UserBaseInfo{}
	if db.Table(model.ClubTable).Preload(model.ManagedClubColumn).Find(club).Error != nil {
		return http.StatusInternalServerError, fmt.Errorf(http.StatusText(http.StatusInternalServerError))
	}
	for _, user := range club.Managed {
		if !isOwner(db, &user, club) {
			clubManagers = append(clubManagers, user.DisplayBaseUserInfo())
		}
	}
	s.Code = status.SuccessCode
	s.Message = status.GetClubManagerSuccess
	s.Data = clubManagers
	return http.StatusOK, nil
}
