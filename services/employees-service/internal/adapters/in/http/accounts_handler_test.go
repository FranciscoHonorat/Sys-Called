package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

func doAuthorizedRequest(router *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSignUpEndpoint(t *testing.T) {
	t.Run("creates an account waiting for approval", func(t *testing.T) {
		router := newTestRouter()

		w := doRequest(router, http.MethodPost, "/auth/signup", `{"name":"Maria Lima","username":"maria","password":"senha-forte"}`)
		require.Equal(t, http.StatusCreated, w.Code)

		login := doRequest(router, http.MethodPost, "/auth/login", `{"username":"maria","password":"senha-forte"}`)
		assert.Equal(t, http.StatusForbidden, login.Code)
		assert.JSONEq(t, `{"error":"account waiting for the administrator approval"}`, login.Body.String())
	})

	t.Run("refuses a username already in use", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/signup", `{"name":"Outra Ana","username":"ana","password":"senha-forte"}`)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("refuses a weak password", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/signup", `{"name":"Maria Lima","username":"maria","password":"123"}`)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPasswordResetRequestEndpoint(t *testing.T) {
	t.Run("always accepts so usernames are not revealed", func(t *testing.T) {
		router := newTestRouter()

		assert.Equal(t, http.StatusAccepted, doRequest(router, http.MethodPost, "/auth/password-reset-requests", `{"username":"ana"}`).Code)
		assert.Equal(t, http.StatusAccepted, doRequest(router, http.MethodPost, "/auth/password-reset-requests", `{"username":"ghost"}`).Code)
	})
}

func TestApproveEndpoint(t *testing.T) {
	t.Run("lets the admin approve a pending account", func(t *testing.T) {
		router := newTestRouter()
		doRequest(router, http.MethodPost, "/auth/signup", `{"name":"Maria Lima","username":"maria","password":"senha-forte"}`)
		id := employeeIDByUsername(t, router, "maria")

		w := doAuthorizedRequest(router, http.MethodPost, "/employees/"+id+"/approve", "admin-token", "")
		require.Equal(t, http.StatusNoContent, w.Code)

		login := doRequest(router, http.MethodPost, "/auth/login", `{"username":"maria","password":"senha-forte"}`)
		assert.Equal(t, http.StatusOK, login.Code)
	})

	t.Run("is only for admins", func(t *testing.T) {
		w := doAuthorizedRequest(newTestRouter(), http.MethodPost, "/employees/agent-2/approve", "agent-token", "")

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("answers 404 for an unknown employee", func(t *testing.T) {
		w := doAuthorizedRequest(newTestRouter(), http.MethodPost, "/employees/ghost/approve", "admin-token", "")

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTemporaryPasswordEndpoint(t *testing.T) {
	t.Run("gives the admin a temporary password once", func(t *testing.T) {
		router := newTestRouter()

		w := doAuthorizedRequest(router, http.MethodPost, "/employees/agent-2/temporary-password", "admin-token", "")
		require.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"temporary_password":"Temp-1234"}`, w.Body.String())

		login := doRequest(router, http.MethodPost, "/auth/login", `{"username":"bruno","password":"Temp-1234"}`)
		assert.Equal(t, http.StatusOK, login.Code)
	})

	t.Run("is only for admins", func(t *testing.T) {
		w := doAuthorizedRequest(newTestRouter(), http.MethodPost, "/employees/agent-2/temporary-password", "agent-token", "")

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestChangePasswordEndpoint(t *testing.T) {
	t.Run("changes the password of whoever is logged in", func(t *testing.T) {
		router := newTestRouter()

		w := doAuthorizedRequest(router, http.MethodPost, "/auth/change-password", "agent-token", `{"current_password":"senha123","new_password":"outra-senha"}`)
		require.Equal(t, http.StatusNoContent, w.Code)

		login := doRequest(router, http.MethodPost, "/auth/login", `{"username":"ana","password":"outra-senha"}`)
		assert.Equal(t, http.StatusOK, login.Code)
	})

	t.Run("refuses a wrong current password", func(t *testing.T) {
		w := doAuthorizedRequest(newTestRouter(), http.MethodPost, "/auth/change-password", "agent-token", `{"current_password":"errada","new_password":"outra-senha"}`)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("refuses a weak new password", func(t *testing.T) {
		w := doAuthorizedRequest(newTestRouter(), http.MethodPost, "/auth/change-password", "agent-token", `{"current_password":"senha123","new_password":"123"}`)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("requires authentication", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/change-password", `{"current_password":"senha123","new_password":"outra-senha"}`)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func employeeIDByUsername(t *testing.T, router *gin.Engine, username string) string {
	t.Helper()
	w := doAuthorizedRequest(router, http.MethodGet, "/employees", "admin-token", "")
	require.Equal(t, http.StatusOK, w.Code)
	var employees []application.EmployeeOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &employees))
	for _, e := range employees {
		if e.Username == username {
			return e.ID
		}
	}
	t.Fatalf("no employee with username %s", username)
	return ""
}
