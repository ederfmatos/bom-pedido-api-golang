package event

import (
	"bom-pedido-api/domain/entity/order"
	"bom-pedido-api/domain/value_object"
	"time"
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
	OrderPaymentFailed    = "ORDER_PAYMENT_FAILED"
)

func NewOrderCreatedEvent(order *order.Order) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderCreated,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
		},
	}
}

func NewOrderApprovedEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderApproved,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderInProgressEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderInProgress,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderDeliveringEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderDelivering,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderAwaitingWithdrawEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderAwaitingWithdraw,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderAwaitingDeliveryEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderAwaitingDelivery,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderRejectedEvent(order *order.Order, by string, at time.Time, reason string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderRejected,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"reason":        reason,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderCancelledEvent(order *order.Order, by string, at time.Time, reason string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderCancelled,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"reason":        reason,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderFinishedEvent(order *order.Order, by string, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderFinished,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            by,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderAwaitingApprovalEvent(order *order.Order, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderAwaitingApproval,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            order.CustomerID,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}

func NewOrderPaymentFailedEvent(order *order.Order, at time.Time) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: order.Id,
		Name:          OrderPaymentFailed,
		Data: map[string]string{
			"orderId":       order.Id,
			"customerId":    order.CustomerID,
			"paymentMethod": order.PaymentMethod.String(),
			"by":            order.CustomerID,
			"at":            at.Format(time.RFC3339),
			"status":        order.GetStatus(),
		},
	}
}
