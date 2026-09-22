package domainerrors

import "errors"

var (
	ErrInvalidTitle       = errors.New("invalid title: title cannot be empty")
	ErrInvalidDescription = errors.New("invalid description: description cannot be empty")
	ErrInvalidStatus      = errors.New("invalid status: status must be one of the allowed values")
	ErrInvalidPriority    = errors.New("invalid priority: priority must be one of the allowed values")
	ErrInvalidAssignee    = errors.New("invalid assignee: assignee cannot be empty")
	ErrInvalidID          = errors.New("invalid ID: ID cannot be empty")
	ErrInvalidUUID        = errors.New("invalid UUID: UUID cannot be empty")

	ErrTicketAlreadyClosed     = errors.New("invalid operation: ticket is already closed")
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	ErrInvalidContent  = errors.New("invalid content: content cannot be empty")
	ErrInvalidAuthorID = errors.New("invalid author: author cannot be empty")
	ErrInvalidTicketID = errors.New("invalid ticket ID: ticket ID cannot be empty")

	ErrInvalidResponse        = errors.New("invalid response: response cannot be nil")
	ErrResponseTicketMismatch = errors.New("invalid response: response does not belong to this ticket")

	ErrEmptyEventHistory   = errors.New("invalid history: event history cannot be empty")
	ErrInvalidEventHistory = errors.New("invalid history: first event must be TicketOpened")
)
