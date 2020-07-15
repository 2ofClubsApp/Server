package model

import "gorm.io/gorm"

type User struct {
	gorm.Model   `json:"-"`
	*Credentials `json:"-"`
	Manages      []Club `gorm:"many2many:user_club;"`
	//Tags    []Tag   `gorm:"many2many:student_tag;association_foreignkey:ID;foreignkey:ID"`
	//Attends []Event `gorm:"many2many:student_event;association_foreignkey:ID;foreignkey:ID"`
	IsAdmin bool
}

func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}}
}

const (
	UserClubTable   = "user_club"
	IsHelpingColumn = "is_helping"
	UserTable       = "user"
	IDColumn        = "id"
	CreatedAtColumn = "created_at"
	DeletedAtColumn = "deleted_at"
)
