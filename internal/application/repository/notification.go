package repository

import (
	"bom-pedido-api/internal/domain/entity/notification"
	"context"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *notification.Notification) error
	Stream(ctx context.Context) <-chan *notification.Notification
	Delete(ctx context.Context, notification *notification.Notification)
	Update(ctx context.Context, notification *notification.Notification)
}
