package events

import (
	"bom-pedido-api/domain/entity/order"
)

var (
	OrderCreatedEventName  = "ORDER_CREATED"
	OrderApprovedEventName = "ORDER_APPROVED"
)

type OrderEventData struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
}

func newOrderEvent(order *order.Order, name string) *Event {
	return &Event{
		Id:   order.Id,
		Name: name,
		Data: OrderEventData{OrderId: order.Id, CustomerId: order.CustomerID},
	}
}

func NewOrderCreatedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderCreatedEventName)
}

func NewOrderApprovedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderApprovedEventName)
}
