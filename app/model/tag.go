package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"UNIQUE" validate:"required,min=1,max=25"`
}

func NewTag() *Tag {
	return &Tag{}
}

const (
	TagTable = "tag"
)
