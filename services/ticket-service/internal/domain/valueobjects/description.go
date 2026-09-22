package valueobjects

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Description struct {
	Description string
}

func NewDescription(description string) (*Description, error) {
	value, err := newNonEmptyString(description, domainErr.ErrInvalidDescription)
	if err != nil {
		return nil, err
	}
	return &Description{Description: value}, nil
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
	return marshalNonEmptyString(d.Description)
}

func (d *Description) UnmarshalJSON(data []byte) error {
	value, err := unmarshalNonEmptyString(data, domainErr.ErrInvalidDescription)
	if err != nil {
		return err
	}
	d.Description = value
	return nil
}
