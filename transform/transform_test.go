package transform

import (
	"github.com/stretchr/testify/assert"
	"switchboard-module-boilerplate/models"
	"testing"
)

func TestWordCount(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		_, result, err := Transform(models.Product{
			Description: "This is a simple sentence. It has multiple words.",
		})
		assert.Nil(t, err)
		assert.Equal(t, result.Attributes["word_count"], 9)
	})

	t.Run("Empty", func(t *testing.T) {
		_, result, err := Transform(models.Product{
			Description: "",
		})
		assert.Nil(t, err)
		assert.Equal(t, result.Attributes["word_count"], 0)
	})

	t.Run("Empty", func(t *testing.T) {
		_, result, err := Transform(models.Product{
			Description: "<p>This is a <b>simple sentence<b>. <br/> <b><i>It</i></b> has multiple words.</p>", // This has a mistake on purpose
		})
		assert.Nil(t, err)
		assert.Equal(t, result.Attributes["word_count"], 9)
	})

	t.Run("Avoid double processing", func(t *testing.T) {
		processed, result, err := Transform(models.Product{
			Description: "<p>This is a <b>simple sentence<b>. <br/> <b><i>It</i></b> has multiple words.</p>", // This has a mistake on purpose
			Attributes: map[string]interface{}{
				wordCountAttributeKey: 9,
			},
		})
		assert.Nil(t, err)
		assert.Equal(t, result.Attributes["word_count"], 9)
		assert.Equal(t, processed, false)
	})
}