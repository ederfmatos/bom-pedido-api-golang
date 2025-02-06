package event

import (
	"bom-pedido-api/internal/domain/value_object"
)

const (
	PixTransactionCreated   Name = "PIX_TRANSACTION_CREATED"
	PixTransactionPaid      Name = "PIX_TRANSACTION_PAID"
	PixTransactionRefunded  Name = "PIX_TRANSACTION_REFUNDED"
	PixTransactionCancelled Name = "PIX_TRANSACTION_CANCELLED"
)

func NewPixTransactionCreated(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixTransactionCreated,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewPixTransactionPaid(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixTransactionPaid,
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
		Name:          PixTransactionRefunded,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}

func NewPixTransactionCancelled(orderId, transactionId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixTransactionCancelled,
		Data: map[string]string{
			"transactionId": transactionId,
			"orderId":       orderId,
		},
	}
}
