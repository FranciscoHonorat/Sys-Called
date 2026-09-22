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
