package events

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/domain/entity/order"
)

var (
	OrderCreatedEventName = "ORDER_CREATED"
)

type OrderCreatedData struct {
	OrderId string `json:"orderId"`
}

func NewOrderCreatedEvent(order *order.Order) *event.Event {
	return &event.Event{
		Id:   order.Id,
		Name: OrderCreatedEventName,
		Data: OrderCreatedData{OrderId: order.Id},
	}
}
