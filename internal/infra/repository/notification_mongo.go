package repository

import (
	"bom-pedido-api/internal/domain/entity/notification"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type NotificationMongoRepository struct {
	collection *mongo.Collection
}

func NewNotificationMongoRepository(database *mongo.Database) *NotificationMongoRepository {
	return &NotificationMongoRepository{collection: database.ForCollection("notifications")}
}

func (r *NotificationMongoRepository) Create(ctx context.Context, notification *notification.Notification) error {
	return r.collection.InsertOne(ctx, notification)
}

func (r *NotificationMongoRepository) Stream(ctx context.Context) <-chan *notification.Notification {
	channel := make(chan *notification.Notification)
	go func() {
		stream, _ := r.collection.FetchStream(ctx)
		for id := range stream {
			aNotification, err := r.FindById(ctx, id)
			if err != nil || aNotification == nil {
				continue
			}
			channel <- aNotification
		}
	}()
	return channel
}

func (r *NotificationMongoRepository) FindById(ctx context.Context, id string) (*notification.Notification, error) {
	var aNotification notification.Notification
	err := r.collection.FindByID(ctx, id, &aNotification)
	if err != nil || aNotification.Id == "" {
		return nil, err
	}
	return &aNotification, nil
}

func (r *NotificationMongoRepository) Delete(ctx context.Context, notification *notification.Notification) {
	_ = r.collection.DeleteByID(ctx, notification.Id)
}

func (r *NotificationMongoRepository) Update(ctx context.Context, notification *notification.Notification) {
	_ = r.collection.UpdateByID(ctx, notification.Id, notification)
}
