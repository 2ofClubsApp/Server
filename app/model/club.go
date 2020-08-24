package model

// Club - Base Club struct
type Club struct {
	Base
	Name    string  `validate:"required,min=3,max=50" json:"name"`
	Email   string  `validate:"required,email" json:"email"`
	Bio     string  `validate:"required,min=1,max=300" json:"bio"`
	Size    int     `validate:"required,gt=0" json:"size"`
	Active  bool    `json:"-"`
	Logo    string  `json:"logo"`
	Sets    []Tag   `gorm:"many2many:club_tag;foreignKey:id;References:Name;" json:"tags"`
	Hosts   []Event `gorm:"many2many:club_event;" json:"hosts"`
	Managed []User  `gorm:"many2many:user_club;" json:"-"`
}

// ClubBaseInfo - Displaying basic club data
type ClubBaseInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ClubUpdateInfo - Payload for updating a club's info
type ClubUpdateInfo struct {
	Size int    `validate:"required,gt=0"`
	Bio  string `validate:"required,max=300"`
}

// NewClub - Create new default Club
func NewClub() *Club {
	return &Club{Sets: []Tag{}, Active: false, Hosts: []Event{}}
}

// NewClubUpdate returns a struct to update club info
func NewClubUpdate() *ClubUpdateInfo {
	return &ClubUpdateInfo{}
}

// DisplayBaseClubInfo displays base club data
func (c *Club) DisplayBaseClubInfo() ClubBaseInfo {
	return ClubBaseInfo{ID: c.ID, Name: c.Name}
}

// Club variables for db columns/route variables
const (
	ManagedClubColumn = "Managed"
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
