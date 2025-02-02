package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/usecase/payment"
	"context"
)

type PaymentHandler struct {
	checkPixPaymentFailedUseCase *payment.CheckPixPaymentFailedUseCase
	createPixPaymentUseCase      *payment.CreatePixPaymentUseCase
	refundPixPaymentUseCase      *payment.RefundPixPaymentUseCase
}

func NewPaymentHandler(
	checkPixPaymentFailedUseCase *payment.CheckPixPaymentFailedUseCase,
	createPixPaymentUseCase *payment.CreatePixPaymentUseCase,
	refundPixPaymentUseCase *payment.RefundPixPaymentUseCase,
) *PaymentHandler {
	return &PaymentHandler{
		checkPixPaymentFailedUseCase: checkPixPaymentFailedUseCase,
		createPixPaymentUseCase:      createPixPaymentUseCase,
		refundPixPaymentUseCase:      refundPixPaymentUseCase,
	}
}

func (h PaymentHandler) CreatePixPayment(ctx context.Context, message *event.MessageEvent) error {
	input := payment.CreatePixPaymentInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.createPixPaymentUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h PaymentHandler) RefundPixPayment(ctx context.Context, message *event.MessageEvent) error {
	input := payment.RefundPixPaymentInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.refundPixPaymentUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h PaymentHandler) CheckPixPaymentFailed(ctx context.Context, message *event.MessageEvent) error {
	input := payment.CheckPixPaymentFailedInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.checkPixPaymentFailedUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}
