package httpapi

import (
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/repository"
)

type Handler struct {
	listEmployees *application.ListEmployeesUseCase
}

func NewHandler(repo repository.EmployeeRepository) *Handler {
	return &Handler{listEmployees: application.NewListEmployeesUseCase(repo)}
}
