package model

import (
	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	DateTime    string
	Description string
	Location    string
	Fee         float64
}

func NewEvent() *Event{
	return &Event{}
}


const (
	DateTimeColumn = "date_time"
	DescriptionColumn= "description"
	LocationColumn = "location"
	DateFeeColumn = "fee"
)
