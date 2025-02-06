package notification

import (
	"bom-pedido-api/internal/application/gateway"
	"bom-pedido-api/internal/domain/entity"
	"bom-pedido-api/pkg/telemetry"
	"context"
)

type TelemetryNotificationGateway struct {
	gateway gateway.NotificationGateway
}

func NewTelemetryNotificationGateway(gateway gateway.NotificationGateway) gateway.NotificationGateway {
	return &TelemetryNotificationGateway{gateway: gateway}
}

func (f *TelemetryNotificationGateway) Send(ctx context.Context, notification *entity.Notification) error {
	return telemetry.StartSpanReturningError(ctx, "NotificationGateway.Send", func(ctx context.Context) error {
		return f.gateway.Send(ctx, notification)
	})
}
