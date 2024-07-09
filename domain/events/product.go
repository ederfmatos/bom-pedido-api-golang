package events

import (
	"bom-pedido-api/domain/entity/product"
)

var (
	ProductCreatedEventName = "PRODUCT_CREATED"
)

type ProductCreatedData struct {
	ProductId string `json:"productId"`
}

func NewProductCreatedEvent(product *product.Product) *Event {
	return &Event{
		Id:   product.Id,
		Name: ProductCreatedEventName,
		Data: ProductCreatedData{ProductId: product.Id},
	}
}
