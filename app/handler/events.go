package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
)

func GetAllEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()
	events := []model.Event{}
	result := db.Find(&events)
	if result.Error != nil {
		status.Message = model.GetAllEventsFailure
	} else {
		allEvents := make(map[string][]model.Event)
		allEvents["Events"] = events
		status.Message = model.AllEventsFound
		status.Code = model.SuccessCode
		status.Data = allEvents
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

func GetClubEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Event")
}

/*
Creating an event for a particular club. The user creating the club must at least be a manager
*/
func CreateEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
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
		status.Code = model.SuccessCode
	} else if !clubExists {
		status.Message = model.ClubNotFound
	} else {
		status.Message = model.CreateEventFailure
		status.Data = model.EventStatus{
			Admin:       model.ManagerOwnerRequired,
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

func AttendEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/*
		Ensure that users can't add multi events
	*/
	status := model.NewStatus()
	eventID := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	if eventExists && userExists {
		db.Model(user).Association(model.AttendsColumn).Append(event)
		status.Message = model.EventFound
		status.Code = model.SuccessCode
	} else {
		status.Message = model.EventNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
