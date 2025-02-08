package event

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/pkg/telemetry"
	"context"
)

const _telemetryEventHandlerName = "TELEMETRY"

type TelemetryEventEmitter struct {
	handler event.Handler
}

func NewTelemetryEventEmitter(handler event.Handler) *TelemetryEventEmitter {
	return &TelemetryEventEmitter{handler: handler}
}

func (t TelemetryEventEmitter) Emit(ctx context.Context, event *event.Event) error {
	return telemetry.StartSpanReturningError(ctx, "EventEmitter.Emit", func(ctx context.Context) error {
		telemetry.IncrementMetricWithTags("emit_event", map[string]string{
			"event_name":    string(event.Name),
			"event_emitter": t.handler.Name(),
		})
		return t.handler.Emit(ctx, event)
	}, "event.name", string(event.Name))
}

func (t TelemetryEventEmitter) OnEvent(eventName string, handler event.HandlerFunc) {
	t.handler.OnEvent(eventName, func(ctx context.Context, message *event.MessageEvent) error {
		tags := map[string]string{
			"event_name":    string(message.Event.Name),
			"event_handler": t.handler.Name(),
			"event_topic":   message.Topic,
		}
		telemetry.IncrementMetricWithTags("handle_event", tags)
		err := handler(ctx, message)
		if err != nil {
			telemetry.IncrementErrorMetricWithTags("handle_event", err, tags)
		} else {
			telemetry.IncrementMetricWithTags("handle_event_success", tags)
		}
		return err
	})
}

func (t TelemetryEventEmitter) Close() {
	t.handler.Close()
}

func (t TelemetryEventEmitter) Name() string {
	return _telemetryEventHandlerName
}
