package out

import (
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type TokenIssuer interface {
	Issue(e employee.Employee) (token string, expiresIn time.Duration, err error)
}
