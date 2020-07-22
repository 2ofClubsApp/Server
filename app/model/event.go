package model

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model  `json:"-"`
	DateTime    time.Time  `validate:"required,gtetoday,datetime"`
	Description string  `validate:"required"`
	Location    string  `validate:"required"`
	Fee         float64 `validate:"required"`
}

func NewEvent() *Event {
	return &Event{}
}

const (
	DateTimeColumn    = "date_time"
	DescriptionColumn = "description"
	LocationColumn    = "location"
	DateFeeColumn     = "fee"
)
