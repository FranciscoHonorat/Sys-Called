package session

import "time"

type RefreshToken struct {
	hash       string
	employeeID string
	expiresAt  time.Time
	revoked    bool
}

func NewRefreshToken(hash, employeeID string, expiresAt time.Time) RefreshToken {
	return RestoreRefreshToken(hash, employeeID, expiresAt, false)
}

func RestoreRefreshToken(hash, employeeID string, expiresAt time.Time, revoked bool) RefreshToken {
	return RefreshToken{hash: hash, employeeID: employeeID, expiresAt: expiresAt, revoked: revoked}
}

func (t RefreshToken) Hash() string {
	return t.hash
}

func (t RefreshToken) EmployeeID() string {
	return t.employeeID
}

func (t RefreshToken) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t RefreshToken) IsExpired(now time.Time) bool {
	return !now.Before(t.expiresAt)
}

func (t RefreshToken) IsRevoked() bool {
	return t.revoked
}

func (t *RefreshToken) Revoke() {
	t.revoked = true
}
