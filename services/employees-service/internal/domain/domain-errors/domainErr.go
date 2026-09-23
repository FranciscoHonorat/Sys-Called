package domainerrors

import "errors"

var (
	ErrInvalidEmployeeID   = errors.New("invalid employee ID: ID cannot be empty")
	ErrInvalidEmployeeName = errors.New("invalid employee name: name cannot be empty")
	ErrInvalidRole         = errors.New("invalid role: role must be user, support or admin")
	ErrInvalidUsername     = errors.New("invalid username: username cannot be empty")
	ErrInvalidPassword     = errors.New("invalid password: password cannot be empty")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrUsernameTaken       = errors.New("username already taken")
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrWeakPassword        = errors.New("password must have at least 8 characters")
	ErrPasswordTooLong     = errors.New("password must have at most 72 bytes")
	ErrPendingApproval     = errors.New("account waiting for the administrator approval")
	ErrForbidden           = errors.New("forbidden: you are not allowed to perform this action")
	ErrUnauthenticated     = errors.New("unauthenticated")
)
