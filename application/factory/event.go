package factory

import (
	"bom-pedido-api/application/event"
)

type EventFactory struct {
	EventDispatcher event.Dispatcher
}

func NewEventFactory(EventDispatcher event.Dispatcher) *EventFactory {
	return &EventFactory{
		EventDispatcher: EventDispatcher,
	}
}
