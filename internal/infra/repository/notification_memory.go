package repository

import (
	"bom-pedido-api/internal/domain/entity/notification"
	"context"
)

type NotificationMemoryRepository struct {
	notifications map[string]*notification.Notification
	channel       chan *notification.Notification
}

func NewNotificationMemoryRepository() *NotificationMemoryRepository {
	return &NotificationMemoryRepository{
		notifications: make(map[string]*notification.Notification),
		channel:       make(chan *notification.Notification),
	}
}

func (r *NotificationMemoryRepository) Create(_ context.Context, notification *notification.Notification) error {
	r.notifications[notification.Id] = notification
	r.channel <- notification
	return nil
}

func (r *NotificationMemoryRepository) Stream(context.Context) <-chan *notification.Notification {
	return r.channel
}

func (r *NotificationMemoryRepository) Delete(_ context.Context, notification *notification.Notification) {
	delete(r.notifications, notification.Id)
}

func (r *NotificationMemoryRepository) Update(_ context.Context, notification *notification.Notification) {
	r.notifications[notification.Id] = notification
}
