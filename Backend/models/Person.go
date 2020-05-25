package models

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
	Username string `gorm: "UNIQUE"`
	Email    string `gorm:"UNIQUE"`
	Password string
}

const (
	ColumnUsername  = "username"
	ColumnEmail     = "email"
	ColumnPassword  = "password"
)
