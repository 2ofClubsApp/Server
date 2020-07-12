package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string
}
func NewTag() *Tag{
	return &Tag{}
}