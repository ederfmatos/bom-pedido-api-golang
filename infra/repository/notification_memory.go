package repository

import (
	"bom-pedido-api/application/repository"
	"bom-pedido-api/domain/entity/notification"
	"context"
)

type NotificationMemoryRepository struct {
	notifications map[string]*notification.Notification
	channel       chan *notification.Notification
}

func NewNotificationMemoryRepository() repository.NotificationRepository {
	return &NotificationMemoryRepository{
		notifications: make(map[string]*notification.Notification),
		channel:       make(chan *notification.Notification),
	}
}

func (repository *NotificationMemoryRepository) Create(_ context.Context, notification *notification.Notification) error {
	repository.notifications[notification.Id] = notification
	repository.channel <- notification
	return nil
}

func (repository *NotificationMemoryRepository) Stream() <-chan *notification.Notification {
	return repository.channel
}

func (repository *NotificationMemoryRepository) Delete(_ context.Context, notification *notification.Notification) {
	delete(repository.notifications, notification.Id)
}

func (repository *NotificationMemoryRepository) Update(_ context.Context, notification *notification.Notification) {
	repository.notifications[notification.Id] = notification
}
