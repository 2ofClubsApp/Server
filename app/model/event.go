package model

import "github.com/2-of-clubs/2ofclubs-server/app/status"

// Event - Basic Event Struct
type Event struct {
	Base
	Name string `validate:"required,min=1,max=50" json:"name"`
	//DateTime    time.Time  `validate:"required,gtetoday,datetime"`
	Description string  `validate:"required,max=300" json:"description"`
	Location    string  `validate:"required,max=100" json:"location"`
	Fee         float64 `validate:"gte=0" json:"fee"`
}

// NewEvent - Create new default Event
func NewEvent() *Event {
	return &Event{}
}

// EventRequirement - Listing out requirements for an event to be successfully created
type EventRequirement struct {
	Admin       string
	Name        string
	Description string
	Location    string
	Fee         string
}

// NewEventRequirement returns the requirements when creating an event
// See model.event or docs for event constraints
func NewEventRequirement() *EventRequirement {
	return &EventRequirement{
		Admin:       status.ManagerOwnerRequired,
		Name:        status.EventNameConstraint,
		Description: status.EventDescriptionConstraint,
		Location:    status.EventLocationConstraint,
		Fee:         status.EventFeeConstraint,
	}
}

// Event variables for db columns/route variables
const (
	EventIDVar = "eid"
	EventTable = "event"
)
