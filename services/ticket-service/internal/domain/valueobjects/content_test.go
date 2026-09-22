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

	t.Run("should validate a valid Content", func(t *testing.T) {
		c := &content.Content{Content: "This is a test response"}

		assert.True(t, c.IsValid())
	})

	t.Run("should invalidate an empty Content", func(t *testing.T) {
		c := &content.Content{Content: ""}

		assert.False(t, c.IsValid())
	})

	t.Run("should compare two equal Contents", func(t *testing.T) {
		contentStr := "This is a test response"
		c1, err := content.NewContent(contentStr)
		assert.NoError(t, err)
		c2, err := content.NewContent(contentStr)
		assert.NoError(t, err)

		assert.True(t, c1.Equals(c2))
	})

	t.Run("should compare two different Contents", func(t *testing.T) {
		c1, err := content.NewContent("This is a test response")
		assert.NoError(t, err)
		c2, err := content.NewContent("This is another test response")
		assert.NoError(t, err)

		assert.False(t, c1.Equals(c2))
	})

	t.Run("should marshal and unmarshal Content to/from JSON", func(t *testing.T) {
		contentStr := "This is a test response"
		c, err := content.NewContent(contentStr)
		assert.NoError(t, err)

		jsonData, err := c.MarshalJSON()
		assert.NoError(t, err)

		var unmarshaled content.Content
		err = unmarshaled.UnmarshalJSON(jsonData)
		assert.NoError(t, err)

		assert.Equal(t, c.GetContent(), unmarshaled.GetContent())
	})

	t.Run("should return an error when unmarshaling an empty Content from JSON", func(t *testing.T) {
		emptyContentJSON := []byte(`""`)

		var unmarshaled content.Content
		err := unmarshaled.UnmarshalJSON(emptyContentJSON)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrInvalidContent)
	})
}
