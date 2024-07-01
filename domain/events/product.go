package events

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/domain/entity"
)

var (
	ProductCreatedEventName = "PRODUCT_CREATED"
)

type ProductCreatedData struct {
	ProductId string `json:"productId"`
}

func NewProductCreatedEvent(product *entity.Product) *event.Event {
	return &event.Event{
		Id:   product.Id,
		Name: ProductCreatedEventName,
		Data: ProductCreatedData{ProductId: product.Id},
	}
}
