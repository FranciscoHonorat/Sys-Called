package valueobjects_test

import (
	"testing"

	title "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	t.Run("should create a new Title with a valid string", func(t *testing.T) {
		titleStr := "Valid Title"
		titleObj := title.NewTitle(titleStr)

		assert.Equal(t, titleStr, titleObj.GetTitle())
	})

	t.Run("should check validity of Title", func(t *testing.T) {
		validTitle := title.NewTitle("Valid Title")
		invalidTitle := title.NewTitle("")

		assert.True(t, validTitle.IsValid())
		assert.False(t, invalidTitle.IsValid())
	})

	t.Run("should check equality of two Titles", func(t *testing.T) {
		title1 := title.NewTitle("Same Title")
		title2 := title.NewTitle("Same Title")
		title3 := title.NewTitle("Different Title")

		assert.True(t, title1.Equals(title2))
		assert.False(t, title1.Equals(title3))
	})

	t.Run("should marshal and unmarshal Title to/from JSON", func(t *testing.T) {
		titleStr := "Valid Title"
		titleObj := title.NewTitle(titleStr)

		jsonData, err := titleObj.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaledTitle title.Title
		err = unmarshaledTitle.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, titleObj.GetTitle(), unmarshaledTitle.GetTitle())
	})
}
