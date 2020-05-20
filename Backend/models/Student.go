package models

type Student struct {
	Person
	IsHelping bool
	Tags      []Tag   `gorm:"many2many:student_tag"`
	Attends   []Event `gorm:"many2many:student_event"`
}

func (s *Student) SetIsHelping(isHelping bool) {
	s.IsHelping = isHelping
}
func (s *Student) GetIsHelping() bool {
	return s.IsHelping
}
