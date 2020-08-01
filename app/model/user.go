package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   `json:"-"`
	*Credentials `json:"-"`
	Manages      []Club  `gorm:"many2many:user_club;"`
	Chooses      []Tag   `gorm:"many2many:user_tag;foreignKey:id;References:Name" json:"-"`
	Attends      []Event `gorm:"many2many:user_event;"`
	IsAdmin      bool    `json:"-"`
	IsApproved   bool    `json:"-"`
}

type UserDisplay struct {
	Email   string
	Manages []*ManagesDisplay
	Tags    []string
	Attends []EventDisplay
}

type ManagesDisplay struct {
	*ClubDisplay
	IsOwner bool
}

func (u *User) Display() *UserDisplay {
	return &UserDisplay{Email: u.Email}
}

func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}, Chooses: []Tag{}, Attends: []Event{}}
}

const (
	ChoosesColumn    = "Chooses"
	AttendsColumn    = "Attends"
	UserTagTable     = "user_tag"
	UserClubTable    = "user_club"
	IsHelpingColumn  = "is_helping"
	UserTable        = "user"
	ManagesColumn    = "Manages"
	IsAdminColumn    = "is_admin"
	IsApprovedColumn = "is_approved"
	IDColumn         = "id"
	CreatedAtColumn  = "created_at"
	DeletedAtColumn  = "deleted_at"
)
