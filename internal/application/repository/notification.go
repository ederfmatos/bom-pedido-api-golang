package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	Stream(ctx context.Context) <-chan *entity.Notification
	Delete(ctx context.Context, notification *entity.Notification)
	Update(ctx context.Context, notification *entity.Notification)
}
