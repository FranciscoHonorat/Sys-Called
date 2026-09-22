package valueobjects

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Title struct {
	Title string
}

func NewTitle(title string) (*Title, error) {
	value, err := newNonEmptyString(title, domainErr.ErrInvalidTitle)
	if err != nil {
		return nil, err
	}
	return &Title{Title: value}, nil
}

func (t *Title) GetTitle() string {
	return t.Title
}

func (t *Title) IsValid() bool {
	return t.Title != ""
}

func (t *Title) Equals(other *Title) bool {
	if other == nil {
		return false
	}
	return t.Title == other.Title
}

func (t *Title) MarshalJSON() ([]byte, error) {
	return marshalNonEmptyString(t.Title)
}

func (t *Title) UnmarshalJSON(data []byte) error {
	value, err := unmarshalNonEmptyString(data, domainErr.ErrInvalidTitle)
	if err != nil {
		return err
	}
	t.Title = value
	return nil
}
