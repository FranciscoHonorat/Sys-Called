package valueobjects

import (
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
	return isOneOf(p, TicketPriorityLow, TicketPriorityMedium, TicketPriorityHigh)
}

func (p Priority) GetPriority() string {
	return string(p)
}
