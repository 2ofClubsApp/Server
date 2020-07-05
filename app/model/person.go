package model

import "github.com/jinzhu/gorm"

type Person struct {
	gorm.Model
	Username string `gorm: "UNIQUE"`
	Email    string `gorm:"UNIQUE"`
	Password string
	//ApiKey   string `gorm:"UNIQUE"`
}

func NewPerson() Person{
	return Person{}
}
const (
	UsernameColumn  = "username"
	EmailColumn   = "email"
	PasswordColumn  = "password"
	IDColumn        = "id"
	CreatedAtColumn = "created_at"
	DeletedAtColumn = "deleted_at"
)

