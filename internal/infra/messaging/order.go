package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/order"
	"context"
	"time"
)

func HandleOrderEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("AWAIT_APPROVAL_ORDER", event.PixTransactionPaid), handleAwaitApprovalOrder(factory))
	factory.EventHandler.Consume(event.OptionsForTopics("ORDER_PAYMENT_FAILED", event.PixPaymentCancelled), handleOrderPaymentFailed(factory))
	factory.EventHandler.Consume(
		event.OptionsForTopics(
			"SAVE_ORDER_STATUS_HISTORY",
			event.OrderAwaitingApproval,
			event.OrderApproved,
			event.OrderAwaitingWithdraw,
			event.OrderAwaitingDelivery,
			event.OrderFinished,
			event.OrderInProgress,
			event.OrderDelivering,
			event.OrderCancelled,
			event.OrderRejected,
			event.OrderPaymentFailed,
		),
		handleOrderStatusChanged(factory),
	)
}

func handleAwaitApprovalOrder(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := order.NewAwaitApprovalOrderUseCase(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := order.AwaitApprovalOrderUseCaseInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleOrderPaymentFailed(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := order.NewFailOrderPayment(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := order.FailOrderPaymentInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleOrderStatusChanged(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := order.NewSaveOrderHistory(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
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
		err = useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
