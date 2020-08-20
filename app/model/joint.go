package model

const (
	ChatIDColumn    = "chat_id"
	LogIDColumn     = "log_id"
	IsOwnerColumn   = "is_owner"
	UserIDColumn    = "user_id"
	ClubIDColumn    = "club_id"
	EventIDColumn   = "event_id"
	TagIDColumn     = "tag_id"
	StudentIDColumn = "student_id"
	TagNameColumn   = "tag_name"
)

// Database format for many to many relation with User and Club
// Keeping a record of users to clubs and whether they're an associated owner of the club or not
type UserClub struct {
	UserId  int
	ClubId  int
	IsOwner bool
}

// Create new default UserClub
func NewUserClub() *UserClub {
	return &UserClub{}
}
