package valueobjects

type Description struct {
	Description string
}

func NewDescription(description string) *Description {
	return &Description{Description: description}
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
	return []byte(`"` + d.Description + `"`), nil
}

func (d *Description) UnmarshalJSON(data []byte) error {
	d.Description = string(data[1 : len(data)-1]) // Remove quotes
	return nil
}
