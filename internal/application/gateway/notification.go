package gateway

import (
	"bom-pedido-api/internal/domain/entity"
	"context"
)

type NotificationGateway interface {
	Send(ctx context.Context, notification *entity.Notification) error
}
