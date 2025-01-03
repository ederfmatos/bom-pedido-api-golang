package factory

import (
	"bom-pedido-api/internal/application/event"
)

type EventFactory struct {
	EventHandler event.Handler
	EventEmitter event.Emitter
}

func NewEventFactory(handler event.Handler) *EventFactory {
	return &EventFactory{
		EventHandler: handler,
		EventEmitter: handler,
	}
}
