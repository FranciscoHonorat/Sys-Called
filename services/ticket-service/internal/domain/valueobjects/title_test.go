package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	title "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	t.Run("should create a new Title with a valid string", func(t *testing.T) {
		titleStr := "Valid Title"
		titleObj, err := title.NewTitle(titleStr)

		assert.NoError(t, err)
		assert.Equal(t, titleStr, titleObj.GetTitle())
	})

	t.Run("should return an error when creating a Title with an empty string", func(t *testing.T) {
		_, err := title.NewTitle("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

}
