package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Title struct {
	Title string
}

func NewTitle(title string) (*Title, error) {
	t := &Title{Title: title}
	if !t.IsValid() {
		return nil, domainErr.ErrInvalidTitle
	}
	return t, nil
}

func (t *Title) GetTitle() string {
	return t.Title
}

func (t *Title) IsValid() bool {
	return len(t.Title) > 0
}

func (t *Title) Equals(other *Title) bool {
	if other == nil {
		return false
	}
	return t.Title == other.Title
}

func (t *Title) MarshalJSON() ([]byte, error) {
	type Alias Title
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

func (t *Title) UnmarshalJSON(data []byte) error {
	type Alias Title
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Title = aux.Title
	return nil
}
