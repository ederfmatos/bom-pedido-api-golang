package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/usecase/transaction"
	"context"
)

type TransactionHandler struct {
	payPixTransactionUseCase    *transaction.PayPixTransactionUseCase
	createPixTransactionUseCase *transaction.CreatePixTransactionUseCase
	refundPixTransactionUseCase *transaction.RefundPixTransactionUseCase
	cancelPixTransactionUseCase *transaction.CancelPixTransactionUseCase
}

func NewTransactionHandler(
	payPixTransactionUseCase *transaction.PayPixTransactionUseCase,
	createPixTransactionUseCase *transaction.CreatePixTransactionUseCase,
	refundPixTransactionUseCase *transaction.RefundPixTransactionUseCase,
	cancelPixTransactionUseCase *transaction.CancelPixTransactionUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		payPixTransactionUseCase:    payPixTransactionUseCase,
		createPixTransactionUseCase: createPixTransactionUseCase,
		refundPixTransactionUseCase: refundPixTransactionUseCase,
		cancelPixTransactionUseCase: cancelPixTransactionUseCase,
	}
}

func (h TransactionHandler) HandlePayPixTransaction(ctx context.Context, message *event.MessageEvent) error {
	input := transaction.PayPixTransactionInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.payPixTransactionUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h TransactionHandler) HandleCreatePixTransaction(ctx context.Context, message *event.MessageEvent) error {
	input := transaction.CreatePixTransactionInput{
		OrderId:        message.Event.Data["orderId"],
		PaymentId:      message.Event.Data["paymentId"],
		PaymentGateway: message.Event.Data["paymentGateway"],
	}
	err := h.createPixTransactionUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h TransactionHandler) HandleRefundPixTransaction(ctx context.Context, message *event.MessageEvent) error {
	input := transaction.RefundPixTransactionInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.refundPixTransactionUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h TransactionHandler) HandleCancelPixTransaction(ctx context.Context, message *event.MessageEvent) error {
	input := transaction.CancelPixTransactionInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.cancelPixTransactionUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}
