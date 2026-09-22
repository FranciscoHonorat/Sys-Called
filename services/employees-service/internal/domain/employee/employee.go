package employee

type Employee struct {
	id     string
	name   string
	events []Event
}

func NewEmployee(id, name string) Employee {
	return Employee{id: id, name: name}
}

func Register(id, name string) Employee {
	e := NewEmployee(id, name)
	e.events = append(e.events, Registered{ID: id, Name: name})
	return e
}

func (e Employee) GetID() string {
	return e.id
}

func (e Employee) GetName() string {
	return e.name
}

func (e Employee) Events() []Event {
	return e.events
}
