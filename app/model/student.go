package model

type Student struct {
	*User
	Tags      []Tag   `gorm:"many2many:student_tag;association_foreignkey:ID;foreignkey:ID"`
	Attends   []Event `gorm:"many2many:student_event;association_foreignkey:ID;foreignkey:ID"`
	Swipes    []Club  `gorm:"many2many:student_swipe;association_foreignkey:ID;foreignkey:ID"`
}

func NewStudent() *Student {
	return &Student{User: NewUser(), Tags: []Tag{}, Attends: []Event{}, Swipes: []Club{}}
}

const (
	IsHelpingColumn = "is_helping"
	StudentTable = "student"
)
