package valueobjects

import (
	"github.com/google/uuid"
)

type ID struct {
	ID uuid.UUID
}

func NewID(id uuid.UUID) *ID {
	if id == uuid.Nil {
		id = uuid.New()
	}
	return &ID{ID: id}
}

func (id *ID) GetID() uuid.UUID {
	return id.ID
}

func (id *ID) String() string {
	return id.ID.String()
}

func (id *ID) Equals(other *ID) bool {
	if other == nil {
		return false
	}
	return id.ID == other.ID
}
