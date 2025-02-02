package messaging

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/internal/application/usecase/notification"
	"context"
)

type NotificationHandler struct {
	notifyCustomerOrderStatusChangedUseCase *notification.NotifyCustomerOrderStatusChangedUseCase
}

func NewNotificationHandler(
	notifyCustomerOrderStatusChangedUseCase *notification.NotifyCustomerOrderStatusChangedUseCase,
	sendNotificationUseCase *notification.SendNotificationUseCase,
) *NotificationHandler {
	go sendNotificationUseCase.Execute(context.Background())
	return &NotificationHandler{
		notifyCustomerOrderStatusChangedUseCase: notifyCustomerOrderStatusChangedUseCase,
	}
}

func (h NotificationHandler) NotifyCustomerOrderStatusChanged(ctx context.Context, message *event.MessageEvent) error {
	input := notification.NotifyCustomerOrderStatusChangedInput{
		OrderId: message.Event.Data["orderId"],
	}
	err := h.notifyCustomerOrderStatusChangedUseCase.Execute(ctx, input)
	return message.AckIfNoError(ctx, err)
}
