package model

import "gorm.io/gorm"

type Club struct {
	// Owners/Administrator
	gorm.Model `json:"-"`
	Name       string `validate:"required,min=3,max=50"`
	Email      string `validate:"required,email"`
	Bio        string `validate:"required,max=300"`
	Size       int    `validate:"required,gt=0"` // Set > 0 as a restriction
	Active     bool   `json:"-"`
	Sets       []Tag  `gorm:"many2many:club_tag;foreignKey:id;References:Name;"`
	//Hosts      []Event `gorm:"many2many:club_event;association_foreignkey:ID;foreignkey:ID"`
}

type ClubDisplay struct {
	ID    uint
	Name  string
	Email string
	Bio   string
	Size  int
	Tags  []string
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
