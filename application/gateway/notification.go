package gateway

import (
	"bom-pedido-api/domain/entity/notification"
	"context"
)

type NotificationGateway interface {
	Send(ctx context.Context, notification *notification.Notification) error
}
