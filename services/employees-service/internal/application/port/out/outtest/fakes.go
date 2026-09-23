package outtest

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

type EmployeeRepository struct {
	Employees   []employee.Employee
	Err         error
	RegisterErr error
}

func (r *EmployeeRepository) List(_ context.Context) ([]employee.Employee, error) {
	return r.Employees, r.Err
}

func (r *EmployeeRepository) FindByUsername(_ context.Context, username string) (employee.Employee, bool, error) {
	return r.find(func(e employee.Employee) bool { return e.GetUsername() == username })
}

func (r *EmployeeRepository) FindByID(_ context.Context, id string) (employee.Employee, bool, error) {
	return r.find(func(e employee.Employee) bool { return e.GetID() == id })
}

func (r *EmployeeRepository) find(match func(employee.Employee) bool) (employee.Employee, bool, error) {
	for _, e := range r.Employees {
		if match(e) {
			return e, true, nil
		}
	}
	return employee.Employee{}, false, r.Err
}

func (r *EmployeeRepository) Register(_ context.Context, e employee.Employee) error {
	if r.RegisterErr != nil {
		return r.RegisterErr
	}
	for _, existing := range r.Employees {
		if existing.GetUsername() == e.GetUsername() && existing.GetID() != e.GetID() {
			return domainErr.ErrUsernameTaken
		}
	}
	r.Employees = append(r.Employees, e)
	return nil
}

func (r *EmployeeRepository) Save(_ context.Context, e employee.Employee) error {
	for i, existing := range r.Employees {
		if existing.GetID() == e.GetID() {
			r.Employees[i] = e
			return nil
		}
	}
	return errors.New("employee not found")
}

type PasswordHasher struct {
	Comparisons int
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (h *PasswordHasher) Matches(hash, password string) bool {
	h.Comparisons++
	return hash == "hashed:"+password
}

type TokenIssuer struct{}

func (TokenIssuer) Issue(e employee.Employee) (string, time.Duration, error) {
	return "token-for-" + e.GetID(), 15 * time.Minute, nil
}

type SigningKeys struct{}

func (SigningKeys) PublicKeys() []out.PublicKey {
	return []out.PublicKey{{ID: "key-1", KeyType: "OKP", Curve: "Ed25519", Algorithm: "EdDSA", X: "public-x"}}
}

type RefreshTokenGenerator struct {
	generated int
}

func (g *RefreshTokenGenerator) Generate() (string, error) {
	g.generated++
	return "refresh-" + strconv.Itoa(g.generated), nil
}

func (g *RefreshTokenGenerator) Hash(token string) string {
	return "sha:" + token
}

type RefreshTokenStore struct {
	Tokens map[string]session.RefreshToken
}

func NewRefreshTokenStore() *RefreshTokenStore {
	return &RefreshTokenStore{Tokens: make(map[string]session.RefreshToken)}
}

func (s *RefreshTokenStore) Save(_ context.Context, token session.RefreshToken) error {
	s.Tokens[token.Hash()] = token
	return nil
}

func (s *RefreshTokenStore) FindByHash(_ context.Context, hash string) (session.RefreshToken, bool, error) {
	token, ok := s.Tokens[hash]
	return token, ok, nil
}

func (s *RefreshTokenStore) Revoke(_ context.Context, hash string) (bool, error) {
	token, ok := s.Tokens[hash]
	if !ok || token.IsRevoked() {
		return false, nil
	}
	token.Revoke()
	s.Tokens[hash] = token
	return true, nil
}

func (s *RefreshTokenStore) RevokeAllForEmployee(_ context.Context, employeeID string) error {
	for hash, token := range s.Tokens {
		if token.EmployeeID() == employeeID {
			token.Revoke()
			s.Tokens[hash] = token
		}
	}
	return nil
}

type TokenVerifier struct {
	Callers map[string]out.Caller
}

func (v TokenVerifier) Verify(token string) (out.Caller, error) {
	caller, ok := v.Callers[token]
	if !ok {
		return out.Caller{}, errors.New("unknown token")
	}
	return caller, nil
}

type PasswordGenerator struct {
	Password string
}

func (g PasswordGenerator) Generate() (string, error) {
	return g.Password, nil
}
