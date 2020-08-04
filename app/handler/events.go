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

func GetEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	eventID := getVar(r, model.EventIDVar)
	event := model.NewEvent()
	status := model.NewStatus()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	if eventExists {
		status.Code = model.SuccessCode
		status.Message = model.EventFound
		status.Data = event
	} else {
		status.Message = model.EventNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)

}

func DeleteClubEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	clubID := getVar(r, model.ClubIDVar)
	eid := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	club := model.NewClub()
	event := model.NewEvent()
	user := model.NewUser()
	status := model.NewStatus()
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eid, event)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	if userExists && eventExists && clubExists && isManager(db, user, club) {
		status.Code = model.SuccessCode
		status.Message = model.EventDeleted
		db.Delete(event)
	} else if !userExists {
		status.Message = model.UserNotFound
	} else if !clubExists {
		status.Message = model.ClubNotFound
	} else if !eventExists {
		status.Message = model.EventNotFound
	} else {
		statusCode = http.StatusForbidden
		status.Message = http.StatusText(http.StatusForbidden)
	}
	WriteData(GetJSON(status), statusCode, w)
}

/*
Creating an event for a particular club. The user creating the club must at least be a manager
*/
func CreateClubEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageClubEvent(db, w, r, model.OpCreate)
}

func UpdateClubEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageClubEvent(db, w, r, model.OpUpdate)
}

func manageClubEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request, operation string) {
	claims := GetTokenClaims(r)
	httpStatus := http.StatusOK
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	validate := validator.New()
	status := model.NewStatus()
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	switch operation {
	case model.OpCreate:
		if userExists && isManager(db, user, club) && clubExists {
			extractBody(r, event)
			err := validate.Struct(event)
			if err == nil {
				db.Model(club).Association(model.HostsColumn).Append(event)
				status.Code = model.SuccessCode
				status.Message = model.CreateEventSuccess
			} else {
				status.Message = model.CreateEventFailure
				status.Data = model.NewEventRequirement()
			}
		}
	case model.OpUpdate:
		eventID := getVar(r, model.EventIDVar)
		eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
		if eventExists && clubExists && userExists && isManager(db, user, club) {
			updatedEvent := model.NewEvent()
			extractBody(r, updatedEvent)
			err := validate.Struct(updatedEvent)
			if err == nil {
				db.Model(event).Updates(updatedEvent)
				status.Code = model.SuccessCode
				status.Message = model.UpdateEventSuccess
			} else {
				status.Message = model.UpdateEventFailure
				status.Data = model.NewEventRequirement()
			}
		} else if !eventExists {
			status.Message = model.EventNotFound
		}
	}
	if !clubExists {
		status.Message = model.ClubNotFound
	} else if !isManager(db, user, club) {
		httpStatus = http.StatusForbidden
		status.Message = http.StatusText(httpStatus)
	}
	WriteData(GetJSON(status), httpStatus, w)
}

func RemoveUserAttendsEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageUserAttends(db, w, r, model.OpRemove)
}

func AddUserAttendsEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageUserAttends(db, w, r, model.OpAdd)
}

func manageUserAttends(db *gorm.DB, w http.ResponseWriter, r *http.Request, operation string) {
	status := model.NewStatus()
	eventID := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	if eventExists && userExists {
		switch operation {
		case model.OpAdd:
			db.Model(user).Association(model.AttendsColumn).Append(event)
			status.Code = model.SuccessCode
			status.Message = model.EventFound
		case model.OpRemove:
			db.Model(user).Association(model.AttendsColumn).Delete(event)
			status.Code = model.SuccessCode
			status.Message = model.EventFound
		}
	} else {
		status.Message = model.EventNotFound
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
