package employee

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusPending Status = "pending"
)

type Snapshot struct {
	ID                     string
	Name                   string
	Username               string
	Role                   Role
	PasswordHash           string
	Status                 Status
	MustChangePassword     bool
	PasswordResetRequested bool
}

type Employee struct {
	id                     string
	name                   string
	username               string
	role                   Role
	passwordHash           string
	status                 Status
	mustChangePassword     bool
	passwordResetRequested bool
	events                 []Event
}

func Rehydrate(s Snapshot) Employee {
	return Employee{
		id:                     s.ID,
		name:                   s.Name,
		username:               s.Username,
		role:                   s.Role,
		passwordHash:           s.PasswordHash,
		status:                 s.Status,
		mustChangePassword:     s.MustChangePassword,
		passwordResetRequested: s.PasswordResetRequested,
	}
}

func NewEmployee(id, name, username string, role Role, passwordHash string) Employee {
	return Rehydrate(Snapshot{ID: id, Name: name, Username: username, Role: role, PasswordHash: passwordHash, Status: StatusActive})
}

func Register(id, name, username string, role Role, passwordHash string) (Employee, error) {
	if err := validateIdentity(id, name, username); err != nil {
		return Employee{}, err
	}

	e := NewEmployee(id, name, username, role, passwordHash)
	e.events = append(e.events, Registered{ID: id, Name: name, Role: role.String()})
	return e, nil
}

func SignUp(id, name, username, passwordHash string) (Employee, error) {
	if err := validateIdentity(id, name, username); err != nil {
		return Employee{}, err
	}

	e := NewEmployee(id, name, username, RoleUser, passwordHash)
	e.status = StatusPending
	e.events = append(e.events, SignedUp{ID: id, Name: name})
	return e, nil
}

func validateIdentity(id, name, username string) error {
	switch {
	case id == "":
		return domainErr.ErrInvalidEmployeeID
	case name == "":
		return domainErr.ErrInvalidEmployeeName
	case username == "":
		return domainErr.ErrInvalidUsername
	}
	return nil
}

func (e *Employee) Approve() {
	e.status = StatusActive
}

func (e *Employee) RequestPasswordReset() {
	if e.passwordResetRequested {
		return
	}
	e.passwordResetRequested = true
	e.events = append(e.events, PasswordResetRequested{ID: e.id, Name: e.name})
}

func (e *Employee) SetTemporaryPassword(passwordHash string) {
	e.passwordHash = passwordHash
	e.mustChangePassword = true
	e.passwordResetRequested = false
}

func (e *Employee) ChangePassword(passwordHash string) {
	e.passwordHash = passwordHash
	e.mustChangePassword = false
}

func (e Employee) CanLogIn() bool {
	return e.status == StatusActive
}

func (e Employee) GetID() string {
	return e.id
}

func (e Employee) GetName() string {
	return e.name
}

func (e Employee) GetUsername() string {
	return e.username
}

func (e Employee) GetRole() Role {
	return e.role
}

func (e Employee) GetPasswordHash() string {
	return e.passwordHash
}

func (e Employee) GetStatus() Status {
	return e.status
}

func (e Employee) MustChangePassword() bool {
	return e.mustChangePassword
}

func (e Employee) HasRequestedPasswordReset() bool {
	return e.passwordResetRequested
}

func (e Employee) Events() []Event {
	return e.events
}
