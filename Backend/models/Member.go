package models

type Student struct {
	Person
	IsHelping bool
	Tags
}

func (s *Student) SetIsHelping(isHelping bool) {
	s.IsHelping = isHelping
}
func (s *Student) GetIsHelping() bool {
	return s.IsHelping
}