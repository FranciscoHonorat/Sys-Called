package employee

type Employee struct {
	id   string
	name string
}

func NewEmployee(id, name string) Employee {
	return Employee{id: id, name: name}
}

func (e Employee) GetID() string {
	return e.id
}

func (e Employee) GetName() string {
	return e.name
}
