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
