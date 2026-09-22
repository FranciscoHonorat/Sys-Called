package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	authorID "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestAuthorID(t *testing.T) {
	t.Run("should create a new AuthorID with a valid string", func(t *testing.T) {
		idStr := "valid-author-id"
		idObj, err := authorID.NewAuthorID(idStr)

		assert.NoError(t, err)
		assert.Equal(t, idStr, idObj.GetAuthorID())
	})

	t.Run("should return an error when creating an AuthorID with an empty string", func(t *testing.T) {
		_, err := authorID.NewAuthorID("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAuthorID)
	})

	t.Run("should check validity of AuthorID", func(t *testing.T) {
		validID, err := authorID.NewAuthorID("valid-author-id")
		assert.NoError(t, err)

		assert.True(t, validID.IsValid())
	})

	t.Run("should check equality of two AuthorIDs", func(t *testing.T) {
		id1, err := authorID.NewAuthorID("same-author-id")
		assert.NoError(t, err)
		id2, err := authorID.NewAuthorID("same-author-id")
		assert.NoError(t, err)
		id3, err := authorID.NewAuthorID("different-author-id")
		assert.NoError(t, err)

		assert.True(t, id1.Equals(id2))
		assert.False(t, id1.Equals(id3))
	})

	t.Run("should marshal and unmarshal AuthorID to/from JSON", func(t *testing.T) {
		idStr := "valid-author-id"
		idObj, err := authorID.NewAuthorID(idStr)
		assert.NoError(t, err)

		jsonData, err := idObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledID authorID.AuthorID
		err = unmarshaledID.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, idObj.GetAuthorID(), unmarshaledID.GetAuthorID())
	})

	t.Run("should return an error when unmarshaling an invalid AuthorID from JSON", func(t *testing.T) {
		invalidIDJSON := []byte(`12345`)

		var unmarshaledID authorID.AuthorID
		err := unmarshaledID.UnmarshalJSON(invalidIDJSON)

		assert.Error(t, err)
	})

	t.Run("should return an error when unmarshaling an empty AuthorID from JSON", func(t *testing.T) {
		emptyIDJSON := []byte(`""`)

		var unmarshaledID authorID.AuthorID
		err := unmarshaledID.UnmarshalJSON(emptyIDJSON)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAuthorID)
	})
}
