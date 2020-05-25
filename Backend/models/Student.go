package models

type Student struct {
	Person
	IsHelping bool
	Chats     []Chat  `gorm:"many2many:student_chat;association_foreignkey:ID;foreignkey:ID"`
	Tags      []Tag   `gorm:"many2many:student_tag;association_foreignkey:ID;foreignkey:ID"`
	Attends   []Event `gorm:"many2many:student_event;association_foreignkey:ID;foreignkey:ID"`
	Swipes    []Club  `gorm:"many2many:student_swipe;association_foreignkey:ID;foreignkey:ID"`
}

const (
	ColumnIsHelping = "is_helping"
)
