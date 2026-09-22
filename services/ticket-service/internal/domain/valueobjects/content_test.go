package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	content "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	t.Run("should create a new Content with a valid string", func(t *testing.T) {
		contentStr := "This is a test response"
		c, err := content.NewContent(contentStr)

		assert.NoError(t, err)
		assert.Equal(t, contentStr, c.GetContent())
	})

	t.Run("should return an error when creating a Content with an empty string", func(t *testing.T) {
		_, err := content.NewContent("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidContent)
	})

}
