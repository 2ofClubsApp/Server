package model

// User base struct
type User struct {
	Base
	*Credentials `json:"-"`
	Manages      []Club  `gorm:"many2many:user_club;" json:"-"`
	Chooses      []Tag   `gorm:"many2many:user_tag;foreignKey:id;References:Name" json:"-"`
	Attends      []Event `gorm:"many2many:user_event;" json:"-"`
	IsAdmin      bool    `json:"-"`
	IsApproved   bool    `json:"-"`
	Swiped       []Club  `gorm:"many2many:user_swipe_club;" json:"-"`
}

// UserDisplay - Displaying public user data
type UserDisplay struct {
	Email   string            `json:"email"`
	Manages []*ManagesDisplay `json:"manages"`
	Tags    []Tag             `json:"tags"`
	Attends []Event           `json:"attends"`
}

// UserBaseInfo - Displaying basic user data
type UserBaseInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// ManagesDisplay is used as a display wrapper for the ClubDisplay
// For a users managed club, whether they're an owner or not is also displayed
type ManagesDisplay struct {
	Club
	IsOwner bool `json:"isOwner"`
}

// DisplayAllInfo public user data
func (u *User) DisplayAllInfo() *UserDisplay {
	return &UserDisplay{
		Email: u.Email,
	}
}

// DisplayBaseUserInfo displays base user data
func (u *User) DisplayBaseUserInfo() UserBaseInfo {
	return UserBaseInfo{
		ID:       u.ID,
		Username: u.Username,
	}
}

// NewUser - Create new default User
func NewUser() *User {
	return &User{
		Credentials: NewCredentials(),
		Manages:     []Club{},
		Chooses:     []Tag{},
		Attends:     []Event{},
		Swiped:      []Club{},
	}
}

// User variables for db columns/route variables
const (
	AllUserInfo         = "all"
	AllUserClubsManage  = "clubs"
	AllUserEventsAttend = "events"
	ChoosesColumn       = "Chooses"
	AttendsColumn       = "Attends"
	SwipedColumn        = "Swiped"
	UserClubTable       = "user_club"
	UserTable           = "user"
	ManagesColumn       = "Manages"
	IsApprovedColumn    = "is_approved"
	IDColumn            = "id"
)
