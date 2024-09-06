package repository

import (
	"bom-pedido-api/domain/entity/customer"
	"context"
)

type CustomerNotificationRepository interface {
	FindByCustomer(ctx context.Context, id string) (*customer.Notification, error)
	Upsert(ctx context.Context, notification *customer.Notification) error
}
