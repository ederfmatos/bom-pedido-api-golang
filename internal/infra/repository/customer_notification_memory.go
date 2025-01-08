package repository

import (
	"bom-pedido-api/internal/domain/entity/customer"
	"context"
)

type CustomerNotificationMemoryRepository struct {
	customers map[string]*customer.Notification
}

func NewCustomerNotificationMemoryRepository() *CustomerNotificationMemoryRepository {
	return &CustomerNotificationMemoryRepository{
		customers: make(map[string]*customer.Notification),
	}
}

func (r *CustomerNotificationMemoryRepository) FindByCustomerId(_ context.Context, id string) (*customer.Notification, error) {
	return r.customers[id], nil
}

func (r *CustomerNotificationMemoryRepository) Upsert(_ context.Context, notification *customer.Notification) error {
	r.customers[notification.CustomerId] = notification
	return nil
}
