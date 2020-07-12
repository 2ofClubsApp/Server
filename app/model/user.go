package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	// Regex's not added in metadata since it can conflict with commas and other symbols
	Username  string `gorm:"UNIQUE" validate:"alpha,min=2,max=15,required"`
	Email     string `gorm:"UNIQUE" validate:"required,email"`
	Password  string `validate:"required,min=3,max=45"`
	// Max 45 due to 50 length limitation of bcrypt
}

func NewUser() *User {
	return &User{}
}


const (
	UsernameColumn  = "username"
	EmailColumn     = "email"
	PasswordColumn  = "password"
	IDColumn        = "id"
	CreatedAtColumn = "created_at"
	DeletedAtColumn = "deleted_at"
)
