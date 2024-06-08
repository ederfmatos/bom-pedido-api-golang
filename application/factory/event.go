package factory

import (
	"bom-pedido-api/application/event"
)

type EventFactory struct {
	EventEmitter event.EventEmitter
}

func NewEventFactory(EventEmitter event.EventEmitter) *EventFactory {
	return &EventFactory{EventEmitter: EventEmitter}
}
