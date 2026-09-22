package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Description struct {
	Description string
}

func NewDescription(description string) (*Description, error) {
	d := &Description{Description: description}
	if !d.IsValid() {
		return nil, domainErr.ErrInvalidDescription
	}
	return d, nil
}

func (d *Description) GetDescription() string {
	return d.Description
}

func (d *Description) IsValid() bool {
	return d.Description != ""
}

func (d *Description) Equals(other *Description) bool {
	if other == nil {
		return false
	}
	return d.Description == other.Description
}

func (d *Description) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Description)
}

func (d *Description) UnmarshalJSON(data []byte) error {
	var description string
	if err := json.Unmarshal(data, &description); err != nil {
		return err
	}
	d.Description = description
	return nil
}
