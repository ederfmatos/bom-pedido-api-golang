package event

import (
	"bom-pedido-api/application/event"
	"context"
)

type MemoryEventDispatcher struct{}

func NewMemoryEventDispatcher() event.Dispatcher {
	return &MemoryEventDispatcher{}
}

func (dispatcher *MemoryEventDispatcher) Emit(context context.Context, event *event.Event) error {
	return nil
}

func (dispatcher *MemoryEventDispatcher) Consume(id string, handler event.Handler) {
	//TODO implement me
	panic("implement me")
}

func (dispatcher *MemoryEventDispatcher) Close() {
	//TODO implement me
	panic("implement me")
}
