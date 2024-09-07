package send_notification

import (
	"bom-pedido-api/application/factory"
	"bom-pedido-api/application/gateway"
	"bom-pedido-api/application/lock"
	"bom-pedido-api/application/repository"
	"bom-pedido-api/infra/telemetry"
	"context"
	"time"
)

type (
	UseCase struct {
		locker                 lock.Locker
		notificationGateway    gateway.NotificationGateway
		notificationRepository repository.NotificationRepository
	}
)

func New(factory *factory.ApplicationFactory) *UseCase {
	return &UseCase{
		locker:                 factory.Locker,
		notificationRepository: factory.NotificationRepository,
		notificationGateway:    factory.NotificationGateway,
	}
}

func (u *UseCase) Execute(ctx context.Context) {
	for notification := range u.notificationRepository.Stream(ctx) {
		ctx, span := telemetry.StartSpan(ctx, "SendNotification", "id", notification.Id)
		_ = u.locker.LockFunc(ctx, notification.Id, time.Second*30, func() {
			err := u.notificationGateway.Send(ctx, notification)
			if err == nil {
				u.notificationRepository.Delete(ctx, notification)
			} else {
				notification.Fail()
				u.notificationRepository.Update(ctx, notification)
			}
		})
		span.End()
	}
}
