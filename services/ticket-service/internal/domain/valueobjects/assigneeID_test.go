package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	assigneeID "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestAssigneeID(t *testing.T) {
	t.Run("should create a new AssigneeID with a valid string", func(t *testing.T) {
		idStr := "valid-assignee-id"
		idObj, err := assigneeID.NewAssigneeID(idStr)

		assert.NoError(t, err)
		assert.Equal(t, idStr, idObj.GetAssigneeID())
	})

	t.Run("should return an error when creating an AssigneeID with an empty string", func(t *testing.T) {
		_, err := assigneeID.NewAssigneeID("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})

	t.Run("should check validity of AssigneeID", func(t *testing.T) {
		validID, err := assigneeID.NewAssigneeID("valid-assignee-id")
		assert.NoError(t, err)

		assert.True(t, validID.IsValid())
	})

	t.Run("should check equality of two AssigneeIDs", func(t *testing.T) {
		id1, err := assigneeID.NewAssigneeID("same-assignee-id")
		assert.NoError(t, err)
		id2, err := assigneeID.NewAssigneeID("same-assignee-id")
		assert.NoError(t, err)
		id3, err := assigneeID.NewAssigneeID("different-assignee-id")
		assert.NoError(t, err)

		assert.True(t, id1.Equals(id2))
		assert.False(t, id1.Equals(id3))
	})

	t.Run("should marshal and unmarshal AssigneeID to/from JSON", func(t *testing.T) {
		idStr := "valid-assignee-id"
		idObj, err := assigneeID.NewAssigneeID(idStr)
		assert.NoError(t, err)

		jsonData, err := idObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledID assigneeID.AssigneeID
		err = unmarshaledID.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, idObj.GetAssigneeID(), unmarshaledID.GetAssigneeID())
	})

	t.Run("should return an error when unmarshaling an invalid AssigneeID from JSON", func(t *testing.T) {
		invalidIDJSON := []byte(`12345`)

		var unmarshaledID assigneeID.AssigneeID
		err := unmarshaledID.UnmarshalJSON(invalidIDJSON)

		assert.Error(t, err)
	})

	t.Run("should return an error when unmarshaling an empty AssigneeID from JSON", func(t *testing.T) {
		emptyIDJSON := []byte(`""`)

		var unmarshaledID assigneeID.AssigneeID
		err := unmarshaledID.UnmarshalJSON(emptyIDJSON)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})
}
