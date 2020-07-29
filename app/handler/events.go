package handler

import (
	"../model"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
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
	if userExists && isManager(db, user, club) && clubExists && err == nil{
		db.Model(club).Association(model.HostsColumn).Append(event)
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


func getVar(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}

func UpdateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Event")
}

func DeleteEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete Event")
}
