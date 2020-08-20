package model

// User base struct
type User struct {
	Base
	*Credentials
	Manages    []Club  `gorm:"many2many:user_club;" json:"-"`
	Chooses    []Tag   `gorm:"many2many:user_tag;foreignKey:id;References:Name" json:"-"`
	Attends    []Event `gorm:"many2many:user_event;" json:"-"`
	IsAdmin    bool    `json:"-"`
	IsApproved bool    `json:"-"`
}

// UserDisplay - Displaying public user data
type UserDisplay struct {
	Email   string            `json:"email"`
	Manages []*ManagesDisplay `json:"manages"`
	Tags    []Tag             `json:"tags"`
	Attends []Event           `json:"attends"`
}

// ManagesDisplay is used as a display wrapper for the ClubDisplay
// For a users managed club, whether they're an owner or not is also displayed
type ManagesDisplay struct {
	Club
	IsOwner bool `json:"isOwner"`
}

// Display - Returns public club data
func (u *User) Display() *UserDisplay {
	return &UserDisplay{Email: u.Email}
}

// NewUser - Create new default User
func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}, Chooses: []Tag{}, Attends: []Event{}}
}

// User variables for db columns/route variables
const (
	AllUserInfo         = "all"
	AllUserClubsManage  = "clubs"
	AllUserEventsAttend = "events"
	ChoosesColumn       = "Chooses"
	AttendsColumn       = "Attends"
	UserClubTable       = "user_club"
	UserTable           = "user"
	ManagesColumn       = "Manages"
	IsApprovedColumn    = "is_approved"
	IDColumn            = "id"
)
