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
)
