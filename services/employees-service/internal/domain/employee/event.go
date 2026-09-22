package employee

type Event interface {
	EventType() string
}

type Registered struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Registered) EventType() string {
	return "EmployeeRegistered"
}
