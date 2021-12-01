package extract

import (
	"errors"
	"switchboard-module-boilerplate/env"
	"switchboard-module-boilerplate/models"
)

func Multiple(event models.TriggerEvent) ([]models.Product, error) {
	switch env.ExtractMethod() {
	default:
		return []models.Product{}, errors.New("invalid extract method")
	}
	return []models.Product{}, nil
}
