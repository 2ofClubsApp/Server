package model

const (
	FAILURE = "Student Not Found"
	SUCCESS = "Student Found"
)

type T struct {}

type status struct {
	Message string
	Data   interface{}
}

func New() *status {
	return &status{Message: SUCCESS, Data: T{}}
}
