package model

import "gorm.io/gorm"

type Club struct {
	//*User Infinite Loop Error here
	// Owners/Administrator
	gorm.Model `json:"-"`
	Name       string `validate:"required,max=50"`
	Email      string `validate:"required,email"`
	Bio        string `validate:"required,max=300"`
	Size       int    `validate:"required"`
	Active     bool   `json:"-"`
	Sets       []Tag  `gorm:"many2many:club_tag;"`
	//Hosts      []Event `gorm:"many2many:club_event;association_foreignkey:ID;foreignkey:ID"`
	//HelpNeeded bool
}

type ClubDisplay struct {
	Name  string
	Email string
	Bio   string
	Size  int
	Tags  []string
}

func (c *Club) Display() *ClubDisplay {
	return &ClubDisplay{
		Name:  c.Name,
		Email: c.Email,
		Bio:   c.Bio,
		Size:  c.Size,
	}
}

func NewClub() *Club {
	return &Club{Sets: []Tag{}}
}

const (
	SetsColumn       = "Sets"
	TagsColumn       = "tags"
	HostsColumn      = "Hosts"
	SizeColumn       = "size"
	BioColumn        = "bio"
	HelpNeededColumn = "help_needed"

	ClubTable  = "club"
	NameColumn = "name"
	OpAdd      = "ADD"
	OpRemove   = "REMOVE"
)
