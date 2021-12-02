package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"reflect"
	"strconv"
	"strings"
	"switchboard-module-boilerplate/env"
	"switchboard-module-boilerplate/models"
)

func UpdateProduct(product models.Product, updatedProduct models.Product, event models.TriggerEvent) error {
	diff, updateProductPayload := CreateUpdateProductPayload(product, updatedProduct)
	if !diff {
		// Two models are the same, no need to update
		return nil
	}

	// Marshal into payload
	productPayloadAsBytes, err := json.Marshal(updateProductPayload)
	if err != nil {
		return nil
	}

	url := fmt.Sprintf("%s/products/%s", env.ConvictionalAPIURL(), product.ID)
	if env.IsBuyer() {
		return errors.New("no endpoint exists yet")
	}

	return PublishToAPI(APIPublishConfig{
		Payload: productPayloadAsBytes,
		Method: "PATCH",
		URL: url,
		Headers: map[string]string{
			"Authorization": env.ConvictionalAPIKeyForLoad(),
			"Content-Type": "application/json",
		},
	}, event)
}

// CreateUpdateProductPayload compares the two objects, and creates an update.
func CreateUpdateProductPayload(product models.Product, updatedProduct models.Product) (bool, models.UpdateProduct) {
	// Get fields that have changed
	fieldsThatDiffer := pretty.Diff(product, updatedProduct)

	// Converted the updated product to reflect so we can retrieve values
	updatedProductAsReflect := reflect.ValueOf(&updatedProduct)

	// Set up the output
	result := models.UpdateProduct{}
	resultAsReflect := reflect.ValueOf(&result)
	fmt.Printf("fieldsThatDiffer :: %+v\n", fieldsThatDiffer)

	for _, fieldThatDiff := range fieldsThatDiffer{
		if fieldThatDiff == "" {
			continue
		}
		// Parse out the field path (Ex. "Variants[1].Title" starts as "Variants[1].Title: "Variant 2a" != "Variant 2"
		path := fieldThatDiff[:strings.Index(fieldThatDiff, ":")]
		fmt.Printf("path :: %+v\n", path)

		// Handle renaming issues
		path = strings.ReplaceAll(path, "Description", "BodyHTML")

		// Update set the field
		newValue := GetValueWithReflect(path, updatedProductAsReflect)
		if newValue == nil {
			// Could not find the field from the path
			continue
		}
		resultAsReflect = SetFieldWithReflect(path, resultAsReflect, newValue)
		fmt.Printf("resultAsReflect :: %+v\n", resultAsReflect)
	}

	// Convert the reflect object back to the struct
	result = resultAsReflect.Interface().(models.UpdateProduct)

	return false, result // TODO - Check if updateProduct is empty
}

// GetValueWithReflect Copied from and changed to return the value. There is an abstraction with finding the object but
// not obvious
func GetValueWithReflect(path string, obj reflect.Value) interface{} {
	if obj.Kind() == reflect.Ptr {
		// Object is a pointer, make reference to the actual value
		obj = obj.Elem()
	}
	// Parse current field
	if strings.Contains(path, ".") || strings.Contains(path, "[") {
		// Not at root (Either currently slice or child object). Find the next
		// field, and iterate.
		if (strings.Index(path, ".") < strings.Index(path, "[")) || !strings.Contains(path, "[") {
			// Confirm the current object is a struct
			// TODO

			// The next iteration is a child object. Get the field, so 0 --> . would be
			// the field name.
			fieldName := path[:strings.Index(path, ".")]

			// Confirm field exists and it's an object
			if obj.FieldByName(fieldName).Kind() == reflect.Invalid {
				// Field does not exist. Just return, and do not make changes.
				return obj
			} else if obj.FieldByName(fieldName).Kind() != reflect.Struct && obj.FieldByName(fieldName).Kind() != reflect.Map  {
				// The field on the object is not a struct or map, therefore we won't be able to set child fields.
				return nil
			}

			// Call the next iteration down the path
			nextPath := path[strings.Index(path, ".") + 1:]
			if nextPath == "" {
				// The path provide was invalid. Ex. MadeUpField.TheMoon.
				return nil
			}

			// Recursively move to child element
			return GetValueWithReflect(nextPath, obj.FieldByName(fieldName))
		} else {
			// The next iteration is an array. Get the field name
			fieldName := path[:strings.Index(path, "[")]

			// Confirm it is an array
			if obj.FieldByName(fieldName).Kind() == reflect.Invalid {
				// Field does not exist. Just return, and do not make changes.
				return nil
			} else if obj.FieldByName(fieldName).Kind() != reflect.Slice {
				// The field on the object is not a slice
				return nil
			}

			// Get the index
			indexAsStr := path[strings.Index(path, "[") + 1: strings.Index(path, "]")] // TODO - Add better validation
			index, err := strconv.Atoi(indexAsStr)
			if err != nil {
				// Provided invalid index
				return nil
			}
			// TODO - Confirm not out of range

			// Get next path
			nextPath := path[strings.Index(path, "]") + 2:]
			if nextPath == "" {
				// Provide an invalid path string. Ex. myArray[2].
				return nil
			}

			// Recurse
			return GetValueWithReflect(nextPath, obj.FieldByName(fieldName).Index(index))
		}
	}

	if obj.Kind() == reflect.Map {
		return obj.MapIndex(reflect.ValueOf(path)).Interface()
	}

	// At this point, just setting the field. First confirm field exists.
	if obj.FieldByName(path).Kind() == reflect.Invalid {
		// Field does not exists
		return nil
	}

	return obj.FieldByName(path).Interface()
}

