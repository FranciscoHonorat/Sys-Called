package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Status string

const (
	TicketStatusOpen       Status = "Open"
	TicketStatusInProgress Status = "In Progress"
	TicketStatusClosed     Status = "Closed"
)

func NewStatus(status string) Status {
	return Status(status)
}

func (s Status) IsValid() bool {
	switch s {
	case TicketStatusOpen, TicketStatusInProgress, TicketStatusClosed:
		return true
	default:
		return false
	}
}

func (s Status) Equals(other Status) bool {
	return s == other
}

func (s Status) String() string {
	return string(s)
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	ticketStatus := Status(status)
	if !ticketStatus.IsValid() {
		return domainErr.ErrInvalidStatus
	}
	*s = ticketStatus
	return nil
}
