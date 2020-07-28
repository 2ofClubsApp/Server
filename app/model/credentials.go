package model

type Credentials struct {
	Username string `gorm:"UNIQUE" validate:"alpha,min=2,max=15,required"`
	Email    string `gorm:"UNIQUE" validate:"required,email"`
	Password string `validate:"required,min=3,max=45"`
	// Max 45 due to 50 length limitation of bcrypt


}

func NewCredentials() *Credentials {
	return &Credentials{}
}

const (

	UsernameColumn  = "username"
	EmailColumn     = "email"
	PasswordColumn  = "password"
)
