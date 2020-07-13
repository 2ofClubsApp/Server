package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	*Credentials
	//Swipes []Club `gorm:"many2many:student_swipe;association_foreignkey:ID;foreignkey:ID"`
	//Tags    []Tag   `gorm:"many2many:student_tag;association_foreignkey:ID;foreignkey:ID"`
	//Attends []Event `gorm:"many2many:student_event;association_foreignkey:ID;foreignkey:ID"`
	IsAdmin bool
}

func NewUser() *User {
	return &User{Credentials: NewCredentials()}
}

const (
	IsHelpingColumn = "is_helping"
	UserTable       = "user"
	IDColumn        = "id"
	CreatedAtColumn = "created_at"
	DeletedAtColumn = "deleted_at"
)
