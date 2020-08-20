package model

// Base Club struct
type Club struct {
	Base
	Name   string  `validate:"required,min=3,max=50" json:"name"`
	Email  string  `validate:"required,email" json:"email"`
	Bio    string  `validate:"required,max=300" json:"bio"`
	Size   int     `validate:"required,gt=0" json:"size"` // Set > 0 as a restriction
	Active bool    `json:"-"`
	Sets   []Tag   `gorm:"many2many:club_tag;foreignKey:id;References:Name;" json:"tags"`
	Hosts  []Event `gorm:"many2many:club_event;" json:"hosts"`
}

// Create new default Club
func NewClub() *Club {
	return &Club{Sets: []Tag{}, Hosts: []Event{}, Active: false}
}

const (
	ClubIDVar         = "cid"
	AllClubInfo       = "all"
	AllClubEventsHost = "events"
	ClubTagTable      = "club_tag"
	SetsColumn        = "Sets"
	TagsColumn        = "tags"
	HostsColumn       = "Hosts"
	SizeColumn        = "size"
	BioColumn         = "bio"
	HelpNeededColumn  = "help_needed"
	ClubTable         = "club"
	NameColumn        = "name"
	ActiveColumn      = "active"
)
