package actor

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Role string

const (
	RoleUser    Role = "user"
	RoleSupport Role = "support"
	RoleAdmin   Role = "admin"
)

func (r Role) isKnown() bool {
	switch r {
	case RoleUser, RoleSupport, RoleAdmin:
		return true
	}
	return false
}

type Actor struct {
	id   string
	name string
	role Role
}

func New(id, name, role string) (Actor, error) {
	if id == "" || !Role(role).isKnown() {
		return Actor{}, domainErr.ErrInvalidActor
	}
	return Actor{id: id, name: name, role: Role(role)}, nil
}

func (a Actor) ID() string {
	return a.id
}

func (a Actor) Name() string {
	return a.name
}

func (a Actor) Role() Role {
	return a.role
}

func (a Actor) IsAdmin() bool {
	return a.role == RoleAdmin
}

func (a Actor) IsSupport() bool {
	return a.role == RoleSupport
}
