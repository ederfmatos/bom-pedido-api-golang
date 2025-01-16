package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type CustomerNotificationRepository interface {
	FindByCustomerId(ctx context.Context, id string) (*entity.CustomerNotification, error)
	Upsert(ctx context.Context, notification *entity.CustomerNotification) error
}
