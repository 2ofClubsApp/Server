package model

const (
	UserNotFound  = "User Not Found"
	UserFound     = "User Found"
	UsernameExists = "Username already exists"
	EmailExists    = "Email already exists"
)

type T struct{}

type Status struct {
	Message string
	Data    interface{}
}

type CredentialStatus struct {
	Username string
	Email    string
}

func NewStatus() *Status {
	return &Status{Message: "", Data: T{}}
}
