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

func (dispatcher *MemoryEventDispatcher) Close() {
}

func (dispatcher *MemoryEventDispatcher) Consume(options *event.ConsumerOptions, handler event.Handler) {
	//TODO implement me
	panic("implement me")
}
