package Status

const (
	FAILURE = 0
	SUCCESS = 1
)

type Status struct {
	Status int
	Data   interface{}
}