func SetFieldWithReflect(path string, obj reflect.Value, newValue interface{}) reflect.Value {
	if obj.Kind() == reflect.Ptr {
		// Object is a pointer, make reference to the actual value
		obj = obj.Elem()
	}
	fmt.Printf("obj.Type() :: %+v\n", obj.Type())
	fmt.Printf("obj.Kind() :: %+v\n", obj.Kind())
	fmt.Printf("reflect.TypeOf(obj) :: %+v\n", reflect.TypeOf(obj))

	// Parse current field
	if strings.Contains(path, ".") || strings.Contains(path, "[") {
		// Not at root (Either currently slice or child object). Find the next
		// field, and iterate.
		if (strings.Index(path, ".") < strings.Index(path, "[")) || !strings.Contains(path, "[") {
			// Confirm the current object is a struct
			// TODO

			// The next iteration is a child object. Get the field, so 0 --> . would be
			// the field name.
			fieldName := path[:strings.Index(path, ".")]

			// Confirm field exists and it's an object
			if obj.FieldByName(fieldName).Kind() == reflect.Invalid {
				// Field does not exist. Just return, and do not make changes.
				return obj
			} else if obj.FieldByName(fieldName).Kind() != reflect.Struct && obj.FieldByName(fieldName).Kind() != reflect.Map  {
				// The field on the object is not a struct or map, therefore we won't be able to set child fields.
				return obj
			}

			// Call the next iteration down the path
			nextPath := path[strings.Index(path, ".") + 1:]
			if nextPath == "" {
				// The path provide was invalid. Ex. MadeUpField.TheMoon.
				return obj
			}


			fmt.Printf("obj.FieldByName(fieldName).Kind() :: %+v\n", obj.FieldByName(fieldName).Kind())
			fmt.Printf("next path :: %+v\n", nextPath)
			// Recursively move to child element
			resultFromChild := SetFieldWithReflect(nextPath, obj.FieldByName(fieldName), newValue)

			// Set outcome of child
			obj.FieldByName(fieldName).Set(resultFromChild)
			return obj
		} else {
			// The next iteration is an array. Get the field name
			fieldName := path[:strings.Index(path, "[")]

			// Confirm it is an array
			if obj.FieldByName(fieldName).Kind() == reflect.Invalid {
				// Field does not exist. Just return, and do not make changes.
				return obj
			} else if obj.FieldByName(fieldName).Kind() != reflect.Slice {
				// The field on the object is not a slice
				return obj
			}

			// Get the index
			indexAsStr := path[strings.Index(path, "[") + 1: strings.Index(path, "]")] // TODO - Add better validation
			index, err := strconv.Atoi(indexAsStr)
			if err != nil {
				// Provided invalid index
				return obj
			}
			// TODO - Confirm not out of range

			// Get next path
			nextPath := path[strings.Index(path, "]") + 2:]
			if nextPath == "" {
				// Provide an invalid path string. Ex. myArray[2].
				return obj
			}

			fmt.Printf("obj.FieldByName(fieldName).Len() :: %d\n", obj.FieldByName(fieldName).Len())
			fmt.Printf("index :: %d\n", index)

			if obj.FieldByName(fieldName).Len() <= index {
				// If the index is after the end of the slice. Therefore this is appending a new item.
				currentSlice := obj.FieldByName(fieldName)
				sliceType := reflect.TypeOf(currentSlice.Interface()).Elem() // Elem() grabs the slice element type compared to slice
				fmt.Printf("sliceType :: %s\n", sliceType)
				// We first need to create a new element of that type. We will pass through (reflect.Value{}) but we
				// need to know the type of struct. It should be an empty one.
				newElement := reflect.New(sliceType)
				fmt.Printf("new Element :: %s\n", newElement.Interface())
				resultFromChild := SetFieldWithReflect(nextPath, newElement, newValue)

				// Set
				fmt.Printf("resultFromChild :: %+v\n", resultFromChild)
				fmt.Printf("resultFromChild.Kind() :: %+v\n", resultFromChild.Kind())

				newSlice := reflect.Append(currentSlice, resultFromChild)
				obj.FieldByName(fieldName).Set(newSlice)
				return obj
			}

			// Recurse
			childElement := obj.FieldByName(fieldName).Index(index)
			resultFromChild := SetFieldWithReflect(nextPath, childElement, newValue)

			// Set the outcome
			obj.FieldByName(fieldName).Index(index).Set(resultFromChild)
			return obj
		}
	}

	if obj.FieldByName(path).Kind() == reflect.Interface {
		// This would be setting a field from being nil
		return obj
	}

	if obj.Kind() == reflect.Map {
		newValueAsReflectValue := reflect.ValueOf(newValue)
		obj.SetMapIndex(reflect.ValueOf(path), newValueAsReflectValue)
		return obj
	}

	// At this point, just setting the field. First confirm field exists.
	if obj.FieldByName(path).Kind() == reflect.Invalid {
		// Field does not exists
		return obj
	}

	// We need to confirm types.
	fmt.Printf("path :: %s\n", path)
	fmt.Printf("obj.FieldByName(path).Kind().String() :: %s\n", obj.FieldByName(path).Kind().String())
	if obj.FieldByName(path).Kind().String() != fmt.Sprintf("%T", newValue) {
		return obj
	}

	// Set value
	newValueAsReflectValue := reflect.ValueOf(newValue)
	obj.FieldByName(path).Set(newValueAsReflectValue)

	return obj
}