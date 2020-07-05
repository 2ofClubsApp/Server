package model

const (
<<<<<<< Updated upstream
	FAILURE = "Student Not Found"
	SUCCESS = "Student Found"
=======
	Failure = "Student Not Found"
	Success = "Student Found"
	UsernameFound = "Username found"
	EmailFound = "Email Found"
>>>>>>> Stashed changes
)

type T struct {}

type Status struct {
	Message string
	Data   interface{}
}

func NewStatus() *Status {
	return &Status{Message: SUCCESS, Data: T{}}
}
