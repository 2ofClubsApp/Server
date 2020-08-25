package handler

import (
	"fmt"
	"github.com/2ofClubsApp/2ofClubs-Server/app/model"
	"github.com/2ofClubsApp/2ofClubs-Server/app/status"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// GetAllEvents - Returning all events from all clubs
func GetAllEvents(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, _ *http.Request, s *status.Status) (int, error) {
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

// GetEvent - Obtaining an event from a specific club
func GetEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	eventID := getVar(r, model.EventIDVar)
	event := model.NewEvent()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	if !eventExists {
		s.Message = status.EventNotFound
		return http.StatusNotFound, nil
	}
	s.Code = status.SuccessCode
	s.Message = status.EventFound
	s.Data = event
	return http.StatusOK, nil

}

// DeleteClubEvent - Deleting a club event
func DeleteClubEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	clubID := getVar(r, model.ClubIDVar)
	eid := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(ExtractToken(r))
	uname := fmt.Sprintf("%v", claims["sub"])
	club := model.NewClub()
	event := model.NewEvent()
	user := model.NewUser()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eid, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !eventExists {
		s.Message = status.EventNotFound
		return http.StatusNotFound, nil
	}
	if !isManager(db, user, club) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	res := db.Delete(event)
	if res.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to delete club event")
	}
	s.Code = status.SuccessCode
	s.Message = status.EventDeleted
	return http.StatusOK, nil
}

// CreateClubEvent - Creating an event for a particular club. The user creating the club must at least be a manager
// See model.Event or docs for the event constraints
func CreateClubEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(ExtractToken(r))
	uname := fmt.Sprintf("%v", claims["sub"])
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	user := model.NewUser()
	event := model.NewEvent()
	validate := validator.New()
	clubExists := IsSingleRecordActive(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !isManager(db, user, club) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	if err := extractBody(r, event); err != nil {
		return http.StatusInternalServerError, fmt.Errorf(err.Error())
	}
	err := validate.Struct(event)
	if !isValidDate(event.DateTime.Format(time.RFC3339)) || err != nil {
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

func isValidDate(datetime string) bool {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return false
	}
	fmt.Println(datetime)
	t, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return false
	}
	eventTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
	now := time.Now().In(location)
	return eventTime.After(now)

}

// UpdateClubEvent - Updating an event for a particular club
func UpdateClubEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	claims := GetTokenClaims(ExtractToken(r))
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
	if !eventExists {
		s.Message = status.EventNotFound
		return http.StatusNotFound, nil
	}
	if !clubExists {
		s.Message = status.ClubNotFound
		return http.StatusNotFound, nil
	}
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	if !isManager(db, user, club) {
		s.Message = http.StatusText(http.StatusForbidden)
		return http.StatusForbidden, nil
	}
	updatedEvent := model.NewEvent()
	if err := extractBody(r, updatedEvent); err != nil {
		return http.StatusInternalServerError, fmt.Errorf(err.Error())
	}
	err := validate.Struct(updatedEvent)
	validDate := isValidDate(updatedEvent.DateTime.Format(time.RFC3339))
	if !validDate || err != nil {
		s.Message = status.UpdateEventFailure
		s.Data = model.NewEventRequirement()
		return http.StatusUnprocessableEntity, nil
	}
	if db.Model(event).Select(model.NameColumn, model.DescriptionColumn, model.LocationColumn, model.FeeColumn, model.DateTimeColumn).Updates(updatedEvent).Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("unable to update event")
	}
	s.Code = status.SuccessCode
	s.Message = status.UpdateEventSuccess
	return http.StatusOK, nil
}

// RemoveUserAttendsEvent - Removing a user attended event
func RemoveUserAttendsEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return manageUserAttends(db, r, model.OpRemove, s)
}

// AddUserAttendsEvent - Adding a user attended event
func AddUserAttendsEvent(db *gorm.DB, _ *redis.Client, _ http.ResponseWriter, r *http.Request, s *status.Status) (int, error) {
	return manageUserAttends(db, r, model.OpAdd, s)
}

// Helper function to add or remove club events
func manageUserAttends(db *gorm.DB, r *http.Request, operation string, s *status.Status) (int, error) {
	eventID := getVar(r, model.EventIDVar)
	claims := GetTokenClaims(ExtractToken(r))
	uname := fmt.Sprintf("%v", claims["sub"])
	event := model.NewEvent()
	user := model.NewUser()
	eventExists := SingleRecordExists(db, model.EventTable, model.IDColumn, eventID, event)
	userExists := IsSingleRecordActive(db, model.UserTable, model.UsernameColumn, uname, user)
	if !eventExists {
		s.Message = status.EventNotFound
		return http.StatusNotFound, nil
	}
	if !userExists {
		s.Message = status.UserNotFound
		return http.StatusNotFound, nil
	}
	switch operation {
	case model.OpAdd:
		err := db.Model(user).Association(model.AttendsColumn).Append(event)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to obtain user attended events")

		}
		s.Message = status.EventAttendSuccess
	case model.OpRemove:
		err := db.Model(user).Association(model.AttendsColumn).Delete(event)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("unable to delete user attended events")
		}
		s.Message = status.EventUnattendSuccess
	}
	s.Code = status.SuccessCode
	return http.StatusOK, nil
}
