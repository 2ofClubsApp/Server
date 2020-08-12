package model

type Credentials struct {
	Username string `gorm:"UNIQUE" validate:"alpha,min=2,max=15,required"`
	Email    string `gorm:"UNIQUE" validate:"required,email"`
	Password string `validate:"required,min=3,max=45"`
	// Max 45 due to 50 length limitation of bcrypt
}

type PasswordChange struct {
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

func NewPasswordChange() *PasswordChange {
	return &PasswordChange{}
}

func NewCredentials() *Credentials {
	return &Credentials{}
}

const (
	EmailVar       = "email"
	UsernameVar    = "username"
	UsernameColumn = "username"
	EmailColumn    = "email"
	PasswordColumn = "password"
)
