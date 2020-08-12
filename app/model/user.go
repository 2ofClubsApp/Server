package model

type User struct {
	Base
	*Credentials
	Manages      []Club  `gorm:"many2many:user_club;" json:"-"`
	Chooses      []Tag   `gorm:"many2many:user_tag;foreignKey:id;References:Name" json:"-"`
	Attends      []Event `gorm:"many2many:user_event;" json:"-"`
	IsAdmin      bool    `json:"-"`
	IsApproved   bool    `json:"-"`
}

type UserDisplay struct {
	Email   string
	Manages []*ManagesDisplay
	Tags    []string
	Attends []Event
}

type ManagesDisplay struct {
	*ClubDisplay
	IsOwner bool
}

func (u *User) Display() *UserDisplay {
	return &UserDisplay{Email: u.Email}
}

func NewUser() *User {
	return &User{Credentials: NewCredentials(), Manages: []Club{}, Chooses: []Tag{}, Attends: []Event{}}
}

const (
	AllUserInfo         = "all"
	AllUserClubsManage  = "clubs"
	AllUserEventsAttend = "events"
	ChoosesColumn       = "Chooses"
	AttendsColumn       = "Attends"
	UserTagTable        = "user_tag"
	UserClubTable       = "user_club"
	IsHelpingColumn     = "is_helping"
	UserTable           = "user"
	ManagesColumn       = "Manages"
	IsAdminColumn       = "is_admin"
	IsApprovedColumn    = "is_approved"
	IDColumn            = "id"
	CreatedAtColumn     = "created_at"
	DeletedAtColumn     = "deleted_at"
)
