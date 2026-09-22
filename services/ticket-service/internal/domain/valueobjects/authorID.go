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

func (a AuthorID) IsValid() bool {
	return a.id != ""
}

func (a AuthorID) Equals(other AuthorID) bool {
	return a.id == other.id
}

func (a AuthorID) MarshalJSON() ([]byte, error) {
	return marshalNonEmptyString(a.id)
}

func (a *AuthorID) UnmarshalJSON(data []byte) error {
	value, err := unmarshalNonEmptyString(data, domainErr.ErrInvalidAuthorID)
	if err != nil {
		return err
	}
	a.id = value
	return nil
}
