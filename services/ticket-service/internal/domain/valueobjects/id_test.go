package valueobjects_test

import (
	"testing"

	"github.com/google/uuid"

	id "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	t.Run("should create a new ID with a valid UUID", func(t *testing.T) {
		uuidValue := uuid.New()
		idObj := id.NewID(uuidValue)

		assert.Equal(t, uuidValue, idObj.GetID())
		assert.Equal(t, uuidValue.String(), idObj.String())
	})

	t.Run("should create a new ID with a new UUID if nil is provided", func(t *testing.T) {
		idObj := id.NewID(uuid.Nil)

		assert.NotEqual(t, uuid.Nil, idObj.GetID())
	})

	t.Run("should check equality of two IDs", func(t *testing.T) {
		uuidValue := uuid.New()
		id1 := id.NewID(uuidValue)
		id2 := id.NewID(uuidValue)
		id3 := id.NewID(uuid.New())

		assert.True(t, id1.Equals(id2))
		assert.False(t, id1.Equals(id3))
	})

	t.Run("should marshal and unmarshal ID to/from JSON", func(t *testing.T) {
		uuidValue := uuid.New()
		idObj := id.NewID(uuidValue)

		jsonData, err := idObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledID id.ID
		err = unmarshaledID.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, idObj.GetID(), unmarshaledID.GetID())
	})

	t.Run("should handle unmarshaling of nil UUID", func(t *testing.T) {
		jsonData := []byte(`{"id":"00000000-0000-0000-0000-000000000000"}`)

		var unmarshaledID id.ID
		err := unmarshaledID.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, uuid.Nil, unmarshaledID.GetID())
	})
}
