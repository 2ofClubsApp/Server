package model

import "gorm.io/gorm"

type Club struct {
	//*User Infinite Loop Error here
	// Owners/Administrator
	gorm.Model `json:"-"`
	Email      string `gorm:"UNIQUE" validate:"required,email"`
	Bio        string `validate:"required,max=300"`
	Size       int    `validate:"required"`
	Name       string `gorm:"UNIQUE" validate:"required,max=50"`
	Approved   bool
	//Tags       []Tag   `gorm:"many2many:club_tag;association_foreignkey:ID;foreignkey:ID"`
	//Hosts      []Event `gorm:"many2many:club_event;association_foreignkey:ID;foreignkey:ID"`
	//HelpNeeded bool
}

func NewClub() *Club {
	return &Club{}
}

const (
	TagsColumn       = "tags"
	HostsColumn      = "Hosts"
	SizeColumn       = "size"
	BioColumn        = "bio"
	HelpNeededColumn = "help_needed"
	ClubTable        = "club"
	NameColumn       = "name"
)
