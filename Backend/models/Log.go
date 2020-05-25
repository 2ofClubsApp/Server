package models

import "github.com/jinzhu/gorm"

type Log struct {
	gorm.Model
	Message string
}

const (
	ColumnMessage = "message"
)