package handler

import (
	"../model"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
)

const (
	SuccessClubCreation = "Club successfully created"
	FailureClubCreation = "Unable to create the Club"
)

func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Clubs")
}

func CreateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/*
		-> Validate User token (Done)
		-> Get Token Claims (User must exist then, unless deleted?) (You can put a check on Record Exists on the deleted column as long as it's null it'll exist then) (Done)
		-> Extract Username and return User struct (Done)
		-> Check if the club is available (Done)
		-> Extract New Club Info to Struct (Done)
		-> Insert to User Manages (Done)
		-> Update user (Done)
		-> Set the user as the owner of the club (Not done yet)
	*/
	claims := GetTokenClaims(r)
	user := model.NewUser()
	uname := fmt.Sprintf("%v", claims["sub"])
	userExists := RecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	club := getClubInfo(r)
	validate := validator.New()
	err := validate.Struct(club)
	clubExists := RecordExists(db, model.ClubTable, model.NameColumn, club.Name, model.NewClub())
	status := model.NewStatus()
	// Keeping userExists as a check even though the user should exist given the valid token because there's a chance that the user is deleted
	// In this case the user will still exist in the database but will be inaccessible.
	if !clubExists && userExists && err == nil {
		user.Manages = append(user.Manages, *club)
		db.Table(model.UserTable).Updates(user)
		// Add setting the user as the owner of the club
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
	fmt.Println("Getting club")
}

func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Club")
}
