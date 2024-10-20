package event

import (
	"bom-pedido-api/internal/domain/value_object"
)

var (
	CheckPixPaymentFailed   = "CHECK_PIX_PAYMENT_FAILED"
	PixPaymentCreated       = "PIX_PAYMENT_CREATED"
	PixPaymentRefunded      = "PIX_PAYMENT_REFUNDED"
	PixPaymentCancelled     = "PIX_PAYMENT_CANCELLED"
	PaymentCallbackReceived = "PAYMENT_CALLBACK_RECEIVED"
)

func NewPixPaymentCancelled(orderId, paymentId, paymentGateway string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixPaymentCancelled,
		Data: map[string]string{
			"paymentId":      paymentId,
			"orderId":        orderId,
			"paymentGateway": paymentGateway,
		},
	}
}

func NewPixPaymentCreated(orderId, paymentId, paymentGateway string) *Event {
	return &Event{
		Id:            value_object.NewID(),
		CorrelationId: orderId,
		Name:          PixPaymentCreated,
		Data: map[string]string{
			"paymentId":      paymentId,
			"orderId":        orderId,
			"paymentGateway": paymentGateway,
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
