package events

import (
	"bom-pedido-api/domain/entity/order"
)

var (
	OrderCreatedEventName = "ORDER_CREATED"
)

type OrderCreatedEventData struct {
	OrderId    string `json:"orderId"`
	CustomerId string `json:"customerId"`
}

func NewOrderCreatedEvent(order *order.Order) *Event {
	return &Event{
		Id:   order.Id,
		Name: OrderCreatedEventName,
		Data: OrderCreatedEventData{
			OrderId:    order.Id,
			CustomerId: order.CustomerID,
		},
	}
}
