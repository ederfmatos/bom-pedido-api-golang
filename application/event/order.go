package event

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/value_object"
)

var (
	OrderCreatedEventName  = "ORDER_CREATED"
	OrderApprovedEventName = "ORDER_APPROVED"
	OrderRejectedEventName = "ORDER_REJECTED"
)

func newOrderEvent(order *order.Order, name string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          name,
		Data: map[string]string{
			"orderId":    order.Id,
			"customerId": order.CustomerID,
		},
	}
}

func NewOrderCreatedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderCreatedEventName)
}

func NewOrderApprovedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderApprovedEventName)
}

func NewOrderRejectedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderRejectedEventName)
}
