package model

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model `json:"-"`
	Name       string `validate:"required,max=50"`
	//DateTime    time.Time  `validate:"required,gtetoday,datetime"`
	Description string  `validate:"required,max=300"`
	Location    string  `validate:"required,max=100"`
	Fee         float64 `validate:"required,gte=0.0"`
}

func NewEvent() *Event {
	return &Event{}
}

type EventDisplay struct {
	ID          uint
	Name        string
	Description string
	Location    string
	Fee         float64
}

func (e *Event) Display() EventDisplay {
	return EventDisplay{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Location:    e.Location,
		Fee:         e.Fee,
	}
}

const (
	DateTimeColumn    = "date_time"
	DescriptionColumn = "description"
	LocationColumn    = "location"
	DateFeeColumn     = "fee"
)
