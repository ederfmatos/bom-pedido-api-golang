package notification

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity/notification"
	"bom-pedido-api/internal/infra/telemetry"
	"context"
	"firebase.google.com/go/v4/messaging"
)

type FirebaseNotificationGateway struct {
	fcmClient *messaging.Client
}

func NewFirebaseNotificationGateway(fcmClient *messaging.Client) gateway.NotificationGateway {
	return &FirebaseNotificationGateway{fcmClient: fcmClient}
}

func (f *FirebaseNotificationGateway) Send(ctx context.Context, notification *notification.Notification) error {
	ctx, span := telemetry.StartSpan(ctx, "FirebaseNotificationGateway.Send")
	defer span.End()
	_, err := f.fcmClient.Send(ctx, &messaging.Message{
		Notification: &messaging.Notification{
			Title: notification.Title,
			Body:  notification.Body,
		},
		Data:  notification.Data,
		Token: notification.Recipient,
	})
	return err
}
