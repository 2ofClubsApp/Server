package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
)

func GetAllEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	s := status.New()
	events := []model.Event{}
	result := db.Find(&events)
	if result.Error != nil {
		s.Message = status.GetAllEventsFailure
	} else {
		allEvents := make(map[string][]model.Event)
		allEvents["Events"] = events
		s.Message = status.AllEventsFound
		s.Code = status.SuccessCode
		s.Data = allEvents
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}

func GetEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	eventID := getVar(r, model.EventIDVar)
	event := model.NewEvent()
	s := status.New()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	if eventExists {
		s.Code = status.SuccessCode
		s.Message = status.EventFound
		s.Data = event
	} else {
		s.Message = status.EventNotFound
	}
	WriteData(GetJSON(s), http.StatusOK, w)

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
	s := status.New()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eid, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if userExists && eventExists && clubExists && isManager(db, user, club) {
		s.Code = status.SuccessCode
		s.Message = status.EventDeleted
		db.Delete(event)
	} else if !userExists {
		s.Message = status.UserNotFound
	} else if !clubExists {
		s.Message = status.ClubNotFound
	} else if !eventExists {
		s.Message = status.EventNotFound
	} else {
		statusCode = http.StatusForbidden
		s.Message = http.StatusText(http.StatusForbidden)
	}
	WriteData(GetJSON(s), statusCode, w)
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
	httpStatusCode := http.StatusOK
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	validate := validator.New()
	s := status.New()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	switch operation {
	case model.OpCreate:
		if userExists && isManager(db, user, club) && clubExists {
			extractBody(r, event)
			err := validate.Struct(event)
			if err == nil {
				db.Model(club).Association(model.HostsColumn).Append(event)
				s.Code = status.SuccessCode
				s.Message = status.CreateEventSuccess
			} else {
				s.Message = status.CreateEventFailure
				s.Data = model.NewEventRequirement()
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
				s.Code = status.SuccessCode
				s.Message = status.UpdateEventSuccess
			} else {
				s.Message = status.UpdateEventFailure
				s.Data = model.NewEventRequirement()
			}
		} else if !eventExists {
			s.Message = status.EventNotFound
		}
	}
	if !clubExists {
		s.Message = status.ClubNotFound
	} else if !isManager(db, user, club) {
		httpStatusCode = http.StatusForbidden
		s.Message = http.StatusText(httpStatusCode)
	}
	WriteData(GetJSON(s), httpStatusCode, w)
}

func RemoveUserAttendsEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageUserAttends(db, w, r, model.OpRemove)
}

func AddUserAttendsEvent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	manageUserAttends(db, w, r, model.OpAdd)
}

func manageUserAttends(db *gorm.DB, w http.ResponseWriter, r *http.Request, operation string) {
	s := status.New()
	eventID := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if eventExists && userExists {
		switch operation {
		case model.OpAdd:
			db.Model(user).Association(model.AttendsColumn).Append(event)
			s.Code = status.SuccessCode
			s.Message = status.EventFound
		case model.OpRemove:
			db.Model(user).Association(model.AttendsColumn).Delete(event)
			s.Code = status.SuccessCode
			s.Message = status.EventFound
		}
	} else {
		s.Message = status.EventNotFound
	}
	WriteData(GetJSON(s), http.StatusOK, w)
}
