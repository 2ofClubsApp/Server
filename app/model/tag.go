package model

import (
	"gorm.io/gorm"
	"time"
)

// Tag - Basic Tag Struct
type Tag struct {
	ID        uint           `gorm:"autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `gorm:"primarykey" validate:"required,min=1,max=25" json:"name"`
	IsActive  bool           `json:"isActive"`
}

// NewTag - Create new default Tag
func NewTag() *Tag {
	return &Tag{}
}

// Tag variables for db columns/route variables
const (
	TagNameVar     = "tagName"
	TagTable       = "tag"
	IsActiveColumn = "IsActive"
)
