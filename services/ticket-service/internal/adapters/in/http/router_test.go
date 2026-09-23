package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	httpapi "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/auth"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	notifications := &outtest.NotificationStore{}
	store := application.NewNotifyingEventStore(outtest.NewEventStore(), notifications)
	c := cache.NewInMemoryTicketCache()
	responsibles := &outtest.ResponsibleDirectory{
		Responsibles: []string{"agent-1", "agent-2", "agent-3"},
		Names:        map[string]string{"agent-1": "Ana Souza", "agent-2": "Bruno Lima", "agent-3": "Carla Melo"},
	}
	admin, _ := actor.New("admin-1", "Administradora", "admin")
	verifier := &outtest.TokenVerifier{Tokens: map[string]actor.Actor{
		"valid-token": admin,
		"user-token":  outtest.Actor("user-1", "user"),
		"agent-token": outtest.Actor("agent-1", "support"),
	}}
	handler := httpapi.NewHandler(httpapi.UseCases{
		Authenticate:           auth.NewAuthenticateUseCase(verifier),
		OpenTicket:             command.NewOpenTicketUseCase(store, c),
		GetTicket:              query.NewGetTicketUseCase(store, c),
		ListTickets:            query.NewListTicketsUseCase(store, c),
		EditTicket:             command.NewEditTicketUseCase(store, c),
		AssignTicket:           command.NewAssignTicketUseCase(store, c),
		AutoAssignTicket:       command.NewAutoAssignTicketUseCase(store, c, responsibles),
		ChangeTicketPriority:   command.NewChangeTicketPriorityUseCase(store, c),
		MoveTicketToInProgress: command.NewMoveTicketToInProgressUseCase(store, c),
		CloseTicket:            command.NewCloseTicketUseCase(store, c),
		AddTicketResponse:      command.NewAddTicketResponseUseCase(store, c),
		ListResponsibles:       query.NewListResponsiblesUseCase(responsibles),
		ListNotifications:      query.NewListNotificationsUseCase(notifications),
		MarkNotificationsRead:  command.NewMarkNotificationsReadUseCase(notifications, time.Now),
		SupportWorkload:        query.NewSupportWorkloadUseCase(store, c, responsibles),
	})
	return httpapi.NewRouter(handler)
}

func doRequest(t *testing.T, router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	return doRequestAs(t, router, "valid-token", method, path, body)
}

func doRequestAs(t *testing.T, router *gin.Engine, token, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func TestRoutesRequireAuthentication(t *testing.T) {
	router := newTestRouter()

	t.Run("should protect ticket routes", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/tickets", nil))

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should keep the health check public", func(t *testing.T) {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestResponsiblesEndpoint(t *testing.T) {
	t.Run("should list the responsibles with names", func(t *testing.T) {
		w := doRequest(t, newTestRouter(), http.MethodGet, "/responsibles", nil)

		require.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":"agent-1","name":"Ana Souza"},{"id":"agent-2","name":"Bruno Lima"},{"id":"agent-3","name":"Carla Melo"}]`, w.Body.String())
	})

	t.Run("should require authentication", func(t *testing.T) {
		w := httptest.NewRecorder()
		newTestRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/responsibles", nil))

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestSupportWorkloadEndpoint(t *testing.T) {
	t.Run("gives admins the workload of each agent", func(t *testing.T) {
		w := doRequest(t, newTestRouter(), http.MethodGet, "/responsibles/workload", nil)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `{"id":"agent-1","name":"Ana Souza","open":0,"in_progress":0,"closed":0}`)
	})

	t.Run("is forbidden to anyone else", func(t *testing.T) {
		w := doRequestAs(t, newTestRouter(), "agent-token", http.MethodGet, "/responsibles/workload", nil)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestNotificationsEndpoints(t *testing.T) {
	router := newTestRouter()
	w := doRequestAs(t, router, "user-token", http.MethodPost, "/tickets", map[string]string{"title": "Impressora", "description": "Não imprime"})
	require.Equal(t, http.StatusCreated, w.Code)

	w = doRequest(t, router, http.MethodGet, "/notifications", nil)
	require.Equal(t, http.StatusOK, w.Code)
	var feed query.NotificationsOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &feed))
	assert.Equal(t, 1, feed.Unread)
	require.Len(t, feed.Items, 1)
	assert.Equal(t, "Novo chamado: Impressora", feed.Items[0].Message)

	w = doRequest(t, router, http.MethodPost, "/notifications/read", nil)
	require.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodGet, "/notifications", nil)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &feed))
	assert.Zero(t, feed.Unread)

	w = doRequestAs(t, router, "user-token", http.MethodGet, "/notifications", nil)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &feed))
	assert.Empty(t, feed.Items)
}

