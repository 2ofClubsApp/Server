package handler

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

const (
	SuccessClubCreation = "Club successfully created"
	FailureClubCreation = "Unable to create the Club"
)

func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs")
}

/*
Check if the email & username is available (RecordExists)

*/
func CreateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/*
		-> Validate User token (Done)
		-> Get Token Claims (User must exist then, unless deleted?) (You can put a check on Record Exists on the deleted column as long as it's null it'll exist then) (Done)
		-> Extract Username and return User struct (Done)
		-> Check if the club is available (Done)
		-> Extract New Club Info to Struct (Done)
		-> Insert to User Manages (Done)
		-> Update user (Done)
		-> Set the user as the owner of the club (Done)
	*/
	claims := GetTokenClaims(r)
	user := model.NewUser()
	uname := fmt.Sprintf("%v", claims["sub"])
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	club := getClubInfo(r)
	validate := validator.New()
	err := validate.Struct(club)
	clubExists := SingleRecordExists(db, model.ClubTable, model.NameColumn, club.Name, model.NewClub())
	emailExists := SingleRecordExists(db, model.ClubTable, model.EmailColumn, club.Email, model.NewClub())
	status := model.NewStatus()
	// Keeping userExists as a check even though the user should exist given the valid token because there's a chance that the user is deleted
	// In this case the user will still exist in the database but will be inaccessible.
	if !emailExists && !clubExists && userExists && err == nil {
		db.Model(user).Association(model.ManagesColumn).Append(club)
		db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", user.ID, club.ID, false).Update(model.IsOwnerColumn, true)
		status.Message = SuccessClubCreation
	} else {
		status.Code = model.FailureCode
		status.Message = FailureClubCreation
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func getClubInfo(r *http.Request) *model.Club {
	decoder := json.NewDecoder(r.Body)
	club := model.NewClub()
	decoder.Decode(club)
	return club
}

func GetClubsTag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs Tag")
}

func GetClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var statusCode int
	var data string
	vars := mux.Vars(r)
	clubName := strings.ToLower(vars["name"])
	status := model.NewStatus()
	c := model.NewClub()
	found := SingleRecordExists(db, model.ClubTable, model.NameColumn, clubName, c)
	if !found {
		status.Message = model.ClubNotFound
		status.Code = model.FailureCode
	} else {
		status.Data = c
		status.Message = model.ClubFound
	}
	statusCode = http.StatusOK
	data = GetJSON(status)
	WriteData(data, statusCode, w)
}

func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Club")
}

func DeleteClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	vars := mux.Vars(r)
	clubName := vars["name"]
	club := model.NewClub()
	if SingleRecordExists(db, model.ClubTable, model.NameColumn, clubName, club) {
		claims := GetTokenClaims(r)
		uname := fmt.Sprintf("%v", claims["sub"])
		user := model.NewUser()
		SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
		if isOwner(db, user, club) || isAdmin(db, r) {
			db.Model(user).Association(model.ManagesColumn).Delete(club)
			db.Delete(club)
			status.Message = model.SuccessClubDelete
		} else {
			status.Code = -1
			status.Message = model.FailureClubDelete
		}
	} else {
		status.Code = -1
		status.Message = model.FailureClubDelete
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Returns true, if the user is the owner of the club and false otherwise
*/
func isOwner(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).Find(userClub)
	return userClub.IsOwner
}

func AddManager(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	claims := GetTokenClaims(r)
	clubOwnerUsername := fmt.Sprintf("%v", claims["sub"])
	vars := mux.Vars(r)
	managerUsername := vars["username"]
	clubname := vars["clubname"]
	owner := model.NewUser()
	manager := model.NewUser()
	club := model.NewClub()
	// Added user must exist
	SingleRecordExists(db, model.UserTable, model.UsernameColumn, clubOwnerUsername, owner)
	// If owner is found, then the owner struct isn't populated, which gives ID=0, but ID's start at 1, so this shouldn't cause any potential security issues
	managerExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, managerUsername, manager)
	clubExists := SingleRecordExists(db, model.ClubTable, model.NameColumn, clubname, club)
	if managerExists && clubExists {
		if isOwner(db, owner, club) && owner.Username != manager.Username {
			db.Model(manager).Association(model.ManagesColumn).Append(club)
			status.Message = model.SuccessManagerAddition
		} else {
			status.Message = model.FailureManagerAddition
			status.Code = model.FailureCode
		}
	} else {
		status.Code = model.FailureCode
		status.Message = model.FailureManagerAddition
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
