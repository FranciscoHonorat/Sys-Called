package ticket

import (
	"encoding/json"
	"time"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/google/uuid"
)

type Ticket struct {
	id          *valueobjects.ID
	title       *valueobjects.Title
	description *valueobjects.Description
	status      valueobjects.Status
	assigneeID  *valueobjects.AssigneeID
	priority    *valueobjects.Priority
	createdAt   time.Time
	responses   []*response.Response

	uncommittedEvents []event.Event
}

func NewTicket(id *valueobjects.ID, title *valueobjects.Title, description *valueobjects.Description, status valueobjects.Status, assigneeID *valueobjects.AssigneeID, priority *valueobjects.Priority) (*Ticket, error) {
	if id == nil {
		return nil, domainErr.ErrInvalidID
	}
	if title == nil {
		return nil, domainErr.ErrInvalidTitle
	}
	if description == nil {
		return nil, domainErr.ErrInvalidDescription
	}
	if !status.IsValid() {
		return nil, domainErr.ErrInvalidStatus
	}

	t := &Ticket{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		assigneeID:  assigneeID,
		priority:    priority,
		createdAt:   time.Now(),
	}
	t.raise(event.NewTicketOpened(id, title, description, status, assigneeID, priority))

	return t, nil
}

func LoadFromHistory(events []event.Event) (*Ticket, error) {
	if len(events) == 0 {
		return nil, domainErr.ErrEmptyEventHistory
	}
	if _, ok := events[0].(event.TicketOpened); !ok {
		return nil, domainErr.ErrInvalidEventHistory
	}

	t := &Ticket{}
	for _, e := range events {
		if err := t.apply(e); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (t *Ticket) GetID() *valueobjects.ID {
	return t.id
}

func (t *Ticket) GetTitle() *valueobjects.Title {
	return t.title
}

func (t *Ticket) GetDescription() *valueobjects.Description {
	return t.description
}

func (t *Ticket) GetStatus() valueobjects.Status {
	return t.status
}

func (t *Ticket) GetAssigneeID() *valueobjects.AssigneeID {
	return t.assigneeID
}

func (t *Ticket) GetPriority() *valueobjects.Priority {
	return t.priority
}

func (t *Ticket) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Ticket) GetResponses() []*response.Response {
	return t.responses
}

func (t *Ticket) GetUncommittedEvents() []event.Event {
	return t.uncommittedEvents
}

func (t *Ticket) ClearUncommittedEvents() {
	t.uncommittedEvents = nil
}

func (t *Ticket) raise(e event.Event) {
	t.uncommittedEvents = append(t.uncommittedEvents, e)
}

func (t *Ticket) apply(e event.Event) error {
	switch ev := e.(type) {
	case event.TicketOpened:
		title, err := valueobjects.NewTitle(ev.Title)
		if err != nil {
			return err
		}
		description, err := valueobjects.NewDescription(ev.Description)
		if err != nil {
			return err
		}
		status := valueobjects.NewStatus(ev.Status)
		if !status.IsValid() {
			return domainErr.ErrInvalidStatus
		}

		t.id = valueobjects.NewID(ev.AggregateID())
		t.title = title
		t.description = description
		t.status = status
		t.createdAt = ev.OccurredAt()

		if ev.AssigneeID != "" {
			assigneeID, err := valueobjects.NewAssigneeID(ev.AssigneeID)
			if err != nil {
				return err
			}
			t.assigneeID = &assigneeID
		}
		if ev.Priority != "" {
			priority, err := valueobjects.NewPriority(ev.Priority)
			if err != nil {
				return err
			}
			t.priority = &priority
		}

	case event.TicketAssigned:
		assigneeID, err := valueobjects.NewAssigneeID(ev.AssigneeID)
		if err != nil {
			return err
		}
		t.assigneeID = &assigneeID

	case event.TicketPriorityChanged:
		priority, err := valueobjects.NewPriority(ev.Priority)
		if err != nil {
			return err
		}
		t.priority = &priority

	case event.TicketMovedToInProgress:
		t.status = valueobjects.TicketStatusInProgress

	case event.TicketClosed:
		t.status = valueobjects.TicketStatusClosed

	case event.TicketResponseAdded:
		responseID, err := uuid.Parse(ev.ResponseID)
		if err != nil {
			return domainErr.ErrInvalidUUID
		}
		authorID, err := valueobjects.NewAuthorID(ev.AuthorID)
		if err != nil {
			return err
		}
		content, err := valueobjects.NewContent(ev.Content)
		if err != nil {
			return err
		}
		r, err := response.NewResponse(valueobjects.NewID(responseID), t.id, &authorID, content)
		if err != nil {
			return err
		}
		t.responses = append(t.responses, r)
	}

	return nil
}

func (t *Ticket) AddResponse(r *response.Response) error {
	if t.status == valueobjects.TicketStatusClosed {
		return domainErr.ErrTicketAlreadyClosed
	}
	if r == nil {
		return domainErr.ErrInvalidResponse
	}
	if !r.GetTicketID().Equals(t.id) {
		return domainErr.ErrResponseTicketMismatch
	}
	t.responses = append(t.responses, r)
	t.raise(event.NewTicketResponseAdded(t.id, r))
	return nil
}

func (t *Ticket) AssignTo(assigneeID *valueobjects.AssigneeID) error {
	if t.status == valueobjects.TicketStatusClosed {
		return domainErr.ErrTicketAlreadyClosed
	}
	if assigneeID == nil {
		return domainErr.ErrInvalidAssignee
	}
	t.assigneeID = assigneeID
	t.raise(event.NewTicketAssigned(t.id, assigneeID))
	return nil
}

func (t *Ticket) ChangePriority(priority *valueobjects.Priority) error {
	if t.status == valueobjects.TicketStatusClosed {
		return domainErr.ErrTicketAlreadyClosed
	}
	if priority == nil {
		return domainErr.ErrInvalidPriority
	}
	t.priority = priority
	t.raise(event.NewTicketPriorityChanged(t.id, priority))
	return nil
}

func (t *Ticket) MoveToInProgress() error {
	if t.status != valueobjects.TicketStatusOpen {
		return domainErr.ErrInvalidStatusTransition
	}
	t.status = valueobjects.TicketStatusInProgress
	t.raise(event.NewTicketMovedToInProgress(t.id))
	return nil
}

func (t *Ticket) Close() error {
	if t.status == valueobjects.TicketStatusClosed {
		return domainErr.ErrTicketAlreadyClosed
	}
	t.status = valueobjects.TicketStatusClosed
	t.raise(event.NewTicketClosed(t.id))
	return nil
}

func (t *Ticket) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID          string               `json:"id"`
		Title       string               `json:"title"`
		Description string               `json:"description"`
		Status      string               `json:"status"`
		AssigneeID  string               `json:"assignee_id,omitempty"`
		Priority    string               `json:"priority,omitempty"`
		CreatedAt   string               `json:"created_at"`
		Responses   []*response.Response `json:"responses,omitempty"`
	}{
		ID:          t.id.String(),
		Title:       t.title.GetTitle(),
		Description: t.description.GetDescription(),
		Status:      t.status.String(),
		CreatedAt:   t.createdAt.Format(time.RFC3339),
		Responses:   t.responses,
	}

	if t.assigneeID != nil {
		aux.AssigneeID = t.assigneeID.GetAssigneeID()
	}
	if t.priority != nil {
		aux.Priority = t.priority.GetPriority()
	}

	return json.Marshal(aux)
}

