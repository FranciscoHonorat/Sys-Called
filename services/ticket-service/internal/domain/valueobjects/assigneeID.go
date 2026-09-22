package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AssigneeID struct {
	id string
}

func NewAssigneeID(id string) (AssigneeID, error) {
	assignee := AssigneeID{id: id}
	if !assignee.IsValid() {
		return AssigneeID{}, domainErr.ErrInvalidAssignee
	}
	return assignee, nil
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
	return json.Marshal(a.id)
}

func (a *AssigneeID) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	assignee := AssigneeID{id: id}
	if !assignee.IsValid() {
		return domainErr.ErrInvalidAssignee
	}
	*a = assignee
	return nil
}
