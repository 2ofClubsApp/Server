package handler

import (
	"../model"
	"fmt"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
)

func GetEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Events")
}

func GetEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Event")
}

/*
Creating an event for a particular club. The user creating the club must at least be a manager
*/
func CreateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, "cid")
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	extractBody(r, event)
	validate := validator.New()
	err := validate.Struct(event)
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	status := model.NewStatus()
	if userExists && isManager(db, user, club) && clubExists && err == nil {
		db.Model(club).Association(model.HostsColumn).Append(event)
		status.Message = model.CreateEventSuccess
	} else if !clubExists {
		status.Code = model.FailureCode
		status.Message = model.ClubNotFound
	} else {
		status.Code = model.FailureCode
		status.Message = model.CreateEventFailure
		status.Data = model.EventStatus{
			Name:        model.EventNameConstraint,
			Description: model.EventDescriptionConstraint,
			Location:    model.EventLocationConstraint,
			Fee:         model.EventFeeConstraint,
		}
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func UpdateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Event")
}

func DeleteEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Event")
}

func AttendEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/*
	Ensure that users can't add multi events
	 */
	status := model.NewStatus()
	eventID := getVar(r, "eid")
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	if eventExists && userExists {
		db.Model(user).Association(model.AttendsColumn).Append(event)
		status.Message = model.EventFound
	} else {
		status.Code = model.FailureCode
		status.Message = model.EventNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
