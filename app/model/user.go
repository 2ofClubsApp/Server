package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model   `json:"-"`
	*Credentials `json:"-"`
	Manages      []Club `gorm:"many2many:user_club;"`
	Chooses      []Tag   `gorm:"many2many:user_tag;foreignKey:id;References:Name" json:"-"`
	Attends      []Event `gorm:"many2many:user_event;"`
	IsAdmin      bool    `json:"-"`
	IsApproved   bool    `json:"-"`
}

type UserDisplay struct {
	Manages []*ManagesDisplay
	Tags    []string
}

type ManagesDisplay struct {
	*ClubDisplay
	IsOwner bool
}

func (u *User) Display() *UserDisplay {
	return &UserDisplay{}
}

func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}, Chooses: []Tag{}}
}

//func (u *User) AfterFind(db *gorm.DB) error {
//var clubsList []Club
//for _, club := range u.Manages {
//	db.Table(ClubTable).Preload(SetsColumn).Find(club)
//	clubsList = append(clubsList, club)
//}
//u.Manages = clubsList
//db.Table(UserTable).Preload(ChoosesColumn).Find(u)
//return nil
//}

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
