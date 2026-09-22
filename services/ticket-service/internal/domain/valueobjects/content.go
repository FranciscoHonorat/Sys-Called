package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Content struct {
	Content string
}

func NewContent(content string) (*Content, error) {
	c := &Content{Content: content}
	if !c.IsValid() {
		return nil, domainErr.ErrInvalidContent
	}
	return c, nil
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
	return json.Marshal(c.Content)
}

func (c *Content) UnmarshalJSON(data []byte) error {
	var content string
	if err := json.Unmarshal(data, &content); err != nil {
		return err
	}
	candidate := Content{Content: content}
	if !candidate.IsValid() {
		return domainErr.ErrInvalidContent
	}
	*c = candidate
	return nil
}
