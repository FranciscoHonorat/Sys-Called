package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	status "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	t.Run("should create a new Status with a valid string", func(t *testing.T) {
		statusStr := "Open"
		statusObj := status.Status(statusStr)

		assert.Equal(t, statusStr, statusObj.String())
	})

	t.Run("should check validity of Status", func(t *testing.T) {
		validStatus := status.TicketStatusOpen
		invalidStatus := status.Status("Invalid")

		assert.True(t, validStatus.IsValid())
		assert.False(t, invalidStatus.IsValid())
	})

	t.Run("should check equality of two Statuses", func(t *testing.T) {
		status1 := status.TicketStatusOpen
		status2 := status.TicketStatusOpen
		status3 := status.TicketStatusClosed

		assert.True(t, status1.Equals(status2))
		assert.False(t, status1.Equals(status3))
	})

	t.Run("should marshal and unmarshal Status to/from JSON", func(t *testing.T) {
		statusStr := "In Progress"
		statusObj := status.Status(statusStr)

		jsonData, err := statusObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledStatus status.Status
		err = unmarshaledStatus.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, statusObj.String(), unmarshaledStatus.String())
	})

	t.Run("should return an error when unmarshaling an invalid Status from JSON", func(t *testing.T) {
		invalidStatusJSON := []byte(`"Invalid"`)

		var unmarshaledStatus status.Status
		err := unmarshaledStatus.UnmarshalJSON(invalidStatusJSON)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidStatus)
	})
}
