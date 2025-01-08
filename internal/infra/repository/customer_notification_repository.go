package repository

import (
	"bom-pedido-api/internal/domain/entity/customer"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type CustomerNotificationMongoRepository struct {
	collection *mongo.Collection
}

func NewCustomerNotificationMongoRepository(database *mongo.Database) *CustomerNotificationMongoRepository {
	return &CustomerNotificationMongoRepository{collection: database.ForCollection("customer_notifications")}
}

func (r *CustomerNotificationMongoRepository) FindByCustomerId(ctx context.Context, id string) (*customer.Notification, error) {
	var notification customer.Notification
	err := r.collection.FindByID(ctx, id, &notification)
	if err != nil || notification.CustomerId == "" {
		return nil, err
	}
	return &notification, nil
}

func (r *CustomerNotificationMongoRepository) Upsert(ctx context.Context, notification *customer.Notification) error {
	return r.collection.Upsert(ctx, notification.CustomerId, notification)
}
