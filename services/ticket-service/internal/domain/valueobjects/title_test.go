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

	t.Run("should check validity of Title", func(t *testing.T) {
		validTitle := &title.Title{Title: "Valid Title"}
		invalidTitle := &title.Title{Title: ""}

		assert.True(t, validTitle.IsValid())
		assert.False(t, invalidTitle.IsValid())
	})

	t.Run("should check equality of two Titles", func(t *testing.T) {
		title1, err := title.NewTitle("Same Title")
		assert.NoError(t, err)
		title2, err := title.NewTitle("Same Title")
		assert.NoError(t, err)
		title3, err := title.NewTitle("Different Title")
		assert.NoError(t, err)

		assert.True(t, title1.Equals(title2))
		assert.False(t, title1.Equals(title3))
	})

	t.Run("should marshal and unmarshal Title to/from JSON", func(t *testing.T) {
		titleStr := "Valid Title"
		titleObj, err := title.NewTitle(titleStr)
		assert.NoError(t, err)

		jsonData, err := titleObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledTitle title.Title
		err = unmarshaledTitle.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, titleObj.GetTitle(), unmarshaledTitle.GetTitle())
	})
}
