package repository

import (
	"bom-pedido-api/internal/domain/entity/customer"
	"context"
)

type CustomerNotificationRepository interface {
	FindByCustomerId(ctx context.Context, id string) (*customer.Notification, error)
	Upsert(ctx context.Context, notification *customer.Notification) error
}
