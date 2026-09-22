package domainerrors

import "errors"

var (
	ErrInvalidEmployeeID   = errors.New("invalid employee ID: ID cannot be empty")
	ErrInvalidEmployeeName = errors.New("invalid employee name: name cannot be empty")
)
