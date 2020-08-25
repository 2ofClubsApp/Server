package model

import (
	"github.com/2-of-clubs/2ofclubs-server/app/status"
	"time"
)

// Event - Basic Event Struct
type Event struct {
	Base
	Name        string    `validate:"required,min=1,max=50" json:"name"`
	DateTime    time.Time `validate:"required" json:"datetime"`
	Description string    `validate:"required,max=300" json:"description"`
	Location    string    `validate:"required,max=300" json:"location"`
	Fee         float64   `validate:"gte=0" json:"fee"`
}

// NewEvent - Create new default Event
func NewEvent() *Event {
	return &Event{}
}

// EventRequirement - Listing out requirements for an event to be successfully created
type EventRequirement struct {
	Admin       string `json:"admin"`
	Name        string `json:"name"`
	DateTime    string `json:"datetime"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Fee         string `json:"fee"`
}

// NewEventRequirement returns the requirements when creating an event
// See model.event or docs for event constraints
func NewEventRequirement() *EventRequirement {
	return &EventRequirement{
		Admin:       status.ManagerOwnerRequired,
		Name:        status.EventNameConstraint,
		DateTime:    status.EventDateTimeConstraint,
		Description: status.EventDescriptionConstraint,
		Location:    status.EventLocationConstraint,
		Fee:         status.EventFeeConstraint,
	}
}

// Event variables for db columns/route variables
const (
	DateTimeColumn = "date_time"
	EventIDVar     = "eid"
	EventTable     = "event"
)
