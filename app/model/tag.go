package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model `json:"-"`
	Name       string `gorm:"UNIQUE" validate:"required,min=1,max=25"`
	IsActive   bool   `json:"-"`
}

type TagDisplay struct {
	ID   uint
	Name string
}

type TagDisplayCollection struct {
	Tags []TagDisplay
}

func NewTag() *Tag {
	return &Tag{}
}

func (t *Tag) Display() *TagDisplay {
	return &TagDisplay{
		ID:   t.ID,
		Name: t.Name,
	}
}

const (
	TagTable       = "tag"
	IsActiveColumn = "is_active"
)
