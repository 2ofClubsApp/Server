package model

// Base Club struct
type Club struct {
	Base
	Name   string  `validate:"required,min=3,max=50"`
	Email  string  `validate:"required,email"`
	Bio    string  `validate:"required,max=300"`
	Size   int     `validate:"required,gt=0"` // Set > 0 as a restriction
	Active bool    `json:"-"`
	Sets   []Tag   `gorm:"many2many:club_tag;foreignKey:id;References:Name;"`
	Hosts  []Event `gorm:"many2many:club_event;"`
}

// Displaying public club data
type ClubDisplay struct {
	ID    uint     `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Bio   string   `json:"bio"`
	Size  int      `json:"size"`
	Tags  []string `json:"tags"`
	Hosts []Event  `json:"hosts"`
}

func (c *Club) Display() *ClubDisplay {
	return &ClubDisplay{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
		Bio:   c.Bio,
		Size:  c.Size,
	}
}

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
