package valueobjects

import (
	"encoding/json"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type Priority string

const (
	TicketPriorityLow    Priority = "Low"
	TicketPriorityMedium Priority = "Medium"
	TicketPriorityHigh   Priority = "High"
)

func NewPriority(priority string) (Priority, error) {
	ticketPriority := Priority(priority)
	if !ticketPriority.IsValid() {
		return "", domainErr.ErrInvalidPriority
	}
	return ticketPriority, nil
}

func (p Priority) IsValid() bool {
	switch p {
	case TicketPriorityLow, TicketPriorityMedium, TicketPriorityHigh:
		return true
	default:
		return false
	}
}

func (p Priority) Equals(other Priority) bool {
	return p == other
}

func (p Priority) GetPriority() string {
	return string(p)
}

func (p Priority) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(p))
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	var priority string
	if err := json.Unmarshal(data, &priority); err != nil {
		return err
	}
	ticketPriority := Priority(priority)
	if !ticketPriority.IsValid() {
		return domainErr.ErrInvalidPriority
	}
	*p = ticketPriority
	return nil
}
