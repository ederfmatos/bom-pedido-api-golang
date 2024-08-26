package event

import (
	"bom-pedido-api/domain/value_object"
)

var (
	refundTransactionEvent  = "REFUND_TRANSACTION"
	pixTransactionCreated   = "PIX_TRANSACTION_CREATED"
	pixTransactionPaid      = "PIX_TRANSACTION_PAID"
	pixTransactionRefunded  = "PIX_TRANSACTION_REFUNDED"
	paymentCallbackReceived = "PAYMENT_CALLBACK_RECEIVED"
)

func NewRefundTransactionEvent(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          refundTransactionEvent,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewPixTransactionCreated(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          pixTransactionCreated,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewTransactionPaid(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          pixTransactionPaid,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewPixTransactionRefunded(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          pixTransactionRefunded,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewPaymentCallbackReceived(gateway, orderId, eventName string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          paymentCallbackReceived,
		Data: map[string]string{
			"gateway":   gateway,
			"orderId":   orderId,
			"eventName": eventName,
		},
	}
}
