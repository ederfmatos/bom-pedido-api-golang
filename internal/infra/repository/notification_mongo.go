package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/mongo"
	"context"
)

type NotificationMongoRepository struct {
	collection mongo.Collection
}

func NewNotificationMongoRepository(database *mongo.Database) *NotificationMongoRepository {
	return &NotificationMongoRepository{collection: database.ForCollection("notifications")}
}

func (r *NotificationMongoRepository) Create(ctx context.Context, notification *entity.Notification) error {
	return r.collection.InsertOne(ctx, notification)
}

func (r *NotificationMongoRepository) Stream(ctx context.Context) <-chan *entity.Notification {
	channel := make(chan *entity.Notification)
	go func() {
		stream, _ := r.collection.FetchStream(ctx)
		for id := range stream {
			notification, err := r.FindById(ctx, id)
			if err != nil || notification == nil {
				continue
			}
			channel <- notification
		}
	}()
	return channel
}

func (r *NotificationMongoRepository) FindById(ctx context.Context, id string) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.collection.FindByID(ctx, id, &notification)
	if err != nil || notification.Id == "" {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationMongoRepository) Delete(ctx context.Context, notification *entity.Notification) {
	_ = r.collection.DeleteByID(ctx, notification.Id)
}

func (r *NotificationMongoRepository) Update(ctx context.Context, notification *entity.Notification) {
	_ = r.collection.UpdateByID(ctx, notification.Id, notification)
}
