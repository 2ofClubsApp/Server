package model

const (
	ChatIDColumn = "chat_id"
	LogIDColumn = "log_id"
	IsOwnerColumn   = "is_owner"
	UserIDColumn    = "user_id"
	ClubIDColumn = "club_id"
	EventIDColumn = "event_id"
	TagIDColumn = "tag_id"
	StudentIDColumn = "student_id"
)

type UserClub struct {
	UserId  int
	ClubID  int
	IsOwner bool
}

func NewUserClub() *UserClub {
	return &UserClub{}
}
