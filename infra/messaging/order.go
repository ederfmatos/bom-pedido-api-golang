package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/order/await_approval_order"
	"bom-pedido-api/application/usecase/order/save_history"
	"context"
	"time"
)

func HandleOrderEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(event.OptionsForTopics("AWAIT_APPROVAL_ORDER", event.PixTransactionPaid), handleAwaitApprovalOrder(factory))
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
		),
		handleOrderStatusChanged(factory),
	)
}

func handleAwaitApprovalOrder(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := await_approval_order.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := await_approval_order.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}

func handleOrderStatusChanged(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := save_history.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		eventTime, err := time.Parse(time.RFC3339, message.Event.Data["at"])
		if err != nil {
			return err
		}
		input := save_history.Input{
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
