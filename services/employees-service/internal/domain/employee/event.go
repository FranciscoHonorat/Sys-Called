package employee

type Event interface {
	EventType() string
}

type Registered struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (Registered) EventType() string {
	return "EmployeeRegistered"
}

type SignedUp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (SignedUp) EventType() string {
	return "EmployeeSignedUp"
}

type PasswordResetRequested struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (PasswordResetRequested) EventType() string {
	return "PasswordResetRequested"
}
