package response

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

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

func (r *Response) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID        string `json:"id"`
		TicketID  string `json:"ticket_id"`
		AuthorID  string `json:"author_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        r.id.String(),
		TicketID:  r.ticketID.String(),
		AuthorID:  r.authorID.GetAuthorID(),
		Content:   r.content.GetContent(),
		CreatedAt: r.createdAt.Format(time.RFC3339),
	}

	return json.Marshal(aux)
}

func (r *Response) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID        string `json:"id"`
		TicketID  string `json:"ticket_id"`
		AuthorID  string `json:"author_id"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedID := uuid.Nil
	if aux.ID != "" {
		var err error
		parsedID, err = uuid.Parse(aux.ID)
		if err != nil {
			return domainErr.ErrInvalidUUID
		}
	}
	r.id = valueobjects.NewID(parsedID)

	if aux.TicketID == "" {
		return domainErr.ErrInvalidTicketID
	}
	parsedTicketID, err := uuid.Parse(aux.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}
	r.ticketID = valueobjects.NewID(parsedTicketID)

	authorID, err := valueobjects.NewAuthorID(aux.AuthorID)
	if err != nil {
		return err
	}
	r.authorID = &authorID

	content, err := valueobjects.NewContent(aux.Content)
	if err != nil {
		return err
	}
	r.content = content

	createdAt, err := time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	r.createdAt = createdAt

	return nil
}
