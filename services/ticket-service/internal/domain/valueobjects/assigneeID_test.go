package valueobjects_test

import (
	"testing"

	assigneeID "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestAssigneeID(t *testing.T) {
	t.Run("should create a new AssigneeID with a valid string", func(t *testing.T) {
		idStr := "valid-assignee-id"
		idObj := assigneeID.NewAssigneeID(idStr)

		assert.Equal(t, idStr, idObj.GetAssigneeID())
	})

	t.Run("should check validity of AssigneeID", func(t *testing.T) {
		validID := assigneeID.NewAssigneeID("valid-assignee-id")
		invalidID := assigneeID.NewAssigneeID("")

		assert.True(t, validID.IsValid())
		assert.False(t, invalidID.IsValid())
	})

	t.Run("should check equality of two AssigneeIDs", func(t *testing.T) {
		id1 := assigneeID.NewAssigneeID("same-assignee-id")
		id2 := assigneeID.NewAssigneeID("same-assignee-id")
		id3 := assigneeID.NewAssigneeID("different-assignee-id")

		assert.True(t, id1.Equals(id2))
		assert.False(t, id1.Equals(id3))
	})

	t.Run("should marshal and unmarshal AssigneeID to/from JSON", func(t *testing.T) {
		idStr := "valid-assignee-id"
		idObj := assigneeID.NewAssigneeID(idStr)

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
}
