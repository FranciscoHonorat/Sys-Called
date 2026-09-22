package valueobjects_test

import (
	"testing"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	des "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	t.Run("should create a new description", func(t *testing.T) {
		description := "This is a test description"
		desc, err := des.NewDescription(description)

		assert.NoError(t, err)
		assert.Equal(t, description, desc.GetDescription())
	})

	t.Run("should return an error when creating a description with an empty string", func(t *testing.T) {
		_, err := des.NewDescription("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidDescription)
	})

}
