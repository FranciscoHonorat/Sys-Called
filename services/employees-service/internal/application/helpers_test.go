package application_test

import (
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
)

var testNow = time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)

const testRefreshTTL = 7 * 24 * time.Hour

func newTestSessions(store *outtest.RefreshTokenStore) *application.Sessions {
	return newTestSessionsAt(store, testNow)
}

func newTestSessionsAt(store *outtest.RefreshTokenStore, now time.Time) *application.Sessions {
	return application.NewSessions(
		outtest.TokenIssuer{},
		&outtest.RefreshTokenGenerator{},
		store,
		testRefreshTTL,
		func() time.Time { return now },
	)
}
