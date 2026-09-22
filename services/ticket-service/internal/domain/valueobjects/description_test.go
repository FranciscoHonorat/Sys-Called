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

	t.Run("should validate a valid description", func(t *testing.T) {
		desc := &des.Description{Description: "This is a test description"}

		assert.True(t, desc.IsValid())
	})

	t.Run("should invalidate an empty description", func(t *testing.T) {
		desc := &des.Description{Description: ""}

		assert.False(t, desc.IsValid())
	})

	t.Run("should compare two equal descriptions", func(t *testing.T) {
		description := "This is a test description"
		desc1, err := des.NewDescription(description)
		assert.NoError(t, err)
		desc2, err := des.NewDescription(description)
		assert.NoError(t, err)

		assert.True(t, desc1.Equals(desc2))
	})

	t.Run("should compare two different descriptions", func(t *testing.T) {
		description1 := "This is a test description"
		description2 := "This is another test description"
		desc1, err := des.NewDescription(description1)
		assert.NoError(t, err)
		desc2, err := des.NewDescription(description2)
		assert.NoError(t, err)

		assert.False(t, desc1.Equals(desc2))
	})
}
