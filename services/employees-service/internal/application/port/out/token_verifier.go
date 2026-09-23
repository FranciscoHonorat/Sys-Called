package out

import (
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type Caller struct {
	ID   string
	Role employee.Role
}

type TokenVerifier interface {
	Verify(token string) (Caller, error)
}
