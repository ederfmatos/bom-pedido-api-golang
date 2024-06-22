package event

import (
	"bom-pedido-api/application/event"
	"context"
)

type MemoryEventEmitter struct{}

func NewMemoryEventEmitter() *MemoryEventEmitter {
	return &MemoryEventEmitter{}
}

func (emitter *MemoryEventEmitter) Emit(context context.Context, event *event.Event) error {
	return nil
}
