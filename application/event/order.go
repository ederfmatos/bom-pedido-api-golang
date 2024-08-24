package event

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/value_object"
)

var (
	OrderCreatedEventName          = "ORDER_CREATED"
	OrderApprovedEventName         = "ORDER_APPROVED"
	OrderFinishedEventName         = "ORDER_FINISHED"
	OrderRejectedEventName         = "ORDER_REJECTED"
	OrderInProgressEventName       = "ORDER_IN_PROGRESS"
	OrderDeliveringEventName       = "ORDER_DELIVERING"
	OrderAwaitingWithdrawEventName = "ORDER_AWAITING_WITHDRAW"
	OrderAwaitingDeliveryEventName = "ORDER_AWAITING_DELIVERY"
	OrderCancelledEventName        = "ORDER_CANCELLED"
)

func newOrderEvent(order *order.Order, name string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          name,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
		},
	}
}

func NewOrderCreatedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderCreatedEventName)
}

func NewOrderApprovedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderApprovedEventName)
}

func NewOrderInProgressEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderInProgressEventName)
}

func NewOrderDeliveringEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderDeliveringEventName)
}

func NewOrderAwaitingWithdrawEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderAwaitingWithdrawEventName)
}

func NewOrderAwaitingDeliveryEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderAwaitingDeliveryEventName)
}

func NewOrderRejectedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderRejectedEventName)
}

func NewOrderCancelledEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderCancelledEventName)
}

func NewOrderFinishedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderFinishedEventName)
}
