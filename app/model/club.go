package model

// Club - Base Club struct
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

// NewClub - Create new default Club
func NewClub() *Club {
	return &Club{Sets: []Tag{}, Hosts: []Event{}, Active: false}
}

// Club variables for db columns/route variables
const (
	ClubIDVar         = "cid"
	AllClubInfo       = "all"
	AllClubEventsHost = "events"
	ClubTagTable      = "club_tag"
	SetsColumn        = "Sets"
	HostsColumn       = "Hosts"
	SizeColumn        = "size"
	BioColumn         = "bio"
	ClubTable         = "club"
	NameColumn        = "name"
	ActiveColumn      = "active"
)
