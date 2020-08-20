package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
)

// Returning all events from all clubs
func GetAllEvents(db *gorm.DB, _ http.ResponseWriter, _ *http.Request, s *status.Status) (int, error) {
	events := []model.Event{}
	result := db.Find(&events)
	if result.Error != nil {
		s.Message = status.GetAllEventsFailure
		return http.StatusInternalServerError, fmt.Errorf("unable to obtain all events")
	}
	allEvents := make(map[string][]model.Event)
	allEvents["events"] = events
	s.Message = status.AllEventsFound
	s.Code = status.SuccessCode
	s.Data = allEvents
	return http.StatusOK, nil
}

// Obtaining an event from a specific club
func GetEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	eventID := getVar(r, model.EventIDVar)
	event := model.NewEvent()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	if eventExists {
		s.Code = status.SuccessCode
		s.Message = status.EventFound
		s.Data = event
		return http.StatusOK, nil
	}
	s.Message = status.EventNotFound
	return http.StatusNotFound, nil

}

// Deleting a club event
func DeleteClubEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	eid := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	club := model.NewClub()
	event := model.NewEvent()
	user := model.NewUser()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eid, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	fmt.Println(clubExists)
	fmt.Println(eventExists)
	fmt.Println(userExists)
	if userExists && eventExists && clubExists && isManager(db, user, club) {
		res := db.Delete(event)
		if res.Error != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to delete club event")
		}
		s.Code = status.SuccessCode
		s.Message = status.EventDeleted
		return http.StatusOK, nil
	} else if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	} else if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	} else if !eventExists {
		s.Message = status.EventNotFound
		return http.StatusNotFound, nil
	}
	s.Message = http.StatusText(http.StatusForbidden)
	return http.StatusForbidden, nil
}

// Creating an event for a particular club. The user creating the club must at least be a manager
// See model.Event or docs for the event constraints
func CreateClubEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	validate := validator.New()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if userExists {
		if clubExists {
			if isManager(db, user, club) {
				if err := extractBody(r, event); err != nil {
					return http.StatusInternalServerError, fmt.Errorf(err.Error())
				}
				err := validate.Struct(event)
				if err != nil {
					s.Message = status.CreateEventFailure
					s.Data = model.NewEventRequirement()
					return http.StatusUnprocessableEntity, nil
				}
				err = db.Model(club).Association(model.HostsColumn).Append(event)
				if err != nil {
					return http.StatusInternalServerError, fmt.Errorf("unable to obtain club events")
				}
				s.Code = status.SuccessCode
				s.Message = status.CreateEventSuccess
				return http.StatusCreated, nil
			}
			s.Message = http.StatusText(http.StatusForbidden)
			return http.StatusForbidden, nil
		}
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	s.Message = status.UserNotFound
	return http.StatusNotFound, nil
}

// Updating an event for a particular club
func UpdateClubEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	validate := validator.New()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	eventID := getVar(r, model.EventIDVar)
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	if eventExists && clubExists && userExists && isManager(db, user, club) {
		updatedEvent := model.NewEvent()
		if err := extractBody(r, updatedEvent); err != nil {
			return http.StatusInternalServerError, fmt.Errorf(err.Error())
		}
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
	if !clubExists {
		s.Message = status.ClubNotFound
	} else if !isManager(db, user, club) {
		//httpStatusCode = http.StatusForbidden
		//s.Message = http.StatusText(httpStatusCode)
	}
	return 403, nil
}

// Removing a user attended event
func RemoveUserAttendsEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return manageUserAttends(db, r, model.OpRemove, s)
}

// Adding a user attended event
func AddUserAttendsEvent(db *gorm.DB, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return manageUserAttends(db, r, model.OpAdd, s)
}

// Helper function to add or remove club events
func manageUserAttends(db *gorm.DB, r *http.Request, operation string, s *status.Status) (int, error) {
	eventID := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(r)
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if eventExists {
		if userExists {
			switch operation {
			case model.OpAdd:
				err := db.Model(user).Association(model.AttendsColumn).Append(event)
				if err != nil {
					return http.StatusInternalServerError, fmt.Errorf("unable to obtain user attended events")

				}
				s.Code = status.SuccessCode
				s.Message = status.EventAttendSuccess
				return http.StatusOK, nil
			case model.OpRemove:
				err := db.Model(user).Association(model.AttendsColumn).Delete(event)
				if err != nil {
					return http.StatusInternalServerError, fmt.Errorf("unable to delete user attended events")
				}
				s.Code = status.SuccessCode
				s.Message = status.EventUnattendSuccess
				return http.StatusOK, nil
			}
		}
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	s.Message = status.EventNotFound
	return http.StatusNotFound, nil
}
