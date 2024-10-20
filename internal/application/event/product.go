package event

import (
	"bom-pedido-api/internal/domain/entity/product"
	"bom-pedido-api/internal/domain/value_object"
)

var (
	ProductCreated = "PRODUCT_CREATED"
)

func NewProductCreatedEvent(product *product.Product) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: product.Id,
		Name:          ProductCreated,
		Data: map[string]string{
			"productId": product.Id,
		},
	}
}
