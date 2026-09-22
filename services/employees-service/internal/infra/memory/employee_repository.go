package memory

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type EmployeeRepository struct {
	employees []employee.Employee
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		employees: []employee.Employee{
			employee.NewEmployee("agent-1", "Ana Souza"),
			employee.NewEmployee("agent-2", "Bruno Lima"),
			employee.NewEmployee("agent-3", "Carla Melo"),
		},
	}
}

func (r *EmployeeRepository) List(_ context.Context) ([]employee.Employee, error) {
	return r.employees, nil
}

func (r *EmployeeRepository) Register(_ context.Context, e employee.Employee) error {
	for _, existing := range r.employees {
		if existing.GetID() == e.GetID() {
			return nil
		}
	}
	r.employees = append(r.employees, e)
	return nil
}