func (t *Ticket) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID          string               `json:"id"`
		Title       string               `json:"title"`
		Description string               `json:"description"`
		Status      string               `json:"status"`
		AssigneeID  string               `json:"assignee_id"`
		Priority    string               `json:"priority"`
		CreatedAt   string               `json:"created_at"`
		Responses   []*response.Response `json:"responses"`
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
	t.id = valueobjects.NewID(parsedID)

	title, err := valueobjects.NewTitle(aux.Title)
	if err != nil {
		return err
	}
	t.title = title

	description, err := valueobjects.NewDescription(aux.Description)
	if err != nil {
		return err
	}
	t.description = description

	status := valueobjects.NewStatus(aux.Status)
	if !status.IsValid() {
		return domainErr.ErrInvalidStatus
	}
	t.status = status

	t.assigneeID = nil
	if aux.AssigneeID != "" {
		assigneeID, err := valueobjects.NewAssigneeID(aux.AssigneeID)
		if err != nil {
			return err
		}
		t.assigneeID = &assigneeID
	}

	t.priority = nil
	if aux.Priority != "" {
		priority, err := valueobjects.NewPriority(aux.Priority)
		if err != nil {
			return err
		}
		t.priority = &priority
	}

	createdAt, err := time.Parse(time.RFC3339, aux.CreatedAt)
	if err != nil {
		return err
	}
	t.createdAt = createdAt

	for _, r := range aux.Responses {
		if r == nil || !r.GetTicketID().Equals(t.id) {
			return domainErr.ErrResponseTicketMismatch
		}
	}
	t.responses = aux.Responses

	return nil
}
