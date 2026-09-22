package repository

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type EmployeeRepository interface {
	List(ctx context.Context) ([]employee.Employee, error)
	Register(ctx context.Context, e employee.Employee) error
}
