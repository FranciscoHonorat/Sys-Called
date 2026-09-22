package valueobjects_test

import (
	"testing"

	des "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	t.Run("should create a new description", func(t *testing.T) {
		description := "This is a test description"
		desc := des.NewDescription(description)

		assert.Equal(t, description, desc.GetDescription())
	})

	t.Run("should validate a valid description", func(t *testing.T) {
		description := "This is a test description"
		desc := des.NewDescription(description)

		assert.True(t, desc.IsValid())
	})

	t.Run("should invalidate an empty description", func(t *testing.T) {
		description := ""
		desc := des.NewDescription(description)

		assert.False(t, desc.IsValid())
	})

	t.Run("should compare two equal descriptions", func(t *testing.T) {
		description := "This is a test description"
		desc1 := des.NewDescription(description)
		desc2 := des.NewDescription(description)

		assert.True(t, desc1.Equals(desc2))
	})

	t.Run("should compare two different descriptions", func(t *testing.T) {
		description1 := "This is a test description"
		description2 := "This is another test description"
		desc1 := des.NewDescription(description1)
		desc2 := des.NewDescription(description2)

		assert.False(t, desc1.Equals(desc2))
	})
}
