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
	ColumnUsername  = "username"
	ColumnEmail     = "email"
	ColumnPassword  = "password"
	ColumnID        = "id"
	ColumnCreatedAt = "created_at"
	ColumnDeletedAt = "deleted_at"
)

