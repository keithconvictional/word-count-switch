package load

import (
	"errors"
	"fmt"
	"switchboard-module-boilerplate/env"
	"switchboard-module-boilerplate/models"
)

const (
	LoadMethodConvictionalAPI = "convictional_api"
)

func Single(product models.Product, event models.TriggerEvent) error {
	fmt.Printf("product :: %+v\n", product)
	switch env.LoadMethod() {
	case LoadMethodConvictionalAPI:
		return UpdateProduct(product, event)
	default:
		return errors.New("invalid load method")
	}
	return nil
}

func Multiple( products []models.Product, event models.TriggerEvent) error {
	return nil
}
