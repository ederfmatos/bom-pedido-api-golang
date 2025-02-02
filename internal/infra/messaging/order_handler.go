package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/usecase/order"
	"context"
	"time"
)

type OrderHandler struct {
	saveOrderHistoryUseCase   *order.SaveOrderHistoryUseCase
	awaitApprovalOrderUseCase *order.AwaitApprovalOrderUseCase
	failOrderPaymentUseCase   *order.FailOrderPaymentUseCase
}

func NewOrderHandler(
	saveOrderHistoryUseCase *order.SaveOrderHistoryUseCase,
	awaitApprovalOrderUseCase *order.AwaitApprovalOrderUseCase,
	failOrderPaymentUseCase *order.FailOrderPaymentUseCase,
) *OrderHandler {
	return &OrderHandler{
		saveOrderHistoryUseCase:   saveOrderHistoryUseCase,
		awaitApprovalOrderUseCase: awaitApprovalOrderUseCase,
		failOrderPaymentUseCase:   failOrderPaymentUseCase,
	}
}

func (h OrderHandler) HandleOrderPaymentFailed(ctx context.Context, message *event.MessageEvent) error {
	input := order.FailOrderPaymentInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.failOrderPaymentUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h OrderHandler) HandleOrderStatusChanged(ctx context.Context, message *event.MessageEvent) error {
	eventTime, err := time.Parse(time.RFC3339, message.Event.Data["at"])
	if err != nil {
		return err
	}
	input := order.SaveOrderHistoryInput{
		Time:      eventTime,
		ChangedBy: message.Event.Data["by"],
		OrderId:   message.Event.Data["orderId"],
		Status:    message.Event.Data["status"],
		Data:      message.Event.Data["reason"],
	}
	err = h.saveOrderHistoryUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}

func (h OrderHandler) HandleAwaitApprovalOrder(ctx context.Context, message *event.MessageEvent) error {
	input := order.AwaitApprovalOrderUseCaseInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.awaitApprovalOrderUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}
