package model

import (
	"gorm.io/gorm"
	"time"
)

// Basic Tag Struct
type Tag struct {
	ID        uint           `gorm:"autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `gorm:"primarykey" validate:"required,min=1,max=25" json:"name"`
	IsActive  bool           `json:"isActive"`
}

func NewTag() *Tag {
	return &Tag{}
}

const (
	TagNameVar     = "tagName"
	TagTable       = "tag"
	IsActiveColumn = "IsActive"
)
