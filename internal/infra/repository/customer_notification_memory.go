package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/customer"
	"context"
)

type CustomerNotificationMemoryRepository struct {
	customers map[string]*customer.Notification
}

func NewCustomerNotificationMemoryRepository() repository.CustomerNotificationRepository {
	return &CustomerNotificationMemoryRepository{
		customers: make(map[string]*customer.Notification),
	}
}

func (repository *CustomerNotificationMemoryRepository) FindByCustomer(_ context.Context, id string) (*customer.Notification, error) {
	return repository.customers[id], nil
}

func (repository *CustomerNotificationMemoryRepository) Upsert(_ context.Context, notification *customer.Notification) error {
	repository.customers[notification.CustomerId] = notification
	return nil
}
