package model

type Club struct {
	*User
	// Owners/Administrator
	Tags  []Tag   `gorm:"many2many:club_tag;association_foreignkey:ID;foreignkey:ID"`
	Hosts []Event `gorm:"many2many:club_event;association_foreignkey:ID;foreignkey:ID"`
	Size       int
	Bio        string
	HelpNeeded bool
}

func NewClub() *Club {
	return &Club{User: NewUser(), Tags: []Tag{}, Hosts: []Event{}}}

const (
	TagsColumn       = "tags"
	HostsColumn      = "Hosts"
	SizeColumn       = "size"
	BioColumn        = "bio"
	HelpNeededColumn = "help_needed"
	ClubTable        = "club"
)
