package model

const (
	Failure = "Student Not Found"
	Success = "Student Found"
)

type T struct {}

type Status struct {
	Message string
	Data   interface{}
}

func NewStatus() *Status {
	return &Status{Message: Success, Data: T{}}
}
