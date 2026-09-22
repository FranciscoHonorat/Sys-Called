package valueobjects

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AuthorID struct {
	id string
}

func NewAuthorID(id string) (AuthorID, error) {
	value, err := newNonEmptyString(id, domainErr.ErrInvalidAuthorID)
	if err != nil {
		return AuthorID{}, err
	}
	return AuthorID{id: value}, nil
}

func (a AuthorID) GetAuthorID() string {
	return a.id
}
