package model

import "github.com/2-of-clubs/2ofclubs-server/app/status"

// Basic Event Struct
type Event struct {
	Base
	Name string `validate:"required,min=1,max=50" json:"name"`
	//DateTime    time.Time  `validate:"required,gtetoday,datetime"`
	Description string  `validate:"required,max=300" json:"description"`
	Location    string  `validate:"required,max=100" json:"location"`
	Fee         float64 `validate:"gte=0" json:"fee"`
}

func NewEvent() *Event {
	return &Event{}
}

// Listing out requirements for an event to be successfully created
type EventRequirement struct {
	Admin       string
	Name        string
	Description string
	Location    string
	Fee         string
}

func NewEventRequirement() *EventRequirement {
	return &EventRequirement{
		Admin:       status.ManagerOwnerRequired,
		Name:        status.EventNameConstraint,
		Description: status.EventDescriptionConstraint,
		Location:    status.EventLocationConstraint,
		Fee:         status.EventFeeConstraint,
	}
}

//type EventDisplay struct {
//	ID          uint
//	Name        string
//	Description string
//	Location    string
//	Fee         float64
//}
//
//func (e *Event) Display() EventDisplay {
//	return EventDisplay{
//		ID:          e.ID,
//		Name:        e.Name,
//		Description: e.Description,
//		Location:    e.Location,
//		Fee:         e.Fee,
//	}
//}

const (
	EventIDVar        = "eid"
	EventTable        = "event"
	DateTimeColumn    = "date_time"
	DescriptionColumn = "description"
	LocationColumn    = "location"
	DateFeeColumn     = "fee"
)
