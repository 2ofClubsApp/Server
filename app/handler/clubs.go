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
	status := model.NewStatus()
	// Keeping userExists as a check even though the user should exist given the valid token because there's a chance that the user is deleted
	// In this case the user will still exist in the database but will be inaccessible.
	if !clubExists && userExists && err == nil {
		db.Model(user).Association("Manages").Append(club)
		db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", user.ID, club.ID, false).Update("is_owner", true)
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
	clubName := strings.ToLower(vars[model.NameColumn])
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
