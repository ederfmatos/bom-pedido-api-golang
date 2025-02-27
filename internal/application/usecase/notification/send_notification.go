package notification

import (
	"bom-pedido-api/internal/application/factory"
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/application/lock"
	"bom-pedido-api/internal/application/repository"
	"context"
)

type (
	SendNotificationUseCase struct {
		locker                 lock.Locker
		notificationGateway    gateway.NotificationGateway
		notificationRepository repository.NotificationRepository
	}
)

func NewSendNotification(factory *factory.ApplicationFactory) *SendNotificationUseCase {
	return &SendNotificationUseCase{
		locker:                 factory.Locker,
		notificationRepository: factory.NotificationRepository,
		notificationGateway:    factory.NotificationGateway,
	}
}

func (u *SendNotificationUseCase) Execute(ctx context.Context) {
	// TODO: Add telemetry
	for notification := range u.notificationRepository.Stream(ctx) {
		_ = u.locker.LockFunc(ctx, notification.Id, func() {
			err := u.notificationGateway.Send(ctx, notification)
			if err == nil {
				u.notificationRepository.Delete(ctx, notification)
			} else {
				notification.Fail()
				u.notificationRepository.Update(ctx, notification)
			}
		})
	}
}
