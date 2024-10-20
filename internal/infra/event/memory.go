package event

import (
	"bom-pedido-api/internal/application/event"
	"context"
)

type MemoryEventHandler struct {
}

func NewMemoryEventHandler() *MemoryEventHandler {
	return &MemoryEventHandler{}
}

func (handler *MemoryEventHandler) Emit(context.Context, *event.Event) error {
	return nil
}

func (handler *MemoryEventHandler) Close() {
}

func (handler *MemoryEventHandler) Consume(*event.ConsumerOptions, event.HandlerFunc) {
}
