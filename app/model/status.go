package model

const (
	Failure = "Student Not Found"
	Success = "Student Found"
	UsernameFound = "Username found"
	EmailFound = "Email Found"
)

type T struct {}

type Status struct {
	Message string
	Data   interface{}
}

func NewStatus() *Status {
	return &Status{Message: "", Data: T{}}
}
