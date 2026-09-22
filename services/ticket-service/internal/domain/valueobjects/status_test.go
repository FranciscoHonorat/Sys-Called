package valueobjects_test

import (
	"testing"

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

}
