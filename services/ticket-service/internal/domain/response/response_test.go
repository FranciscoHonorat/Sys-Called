package response_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func validResponseParts(t *testing.T) (*valueobjects.ID, *valueobjects.ID, *valueobjects.AuthorID, *valueobjects.Content) {
	t.Helper()

	id := valueobjects.NewID(uuid.New())
	ticketID := valueobjects.NewID(uuid.New())

	author, err := valueobjects.NewAuthorID("agent-1")
	assert.NoError(t, err)

	content, err := valueobjects.NewContent("Valid response content")
	assert.NoError(t, err)

	return id, ticketID, &author, content
}

func TestResponse(t *testing.T) {
	t.Run("should create a new Response with valid fields", func(t *testing.T) {
		id, ticketID, author, content := validResponseParts(t)

		r, err := response.NewResponse(id, ticketID, author, content)

		assert.NoError(t, err)
		assert.Equal(t, id, r.GetID())
		assert.Equal(t, ticketID, r.GetTicketID())
		assert.Equal(t, author, r.GetAuthorID())
		assert.Equal(t, content, r.GetContent())
		assert.WithinDuration(t, time.Now(), r.GetCreatedAt(), time.Second)
	})

	t.Run("should return an error when id is nil", func(t *testing.T) {
		_, ticketID, author, content := validResponseParts(t)

		r, err := response.NewResponse(nil, ticketID, author, content)

		assert.Nil(t, r)
		assert.ErrorIs(t, err, domainErr.ErrInvalidID)
	})

	t.Run("should return an error when ticketID is nil", func(t *testing.T) {
		id, _, author, content := validResponseParts(t)

		r, err := response.NewResponse(id, nil, author, content)

		assert.Nil(t, r)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTicketID)
	})

	t.Run("should return an error when authorID is nil", func(t *testing.T) {
		id, ticketID, _, content := validResponseParts(t)

		r, err := response.NewResponse(id, ticketID, nil, content)

		assert.Nil(t, r)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAuthorID)
	})

	t.Run("should return an error when content is nil", func(t *testing.T) {
		id, ticketID, author, _ := validResponseParts(t)

		r, err := response.NewResponse(id, ticketID, author, nil)

		assert.Nil(t, r)
		assert.ErrorIs(t, err, domainErr.ErrInvalidContent)
	})

	t.Run("should marshal a Response to JSON", func(t *testing.T) {
		id, ticketID, author, content := validResponseParts(t)
		r, err := response.NewResponse(id, ticketID, author, content)
		assert.NoError(t, err)

		data, err := r.MarshalJSON()

		assert.NoError(t, err)
		assert.Contains(t, string(data), `"id":"`+id.String()+`"`)
		assert.Contains(t, string(data), `"ticket_id":"`+ticketID.String()+`"`)
		assert.Contains(t, string(data), `"author_id":"agent-1"`)
		assert.Contains(t, string(data), `"content":"Valid response content"`)
	})

	t.Run("should round-trip marshal and unmarshal", func(t *testing.T) {
		id, ticketID, author, content := validResponseParts(t)
		original, err := response.NewResponse(id, ticketID, author, content)
		assert.NoError(t, err)

		data, err := original.MarshalJSON()
		assert.NoError(t, err)

		var decoded response.Response
		err = decoded.UnmarshalJSON(data)
		assert.NoError(t, err)

		assert.Equal(t, original.GetID().GetID(), decoded.GetID().GetID())
		assert.Equal(t, original.GetTicketID().GetID(), decoded.GetTicketID().GetID())
		assert.Equal(t, original.GetAuthorID().GetAuthorID(), decoded.GetAuthorID().GetAuthorID())
		assert.Equal(t, original.GetContent().GetContent(), decoded.GetContent().GetContent())
	})

	t.Run("should return an error when unmarshaling without a ticket_id", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","author_id":"agent-1","content":"Valid content","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var r response.Response
		err := r.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidTicketID)
	})

	t.Run("should return an error when unmarshaling an invalid ticket_id", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","ticket_id":"not-a-uuid","author_id":"agent-1","content":"Valid content","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var r response.Response
		err := r.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when unmarshaling an empty author_id", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","ticket_id":"` + uuid.New().String() + `","author_id":"","content":"Valid content","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var r response.Response
		err := r.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidAuthorID)
	})

	t.Run("should return an error when unmarshaling an empty content", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","ticket_id":"` + uuid.New().String() + `","author_id":"agent-1","content":"","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var r response.Response
		err := r.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidContent)
	})

	t.Run("should return an error when unmarshaling a malformed created_at", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","ticket_id":"` + uuid.New().String() + `","author_id":"agent-1","content":"Valid content","created_at":"not-a-date"}`)

		var r response.Response
		err := r.UnmarshalJSON(data)

		assert.Error(t, err)
	})
}
