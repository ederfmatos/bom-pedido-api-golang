package repository

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type NotificationMemoryRepository struct {
	notifications map[string]*entity.Notification
	channel       chan *entity.Notification
}

func NewNotificationMemoryRepository() *NotificationMemoryRepository {
	return &NotificationMemoryRepository{
		notifications: make(map[string]*entity.Notification),
		channel:       make(chan *entity.Notification),
	}
}

func (r *NotificationMemoryRepository) Create(_ context.Context, notification *entity.Notification) error {
	r.notifications[notification.Id] = notification
	r.channel <- notification
	return nil
}

func (r *NotificationMemoryRepository) Stream(context.Context) <-chan *entity.Notification {
	return r.channel
}

func (r *NotificationMemoryRepository) Delete(_ context.Context, notification *entity.Notification) {
	delete(r.notifications, notification.Id)
}

func (r *NotificationMemoryRepository) Update(_ context.Context, notification *entity.Notification) {
	r.notifications[notification.Id] = notification
}
