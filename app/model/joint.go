package model

// IsOwnerColumn represents the is_owner column in the UserClub joint table
var IsOwnerColumn = "is_owner"

// UserClub - Database format for many to many relation with User and Club
// Keeping a record of users to clubs and whether they're an associated owner of the club or not
type UserClub struct {
	UserId  int
	ClubId  int
	IsOwner bool
}

// NewUserClub - Create new default UserClub
func NewUserClub() *UserClub {
	return &UserClub{}
}
