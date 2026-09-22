package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	httpapi "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	store := outtest.NewEventStore()
	c := cache.NewInMemoryTicketCache()
	responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1", "agent-2", "agent-3"}}
	handler := httpapi.NewHandler(httpapi.UseCases{
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
	})
	return httpapi.NewRouter(handler)
}

func doRequest(t *testing.T, router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
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
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
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

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/start", nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/responses", map[string]string{"author_id": "agent-1", "content": "Working on it"})
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = doRequest(t, router, http.MethodGet, "/tickets/"+opened.TicketID, nil)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "In Progress", got.Status)
	assert.Equal(t, "agent-1", got.AssigneeID)
	assert.Equal(t, "High", got.Priority)
	require.Len(t, got.Responses, 1)
	assert.Equal(t, "Working on it", got.Responses[0].Content)

	w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/close", nil)
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

	t.Run("should return 409 when closing an already closed ticket", func(t *testing.T) {
		router := newTestRouter()

		w := doRequest(t, router, http.MethodPost, "/tickets", map[string]string{
			"title":       "Valid Title",
			"description": "Valid Description",
		})
		require.Equal(t, http.StatusCreated, w.Code)
		var opened command.OpenTicketOutput
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &opened))

		w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/close", nil)
		require.Equal(t, http.StatusNoContent, w.Code)

		w = doRequest(t, router, http.MethodPost, "/tickets/"+opened.TicketID+"/close", nil)
		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
