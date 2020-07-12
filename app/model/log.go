package model

import "gorm.io/gorm"

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