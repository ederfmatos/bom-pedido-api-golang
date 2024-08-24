package event

import (
	"bom-pedido-api/domain/value_object"
)

var (
	refundTransactionEvent = "REFUND_TRANSACTION"
	pixTransactionCreated  = "PIX_TRANSACTION_CREATED"
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
