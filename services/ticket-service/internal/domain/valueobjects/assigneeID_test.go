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

}
