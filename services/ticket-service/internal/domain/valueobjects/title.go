package valueobjects

import "encoding/json"

type Title struct {
	Title string
}

func NewTitle(title string) *Title {
	return &Title{Title: title}
}

func (t *Title) GetTitle() string {
	return t.Title
}

func (t *Title) IsValid() bool {
	return len(t.Title) > 0
}

func (t *Title) Equals(other *Title) bool {
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
