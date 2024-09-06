package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/notification"
	"bom-pedido-api/infra/event"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationMongoRepository struct {
	collection *mongo.Collection
	stream     event.Stream
}

func NewNotificationMongoRepository(database *mongo.Database) repository.NotificationRepository {
	collection := database.Collection("notifications")
	return &NotificationMongoRepository{
		collection: collection,
		stream:     event.NewMongoStream(collection),
	}
}

func (repository *NotificationMongoRepository) Create(ctx context.Context, notification *notification.Notification) error {
	_, err := repository.collection.InsertOne(ctx, notification)
	return err
}

func (repository *NotificationMongoRepository) Stream() <-chan *notification.Notification {
	channel := make(chan *notification.Notification)
	stream, _ := repository.stream.FetchStream()
	go func() {
		for id := range stream {
			aNotification, err := repository.FindById(context.Background(), id)
			if err != nil || aNotification == nil {
				continue
			}
			channel <- aNotification
		}
	}()
	return channel
}

func (repository *NotificationMongoRepository) FindById(ctx context.Context, id string) (*notification.Notification, error) {
	result := repository.collection.FindOne(ctx, bson.M{"_id": id})
	if err := result.Err(); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	aNotification := &notification.Notification{}
	err := result.Decode(aNotification)
	if err != nil {
		return nil, err
	}
	if aNotification.Id == "" {
		return nil, nil
	}
	return aNotification, nil
}

func (repository *NotificationMongoRepository) Delete(ctx context.Context, notification *notification.Notification) {
	_, _ = repository.collection.DeleteOne(ctx, bson.M{"_id": notification.Id})
}

func (repository *NotificationMongoRepository) Update(ctx context.Context, notification *notification.Notification) {
	update := bson.M{"$set": notification}
	updateOptions := options.Update()
	filter := bson.D{{"_id", notification.Id}}
	_, _ = repository.collection.UpdateOne(ctx, filter, update, updateOptions)
}