func TestForbiddenAction(t *testing.T) {
	router := newTestRouter()
	w := doRequestAs(t, router, "user-token", http.MethodPost, "/tickets", map[string]string{
		"title": "Printer", "description": "Broken",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	var opened command.OpenTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))

	w = doRequestAs(t, router, "user-token", http.MethodPost, "/tickets/"+opened.TicketID+"/assign", map[string]string{
		"assignee_id": "agent-1",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestTicketLifecycle(t *testing.T) {
	router := newTestRouter()

	w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
		"title":       "Valid Title",
		"description": "Valid Description",
	})
	require.Equal(t, http.StatusCreated, w.Code)

	var opened command.OpenTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))
	require.NotEmpty(t, opened.TicketID)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	var got query.GetTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "Valid Title", got.Title)
	assert.Equal(t, "Open", got.Status)

	w = doRequest(t, router, http.MethodPut, "/tickets/"+opened.TicketID, map[string]string{
		"title":       "Edited Title",
		"description": "Edited Description",
	})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "Edited Title", got.Title)
	assert.Equal(t, "Edited Description", got.Description)

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/assign", map[string]string{"assignee_id": "agent-1"})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/priority", map[string]string{"priority": "High"})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequestAs(t, router, "agent-token", http.MethodPost, "/tickets/"+opened.TicketID+"/start", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/responses", map[string]string{"content": "Working on it"})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "In Progress", got.Status)
	assert.Equal(t, "agent-1", got.AssigneeID)
	assert.Equal(t, "High", got.Priority)
	require.Len(t, got.Responses, 1)
	assert.Equal(t, "Working on it", got.Responses[0].Content)

	w = doRequestAs(t, router, "agent-token", http.MethodPost, "/tickets/"+opened.TicketID+"/close", map[string]string{"resolution": "Troquei o cabo"})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "Closed", got.Status)
}

func TestListTickets(t *testing.T) {
	router := newTestRouter()

	w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
		"title":       "Valid Title",
		"description": "Valid Description",
	})
	require.Equal(t, http.StatusCreated, w.Code)

	w = doRequest(t, router, http.MethodGet, "/tickets", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var tickets []query.GetTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tickets))
	require.Len(t, tickets, 1)
	assert.Equal(t, "Valid Title", tickets[0].Title)
}

func TestAutoAssignTicket(t *testing.T) {
	router := newTestRouter()

	w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
		"title":       "Valid Title",
		"description": "Valid Description",
	})
	require.Equal(t, http.StatusCreated, w.Code)
	var opened command.OpenTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/assign/auto", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var autoAssigned command.AutoAssignTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &autoAssigned))
	assert.Contains(t, []string{"agent-1", "agent-2", "agent-3"}, autoAssigned.AssigneeID)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	var got query.GetTicketOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, autoAssigned.AssigneeID, got.AssigneeID)
}

func TestTicketErrors(t *testing.T) {
	t.Run("should return 400 for an invalid title", func(t *testing.T) {
		router := newTestRouter()

		w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
			"title":       "",
			"description": "Valid Description",
		})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 404 for a ticket that does not exist", func(t *testing.T) {
		router := newTestRouter()

		w := doRequest(t, router, http.MethodGet, "/tickets/00000000-0000-0000-0000-000000000001", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 when closing without saying what was done", func(t *testing.T) {
		router := newTestRouter()
		w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{"title": "Valid Title", "description": "Valid Description"})
		require.Equal(t, http.StatusCreated, w.Code)
		var opened command.OpenTicketOutput
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))
		w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/assign", map[string]string{"assignee_id": "agent-1"})
		require.Equal(t, http.StatusNoContent, w.Code)

		w = doRequestAs(t, router, "agent-token", http.MethodPost, "/tickets/"+opened.TicketID+"/close", map[string]string{"resolution": ""})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 409 when closing an already closed ticket", func(t *testing.T) {
		router := newTestRouter()

		w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
			"title":       "Valid Title",
			"description": "Valid Description",
		})
		require.Equal(t, http.StatusCreated, w.Code)
		var opened command.OpenTicketOutput
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))

		w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/assign", map[string]string{"assignee_id": "agent-1"})
		require.Equal(t, http.StatusNoContent, w.Code)
		w = doRequestAs(t, router, "agent-token", http.MethodPost, "/tickets/"+opened.TicketID+"/close", map[string]string{"resolution": "Troquei o cabo"})
		require.Equal(t, http.StatusNoContent, w.Code)

		w = doRequestAs(t, router, "agent-token", http.MethodPost, "/tickets/"+opened.TicketID+"/close", map[string]string{"resolution": "Troquei o cabo"})
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
