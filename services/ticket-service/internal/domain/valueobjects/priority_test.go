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
		priorityObj, err := priority.NewPriority(priorityStr)

		assert.NoError(t, err)
		assert.Equal(t, priorityStr, priorityObj.GetPriority())
	})

	t.Run("should return an error when creating a Priority with an invalid string", func(t *testing.T) {
		_, err := priority.NewPriority("Invalid")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})

	t.Run("should check validity of Priority", func(t *testing.T) {
		validPriority := priority.TicketPriorityMedium
		invalidPriority := priority.Priority("Invalid")

		assert.True(t, validPriority.IsValid())
		assert.False(t, invalidPriority.IsValid())
	})

}
