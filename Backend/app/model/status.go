package model

const (
	FAILURE = "Student Not Found"
	SUCCESS = "Student Found"
)

type T struct {}

type Status struct {
	Message string
	Data   interface{}
}

func NewStatus() *Status {
	return &Status{Message: SUCCESS, Data: T{}}
}
