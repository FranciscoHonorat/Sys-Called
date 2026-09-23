package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func refreshCookie(t *testing.T, w *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	for _, c := range w.Result().Cookies() {
		if c.Name == "refresh_token" {
			return c
		}
	}
	t.Fatal("refresh_token cookie not set")
	return nil
}

func login(t *testing.T, router *gin.Engine) *http.Cookie {
	t.Helper()
	w := doRequest(router, http.MethodPost, "/auth/login", `{"username":"ana","password":"senha123"}`)
	require.Equal(t, http.StatusOK, w.Code)
	return refreshCookie(t, w)
}

func TestSessionCookie(t *testing.T) {
	t.Run("should keep the refresh token in a secure HttpOnly cookie and out of the body", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/login", `{"username":"ana","password":"senha123"}`)

		cookie := refreshCookie(t, w)
		assert.Equal(t, "refresh-1", cookie.Value)
		assert.True(t, cookie.HttpOnly)
		assert.True(t, cookie.Secure)
		assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
		assert.Equal(t, "/auth", cookie.Path)
		assert.Equal(t, 7*24*3600, cookie.MaxAge)
		assert.NotContains(t, w.Body.String(), "refresh-1")
	})
}

func TestRefreshEndpoint(t *testing.T) {
	t.Run("should rotate the session from the refresh cookie", func(t *testing.T) {
		router := newTestRouter()
		cookie := login(t, router)

		w := doRequest(router, http.MethodPost, "/auth/refresh", "", cookie)

		require.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"access_token":"token-for-agent-1","token_type":"Bearer","expires_in":900}`, w.Body.String())
		assert.Equal(t, "refresh-2", refreshCookie(t, w).Value)
	})

	t.Run("should return 401 without a refresh cookie", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/refresh", "")

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"invalid refresh token"}`, w.Body.String())
	})

	t.Run("should return 401 and clear the cookie for an invalid refresh token", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/refresh", "", &http.Cookie{Name: "refresh_token", Value: "forged"})

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Negative(t, refreshCookie(t, w).MaxAge)
	})
}

func TestLogoutEndpoint(t *testing.T) {
	t.Run("should end the session and clear the cookie", func(t *testing.T) {
		router := newTestRouter()
		cookie := login(t, router)

		w := doRequest(router, http.MethodPost, "/auth/logout", "", cookie)

		require.Equal(t, http.StatusNoContent, w.Code)
		assert.Negative(t, refreshCookie(t, w).MaxAge)
		assert.Equal(t, http.StatusUnauthorized, doRequest(router, http.MethodPost, "/auth/refresh", "", cookie).Code)
	})

	t.Run("should succeed even without a session", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/logout", "")

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
