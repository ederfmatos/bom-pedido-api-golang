package event

import (
	"bom-pedido-api/domain/entity/product"
	"bom-pedido-api/domain/value_object"
)

var (
	ProductCreatedEventName = "PRODUCT_CREATED"
)

type ProductCreatedData struct {
	ProductId string `json:"productId"`
}

func NewProductCreatedEvent(product *product.Product) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: product.Id,
		Name:          ProductCreatedEventName,
		Data: map[string]string{
			"productId": product.Id,
		},
	}
}
