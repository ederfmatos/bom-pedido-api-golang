package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type CustomerNotificationMemoryRepository struct {
	customers map[string]*entity.CustomerNotification
}

func NewCustomerNotificationMemoryRepository() *CustomerNotificationMemoryRepository {
	return &CustomerNotificationMemoryRepository{
		customers: make(map[string]*entity.CustomerNotification),
	}
}

func (r *CustomerNotificationMemoryRepository) FindByCustomerId(_ context.Context, id string) (*entity.CustomerNotification, error) {
	return r.customers[id], nil
}

func (r *CustomerNotificationMemoryRepository) Upsert(_ context.Context, notification *entity.CustomerNotification) error {
	r.customers[notification.CustomerId] = notification
	return nil
}
