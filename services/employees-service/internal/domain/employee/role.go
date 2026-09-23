package employee

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type Role string

const (
	RoleUser    Role = "user"
	RoleSupport Role = "support"
	RoleAdmin   Role = "admin"
)

func NewRole(raw string) (Role, error) {
	role := Role(raw)
	switch role {
	case RoleUser, RoleSupport, RoleAdmin:
		return role, nil
	default:
		return "", domainErr.ErrInvalidRole
	}
}

func (r Role) String() string {
	return string(r)
}
