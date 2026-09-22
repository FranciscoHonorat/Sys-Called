package valueobjects

import (
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Content struct {
	Content string
}

func NewContent(content string) (*Content, error) {
	value, err := newNonEmptyString(content, domainErr.ErrInvalidContent)
	if err != nil {
		return nil, err
	}
	return &Content{Content: value}, nil
}

func (c *Content) GetContent() string {
	return c.Content
}

func (c *Content) IsValid() bool {
	return c.Content != ""
}

func (c *Content) Equals(other *Content) bool {
	if other == nil {
		return false
	}
	return c.Content == other.Content
}

func (c *Content) MarshalJSON() ([]byte, error) {
	return marshalNonEmptyString(c.Content)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	value, err := unmarshalNonEmptyString(data, domainErr.ErrInvalidContent)
	if err != nil {
		return err
	}
	c.Content = value
	return nil
}
