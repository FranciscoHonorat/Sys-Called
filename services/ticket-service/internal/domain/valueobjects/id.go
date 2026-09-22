package valueobjects

import (
	"encoding/json"

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
	return id.ID == other.ID
}

func (id *ID) MarshalJSON() ([]byte, error) {
	auxID := struct {
		ID string `json:"id"`
	}{
		ID: id.ID.String(),
	}
	return json.Marshal(auxID)
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var auxID struct {
		ID uuid.UUID
	}
	if err := json.Unmarshal(data, &auxID); err != nil {
		return err
	}
	if auxID.ID == uuid.Nil {
		id.ID = uuid.Nil
		return nil
	}
	id.ID = auxID.ID
	return nil
}
