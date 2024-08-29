package event

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/value_object"
)

var (
	OrderCreated          = "ORDER_CREATED"
	OrderAwaitingApproval = "ORDER_AWAITING_APPROVAL"
	OrderApproved         = "ORDER_APPROVED"
	OrderFinished         = "ORDER_FINISHED"
	OrderRejected         = "ORDER_REJECTED"
	OrderInProgress       = "ORDER_IN_PROGRESS"
	OrderDelivering       = "ORDER_DELIVERING"
	OrderAwaitingWithdraw = "ORDER_AWAITING_WITHDRAW"
	OrderAwaitingDelivery = "ORDER_AWAITING_DELIVERY"
	OrderCancelled        = "ORDER_CANCELLED"
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
	return newOrderEvent(order, OrderCreated)
}

func NewOrderApprovedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderApproved)
}

func NewOrderInProgressEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderInProgress)
}

func NewOrderDeliveringEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderDelivering)
}

func NewOrderAwaitingWithdrawEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderAwaitingWithdraw)
}

func NewOrderAwaitingDeliveryEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderAwaitingDelivery)
}

func NewOrderRejectedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderRejected)
}

func NewOrderCancelledEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderCancelled)
}

func NewOrderFinishedEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderFinished)
}

func NewOrderAwaitingApprovalEvent(order *order.Order) *Event {
	return newOrderEvent(order, OrderAwaitingApproval)
}
