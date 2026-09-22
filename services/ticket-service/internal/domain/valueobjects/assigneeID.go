package valueobjects

import (
	"encoding/json"
	"reflect"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AssigneeID struct {
	id string
}

func NewAssigneeID(id string) AssigneeID {
	return AssigneeID{id: id}
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
	return []byte(`"` + a.id + `"`), nil
}

func (a *AssigneeID) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return &json.UnmarshalTypeError{Value: string(data), Type: reflect.TypeOf(a)}
	}
	a.id = string(data[1 : len(data)-1])
	return domainErr.ErrInvalidAssignee
}
