package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/pkg/telemetry"
	"context"
)

type TelemetryEventEmitter struct {
	handler event.Handler
}

func NewTelemetryEventEmitter(handler event.Handler) *TelemetryEventEmitter {
	return &TelemetryEventEmitter{handler: handler}
}

func (t TelemetryEventEmitter) Emit(ctx context.Context, event *event.Event) error {
	return telemetry.StartSpanReturningError(ctx, "EventEmitter.Emit", func(ctx context.Context) error {
		return t.handler.Emit(ctx, event)
	}, "event.name", event.Name)
}

func (t TelemetryEventEmitter) OnEvent(eventName string, handler event.HandlerFunc) {
	t.handler.OnEvent(eventName, handler)
}

func (t TelemetryEventEmitter) Close() {
	t.handler.Close()
}
