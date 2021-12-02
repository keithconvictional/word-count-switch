package load

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"switchboard-module-boilerplate/models"
	"testing"
)

func Test_CreateUpdateProductPayload(t *testing.T) {
	t.Run("Temp", func(t *testing.T) {
		product1 := models.Product{
			Description: "Sample 1",
			Variants: []models.Variant{
				{
					Title: "Variant 1",
				},
				{
					Title: "Variant 2a",
				},
			},
		}
		product2 := models.Product{
			Description: "Sample 2",
			Variants: []models.Variant{
				{
					Title: "Variant 1",
					Attributes: map[string]interface{}{
						"sample": "word_count",
					},
				},
				{
					Title: "Variant 2",
				},
			},
		}
		_, updateProductPayload := CreateUpdateProductPayload(product1, product2)
		fmt.Printf("updateProductPayload :: %+v\n", updateProductPayload)
		assert.Equal(t, "Sample 2", updateProductPayload.BodyHTML)
		assert.Equal(t, 2, len(updateProductPayload.Variants))
		assert.Equal(t, "Variant 1", updateProductPayload.Variants[0].Title)
		assert.Equal(t, map[string]interface{}{
			"sample": "word_count",
		}, updateProductPayload.Variants[0].Attributes)
		assert.Equal(t, "Variant 2", updateProductPayload.Variants[1].Title)
	})

	t.Run("Different", func(t *testing.T) {

	})

	t.Run("Different (Nested)", func(t *testing.T) {

	})

	t.Run("Different (Nested Array)", func(t *testing.T) {

	})

	t.Run("Same", func(t *testing.T) {

	})
}

func Test_SetFieldWithReflect(t *testing.T) {
	t.Run("Root Field (Handle field does not exist on object)", func(t *testing.T) {
		updatedProduct := models.UpdateProduct{
			Attributes: map[string]interface{}{
				"sample": "info",
			},
		}
		updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
		// Nothing should be updated, but should not fail (original object is returned)
		_ = SetFieldWithReflect("MadeUpField.ChildField", updatedProductWithReflect, "NEW VALUE")
	})

	t.Run("Update on object", func(t *testing.T) {
		updatedProduct := models.UpdateProduct{
			GoogleProductCategory: models.GoogleProductCategory{
				Name: "Old Category",
			},
		}
		updatedProductWithReflect := reflect.ValueOf(&updatedProduct)

		// Nothing should be updated, but should not fail (original object is returned)
		result := SetFieldWithReflect("GoogleProductCategory.Name", updatedProductWithReflect, "NEW VALUE")
		resultAsType := result.Interface().(models.UpdateProduct)
		assert.Equal(t, resultAsType.GoogleProductCategory.Name, "NEW VALUE")
	})

	t.Run("Update on map (Update)", func(t *testing.T) {
		updatedProduct := models.UpdateProduct{
			Attributes: map[string]interface{}{
				"sample": "info",
			},
		}
		updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
		// Nothing should be updated, but should not fail (original object is returned)
		result := SetFieldWithReflect("Attributes.sample", updatedProductWithReflect, "NEW VALUE")
		resultAsType := result.Interface().(models.UpdateProduct)
		assert.Equal(t, resultAsType.Attributes["sample"], "NEW VALUE")
	})

	t.Run("Update on map (Create)", func(t *testing.T) {
		updatedProduct := models.UpdateProduct{
			Attributes: map[string]interface{}{
				"sample": "info",
			},
		}
		updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
		// Nothing should be updated, but should not fail (original object is returned)
		result := SetFieldWithReflect("Attributes.newField", updatedProductWithReflect, "NEW VALUE")
		resultAsType := result.Interface().(models.UpdateProduct)
		assert.Equal(t, resultAsType.Attributes["newField"], "NEW VALUE")
	})

	t.Run("Nested Map", func(t *testing.T) {

	})
	
	t.Run("Set Nest (Success)", func(t *testing.T) {
		updatedProduct := models.UpdateProduct{
			Variants: []models.Variant{
				{
					Title: "Variant 1", // Should stay the same
				},
				{
					Title: "Variant 2", // Should change
				},
			},
		}
		updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
		result := SetFieldWithReflect("Variants[1].Title", updatedProductWithReflect, "NEW VALUE")
		resultAsType := result.Interface().(models.UpdateProduct)
		assert.Equal(t, "Variant 1", resultAsType.Variants[0].Title)
		assert.Equal(t, "NEW VALUE", resultAsType.Variants[1].Title)
	})

	//t.Run("Set Nest (Array index does not exist)", func(t *testing.T) {
	//	updatedProduct := models.UpdateProduct{}
	//	updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
	//	t.SetFieldWithReflect()
	//})
	//
	//t.Run("Set Nest (Field on array object does not exist)", func(t *testing.T) {
	//	updatedProduct := models.UpdateProduct{}
	//	updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
	//	t.SetFieldWithReflect()
	//})
	//
	//t.Run("Root field", func(t *testing.T) {
	//	updatedProduct := models.UpdateProduct{}
	//	updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
	//	t.SetFieldWithReflect()
	//})

	//t.Run("Root field (Data type misaligned)", func(t *testing.T) {
	//	updatedProduct := models.UpdateProduct{}
	//	updatedProductWithReflect := reflect.ValueOf(&updatedProduct)
	//	t.SetFieldWithReflect()
	//})
}