package models

type Club struct {
	Person
	Tags       []Tag   `gorm:"many2many:club_tag;association_foreignkey:ID;foreignkey:ID"`
	Hosts      []Event `gorm:"many2many:club_event;association_foreignkey:ID;foreignkey:ID"`
	Chats      []Chat  `gorm:"many2many:club_chat;association_foreignkey:ID;foreignkey:ID"`
	Size       int
	Bio        string
	HelpNeeded bool
}

func NewClub() *Club {
	return &Club{Person: NewPerson(), Tags: []Tag{}, Hosts: []Event{}, Chats: []Chat{}}
}

const (
	ColumnSize       = "size"
	ColumnBio        = "bio"
	ColumnHelpNeeded = "help_needed"
)
