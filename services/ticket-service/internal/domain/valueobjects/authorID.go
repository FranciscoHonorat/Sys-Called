package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AuthorID struct {
	id string
}

func NewAuthorID(id string) (AuthorID, error) {
	author := AuthorID{id: id}
	if !author.IsValid() {
		return AuthorID{}, domainErr.ErrInvalidAuthorID
	}
	return author, nil
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
	return json.Marshal(a.id)
}

func (a *AuthorID) UnmarshalJSON(data []byte) error {
	var id string
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	author := AuthorID{id: id}
	if !author.IsValid() {
		return domainErr.ErrInvalidAuthorID
	}
	*a = author
	return nil
}
