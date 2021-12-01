package transform

import (
	"regexp"
	"strings"
	"switchboard-module-boilerplate/models"
)

const regex = `<.*?>`
const wordCountAttributeKey = "word_count"

// Transform returns processed flagged, the updated model, and error.
// Products that have already been processed will be returned as false.
func Transform(product models.Product) (bool, models.Product, error) {
	if _, ok := product.Attributes[wordCountAttributeKey]; ok {
		// Already exists
		return false, models.Product{}, nil
	}

	description := stripHtmlRegex(product.Description)

	// Count words
	if product.Attributes == nil {
		product.Attributes = map[string]interface{}{}
	}
	product.Attributes[wordCountAttributeKey] = CountWords(description)

	return true, product, nil
}

// This method uses a regular expresion to remove HTML tags.
func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(regex)
	return r.ReplaceAllString(s, "")
}


func CountWords(s string)  int {
	return len(strings.Fields(s))
}