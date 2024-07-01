package events

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/domain/entity"
)

var (
	OrderCreatedEventName = "ORDER_CREATED"
)

type OrderCreatedData struct {
	OrderId string `json:"orderId"`
}

func NewOrderCreatedEvent(order *entity.Order) *event.Event {
	return &event.Event{
		Id:   order.Id,
		Name: OrderCreatedEventName,
		Data: OrderCreatedData{OrderId: order.Id},
	}
}
