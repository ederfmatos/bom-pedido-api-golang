package event

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/internal/domain/value_object"
	"time"
)

const (
	OrderCreated          Name = "ORDER_CREATED"
	OrderAwaitingApproval Name = "ORDER_AWAITING_APPROVAL"
	OrderApproved         Name = "ORDER_APPROVED"
	OrderFinished         Name = "ORDER_FINISHED"
	OrderRejected         Name = "ORDER_REJECTED"
	OrderInProgress       Name = "ORDER_IN_PROGRESS"
	OrderDelivering       Name = "ORDER_DELIVERING"
	OrderAwaitingWithdraw Name = "ORDER_AWAITING_WITHDRAW"
	OrderAwaitingDelivery Name = "ORDER_AWAITING_DELIVERY"
	OrderCancelled        Name = "ORDER_CANCELLED"
	OrderPaymentFailed    Name = "ORDER_PAYMENT_FAILED"
)

func NewOrderCreatedEvent(order *entity.Order) *Event {
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

func NewOrderApprovedEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderInProgressEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderDeliveringEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderAwaitingWithdrawEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderAwaitingDeliveryEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderRejectedEvent(order *entity.Order, by string, at time.Time, reason string) *Event {
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

func NewOrderCancelledEvent(order *entity.Order, by string, at time.Time, reason string) *Event {
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

func NewOrderFinishedEvent(order *entity.Order, by string, at time.Time) *Event {
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

func NewOrderAwaitingApprovalEvent(order *entity.Order, at time.Time) *Event {
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

func NewOrderPaymentFailedEvent(order *entity.Order, at time.Time) *Event {
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
