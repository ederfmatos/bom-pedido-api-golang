package event

import (
	"bom-pedido-api/domain/value_object"
)

var (
	PixPaymentCreated       = "PIX_PAYMENT_CREATED"
	PixPaymentRefunded      = "PIX_PAYMENT_REFUNDED"
	PaymentCallbackReceived = "PAYMENT_CALLBACK_RECEIVED"
)

func NewPixPaymentCreated(orderId, paymentId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixPaymentCreated,
		Data: map[string]string{
			"paymentId": paymentId,
			"orderId":   orderId,
		},
	}
}

func NewPixPaymentRefunded(orderId, paymentId string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixPaymentRefunded,
		Data: map[string]string{
			"paymentId": paymentId,
			"orderId":   orderId,
		},
	}
}

func NewPaymentCallbackReceived(gateway, orderId, eventName string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PaymentCallbackReceived,
		Data: map[string]string{
			"gateway":   gateway,
			"orderId":   orderId,
			"eventName": eventName,
		},
	}
}
