package application

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

type SessionOutput struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"-"`
	RefreshExpiresIn int    `json:"-"`
}

type Sessions struct {
	issuer     out.TokenIssuer
	generator  out.RefreshTokenGenerator
	store      out.RefreshTokenStore
	refreshTTL time.Duration
	now        func() time.Time
}

func NewSessions(issuer out.TokenIssuer, generator out.RefreshTokenGenerator, store out.RefreshTokenStore, refreshTTL time.Duration, now func() time.Time) *Sessions {
	return &Sessions{issuer: issuer, generator: generator, store: store, refreshTTL: refreshTTL, now: now}
}

func (s *Sessions) Start(ctx context.Context, e employee.Employee) (SessionOutput, error) {
	accessToken, expiresIn, err := s.issuer.Issue(e)
	if err != nil {
		return SessionOutput{}, err
	}

	refreshToken, err := s.generator.Generate()
	if err != nil {
		return SessionOutput{}, err
	}
	stored := session.NewRefreshToken(s.generator.Hash(refreshToken), e.GetID(), s.now().Add(s.refreshTTL))
	if err := s.store.Save(ctx, stored); err != nil {
		return SessionOutput{}, err
	}

	return SessionOutput{
		AccessToken:      accessToken,
		TokenType:        "Bearer",
		ExpiresIn:        int(expiresIn.Seconds()),
		RefreshToken:     refreshToken,
		RefreshExpiresIn: int(s.refreshTTL.Seconds()),
	}, nil
}

func (s *Sessions) Consume(ctx context.Context, refreshToken string) (string, error) {
	hash := s.generator.Hash(refreshToken)
	stored, found, err := s.store.FindByHash(ctx, hash)
	if err != nil {
		return "", err
	}
	if !found || stored.IsExpired(s.now()) {
		return "", domainErr.ErrInvalidRefreshToken
	}
	revoked, err := s.store.Revoke(ctx, hash)
	if err != nil {
		return "", err
	}
	if !revoked {
		if err := s.store.RevokeAllForEmployee(ctx, stored.EmployeeID()); err != nil {
			return "", err
		}
		return "", domainErr.ErrInvalidRefreshToken
	}
	return stored.EmployeeID(), nil
}

func (s *Sessions) EndAll(ctx context.Context, employeeID string) error {
	return s.store.RevokeAllForEmployee(ctx, employeeID)
}

func (s *Sessions) End(ctx context.Context, refreshToken string) error {
	_, err := s.store.Revoke(ctx, s.generator.Hash(refreshToken))
	return err
}
