package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   `json:"-"`
	*Credentials `json:"-"`
	Manages      []Club `gorm:"many2many:user_club;"`
	Chooses      []Tag  `gorm:"many2many:user_tag;"`
	//Attends      []Event `gorm:"many2many:User_Event;`
	IsAdmin    bool `json:"-"`
	IsApproved bool `json:"-"`
}

func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}, Chooses: []Tag{}}
}

/*
Add IsOwner in Manages
*/
//func (u *User) AfterFind(tx *gorm.DB) error{
//	fmt.Println(tx)
//	return nil
//}

const (
	ChoosesColumn   = "Chooses"
	UserTagTable    = "user_tag"
	UserClubTable   = "user_club"
	IsHelpingColumn = "is_helping"
	UserTable       = "user"
	ManagesColumn   = "Manages"
	IsAdminColumn   = "is_admin"
	IDColumn        = "id"
	CreatedAtColumn = "created_at"
	DeletedAtColumn = "deleted_at"
)
