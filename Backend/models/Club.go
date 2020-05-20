package models

type Club struct {
	Person
	Tags       []Tag `gorm:"many2many:club_tag"`
	Hosts      []Event `gorm:"many2many:club_event"`
	Size       int
	Bio        string
	HelpNeeded bool
}

func (c *Club) GetSize() int {
	return c.Size
}

func (c *Club) SetSize(size int) {
	c.Size = size
}

func (c *Club) GetBio() string {
	return c.Bio
}

func (c *Club) SetBio(bio string) {
	c.Bio = bio
}

func (c *Club) isHelpNeeded() bool {
	return c.HelpNeeded
}

func (c *Club) setHelpNeeded(helpNeeded bool) {
	c.HelpNeeded = helpNeeded
}
