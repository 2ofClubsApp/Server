package model

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	Message string
}
func NewLog() *Log{
	return &Log{}
}
const (
	MessageColumn = "message"
)