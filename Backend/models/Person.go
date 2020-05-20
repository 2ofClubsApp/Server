package models

/*
Notes:
	- The password needs to be hashed later
 */

/*
	"Abstract" models
*/
type Person struct {
	ID       uint    `gorm:"AUTO_INCREMENT"`
	Username string `gorm:"primary_key"`
	Email    string
	Password string
}

func (p *Person) SetUsername(username string) {
	p.Username = username
}

func (p *Person) SetEmail(email string) {
	p.Email = email
}

func (p *Person) SetPassword(password string) {
	p.Password = password
}

func (p Person) GetUsername() string {
	return p.Username
}

func (p Person) GetEmail() string {
	return p.Email
}

func (p Person) GetPassword() string {
	return p.Password
}
