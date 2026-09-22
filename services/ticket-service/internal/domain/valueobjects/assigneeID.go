package valueobjects

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AssigneeID struct {
	id string
}

func NewAssigneeID(id string) (AssigneeID, error) {
	value, err := newNonEmptyString(id, domainErr.ErrInvalidAssignee)
	if err != nil {
		return AssigneeID{}, err
	}
	return AssigneeID{id: value}, nil
}

func (a AssigneeID) GetAssigneeID() string {
	return a.id
}

func (a AssigneeID) IsValid() bool {
	return a.id != ""
}

func (a AssigneeID) Equals(other AssigneeID) bool {
	return a.id == other.id
}

func (a AssigneeID) MarshalJSON() ([]byte, error) {
	return marshalNonEmptyString(a.id)
}

func (a *AssigneeID) UnmarshalJSON(data []byte) error {
	value, err := unmarshalNonEmptyString(data, domainErr.ErrInvalidAssignee)
	if err != nil {
		return err
	}
	a.id = value
	return nil
}
