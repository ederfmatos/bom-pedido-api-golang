package repository

import (
	"bom-pedido-api/internal/application/repository"
	"bom-pedido-api/internal/domain/entity/customer"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerNotificationMongoRepository struct {
	collection *mongo.Collection
}

func NewCustomerNotificationMongoRepository(database *mongo.Database) repository.CustomerNotificationRepository {
	return &CustomerNotificationMongoRepository{collection: database.Collection("customer_notification_settings")}
}

func (repository *CustomerNotificationMongoRepository) FindByCustomer(ctx context.Context, id string) (*customer.Notification, error) {
	result := repository.collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	notification := &customer.Notification{}
	err := result.Decode(notification)
	if err != nil || notification.CustomerId == "" {
		return nil, err
	}
	return notification, nil
}

func (repository *CustomerNotificationMongoRepository) Upsert(ctx context.Context, notification *customer.Notification) error {
	update := bson.M{"$set": notification}
	updateOptions := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "_id", Value: notification.CustomerId}}
	_, err := repository.collection.UpdateOne(ctx, filter, update, updateOptions)
	return err
}
