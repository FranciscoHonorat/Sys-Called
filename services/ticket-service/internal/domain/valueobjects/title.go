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
