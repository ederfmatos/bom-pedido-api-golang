package event

import (
	"bom-pedido-api/application/event"
	"fmt"
)

type MemoryEventEmitter struct{}

func NewMemoryEventEmitter() *MemoryEventEmitter {
	return &MemoryEventEmitter{}
}

func (emitter *MemoryEventEmitter) Emit(event *event.Event) error {
	fmt.Printf("Emitting event %v", event)
	return nil
}
