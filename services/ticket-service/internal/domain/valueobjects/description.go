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
