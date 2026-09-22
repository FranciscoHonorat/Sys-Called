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

}
