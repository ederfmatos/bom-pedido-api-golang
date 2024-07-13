package event

import (
	"bom-pedido-api/application/event"
	"context"
)

type MemoryEventHandler struct{}

func NewMemoryEventHandler() event.Handler {
	return &MemoryEventHandler{}
}

func (handler *MemoryEventHandler) Emit(context context.Context, event *event.Event) error {
	return nil
}

func (handler *MemoryEventHandler) Close() {
}

func (handler *MemoryEventHandler) Consume(options *event.ConsumerOptions, handlerFunc event.HandlerFunc) {
	//TODO implement me
	panic("implement me")
}
