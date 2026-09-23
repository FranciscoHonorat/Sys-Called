package out

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type EmployeeRepository interface {
	List(ctx context.Context) ([]employee.Employee, error)
	FindByUsername(ctx context.Context, username string) (employee.Employee, bool, error)
	FindByID(ctx context.Context, id string) (employee.Employee, bool, error)
	Register(ctx context.Context, e employee.Employee) error
	Save(ctx context.Context, e employee.Employee) error
}
