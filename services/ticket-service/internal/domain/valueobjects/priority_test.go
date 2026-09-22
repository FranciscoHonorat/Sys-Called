package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	priority "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestPriority(t *testing.T) {
	t.Run("should create a new Priority with a valid string", func(t *testing.T) {
		priorityStr := "High"
		priorityObj := priority.NewPriority(priorityStr)

		assert.Equal(t, priorityStr, priorityObj.GetPriority())
	})

	t.Run("should check validity of Priority", func(t *testing.T) {
		validPriority := priority.TicketPriorityMedium
		invalidPriority := priority.Priority("Invalid")

		assert.True(t, validPriority.IsValid())
		assert.False(t, invalidPriority.IsValid())
	})

	t.Run("should check equality of two Priorities", func(t *testing.T) {
		priority1 := priority.TicketPriorityLow
		priority2 := priority.TicketPriorityLow
		priority3 := priority.TicketPriorityHigh

		assert.True(t, priority1.Equals(priority2))
		assert.False(t, priority1.Equals(priority3))
	})

	t.Run("should marshal and unmarshal Priority to/from JSON", func(t *testing.T) {
		priorityStr := "Medium"
		priorityObj := priority.NewPriority(priorityStr)

		jsonData, err := priorityObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledPriority priority.Priority
		err = unmarshaledPriority.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, priorityObj.GetPriority(), unmarshaledPriority.GetPriority())
	})

	t.Run("should return an error when unmarshaling an invalid Priority from JSON", func(t *testing.T) {
		invalidPriorityJSON := []byte(`"Invalid"`)

		var unmarshaledPriority priority.Priority
		err := unmarshaledPriority.UnmarshalJSON(invalidPriorityJSON)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})
}
