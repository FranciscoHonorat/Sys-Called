package response

import (
	"time"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type Response struct {
	id        *valueobjects.ID
	ticketID  *valueobjects.ID
	authorID  *valueobjects.AuthorID
	content   *valueobjects.Content
	createdAt time.Time
}

func NewResponse(id *valueobjects.ID, ticketID *valueobjects.ID, authorID *valueobjects.AuthorID, content *valueobjects.Content) (*Response, error) {
	if id == nil {
		return nil, domainErr.ErrInvalidID
	}
	if ticketID == nil {
		return nil, domainErr.ErrInvalidTicketID
	}
	if authorID == nil {
		return nil, domainErr.ErrInvalidAuthorID
	}
	if content == nil {
		return nil, domainErr.ErrInvalidContent
	}

	return &Response{
		id:        id,
		ticketID:  ticketID,
		authorID:  authorID,
		content:   content,
		createdAt: time.Now(),
	}, nil
}

func RestoreResponse(id *valueobjects.ID, ticketID *valueobjects.ID, authorID *valueobjects.AuthorID, content *valueobjects.Content, createdAt time.Time) (*Response, error) {
	r, err := NewResponse(id, ticketID, authorID, content)
	if err != nil {
		return nil, err
	}
	r.createdAt = createdAt
	return r, nil
}

func (r *Response) GetID() *valueobjects.ID {
	return r.id
}

func (r *Response) GetTicketID() *valueobjects.ID {
	return r.ticketID
}

func (r *Response) GetAuthorID() *valueobjects.AuthorID {
	return r.authorID
}

func (r *Response) GetContent() *valueobjects.Content {
	return r.content
}

func (r *Response) GetCreatedAt() time.Time {
	return r.createdAt
}
