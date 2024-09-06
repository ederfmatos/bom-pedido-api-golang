package messaging

import (
	"bom-pedido-api/application/event"
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/usecase/notification/notify_customer_order_status_changed"
	"bom-pedido-api/application/usecase/notification/send_notification"
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

	sendNotificationsUseCase := send_notification.New(factory)
	go sendNotificationsUseCase.Execute(context.Background())
}

func handleNotifyCustomerOrderStatusChanged(factory *factory.ApplicationFactory) event.HandlerFunc {
	useCase := notify_customer_order_status_changed.New(factory)
	return func(ctx context.Context, message *event.MessageEvent) error {
		input := notify_customer_order_status_changed.Input{
			OrderId: message.Event.Data["orderId"],
		}
		err := useCase.Execute(ctx, input)
		return message.AckIfNoError(ctx, err)
	}
}
