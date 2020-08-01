package model

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `gorm:"UNIQUE;primarykey" validate:"required,min=1,max=25"`
	IsActive  bool           `json:"-"`
}

func NewTag() *Tag {
	return &Tag{}
}

const (
	TagTable       = "tag"
	IsActiveColumn = "is_active"
)
