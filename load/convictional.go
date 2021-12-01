package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"switchboard-module-boilerplate/env"
	"switchboard-module-boilerplate/models"
)

func UpdateProduct(product models.Product, event models.TriggerEvent) error {
	productPayload, err := json.Marshal(product)
	if err != nil {
		return nil
	}

	url := fmt.Sprintf("https://api.convictional.com/products/%s", product.ID)
	if env.IsBuyer() {
		return errors.New("no endpoint exists yet")
	}

	return PublishToAPI(APIPublishConfig{
		Payload: productPayload,
		Method: "PATCH",
		URL: url,
		Headers: map[string]string{
			"Authorization": env.ConvictionalAPIKeyForLoad(),
			"Content-Type": "application/json",
		},
	}, event)
}