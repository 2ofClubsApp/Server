package model

import (
	"github.com/jinzhu/gorm"
)

type Chat struct {
	gorm.Model
	Logs []Log `gorm:"many2many:chat_log"`
}

func NewChat() *Chat{
	return &Chat{Logs: []Log{}}
}