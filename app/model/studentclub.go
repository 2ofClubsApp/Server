package model

type StudentClub struct {
	StudentID int
	ClubID    int
	IsOwner   bool
}

func NewStudentClub() *StudentClub {
	return &StudentClub{}
}

const (
	StudentClubTable = "student_club"
)
