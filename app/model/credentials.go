package model

// Credentials struct for user login/signup
type Credentials struct {
	Username string `gorm:"UNIQUE" validate:"alpha,min=2,max=15,required" json:"username"`
	Email    string `gorm:"UNIQUE" validate:"required,email" json:"email"`

	// Max 45 due to 50 length limitation of bcrypt
	Password string `validate:"required,min=3,max=45" json:"password"`
}

// Resetting a user password given the old and new passwords
type PasswordChange struct {
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

// Create new default PasswordChange
func NewPasswordChange() *PasswordChange {
	return &PasswordChange{}
}

// Create new default Credentials
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
