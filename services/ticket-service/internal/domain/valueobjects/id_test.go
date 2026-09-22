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

}
