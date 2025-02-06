package event

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
)

const (
	ProductCreated Name = "PRODUCT_CREATED"
)

func NewProductCreatedEvent(product *entity.Product) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: product.Id,
		Name:          ProductCreated,
		Data: map[string]string{
			"productId": product.Id,
		},
	}
}
