package factory

import (
	"bom-pedido-api/application/event"
)

type EventFactory struct {
	EventEmitter event.EventEmitter
	EventHandler event.EventHandler
}

func NewEventFactory(EventEmitter event.EventEmitter, EventHandler event.EventHandler) *EventFactory {
	return &EventFactory{
		EventEmitter: EventEmitter,
		EventHandler: EventHandler,
	}
}
