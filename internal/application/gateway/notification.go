package gateway

import (
	"bom-pedido-api/internal/domain/entity/notification"
	"context"
)

type NotificationGateway interface {
	Send(ctx context.Context, notification *notification.Notification) error
}
