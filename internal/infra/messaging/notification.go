package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/usecase/notification"
	"context"
)

func HandleNotificationEvents(factory *factory.ApplicationFactory) {
	factory.EventHandler.Consume(
		event.OptionsForTopics(
			"NOTIFY_CUSTOMER_ORDER_STATUS_CHANGED",
			event.OrderAwaitingApproval,
			event.OrderApproved,
			event.OrderInProgress,
			event.OrderRejected,
			event.OrderCancelled,
			event.OrderDelivering,
			event.OrderAwaitingWithdraw,
			event.OrderAwaitingDelivery,
			event.OrderFinished,
		),
		handleNotifyCustomerOrderStatusChanged(factory),
	)

	sendNotificationsUseCase := notification.NewSendNotification(factory)
	go sendNotificationsUseCase.Execute(context.Background())
}

func handleNotifyCustomerOrderStatusChanged(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := notification.NewNotifyCustomerOrderStatusChanged(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := notification.NotifyCustomerOrderStatusChangedInput{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
